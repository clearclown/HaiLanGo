# HaiLanGo - Design Specification

**Project**: HaiLanGo AI Language Learning Platform
**Generated**: 2026-02-02
**Source**: docs/architecture/*.md
**Framework**: Reinhardt (Rust Full-Stack)

---

## 1. System Architecture Overview

### 1.1 Architecture Principles

1. **Let LLM Handle It**: Pass natural language to AI; don't hardcode templates
2. **Privacy First**: E2E encryption, no external sharing
3. **Offline Capable**: PWA with service worker caching
4. **Modular Design**: Composable Reinhardt crates
5. **API-Driven**: External services abstracted behind traits

### 1.2 Technology Decisions

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| **Backend** | Rust + Reinhardt | Type safety, performance, full-stack |
| **ORM** | reinhardt-db (SeaQuery + sqlx) | Type-safe queries, migrations |
| **API** | reinhardt-rest | Django-style ViewSets, auto-serialization |
| **Frontend** | reinhardt-pages (WASM + SSR) | Reactive UI, Leptos-style components |
| **Auth** | reinhardt-auth | JWT + Session, OAuth support |
| **WebSocket** | reinhardt-websockets | Real-time Teacher Mode |
| **Database** | PostgreSQL 16 | ACID, JSONB, full-text search |
| **Cache** | Redis 7 | Session storage, rate limiting |

---

## 2. Component Architecture

### 2.1 Layer Diagram

```
┌─────────────────────────────────────────────────────┐
│               Client Layer                          │
│  reinhardt-pages (WASM) + Service Worker (PWA)     │
└────────────────┬────────────────────────────────────┘
                 │ HTTP/WebSocket
┌────────────────┴────────────────────────────────────┐
│               API Gateway                           │
│  reinhardt-rest (REST) + reinhardt-websockets      │
│  reinhardt-auth (JWT Middleware)                   │
└────────────────┬────────────────────────────────────┘
                 │
┌────────────────┴────────────────────────────────────┐
│            Application Layer                        │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐     │
│  │  auth  │ │ books  │ │learning│ │  tts   │     │
│  └────────┘ └────────┘ └────────┘ └────────┘     │
│  ┌────────┐ ┌────────┐ ┌────────┐                 │
│  │  stt   │ │ review │ │teacher │                 │
│  └────────┘ └────────┘ └────────┘                 │
└────────────────┬────────────────────────────────────┘
                 │
┌────────────────┴────────────────────────────────────┐
│              Data Layer                             │
│  PostgreSQL (reinhardt-db) + Redis                 │
└─────────────────────────────────────────────────────┘
                 │
┌────────────────┴────────────────────────────────────┐
│           External Services                         │
│  OCR │ TTS │ STT │ LLM │ Stripe                    │
└─────────────────────────────────────────────────────┘
```

---

## 3. Module Structure (Rust 2024 Edition)

### 3.1 Directory Layout

```
src/
├── main.rs                         # App entry point
├── lib.rs                          # Library root
├── config/
│   ├── module.rs                   # Config module (NOT mod.rs)
│   ├── settings/
│   │   ├── module.rs
│   │   ├── base.rs                 # Base settings (TOML)
│   │   ├── development.rs          # Dev overrides
│   │   └── production.rs           # Prod overrides
│   ├── urls.rs                     # URL routing
│   └── apps.rs                     # App registration
├── apps/
│   ├── auth/
│   │   ├── module.rs               # ✅ Always module.rs
│   │   ├── models.rs               # User, Session
│   │   ├── serializers.rs          # UserSerializer
│   │   ├── viewsets.rs             # AuthViewSet
│   │   ├── middleware.rs           # JWT middleware
│   │   ├── oauth.rs                # Google OAuth
│   │   └── tests.rs
│   ├── books/
│   │   ├── module.rs
│   │   ├── models.rs               # Book, Page
│   │   ├── serializers.rs
│   │   ├── viewsets.rs             # BookViewSet
│   │   ├── ocr/
│   │   │   ├── module.rs
│   │   │   ├── provider.rs         # OcrProvider trait
│   │   │   ├── google.rs           # Google Vision
│   │   │   ├── azure.rs            # Azure Vision
│   │   │   └── mock.rs             # Test mock
│   │   ├── encryption.rs           # E2E encryption
│   │   └── tests.rs
│   ├── learning/
│   │   ├── module.rs
│   │   ├── models.rs               # LearningSession, Progress
│   │   ├── serializers.rs
│   │   ├── viewsets.rs
│   │   ├── stats.rs                # Analytics
│   │   └── tests.rs
│   ├── tts/
│   │   ├── module.rs
│   │   ├── viewsets.rs
│   │   ├── providers/
│   │   │   ├── module.rs
│   │   │   ├── trait.rs            # TtsProvider trait
│   │   │   ├── google.rs
│   │   │   ├── azure.rs
│   │   │   └── mock.rs
│   │   ├── batch.rs                # Batch generation
│   │   └── tests.rs
│   ├── stt/
│   │   ├── module.rs
│   │   ├── viewsets.rs
│   │   ├── providers/
│   │   │   ├── module.rs
│   │   │   ├── trait.rs            # SttProvider trait
│   │   │   ├── whisper.rs
│   │   │   ├── azure.rs
│   │   │   └── mock.rs
│   │   ├── evaluation.rs           # Scoring logic
│   │   ├── feedback.rs             # AI feedback
│   │   └── tests.rs
│   ├── review/
│   │   ├── module.rs
│   │   ├── models.rs               # Vocabulary, SrsSchedule
│   │   ├── serializers.rs
│   │   ├── viewsets.rs
│   │   ├── sm2.rs                  # SM-2 algorithm
│   │   └── tests.rs
│   └── teacher_mode/
│       ├── module.rs
│       ├── websocket.rs            # WS handler
│       ├── session.rs              # Session state
│       ├── commands.rs             # Command handlers
│       ├── events.rs               # Event types
│       └── tests.rs
├── pages/                          # WASM frontend
│   ├── module.rs
│   ├── components/
│   │   ├── module.rs
│   │   ├── layout.rs
│   │   ├── auth.rs
│   │   ├── books.rs
│   │   ├── learning.rs
│   │   ├── teacher.rs
│   │   └── review.rs
│   ├── routes.rs
│   └── state.rs
├── services/
│   ├── module.rs
│   ├── circuit_breaker.rs
│   ├── usage_tracker.rs
│   └── llm.rs
└── utils/
    ├── module.rs
    ├── crypto.rs
    └── validators.rs
```

### 3.2 Module Declaration (MANDATORY)

**CORRECT** ✅:
```rust
// src/apps/auth/module.rs
pub mod models;
pub mod serializers;
pub mod viewsets;
pub mod middleware;
pub mod oauth;

#[cfg(test)]
mod tests;
```

**INCORRECT** ❌:
```rust
// src/apps/auth/mod.rs  ← NEVER USE mod.rs
```

---

## 4. Database Design

### 4.1 Entity-Relationship Overview

```mermaid
erDiagram
    users ||--o{ books : owns
    users ||--o{ learning_sessions : participates
    users ||--o{ srs_schedules : has
    users ||--o{ subscriptions : has

    books ||--o{ pages : contains
    books ||--o{ learning_sessions : subject_of

    pages ||--o{ vocabularies : extracted_from
    pages ||--o{ learning_progress : tracked_in

    vocabularies ||--o{ srs_schedules : scheduled_for

    learning_sessions ||--o{ learning_progress : records
```

### 4.2 Core Models (reinhardt-db)

#### User Model
```rust
use reinhardt_db::prelude::*;

#[derive(Model, Debug, Clone)]
#[model(table_name = "users")]
pub struct User {
    #[pk]
    pub id: Uuid,

    #[unique]
    pub email: String,

    pub password_hash: Option<String>,  // None for OAuth users
    pub display_name: String,

    #[default("en")]
    pub native_language: String,

    pub avatar_url: Option<String>,
    pub oauth_provider: Option<String>,  // "google", "github"
    pub oauth_id: Option<String>,

    #[default(false)]
    pub email_verified: bool,

    #[auto_now_add]
    pub created_at: DateTime<Utc>,

    #[auto_now]
    pub updated_at: DateTime<Utc>,

    pub last_login_at: Option<DateTime<Utc>>,
}

impl User {
    pub fn verify_password(&self, password: &str) -> bool {
        match &self.password_hash {
            Some(hash) => {
                let parsed = PasswordHash::new(hash).expect("invalid hash");
                Argon2::default().verify_password(password.as_bytes(), &parsed).is_ok()
            }
            None => false,
        }
    }

    pub async fn find_by_email(
        conn: &DatabaseConnection,
        email: &str,
    ) -> Result<Option<Self>, DbError> {
        Self::query()
            .filter(UserColumn::Email.eq(email))
            .one(conn)
            .await
    }
}
```

#### Book Model
```rust
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum BookStatus {
    Pending,      // Upload received, OCR not started
    Processing,   // OCR in progress
    Ready,        // All pages processed
    Error,        // OCR failed
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct BookSettings {
    pub tts_language: Option<String>,
    pub tts_speed: Option<f32>,
    pub auto_play: Option<bool>,
}

#[derive(Model, Debug, Clone)]
#[model(table_name = "books")]
pub struct Book {
    #[pk]
    pub id: Uuid,

    #[foreign_key(User)]
    pub user_id: Uuid,

    pub title: String,
    pub source_language: String,    // ISO 639-1
    pub target_language: String,
    pub reference_language: Option<String>,

    #[default(0)]
    pub total_pages: i32,

    #[default(BookStatus::Pending)]
    pub status: BookStatus,

    pub encryption_key_hash: Option<String>,

    #[json]
    #[default(BookSettings::default())]
    pub settings: BookSettings,

    #[auto_now_add]
    pub created_at: DateTime<Utc>,

    #[auto_now]
    pub updated_at: DateTime<Utc>,
}

impl Book {
    pub async fn find_by_user(
        conn: &DatabaseConnection,
        user_id: Uuid,
    ) -> Result<Vec<Self>, DbError> {
        Self::query()
            .filter(BookColumn::UserId.eq(user_id))
            .order_by(BookColumn::CreatedAt, Order::Desc)
            .all(conn)
            .await
    }
}
```

#### SRS Schedule Model (SM-2 Algorithm)
```rust
#[derive(Model, Debug, Clone)]
#[model(table_name = "srs_schedules")]
pub struct SrsSchedule {
    #[pk]
    pub id: Uuid,

    #[foreign_key(User)]
    pub user_id: Uuid,

    #[foreign_key(Vocabulary)]
    pub vocabulary_id: Uuid,

    pub next_review_date: NaiveDate,

    #[default(1)]
    pub interval_days: i32,

    #[default(2.5)]
    pub easiness_factor: f32,

    #[default(0)]
    pub repetitions: i32,

    #[default(0)]
    pub correct_count: i32,

    #[default(0)]
    pub incorrect_count: i32,

    pub last_reviewed_at: Option<DateTime<Utc>>,

    #[auto_now_add]
    pub created_at: DateTime<Utc>,
}

impl SrsSchedule {
    /// SM-2 algorithm: quality 0-5 (0-2=fail, 3-5=pass)
    pub fn update_after_review(&mut self, quality: u8) {
        let q = quality.min(5) as f32;

        // Update easiness factor
        self.easiness_factor = (self.easiness_factor
            + (0.1 - (5.0 - q) * (0.08 + (5.0 - q) * 0.02)))
            .max(1.3);

        if quality >= 3 {
            // Correct answer
            self.correct_count += 1;
            self.repetitions += 1;

            self.interval_days = match self.repetitions {
                1 => 1,
                2 => 6,
                _ => (self.interval_days as f32 * self.easiness_factor).round() as i32,
            };
        } else {
            // Incorrect answer
            self.incorrect_count += 1;
            self.repetitions = 0;
            self.interval_days = 1;
        }

        self.next_review_date = Utc::now().date_naive() + Duration::days(self.interval_days as i64);
        self.last_reviewed_at = Some(Utc::now());
    }

    pub async fn find_due_reviews(
        conn: &DatabaseConnection,
        user_id: Uuid,
        limit: i64,
    ) -> Result<Vec<Self>, DbError> {
        Self::query()
            .filter(SrsScheduleColumn::UserId.eq(user_id))
            .filter(SrsScheduleColumn::NextReviewDate.lte(Utc::now().date_naive()))
            .order_by(SrsScheduleColumn::NextReviewDate, Order::Asc)
            .limit(limit)
            .all(conn)
            .await
    }
}
```

### 4.3 Index Strategy

```sql
-- High-frequency lookups
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_oauth ON users(oauth_provider, oauth_id) WHERE oauth_provider IS NOT NULL;

-- Book queries
CREATE INDEX idx_books_user_id ON books(user_id);
CREATE INDEX idx_books_created_at ON books(user_id, created_at DESC);
CREATE INDEX idx_books_status ON books(status);

-- SRS due reviews (partial index)
CREATE INDEX idx_srs_due_today ON srs_schedules(user_id, next_review_date)
    WHERE next_review_date <= CURRENT_DATE;

-- Page lookups
CREATE UNIQUE INDEX idx_pages_book_page ON pages(book_id, page_number);

-- JSONB indexes
CREATE INDEX idx_books_tts_lang ON books USING gin ((settings->'tts_language'));
```

---

## 5. API Design (reinhardt-rest)

### 5.1 ViewSet Pattern

```rust
use reinhardt_rest::prelude::*;

#[derive(Serialize, Deserialize)]
pub struct BookSerializer {
    pub id: Uuid,
    pub title: String,
    pub source_language: String,
    pub target_language: String,
    pub total_pages: i32,
    pub status: BookStatus,
    pub progress: BookProgress,
    pub created_at: DateTime<Utc>,
}

#[viewset]
impl BookViewSet {
    type Model = Book;
    type Serializer = BookSerializer;

    #[action(detail = false, methods = ["GET"])]
    async fn list(&self, request: Request) -> Response {
        let user = request.user()?;
        let books = Book::find_by_user(&self.db, user.id).await?;
        let serialized = books.into_iter()
            .map(|b| BookSerializer::from(b))
            .collect::<Vec<_>>();
        Response::json(serialized)
    }

    #[action(detail = false, methods = ["POST"])]
    async fn upload(&self, request: Request) -> Response {
        let user = request.user()?;
        let form = request.multipart().await?;

        // Queue OCR job
        let book = self.queue_ocr_job(user.id, form).await?;

        Response::accepted(json!({
            "id": book.id,
            "status": "pending",
            "job_id": book.ocr_job_id
        }))
    }

    #[action(detail = true, methods = ["GET"])]
    async fn retrieve(&self, request: Request, id: Uuid) -> Response {
        let user = request.user()?;
        let book = Book::find_by_id(&self.db, id).await?;

        // Permission check
        if book.user_id != user.id {
            return Response::forbidden();
        }

        Response::json(BookSerializer::from(book))
    }

    #[action(detail = true, methods = ["PATCH"])]
    async fn update(&self, request: Request, id: Uuid) -> Response {
        let user = request.user()?;
        let payload: UpdateBookRequest = request.json().await?;

        let mut book = Book::find_by_id(&self.db, id).await?;
        if book.user_id != user.id {
            return Response::forbidden();
        }

        book.apply_updates(payload);
        book.save(&self.db).await?;

        Response::json(BookSerializer::from(book))
    }

    #[action(detail = true, methods = ["DELETE"])]
    async fn destroy(&self, request: Request, id: Uuid) -> Response {
        let user = request.user()?;
        let book = Book::find_by_id(&self.db, id).await?;

        if book.user_id != user.id {
            return Response::forbidden();
        }

        book.delete(&self.db).await?;
        Response::no_content()
    }

    #[action(detail = true, methods = ["GET"], path = "pages")]
    async fn pages(&self, request: Request, id: Uuid) -> Response {
        // GET /api/books/{id}/pages
        let user = request.user()?;
        let pages = Page::find_by_book(&self.db, id).await?;
        Response::json_paginated(pages, request.pagination())
    }

    #[action(detail = true, methods = ["GET"], path = "status")]
    async fn ocr_status(&self, request: Request, id: Uuid) -> Response {
        // GET /api/books/{id}/status
        let user = request.user()?;
        let book = Book::find_by_id(&self.db, id).await?;

        Response::json(json!({
            "status": book.status,
            "progress": {
                "processed_pages": book.processed_pages,
                "total_pages": book.total_pages,
                "percentage": (book.processed_pages as f32 / book.total_pages as f32 * 100.0)
            }
        }))
    }
}
```

### 5.2 URL Configuration

```rust
// src/config/urls.rs
use reinhardt_rest::router::Router;

pub fn configure_routes() -> Router {
    Router::new()
        .namespace("/api", |api| {
            api
                .include("auth", auth::urls())
                .include("books", books::urls())
                .include("learning", learning::urls())
                .include("tts", tts::urls())
                .include("stt", stt::urls())
                .include("review", review::urls())
        })
        .websocket("/ws/teacher/{book_id}", teacher_mode::websocket_handler)
}

// src/apps/books/module.rs
use reinhardt_rest::routes;

#[routes]
pub fn urls() -> Router {
    Router::new()
        .viewset("", BookViewSet)
        .post("upload", upload_handler)
}
```

---

## 6. Authentication & Authorization (reinhardt-auth)

### 6.1 JWT Configuration

```rust
use reinhardt_auth::prelude::*;

#[derive(Clone)]
pub struct AuthConfig {
    pub jwt_secret: String,           // 256-bit secret
    pub token_expiry: Duration,       // 1 hour
    pub refresh_expiry: Duration,     // 30 days
}

pub fn configure_auth(app: &mut App, config: AuthConfig) {
    // JWT middleware for Bearer token auth
    app.middleware(JwtAuthMiddleware::new(config.clone()));

    // Session middleware for cookie-based auth (optional)
    app.middleware(SessionMiddleware::new(RedisBackend::new()));
}

// Protected route
#[get("/api/users/me")]
#[authenticated]
async fn get_current_user(user: AuthenticatedUser) -> Response {
    Response::json(UserSerializer::from(user))
}
```

### 6.2 OAuth Providers

```rust
// src/apps/auth/oauth.rs
use oauth2::{AuthorizationCode, CsrfToken, PkceCodeChallenge, RedirectUrl, Scope};

pub struct GoogleOAuthProvider {
    client_id: String,
    client_secret: String,
    redirect_uri: String,
}

impl GoogleOAuthProvider {
    pub fn authorization_url(&self) -> (String, CsrfToken) {
        // Generate OAuth authorization URL
    }

    pub async fn exchange_code(&self, code: &str) -> Result<GoogleUserInfo, OAuthError> {
        // Exchange authorization code for access token
        // Fetch user info from Google
    }
}

// ViewSet endpoint
#[post("/api/auth/oauth/google")]
pub async fn google_oauth_callback(
    db: DatabaseConnection,
    payload: OAuthCallbackRequest,
) -> Response {
    let provider = GoogleOAuthProvider::from_config();
    let user_info = provider.exchange_code(&payload.code).await?;

    // Find or create user
    let user = User::find_or_create_oauth(
        &db,
        "google",
        &user_info.id,
        &user_info.email,
        &user_info.name,
    ).await?;

    // Generate JWT
    let tokens = generate_tokens(&user)?;

    Response::json(json!({
        "user": UserSerializer::from(user),
        "tokens": tokens
    }))
}
```

### 6.3 Permission Checking

```rust
// Middleware-level
#[authenticated]
async fn handler(user: AuthenticatedUser) -> Response { }

// Method-level
impl BookViewSet {
    fn check_permission(&self, request: &Request, book: &Book) -> Result<(), PermissionError> {
        let user = request.user()?;
        if book.user_id != user.id {
            return Err(PermissionError::Forbidden);
        }
        Ok(())
    }
}

// Role-based (future)
#[authenticated(roles = ["admin"])]
async fn admin_only() -> Response { }
```

---

## 7. WebSocket Design (reinhardt-websockets)

### 7.1 Teacher Mode WebSocket

```rust
use reinhardt_websockets::prelude::*;

#[websocket("/ws/teacher/{book_id}")]
pub async fn teacher_mode_handler(
    ws: WebSocket,
    book_id: Uuid,
    user: AuthenticatedUser,
) -> Result<(), WsError> {
    let (tx, rx) = ws.split();

    let session = TeacherSession::new(book_id, user.id);

    // Handle incoming commands (pause, resume, skip, settings)
    let cmd_handler = tokio::spawn(handle_commands(rx, session.clone()));

    // Stream audio and page updates
    let stream_handler = tokio::spawn(stream_lesson(tx, session));

    tokio::select! {
        _ = cmd_handler => {},
        _ = stream_handler => {},
    }

    Ok(())
}

// Event types
#[derive(Serialize)]
#[serde(tag = "type")]
pub enum TeacherEvent {
    PageChange {
        page_index: i32,
        page_id: Uuid,
        content: PageContent,
    },
    AudioChunk {
        page_index: i32,
        chunk_index: i32,
        data: Vec<u8>,
        is_last: bool,
    },
    SessionEnd {
        completed_pages: i32,
        total_duration_seconds: i32,
    },
    Error {
        code: String,
        message: String,
    },
}

// Command types
#[derive(Deserialize)]
#[serde(tag = "command")]
pub enum TeacherCommand {
    Pause,
    Resume,
    Skip { page_index: i32 },
    UpdateSettings { settings: SessionSettings },
    Stop,
}

async fn stream_lesson(
    tx: SplitSink<WebSocket, Message>,
    session: TeacherSession,
) -> Result<(), WsError> {
    let pages = session.load_pages().await?;

    for (index, page) in pages.iter().enumerate() {
        // Send page change event
        tx.send(Message::Text(serde_json::to_string(&TeacherEvent::PageChange {
            page_index: index as i32,
            page_id: page.id,
            content: page.content.clone(),
        })?)).await?;

        // Generate and stream audio
        let audio_chunks = session.generate_audio(page).await?;
        for (chunk_idx, chunk) in audio_chunks.iter().enumerate() {
            tx.send(Message::Binary(chunk.clone())).await?;
        }

        // Wait for page interval
        tokio::time::sleep(Duration::from_secs(session.settings.page_interval as u64)).await;
    }

    // Send session end
    tx.send(Message::Text(serde_json::to_string(&TeacherEvent::SessionEnd {
        completed_pages: pages.len() as i32,
        total_duration_seconds: session.elapsed_seconds(),
    })?)).await?;

    Ok(())
}
```

---

## 8. External Service Integration

### 8.1 Trait-Based Provider Pattern

```rust
// OCR Provider Trait
#[async_trait]
pub trait OcrProvider: Send + Sync {
    async fn extract_text(&self, image: &[u8]) -> Result<OcrResult, OcrError>;
    async fn extract_text_pdf(&self, pdf: &[u8]) -> Result<Vec<OcrResult>, OcrError>;
}

// Production: Google Vision
pub struct GoogleVisionClient {
    api_key: String,
    http_client: reqwest::Client,
}

#[async_trait]
impl OcrProvider for GoogleVisionClient {
    async fn extract_text(&self, image: &[u8]) -> Result<OcrResult, OcrError> {
        let response = self.http_client
            .post("https://vision.googleapis.com/v1/images:annotate")
            .header("X-Goog-Api-Key", &self.api_key)
            .json(&json!({
                "requests": [{
                    "image": { "content": base64::encode(image) },
                    "features": [{ "type": "DOCUMENT_TEXT_DETECTION" }]
                }]
            }))
            .send()
            .await?;

        let result: GoogleVisionResponse = response.json().await?;
        Ok(OcrResult::from(result))
    }
}

// Test: Mock
pub struct MockOcrClient {
    pub fixed_response: Option<OcrResult>,
}

#[async_trait]
impl OcrProvider for MockOcrClient {
    async fn extract_text(&self, _image: &[u8]) -> Result<OcrResult, OcrError> {
        Ok(self.fixed_response.clone().unwrap_or_default())
    }
}
```

### 8.2 Circuit Breaker Pattern

```rust
use std::sync::atomic::{AtomicU32, AtomicUsize, Ordering};
use std::time::{Duration, Instant};

pub struct CircuitBreaker {
    failure_count: AtomicU32,
    threshold: u32,
    timeout: Duration,
    last_failure: Arc<Mutex<Option<Instant>>>,
}

impl CircuitBreaker {
    pub async fn call<F, T, E>(&self, f: F) -> Result<T, CircuitError<E>>
    where
        F: FnOnce() -> Result<T, E>,
    {
        // Check if circuit is open
        if self.is_open() {
            return Err(CircuitError::Open);
        }

        match f() {
            Ok(result) => {
                self.reset();
                Ok(result)
            }
            Err(error) => {
                self.record_failure();
                Err(CircuitError::Inner(error))
            }
        }
    }

    fn is_open(&self) -> bool {
        let failures = self.failure_count.load(Ordering::Relaxed);
        if failures >= self.threshold {
            if let Some(last) = *self.last_failure.lock().unwrap() {
                return last.elapsed() < self.timeout;
            }
        }
        false
    }

    fn record_failure(&self) {
        self.failure_count.fetch_add(1, Ordering::Relaxed);
        *self.last_failure.lock().unwrap() = Some(Instant::now());
    }

    fn reset(&self) {
        self.failure_count.store(0, Ordering::Relaxed);
    }
}

// Usage
pub struct ResilientOcrClient {
    provider: Arc<dyn OcrProvider>,
    circuit: CircuitBreaker,
}

impl ResilientOcrClient {
    pub async fn extract_text(&self, image: &[u8]) -> Result<OcrResult, OcrError> {
        self.circuit
            .call(|| self.provider.extract_text(image))
            .await
            .map_err(|e| match e {
                CircuitError::Open => OcrError::ServiceUnavailable,
                CircuitError::Inner(e) => e,
            })
    }
}
```

### 8.3 Usage Tracking (Rate Limiting)

```rust
#[derive(Clone)]
pub struct UsageTracker {
    redis: RedisPool,
}

impl UsageTracker {
    pub async fn check_quota(
        &self,
        user_id: Uuid,
        service: ApiService,
        units: u32,
    ) -> Result<(), QuotaError> {
        let key = format!("quota:{}:{}:{}", user_id, service, today());
        let current: u32 = self.redis.get(&key).await.unwrap_or(0);

        let limit = self.get_limit(user_id, service).await;

        if current + units > limit {
            return Err(QuotaError::Exceeded { current, limit });
        }

        // Increment usage
        self.redis.incr(&key, units).await?;
        self.redis.expire(&key, 25 * 3600).await?;  // 25 hour TTL

        Ok(())
    }

    async fn get_limit(&self, user_id: Uuid, service: ApiService) -> u32 {
        let sub = Subscription::find_by_user(&self.db, user_id).await;

        match (service, sub.plan_type) {
            (ApiService::Ocr, PlanType::Free) => 100,           // 100 pages/day
            (ApiService::Ocr, PlanType::Premium) => 1000,       // 1000 pages/day
            (ApiService::TtsMinutes, PlanType::Free) => 30,     // 30 min/day
            (ApiService::TtsMinutes, PlanType::Premium) => u32::MAX,  // Unlimited
            (ApiService::SttEvaluations, PlanType::Free) => 50,
            (ApiService::SttEvaluations, PlanType::Premium) => 500,
            _ => 0,
        }
    }
}

// Middleware integration
#[inject]
async fn upload_handler(
    usage: UsageTracker,
    user: AuthenticatedUser,
    payload: UploadPayload,
) -> Response {
    // Check quota before processing
    usage.check_quota(user.id, ApiService::Ocr, 1).await?;

    // Process upload
}
```

---

## 9. Frontend Architecture (reinhardt-pages)

### 9.1 Component Pattern

```rust
use reinhardt_pages::prelude::*;

#[component]
pub fn LearningPage(book_id: Uuid) -> impl IntoView {
    let (page_index, set_page_index) = create_signal(0);
    let (is_playing, set_is_playing) = create_signal(false);

    let pages = create_resource(
        move || book_id,
        |id| async move {
            fetch_pages(id).await
        }
    );

    let play_audio = move |_| {
        set_is_playing(true);
        spawn_local(async move {
            let audio = fetch_audio(page_index()).await;
            audio.play();
        });
    };

    view! {
        <div class="learning-container">
            <Suspense fallback=|| view! { <LoadingSpinner/> }>
                {move || pages.get().map(|p| view! {
                    <PageViewer
                        page=p[page_index()]
                        on_play=play_audio
                        on_next=move |_| set_page_index.update(|i| *i += 1)
                        on_prev=move |_| set_page_index.update(|i| *i = i.saturating_sub(1))
                    />
                })}
            </Suspense>
        </div>
    }
}

#[component]
fn PageViewer<F>(
    page: Page,
    on_play: F,
    on_next: F,
    on_prev: F,
) -> impl IntoView
where
    F: Fn(MouseEvent) + 'static,
{
    view! {
        <div class="page-viewer">
            <div class="page-content">
                {page.processed_content}
            </div>
            <div class="controls">
                <button on:click=on_prev>"Previous"</button>
                <button on:click=on_play>"Play"</button>
                <button on:click=on_next>"Next"</button>
            </div>
        </div>
    }
}
```

### 9.2 State Management

```rust
// Global state
#[derive(Clone)]
pub struct AppState {
    pub auth: AuthState,
    pub books: BooksState,
    pub learning: LearningState,
}

#[derive(Clone)]
pub struct AuthState {
    pub user: Signal<Option<User>>,
    pub tokens: Signal<Option<TokenPair>>,
}

impl AppState {
    pub fn new() -> Self {
        Self {
            auth: AuthState::new(),
            books: BooksState::new(),
            learning: LearningState::new(),
        }
    }
}

// Provide to component tree
#[component]
pub fn App() -> impl IntoView {
    let state = create_rw_signal(AppState::new());

    view! {
        <Provider value=state>
            <Router>
                <Routes>
                    <Route path="/" view=HomePage/>
                    <Route path="/books" view=BooksPage/>
                    <Route path="/books/:id" view=LearningPage/>
                    <Route path="/review" view=ReviewPage/>
                </Routes>
            </Router>
        </Provider>
    }
}

// Consume in child component
#[component]
fn BooksList() -> impl IntoView {
    let state = expect_context::<RwSignal<AppState>>();
    let books = create_resource(
        move || (),
        |_| async move {
            let token = state.with(|s| s.auth.tokens.get().unwrap());
            fetch_books(&token).await
        }
    );

    view! {
        <For
            each=move || books.get().unwrap_or_default()
            key=|book| book.id
            children=|book| view! { <BookCard book=book/> }
        />
    }
}
```

---

## 10. Deployment Architecture

### 10.1 Development (Podman/Docker Compose)

```yaml
# compose.yaml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=development
      - DATABASE_URL=postgresql://hailango:password@db:5432/hailango
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=${JWT_SECRET}
      - GOOGLE_CLOUD_VISION_API_KEY=${GOOGLE_CLOUD_VISION_API_KEY}
      - GOOGLE_CLOUD_TTS_API_KEY=${GOOGLE_CLOUD_TTS_API_KEY}
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
      - STRIPE_SECRET_KEY=${STRIPE_SECRET_KEY}
    depends_on:
      - db
      - redis
    volumes:
      - ./data:/app/data

  db:
    image: postgres:16
    volumes:
      - pgdata:/var/lib/postgresql/data
    environment:
      POSTGRES_DB: hailango
      POSTGRES_USER: hailango
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"

  redis:
    image: redis:7
    volumes:
      - redisdata:/data
    ports:
      - "6379:6379"

volumes:
  pgdata:
  redisdata:
```

### 10.2 Production (Future)

```mermaid
graph TB
    subgraph "Edge"
        CF[Cloudflare CDN]
    end

    subgraph "Compute"
        LB[Load Balancer]
        APP1[App Instance 1]
        APP2[App Instance 2]
        APP3[App Instance 3]
    end

    subgraph "Data"
        PG[(PostgreSQL Primary)]
        PG_R[(PostgreSQL Replica)]
        REDIS[(Redis Cluster)]
    end

    subgraph "Storage"
        R2[Cloudflare R2<br/>Audio/Images]
    end

    CF --> LB
    LB --> APP1
    LB --> APP2
    LB --> APP3

    APP1 --> PG
    APP2 --> PG
    APP3 --> PG
    PG --> PG_R

    APP1 --> REDIS
    APP2 --> REDIS
    APP3 --> REDIS

    APP1 --> R2
    APP2 --> R2
    APP3 --> R2
```

---

## 11. Testing Architecture

### 11.1 Test Structure

```rust
// tests/common/mod.rs
use reinhardt_test::prelude::*;
use testcontainers::{clients::Cli, images::{postgres::Postgres, redis::Redis}};

pub struct TestApp {
    pub app: App,
    pub db: DatabaseConnection,
    pub redis: RedisPool,
    pub ocr_mock: Arc<MockOcrClient>,
    pub tts_mock: Arc<MockTtsClient>,
    pub stt_mock: Arc<MockSttClient>,
}

impl TestApp {
    pub async fn new() -> Self {
        let docker = Cli::default();

        let postgres = docker.run(Postgres::default());
        let redis = docker.run(Redis::default());

        let db_url = format!(
            "postgresql://test:test@localhost:{}/hailango_test",
            postgres.get_host_port_ipv4(5432)
        );

        let redis_url = format!(
            "redis://localhost:{}",
            redis.get_host_port_ipv4(6379)
        );

        // Run migrations
        let pool = PgPool::connect(&db_url).await.unwrap();
        sqlx::migrate!("./migrations").run(&pool).await.unwrap();

        // Create app with mocks
        let ocr_mock = Arc::new(MockOcrClient::default());
        let tts_mock = Arc::new(MockTtsClient::default());
        let stt_mock = Arc::new(MockSttClient::default());

        let app = App::builder()
            .database_url(&db_url)
            .redis_url(&redis_url)
            .ocr_provider(ocr_mock.clone())
            .tts_provider(tts_mock.clone())
            .stt_provider(stt_mock.clone())
            .build()
            .await
            .unwrap();

        Self {
            app,
            db: pool,
            redis: RedisPool::new(&redis_url).await.unwrap(),
            ocr_mock,
            tts_mock,
            stt_mock,
        }
    }

    pub async fn create_test_user(&self) -> User {
        User::create(&self.db, CreateUserRequest {
            email: format!("test_{}@example.com", Uuid::new_v4()),
            password: "TestPassword123!".to_string(),
            display_name: "Test User".to_string(),
            native_language: "en".to_string(),
        }).await.unwrap()
    }

    pub async fn authenticate(&self, user: &User) -> TokenPair {
        generate_tokens(user).unwrap()
    }
}
```

### 11.2 Test Example

```rust
// tests/integration/books_tests.rs
use crate::common::TestApp;

#[tokio::test]
async fn test_book_upload_and_ocr() {
    // Arrange
    let app = TestApp::new().await;
    let user = app.create_test_user().await;
    let auth = app.authenticate(&user).await;

    app.ocr_mock.set_response(OcrResult {
        text: "Page 1 content".to_string(),
        confidence: 0.95,
        layout: vec![],
    });

    // Act
    let response = app.app
        .post("/api/books/upload")
        .bearer_auth(&auth.access_token)
        .multipart(
            Form::new()
                .file("file", "tests/fixtures/pdfs/sample.pdf")
                .text("title", "Test Book")
                .text("source_language", "en")
                .text("target_language", "ja"),
        )
        .send()
        .await;

    // Assert
    assert_eq!(response.status(), StatusCode::ACCEPTED);

    let body: UploadResponse = response.json().await;
    assert_eq!(body.status, "pending");

    // Wait for OCR to complete (in real app this would be async)
    tokio::time::sleep(Duration::from_secs(1)).await;

    // Verify book was created
    let book = Book::find_by_id(&app.db, body.id).await.unwrap();
    assert_eq!(book.title, "Test Book");
    assert_eq!(book.status, BookStatus::Ready);
    assert_eq!(book.total_pages, 1);

    // Verify page was created
    let pages = Page::find_by_book(&app.db, book.id).await.unwrap();
    assert_eq!(pages.len(), 1);
    assert_eq!(pages[0].original_content.as_ref().unwrap(), "Page 1 content");
}

#[tokio::test]
#[serial(database)]
async fn test_srs_sm2_algorithm() {
    // Arrange
    let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());

    // Act: Perfect recall progression
    schedule.update_after_review(5);  // First review
    assert_eq!(schedule.interval_days, 1);
    assert_eq!(schedule.repetitions, 1);

    schedule.update_after_review(5);  // Second review
    assert_eq!(schedule.interval_days, 6);
    assert_eq!(schedule.repetitions, 2);

    schedule.update_after_review(5);  // Third review
    assert!(schedule.interval_days > 6);
    assert_eq!(schedule.repetitions, 3);

    // Act: Failed review resets
    schedule.update_after_review(1);
    assert_eq!(schedule.interval_days, 1);
    assert_eq!(schedule.repetitions, 0);

    // Assert: Easiness factor decreased
    assert!(schedule.easiness_factor < 2.5);
    assert!(schedule.easiness_factor >= 1.3);
}
```

---

## 12. Security Design

### 12.1 Encryption

```rust
use aes_gcm::{Aes256Gcm, Key, Nonce};
use argon2::{Argon2, PasswordHasher};

pub struct ContentEncryption {
    cipher: Aes256Gcm,
}

impl ContentEncryption {
    pub fn from_user_password(password: &str) -> Self {
        // Derive encryption key from password
        let key_bytes = Self::derive_key(password);
        let key = Key::<Aes256Gcm>::from_slice(&key_bytes);

        Self {
            cipher: Aes256Gcm::new(key),
        }
    }

    fn derive_key(password: &str) -> [u8; 32] {
        let salt = b"hailango_encryption_salt_v1";  // Should be per-user
        let mut key = [0u8; 32];
        pbkdf2::pbkdf2::<Hmac<Sha256>>(password.as_bytes(), salt, 100_000, &mut key);
        key
    }

    pub fn encrypt(&self, plaintext: &[u8]) -> Vec<u8> {
        let nonce = Nonce::from_slice(b"unique_nonce");  // Should be random
        self.cipher.encrypt(nonce, plaintext).expect("encryption failed")
    }

    pub fn decrypt(&self, ciphertext: &[u8]) -> Vec<u8> {
        let nonce = Nonce::from_slice(b"unique_nonce");
        self.cipher.decrypt(nonce, ciphertext).expect("decryption failed")
    }
}
```

### 12.2 Rate Limiting

```rust
// Middleware
pub struct RateLimitMiddleware {
    redis: RedisPool,
    limits: HashMap<String, RateLimit>,
}

impl Middleware for RateLimitMiddleware {
    async fn handle(&self, req: Request, next: Next) -> Response {
        let key = format!("ratelimit:{}:{}", req.path(), req.ip());
        let limit = self.limits.get(req.path()).unwrap_or(&RateLimit::default());

        let current: u32 = self.redis.get(&key).await.unwrap_or(0);

        if current >= limit.max_requests {
            return Response::too_many_requests(json!({
                "error": "RATE_LIMIT_EXCEEDED",
                "retry_after": limit.window_seconds
            }));
        }

        self.redis.incr(&key, 1).await;
        self.redis.expire(&key, limit.window_seconds).await;

        let mut response = next.run(req).await;
        response.headers_mut().insert(
            "X-RateLimit-Limit",
            limit.max_requests.to_string().parse().unwrap(),
        );
        response.headers_mut().insert(
            "X-RateLimit-Remaining",
            (limit.max_requests - current - 1).to_string().parse().unwrap(),
        );

        response
    }
}
```

---

## References
- [system_architecture.md](../docs/architecture/system_architecture.md)
- [database_schema.md](../docs/architecture/database_schema.md)
- [api_specification.md](../docs/architecture/api_specification.md)
- [test_strategy.md](../docs/testing/test_strategy.md)
- [CLAUDE.md](../CLAUDE.md)
