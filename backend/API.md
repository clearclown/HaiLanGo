# HaiLanGo API Documentation

## 概要

HaiLanGo APIは、言語学習本のアップロードと管理機能を提供します。

**ベースURL**: `http://localhost:8080/api/v1`

## 認証

現在の実装では認証は未実装です。将来的にはJWTトークンによる認証を実装する予定です。

## エンドポイント

### ヘルスチェック

サーバーの状態を確認します。

**エンドポイント**: `GET /api/v1/health`

**レスポンス**:
```json
{
  "status": "ok",
  "message": "HaiLanGo API is running"
}
```

---

### 書籍の作成

新しい書籍を作成します。

**エンドポイント**: `POST /api/v1/books`

**リクエストボディ**:
```json
{
  "title": "ロシア語入門",
  "target_language": "ru",
  "native_language": "ja",
  "reference_language": "en"
}
```

**フィールド**:
- `title` (必須): 書籍のタイトル
- `target_language` (必須): 学習先言語（ISO 639-1コード）
- `native_language` (必須): 母国語（ISO 639-1コード）
- `reference_language` (オプション): 参照言語（ISO 639-1コード）

**レスポンス** (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "660e8400-e29b-41d4-a716-446655440000",
  "title": "ロシア語入門",
  "target_language": "ru",
  "native_language": "ja",
  "reference_language": "en",
  "cover_image_url": "",
  "total_pages": 0,
  "status": "uploading",
  "created_at": "2025-11-13T08:00:00Z",
  "updated_at": "2025-11-13T08:00:00Z"
}
```

**エラーレスポンス**:
- `400 Bad Request`: 無効なリクエストボディ
- `500 Internal Server Error`: サーバーエラー

---

### ファイルのアップロード

書籍のファイル（PDF、PNG、JPG、HEIC）をアップロードします。

**エンドポイント**: `POST /api/v1/books/:book_id/upload`

**リクエスト**:
- Content-Type: `multipart/form-data`
- フィールド: `files` (複数ファイル可能)

**パラメータ**:
- `book_id` (必須): 書籍ID（UUID）

**制限事項**:
- 最大ファイルサイズ: 100MB/ファイル
- 対応フォーマット: PDF, PNG, JPG, JPEG, HEIC

**curl例**:
```bash
curl -X POST \
  http://localhost:8080/api/v1/books/550e8400-e29b-41d4-a716-446655440000/upload \
  -F "files=@/path/to/page1.png" \
  -F "files=@/path/to/page2.png"
```

**レスポンス** (200 OK):
```json
{
  "message": "files uploaded successfully",
  "count": 2,
  "files": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440000",
      "book_id": "550e8400-e29b-41d4-a716-446655440000",
      "file_name": "page1.png",
      "file_type": "png",
      "file_size": 1048576,
      "storage_path": "users/660e8400.../books/550e8400.../page1_a1b2c3d4.png",
      "uploaded_at": "2025-11-13T08:00:00Z"
    },
    {
      "id": "880e8400-e29b-41d4-a716-446655440000",
      "book_id": "550e8400-e29b-41d4-a716-446655440000",
      "file_name": "page2.png",
      "file_type": "png",
      "file_size": 2097152,
      "storage_path": "users/660e8400.../books/550e8400.../page2_e5f6g7h8.png",
      "uploaded_at": "2025-11-13T08:00:00Z"
    }
  ]
}
```

**エラーレスポンス**:
- `400 Bad Request`:
  - 無効なbook_id
  - ファイルが提供されていない
  - サポートされていないファイル形式
  - ファイルサイズが制限を超えている
- `500 Internal Server Error`: サーバーエラー

---

### アップロード進捗の取得

ファイルアップロードの進捗状況を取得します。

**エンドポイント**: `GET /api/v1/books/:book_id/upload-status`

**パラメータ**:
- `book_id` (必須): 書籍ID（UUID）

**レスポンス** (200 OK):
```json
{
  "book_id": "550e8400-e29b-41d4-a716-446655440000",
  "total_files": 10,
  "uploaded_files": 7,
  "total_bytes": 10485760,
  "uploaded_bytes": 7340032,
  "status": "uploading",
  "message": ""
}
```

**ステータス値**:
- `uploading`: アップロード中
- `completed`: 完了
- `failed`: 失敗

**エラーレスポンス**:
- `400 Bad Request`: 無効なbook_id
- `404 Not Found`: 進捗情報が見つからない
- `500 Internal Server Error`: サーバーエラー

---

## エラーハンドリング

すべてのエラーレスポンスは以下の形式で返されます：

```json
{
  "error": "エラーメッセージ",
  "details": "詳細なエラー情報（オプション）"
}
```

---

## 使用例

### 1. 書籍を作成してファイルをアップロード

```bash
# 1. 書籍を作成
BOOK_RESPONSE=$(curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "ロシア語入門",
    "target_language": "ru",
    "native_language": "ja"
  }')

# book_idを取得
BOOK_ID=$(echo $BOOK_RESPONSE | jq -r '.id')

# 2. ファイルをアップロード
curl -X POST http://localhost:8080/api/v1/books/$BOOK_ID/upload \
  -F "files=@page1.png" \
  -F "files=@page2.png" \
  -F "files=@page3.png"

# 3. 進捗を確認
curl http://localhost:8080/api/v1/books/$BOOK_ID/upload-status
```

---

## 開発・テスト

### サーバーの起動

```bash
cd backend
go run cmd/server/main.go
```

環境変数:
- `BACKEND_PORT`: サーバーポート（デフォルト: 8080）
- `STORAGE_PATH`: ファイル保存先（デフォルト: ./storage）

### テストの実行

```bash
# すべてのテストを実行
go test ./...

# カバレッジ付き
go test -cover ./...

# 特定のパッケージのみ
go test ./pkg/file/...
go test ./pkg/storage/...
go test ./internal/service/...
```

---

## 今後の実装予定

- [ ] JWT認証
- [ ] チャンクアップロード対応
- [ ] S3互換ストレージ対応
- [ ] WebSocket による進捗のリアルタイム通知
- [ ] OCR処理キューへの自動追加
- [ ] データベース連携（PostgreSQL）
- [ ] Redis キャッシュ
- [ ] レート制限
- [ ] ファイル圧縮・最適化

---

## ライセンス

MIT License
