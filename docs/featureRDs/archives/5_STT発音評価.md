# 機能実装: STT発音評価

## 要件

### 機能要件

#### STT（Speech-to-Text）機能
1. **音声認識**
   - 主要12言語対応
   - リアルタイム音声認識
   - 録音時間制限：30秒

2. **発音評価**
   - **評価粒度**: 単語レベル
   - **スコア**: 0-100点
   - **フィードバック**:
     - 具体的な改善点提示（英会話教室風）
     - 発音のビジュアル化（波形など）
     - 良かった点の提示

3. **評価項目**
   - **正確性**: 発音の正確さ（0-100点）
   - **流暢性**: 話す速度・リズム（0-100点）
   - **発音**: 個別音素の正確さ（0-100点）

#### API連携
- Google Cloud STT（優先）
- Whisper API（フォールバック）

#### 音声処理
- ノイズ除去（オプション）
- 音量正規化
- サンプリングレート変換（16kHz推奨）

### 非機能要件

- **パフォーマンス**:
  - STT処理時間: 録音終了後3秒以内のフィードバック
  - リアルタイム認識: 100ms以内の遅延
  - 同時処理: 10ユーザーまで

- **セキュリティ**:
  - ユーザー認証必須
  - 音声データの暗号化（転送時）
  - プライバシー保護（音声データは即座に削除）

- **拡張性**:
  - 複数STT APIの切り替え対応
  - 将来的な評価精度向上への対応

- **エラーハンドリング**:
  - STT APIエラー時のフォールバック
  - 低品質音声の検出
  - ユーザーへのエラー通知

## 実装指示

### Step 1: テスト設計

以下の順でテストを作成：

1. **ユニットテスト（関数/メソッドレベル）**
   - STT API呼び出し（モック）
   - 発音評価アルゴリズム
   - スコア計算
   - フィードバック生成

2. **統合テスト（モジュール間）**
   - 音声録音から評価まで
   - スコア計算の正確性
   - フィードバック生成
   - エラーハンドリング

3. **エッジケースのテスト**
   - 無音録音
   - ノイズが多い録音
   - 短すぎる録音（1秒未満）
   - 長すぎる録音（30秒超過）
   - 未対応言語

テストファイルは `backend/internal/service/stt/stt_test.go` および `backend/pkg/stt/stt_test.go` に配置。

テストは **実装前に実行してすべて失敗すること** を確認。

#### テスト例（Go）

```go
// backend/internal/service/stt/stt_test.go
package stt

import (
    "testing"
    "context"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestRecognizeSpeech(t *testing.T) {
    ctx := context.Background()
    audioData := []byte("test audio data")
    lang := "en"
    expectedText := "Hello, world!"

    result, err := RecognizeSpeech(ctx, audioData, lang)
    require.NoError(t, err)
    assert.Equal(t, expectedText, result.Text)
}

func TestEvaluatePronunciation(t *testing.T) {
    // 発音評価のテスト
    expectedText := "Здравствуйте"
    recognizedText := "Здравствуйте"

    score := EvaluatePronunciation(expectedText, recognizedText)
    assert.GreaterOrEqual(t, score, 0)
    assert.LessOrEqual(t, score, 100)
}

func TestCalculateAccuracyScore(t *testing.T) {
    // 正確性スコアの計算テスト
    expected := "Hello"
    recognized := "Hello"

    score := CalculateAccuracyScore(expected, recognized)
    assert.Equal(t, 100, score)

    // 部分的な一致
    recognized = "Hallo"
    score = CalculateAccuracyScore(expected, recognized)
    assert.Greater(t, score, 0)
    assert.Less(t, score, 100)
}

func TestGenerateFeedback(t *testing.T) {
    // フィードバック生成のテスト
    score := 85
    expectedText := "Здравствуйте"
    recognizedText := "Здравствуйте"

    feedback := GenerateFeedback(score, expectedText, recognizedText)
    assert.NotEmpty(t, feedback.Improvements)
    assert.NotEmpty(t, feedback.PositivePoints)
}

func TestNoiseReduction(t *testing.T) {
    // ノイズ除去のテスト
    noisyAudio := []byte("noisy audio data")
    cleanedAudio := ReduceNoise(noisyAudio)
    assert.NotNil(t, cleanedAudio)
}
```

### Step 2: 実装

- テストが通るように実装
- コードコメントは日本語で
- 型安全性を最優先
- エラーハンドリングを必ず含める

#### 実装ファイル構造

```
backend/
├── internal/
│   ├── service/
│   │   ├── stt/
│   │   │   ├── service.go          # STTサービス
│   │   │   ├── evaluation.go       # 発音評価
│   │   │   └── service_test.go
│   │   └── audio/
│   │       ├── processor.go        # 音声処理
│   │       └── processor_test.go
├── pkg/
│   ├── stt/
│   │   ├── google_stt.go           # Google Cloud STT
│   │   ├── whisper.go              # Whisper API
│   │   └── stt_test.go
│   └── audio/
│       ├── processor.go            # 音声処理（ノイズ除去など）
│       └── processor_test.go
```

#### 主要な実装ポイント

1. **STT API呼び出し**
   ```go
   // pkg/stt/google_stt.go
   func RecognizeWithGoogleSTT(ctx context.Context, audioData []byte, lang string) (*STTResult, error) {
       // Google Cloud STT API呼び出し
       // エラーハンドリング
       // リトライロジック
   }
   ```

2. **発音評価**
   ```go
   // internal/service/stt/evaluation.go
   func EvaluatePronunciation(expectedText string, recognizedText string, audioData []byte) (*PronunciationScore, error) {
       // 1. テキスト比較（正確性）
       // 2. 音素レベルの比較（発音）
       // 3. リズム・速度の分析（流暢性）
       // 4. スコア計算（0-100点）
       // 5. フィードバック生成
   }
   ```

3. **フィードバック生成**
   ```go
   // internal/service/stt/evaluation.go
   func GenerateFeedback(score *PronunciationScore, expectedText string, recognizedText string) *Feedback {
       // 良かった点の抽出
       // 改善ポイントの特定
       // 具体的なアドバイス生成
   }
   ```

### Step 3: リファクタリング

- DRY原則の適用
- パフォーマンス最適化（並列処理）
- コードの可読性向上

### Step 4: ドキュメント

- README.md への機能説明追加
- APIドキュメント
- 発音評価アルゴリズムの説明

#### APIエンドポイント

```
POST /api/v1/stt/recognize
POST /api/v1/stt/evaluate-pronunciation
GET  /api/v1/stt/evaluation/{evaluation_id}
```

## 制約事項

- 既存の `backend/internal/models/` は変更しない（新規追加のみ）
- `backend/internal/service/stt/` 配下のみ編集可能
- 依存関係の追加は `go.mod` のみ

## 完了条件

- [ ] すべてのテストが通る
- [ ] lintエラーがない
- [ ] タイプエラーがない
- [ ] ドキュメントが更新されている
- [ ] 発音評価の精度検証（テスト音声で90%以上）
- [ ] フィードバック生成の動作確認

## 追加のテスト要件

### セキュリティテスト
- [ ] ユーザー認証の動作確認
- [ ] 音声データの暗号化
- [ ] プライバシー保護（音声データの即座削除）

### パフォーマンステスト
- [ ] STT処理時間（3秒以内）
- [ ] リアルタイム認識の遅延（100ms以内）
- [ ] 同時処理の動作確認

### E2Eテスト
- [ ] ブラウザでの音声録音
- [ ] 発音評価の動作確認
- [ ] フィードバック表示の動作確認
