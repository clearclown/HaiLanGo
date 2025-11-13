# 機能実装: OCR処理

## 要件

### 機能要件

#### OCR処理
1. **対応フォーマット**
   - PDF（メイン）
   - PNG, JPG, HEIC（画像）

2. **OCR精度要件**
   - 複雑なレイアウト対応（表、ルビ、縦書き）
   - 多言語対応（12言語以上）
   - ユーザーによる手動修正機能（任意）

3. **OCR API連携**
   - Google Vision API（優先）
   - Azure Computer Vision（フォールバック）
   - Tesseract（ローカル、オプション）

4. **言語設定**
   - **母国語**: ユーザーのインターフェース言語
   - **学習先言語**: 習得したい言語
   - **参照言語**: 本の仲介言語

#### 処理フロー
1. アップロードされたファイルをOCR APIに送信
2. OCR結果を取得・解析
3. ページ単位でテキストを抽出
4. テキストをDBに保存（PostgreSQL）
5. Redisにキャッシュ（7日間）

#### 非同期処理
- バックグラウンドジョブキュー（Redis/BullMQ）
- 進捗状況のリアルタイム更新（WebSocket）
- エラー時のリトライ機能

### 非機能要件

- **パフォーマンス**:
  - OCR処理時間: 1ページ5-10秒（API依存）
  - バッチ処理: 複数ページの並列処理
  - キャッシュヒット率: 80%以上

- **セキュリティ**:
  - 画像データの暗号化（転送時）
  - ユーザー認証必須
  - OCR結果のアクセス制御

- **拡張性**:
  - 複数OCR APIの切り替え対応
   - 将来的なOCR精度向上への対応

- **エラーハンドリング**:
  - OCR APIエラー時のフォールバック
  - リトライロジック（指数バックオフ）
  - ユーザーへのエラー通知

## 実装指示

### Step 1: テスト設計

以下の順でテストを作成：

1. **ユニットテスト（関数/メソッドレベル）**
   - OCR API呼び出し（モック）
   - テキスト抽出ロジック
   - 言語検出
   - キャッシュ操作

2. **統合テスト（モジュール間）**
   - OCR処理フロー（モックAPI）
   - バッチ処理
   - キャッシュの動作
   - エラーハンドリング

3. **エッジケースのテスト**
   - 低品質画像
   - 複雑なレイアウト
   - 多言語混在ページ
   - APIレート制限
   - ネットワークエラー

テストファイルは `backend/internal/service/ocr/ocr_test.go` および `backend/pkg/ocr/ocr_test.go` に配置。

テストは **実装前に実行してすべて失敗すること** を確認。

#### テスト例（Go）

```go
// backend/internal/service/ocr/ocr_test.go
package ocr

import (
    "testing"
    "context"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestProcessPage(t *testing.T) {
    // OCR処理のテスト（モックAPI使用）
    ctx := context.Background()
    imageData := []byte("test image data")

    result, err := ProcessPage(ctx, imageData, "ru", "ja")
    require.NoError(t, err)
    assert.NotEmpty(t, result.Text)
    assert.Equal(t, "ru", result.DetectedLanguage)
}

func TestDetectLanguage(t *testing.T) {
    // 言語検出のテスト
    text := "Здравствуйте!"
    lang := DetectLanguage(text)
    assert.Equal(t, "ru", lang)
}

func TestExtractTextFromPDF(t *testing.T) {
    // PDFからのテキスト抽出テスト
    // ...
}

func TestCacheOCRResult(t *testing.T) {
    // OCR結果のキャッシュテスト
    // ...
}

func TestRetryOnError(t *testing.T) {
    // エラー時のリトライテスト
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
│   │   ├── ocr/
│   │   │   ├── service.go          # OCRサービス
│   │   │   └── service_test.go
│   │   └── job/
│   │       ├── ocr_job.go          # バックグラウンドジョブ
│   │       └── ocr_job_test.go
├── pkg/
│   ├── ocr/
│   │   ├── google_vision.go        # Google Vision API
│   │   ├── azure_vision.go         # Azure Vision API
│   │   ├── tesseract.go            # Tesseract（ローカル）
│   │   └── ocr_test.go
│   └── cache/
│       ├── redis.go                # Redisキャッシュ
│       └── redis_test.go
```

#### 主要な実装ポイント

1. **OCR API呼び出し**
   ```go
   // pkg/ocr/google_vision.go
   func ProcessWithGoogleVision(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
       // Google Vision API呼び出し
       // エラーハンドリング
       // リトライロジック
   }
   ```

2. **OCRサービス**
   ```go
   // internal/service/ocr/service.go
   func (s *OCRService) ProcessBook(ctx context.Context, bookID string) error {
       // 1. 書籍の全ページ取得
       // 2. 各ページをOCR処理（並列）
       // 3. 結果をDBに保存
       // 4. キャッシュに保存
   }
   ```

3. **バックグラウンドジョブ**
   ```go
   // internal/service/job/ocr_job.go
   func ProcessOCRJob(ctx context.Context, bookID string) error {
       // ジョブキューから取得
       // OCR処理実行
       // 進捗更新（WebSocket）
   }
   ```

### Step 3: リファクタリング

- DRY原則の適用
- パフォーマンス最適化（並列処理）
- コードの可読性向上

### Step 4: ドキュメント

- README.md への機能説明追加
- APIドキュメント
- OCR精度向上のためのガイドライン

#### APIエンドポイント

```
POST /api/v1/books/{book_id}/ocr/process
GET  /api/v1/books/{book_id}/ocr/status
POST /api/v1/books/{book_id}/pages/{page_id}/ocr/retry
```

## 制約事項

- 既存の `backend/internal/models/book.go` は変更しない（新規追加のみ）
- `backend/internal/service/ocr/` 配下のみ編集可能
- 依存関係の追加は `go.mod` のみ

## 完了条件

- [ ] すべてのテストが通る（Vitest + Playwright）
- [ ] BiomeJSエラーがない（フロントエンド）
- [ ] タイプエラーがない
- [ ] GitHub CIが通る
- [ ] ドキュメントが更新されている
- [ ] OCR精度の検証（テスト画像で90%以上）
- [ ] キャッシュの動作確認

## 追加のテスト要件

### セキュリティテスト
- [ ] 画像データの暗号化
- [ ] ユーザー認証の動作確認
- [ ] OCR結果のアクセス制御

### パフォーマンステスト
- [ ] 1ページのOCR処理時間（5-10秒以内）
- [ ] バッチ処理の効率性
- [ ] キャッシュヒット率（80%以上）

### E2Eテスト
- [ ] 書籍アップロードからOCR完了まで
- [ ] 進捗表示の動作確認
- [ ] エラー時のリトライ動作
