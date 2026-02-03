# API Specification

## 1. Overview

HaiLanGo exposes a RESTful API built with `reinhardt-rest` and real-time WebSocket endpoints via `reinhardt-websockets`.

### Base URL

- **Development**: `http://localhost:8080/api`
- **Production**: `https://api.hailango.com/api`

### Authentication

All endpoints (except `/auth/*`) require a valid JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

### Response Format

All responses follow a consistent JSON structure:

```json
// Success
{
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100
  }
}

// Error
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "details": {
      "field": "email",
      "constraint": "email"
    }
  }
}
```

---

## 2. Authentication Endpoints

### 2.1 Register

Create a new user account.

```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123",
  "display_name": "John Doe",
  "native_language": "ja"
}
```

**Response (201 Created):**
```json
{
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "display_name": "John Doe",
      "native_language": "ja",
      "email_verified": false,
      "created_at": "2025-01-15T10:30:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g...",
      "expires_in": 3600
    }
  }
}
```

### 2.2 Login

Authenticate with email and password.

```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "display_name": "John Doe"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g...",
      "expires_in": 3600
    }
  }
}
```

### 2.3 OAuth Login

Authenticate via OAuth provider (Google).

```http
POST /api/auth/oauth/google
Content-Type: application/json

{
  "code": "4/0AX4XfWh...",
  "redirect_uri": "https://hailango.com/auth/callback"
}
```

### 2.4 Refresh Token

Obtain new access token using refresh token.

```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g..."
}
```

### 2.5 Logout

Invalidate current session.

```http
POST /api/auth/logout
Authorization: Bearer <access_token>
```

### 2.6 Password Reset

Request password reset email.

```http
POST /api/auth/password-reset
Content-Type: application/json

{
  "email": "user@example.com"
}
```

---

## 3. User Endpoints

### 3.1 Get Current User

```http
GET /api/users/me
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "display_name": "John Doe",
    "native_language": "ja",
    "avatar_url": "https://...",
    "email_verified": true,
    "subscription": {
      "plan_type": "premium_monthly",
      "status": "active",
      "current_period_end": "2025-02-15T00:00:00Z"
    },
    "stats": {
      "total_books": 5,
      "total_study_time_minutes": 1250,
      "current_streak_days": 7
    },
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

### 3.2 Update User Profile

```http
PATCH /api/users/me
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "display_name": "Jane Doe",
  "native_language": "en"
}
```

### 3.3 Delete Account

```http
DELETE /api/users/me
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "password": "currentpassword123"
}
```

---

## 4. Book Management API

### 4.1 List Books

```http
GET /api/books
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `per_page` | integer | 20 | Items per page (max 100) |
| `status` | string | - | Filter by status |
| `language` | string | - | Filter by target language |

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "book-uuid-1",
      "title": "Kurdish for Beginners",
      "source_language": "en",
      "target_language": "ku",
      "total_pages": 120,
      "status": "ready",
      "progress": {
        "completed_pages": 45,
        "percentage": 37.5
      },
      "created_at": "2025-01-10T08:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 5
  }
}
```

### 4.2 Upload Book

```http
POST /api/books/upload
Authorization: Bearer <access_token>
Content-Type: multipart/form-data

file: <PDF or image file>
title: "My Language Book"
source_language: "en"
target_language: "ja"
reference_language: null
```

**Response (202 Accepted):**
```json
{
  "data": {
    "id": "book-uuid-new",
    "title": "My Language Book",
    "status": "pending",
    "job_id": "ocr-job-uuid"
  }
}
```

### 4.3 Get Book Details

```http
GET /api/books/{book_id}
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "data": {
    "id": "book-uuid-1",
    "title": "Kurdish for Beginners",
    "source_language": "en",
    "target_language": "ku",
    "reference_language": null,
    "total_pages": 120,
    "status": "ready",
    "settings": {
      "tts_language": "ku",
      "tts_speed": 1.0,
      "auto_play": true
    },
    "progress": {
      "completed_pages": 45,
      "last_page": 46,
      "percentage": 37.5,
      "total_study_time_minutes": 320
    },
    "created_at": "2025-01-10T08:00:00Z",
    "updated_at": "2025-01-14T15:30:00Z"
  }
}
```

### 4.4 Update Book Settings

```http
PATCH /api/books/{book_id}
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "title": "Kurdish - Updated Title",
  "settings": {
    "tts_speed": 0.8,
    "auto_play": false
  }
}
```

### 4.5 Delete Book

```http
DELETE /api/books/{book_id}
Authorization: Bearer <access_token>
```

### 4.6 Get Book Pages

```http
GET /api/books/{book_id}/pages
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `per_page` | integer | 20 | Items per page |

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "page-uuid-1",
      "page_number": 1,
      "content_preview": "Chapter 1: Greetings...",
      "has_audio": true,
      "is_completed": true
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 120
  }
}
```

### 4.7 Get Single Page

```http
GET /api/books/{book_id}/pages/{page_number}
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "data": {
    "id": "page-uuid-1",
    "book_id": "book-uuid-1",
    "page_number": 1,
    "original_content": "سڵاو! چۆنی؟",
    "processed_content": {
      "sentences": [
        {
          "text": "سڵاو!",
          "translation": "Hello!",
          "phonetic": "Slav!"
        },
        {
          "text": "چۆنی؟",
          "translation": "How are you?",
          "phonetic": "Choni?"
        }
      ]
    },
    "layout_data": {
      "bounding_boxes": [...]
    },
    "audio_url": "/api/tts/pages/page-uuid-1",
    "vocabulary": [
      {
        "id": "vocab-1",
        "word": "سڵاو",
        "meaning": "hello",
        "part_of_speech": "interjection"
      }
    ]
  }
}
```

### 4.8 Check OCR Status

```http
GET /api/books/{book_id}/status
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "data": {
    "status": "processing",
    "progress": {
      "processed_pages": 45,
      "total_pages": 120,
      "percentage": 37.5
    },
    "estimated_completion": "2025-01-14T16:00:00Z"
  }
}
```

---

## 5. Learning API

### 5.1 Start Learning Session

```http
POST /api/learning/sessions
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "book_id": "book-uuid-1",
  "session_type": "page_by_page",
  "start_page": 46,
  "settings": {
    "tts_speed": 1.0,
    "include_translation": true
  }
}
```

**Response (201 Created):**
```json
{
  "data": {
    "id": "session-uuid-1",
    "book_id": "book-uuid-1",
    "session_type": "page_by_page",
    "start_page": 46,
    "status": "active",
    "started_at": "2025-01-14T15:30:00Z"
  }
}
```

### 5.2 Get Active Session

```http
GET /api/learning/sessions/active
Authorization: Bearer <access_token>
```

### 5.3 Update Session Progress

```http
POST /api/learning/sessions/{session_id}/progress
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "page_id": "page-uuid-46",
  "time_spent_seconds": 180,
  "pronunciation_score": 85,
  "feedback_data": {
    "problematic_words": ["سڵاو"]
  }
}
```

### 5.4 End Session

```http
POST /api/learning/sessions/{session_id}/end
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "end_page": 50
}
```

### 5.5 Get Learning Statistics

```http
GET /api/learning/stats
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `period` | string | week | 'day', 'week', 'month', 'year' |
| `book_id` | uuid | - | Filter by specific book |

**Response (200 OK):**
```json
{
  "data": {
    "total_study_time_minutes": 1250,
    "average_daily_minutes": 45,
    "current_streak_days": 7,
    "longest_streak_days": 14,
    "pages_completed": 245,
    "vocabulary_learned": 892,
    "average_pronunciation_score": 78,
    "daily_breakdown": [
      {"date": "2025-01-14", "minutes": 45, "pages": 5},
      {"date": "2025-01-13", "minutes": 30, "pages": 3}
    ]
  }
}
```

---

## 6. TTS API

### 6.1 Generate Page Audio

```http
GET /api/tts/pages/{page_id}
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `speed` | float | 1.0 | Playback speed (0.5-2.0) |
| `quality` | string | standard | 'standard' or 'premium' |

**Response:** Audio stream (`audio/mpeg`)

### 6.2 Generate Text Audio

```http
POST /api/tts/synthesize
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "text": "سڵاو! چۆنی؟",
  "language": "ku",
  "speed": 0.8,
  "quality": "premium"
}
```

**Response:** Audio stream (`audio/mpeg`)

### 6.3 Batch Download (Teacher Mode)

```http
POST /api/tts/batch
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "book_id": "book-uuid-1",
  "start_page": 1,
  "end_page": 50,
  "settings": {
    "speed": 1.0,
    "include_translation": true
  }
}
```

**Response (202 Accepted):**
```json
{
  "data": {
    "job_id": "tts-batch-uuid",
    "estimated_size_mb": 125,
    "estimated_duration_minutes": 5
  }
}
```

---

## 7. STT API

### 7.1 Evaluate Pronunciation

```http
POST /api/stt/evaluate
Authorization: Bearer <access_token>
Content-Type: multipart/form-data

audio: <audio blob>
reference_text: "سڵاو! چۆنی؟"
language: "ku"
```

**Response (200 OK):**
```json
{
  "data": {
    "transcription": "سڵاو چۆنی",
    "overall_score": 85,
    "word_scores": [
      {
        "word": "سڵاو",
        "expected": "سڵاو",
        "actual": "سڵاو",
        "score": 92,
        "timing": {"start": 0.0, "end": 0.8}
      },
      {
        "word": "چۆنی",
        "expected": "چۆنی؟",
        "actual": "چۆنی",
        "score": 78,
        "timing": {"start": 0.9, "end": 1.5},
        "feedback": "Intonation should rise for question"
      }
    ],
    "suggestions": [
      "Practice the rising intonation for questions",
      "Good pronunciation of 'سڵاو'"
    ]
  }
}
```

---

## 8. SRS Review API

### 8.1 Get Due Reviews

```http
GET /api/review/due
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | integer | 20 | Max items to return |
| `book_id` | uuid | - | Filter by book |

**Response (200 OK):**
```json
{
  "data": {
    "due_count": 15,
    "items": [
      {
        "id": "srs-uuid-1",
        "vocabulary": {
          "id": "vocab-uuid-1",
          "word": "سڵاو",
          "meaning": "hello",
          "example_sentence": "سڵاو! چۆنی؟"
        },
        "interval_days": 3,
        "easiness_factor": 2.5,
        "repetitions": 2
      }
    ]
  }
}
```

### 8.2 Submit Review Result

```http
POST /api/review/{srs_id}/submit
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "quality": 4
}
```

**Quality Scale (SM-2):**
- `0` - Complete blackout
- `1` - Incorrect, remembered after seeing answer
- `2` - Incorrect, easy to recall after seeing
- `3` - Correct with difficulty
- `4` - Correct with hesitation
- `5` - Perfect recall

**Response (200 OK):**
```json
{
  "data": {
    "id": "srs-uuid-1",
    "next_review_date": "2025-01-21",
    "new_interval_days": 7,
    "new_easiness_factor": 2.6
  }
}
```

### 8.3 Add Vocabulary to SRS

```http
POST /api/review/add
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "vocabulary_id": "vocab-uuid-new"
}
```

---

## 9. WebSocket API (Teacher Mode)

### 9.1 Connection

```
ws://localhost:8080/ws/teacher/{book_id}
Authorization: Bearer <access_token>
```

### 9.2 Server Events

**PageChange** - New page started:
```json
{
  "type": "page_change",
  "page_index": 46,
  "page_id": "page-uuid-46",
  "content": {
    "original": "سڵاو! چۆنی؟",
    "translation": "Hello! How are you?"
  }
}
```

**AudioChunk** - Audio data:
```json
{
  "type": "audio_chunk",
  "page_index": 46,
  "chunk_index": 0,
  "data": "<base64 encoded audio>",
  "is_last": false
}
```

**SessionEnd** - Session completed:
```json
{
  "type": "session_end",
  "completed_pages": 10,
  "total_duration_seconds": 1800
}
```

**Error**:
```json
{
  "type": "error",
  "code": "TTS_FAILED",
  "message": "Failed to generate audio for page 47"
}
```

### 9.3 Client Commands

**Pause**:
```json
{
  "command": "pause"
}
```

**Resume**:
```json
{
  "command": "resume"
}
```

**Skip** - Skip to specific page:
```json
{
  "command": "skip",
  "page_index": 50
}
```

**UpdateSettings**:
```json
{
  "command": "update_settings",
  "settings": {
    "tts_speed": 0.8,
    "page_interval": 10
  }
}
```

**Stop**:
```json
{
  "command": "stop"
}
```

---

## 10. Rate Limiting

### Limits by Plan

| Endpoint Category | Free | Premium |
|-------------------|------|---------|
| Auth endpoints | 10/min | 10/min |
| Book upload | 5/day | 50/day |
| OCR processing | 100 pages/day | 1000 pages/day |
| TTS synthesis | 30 min/day | Unlimited |
| STT evaluation | 50/day | 500/day |
| API calls (general) | 1000/hour | 10000/hour |

### Rate Limit Headers

```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 950
X-RateLimit-Reset: 1705250400
```

### Rate Limit Error

**Response (429 Too Many Requests):**
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests",
    "retry_after": 60
  }
}
```

---

## 11. Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Missing or invalid token |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 422 | Invalid request data |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `QUOTA_EXCEEDED` | 403 | Plan quota exceeded |
| `OCR_FAILED` | 500 | OCR processing error |
| `TTS_FAILED` | 500 | TTS generation error |
| `STT_FAILED` | 500 | STT processing error |
| `PAYMENT_REQUIRED` | 402 | Premium feature requires subscription |

---

## 12. OpenAPI Specification

The complete OpenAPI 3.0 specification is available at:

- **Development**: `http://localhost:8080/api/openapi.json`
- **Swagger UI**: `http://localhost:8080/api/docs`
- **ReDoc**: `http://localhost:8080/api/redoc`

---

## References

- [System Architecture](system_architecture.md)
- [Database Schema](database_schema.md)
- [Requirements Definition](../requirements_definition.md)
- [reinhardt-rest Documentation](https://docs.rs/reinhardt-rest)
