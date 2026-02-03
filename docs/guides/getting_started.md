# Getting Started

This guide will help you set up the HaiLanGo development environment and run the application locally.

## 1. Prerequisites

### Required Software

| Software | Version | Installation |
|----------|---------|--------------|
| **Rust** | Nightly (1.85+) | [rustup.rs](https://rustup.rs) |
| **PostgreSQL** | 16+ | Package manager or container |
| **Redis** | 7+ | Package manager or container |
| **Podman** or Docker | Latest | Package manager |
| **podman-compose** | Latest | `pip install podman-compose` |

> **Note**: Reinhardt framework requires nightly Rust. The project includes `rust-toolchain.toml` which automatically selects the nightly toolchain.

### Verify Installation

```bash
# Check Rust toolchain
rustc --version    # Should be nightly-2024-xx-xx or later
cargo --version
rustup show        # Verify nightly is active

# Check container runtime
podman --version   # or docker --version
podman-compose --version

# Check database clients (optional, for debugging)
psql --version
redis-cli --version
```

### Rust Toolchain Setup

```bash
# Install Rust via rustup
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Add WASM target for frontend
rustup target add wasm32-unknown-unknown

# Install useful tools
cargo install cargo-watch    # Auto-rebuild on file changes
cargo install sqlx-cli       # Database migrations
cargo install trunk          # WASM bundler
```

---

## 2. Quick Start (5 Minutes)

### Step 1: Clone Repository

```bash
git clone https://github.com/clearclown/HaiLanGo.git
cd HaiLanGo
```

### Step 2: Start Infrastructure

```bash
# Start PostgreSQL and Redis
podman-compose up -d db redis

# Wait for databases to be ready
sleep 5

# Verify containers are running
podman ps
```

### Step 3: Configure Environment

```bash
# Copy example environment file
cp .env.example .env

# Edit configuration (use your preferred editor)
nano .env
```

Minimum required configuration:
```bash
# .env
APP_ENV=development
DATABASE_URL=postgresql://hailango:password@localhost:5432/hailango
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-development-secret-key-min-32-chars
```

### Step 4: Initialize Database

```bash
# Run migrations
cargo sqlx migrate run

# Verify database schema
cargo sqlx prepare --check
```

### Step 5: Run Application

```bash
# Development mode with auto-reload
cargo watch -x run

# Or standard run
cargo run
```

### Step 6: Verify Installation

```bash
# Check API health
curl http://localhost:8080/api/health

# Expected response:
# {"status":"ok","version":"0.1.0"}
```

Open your browser to `http://localhost:8080` to see the web interface.

---

## 3. Configuration Files

### 3.1 Environment Variables (.env)

```bash
# Application
APP_ENV=development          # development, staging, production
APP_HOST=0.0.0.0
APP_PORT=8080
APP_LOG_LEVEL=debug          # trace, debug, info, warn, error

# Database
DATABASE_URL=postgresql://hailango:password@localhost:5432/hailango
DATABASE_MAX_CONNECTIONS=10
DATABASE_TIMEOUT_SECONDS=30

# Redis
REDIS_URL=redis://localhost:6379
REDIS_MAX_CONNECTIONS=5

# Authentication
JWT_SECRET=your-secret-key-at-least-32-characters
JWT_ACCESS_EXPIRY=3600       # 1 hour
JWT_REFRESH_EXPIRY=604800    # 7 days

# OAuth (optional for development)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URI=http://localhost:8080/auth/callback/google

# External APIs (optional - uses mock in development)
GOOGLE_CLOUD_VISION_API_KEY=
GOOGLE_CLOUD_TTS_API_KEY=
OPENAI_API_KEY=
ANTHROPIC_API_KEY=
STRIPE_SECRET_KEY=
STRIPE_WEBHOOK_SECRET=

# Storage
STORAGE_PATH=./storage       # Local file storage path
MAX_UPLOAD_SIZE_MB=50

# Feature Flags
ENABLE_MOCK_APIS=true        # Use mock APIs in development
ENABLE_ADMIN_PANEL=true
ENABLE_SWAGGER_UI=true
```

### 3.2 Application Settings (config/settings/)

Settings are organized by environment in TOML format:

```
config/
├── settings/
│   ├── base.toml           # Shared settings
│   ├── development.toml    # Development overrides
│   ├── staging.toml        # Staging overrides
│   └── production.toml     # Production overrides
```

**base.toml:**
```toml
[server]
host = "0.0.0.0"
port = 8080
request_timeout_secs = 30

[database]
max_connections = 10
min_connections = 2
connect_timeout_secs = 10
idle_timeout_secs = 600

[auth]
password_min_length = 8
session_duration_secs = 3600
max_login_attempts = 5
lockout_duration_secs = 900

[tts]
default_speed = 1.0
min_speed = 0.5
max_speed = 2.0
cache_duration_secs = 86400

[srs]
initial_interval_days = 1
initial_easiness = 2.5
min_easiness = 1.3
```

**development.toml:**
```toml
[server]
enable_cors_all = true
enable_request_logging = true

[features]
mock_external_apis = true
enable_debug_endpoints = true
```

---

## 4. Mock Mode Development

When `ENABLE_MOCK_APIS=true`, external APIs are replaced with mock implementations:

### Mock OCR

- Returns predefined text based on filename patterns
- Simulates processing delay (1-3 seconds)
- Useful for testing upload flow without API costs

```rust
// Automatic mock when environment variable is set
let ocr_client: Arc<dyn OcrProvider> = if config.mock_apis {
    Arc::new(MockOcrClient::new())
} else {
    Arc::new(GoogleVisionClient::new(&config.google_api_key))
};
```

### Mock TTS

- Returns silent audio or pre-recorded samples
- Instant response for faster development
- Sample audio files in `tests/fixtures/audio/`

### Mock STT

- Returns configurable transcription results
- Simulates scoring based on input length
- Perfect for testing pronunciation flow

### Enabling/Disabling Mocks

```bash
# In .env file
ENABLE_MOCK_APIS=true   # Development
ENABLE_MOCK_APIS=false  # Use real APIs
```

Or per-service:
```bash
MOCK_OCR=true
MOCK_TTS=false
MOCK_STT=true
```

---

## 5. Available Services

### 5.1 Cache Service (Redis)

Redis-based caching for sessions, rate limiting, and job queues.

```rust
use hailango::services::CacheService;

// Initialize cache
let cache = CacheService::from_env().await?;

// Session management
cache.set_session("session_id", &user_data, Duration::from_secs(3600)).await?;
let session: Option<UserSession> = cache.get_session("session_id").await?;

// Rate limiting
let remaining = cache.check_rate_limit("user:123", 100, Duration::from_secs(60)).await?;

// Job queue
cache.push_job("ocr_jobs", &job_data).await?;
let job: Option<OcrJob> = cache.pop_job("ocr_jobs", Duration::from_secs(30)).await?;
```

### 5.2 WebSocket Service (Teacher Mode)

Real-time communication for automated lesson playback.

```rust
use hailango::services::{WsConnectionManager, PlaybackState};

let manager = WsConnectionManager::new();

// Start lesson session
let session = manager.start_lesson(user_id, book_id, total_pages).await;

// Control playback
manager.update_state(user_id, PlaybackState::Playing).await;
manager.next_page(user_id).await;

// Stream audio
manager.send_audio_chunk(user_id, base64_audio, false).await;

// End session
manager.end_session(user_id).await;
```

### 5.3 i18n Service (Internationalization)

Multi-language support for 7 languages.

```rust
use hailango::services::{Language, translate, keys};

// Get language from code
let lang = Language::from_code("ja").unwrap();

// Translate UI strings
let welcome = translate(lang, keys::WELCOME);  // "HaiLanGoへようこそ"
let login = translate(lang, keys::LOGIN);      // "ログイン"

// Supported languages
let all_langs = Language::all();  // EN, JA, ZH, KO, ES, FR, DE
```

**Supported Languages:**

| Code | Language |
|------|----------|
| `en` | English |
| `ja` | Japanese |
| `zh` | Chinese |
| `ko` | Korean |
| `es` | Spanish |
| `fr` | French |
| `de` | German |

---

## 6. Development Workflow

### 6.1 Running with Auto-Reload

```bash
# Watch for changes and rebuild
cargo watch -x run

# Watch with custom ignore patterns
cargo watch -x run -i "*.md" -i "tests/*"
```

### 6.2 Running Tests

```bash
# Run all tests
cargo test --workspace --all-features

# Run specific test
cargo test test_user_registration

# Run with output
cargo test -- --nocapture

# Run integration tests only
cargo test --test integration
```

### 6.3 Database Operations

```bash
# Create new migration
cargo sqlx migrate add <migration_name>

# Run pending migrations
cargo sqlx migrate run

# Rollback last migration
cargo sqlx migrate revert

# Reset database (careful!)
cargo sqlx database drop
cargo sqlx database create
cargo sqlx migrate run
```

### 6.4 Code Quality

```bash
# Format code
cargo fmt

# Check formatting
cargo fmt --check

# Run linter
cargo clippy -- -D warnings

# Full check
cargo check --workspace --all-features
```

---

## 7. Project Structure

```
HaiLanGo/
├── src/
│   ├── main.rs              # Application entry point
│   ├── lib.rs               # Library root
│   ├── config/
│   │   ├── mod.rs           # Configuration module
│   │   ├── settings.rs      # Settings loader
│   │   └── urls.rs          # URL routing
│   ├── apps/
│   │   ├── auth/
│   │   │   ├── mod.rs
│   │   │   ├── models.rs    # User model
│   │   │   ├── views.rs     # Auth endpoints
│   │   │   └── services.rs  # Auth logic
│   │   ├── books/
│   │   ├── learning/
│   │   ├── tts/
│   │   ├── stt/
│   │   ├── review/
│   │   └── teacher_mode/
│   └── pages/               # WASM frontend components
├── migrations/              # SQL migrations
├── config/
│   └── settings/            # TOML configuration
├── tests/
│   ├── fixtures/            # Test data
│   └── integration/         # Integration tests
├── storage/                 # Local file storage (gitignored)
├── Cargo.toml
├── compose.yaml             # Podman/Docker compose
└── .env                     # Environment variables (gitignored)
```

---

## 8. Troubleshooting

### Database Connection Failed

```
Error: Connection refused (os error 111)
```

**Solution:**
```bash
# Check if PostgreSQL is running
podman ps | grep postgres

# Start if not running
podman-compose up -d db

# Check logs
podman logs hailango-db
```

### Redis Connection Failed

```
Error: Could not connect to Redis
```

**Solution:**
```bash
# Check if Redis is running
podman ps | grep redis

# Test connection
redis-cli ping  # Should return PONG
```

### Migration Failed

```
Error: relation "users" does not exist
```

**Solution:**
```bash
# Run migrations
cargo sqlx migrate run

# If migrations are corrupted
cargo sqlx database drop
cargo sqlx database create
cargo sqlx migrate run
```

### WASM Build Failed

```
Error: target wasm32-unknown-unknown not found
```

**Solution:**
```bash
# Add WASM target
rustup target add wasm32-unknown-unknown

# Install trunk
cargo install trunk
```

### Port Already in Use

```
Error: Address already in use (os error 98)
```

**Solution:**
```bash
# Find process using port
lsof -i :8080

# Kill process or use different port
APP_PORT=8081 cargo run
```

### API Key Not Working

```
Error: Invalid API key for Google Vision
```

**Solution:**
- Verify API key in `.env`
- Check API is enabled in Google Cloud Console
- Enable mock mode for development: `ENABLE_MOCK_APIS=true`

---

## 9. Next Steps

1. **Explore the Codebase**: Start with `src/apps/auth/` to understand the module structure
2. **Read Architecture Docs**: See [System Architecture](../architecture/system_architecture.md)
3. **Check API Spec**: Review [API Specification](../architecture/api_specification.md)
4. **Run Tests**: Familiarize yourself with the test suite
5. **Try Mock Mode**: Upload a test PDF and see the OCR flow

---

## References

- [CLAUDE.md](../../CLAUDE.md) - Coding guidelines
- [Requirements Definition](../requirements_definition.md)
- [Reinhardt Framework](https://github.com/kent8192/reinhardt-web)
- [SeaQuery Documentation](https://www.sea-ql.org/SeaQuery/)
