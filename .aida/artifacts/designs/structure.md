# Phase 2: Structure Definition

## Execution Date
2026-02-02

## Directory Structure (Reinhardt Conventions)

```
HaiLanGo/
├── Cargo.toml                      # Workspace manifest
├── .aida/                          # AIDA pipeline artifacts
│   ├── artifacts/
│   ├── specs/
│   ├── results/
│   └── state/
├── src/
│   ├── main.rs                     # Application entry point
│   ├── lib.rs                      # Library root
│   ├── config/
│   │   ├── module.rs               # Config module (NOT mod.rs)
│   │   ├── settings/
│   │   │   ├── module.rs
│   │   │   ├── base.rs             # Base settings (TOML)
│   │   │   ├── development.rs      # Dev-specific settings
│   │   │   └── production.rs       # Prod-specific settings
│   │   ├── urls.rs                 # URL routing configuration
│   │   └── apps.rs                 # App registration (installed_apps!)
│   ├── apps/
│   │   ├── auth/
│   │   │   ├── module.rs           # Auth app module
│   │   │   ├── models.rs           # User, Session models
│   │   │   ├── serializers.rs      # User serializers
│   │   │   ├── viewsets.rs         # AuthViewSet (REST)
│   │   │   ├── middleware.rs       # JWT/Session middleware
│   │   │   ├── oauth.rs            # OAuth providers (Google)
│   │   │   └── tests.rs            # Auth tests
│   │   ├── books/
│   │   │   ├── module.rs
│   │   │   ├── models.rs           # Book, Page models
│   │   │   ├── serializers.rs      # Book serializers
│   │   │   ├── viewsets.rs         # BookViewSet
│   │   │   ├── ocr/
│   │   │   │   ├── module.rs
│   │   │   │   ├── provider.rs     # OcrProvider trait
│   │   │   │   ├── google.rs       # Google Vision client
│   │   │   │   ├── azure.rs        # Azure Vision client
│   │   │   │   └── mock.rs         # Mock OCR (for tests)
│   │   │   ├── encryption.rs       # E2E encryption for books
│   │   │   └── tests.rs
│   │   ├── learning/
│   │   │   ├── module.rs
│   │   │   ├── models.rs           # LearningSession, Progress
│   │   │   ├── serializers.rs
│   │   │   ├── viewsets.rs         # SessionViewSet
│   │   │   ├── stats.rs            # Learning statistics
│   │   │   └── tests.rs
│   │   ├── tts/
│   │   │   ├── module.rs
│   │   │   ├── models.rs           # Audio cache models
│   │   │   ├── viewsets.rs         # TtsViewSet
│   │   │   ├── providers/
│   │   │   │   ├── module.rs
│   │   │   │   ├── trait.rs        # TtsProvider trait
│   │   │   │   ├── google.rs       # Google Cloud TTS
│   │   │   │   ├── azure.rs        # Azure Speech
│   │   │   │   └── mock.rs
│   │   │   ├── batch.rs            # Batch audio generation
│   │   │   └── tests.rs
│   │   ├── stt/
│   │   │   ├── module.rs
│   │   │   ├── viewsets.rs         # SttViewSet
│   │   │   ├── providers/
│   │   │   │   ├── module.rs
│   │   │   │   ├── trait.rs        # SttProvider trait
│   │   │   │   ├── whisper.rs      # OpenAI Whisper
│   │   │   │   ├── azure.rs        # Azure Speech
│   │   │   │   └── mock.rs
│   │   │   ├── evaluation.rs       # Pronunciation scoring
│   │   │   ├── feedback.rs         # AI feedback generation
│   │   │   └── tests.rs
│   │   ├── review/
│   │   │   ├── module.rs
│   │   │   ├── models.rs           # Vocabulary, SrsSchedule
│   │   │   ├── serializers.rs
│   │   │   ├── viewsets.rs         # ReviewViewSet
│   │   │   ├── sm2.rs              # SM-2 algorithm
│   │   │   └── tests.rs            # SRS algorithm tests
│   │   └── teacher_mode/
│   │       ├── module.rs
│   │       ├── websocket.rs        # WebSocket handler
│   │       ├── session.rs          # TeacherSession state
│   │       ├── commands.rs         # WS command handlers
│   │       ├── events.rs           # WS event types
│   │       └── tests.rs
│   ├── pages/                      # WASM frontend (reinhardt-pages)
│   │   ├── module.rs
│   │   ├── components/
│   │   │   ├── module.rs
│   │   │   ├── layout.rs           # Layout components
│   │   │   ├── auth.rs             # Login/Register forms
│   │   │   ├── books.rs            # Book list/detail
│   │   │   ├── learning.rs         # Learning page
│   │   │   ├── teacher.rs          # Teacher mode UI
│   │   │   └── review.rs           # SRS review UI
│   │   ├── routes.rs               # Client-side routing
│   │   └── state.rs                # Global state management
│   ├── services/                   # Shared services
│   │   ├── module.rs
│   │   ├── circuit_breaker.rs      # Circuit breaker pattern
│   │   ├── usage_tracker.rs        # API quota tracking
│   │   └── llm.rs                  # Claude API client
│   └── utils/                      # Utility functions
│       ├── module.rs
│       ├── crypto.rs               # Encryption utilities
│       └── validators.rs           # Input validation
├── migrations/                     # SQL migrations
│   ├── 0001_initial.sql
│   ├── 0002_srs_schedules.sql
│   └── ...
├── templates/                      # Server-side templates (if needed)
├── static/                         # Static files (CSS, JS, images)
├── tests/
│   ├── common/
│   │   ├── module.rs
│   │   ├── test_app.rs             # TestApp helper
│   │   └── factories.rs            # Test data factories
│   ├── fixtures/
│   │   ├── audio/
│   │   ├── images/
│   │   ├── pdfs/
│   │   └── json/
│   ├── integration/
│   │   ├── auth_tests.rs
│   │   ├── books_tests.rs
│   │   ├── learning_tests.rs
│   │   ├── tts_tests.rs
│   │   ├── stt_tests.rs
│   │   ├── review_tests.rs
│   │   └── teacher_mode_tests.rs
│   └── load/
│       └── basic_load.js           # k6 load tests
├── docs/                           # Documentation (already exists)
├── .github/
│   └── workflows/
│       └── ci.yml                  # GitHub Actions CI/CD
└── compose.yaml                    # Podman/Docker Compose

```

## reinhardt-db Models

### User Model
```rust
use reinhardt_db::prelude::*;

#[derive(Model, Debug, Clone)]
#[model(table_name = "users")]
pub struct User {
    #[pk]
    pub id: Uuid,

    #[unique]
    pub email: String,

    pub password_hash: Option<String>,
    pub display_name: String,

    #[default("en")]
    pub native_language: String,

    pub avatar_url: Option<String>,
    pub oauth_provider: Option<String>,
    pub oauth_id: Option<String>,

    #[default(false)]
    pub email_verified: bool,

    #[auto_now_add]
    pub created_at: DateTime<Utc>,

    #[auto_now]
    pub updated_at: DateTime<Utc>,

    pub last_login_at: Option<DateTime<Utc>>,
}
```

### Book Model
```rust
#[derive(Model, Debug, Clone)]
#[model(table_name = "books")]
pub struct Book {
    #[pk]
    pub id: Uuid,

    #[foreign_key(User)]
    pub user_id: Uuid,

    pub title: String,
    pub source_language: String,
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
```

### Page Model
```rust
#[derive(Model, Debug, Clone)]
#[model(table_name = "pages")]
pub struct Page {
    #[pk]
    pub id: Uuid,

    #[foreign_key(Book)]
    pub book_id: Uuid,

    pub page_number: i32,
    pub original_content: Option<String>,
    pub processed_content: Option<String>,

    #[json]
    pub layout_data: Option<LayoutData>,

    pub audio_url: Option<String>,

    #[default(false)]
    pub is_processed: bool,

    #[auto_now_add]
    pub created_at: DateTime<Utc>,
}
```

### SRS Schedule Model (SM-2)
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
```

### Learning Session Model
```rust
#[derive(Model, Debug, Clone)]
#[model(table_name = "learning_sessions")]
pub struct LearningSession {
    #[pk]
    pub id: Uuid,

    #[foreign_key(User)]
    pub user_id: Uuid,

    #[foreign_key(Book)]
    pub book_id: Option<Uuid>,

    pub session_type: SessionType,
    pub start_page: Option<i32>,
    pub end_page: Option<i32>,

    #[default(0)]
    pub duration_seconds: i32,

    #[json]
    #[default(SessionSettings::default())]
    pub settings: SessionSettings,

    #[default(SessionStatus::Active)]
    pub status: SessionStatus,

    #[auto_now_add]
    pub started_at: DateTime<Utc>,

    pub ended_at: Option<DateTime<Utc>>,
}
```

## reinhardt-rest API Structure

### ViewSet Pattern
```rust
use reinhardt_rest::prelude::*;

#[viewset]
impl BookViewSet {
    type Model = Book;
    type Serializer = BookSerializer;

    #[action(detail = false, methods = ["GET"])]
    async fn list(&self, request: Request) -> Response {
        // GET /api/books
    }

    #[action(detail = false, methods = ["POST"])]
    async fn create(&self, request: Request) -> Response {
        // POST /api/books/upload
    }

    #[action(detail = true, methods = ["GET"])]
    async fn retrieve(&self, request: Request, id: Uuid) -> Response {
        // GET /api/books/{id}
    }

    #[action(detail = true, methods = ["PATCH"])]
    async fn update(&self, request: Request, id: Uuid) -> Response {
        // PATCH /api/books/{id}
    }

    #[action(detail = true, methods = ["DELETE"])]
    async fn destroy(&self, request: Request, id: Uuid) -> Response {
        // DELETE /api/books/{id}
    }

    #[action(detail = true, methods = ["GET"], path = "pages")]
    async fn pages(&self, request: Request, id: Uuid) -> Response {
        // GET /api/books/{id}/pages
    }

    #[action(detail = true, methods = ["GET"], path = "status")]
    async fn status(&self, request: Request, id: Uuid) -> Response {
        // GET /api/books/{id}/status
    }
}
```

### URL Configuration
```rust
// src/config/urls.rs
use reinhardt_rest::router::Router;

pub fn configure_routes() -> Router {
    Router::new()
        .namespace("/api", |api| {
            api.include("auth", auth::urls())
                .include("books", books::urls())
                .include("learning", learning::urls())
                .include("tts", tts::urls())
                .include("stt", stt::urls())
                .include("review", review::urls())
        })
        .websocket("/ws/teacher/{book_id}", teacher_mode::websocket_handler)
}
```

## reinhardt-auth Configuration

### JWT Configuration
```rust
use reinhardt_auth::prelude::*;

#[derive(Clone)]
pub struct AuthConfig {
    pub jwt_secret: String,
    pub token_expiry: Duration,
    pub refresh_expiry: Duration,
}

pub fn configure_auth(app: &mut App, config: AuthConfig) {
    app.middleware(JwtAuthMiddleware::new(config.clone()));
    app.middleware(SessionMiddleware::new(RedisBackend::new()));
}
```

### OAuth Providers
```rust
pub struct OAuthConfig {
    pub google: GoogleOAuthProvider {
        client_id: String,
        client_secret: String,
        redirect_uri: String,
    },
}
```

## reinhardt-websockets Structure

### Teacher Mode WebSocket
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

    // Handle incoming commands (pause, resume, skip)
    let cmd_handler = spawn(handle_commands(rx, session.clone()));

    // Stream audio and page updates
    let stream_handler = spawn(stream_lesson(tx, session));

    tokio::select! {
        _ = cmd_handler => {},
        _ = stream_handler => {},
    }

    Ok(())
}
```

## External Service Abstraction

### Trait-Based Providers
```rust
// OCR Provider
#[async_trait]
pub trait OcrProvider: Send + Sync {
    async fn extract_text(&self, image: &[u8]) -> Result<OcrResult, OcrError>;
    async fn extract_text_pdf(&self, pdf: &[u8]) -> Result<Vec<OcrResult>, OcrError>;
}

// TTS Provider
#[async_trait]
pub trait TtsProvider: Send + Sync {
    async fn synthesize(&self, text: &str, language: &str, options: TtsOptions) -> Result<AudioData, TtsError>;
    async fn synthesize_batch(&self, requests: Vec<TtsRequest>) -> Result<Vec<AudioData>, TtsError>;
}

// STT Provider
#[async_trait]
pub trait SttProvider: Send + Sync {
    async fn transcribe(&self, audio: &[u8], language: &str) -> Result<Transcription, SttError>;
    async fn evaluate_pronunciation(&self, audio: &[u8], reference: &str, language: &str) -> Result<PronunciationScore, SttError>;
}
```

## Dependency Injection

```rust
use reinhardt_di::inject;

#[inject]
async fn process_book_upload(
    ocr: Arc<dyn OcrProvider>,
    db: DatabaseConnection,
    redis: RedisPool,
    user: AuthenticatedUser,
    payload: UploadPayload,
) -> Result<Book, ApiError> {
    // Dependencies automatically resolved by DI container
}
```

## Testing Structure

### Test App Helper
```rust
use reinhardt_test::prelude::*;

pub struct TestApp {
    app: App,
    db: DatabaseConnection,
    redis: RedisPool,
    ocr_mock: Arc<MockOcrClient>,
    tts_mock: Arc<MockTtsClient>,
    stt_mock: Arc<MockSttClient>,
}

impl TestApp {
    pub async fn new() -> Self {
        let infra = TestInfra::new().await;
        // Configure with mocks
    }

    pub async fn create_test_user(&self) -> User { /* ... */ }
    pub async fn authenticate(&self, user: &User) -> AuthTokens { /* ... */ }
}
```

## Module System (Rust 2024)

### ALWAYS Use module.rs (NOT mod.rs)
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

### Route Registration
```rust
use reinhardt_rest::routes;

#[routes]
pub fn urls() -> Router {
    Router::new()
        .viewset("", AuthViewSet)
        .route("oauth/google", oauth::google_login)
        .route("refresh", auth::refresh_token)
}
```

## Next Steps
Proceed to Phase 3: Alignment Verification
