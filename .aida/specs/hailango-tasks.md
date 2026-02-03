# HaiLanGo - Implementation Tasks (TDD)

**Project**: HaiLanGo AI Language Learning Platform
**Generated**: 2026-02-02
**Methodology**: Test-Driven Development (RED → GREEN → REFACTOR)
**Framework**: Reinhardt (Rust Full-Stack)

---

## TDD Workflow

For each task:
1. **RED**: Write failing test first
2. **GREEN**: Implement minimum code to pass test
3. **REFACTOR**: Clean up implementation while keeping tests green
4. **COMMIT**: Commit with test evidence

---

## Phase 1: MVP Foundation (P0)

### Epic 1.1: Project Setup

#### Task 1.1.1: Initialize Cargo Workspace
**Priority**: P0
**Estimated Time**: 1 hour

**TDD Steps**:
1. **RED**: N/A (setup task)
2. **GREEN**:
   - Create `Cargo.toml` workspace manifest
   - Add dependencies:
     ```toml
     [workspace]
     members = ["src"]

     [workspace.dependencies]
     reinhardt-db = "0.1"
     reinhardt-rest = "0.1"
     reinhardt-pages = "0.1"
     reinhardt-auth = "0.1"
     reinhardt-websockets = "0.1"
     reinhardt-test = "0.1"
     sqlx = { version = "0.7", features = ["postgres", "uuid", "chrono", "json"] }
     tokio = { version = "1", features = ["full"] }
     serde = { version = "1", features = ["derive"] }
     serde_json = "1"
     uuid = { version = "1", features = ["v4", "serde"] }
     chrono = { version = "0.4", features = ["serde"] }
     ```
3. **REFACTOR**: Organize dependencies by category

**Acceptance**:
- [ ] `cargo check --workspace` succeeds
- [ ] All reinhardt crates resolve correctly

---

#### Task 1.1.2: Setup Project Structure
**Priority**: P0
**Estimated Time**: 1 hour

**TDD Steps**:
1. **RED**: N/A (setup task)
2. **GREEN**:
   - Create directory structure following Rust 2024 Edition:
     ```
     src/
     ├── main.rs
     ├── lib.rs
     ├── config/
     │   ├── module.rs
     │   ├── settings/
     │   │   ├── module.rs
     │   │   ├── base.rs
     │   │   ├── development.rs
     │   │   └── production.rs
     │   ├── urls.rs
     │   └── apps.rs
     ├── apps/
     │   ├── auth/module.rs
     │   ├── books/module.rs
     │   ├── learning/module.rs
     │   ├── tts/module.rs
     │   ├── stt/module.rs
     │   ├── review/module.rs
     │   └── teacher_mode/module.rs
     ├── pages/module.rs
     ├── services/module.rs
     └── utils/module.rs
     tests/
     ├── common/mod.rs
     └── fixtures/
     ```
   - Create empty `module.rs` files (NOT `mod.rs`)
3. **REFACTOR**: Add module declarations to `lib.rs`

**Acceptance**:
- [ ] No `mod.rs` files exist
- [ ] All modules use `module.rs` pattern
- [ ] `cargo build --workspace` succeeds

---

#### Task 1.1.3: Configure Development Environment
**Priority**: P0
**Estimated Time**: 2 hours

**TDD Steps**:
1. **RED**: N/A (setup task)
2. **GREEN**:
   - Create `compose.yaml` for Podman/Docker:
     ```yaml
     services:
       app:
         build: .
         ports: ["8080:8080"]
         environment:
           - DATABASE_URL=postgresql://hailango:password@db:5432/hailango
           - REDIS_URL=redis://redis:6379
         depends_on: [db, redis]
       db:
         image: postgres:16
         environment:
           POSTGRES_DB: hailango
           POSTGRES_USER: hailango
           POSTGRES_PASSWORD: password
       redis:
         image: redis:7
     ```
   - Create `.env.example`
   - Create `Dockerfile`
3. **REFACTOR**: Add health checks, restart policies

**Acceptance**:
- [ ] `podman-compose up` starts all services
- [ ] Can connect to PostgreSQL and Redis
- [ ] App container builds successfully

---

### Epic 1.2: Database Setup

#### Task 1.2.1: Initial Migration (Users Table)
**Priority**: P0
**Estimated Time**: 2 hours

**Test File**: `migrations/0001_initial.sql`
**Implementation File**: SQL migration

**TDD Steps**:
1. **RED**: Create test that expects `users` table to exist:
   ```rust
   // tests/integration/database_tests.rs
   #[tokio::test]
   async fn test_users_table_exists() {
       let pool = PgPool::connect(&db_url).await.unwrap();
       sqlx::migrate!("./migrations").run(&pool).await.unwrap();

       let result: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM users")
           .fetch_one(&pool)
           .await
           .unwrap();

       assert_eq!(result.0, 0);
   }
   ```
2. **GREEN**: Create migration:
   ```sql
   -- migrations/0001_initial.sql
   CREATE EXTENSION IF NOT EXISTS "pgcrypto";

   CREATE TABLE users (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       email VARCHAR(255) NOT NULL UNIQUE,
       password_hash VARCHAR(255),
       display_name VARCHAR(100) NOT NULL,
       native_language VARCHAR(10) NOT NULL DEFAULT 'en',
       avatar_url TEXT,
       oauth_provider VARCHAR(50),
       oauth_id VARCHAR(255),
       email_verified BOOLEAN DEFAULT FALSE,
       created_at TIMESTAMPTZ DEFAULT NOW(),
       updated_at TIMESTAMPTZ DEFAULT NOW(),
       last_login_at TIMESTAMPTZ
   );

   CREATE UNIQUE INDEX idx_users_email ON users(email);
   CREATE UNIQUE INDEX idx_users_oauth ON users(oauth_provider, oauth_id)
       WHERE oauth_provider IS NOT NULL;
   ```
3. **REFACTOR**: Add trigger for `updated_at`

**Acceptance**:
- [ ] Migration runs without errors
- [ ] Test passes
- [ ] Indexes created correctly

---

#### Task 1.2.2: Books & Pages Migration
**Priority**: P0
**Estimated Time**: 2 hours

**Test File**: `tests/integration/database_tests.rs`
**Implementation File**: `migrations/0002_books_and_pages.sql`

**TDD Steps**:
1. **RED**:
   ```rust
   #[tokio::test]
   async fn test_books_table_exists() {
       let pool = setup_test_db().await;

       let result: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM books")
           .fetch_one(&pool)
           .await
           .unwrap();

       assert_eq!(result.0, 0);
   }

   #[tokio::test]
   async fn test_pages_foreign_key_cascade() {
       let pool = setup_test_db().await;

       // Insert user and book
       let user_id = insert_test_user(&pool).await;
       let book_id = insert_test_book(&pool, user_id).await;
       insert_test_page(&pool, book_id).await;

       // Delete book
       sqlx::query("DELETE FROM books WHERE id = $1")
           .bind(book_id)
           .execute(&pool)
           .await
           .unwrap();

       // Verify pages were cascade deleted
       let count: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM pages WHERE book_id = $1")
           .bind(book_id)
           .fetch_one(&pool)
           .await
           .unwrap();

       assert_eq!(count.0, 0);
   }
   ```
2. **GREEN**: Create migration
3. **REFACTOR**: Optimize indexes

**Acceptance**:
- [ ] All tests pass
- [ ] Foreign key constraints work
- [ ] Cascade deletes work correctly

---

### Epic 1.3: Authentication System

#### Task 1.3.1: User Model
**Priority**: P0
**Estimated Time**: 3 hours

**Test File**: `src/apps/auth/tests.rs`
**Implementation File**: `src/apps/auth/models.rs`

**TDD Steps**:
1. **RED**:
   ```rust
   #[tokio::test]
   async fn test_user_password_hashing() {
       let mut user = User {
           id: Uuid::new_v4(),
           email: "test@example.com".to_string(),
           password_hash: None,
           display_name: "Test".to_string(),
           ..Default::default()
       };

       user.set_password("SecurePassword123!");

       assert!(user.password_hash.is_some());
       assert_ne!(user.password_hash.as_ref().unwrap(), "SecurePassword123!");
   }

   #[tokio::test]
   async fn test_user_password_verification() {
       let mut user = User::default();
       user.set_password("SecurePassword123!");

       assert!(user.verify_password("SecurePassword123!"));
       assert!(!user.verify_password("WrongPassword"));
   }

   #[tokio::test]
   async fn test_user_find_by_email() {
       let pool = setup_test_db().await;

       let user = User::create(&pool, CreateUserRequest {
           email: "test@example.com".to_string(),
           password: "Password123!".to_string(),
           display_name: "Test User".to_string(),
           native_language: "en".to_string(),
       }).await.unwrap();

       let found = User::find_by_email(&pool, "test@example.com")
           .await
           .unwrap()
           .unwrap();

       assert_eq!(found.id, user.id);
       assert_eq!(found.email, "test@example.com");
   }
   ```
2. **GREEN**: Implement User model with reinhardt-db
3. **REFACTOR**: Extract password hashing to utility

**Acceptance**:
- [ ] All tests pass
- [ ] Argon2id hashing works
- [ ] Database queries use reinhardt-db

---

#### Task 1.3.2: JWT Authentication
**Priority**: P0
**Estimated Time**: 4 hours

**Test File**: `src/apps/auth/tests.rs`
**Implementation File**: `src/apps/auth/middleware.rs`

**TDD Steps**:
1. **RED**:
   ```rust
   #[tokio::test]
   async fn test_jwt_token_generation() {
       let user = create_test_user();
       let tokens = generate_tokens(&user).unwrap();

       assert!(!tokens.access_token.is_empty());
       assert!(!tokens.refresh_token.is_empty());
       assert_eq!(tokens.token_type, "Bearer");
       assert_eq!(tokens.expires_in, 3600);
   }

   #[tokio::test]
   async fn test_jwt_token_validation() {
       let user = create_test_user();
       let tokens = generate_tokens(&user).unwrap();

       let claims = validate_token(&tokens.access_token).unwrap();

       assert_eq!(claims.sub, user.id);
       assert_eq!(claims.email, user.email);
   }

   #[tokio::test]
   async fn test_jwt_token_expiry() {
       let user = create_test_user();
       let tokens = generate_tokens_with_expiry(&user, -3600).unwrap();  // Expired 1 hour ago

       let result = validate_token(&tokens.access_token);

       assert!(result.is_err());
       assert_eq!(result.unwrap_err().kind(), TokenErrorKind::Expired);
   }

   #[tokio::test]
   async fn test_protected_endpoint_requires_auth() {
       let app = TestApp::new().await;

       let response = app
           .get("/api/users/me")
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::UNAUTHORIZED);
   }

   #[tokio::test]
   async fn test_protected_endpoint_with_valid_token() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;
       let tokens = app.authenticate(&user).await;

       let response = app
           .get("/api/users/me")
           .bearer_auth(&tokens.access_token)
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::OK);

       let body: UserResponse = response.json().await;
       assert_eq!(body.email, user.email);
   }
   ```
2. **GREEN**: Implement JWT middleware with reinhardt-auth
3. **REFACTOR**: Extract token utilities

**Acceptance**:
- [ ] All tests pass
- [ ] JWT tokens expire correctly
- [ ] Protected endpoints enforce authentication

---

#### Task 1.3.3: Registration & Login Endpoints
**Priority**: P0
**Estimated Time**: 4 hours

**Test File**: `tests/integration/auth_tests.rs`
**Implementation File**: `src/apps/auth/viewsets.rs`

**TDD Steps**:
1. **RED**:
   ```rust
   #[tokio::test]
   async fn test_register_success() {
       let app = TestApp::new().await;

       let response = app
           .post("/api/auth/register")
           .json(&json!({
               "email": "newuser@example.com",
               "password": "SecurePassword123!",
               "display_name": "New User",
               "native_language": "en"
           }))
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::CREATED);

       let body: RegisterResponse = response.json().await;
       assert_eq!(body.data.user.email, "newuser@example.com");
       assert!(!body.data.tokens.access_token.is_empty());

       // Verify user in database
       let user = User::find_by_email(&app.db, "newuser@example.com")
           .await
           .unwrap()
           .unwrap();
       assert_eq!(user.display_name, "New User");
   }

   #[tokio::test]
   async fn test_register_duplicate_email() {
       let app = TestApp::new().await;

       // First registration
       app.post("/api/auth/register")
           .json(&json!({
               "email": "duplicate@example.com",
               "password": "Password123!",
               "display_name": "User 1"
           }))
           .send()
           .await;

       // Second registration with same email
       let response = app
           .post("/api/auth/register")
           .json(&json!({
               "email": "duplicate@example.com",
               "password": "Password456!",
               "display_name": "User 2"
           }))
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::UNPROCESSABLE_ENTITY);

       let body: ErrorResponse = response.json().await;
       assert_eq!(body.error.code, "VALIDATION_ERROR");
   }

   #[tokio::test]
   async fn test_login_success() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;

       let response = app
           .post("/api/auth/login")
           .json(&json!({
               "email": user.email,
               "password": "TestPassword123!"
           }))
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::OK);

       let body: LoginResponse = response.json().await;
       assert_eq!(body.data.user.id, user.id);
       assert!(!body.data.tokens.access_token.is_empty());
   }

   #[tokio::test]
   async fn test_login_wrong_password() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;

       let response = app
           .post("/api/auth/login")
           .json(&json!({
               "email": user.email,
               "password": "WrongPassword!"
           }))
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::UNAUTHORIZED);
   }
   ```
2. **GREEN**: Implement AuthViewSet with reinhardt-rest
3. **REFACTOR**: Extract validation logic

**Acceptance**:
- [ ] All tests pass
- [ ] Email validation works (RFC 5322)
- [ ] Password validation enforces complexity
- [ ] Duplicate email returns 422

---

#### Task 1.3.4: OAuth (Google) Integration
**Priority**: P0
**Estimated Time**: 6 hours

**Test File**: `tests/integration/oauth_tests.rs`
**Implementation File**: `src/apps/auth/oauth.rs`

**TDD Steps**:
1. **RED**:
   ```rust
   #[tokio::test]
   async fn test_oauth_google_new_user() {
       let app = TestApp::new().await;

       // Mock Google OAuth response
       let mock_server = wiremock::MockServer::start().await;
       wiremock::Mock::given(wiremock::matchers::method("POST"))
           .and(wiremock::matchers::path("/token"))
           .respond_with(wiremock::ResponseTemplate::new(200).set_body_json(json!({
               "access_token": "mock_access_token",
               "token_type": "Bearer",
               "expires_in": 3600
           })))
           .mount(&mock_server)
           .await;

       let response = app
           .post("/api/auth/oauth/google")
           .json(&json!({
               "code": "mock_authorization_code",
               "redirect_uri": "http://localhost:3000/callback"
           }))
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::OK);

       let body: OAuthResponse = response.json().await;
       assert_eq!(body.data.user.oauth_provider, Some("google".to_string()));
       assert!(!body.data.tokens.access_token.is_empty());

       // Verify user created
       let user = User::find_by_email(&app.db, body.data.user.email)
           .await
           .unwrap()
           .unwrap();
       assert_eq!(user.oauth_provider, Some("google".to_string()));
       assert!(user.password_hash.is_none());
   }

   #[tokio::test]
   async fn test_oauth_google_existing_user() {
       let app = TestApp::new().await;

       // Create user via OAuth first time
       let first_response = oauth_login(&app, "google_user_123").await;
       let first_user_id = first_response.data.user.id;

       // Login again with same Google account
       let second_response = oauth_login(&app, "google_user_123").await;

       assert_eq!(second_response.data.user.id, first_user_id);
       assert_ne!(
           second_response.data.tokens.access_token,
           first_response.data.tokens.access_token
       );
   }
   ```
2. **GREEN**: Implement Google OAuth provider
3. **REFACTOR**: Abstract OAuth provider trait

**Acceptance**:
- [ ] All tests pass
- [ ] New users created via OAuth
- [ ] Existing users login via OAuth
- [ ] OAuth state validation works

---

### Epic 1.4: Book Upload & OCR

#### Task 1.4.1: Book Model & ViewSet
**Priority**: P0
**Estimated Time**: 4 hours

**Test File**: `tests/integration/books_tests.rs`
**Implementation File**: `src/apps/books/models.rs`, `src/apps/books/viewsets.rs`

**TDD Steps**:
1. **RED**:
   ```rust
   #[tokio::test]
   async fn test_list_books_empty() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;
       let auth = app.authenticate(&user).await;

       let response = app
           .get("/api/books")
           .bearer_auth(&auth.access_token)
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::OK);

       let body: ListBooksResponse = response.json().await;
       assert_eq!(body.data.len(), 0);
   }

   #[tokio::test]
   async fn test_list_books_with_data() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;
       let auth = app.authenticate(&user).await;

       // Create test books
       let book1 = Book::create(&app.db, user.id, CreateBookRequest {
           title: "Book 1".to_string(),
           source_language: "en".to_string(),
           target_language: "ja".to_string(),
           reference_language: None,
       }).await.unwrap();

       let book2 = Book::create(&app.db, user.id, CreateBookRequest {
           title: "Book 2".to_string(),
           source_language: "en".to_string(),
           target_language: "es".to_string(),
           reference_language: None,
       }).await.unwrap();

       let response = app
           .get("/api/books")
           .bearer_auth(&auth.access_token)
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::OK);

       let body: ListBooksResponse = response.json().await;
       assert_eq!(body.data.len(), 2);
       assert_eq!(body.data[0].title, "Book 2");  // Newest first
       assert_eq!(body.data[1].title, "Book 1");
   }

   #[tokio::test]
   async fn test_get_book_details() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;
       let auth = app.authenticate(&user).await;

       let book = Book::create(&app.db, user.id, CreateBookRequest {
           title: "Test Book".to_string(),
           source_language: "en".to_string(),
           target_language: "ja".to_string(),
           reference_language: None,
       }).await.unwrap();

       let response = app
           .get(&format!("/api/books/{}", book.id))
           .bearer_auth(&auth.access_token)
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::OK);

       let body: BookDetailResponse = response.json().await;
       assert_eq!(body.data.id, book.id);
       assert_eq!(body.data.title, "Test Book");
   }

   #[tokio::test]
   async fn test_get_book_forbidden_other_user() {
       let app = TestApp::new().await;

       let user1 = app.create_test_user().await;
       let book = Book::create(&app.db, user1.id, CreateBookRequest {
           title: "User 1 Book".to_string(),
           source_language: "en".to_string(),
           target_language: "ja".to_string(),
           reference_language: None,
       }).await.unwrap();

       let user2 = app.create_test_user().await;
       let auth2 = app.authenticate(&user2).await;

       let response = app
           .get(&format!("/api/books/{}", book.id))
           .bearer_auth(&auth2.access_token)
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::FORBIDDEN);
   }
   ```
2. **GREEN**: Implement Book model and BookViewSet
3. **REFACTOR**: Add pagination support

**Acceptance**:
- [ ] All tests pass
- [ ] Books scoped to authenticated user
- [ ] Permission checks enforced

---

#### Task 1.4.2: PDF Upload Handler
**Priority**: P0
**Estimated Time**: 6 hours

**Test File**: `tests/integration/upload_tests.rs`
**Implementation File**: `src/apps/books/viewsets.rs`

**TDD Steps**:
1. **RED**:
   ```rust
   #[tokio::test]
   async fn test_upload_pdf_success() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;
       let auth = app.authenticate(&user).await;

       let response = app
           .post("/api/books/upload")
           .bearer_auth(&auth.access_token)
           .multipart(
               multipart::Form::new()
                   .file("file", "tests/fixtures/pdfs/sample.pdf")
                   .text("title", "Sample Book")
                   .text("source_language", "en")
                   .text("target_language", "ja"),
           )
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::ACCEPTED);

       let body: UploadResponse = response.json().await;
       assert_eq!(body.data.status, "pending");
       assert!(!body.data.job_id.is_empty());

       // Verify book created in database
       let book = Book::find_by_id(&app.db, body.data.id).await.unwrap();
       assert_eq!(book.title, "Sample Book");
       assert_eq!(book.status, BookStatus::Pending);
   }

   #[tokio::test]
   async fn test_upload_file_too_large() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;
       let auth = app.authenticate(&user).await;

       // Generate 51MB file
       let large_file = vec![0u8; 51 * 1024 * 1024];

       let response = app
           .post("/api/books/upload")
           .bearer_auth(&auth.access_token)
           .multipart(
               multipart::Form::new()
                   .part("file", multipart::Part::bytes(large_file).file_name("large.pdf"))
                   .text("title", "Large Book")
                   .text("source_language", "en")
                   .text("target_language", "ja"),
           )
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::PAYLOAD_TOO_LARGE);
   }

   #[tokio::test]
   async fn test_upload_invalid_language_code() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;
       let auth = app.authenticate(&user).await;

       let response = app
           .post("/api/books/upload")
           .bearer_auth(&auth.access_token)
           .multipart(
               multipart::Form::new()
                   .file("file", "tests/fixtures/pdfs/sample.pdf")
                   .text("title", "Sample Book")
                   .text("source_language", "invalid")
                   .text("target_language", "ja"),
           )
           .send()
           .await;

       assert_eq!(response.status(), StatusCode::UNPROCESSABLE_ENTITY);

       let body: ErrorResponse = response.json().await;
       assert_eq!(body.error.code, "VALIDATION_ERROR");
   }
   ```
2. **GREEN**: Implement upload handler with validation
3. **REFACTOR**: Extract file handling utilities

**Acceptance**:
- [ ] All tests pass
- [ ] File size validation (max 50MB)
- [ ] Language code validation (ISO 639-1)
- [ ] Book created with status "pending"

---

#### Task 1.4.3: OCR Provider Trait & Google Vision Implementation
**Priority**: P0
**Estimated Time**: 8 hours

**Test File**: `src/apps/books/ocr/tests.rs`
**Implementation File**: `src/apps/books/ocr/provider.rs`, `src/apps/books/ocr/google.rs`

**TDD Steps**:
1. **RED**:
   ```rust
   #[tokio::test]
   async fn test_google_vision_extract_text_success() {
       let client = GoogleVisionClient::new("test_api_key");

       let mock_server = wiremock::MockServer::start().await;
       wiremock::Mock::given(wiremock::matchers::method("POST"))
           .and(wiremock::matchers::path("/v1/images:annotate"))
           .respond_with(wiremock::ResponseTemplate::new(200).set_body_json(json!({
               "responses": [{
                   "fullTextAnnotation": {
                       "text": "Sample extracted text\nLine 2\nLine 3",
                       "pages": [...]
                   }
               }]
           })))
           .mount(&mock_server)
           .await;

       let image_data = std::fs::read("tests/fixtures/images/sample_page.png").unwrap();
       let result = client.extract_text(&image_data).await.unwrap();

       assert_eq!(result.text, "Sample extracted text\nLine 2\nLine 3");
       assert!(result.confidence > 0.8);
   }

   #[tokio::test]
   async fn test_google_vision_extract_text_pdf() {
       let client = GoogleVisionClient::new("test_api_key");

       let pdf_data = std::fs::read("tests/fixtures/pdfs/sample.pdf").unwrap();
       let results = client.extract_text_pdf(&pdf_data).await.unwrap();

       assert_eq!(results.len(), 3);  // 3-page PDF
       assert!(!results[0].text.is_empty());
       assert!(!results[1].text.is_empty());
       assert!(!results[2].text.is_empty());
   }

   #[tokio::test]
   async fn test_google_vision_api_error() {
       let client = GoogleVisionClient::new("invalid_key");

       let mock_server = wiremock::MockServer::start().await;
       wiremock::Mock::given(wiremock::matchers::method("POST"))
           .respond_with(wiremock::ResponseTemplate::new(403).set_body_json(json!({
               "error": {
                   "code": 403,
                   "message": "API key not valid"
               }
           })))
           .mount(&mock_server)
           .await;

       let image_data = vec![1, 2, 3];
       let result = client.extract_text(&image_data).await;

       assert!(result.is_err());
       assert_eq!(result.unwrap_err().kind(), OcrErrorKind::Unauthorized);
   }
   ```
2. **GREEN**: Implement OcrProvider trait and Google Vision client
3. **REFACTOR**: Add circuit breaker, retry logic

**Acceptance**:
- [ ] All tests pass
- [ ] Trait abstraction works
- [ ] Google Vision API integration works
- [ ] Error handling complete

---

#### Task 1.4.4: OCR Job Queue & Processing
**Priority**: P0
**Estimated Time**: 8 hours

**Test File**: `tests/integration/ocr_processing_tests.rs`
**Implementation File**: `src/apps/books/ocr/processor.rs`

**TDD Steps**:
1. **RED**:
   ```rust
   #[tokio::test]
   async fn test_ocr_job_processing_success() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;

       // Configure mock OCR
       app.ocr_mock.set_response(OcrResult {
           text: "Page 1 content".to_string(),
           confidence: 0.95,
           layout: vec![],
       });

       // Create book with OCR job
       let book = Book::create(&app.db, user.id, CreateBookRequest {
           title: "Test Book".to_string(),
           source_language: "en".to_string(),
           target_language: "ja".to_string(),
           reference_language: None,
       }).await.unwrap();

       let job_id = queue_ocr_job(&app.redis, book.id, vec![/* page data */]).await.unwrap();

       // Process job
       process_ocr_job(&app, &job_id).await.unwrap();

       // Verify book status updated
       let updated_book = Book::find_by_id(&app.db, book.id).await.unwrap();
       assert_eq!(updated_book.status, BookStatus::Ready);
       assert_eq!(updated_book.total_pages, 1);

       // Verify page created
       let pages = Page::find_by_book(&app.db, book.id).await.unwrap();
       assert_eq!(pages.len(), 1);
       assert_eq!(pages[0].original_content.as_ref().unwrap(), "Page 1 content");
       assert!(pages[0].is_processed);
   }

   #[tokio::test]
   async fn test_ocr_job_processing_failure() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;

       // Configure mock OCR to fail
       app.ocr_mock.set_error(OcrError::ServiceUnavailable);

       let book = Book::create(&app.db, user.id, CreateBookRequest {
           title: "Test Book".to_string(),
           source_language: "en".to_string(),
           target_language: "ja".to_string(),
           reference_language: None,
       }).await.unwrap();

       let job_id = queue_ocr_job(&app.redis, book.id, vec![/* page data */]).await.unwrap();

       let result = process_ocr_job(&app, &job_id).await;

       assert!(result.is_err());

       // Verify book status updated to error
       let updated_book = Book::find_by_id(&app.db, book.id).await.unwrap();
       assert_eq!(updated_book.status, BookStatus::Error);
   }

   #[tokio::test]
   async fn test_ocr_status_polling() {
       let app = TestApp::new().await;
       let user = app.create_test_user().await;
       let auth = app.authenticate(&user).await;

       let book = Book::create(&app.db, user.id, CreateBookRequest {
           title: "Test Book".to_string(),
           source_language: "en".to_string(),
           target_language: "ja".to_string(),
           reference_language: None,
       }).await.unwrap();

       // Book initially pending
       let response = app
           .get(&format!("/api/books/{}/status", book.id))
           .bearer_auth(&auth.access_token)
           .send()
           .await;

       let body: StatusResponse = response.json().await;
       assert_eq!(body.data.status, "pending");
       assert_eq!(body.data.progress.percentage, 0.0);
   }
   ```
2. **GREEN**: Implement job queue and processor
3. **REFACTOR**: Add background worker pattern

**Acceptance**:
- [ ] All tests pass
- [ ] Jobs queued in Redis
- [ ] Background processing works
- [ ] Status polling works

---

## Phase 2: Core Features (P1)

### Epic 2.1: TTS System

#### Task 2.1.1: TTS Provider Trait & Google Cloud TTS
**Priority**: P1
**Estimated Time**: 6 hours

**Test File**: `src/apps/tts/tests.rs`
**Implementation File**: `src/apps/tts/providers/google.rs`

**TDD Steps**:
1. **RED**: Write tests for TTS synthesis
2. **GREEN**: Implement Google Cloud TTS provider
3. **REFACTOR**: Add audio caching

#### Task 2.1.2: TTS ViewSet & Audio Streaming
**Priority**: P1
**Estimated Time**: 6 hours

**Test File**: `tests/integration/tts_tests.rs`
**Implementation File**: `src/apps/tts/viewsets.rs`

---

### Epic 2.2: STT & Pronunciation Evaluation

#### Task 2.2.1: STT Provider Trait & Whisper Integration
**Priority**: P1
**Estimated Time**: 8 hours

**Test File**: `src/apps/stt/tests.rs`
**Implementation File**: `src/apps/stt/providers/whisper.rs`

#### Task 2.2.2: Pronunciation Scoring Algorithm
**Priority**: P1
**Estimated Time**: 8 hours

**Test File**: `src/apps/stt/evaluation/tests.rs`
**Implementation File**: `src/apps/stt/evaluation.rs`

---

### Epic 2.3: SRS System

#### Task 2.3.1: SM-2 Algorithm Implementation
**Priority**: P1
**Estimated Time**: 8 hours

**Test File**: `src/apps/review/tests.rs`
**Implementation File**: `src/apps/review/sm2.rs`

**TDD Steps**:
1. **RED**:
   ```rust
   #[test]
   fn test_sm2_first_review_perfect() {
       let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());
       schedule.update_after_review(5);

       assert_eq!(schedule.interval_days, 1);
       assert_eq!(schedule.repetitions, 1);
       assert_eq!(schedule.correct_count, 1);
   }

   #[test]
   fn test_sm2_second_review_perfect() {
       let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());
       schedule.update_after_review(5);
       schedule.update_after_review(5);

       assert_eq!(schedule.interval_days, 6);
       assert_eq!(schedule.repetitions, 2);
   }

   #[test]
   fn test_sm2_failed_review_resets() {
       let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());
       schedule.update_after_review(5);
       schedule.update_after_review(5);
       schedule.update_after_review(5);

       let interval_before = schedule.interval_days;
       schedule.update_after_review(1);  // Fail

       assert_eq!(schedule.interval_days, 1);
       assert_eq!(schedule.repetitions, 0);
       assert!(schedule.easiness_factor < 2.5);
   }

   #[test]
   fn test_sm2_easiness_factor_bounds() {
       let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());

       // Many failures
       for _ in 0..20 {
           schedule.update_after_review(0);
       }

       assert!(schedule.easiness_factor >= 1.3);
       assert!(schedule.easiness_factor <= 2.5);
   }
   ```
2. **GREEN**: Implement SM-2 algorithm
3. **REFACTOR**: Optimize calculation

#### Task 2.3.2: Vocabulary Extraction
**Priority**: P1
**Estimated Time**: 6 hours

#### Task 2.3.3: Review ViewSet
**Priority**: P1
**Estimated Time**: 4 hours

---

### Epic 2.4: Subscription & Payments

#### Task 2.4.1: Subscription Model & Migration
**Priority**: P1
**Estimated Time**: 3 hours

#### Task 2.4.2: Stripe Integration
**Priority**: P1
**Estimated Time**: 10 hours

#### Task 2.4.3: Usage Tracking & Quota Enforcement
**Priority**: P1
**Estimated Time**: 6 hours

---

## Phase 3: Extended Features (P2)

### Epic 3.1: Teacher Mode (WebSocket)

#### Task 3.1.1: WebSocket Handler Setup
**Priority**: P2
**Estimated Time**: 6 hours

**Test File**: `tests/integration/teacher_mode_tests.rs`
**Implementation File**: `src/apps/teacher_mode/websocket.rs`

#### Task 3.1.2: Session State Management
**Priority**: P2
**Estimated Time**: 6 hours

#### Task 3.1.3: Audio Streaming
**Priority**: P2
**Estimated Time**: 8 hours

#### Task 3.1.4: Command Handling (Pause/Resume/Skip)
**Priority**: P2
**Estimated Time**: 6 hours

---

## Continuous Tasks

### Task C.1: CI/CD Pipeline
**Priority**: P0
**Ongoing**: Throughout development

- Setup GitHub Actions workflow
- Run tests on every push
- Code coverage reporting
- Clippy linting
- Format checking

### Task C.2: Documentation
**Priority**: P1
**Ongoing**: Throughout development

- Update API docs after each endpoint
- Document Reinhardt patterns used
- Update architecture diagrams
- Maintain CHANGELOG.md

---

## Testing Standards

For every task:
1. Write tests FIRST (RED phase)
2. Run `cargo test` - verify test fails
3. Implement feature (GREEN phase)
4. Run `cargo test` - verify test passes
5. Refactor code (REFACTOR phase)
6. Run `cargo test` - verify still passes
7. Run `cargo clippy` - fix warnings
8. Run `cargo fmt` - format code
9. Commit with message: `feat: <description> (TDD)`

---

## Estimated Timeline

- **Phase 1 (MVP)**: 12-14 weeks
- **Phase 2 (Core)**: 8-10 weeks
- **Phase 3 (Extended)**: 8-10 weeks
- **Total**: 28-34 weeks (~7-8 months)

---

## References
- [hailango-requirements.md](./.aida/specs/hailango-requirements.md)
- [hailango-design.md](./.aida/specs/hailango-design.md)
- [test_strategy.md](../docs/testing/test_strategy.md)
