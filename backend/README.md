# HaiLanGo Backend

AI言語学習プラットフォームのバックエンドAPI実装です。

## 概要

HaiLanGoは、既存の言語学習本とAI技術を組み合わせ、個人に最適化された言語学習環境を提供します。

### 主な機能

- ✅ **ユーザー認証**: JWT認証、OAuth対応
- ✅ **OCR処理**: 書籍画像からテキスト抽出
- ✅ **TTS（音声合成）**: 多言語対応の音声読み上げ
- ✅ **STT（音声認識）**: 発音評価・フィードバック
- ✅ **SRS（間隔反復学習）**: 効率的な復習スケジューリング
- ✅ **辞書API**: 単語検索・詳細解説
- ✅ **教師モード**: 自動学習モード
- ✅ **WebSocket通知**: リアルタイム進捗通知
- ✅ **動的言語サポート**: 無制限の言語対応 ⭐ NEW

## 技術スタック

- **言語**: Go 1.21+
- **フレームワーク**: Gin Web Framework
- **データベース**: PostgreSQL 15+ / InMemory（開発用）
- **キャッシュ**: Redis 7+
- **認証**: JWT (RS256)
- **コンテナ**: Podman / Docker

## クイックスタート

### 1. 依存関係のインストール

```bash
go mod download
```

### 2. 環境変数の設定

```bash
cp .env.example .env
```

#### 最小構成（APIキーなし）

```bash
# アプリケーション
APP_ENV=development
BACKEND_PORT=8080

# モック使用（APIキーなしで開発可能）
USE_MOCK_APIS=true

# JWT（開発用）
JWT_SECRET=dev-secret-key-change-in-production
```

#### 本番構成（外部API使用）

```bash
# Google Cloud APIs
GOOGLE_CLOUD_VISION_API_KEY=your_key_here      # OCR用
GOOGLE_CLOUD_TTS_API_KEY=your_key_here         # 音声合成用
GOOGLE_CLOUD_STT_API_KEY=your_key_here         # 音声認識用

# OpenAI (オプション)
OPENAI_API_KEY=your_key_here                   # Whisper STT用

# Stripe (決済)
STRIPE_SECRET_KEY=sk_test_your_key_here
STRIPE_PUBLISHABLE_KEY=pk_test_your_key_here

# データベース
DATABASE_URL=postgresql://HaiLanGo:password@localhost:5432/HaiLanGo_dev
REDIS_URL=redis://localhost:6379
```

### 3. サーバー起動

```bash
go run cmd/server/main.go
```

サーバーは `http://localhost:8080` で起動します。

### 4. テスト実行

```bash
# すべてのテスト
go test ./...

# カバレッジ付き
go test ./... -cover

# 特定パッケージ
go test ./pkg/language/... -v
```

## 動的言語サポート ⭐ NEW

### 概要

HaiLanGoは**無制限の言語サポート**を提供します。LLM/TTS/STT APIがサポートする言語であれば、どんなマイナー言語でも使用可能です。

### サポートティア

| ティア | 説明 | 言語数 |
|--------|------|--------|
| **Verified** | テスト済み・品質保証 | 9言語 |
| **Supported** | 動作確認済み | 21言語 |
| **Experimental** | 未検証（APIが対応すれば動作） | 無制限 |

### Verified言語（9言語）

- 日本語 (ja), 英語 (en), 中国語 (zh), ロシア語 (ru)
- スペイン語 (es), フランス語 (fr), ドイツ語 (de)
- ポルトガル語 (pt), イタリア語 (it)

### Supported言語（21言語）

- ペルシャ語 (fa), ヘブライ語 (he), トルコ語 (tr), 韓国語 (ko)
- アラビア語 (ar), ヒンディー語 (hi), タイ語 (th), ベトナム語 (vi)
- オランダ語 (nl), ポーランド語 (pl), ウクライナ語 (uk), チェコ語 (cs)
- スウェーデン語 (sv), デンマーク語 (da), フィンランド語 (fi), ノルウェー語 (no)
- ギリシャ語 (el), ハンガリー語 (hu), ルーマニア語 (ro)
- インドネシア語 (id), マレー語 (ms)

### Experimental言語

以下のようなマイナー言語も使用可能（APIが対応していれば）：

- クルド語 (ku), エスペラント語 (eo), ウェールズ語 (cy)
- アイルランド語 (ga), マルタ語 (mt), アイスランド語 (is)
- その他、有効なISO 639-1/2コードを持つすべての言語

### 使用方法

```go
import "github.com/clearclown/HaiLanGo/backend/pkg/language"

// レジストリ取得
registry := language.GetRegistry()

// 言語情報取得（未知の言語でもexperimentalとして返される）
info := registry.Get("ku")  // クルド語
fmt.Println(info.SupportTier)  // "experimental"

// 言語コード検証
if language.IsValidCode("ja") {
    // 有効なISO 639-1/2コード
}
```

## 外部API

### 必須API（本番環境）

| API | 用途 | 取得方法 |
|-----|------|---------|
| **Google Cloud Vision** | OCR（テキスト認識） | [Google Cloud Console](https://console.cloud.google.com/) |
| **Google Cloud TTS** | 音声合成 | 同上 |
| **Google Cloud STT** | 音声認識 | 同上 |

### オプションAPI

| API | 用途 | 取得方法 |
|-----|------|---------|
| **OpenAI Whisper** | 高精度STT | [OpenAI Platform](https://platform.openai.com/) |
| **Stripe** | 決済処理 | [Stripe Dashboard](https://dashboard.stripe.com/) |
| **DeepL** | 高品質翻訳 | [DeepL API](https://www.deepl.com/pro-api) |

### モック開発

**APIキーなしでも開発可能です！**

```bash
USE_MOCK_APIS=true go run cmd/server/main.go
```

詳細は [モック構築戦略](../docs/mocking_strategy.md) を参照。

## APIエンドポイント

### 認証

```
POST /api/v1/auth/register     # ユーザー登録
POST /api/v1/auth/login        # ログイン
POST /api/v1/auth/refresh      # トークン更新
POST /api/v1/auth/logout       # ログアウト
```

### 書籍管理

```
POST   /api/v1/books           # 書籍作成
GET    /api/v1/books           # 一覧取得
GET    /api/v1/books/:id       # 詳細取得
DELETE /api/v1/books/:id       # 削除
```

### OCR処理

```
POST /api/v1/ocr/process       # ページ処理
GET  /api/v1/ocr/jobs/:id      # ジョブ状態取得
```

### TTS（音声合成）

```
POST /api/v1/tts/synthesize    # 音声生成
GET  /api/v1/tts/languages     # 対応言語一覧
GET  /api/v1/tts/jobs/:id      # ジョブ状態取得
```

### STT（発音評価）

```
POST /api/v1/stt/recognize     # 音声認識・発音評価
GET  /api/v1/stt/languages     # 対応言語一覧
GET  /api/v1/stt/jobs/:id      # ジョブ状態取得
```

### SRS（間隔反復学習）

```
GET  /api/v1/review/items      # 復習項目取得
POST /api/v1/review/complete   # 復習完了
GET  /api/v1/review/stats      # 統計情報
```

### ヘルスチェック

```
GET /health                    # サーバー状態確認
```

## プロジェクト構造

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # エントリーポイント
├── internal/
│   ├── api/
│   │   ├── handler/             # HTTPハンドラー
│   │   ├── middleware/          # ミドルウェア
│   │   └── router/              # ルーティング
│   ├── service/                 # ビジネスロジック
│   │   ├── tts/                 # TTS処理
│   │   ├── stt/                 # STT処理
│   │   ├── ocr/                 # OCR処理
│   │   └── srs/                 # SRS処理
│   ├── repository/              # データアクセス層
│   │   ├── *_inmemory.go        # InMemory実装
│   │   └── *_postgres.go        # PostgreSQL実装
│   └── models/                  # データモデル
├── pkg/
│   ├── language/                # 動的言語サポート ⭐ NEW
│   │   ├── registry.go          # LanguageRegistry
│   │   └── registry_test.go
│   ├── tts/                     # TTSクライアント
│   ├── stt/                     # STTクライアント
│   ├── ocr/                     # OCRクライアント
│   ├── jwt/                     # JWT処理
│   └── validator/               # バリデーション
└── mocks/                       # モックデータ
```

## 開発ガイド

### コーディング規約

- `gofmt` でフォーマット
- `golangci-lint` でリント
- エラーハンドリングは必須
- コメントは日本語でOK

### 新機能の追加

1. テストを書く（TDD）
2. 実装する
3. テストが通ることを確認
4. リファクタリング

### 言語サポートの追加

新しい言語を**Verified**または**Supported**に追加する場合：

```go
// pkg/language/registry.go の initWellKnownLanguages() を編集

verified := []*Info{
    // 既存の言語...
    {Code: "新コード", Name: "言語名", NativeName: "ネイティブ名",
     SupportTier: TierVerified, SupportsPronunciation: true},
}
```

## テスト

### 現在のテストカバレッジ

```
pkg/language         ✅ 100%  # 動的言語サポート
pkg/tts              ✅ 95%   # TTS処理
pkg/stt              ✅ 95%   # STT処理
pkg/srs              ✅ 100%  # SRSアルゴリズム
internal/repository  ✅ 90%   # リポジトリ
internal/service     ✅ 85%   # ビジネスロジック
internal/api/handler ✅ 80%   # APIハンドラー
```

## トラブルシューティング

### よくある問題

**Q: APIキーがないとテストが動かない**
```bash
# モックを使用してテスト
USE_MOCK_APIS=true go test ./...
```

**Q: 言語がサポートされていないエラー**
```
動的言語サポートにより、有効なISO 639-1/2コードであれば
すべての言語がexperimentalとして使用可能です。
```

**Q: PostgreSQLがない**
```
InMemoryリポジトリが自動的に使用されます。
開発・テストには影響ありません。
```

## ライセンス

MIT License

## サポート

- GitHub Issues: バグ報告・機能リクエスト
- ドキュメント: [docs/](../docs/)
