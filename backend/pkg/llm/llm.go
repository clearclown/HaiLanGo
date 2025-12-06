package llm

import (
	"context"
)

// LLMClient はLLM APIのインターフェース
// HaiLanGoの価値提案: ドメイン特化学習のためのコンテンツ生成
type LLMClient interface {
	// Generate はプロンプトに基づいてテキストを生成する
	Generate(ctx context.Context, prompt string, options *GenerateOptions) (*GenerateResult, error)

	// Chat はチャット形式のやり取りを行う
	Chat(ctx context.Context, messages []Message, options *GenerateOptions) (*GenerateResult, error)

	// GetName はプロバイダー名を返す
	GetName() string
}

// Message はチャットメッセージ
type Message struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"` // メッセージ内容
}

// GenerateOptions は生成オプション
type GenerateOptions struct {
	MaxTokens   int      `json:"max_tokens"`   // 最大トークン数
	Temperature float64  `json:"temperature"`  // 温度 (0.0-2.0)
	TopP        float64  `json:"top_p"`        // Top-P サンプリング
	StopWords   []string `json:"stop_words"`   // 停止ワード
	SystemPrompt string  `json:"system_prompt"` // システムプロンプト
}

// GenerateResult は生成結果
type GenerateResult struct {
	Content      string `json:"content"`       // 生成テキスト
	FinishReason string `json:"finish_reason"` // 終了理由
	TokensUsed   int    `json:"tokens_used"`   // 使用トークン数
	Provider     string `json:"provider"`      // 使用プロバイダー
	Model        string `json:"model"`         // 使用モデル
}

// LLMProvider はLLMプロバイダーの種類
type LLMProvider string

const (
	// ProviderClaude はAnthropic Claude（推奨: 多言語＆ドメイン特化に最適）
	ProviderClaude LLMProvider = "claude"

	// ProviderOpenAI はOpenAI GPT
	ProviderOpenAI LLMProvider = "openai"

	// ProviderGemini はGoogle Gemini
	ProviderGemini LLMProvider = "gemini"

	// ProviderLocal はローカルLLM（Llama, Mistral等）
	ProviderLocal LLMProvider = "local"
)

// DomainType は学習ドメインの種類
// HaiLanGoの価値提案: 特定分野に特化した学習
type DomainType string

const (
	DomainGeneral    DomainType = "general"    // 一般
	DomainPolitics   DomainType = "politics"   // 政治
	DomainReligion   DomainType = "religion"   // 宗教（クルアーン等）
	DomainSNS        DomainType = "sns"        // SNSスラング
	DomainBusiness   DomainType = "business"   // ビジネス
	DomainMedical    DomainType = "medical"    // 医療
	DomainLegal      DomainType = "legal"      // 法律
	DomainTech       DomainType = "tech"       // テクノロジー
	DomainAcademic   DomainType = "academic"   // 学術
	DomainLiterature DomainType = "literature" // 文学
)

// DomainPromptGenerator はドメイン特化プロンプトを生成
type DomainPromptGenerator struct {
	Domain   DomainType
	Language string
}

// GenerateExplanationPrompt は単語/フレーズ解説プロンプトを生成
func (g *DomainPromptGenerator) GenerateExplanationPrompt(text string) string {
	templates := map[DomainType]string{
		DomainGeneral: `
言語: %s
テキスト: %s

このテキストについて以下を説明してください：
1. 発音のポイント
2. 意味と使用場面
3. 関連する表現
4. 例文`,

		DomainPolitics: `
あなたは政治用語・時事用語の専門家です。

言語: %s
テキスト: %s

この政治用語について：
1. 正式な定義と歴史的背景
2. 現代の政治的文脈での用法
3. 関連する政治用語・概念
4. ニュースでよく見る使用例
5. 発音のポイント（フォーマルな場面向け）`,

		DomainReligion: `
あなたは宗教テキスト・神学用語の専門家です。

言語: %s
テキスト: %s

この宗教用語/テキストについて：
1. 原語の正確な発音と意味
2. 宗教的・神学的な意義
3. 伝統的な解釈と現代の理解
4. 日常用法との違い
5. 類似する表現や関連する経典の引用`,

		DomainSNS: `
あなたはSNSスラング・ネット用語の専門家です。

言語: %s
テキスト: %s

このスラング/ネット用語について：
1. 意味と使用される文脈
2. 由来・語源（ミームなど）
3. フォーマル度（使える場面・使えない場面）
4. 派生語・関連スラング
5. 流行の時期（現在も使われているか）`,

		DomainBusiness: `
あなたはビジネス用語・ビジネスコミュニケーションの専門家です。

言語: %s
テキスト: %s

このビジネス用語について：
1. 正式な定義とビジネス文脈での意味
2. 使用される場面（会議、メール、契約等）
3. 敬語・フォーマル表現との関係
4. 関連するビジネス用語
5. 正確な発音（プレゼンテーション向け）`,

		DomainMedical: `
あなたは医療用語の専門家です。

言語: %s
テキスト: %s

この医療用語について：
1. 医学的定義
2. 一般向け説明
3. 関連する医療用語
4. 正確な発音（医療従事者向け）
5. 患者とのコミュニケーションでの使い方`,

		DomainLegal: `
あなたは法律用語の専門家です。

言語: %s
テキスト: %s

この法律用語について：
1. 法的定義
2. 適用される法域・分野
3. 関連する法律用語
4. 一般向け説明
5. 正確な発音（法廷向け）`,

		DomainTech: `
あなたはIT・テクノロジー用語の専門家です。

言語: %s
テキスト: %s

このIT/テクノロジー用語について：
1. 技術的定義
2. 使用される文脈（プログラミング、ネットワーク等）
3. 関連する技術用語
4. 発音（英語由来の場合の現地語発音）
5. 例：コードやコマンドでの使用例`,

		DomainAcademic: `
あなたは学術用語の専門家です。

言語: %s
テキスト: %s

この学術用語について：
1. 学術的定義
2. 使用される学問分野
3. 語源（ラテン語、ギリシャ語等）
4. 関連する学術用語
5. 論文や講義での使用例`,

		DomainLiterature: `
あなたは文学・古典の専門家です。

言語: %s
テキスト: %s

この文学的表現について：
1. 意味と文学的ニュアンス
2. 出典・引用元（古典、詩、小説等）
3. 歴史的背景
4. 現代での使用
5. 音読のポイント（韻律、リズム）`,
	}

	template, ok := templates[g.Domain]
	if !ok {
		template = templates[DomainGeneral]
	}

	return template
}

// GeneratePronunciationEvalPrompt は発音評価プロンプトを生成
// Whisper + LLMによるマイナー言語対応発音評価
func (g *DomainPromptGenerator) GeneratePronunciationEvalPrompt(expectedText, recognizedText string) string {
	basePrompt := `
あなたは言語教育の専門家で、発音評価を担当しています。

言語: %s
期待されたテキスト: %s
認識されたテキスト: %s

以下の観点で発音を評価してください：

1. 総合スコア (0-100):
2. 正確性スコア (0-100): 単語が正しく発音されているか
3. 流暢性スコア (0-100): 自然な流れで話せているか
4. 改善が必要な単語: 具体的にどの単語をどう改善すべきか

JSON形式で回答してください:
{
  "overall_score": 85,
  "accuracy_score": 88,
  "fluency_score": 82,
  "feedback": "全体的に良い発音です。",
  "improvements": [
    {"word": "xxx", "issue": "...", "suggestion": "..."}
  ]
}
`

	// ドメイン特化の追加指示
	domainNotes := map[DomainType]string{
		DomainReligion: `
追加指示：宗教テキストとして、伝統的な発音規則（タジュウィード等）への準拠も評価してください。`,
		DomainBusiness: `
追加指示：ビジネス場面として、プロフェッショナルな印象を与える発音かどうかも評価してください。`,
		DomainAcademic: `
追加指示：学術発表として、明瞭さと権威ある印象を評価してください。`,
	}

	prompt := basePrompt
	if note, ok := domainNotes[g.Domain]; ok {
		prompt += note
	}

	return prompt
}

// DefaultGenerateOptions はデフォルトの生成オプション
func DefaultGenerateOptions() *GenerateOptions {
	return &GenerateOptions{
		MaxTokens:   1024,
		Temperature: 0.7,
		TopP:        0.9,
	}
}
