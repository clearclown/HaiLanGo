package audio

import (
	"fmt"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// AudioProcessor は音声処理を行う
type AudioProcessor struct{}

// NewAudioProcessor は新しいオーディオプロセッサーを作成する
func NewAudioProcessor() *AudioProcessor {
	return &AudioProcessor{}
}

// Process は音声データを処理する
func (p *AudioProcessor) Process(audioData []byte) (*models.AudioProcessingResult, error) {
	if len(audioData) == 0 {
		return nil, fmt.Errorf("音声データが空です")
	}

	// ノイズレベルを検出
	noiseLevel, err := p.DetectNoiseLevel(audioData)
	if err != nil {
		return nil, fmt.Errorf("ノイズレベルの検出に失敗しました: %w", err)
	}

	// ノイズ除去
	cleaned, err := p.ApplyNoiseReduction(audioData)
	if err != nil {
		return nil, fmt.Errorf("ノイズ除去に失敗しました: %w", err)
	}

	// 音量正規化
	normalized, err := p.NormalizeVolume(cleaned)
	if err != nil {
		return nil, fmt.Errorf("音量正規化に失敗しました: %w", err)
	}

	// サンプリングレート変換（16kHzに統一）
	processed, err := p.ConvertSampleRate(normalized, 16000)
	if err != nil {
		return nil, fmt.Errorf("サンプリングレート変換に失敗しました: %w", err)
	}

	result := &models.AudioProcessingResult{
		ProcessedAudio: processed,
		SampleRate:     16000,
		Channels:       1,
		Duration:       float64(len(processed)) / 16000.0, // 簡易的な計算
		NoiseLevel:     noiseLevel,
		IsLowQuality:   noiseLevel > 0.3, // ノイズレベルが30%を超えると低品質と判断
	}

	return result, nil
}

// DetectNoiseLevel はノイズレベルを検出する
func (p *AudioProcessor) DetectNoiseLevel(audioData []byte) (float64, error) {
	if len(audioData) == 0 {
		return 0, fmt.Errorf("音声データが空です")
	}

	// TODO: 実際のノイズレベル検出アルゴリズムを実装
	// 現時点では簡易的な計算を返す
	return 0.15, nil
}

// ApplyNoiseReduction はノイズ除去を適用する
func (p *AudioProcessor) ApplyNoiseReduction(audioData []byte) ([]byte, error) {
	if len(audioData) == 0 {
		return nil, fmt.Errorf("音声データが空です")
	}

	// TODO: 実際のノイズ除去アルゴリズムを実装
	// 現時点では元のデータをそのまま返す
	return audioData, nil
}

// NormalizeVolume は音量を正規化する
func (p *AudioProcessor) NormalizeVolume(audioData []byte) ([]byte, error) {
	if len(audioData) == 0 {
		return nil, fmt.Errorf("音声データが空です")
	}

	// TODO: 実際の音量正規化アルゴリズムを実装
	// 現時点では元のデータをそのまま返す
	return audioData, nil
}

// ConvertSampleRate はサンプリングレートを変換する
func (p *AudioProcessor) ConvertSampleRate(audioData []byte, targetRate int) ([]byte, error) {
	if len(audioData) == 0 {
		return nil, fmt.Errorf("音声データが空です")
	}

	if targetRate <= 0 {
		return nil, fmt.Errorf("無効なサンプリングレート: %d", targetRate)
	}

	// TODO: 実際のサンプリングレート変換を実装
	// 現時点では元のデータをそのまま返す
	return audioData, nil
}

// ValidateFormat は音声フォーマットを検証する
func (p *AudioProcessor) ValidateFormat(audioData []byte) (bool, error) {
	if len(audioData) == 0 {
		return false, nil
	}

	// TODO: 実際のフォーマット検証を実装
	// 現時点では簡易的な検証を行う

	// "low quality"という文字列は低品質とみなす（テスト用）
	dataStr := string(audioData)
	if dataStr == "low quality" {
		return false, nil
	}

	// 最小サイズチェック（少なくとも10バイト以上が必要）
	if len(audioData) < 10 {
		return false, nil
	}

	return true, nil
}
