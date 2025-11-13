<<<<<<< HEAD
# HaiLanGo Backend - 間隔反復学習（SRS）機能

## 概要

間隔反復学習（Spaced Repetition System）機能のバックエンド実装です。

## 機能

- ✅ SRSアルゴリズム実装
- ✅ 復習項目管理
- ✅ 優先度別復習項目取得（緊急・推奨・余裕あり）
- ✅ 復習完了処理
- ✅ 学習統計
- ✅ RESTful API
=======
# HaiLanGo Backend - User Authentication API

## 概要

HaiLanGoプロジェクトのバックエンドAPI実装です。JWT認証、ユーザー管理、セキュアなパスワードハッシュ化を提供します。
>>>>>>> origin/main

## 技術スタック

- **言語**: Go 1.21+
<<<<<<< HEAD
- **フレームワーク**: Gin
- **テスト**: testify
- **データベース**: PostgreSQL（予定）、Redis（予定）

## セットアップ

### 依存関係のインストール
=======
- **フレームワーク**: Gin Web Framework
- **データベース**: PostgreSQL 15+
- **認証**: JWT (RS256)
- **パスワードハッシュ**: bcrypt (cost 10)

## セットアップ

### 1. 依存関係のインストール
>>>>>>> origin/main

```bash
go mod download
```

<<<<<<< HEAD
### サーバーの起動
=======
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
>>>>>>> origin/main

```bash
go run cmd/server/main.go
```

<<<<<<< HEAD
サーバーはポート8080で起動します。

### テストの実行

```bash
# すべてのテスト
go test ./... -v

# カバレッジ付き
go test ./... -cover

# 特定のパッケージ
go test ./pkg/srs/... -v
```

## APIエンドポイント

### 復習項目取得
```
GET /api/v1/review/items/:user_id
```

### 復習完了
```
POST /api/v1/review/items/:item_id/complete
Content-Type: application/json

{
  "score": 85,
  "time_spent_sec": 30
}
```

### 統計情報取得
```
GET /api/v1/review/stats/:user_id
```

### ヘルスチェック
=======
サーバーは `http://localhost:8080` で起動します。

## API エンドポイント

### ヘルスチェック

>>>>>>> origin/main
```
GET /health
```

<<<<<<< HEAD
## ディレクトリ構造
=======
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
>>>>>>> origin/main

```
backend/
├── cmd/
│   └── server/
<<<<<<< HEAD
│       └── main.go              # エントリーポイント
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   └── review_handler.go  # APIハンドラー
│   │   └── router/
│   │       └── router.go           # ルーティング
│   ├── service/
│   │   └── srs/
│   │       ├── srs.go              # ビジネスロジック
│   │       └── srs_test.go
│   ├── repository/
│   │   ├── mock_review_item.go     # リポジトリ
│   │   └── review_item_test.go
│   └── models/
│       └── review_item.go          # データモデル
└── pkg/
    └── srs/
        ├── algorithm.go            # SRSアルゴリズム
        └── algorithm_test.go
```

## SRSアルゴリズム

### 基本間隔
- 初回: 1日後
- 2回目: 3日後
- 3回目: 7日後
- 4回目: 14日後
- 5回目: 30日後
- 6回目以降: 60日後

### スコア調整
- 85点以上: 間隔を1.5倍に延長
- 70-84点: 通常の間隔
- 50-69点: 間隔を半分に短縮
- 50点未満: 翌日に再度復習

## テスト

すべてのテストがPASSしています：

- ✅ アルゴリズムテスト（6/6）
- ✅ リポジトリテスト（6/6）
- ✅ サービス層テスト（5/5）
=======
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
>>>>>>> origin/main

## 開発

### コーディング規約

<<<<<<< HEAD
- コメントは日本語で記述
- `gofmt`でフォーマット
- テストファーストで実装（TDD）

### 新機能の追加

1. テストを書く
2. 実装する
3. テストが通ることを確認
4. リファクタリング

## 今後の予定

- [ ] PostgreSQL実装
- [ ] Redis キャッシュ
- [ ] JWT認証
- [ ] WebSocket通知
- [ ] カスタマイズ設定
=======
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
>>>>>>> origin/main

## ライセンス

MIT License
<<<<<<< HEAD
=======

## サポート

質問や問題がある場合は、GitHub Issuesを開いてください。
>>>>>>> origin/main
