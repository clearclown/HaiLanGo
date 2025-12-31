package stt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/clearclown/HaiLanGo/backend/pkg/llm"
)

// WhisperLLMEvaluator はWhisper + LLMによる発音評価器
// HaiLanGoの価値提案: マイナー言語でも発音評価が可能
type WhisperLLMEvaluator struct {
	sttClient STTClient
	llmClient llm.LLMClient
}

// NewWhisperLLMEvaluator は新しいWhisperLLM発音評価器を作成する
func NewWhisperLLMEvaluator(sttClient STTClient, llmClient llm.LLMClient) *WhisperLLMEvaluator {
	return &WhisperLLMEvaluator{
		sttClient: sttClient,
		llmClient: llmClient,
	}
}

// EvaluatePronunciation は発音を評価する
func (e *WhisperLLMEvaluator) EvaluatePronunciation(ctx context.Context, audio io.Reader, expectedText string, language string) (*PronunciationResult, error) {
	// Step 1: Whisperで音声を書き起こす
	transcription, err := e.sttClient.Transcribe(ctx, audio, language)
	if err != nil {
		return nil, fmt.Errorf("failed to transcribe audio: %w", err)
	}

	// Step 2: LLMで発音評価を行う
	evaluation, err := e.evaluateWithLLM(ctx, transcription.Text, expectedText, language)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate pronunciation: %w", err)
	}

	evaluation.RecognizedText = transcription.Text
	evaluation.ExpectedText = expectedText
	evaluation.EvaluationMethod = "whisper_llm"

	return evaluation, nil
}

// llmEvaluationResult はLLMからの評価結果をパースするための構造体
type llmEvaluationResult struct {
	OverallScore    float64 `json:"overall_score"`
	AccuracyScore   float64 `json:"accuracy_score"`
	FluencyScore    float64 `json:"fluency_score"`
	ProsodyScore    float64 `json:"prosody_score"`
	Feedback        string  `json:"feedback"`
	ImprovementTips []string `json:"improvement_tips"`
	WordAnalysis    []struct {
		Word       string  `json:"word"`
		Score      float64 `json:"score"`
		ErrorType  string  `json:"error_type,omitempty"`
		Suggestion string  `json:"suggestion,omitempty"`
	} `json:"word_analysis"`
}

// evaluateWithLLM はLLMを使用して発音評価を行う
func (e *WhisperLLMEvaluator) evaluateWithLLM(ctx context.Context, recognizedText, expectedText, language string) (*PronunciationResult, error) {
	prompt := e.buildEvaluationPrompt(recognizedText, expectedText, language)

	options := &llm.GenerateOptions{
		MaxTokens:   1024,
		Temperature: 0.3, // より一貫した評価のため低めの温度
		SystemPrompt: `You are an expert pronunciation evaluator for language learning.
You analyze speech recognition results and provide detailed pronunciation feedback.
You always respond with valid JSON following the exact structure requested.
Be encouraging but honest in your feedback.`,
	}

	result, err := e.llmClient.Generate(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("LLM evaluation failed: %w", err)
	}

	// JSONをパース
	evaluation, err := e.parseEvaluationResult(result.Content)
	if err != nil {
		// パース失敗時はデフォルトの評価を返す
		return e.createDefaultEvaluation(recognizedText, expectedText), nil
	}

	return evaluation, nil
}

// buildEvaluationPrompt は評価用のプロンプトを構築する
func (e *WhisperLLMEvaluator) buildEvaluationPrompt(recognizedText, expectedText, language string) string {
	return fmt.Sprintf(`Evaluate the pronunciation accuracy by comparing what was spoken (recognized) with what was expected.

Language: %s
Expected text: "%s"
Recognized text: "%s"

Please analyze and provide a JSON response with the following structure:
{
  "overall_score": <0-100>,
  "accuracy_score": <0-100>,
  "fluency_score": <0-100>,
  "prosody_score": <0-100>,
  "feedback": "<brief encouraging feedback in the target language or English>",
  "improvement_tips": ["<specific tip 1>", "<specific tip 2>"],
  "word_analysis": [
    {
      "word": "<word from expected text>",
      "score": <0-100>,
      "error_type": "<omission|mispronunciation|insertion|none>",
      "suggestion": "<how to improve pronunciation of this word>"
    }
  ]
}

Scoring guidelines:
- 90-100: Near-native pronunciation
- 75-89: Good pronunciation with minor issues
- 60-74: Understandable but needs improvement
- 40-59: Significant pronunciation issues
- 0-39: Major pronunciation problems

Focus on:
1. Word accuracy (were all words pronounced?)
2. Phoneme accuracy (were sounds correct?)
3. Fluency (smooth delivery?)
4. Prosody (rhythm and intonation?)

Respond ONLY with valid JSON, no additional text.`, language, expectedText, recognizedText)
}

// parseEvaluationResult はLLMの応答からPronunciationResultを構築する
func (e *WhisperLLMEvaluator) parseEvaluationResult(content string) (*PronunciationResult, error) {
	// JSONブロックを抽出（マークダウンコードブロック対応）
	jsonContent := content
	if strings.Contains(content, "```json") {
		start := strings.Index(content, "```json") + 7
		end := strings.LastIndex(content, "```")
		if end > start {
			jsonContent = strings.TrimSpace(content[start:end])
		}
	} else if strings.Contains(content, "```") {
		start := strings.Index(content, "```") + 3
		end := strings.LastIndex(content, "```")
		if end > start {
			jsonContent = strings.TrimSpace(content[start:end])
		}
	}

	var eval llmEvaluationResult
	if err := json.Unmarshal([]byte(jsonContent), &eval); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	// WordScoresに変換
	wordScores := make([]WordScore, len(eval.WordAnalysis))
	for i, wa := range eval.WordAnalysis {
		wordScores[i] = WordScore{
			Word:       wa.Word,
			Score:      wa.Score,
			ErrorType:  wa.ErrorType,
			Suggestion: wa.Suggestion,
		}
	}

	return &PronunciationResult{
		OverallScore:    eval.OverallScore,
		AccuracyScore:   eval.AccuracyScore,
		FluencyScore:    eval.FluencyScore,
		ProsodyScore:    eval.ProsodyScore,
		Feedback:        eval.Feedback,
		ImprovementTips: eval.ImprovementTips,
		WordScores:      wordScores,
	}, nil
}

// createDefaultEvaluation はパース失敗時のデフォルト評価を作成
func (e *WhisperLLMEvaluator) createDefaultEvaluation(recognizedText, expectedText string) *PronunciationResult {
	// シンプルな類似度計算
	similarity := e.calculateSimilarity(recognizedText, expectedText)
	score := similarity * 100

	feedback := "Your pronunciation was evaluated."
	if score >= 80 {
		feedback = "Excellent pronunciation! Keep up the great work."
	} else if score >= 60 {
		feedback = "Good effort! Some words need a bit more practice."
	} else if score >= 40 {
		feedback = "You're making progress. Focus on clearer pronunciation."
	} else {
		feedback = "Keep practicing! Listen carefully and try again."
	}

	return &PronunciationResult{
		OverallScore:   score,
		AccuracyScore:  score,
		FluencyScore:   score * 0.9, // やや低め
		ProsodyScore:   score * 0.85,
		Feedback:       feedback,
		ImprovementTips: []string{
			"Listen to the native pronunciation carefully",
			"Practice each word slowly before speaking the full phrase",
			"Record yourself and compare with the original",
		},
	}
}

// calculateSimilarity は2つのテキストの類似度を計算する（0.0-1.0）
func (e *WhisperLLMEvaluator) calculateSimilarity(s1, s2 string) float64 {
	// 正規化
	s1 = strings.ToLower(strings.TrimSpace(s1))
	s2 = strings.ToLower(strings.TrimSpace(s2))

	if s1 == s2 {
		return 1.0
	}

	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// 単語ベースのジャカード類似度
	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)

	set1 := make(map[string]bool)
	for _, w := range words1 {
		set1[w] = true
	}

	set2 := make(map[string]bool)
	for _, w := range words2 {
		set2[w] = true
	}

	// 共通要素
	intersection := 0
	for w := range set1 {
		if set2[w] {
			intersection++
		}
	}

	// 和集合のサイズ
	union := len(set1) + len(set2) - intersection

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// Ensure WhisperLLMEvaluator implements PronunciationEvaluator interface
var _ PronunciationEvaluator = (*WhisperLLMEvaluator)(nil)
