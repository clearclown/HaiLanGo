# HaiLanGo Backend

HaiLanGo（ハイランゴー）のバックエンドAPIサーバーです。

## 技術スタック

- **言語**: Go 1.21+
- **Webフレームワーク**: Gin
- **テストフレームワーク**: testify
- **データベース**: PostgreSQL（今後実装）
- **キャッシュ**: Redis（今後実装）
- **ストレージ**: ローカルファイルシステム / S3互換ストレージ

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
│   ├── repository/              # データアクセス層（今後実装）
│   └── models/                  # データモデル
├── pkg/
│   ├── storage/                 # ストレージ操作
│   └── file/                    # ファイル検証
├── go.mod
├── go.sum
├── API.md                       # API ドキュメント
└── README.md                    # このファイル
```

## セットアップ

### 前提条件

- Go 1.21 以上
- Git

### インストール

```bash
# リポジトリをクローン
git clone https://github.com/clearclown/HaiLanGo.git
cd HaiLanGo/backend

# 依存関係をインストール
go mod download

# テストを実行
go test ./...
```

### 環境変数

`.env` ファイルを作成して以下の環境変数を設定（オプション）：

```bash
# サーバーポート（デフォルト: 8080）
BACKEND_PORT=8080

# ファイル保存先（デフォルト: ./storage）
STORAGE_PATH=./storage
```

### サーバーの起動

```bash
# 開発モードで起動
go run cmd/server/main.go

# ビルドして実行
go build -o bin/server cmd/server/main.go
./bin/server
```

サーバーは `http://localhost:8080` で起動します。

## API ドキュメント

詳細なAPI仕様は [API.md](./API.md) を参照してください。

### 主要なエンドポイント

- `GET /api/v1/health` - ヘルスチェック
- `POST /api/v1/books` - 書籍の作成
- `POST /api/v1/books/:book_id/upload` - ファイルのアップロード
- `GET /api/v1/books/:book_id/upload-status` - アップロード進捗の取得

### 使用例

```bash
# ヘルスチェック
curl http://localhost:8080/api/v1/health

# 書籍を作成
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "ロシア語入門",
    "target_language": "ru",
    "native_language": "ja"
  }'

# ファイルをアップロード
curl -X POST http://localhost:8080/api/v1/books/550e8400-e29b-41d4-a716-446655440000/upload \
  -F "files=@page1.png" \
  -F "files=@page2.png"
```

## テスト

### すべてのテストを実行

```bash
go test ./...
```

### カバレッジ付きテスト

```bash
go test -cover ./...
```

### 特定のパッケージのテスト

```bash
# ファイル検証
go test ./pkg/file/... -v

# ストレージ
go test ./pkg/storage/... -v

# アップロードサービス
go test ./internal/service/... -v
```

### テスト結果（2025-11-13現在）

```
✅ pkg/file: 全テスト合格
✅ pkg/storage: 全テスト合格
✅ internal/service: 全テスト合格
```

## 実装済み機能

### ✅ Phase 1: 書籍アップロード機能

- [x] ファイル検証（サイズ、形式）
  - 対応フォーマット: PDF, PNG, JPG, JPEG, HEIC
  - 最大ファイルサイズ: 100MB
- [x] ローカルストレージへの保存
  - ファイル名のサニタイズ
  - ユニークなファイルパスの生成
  - ユーザーごとのディレクトリ分離
- [x] 書籍メタデータ管理
  - タイトル、学習先言語、母国語、参照言語
  - 書籍ステータス管理（アップロード中、処理中、準備完了）
- [x] アップロード進捗追跡
  - リアルタイムの進捗状況
  - ファイル数とバイト数の追跡
- [x] RESTful API
  - 書籍作成
  - ファイルアップロード
  - 進捗取得
- [x] 包括的なテストカバレッジ
  - ユニットテスト
  - 統合テスト

## 今後の実装予定

### Phase 2: データベース統合

- [ ] PostgreSQL データベースのセットアップ
- [ ] 書籍リポジトリの実装
- [ ] マイグレーション管理
- [ ] トランザクション管理

### Phase 3: 認証・認可

- [ ] JWT 認証
- [ ] ユーザー登録・ログイン
- [ ] OAuth 2.0（Google Login）
- [ ] セッション管理

### Phase 4: OCR 統合

- [ ] OCR 処理キューの実装
- [ ] Google Vision API 統合
- [ ] OCR 結果のキャッシュ（Redis）
- [ ] 手動修正機能

### Phase 5: 高度な機能

- [ ] チャンクアップロード対応
- [ ] S3互換ストレージ対応
- [ ] WebSocket による進捗のリアルタイム通知
- [ ] レート制限
- [ ] ファイル圧縮・最適化

## コーディング規約

### Go 言語

- `gofmt` でフォーマット（保存時に自動実行）
- `golangci-lint` を使用してリント
- エラーハンドリングは必ず行う（`if err != nil` を省略しない）
- コメントは日本語で記述

### 命名規則

- パッケージ名: 小文字、単語区切りなし（例: `userservice`, `bookrepository`）
- インターフェース名: 動詞 + "er" または 名詞 + "Service" （例: `BookReader`, `UserService`）
- 構造体名: PascalCase（例: `UserProfile`, `BookMetadata`）
- メソッド/関数名: PascalCase（外部公開）、camelCase（内部のみ）

### テスト駆動開発（TDD）

1. テストを先に作成
2. テストが失敗することを確認（Red）
3. 実装を作成してテストを通す（Green）
4. リファクタリング（Refactor）

## トラブルシューティング

### ポートが既に使用されている

```bash
# ポート8080を使用しているプロセスを確認
lsof -i :8080

# プロセスを終了
kill -9 <PID>

# または別のポートを使用
BACKEND_PORT=8081 go run cmd/server/main.go
```

### ストレージディレクトリが作成できない

```bash
# 権限を確認
ls -la ./

# ストレージディレクトリを手動で作成
mkdir -p ./storage
chmod 755 ./storage
```

## コントリビューション

プルリクエストは大歓迎です！

1. フィーチャーブランチを作成
2. 変更をコミット
3. テストを実行して通ることを確認
4. プルリクエストを作成

## ライセンス

MIT License

---

開発者: [clearclown](https://github.com/clearclown)
