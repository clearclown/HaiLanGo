# HaiLanGo Backend - User Authentication API

## 概要

HaiLanGoプロジェクトのバックエンドAPI実装です。JWT認証、ユーザー管理、セキュアなパスワードハッシュ化を提供します。

## 技術スタック

- **言語**: Go 1.21+
- **フレームワーク**: Gin Web Framework
- **データベース**: PostgreSQL 15+
- **認証**: JWT (RS256)
- **パスワードハッシュ**: bcrypt (cost 10)

## セットアップ

### 1. 依存関係のインストール

```bash
go mod download
```

### 2. 環境変数の設定

`.env`ファイルを作成して以下を設定：

```bash
# サーバー設定
BACKEND_PORT=8080

# データベース設定
DATABASE_URL=postgresql://HaiLanGo:password@localhost:5432/HaiLanGo_dev?sslmode=disable

# Redis設定（将来使用）
REDIS_URL=redis://localhost:6379
```

### 3. データベースマイグレーション

```bash
# PostgreSQLに接続
psql -U HaiLanGo -d HaiLanGo_dev

# マイグレーション実行
\i migrations/001_create_users_table.up.sql
\i migrations/002_create_refresh_tokens_table.up.sql
```

### 4. サーバー起動

```bash
go run cmd/server/main.go
```

サーバーは `http://localhost:8080` で起動します。

## API エンドポイント

### ヘルスチェック

```
GET /health
```

**レスポンス:**
```json
{
  "status": "ok"
}
```

### ユーザー登録

```
POST /api/v1/auth/register
```

**リクエスト:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "display_name": "Test User"
}
```

**レスポンス:**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "display_name": "Test User",
    "email_verified": false,
    "created_at": "2025-11-13T08:00:00Z",
    "updated_at": "2025-11-13T08:00:00Z"
  },
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "random_token_string",
  "expires_in": 900
}
```

### ログイン

```
POST /api/v1/auth/login
```

**リクエスト:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**レスポンス:** ユーザー登録と同じ

### トークンリフレッシュ

```
POST /api/v1/auth/refresh
```

**リクエスト:**
```json
{
  "refresh_token": "random_token_string"
}
```

**レスポンス:** ユーザー登録と同じ

### ログアウト

```
POST /api/v1/auth/logout
```

**リクエスト:**
```json
{
  "refresh_token": "random_token_string"
}
```

**レスポンス:**
```json
{
  "message": "ログアウトしました"
}
```

## セキュリティ

### パスワード要件

- 最小8文字
- 大文字、小文字、数字、記号のうち3種類以上を含む

### JWT トークン

- **アクセストークン**: 15分有効
- **リフレッシュトークン**: 7日有効
- **署名方式**: RS256（RSA）

### レート制限

- 1分あたり100リクエストまで
- IPアドレスベース

## テスト

### すべてのテストを実行

```bash
go test ./... -v
```

### カバレッジ付きテスト

```bash
go test ./... -cover
```

### 特定パッケージのテスト

```bash
# パスワードパッケージ
go test ./pkg/password/... -v

# JWTパッケージ
go test ./pkg/jwt/... -v

# リポジトリ
go test ./internal/repository/... -v
```

## プロジェクト構造

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # エントリーポイント
├── internal/
│   ├── api/
│   │   ├── handler/          # HTTPハンドラー
│   │   ├── middleware/       # ミドルウェア
│   │   └── router/           # ルーティング
│   ├── service/              # ビジネスロジック
│   ├── repository/           # データアクセス層
│   └── models/               # データモデル
├── pkg/
│   ├── jwt/                  # JWT生成・検証
│   └── password/             # パスワードハッシュ化
├── migrations/               # データベースマイグレーション
├── go.mod
└── go.sum
```

## エラーハンドリング

### 一般的なエラーレスポンス

```json
{
  "error": "エラーメッセージ"
}
```

### HTTPステータスコード

- `200 OK`: 成功
- `201 Created`: リソース作成成功
- `400 Bad Request`: リクエストが不正
- `401 Unauthorized`: 認証エラー
- `409 Conflict`: リソースの競合（重複メールアドレスなど）
- `429 Too Many Requests`: レート制限超過
- `500 Internal Server Error`: サーバーエラー

## 開発

### コーディング規約

- `gofmt`でフォーマット
- `golangci-lint`でリント
- エラーハンドリングは必須
- コメントは日本語で記述

### Git ワークフロー

```bash
# 機能ブランチの作成
git checkout -b feature/your-feature

# 変更をコミット
git add .
git commit -m "feat: 機能の説明"

# プッシュ
git push origin feature/your-feature
```

## ライセンス

MIT License

## サポート

質問や問題がある場合は、GitHub Issuesを開いてください。
