# Test Strategy

## 1. Overview

HaiLanGo follows a comprehensive testing strategy to ensure reliability, maintainability, and quality. This document outlines the testing approach, coverage goals, tools, and CI/CD integration.

### Testing Philosophy

- **Test Behavior, Not Implementation**: Focus on what the code does, not how it does it
- **Meaningful Assertions**: Every test must have clear, strict assertions
- **Fast Feedback**: Unit tests should run in milliseconds
- **Realistic Integration Tests**: Use real databases via TestContainers
- **Mock External APIs**: Isolate tests from third-party service availability

---

## 2. Test Pyramid

```
                    ╱╲
                   ╱  ╲
                  ╱ E2E╲           ~5% of tests
                 ╱──────╲          UI flows, critical paths
                ╱        ╲
               ╱Integration╲       ~25% of tests
              ╱────────────╲       API endpoints, DB operations
             ╱              ╲
            ╱   Unit Tests   ╲     ~70% of tests
           ╱──────────────────╲    Business logic, pure functions
          ╱                    ╲
```

### Test Distribution

| Layer | Coverage Target | Focus Areas |
|-------|-----------------|-------------|
| **Unit** | 80%+ | Business logic, models, utilities, SRS algorithm |
| **Integration** | 60%+ | API endpoints, database operations, auth flows |
| **E2E** | Critical paths | User registration, book upload, learning flow |

---

## 3. Coverage Goals

### Overall Target: 75%+ Line Coverage

### Per-Module Targets

| Module | Target | Rationale |
|--------|--------|-----------|
| `apps/auth` | 90% | Security-critical |
| `apps/books` | 80% | Core functionality |
| `apps/review` (SRS) | 95% | Algorithm correctness |
| `apps/learning` | 80% | User-facing features |
| `apps/tts` | 70% | External API wrappers |
| `apps/stt` | 70% | External API wrappers |
| `config` | 60% | Boilerplate code |
| `pages` (WASM) | 50% | UI components |

### Exclusions from Coverage

- Generated code (migrations, OpenAPI types)
- Main entry point (`main.rs`)
- Debug/development utilities
- Third-party wrapper glue code

---

## 4. reinhardt-test Usage

### Basic Test Setup

```rust
use reinhardt_test::prelude::*;

#[tokio::test]
async fn test_user_creation() {
    // Arrange
    let app = TestApp::new().await;
    let user_data = CreateUserRequest {
        email: "test@example.com".to_string(),
        password: "securepassword123".to_string(),
        display_name: "Test User".to_string(),
    };

    // Act
    let response = app
        .post("/api/auth/register")
        .json(&user_data)
        .send()
        .await;

    // Assert
    assert_eq!(response.status(), StatusCode::CREATED);
    let body: RegisterResponse = response.json().await;
    assert_eq!(body.data.user.email, "test@example.com");
    assert!(!body.data.tokens.access_token.is_empty());
}
```

### Authenticated Requests

```rust
#[tokio::test]
async fn test_get_user_books() {
    let app = TestApp::new().await;

    // Create and authenticate user
    let user = app.create_test_user().await;
    let auth = app.authenticate(&user).await;

    // Make authenticated request
    let response = app
        .get("/api/books")
        .bearer_auth(&auth.access_token)
        .send()
        .await;

    assert_eq!(response.status(), StatusCode::OK);
}
```

### Database Assertions

```rust
#[tokio::test]
async fn test_book_persisted_to_database() {
    let app = TestApp::new().await;
    let user = app.create_test_user().await;
    let auth = app.authenticate(&user).await;

    // Upload book
    let response = app
        .post("/api/books/upload")
        .bearer_auth(&auth.access_token)
        .multipart(
            Form::new()
                .file("file", "tests/fixtures/sample.pdf")
                .text("title", "Test Book")
                .text("source_language", "en")
                .text("target_language", "ja"),
        )
        .send()
        .await;

    assert_eq!(response.status(), StatusCode::ACCEPTED);

    // Verify in database
    let book = app
        .db()
        .query_one::<Book>("SELECT * FROM books WHERE user_id = $1", &[&user.id])
        .await
        .unwrap();

    assert_eq!(book.title, "Test Book");
    assert_eq!(book.status, BookStatus::Pending);
}
```

---

## 5. TestContainers Setup

### Configuration

```rust
// tests/common/mod.rs
use testcontainers::{clients::Cli, images::{postgres::Postgres, redis::Redis}};

pub struct TestInfra {
    pub postgres: Container<Postgres>,
    pub redis: Container<Redis>,
    pub db_url: String,
    pub redis_url: String,
}

impl TestInfra {
    pub async fn new() -> Self {
        let docker = Cli::default();

        let postgres = docker.run(
            Postgres::default()
                .with_db_name("hailango_test")
                .with_username("test")
                .with_password("test"),
        );

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

        Self {
            postgres,
            redis,
            db_url,
            redis_url,
        }
    }
}
```

### Using in Tests

```rust
use crate::common::TestInfra;

#[tokio::test]
async fn test_with_real_database() {
    let infra = TestInfra::new().await;

    let app = TestApp::with_config(AppConfig {
        database_url: infra.db_url.clone(),
        redis_url: infra.redis_url.clone(),
        ..Default::default()
    })
    .await;

    // Test with real PostgreSQL and Redis
    // ...
}
```

### Parallel Test Isolation

```rust
use serial_test::serial;

// Tests that share global state run serially
#[tokio::test]
#[serial(database)]
async fn test_migration_rollback() {
    // This test modifies migration state
}

#[tokio::test]
#[serial(database)]
async fn test_seed_data() {
    // This test seeds specific data
}

// Independent tests run in parallel
#[tokio::test]
async fn test_password_hashing() {
    // Pure function test, no external deps
}
```

---

## 6. Mock API Strategy

### Mock Trait Pattern

```rust
// src/services/ocr/mod.rs
#[async_trait]
pub trait OcrProvider: Send + Sync {
    async fn extract_text(&self, image: &[u8]) -> Result<OcrResult, OcrError>;
}

// Production implementation
pub struct GoogleVisionClient { /* ... */ }

#[async_trait]
impl OcrProvider for GoogleVisionClient {
    async fn extract_text(&self, image: &[u8]) -> Result<OcrResult, OcrError> {
        // Real API call
    }
}

// Test mock
pub struct MockOcrClient {
    pub responses: HashMap<String, OcrResult>,
}

#[async_trait]
impl OcrProvider for MockOcrClient {
    async fn extract_text(&self, _image: &[u8]) -> Result<OcrResult, OcrError> {
        Ok(OcrResult {
            text: "Mock extracted text".to_string(),
            confidence: 0.95,
            ..Default::default()
        })
    }
}
```

### Configurable Mock Responses

```rust
#[tokio::test]
async fn test_ocr_error_handling() {
    let mock = MockOcrClient::new()
        .with_error(OcrError::ServiceUnavailable);

    let app = TestApp::with_ocr(Arc::new(mock)).await;

    let response = app
        .post("/api/books/upload")
        .multipart(/* ... */)
        .send()
        .await;

    assert_eq!(response.status(), StatusCode::SERVICE_UNAVAILABLE);
}

#[tokio::test]
async fn test_ocr_low_confidence() {
    let mock = MockOcrClient::new()
        .with_response(OcrResult {
            text: "Partially readable".to_string(),
            confidence: 0.3,  // Low confidence
            ..Default::default()
        });

    let app = TestApp::with_ocr(Arc::new(mock)).await;

    let response = app
        .post("/api/books/upload")
        .multipart(/* ... */)
        .send()
        .await;

    // Should still succeed but flag low confidence
    let body: UploadResponse = response.json().await;
    assert!(body.warnings.contains(&"Low OCR confidence"));
}
```

### Mock HTTP Server (for Integration Tests)

```rust
use wiremock::{MockServer, Mock, ResponseTemplate};
use wiremock::matchers::{method, path};

#[tokio::test]
async fn test_stripe_webhook() {
    let mock_server = MockServer::start().await;

    Mock::given(method("POST"))
        .and(path("/v1/customers"))
        .respond_with(ResponseTemplate::new(200).set_body_json(json!({
            "id": "cus_test123",
            "email": "test@example.com"
        })))
        .mount(&mock_server)
        .await;

    let app = TestApp::with_stripe_url(&mock_server.uri()).await;

    // Test Stripe integration
}
```

---

## 7. Test Fixtures

### Directory Structure

```
tests/
├── fixtures/
│   ├── audio/
│   │   ├── sample_speech.mp3
│   │   └── silence_1s.mp3
│   ├── images/
│   │   ├── sample_page.png
│   │   └── complex_layout.jpg
│   ├── pdfs/
│   │   ├── sample_book.pdf
│   │   └── multilingual.pdf
│   └── json/
│       ├── ocr_response.json
│       └── tts_response.json
├── common/
│   ├── mod.rs
│   ├── test_app.rs
│   └── factories.rs
└── integration/
    ├── auth_tests.rs
    ├── books_tests.rs
    └── learning_tests.rs
```

### Factory Pattern

```rust
// tests/common/factories.rs
pub struct UserFactory;

impl UserFactory {
    pub fn build() -> CreateUserRequest {
        CreateUserRequest {
            email: format!("user_{}@test.com", Uuid::new_v4()),
            password: "TestPassword123!".to_string(),
            display_name: "Test User".to_string(),
            native_language: "en".to_string(),
        }
    }

    pub fn with_email(email: &str) -> CreateUserRequest {
        CreateUserRequest {
            email: email.to_string(),
            ..Self::build()
        }
    }
}

pub struct BookFactory;

impl BookFactory {
    pub fn build(user_id: Uuid) -> Book {
        Book {
            id: Uuid::new_v4(),
            user_id,
            title: "Test Book".to_string(),
            source_language: "en".to_string(),
            target_language: "ja".to_string(),
            status: BookStatus::Ready,
            total_pages: 10,
            ..Default::default()
        }
    }
}
```

---

## 8. Testing Specific Components

### SRS Algorithm Tests

```rust
#[cfg(test)]
mod srs_tests {
    use super::*;

    #[test]
    fn test_sm2_perfect_recall() {
        let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());

        // First review - perfect
        schedule.update_after_review(5);
        assert_eq!(schedule.interval_days, 1);
        assert_eq!(schedule.repetitions, 1);

        // Second review - perfect
        schedule.update_after_review(5);
        assert_eq!(schedule.interval_days, 6);
        assert_eq!(schedule.repetitions, 2);

        // Third review - perfect
        schedule.update_after_review(5);
        assert!(schedule.interval_days > 6);
        assert_eq!(schedule.repetitions, 3);
    }

    #[test]
    fn test_sm2_failed_review_resets() {
        let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());

        // Build up interval
        schedule.update_after_review(5);
        schedule.update_after_review(5);
        schedule.update_after_review(5);
        let old_interval = schedule.interval_days;

        // Fail
        schedule.update_after_review(1);

        assert_eq!(schedule.interval_days, 1);  // Reset to 1
        assert_eq!(schedule.repetitions, 0);    // Reset repetitions
        assert!(schedule.easiness_factor < 2.5); // EF decreased
    }

    #[test]
    fn test_sm2_easiness_bounds() {
        let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());

        // Many failures should not drop EF below 1.3
        for _ in 0..20 {
            schedule.update_after_review(0);
        }

        assert!(schedule.easiness_factor >= 1.3);
    }
}
```

### Authentication Tests

```rust
#[tokio::test]
async fn test_jwt_token_expiry() {
    let app = TestApp::new().await;
    let user = app.create_test_user().await;
    let auth = app.authenticate(&user).await;

    // Token should work immediately
    let response = app
        .get("/api/users/me")
        .bearer_auth(&auth.access_token)
        .send()
        .await;
    assert_eq!(response.status(), StatusCode::OK);

    // Fast-forward time (mock time in tests)
    app.advance_time(Duration::hours(2)).await;

    // Token should be expired
    let response = app
        .get("/api/users/me")
        .bearer_auth(&auth.access_token)
        .send()
        .await;
    assert_eq!(response.status(), StatusCode::UNAUTHORIZED);
}

#[tokio::test]
async fn test_refresh_token_rotation() {
    let app = TestApp::new().await;
    let user = app.create_test_user().await;
    let auth = app.authenticate(&user).await;

    // Refresh token
    let response = app
        .post("/api/auth/refresh")
        .json(&json!({ "refresh_token": auth.refresh_token }))
        .send()
        .await;

    let new_auth: TokenResponse = response.json().await;

    // Old refresh token should be invalidated
    let response = app
        .post("/api/auth/refresh")
        .json(&json!({ "refresh_token": auth.refresh_token }))
        .send()
        .await;
    assert_eq!(response.status(), StatusCode::UNAUTHORIZED);

    // New refresh token should work
    let response = app
        .post("/api/auth/refresh")
        .json(&json!({ "refresh_token": new_auth.refresh_token }))
        .send()
        .await;
    assert_eq!(response.status(), StatusCode::OK);
}
```

---

## 9. CI/CD Integration

### GitHub Actions Workflow

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

env:
  CARGO_TERM_COLOR: always
  RUSTFLAGS: "-Dwarnings"

jobs:
  check:
    name: Check
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: dtolnay/rust-toolchain@stable
      - uses: Swatinem/rust-cache@v2
      - run: cargo check --workspace --all-features

  fmt:
    name: Format
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: dtolnay/rust-toolchain@stable
        with:
          components: rustfmt
      - run: cargo fmt --all --check

  clippy:
    name: Clippy
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: dtolnay/rust-toolchain@stable
        with:
          components: clippy
      - uses: Swatinem/rust-cache@v2
      - run: cargo clippy --workspace --all-features -- -D warnings

  test:
    name: Test
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: hailango_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:7
        ports:
          - 6379:6379
    env:
      DATABASE_URL: postgresql://test:test@localhost:5432/hailango_test
      REDIS_URL: redis://localhost:6379
    steps:
      - uses: actions/checkout@v4
      - uses: dtolnay/rust-toolchain@stable
      - uses: Swatinem/rust-cache@v2
      - name: Install sqlx-cli
        run: cargo install sqlx-cli --no-default-features --features postgres
      - name: Run migrations
        run: sqlx migrate run
      - name: Run tests
        run: cargo test --workspace --all-features

  coverage:
    name: Coverage
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: hailango_test
        ports:
          - 5432:5432
      redis:
        image: redis:7
        ports:
          - 6379:6379
    env:
      DATABASE_URL: postgresql://test:test@localhost:5432/hailango_test
      REDIS_URL: redis://localhost:6379
    steps:
      - uses: actions/checkout@v4
      - uses: dtolnay/rust-toolchain@stable
        with:
          components: llvm-tools-preview
      - uses: Swatinem/rust-cache@v2
      - name: Install cargo-llvm-cov
        run: cargo install cargo-llvm-cov
      - name: Install sqlx-cli
        run: cargo install sqlx-cli --no-default-features --features postgres
      - name: Run migrations
        run: sqlx migrate run
      - name: Generate coverage report
        run: cargo llvm-cov --workspace --all-features --lcov --output-path lcov.info
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: lcov.info
          fail_ci_if_error: true
```

### Pre-commit Hooks

```bash
# .git/hooks/pre-commit
#!/bin/sh

# Format check
cargo fmt --check || {
    echo "Run 'cargo fmt' to fix formatting"
    exit 1
}

# Clippy
cargo clippy --workspace --all-features -- -D warnings || {
    echo "Fix clippy warnings before committing"
    exit 1
}

# Quick test
cargo test --workspace --lib || {
    echo "Tests failed"
    exit 1
}
```

---

## 10. Performance Testing

### Load Testing with k6

```javascript
// tests/load/basic_load.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 20 },  // Ramp up
    { duration: '1m', target: 20 },   // Stay at 20
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% under 500ms
    http_req_failed: ['rate<0.01'],    // Error rate under 1%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Login
  const loginRes = http.post(`${BASE_URL}/api/auth/login`, JSON.stringify({
    email: 'loadtest@example.com',
    password: 'LoadTest123!',
  }), { headers: { 'Content-Type': 'application/json' } });

  check(loginRes, { 'login succeeded': (r) => r.status === 200 });

  const token = loginRes.json('data.tokens.access_token');

  // Get books
  const booksRes = http.get(`${BASE_URL}/api/books`, {
    headers: { Authorization: `Bearer ${token}` },
  });

  check(booksRes, { 'get books succeeded': (r) => r.status === 200 });

  sleep(1);
}
```

### Running Load Tests

```bash
# Run load test
k6 run tests/load/basic_load.js

# With environment variable
k6 run -e BASE_URL=http://staging.hailango.com tests/load/basic_load.js
```

---

## References

- [Getting Started](../guides/getting_started.md)
- [API Specification](../architecture/api_specification.md)
- [reinhardt-test Documentation](https://docs.rs/reinhardt-test)
- [TestContainers-rs](https://github.com/testcontainers/testcontainers-rs)
