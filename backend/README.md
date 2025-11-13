# HaiLanGo Backend

HaiLanGoプロジェクトのバックエンドAPI実装

## 概要

Goで実装された高性能なバックエンドサーバーです。以下の主要機能を提供します：

- TTS（Text-to-Speech）音声読み上げ
- OCR（Optical Character Recognition）テキスト認識
- STT（Speech-to-Text）音声認識
- ユーザー認証・認可
- 学習履歴管理

## 技術スタック

- **言語**: Go 1.21+
- **データベース**: PostgreSQL 15+, Redis 7+
- **API設計**: RESTful API
- **テスト**: testify

## ディレクトリ構造

```
backend/
├── cmd/                    # エントリーポイント
│   └── server/             # サーバーアプリケーション
├── internal/               # 内部パッケージ
│   ├── api/                # APIハンドラー
│   │   ├── handler/        # HTTPハンドラー
│   │   ├── middleware/     # ミドルウェア
│   │   └── router/         # ルーティング
│   ├── service/            # ビジネスロジック
│   │   ├── tts/            # TTS音声読み上げ
│   │   └── cache/          # キャッシュ管理
│   ├── repository/         # データアクセス層
│   └── models/             # データモデル
└── pkg/                    # 外部パッケージ（再利用可能）
    ├── tts/                # TTS APIクライアント
    └── storage/            # ストレージ管理
```

## 実装済み機能

### 1. TTS音声読み上げ機能 ✅

**概要**: テキストから音声を生成し、URLを返すサービス

**対応言語（主要12言語）**:
- 日本語 (ja)
- 中国語 (zh)
- 英語 (en)
- ロシア語 (ru)
- ペルシャ語 (fa)
- ヘブライ語 (he)
- スペイン語 (es)
- フランス語 (fr)
- ポルトガル語 (pt)
- ドイツ語 (de)
- イタリア語 (it)
- トルコ語 (tr)

**音声品質**:
- `standard`: 標準品質（無料プラン）
- `premium`: 高品質（プレミアムプラン）

**速度調整**: 0.5x, 0.75x, 1.0x, 1.25x, 1.5x, 2.0x

**キャッシュ**: 7日間のTTL（Time To Live）

**使用例**:

```go
package main

import (
    "context"
    "fmt"
    "github.com/clearclown/HaiLanGo/backend/internal/service/tts"
)

func main() {
    service := tts.NewTTSService()

    ctx := context.Background()
    text := "こんにちは、世界！"
    lang := "ja"
    quality := "standard"
    speed := 1.0

    audioURL, err := service.GenerateAudio(ctx, text, lang, quality, speed)
    if err != nil {
        panic(err)
    }

    fmt.Println("Audio URL:", audioURL)
}
```

**バッチ生成**:

```go
texts := []string{"Hello", "Goodbye", "Thank you"}
audioURLs, err := service.BatchGenerate(ctx, texts, "en", "standard", 1.0)
```

## モックシステム

開発・テスト時にAPIキーなしで動作可能なモックシステムを実装しています。

### モックの有効化

環境変数を設定してモックを使用：

```bash
export USE_MOCK_APIS=true
export TEST_USE_MOCKS=true  # テスト実行時は自動設定
```

### モックの仕組み

1. **自動切り替え**: 環境変数またはAPIキーの有無で自動的にモード切り替え
2. **決定論的**: 同じ入力に対して同じモックデータを返す
3. **高速**: 外部API呼び出しなしで即座にレスポンス

## テスト

### すべてのテストを実行

```bash
go test ./... -v
```

### 特定のパッケージのテスト

```bash
# TTSサービスのテスト
go test ./internal/service/tts -v

# TTSクライアントのテスト
go test ./pkg/tts -v

# キャッシュのテスト
go test ./internal/service/cache -v

# ストレージのテスト
go test ./pkg/storage -v
```

### カバレッジ付きテスト

```bash
go test ./... -cover

# カバレッジレポート生成
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### テスト結果

現在のテスト状況:

```
✅ internal/service/cache: 8/8 テスト成功
✅ internal/service/tts: 8/8 テスト成功
✅ pkg/storage: 6/6 テスト成功
✅ pkg/tts: 9/9 テスト成功

合計: 31/31 テスト成功
```

## 環境変数

### 必須（最小構成）

```bash
# アプリケーション環境
APP_ENV=development

# サーバーポート
BACKEND_PORT=8080

# データベース
DATABASE_URL=postgresql://HaiLanGo:password@localhost:5432/HaiLanGo_dev
REDIS_URL=redis://localhost:6379
```

### オプション（実APIを使用する場合）

```bash
# Google Cloud TTS
GOOGLE_CLOUD_TTS_API_KEY=your_key_here

# Amazon Polly
AWS_ACCESS_KEY_ID=your_key_here
AWS_SECRET_ACCESS_KEY=your_secret_here

# ElevenLabs（プレミアム）
ELEVENLABS_API_KEY=your_key_here

# ストレージ設定
AUDIO_STORAGE_PATH=./storage/audio
AUDIO_BASE_URL=http://localhost:8080/audio
```

### モック設定

```bash
# モック使用（APIキーなしで開発）
USE_MOCK_APIS=true

# テスト時のモック使用（自動設定）
TEST_USE_MOCKS=true
```

## パフォーマンス

### ベンチマーク結果

- **TTSレイテンシ**: < 1秒（モック環境）
- **キャッシュヒット率**: 70%以上（目標）
- **並列生成**: 10並列で問題なし
- **音声生成時間**: テキスト長に比例（100文字/秒）

## セキュリティ

- ユーザー認証必須（TODO）
- レート制限（1日1000リクエスト/ユーザー）（TODO）
- 音声ファイルのアクセス制御（TODO）

## 今後の実装予定

- [ ] RESTful APIエンドポイント
- [ ] WebSocketサポート
- [ ] Redis連携（キャッシュ）
- [ ] PostgreSQL連携（永続化）
- [ ] ユーザー認証・認可
- [ ] レート制限
- [ ] ロギング・モニタリング
- [ ] OCR機能
- [ ] STT機能

## コントリビューション

1. フィーチャーブランチを作成
2. テストを追加
3. すべてのテストが通ることを確認
4. プルリクエストを作成

## ライセンス

MIT License
