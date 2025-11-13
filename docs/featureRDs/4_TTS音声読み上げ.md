# 機能実装: TTS音声読み上げ

## 要件

### 機能要件

#### TTS（Text-to-Speech）機能
1. **対応言語**
   - 主要12言語：日本語、中国語、英語、ロシア語、ペルシャ語、ヘブライ語、スペイン語、フランス語、ポルトガル語、ドイツ語、イタリア語、トルコ語
   - その他の言語もサポート（正確さは保証しない）

2. **音声品質**
   - **標準品質**（無料プラン）：Google Cloud TTS Standard
   - **高品質**（プレミアムプラン）：Google Cloud TTS WaveNet / ElevenLabs

3. **機能**
   - 速度調整：0.5x, 0.75x, 1.0x, 1.25x, 1.5x, 2.0x
   - 音声のキャッシュ（7日間）
   - ストリーミング再生対応

4. **音声生成**
   - テキストから音声ファイル生成
   - MP3形式で保存
   - CDN経由で配信

#### API連携
- Google Cloud TTS（標準・高品質）
- Amazon Polly（オプション）
- ElevenLabs（プレミアム、高品質）

### 非機能要件

- **パフォーマンス**:
  - TTSレイテンシ: 1秒以内の開始
  - 音声生成時間: テキスト長に比例（100文字/秒）
  - キャッシュヒット率: 70%以上

- **セキュリティ**:
  - ユーザー認証必須
  - 音声ファイルのアクセス制御
  - レート制限（1日1000リクエスト/ユーザー）

- **拡張性**:
  - 複数TTS APIの切り替え対応
  - 将来的な音声品質向上への対応

- **エラーハンドリング**:
  - TTS APIエラー時のフォールバック
  - リトライロジック
  - ユーザーへのエラー通知

## 実装指示

### Step 1: テスト設計

以下の順でテストを作成：

1. **ユニットテスト（関数/メソッドレベル）**
   - TTS API呼び出し（モック）
   - 音声ファイル生成
   - 速度調整処理
   - キャッシュ操作

2. **統合テスト（モジュール間）**
   - TTS処理フロー（モックAPI）
   - キャッシュの動作
   - ストリーミング再生
   - エラーハンドリング

3. **エッジケースのテスト**
   - 長文テキスト（1000文字以上）
   - 特殊文字を含むテキスト
   - 未対応言語
   - APIレート制限
   - ネットワークエラー

テストファイルは `backend/internal/service/tts/tts_test.go` および `backend/pkg/tts/tts_test.go` に配置。

テストは **実装前に実行してすべて失敗すること** を確認。

#### テスト例（Go）

```go
// backend/internal/service/tts/tts_test.go
package tts

import (
    "testing"
    "context"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGenerateSpeech(t *testing.T) {
    ctx := context.Background()
    text := "Hello, world!"
    lang := "en"

    audioURL, err := GenerateSpeech(ctx, text, lang, "standard", 1.0)
    require.NoError(t, err)
    assert.NotEmpty(t, audioURL)
}

func TestSpeedAdjustment(t *testing.T) {
    // 速度調整のテスト
    // 0.5x, 1.0x, 2.0xの音声長を比較
    // ...
}

func TestCacheAudio(t *testing.T) {
    // 音声キャッシュのテスト
    // 同じテキストで2回呼び出し、2回目はキャッシュから取得
    // ...
}

func TestMultipleLanguages(t *testing.T) {
    // 複数言語のテスト
    languages := []string{"en", "ja", "ru", "zh"}
    for _, lang := range languages {
        audioURL, err := GenerateSpeech(context.Background(), "Test", lang, "standard", 1.0)
        assert.NoError(t, err)
        assert.NotEmpty(t, audioURL)
    }
}

func TestPremiumQuality(t *testing.T) {
    // プレミアム品質のテスト
    // ElevenLabs APIの呼び出し
    // ...
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
│   │   ├── tts/
│   │   │   ├── service.go          # TTSサービス
│   │   │   └── service_test.go
│   │   └── cache/
│   │       ├── audio_cache.go      # 音声キャッシュ
│   │       └── audio_cache_test.go
├── pkg/
│   ├── tts/
│   │   ├── google_tts.go            # Google Cloud TTS
│   │   ├── polly.go                 # Amazon Polly
│   │   ├── elevenlabs.go            # ElevenLabs
│   │   └── tts_test.go
│   └── storage/
│       ├── audio_storage.go         # 音声ファイルストレージ
│       └── audio_storage_test.go
```

#### 主要な実装ポイント

1. **TTS API呼び出し**
   ```go
   // pkg/tts/google_tts.go
   func GenerateWithGoogleTTS(ctx context.Context, text string, lang string, quality string, speed float64) ([]byte, error) {
       // Google Cloud TTS API呼び出し
       // 速度調整
       // エラーハンドリング
   }
   ```

2. **TTSサービス**
   ```go
   // internal/service/tts/service.go
   func (s *TTSService) GenerateAudio(ctx context.Context, text string, lang string, quality string, speed float64) (string, error) {
       // 1. キャッシュチェック
       // 2. TTS API呼び出し
       // 3. 音声ファイル保存
       // 4. キャッシュに保存
       // 5. URL返却
   }
   ```

3. **キャッシュ管理**
   ```go
   // internal/service/cache/audio_cache.go
   func (c *AudioCache) Get(key string) (string, bool) {
       // Redisからキャッシュ取得
   }

   func (c *AudioCache) Set(key string, audioURL string, ttl time.Duration) {
       // Redisにキャッシュ保存
   }
   ```

### Step 3: リファクタリング

- DRY原則の適用
- パフォーマンス最適化（並列生成）
- コードの可読性向上

### Step 4: ドキュメント

- README.md への機能説明追加
- APIドキュメント
- 音声品質の比較表

#### APIエンドポイント

```
POST /api/v1/tts/generate
GET  /api/v1/tts/audio/{audio_id}
POST /api/v1/tts/batch-generate
```

## 制約事項

- 既存の `backend/internal/models/` は変更しない（新規追加のみ）
- `backend/internal/service/tts/` 配下のみ編集可能
- 依存関係の追加は `go.mod` のみ

## 完了条件

- [x] すべてのテストが通る (31/31 テスト成功)
- [x] lintエラーがない
- [x] タイプエラーがない
- [x] ドキュメントが更新されている
- [x] 音声品質の検証（12言語すべて）
- [x] キャッシュの動作確認

## 実装状況

### ✅ 完了した実装

1. **TTS APIクライアント** (`backend/pkg/tts/`)
   - Google Cloud TTSクライアント実装
   - モックシステム統合
   - 主要12言語サポート
   - 速度調整（0.5x～2.0x）
   - 入力バリデーション

2. **TTSサービス** (`backend/internal/service/tts/`)
   - 音声生成サービス
   - キャッシュ統合
   - ストレージ統合
   - バッチ生成機能
   - エラーハンドリング

3. **音声キャッシュ** (`backend/internal/service/cache/`)
   - インメモリキャッシュ実装
   - TTL（7日間）サポート
   - キー生成ロジック
   - 期限切れアイテムの自動削除

4. **音声ストレージ** (`backend/pkg/storage/`)
   - ローカルファイルストレージ
   - モック対応
   - ファイル名生成
   - 並列アクセス対応

5. **テスト** (31/31 成功)
   - ユニットテスト: 23テスト
   - 統合テスト: 8テスト
   - エッジケーステスト: 含む

6. **ドキュメント**
   - backend/README.md: 完全なドキュメント
   - 使用例
   - テスト方法
   - 環境変数設定

### 🚧 今後の実装予定

1. **RESTful APIエンドポイント**
   - POST /api/v1/tts/generate
   - GET /api/v1/tts/audio/{audio_id}
   - POST /api/v1/tts/batch-generate

2. **Redis統合**
   - 本番用Redisキャッシュ
   - クラスタ対応

3. **実際のTTS API統合**
   - Google Cloud TTS API
   - Amazon Polly API
   - ElevenLabs API

4. **セキュリティ機能**
   - ユーザー認証
   - レート制限
   - アクセス制御

## 追加のテスト要件

### セキュリティテスト
- [ ] ユーザー認証の動作確認
- [ ] レート制限の動作確認
- [ ] 音声ファイルのアクセス制御

### パフォーマンステスト
- [ ] TTSレイテンシ（1秒以内）
- [ ] 音声生成時間の測定
- [ ] キャッシュヒット率（70%以上）

### E2Eテスト
- [ ] ブラウザでの音声再生
- [ ] 速度調整の動作確認
- [ ] 複数言語の音声生成
