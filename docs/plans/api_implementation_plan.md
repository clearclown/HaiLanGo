# API実装計画書 - TDD アプローチ

## 概要

TTS、STT、LLMの各APIをTDD（テスト駆動開発）で実装する。
既存のOCRパッケージの実装パターンを踏襲する。

## 実装パターン（OCRから継承）

各パッケージは以下の構造を持つ：

```
pkg/{api}/
├── {api}.go          # インターフェース定義 ✅ 完了
├── factory.go        # ファクトリー（環境変数による切り替え）
├── mock.go           # モック実装
├── {provider}.go     # 各プロバイダー実装
└── {api}_test.go     # テスト
```

## TDD サイクル

各機能について以下のサイクルを繰り返す：
1. **Red**: 失敗するテストを書く
2. **Green**: テストを通す最小限のコードを書く
3. **Refactor**: コードをリファクタリング

---

## Phase 1: TTS (Text-to-Speech)

### 1.1 ファイル構造

```
backend/pkg/tts/
├── tts.go           # ✅ 完了（インターフェース定義）
├── factory.go       # 🔨 実装予定
├── mock.go          # 🔨 実装予定
├── azure.go         # 🔨 実装予定（推奨: 140言語）
├── google.go        # 🔨 実装予定
├── edge.go          # 🔨 実装予定（無料オプション）
└── tts_test.go      # 🔨 実装予定
```

### 1.2 TDD ステップ

#### Step 1: テストファイル作成（Red）
```go
// tts_test.go
func TestNewTTSClient_WithMocks(t *testing.T)
func TestNewTTSClient_WithoutAPIKey(t *testing.T)
func TestMockTTSClient_Synthesize(t *testing.T)
func TestMockTTSClient_GetVoices(t *testing.T)
func TestMockTTSClient_GetSupportedLanguages(t *testing.T)
func TestGetRecommendedProvider(t *testing.T)
```

#### Step 2: モック実装（Green）
```go
// mock.go
type MockTTSClient struct {
    dataDir string
}
func NewMockTTSClient() *MockTTSClient
func (m *MockTTSClient) Synthesize(...) (io.ReadCloser, error)
func (m *MockTTSClient) GetVoices(...) ([]Voice, error)
func (m *MockTTSClient) GetSupportedLanguages(...) ([]string, error)
func (m *MockTTSClient) GetName() string
```

#### Step 3: ファクトリー実装
```go
// factory.go
func NewTTSClient() (TTSClient, error)
// 環境変数: TTS_PROVIDER, AZURE_SPEECH_KEY, GOOGLE_TTS_KEY
```

#### Step 4: Azure TTS 実装（優先）
```go
// azure.go
type AzureTTSClient struct
func NewAzureTTSClient(subscriptionKey, region string) *AzureTTSClient
// REST API使用: cognitiveservices.azure.com
```

### 1.3 優先プロバイダー

| 優先度 | プロバイダー | 言語数 | 理由 |
|-------|-------------|-------|------|
| 1 | Azure | 140 | マイナー言語最大カバレッジ |
| 2 | edge-tts | 50+ | 無料、オフライン可 |
| 3 | Google | 50+ | 安定性 |

---

## Phase 2: STT (Speech-to-Text)

### 2.1 ファイル構造

```
backend/pkg/stt/
├── stt.go           # ✅ 完了（インターフェース定義）
├── factory.go       # 🔨 実装予定
├── mock.go          # 🔨 実装予定
├── whisper.go       # 🔨 実装予定（推奨: 99言語）
├── azure.go         # 🔨 実装予定
├── evaluator.go     # 🔨 実装予定（発音評価: Whisper+LLM）
└── stt_test.go      # 🔨 実装予定
```

### 2.2 TDD ステップ

#### Step 1: テストファイル作成（Red）
```go
// stt_test.go
func TestNewSTTClient_WithMocks(t *testing.T)
func TestNewSTTClient_WithoutAPIKey(t *testing.T)
func TestMockSTTClient_Transcribe(t *testing.T)
func TestMockSTTClient_GetSupportedLanguages(t *testing.T)
func TestGetRecommendedSTTProvider(t *testing.T)
func TestGetRecommendedEvalMethod(t *testing.T)
// 発音評価テスト
func TestWhisperLLMEvaluator_EvaluatePronunciation(t *testing.T)
```

#### Step 2: モック実装（Green）
```go
// mock.go
type MockSTTClient struct
func (m *MockSTTClient) Transcribe(...) (*TranscriptionResult, error)
func (m *MockSTTClient) GetSupportedLanguages(...) ([]string, error)

type MockPronunciationEvaluator struct
func (m *MockPronunciationEvaluator) EvaluatePronunciation(...) (*PronunciationResult, error)
```

#### Step 3: ファクトリー実装
```go
// factory.go
func NewSTTClient() (STTClient, error)
func NewPronunciationEvaluator(llmClient llm.LLMClient) (PronunciationEvaluator, error)
```

#### Step 4: Whisper 実装（優先）
```go
// whisper.go
type WhisperClient struct
func NewWhisperClient(apiKey string) *WhisperClient
// OpenAI Whisper API使用
```

#### Step 5: 発音評価実装
```go
// evaluator.go
type WhisperLLMEvaluator struct {
    sttClient STTClient
    llmClient llm.LLMClient
}
func (e *WhisperLLMEvaluator) EvaluatePronunciation(...) (*PronunciationResult, error)
// Whisperで認識 → LLMで評価（マイナー言語対応）
```

### 2.3 優先プロバイダー

| 優先度 | プロバイダー | 言語数 | 理由 |
|-------|-------------|-------|------|
| 1 | Whisper | 99 | マイナー言語最強 |
| 2 | Azure | 100+ | リアルタイム発音評価 |
| 3 | whisper.cpp | 99 | オフライン |

---

## Phase 3: LLM (Large Language Model)

### 3.1 ファイル構造

```
backend/pkg/llm/
├── llm.go           # ✅ 完了（インターフェース定義）
├── factory.go       # 🔨 実装予定
├── mock.go          # 🔨 実装予定
├── claude.go        # 🔨 実装予定（推奨）
├── openai.go        # 🔨 実装予定
└── llm_test.go      # 🔨 実装予定
```

### 3.2 TDD ステップ

#### Step 1: テストファイル作成（Red）
```go
// llm_test.go
func TestNewLLMClient_WithMocks(t *testing.T)
func TestNewLLMClient_WithoutAPIKey(t *testing.T)
func TestMockLLMClient_Generate(t *testing.T)
func TestMockLLMClient_Chat(t *testing.T)
func TestDefaultGenerateOptions(t *testing.T)
```

#### Step 2: モック実装（Green）
```go
// mock.go
type MockLLMClient struct
func (m *MockLLMClient) Generate(...) (*GenerateResult, error)
func (m *MockLLMClient) Chat(...) (*GenerateResult, error)
func (m *MockLLMClient) GetName() string
```

#### Step 3: ファクトリー実装
```go
// factory.go
func NewLLMClient() (LLMClient, error)
// 環境変数: LLM_PROVIDER, ANTHROPIC_API_KEY, OPENAI_API_KEY
```

#### Step 4: Claude 実装（優先）
```go
// claude.go
type ClaudeClient struct
func NewClaudeClient(apiKey string) *ClaudeClient
// Anthropic API使用
```

### 3.3 優先プロバイダー

| 優先度 | プロバイダー | 理由 |
|-------|-------------|------|
| 1 | Claude 3.5 | 多言語対応、長文理解 |
| 2 | GPT-4 | 安定性、幅広い対応 |
| 3 | Gemini | コスト効率 |

---

## Phase 4: 統合テスト

### 4.1 発音評価パイプラインテスト

```go
// backend/internal/service/pronunciation/service_test.go
func TestPronunciationService_Evaluate(t *testing.T)
// Whisper STT → LLM 評価の統合テスト
```

### 4.2 言語別テスト

```go
// マイナー言語での動作確認
func TestMinorLanguageSupport(t *testing.T) {
    languages := []string{"ku", "am", "ti", "bo"} // クルド語、アムハラ語等
    // 各言語でTTS/STTが動作することを確認
}
```

---

## 実装順序

```
Day 1-2: TTS
├── [ ] tts_test.go (テスト作成)
├── [ ] mock.go (モック実装)
├── [ ] factory.go (ファクトリー)
└── [ ] azure.go (Azure実装)

Day 3-4: STT
├── [ ] stt_test.go (テスト作成)
├── [ ] mock.go (モック実装)
├── [ ] factory.go (ファクトリー)
├── [ ] whisper.go (Whisper実装)
└── [ ] evaluator.go (発音評価)

Day 5: LLM
├── [ ] llm_test.go (テスト作成)
├── [ ] mock.go (モック実装)
├── [ ] factory.go (ファクトリー)
└── [ ] claude.go (Claude実装)

Day 6: 統合
├── [ ] 統合テスト
├── [ ] ドキュメント更新
└── [ ] CI/CD設定
```

---

## 環境変数まとめ

```bash
# 共通
USE_MOCK_APIS=true        # モック使用
TEST_USE_MOCKS=true       # テスト時モック

# TTS
TTS_PROVIDER=azure        # azure, google, edge
AZURE_SPEECH_KEY=xxx
AZURE_SPEECH_REGION=eastasia
GOOGLE_TTS_KEY=xxx

# STT
STT_PROVIDER=whisper      # whisper, azure, google
OPENAI_API_KEY=xxx        # Whisper用
AZURE_SPEECH_KEY=xxx      # Azure STT用

# LLM
LLM_PROVIDER=claude       # claude, openai, gemini
ANTHROPIC_API_KEY=xxx
OPENAI_API_KEY=xxx
```

---

## テスト実行コマンド

```bash
# 全テスト（モック使用）
go test ./backend/pkg/tts/... ./backend/pkg/stt/... ./backend/pkg/llm/...

# カバレッジ付き
go test -cover ./backend/pkg/...

# 特定パッケージ
go test -v ./backend/pkg/tts/...

# 実API使用（APIキー必要）
USE_MOCK_APIS=false go test ./backend/pkg/...
```

---

## 成功基準

1. **テストカバレッジ**: 80%以上
2. **モック動作**: APIキーなしで全テストパス
3. **マイナー言語**: クルド語(ku)、アムハラ語(am)でTTS/STT動作
4. **発音評価**: Whisper+LLMパイプラインが動作
5. **CI/CD**: GitHub Actionsでテスト自動実行
