# API仕様書

## 概要

HaiLanGoバックエンドのREST API仕様を定義する。
すべてのAPIは`/api/v1`プレフィックスを持つ。

---

## 共通仕様

### ベースURL

```
開発: http://localhost:8080/api/v1
本番: https://api.hailango.com/api/v1
```

### 認証

JWT Bearer Token認証を使用。

```http
Authorization: Bearer <jwt_token>
```

### レスポンス形式

**成功時**
```json
{
  "success": true,
  "data": { ... }
}
```

**エラー時**
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

### HTTPステータスコード

| コード | 意味 |
|--------|------|
| 200 | 成功 |
| 201 | 作成成功 |
| 400 | リクエスト不正 |
| 401 | 認証エラー |
| 403 | 権限なし |
| 404 | リソースなし |
| 422 | バリデーションエラー |
| 429 | レート制限 |
| 500 | サーバーエラー |

---

## 認証 API

### POST /auth/register

ユーザー登録

**Request**
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "display_name": "太郎"
}
```

**Response** `201 Created`
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "display_name": "太郎"
    },
    "token": "jwt_token"
  }
}
```

---

### POST /auth/login

ログイン

**Request**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "display_name": "太郎"
    },
    "token": "jwt_token"
  }
}
```

---

### POST /auth/logout

ログアウト

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "message": "Logged out successfully"
  }
}
```

---

### GET /auth/me

現在のユーザー情報

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "display_name": "太郎",
    "native_language": "ja",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 教材 (Books) API

### GET /books

教材一覧取得

**Query Parameters**
- `page` (int, default: 1)
- `limit` (int, default: 20, max: 100)
- `status` (string, optional): "processing" | "ready" | "error"

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "books": [
      {
        "id": "uuid",
        "title": "ロシア語入門",
        "target_language": "ru",
        "native_language": "ja",
        "cover_image_url": "https://...",
        "total_pages": 150,
        "status": "ready",
        "progress": {
          "completed_pages": 12,
          "percentage": 8
        },
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

---

### POST /books

教材作成

**Request**
```json
{
  "title": "ロシア語入門",
  "target_language": "ru",
  "native_language": "ja",
  "reference_language": "en"
}
```

**Response** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "title": "ロシア語入門",
    "target_language": "ru",
    "native_language": "ja",
    "reference_language": "en",
    "status": "processing",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### GET /books/:id

教材詳細取得

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "title": "ロシア語入門",
    "target_language": "ru",
    "native_language": "ja",
    "reference_language": "en",
    "total_pages": 150,
    "status": "ready",
    "pages": [
      {
        "id": "uuid",
        "page_number": 1,
        "thumbnail_url": "https://...",
        "ocr_status": "completed"
      }
    ],
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### DELETE /books/:id

教材削除

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "message": "Book deleted successfully"
  }
}
```

---

## アップロード API

### POST /upload/init

チャンクアップロード初期化

**Request**
```json
{
  "book_id": "uuid",
  "filename": "textbook.pdf",
  "file_size": 52428800,
  "content_type": "application/pdf",
  "total_chunks": 50
}
```

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "upload_id": "uuid",
    "chunk_size": 1048576
  }
}
```

---

### POST /upload/chunk

チャンクアップロード

**Request** `multipart/form-data`
- `upload_id`: string
- `chunk_index`: int
- `chunk`: file

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "chunk_index": 5,
    "received": true
  }
}
```

---

### POST /upload/complete

アップロード完了

**Request**
```json
{
  "upload_id": "uuid"
}
```

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "book_id": "uuid",
    "pages_created": 150,
    "ocr_job_id": "uuid"
  }
}
```

---

## ページ API

### GET /books/:book_id/pages/:page_number

ページ詳細取得

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "book_id": "uuid",
    "page_number": 12,
    "image_url": "https://...",
    "ocr_text": "Здравствуйте! Меня зовут...",
    "ocr_confidence": 0.95,
    "ocr_status": "completed",
    "vocabularies": [
      {
        "id": "uuid",
        "word": "Здравствуйте",
        "meaning": "こんにちは（丁寧）"
      }
    ]
  }
}
```

---

### PUT /books/:book_id/pages/:page_number/ocr

OCRテキスト修正

**Request**
```json
{
  "ocr_text": "修正後のテキスト..."
}
```

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "ocr_text": "修正後のテキスト...",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 学習 API

### POST /learning/sessions

学習セッション開始

**Request**
```json
{
  "book_id": "uuid",
  "session_type": "page_learning"
}
```

**Response** `201 Created`
```json
{
  "success": true,
  "data": {
    "session_id": "uuid",
    "started_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### POST /learning/sessions/:session_id/end

学習セッション終了

**Request**
```json
{
  "pages_studied": 3,
  "words_learned": 15
}
```

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "session_id": "uuid",
    "duration_seconds": 1800,
    "pages_studied": 3,
    "words_learned": 15
  }
}
```

---

### POST /learning/log

学習アクションログ

**Request**
```json
{
  "session_id": "uuid",
  "page_id": "uuid",
  "vocabulary_id": "uuid",
  "action_type": "stt_pronounce",
  "score": 85.5,
  "feedback": "良い発音です。Рの巻き舌をもう少し..."
}
```

**Response** `201 Created`
```json
{
  "success": true,
  "data": {
    "log_id": "uuid",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## AI会話 API

### POST /ai/chat

AI会話（ストリーミング対応）

**Request**
```json
{
  "session_id": "uuid",
  "page_id": "uuid",
  "message": "この単語の意味を教えて",
  "context": {
    "target_language": "ru",
    "native_language": "ja",
    "current_text": "Здравствуйте!"
  }
}
```

**Response** `200 OK` (Server-Sent Events)
```
data: {"type": "start"}

data: {"type": "content", "content": "「Здрав"}

data: {"type": "content", "content": "ствуйте」は"}

data: {"type": "content", "content": "丁寧な挨拶で..."}

data: {"type": "end", "total_tokens": 150}
```

---

## TTS API

### POST /tts/synthesize

音声合成

**Request**
```json
{
  "text": "Здравствуйте!",
  "language": "ru",
  "voice": "default",
  "speed": 1.0
}
```

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "audio_url": "https://...",
    "duration_seconds": 1.5,
    "cached": true
  }
}
```

---

## STT API

### POST /stt/evaluate

発音評価

**Request** `multipart/form-data`
- `audio`: file (wav, mp3, webm)
- `expected_text`: string
- `language`: string

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "recognized_text": "Здравствуйте",
    "score": 85,
    "details": {
      "accuracy": 90,
      "fluency": 80,
      "pronunciation": 85
    },
    "feedback": "良い発音です。「вст」の部分をもう少し明確に。",
    "word_scores": [
      {
        "word": "Здравствуйте",
        "score": 85,
        "issue": "вст部分が不明瞭"
      }
    ]
  }
}
```

---

## SRS (復習) API

### GET /review/due

復習対象取得

**Query Parameters**
- `limit` (int, default: 20)

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "vocabulary": {
          "id": "uuid",
          "word": "Здравствуйте",
          "meaning": "こんにちは",
          "context": "..."
        },
        "priority": "urgent",
        "due_at": "2024-01-01T00:00:00Z",
        "repetitions": 3
      }
    ],
    "counts": {
      "urgent": 3,
      "recommended": 5,
      "optional": 4,
      "total": 12
    }
  }
}
```

---

### POST /review/submit

復習結果送信

**Request**
```json
{
  "srs_item_id": "uuid",
  "quality": 4,
  "response_time_ms": 3500
}
```

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "srs_item_id": "uuid",
    "new_interval_days": 7,
    "next_review_at": "2024-01-08T00:00:00Z",
    "easiness_factor": 2.6
  }
}
```

---

## 統計 API

### GET /stats/overview

統計概要

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "streak": {
      "current": 7,
      "longest": 15
    },
    "study_time": {
      "today_seconds": 1800,
      "this_week_seconds": 10800,
      "total_seconds": 86400
    },
    "progress": {
      "words_learned": 230,
      "pages_completed": 45,
      "books_in_progress": 3
    },
    "review": {
      "due_today": 12,
      "completed_today": 8
    }
  }
}
```

---

### GET /stats/history

学習履歴

**Query Parameters**
- `period`: "week" | "month" | "year"

**Response** `200 OK`
```json
{
  "success": true,
  "data": {
    "period": "week",
    "data_points": [
      {
        "date": "2024-01-01",
        "study_time_seconds": 3600,
        "words_learned": 15,
        "pages_completed": 2
      }
    ]
  }
}
```

---

## WebSocket API

### /ws/notifications

リアルタイム通知

**接続**
```
ws://localhost:8080/ws/notifications?token=<jwt>
```

**メッセージ形式**
```json
{
  "type": "ocr_progress",
  "data": {
    "book_id": "uuid",
    "current_page": 50,
    "total_pages": 150,
    "percentage": 33
  }
}
```

**イベントタイプ**
- `ocr_progress`: OCR処理進捗
- `ocr_complete`: OCR完了
- `ocr_error`: OCRエラー
- `review_reminder`: 復習リマインダー

---

## エラーコード一覧

| コード | HTTP | 説明 |
|--------|------|------|
| AUTH_INVALID_CREDENTIALS | 401 | 認証情報が不正 |
| AUTH_TOKEN_EXPIRED | 401 | トークン期限切れ |
| AUTH_REQUIRED | 401 | 認証が必要 |
| VALIDATION_ERROR | 422 | バリデーションエラー |
| BOOK_NOT_FOUND | 404 | 教材が見つからない |
| PAGE_NOT_FOUND | 404 | ページが見つからない |
| UPLOAD_TOO_LARGE | 413 | ファイルサイズ超過 |
| UPLOAD_INVALID_TYPE | 422 | 非対応ファイル形式 |
| OCR_PROCESSING | 202 | OCR処理中 |
| OCR_FAILED | 500 | OCR処理失敗 |
| TTS_FAILED | 500 | TTS処理失敗 |
| STT_FAILED | 500 | STT処理失敗 |
| RATE_LIMITED | 429 | レート制限 |
| INTERNAL_ERROR | 500 | 内部エラー |

---

## 次のドキュメント

- [05_TEST_STRATEGY.md](./05_TEST_STRATEGY.md) - テスト戦略（Testcontainers）
