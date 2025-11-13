# HaiLanGo Backend

HaiLanGoプロジェクトのバックエンド実装（Go）

## 概要

このディレクトリには、HaiLanGoプロジェクトのバックエンドコードが含まれています。

## ディレクトリ構造

```
backend/
├── cmd/                       # エントリーポイント
│   └── server/
│       └── main.go
├── internal/                  # 内部パッケージ（他プロジェクトから使用不可）
│   ├── api/                   # APIハンドラー
│   │   ├── handler/           # HTTPハンドラー
│   │   ├── middleware/        # ミドルウェア
│   │   └── router/            # ルーティング
│   ├── service/               # ビジネスロジック
│   │   ├── ocr/               # OCR処理サービス
│   │   └── job/               # バックグラウンドジョブ
│   ├── repository/            # データアクセス層
│   └── models/                # データモデル
└── pkg/                       # 外部パッケージ（再利用可能）
    ├── ocr/                   # OCRクライアント
    └── cache/                 # キャッシュ
```

## 技術スタック

- **言語**: Go 1.21+
- **データベース**: PostgreSQL 15+
- **キャッシュ**: Redis 7+
- **テスト**: testify
- **外部API**: Google Vision API, Azure Computer Vision

## セットアップ

### 1. 依存関係のインストール

```bash
go mod download
```

### 2. 環境変数の設定

`.env`ファイルを作成して、以下の環境変数を設定します：

```bash
# APIキーなしでも開発可能
USE_MOCK_APIS=true

# OCRプロバイダー（google_vision, azure_vision, tesseract）
OCR_PROVIDER=google_vision

# Google Cloud Vision API（オプション）
GOOGLE_CLOUD_VISION_API_KEY=your_api_key_here

# Azure Computer Vision（オプション）
AZURE_COMPUTER_VISION_ENDPOINT=https://your-resource.cognitiveservices.azure.com/
AZURE_COMPUTER_VISION_API_KEY=your_api_key_here

# データベース
DATABASE_URL=postgresql://HaiLanGo:password@localhost:5432/HaiLanGo_dev

# Redis
REDIS_URL=redis://localhost:6379
```

### 3. テストの実行

```bash
# すべてのテストを実行
go test ./...

# カバレッジ付きでテスト
go test ./... -cover

# 詳細なテスト結果
go test ./... -v
```

## OCR処理機能

### 概要

OCR（Optical Character Recognition）処理機能は、書籍のページ画像からテキストを抽出します。

### 主な機能

- **多言語対応**: 12言語以上に対応（日本語、中国語、ロシア語、英語など）
- **複数プロバイダー対応**: Google Vision API、Azure Computer Vision、Tesseract
- **キャッシュ機能**: Redisを使用した7日間のキャッシュ
- **モックシステム**: APIキーなしでも開発・テスト可能

### 使用方法

#### OCRクライアントの作成

```go
import (
    "github.com/clearclown/HaiLanGo/backend/pkg/ocr"
    "github.com/clearclown/HaiLanGo/backend/pkg/cache"
    ocrService "github.com/clearclown/HaiLanGo/backend/internal/service/ocr"
)

// OCRクライアントを作成
ocrClient, err := ocr.NewOCRClient()
if err != nil {
    log.Fatal(err)
}

// キャッシュを作成
cache := cache.NewMockCache() // または実際のRedisクライアント

// OCRサービスを作成
service := ocrService.NewOCRService(ocrClient, cache)
```

#### ページのOCR処理

```go
import (
    "context"
    "github.com/google/uuid"
)

ctx := context.Background()
pageID := uuid.New()
imageData := []byte("...") // 画像データ
languages := []string{"ru", "en"} // 検出する言語

// OCR処理を実行
page, err := service.ProcessPage(ctx, pageID, imageData, languages)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("OCR Text: %s\n", page.OCRText)
fmt.Printf("Detected Language: %s\n", page.DetectedLang)
fmt.Printf("Confidence: %.2f\n", page.OCRConfidence)
```

#### 書籍全体の処理（複数ページ並列処理）

```go
import (
    "context"
    "github.com/google/uuid"
    ocrService "github.com/clearclown/HaiLanGo/backend/internal/service/ocr"
)

ctx := context.Background()
bookID := uuid.New()

// ページデータを準備
pages := []ocrService.PageData{
    {PageID: uuid.New(), ImageData: []byte("page1 data")},
    {PageID: uuid.New(), ImageData: []byte("page2 data")},
    {PageID: uuid.New(), ImageData: []byte("page3 data")},
}

languages := []string{"ru", "en"}
maxConcurrency := 5 // 最大5ページを並列処理

// 書籍全体をOCR処理
result, err := service.ProcessBook(ctx, bookID, pages, languages, maxConcurrency)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Pages: %d\n", result.TotalPages)
fmt.Printf("Processed: %d\n", result.ProcessedPages)
fmt.Printf("Failed: %d\n", result.FailedPages)
fmt.Printf("Processing Time: %v\n", result.ProcessingTime)
```

### モックシステム

APIキーがなくても開発・テストが可能です。環境変数`USE_MOCK_APIS=true`を設定すると、すべての外部API呼び出しが自動的にモックに置き換わります。

```bash
# モックを使用して開発
USE_MOCK_APIS=true go run cmd/server/main.go

# モックを使用してテスト
USE_MOCK_APIS=true go test ./...
```

### 対応言語

- 日本語 (ja)
- 中国語 (zh)
- ロシア語 (ru)
- 英語 (en)
- スペイン語 (es)
- フランス語 (fr)
- ドイツ語 (de)
- アラビア語 (ar)
- ヘブライ語 (he)
- ペルシャ語 (fa)

その他の言語もサポートされていますが、正確性は言語によって異なります。

## テストカバレッジ

```bash
# カバレッジレポートを生成
go test ./... -coverprofile=coverage.out

# カバレッジをブラウザで表示
go tool cover -html=coverage.out
```

現在のカバレッジ:
- `pkg/cache`: 100%
- `internal/service/ocr`: 94.7%
- `pkg/ocr`: 52.8%

## トラブルシューティング

### モックが使用されない

環境変数を確認してください：

```bash
export USE_MOCK_APIS=true
```

### テストが失敗する

依存関係を更新してください：

```bash
go mod tidy
go test ./...
```

### OCR APIエラー

APIキーが正しく設定されているか確認してください。または、モックを使用してください：

```bash
USE_MOCK_APIS=true
```

## ライセンス

MIT License
