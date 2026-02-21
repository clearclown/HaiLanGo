//! API Integration Tests
//!
//! End-to-end tests using reinhardt-test's `spawn_test_server` and `APIClient`.
//! These tests run the full middleware stack (CORS, Logging, JWT) against a real
//! HTTP server on a random local port.
//!
//! Run with: cargo test --test api_integration

use std::sync::Arc;

use reinhardt_test::client::APIClient;
use reinhardt_test::server::{shutdown_test_server, spawn_test_server};
use serde_json::json;

use hailango::MiddlewareChain;
use hailango::api::middleware::{CorsMiddleware, JwtAuthMiddleware, LoggingMiddleware};
use hailango::config::urls::configure_urls;

/// Build the API application with full middleware stack (no DB state).
fn build_test_app() -> MiddlewareChain {
    let router = configure_urls();
    MiddlewareChain::new(Arc::new(router))
        .with_middleware(Arc::new(LoggingMiddleware::new()))
        .with_middleware(Arc::new(CorsMiddleware::permissive()))
        .with_middleware(Arc::new(JwtAuthMiddleware::new()))
}

// ---------------------------------------------------------------------------
// Auth API
// ---------------------------------------------------------------------------

#[tokio::test]
async fn test_auth_register() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client
        .post(
            "/api/auth/register/",
            &json!({"email": "user@example.com", "password": "pass1234", "display_name": "User"}),
            "json",
        )
        .await
        .unwrap();

    assert_eq!(response.status_code(), 201, "register should return 201");
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_auth_login_invalid_credentials() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client
        .post(
            "/api/auth/login/",
            &json!({"email": "nobody@example.com", "password": "wrong"}),
            "json",
        )
        .await
        .unwrap();

    // Login with unknown user returns 404 (user not found)
    assert_eq!(
        response.status_code(),
        404,
        "unknown user should return 404"
    );
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_auth_oauth_providers() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client.get("/api/auth/oauth/providers/").await.unwrap();
    assert_eq!(response.status_code(), 200);

    let body = response.json_value().unwrap();
    assert!(
        body.is_array() || body.is_object(),
        "providers should be array or object"
    );
    shutdown_test_server(handle).await;
}

// ---------------------------------------------------------------------------
// Books API
// ---------------------------------------------------------------------------

#[tokio::test]
async fn test_books_list_empty() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client.get("/api/books/").await.unwrap();
    assert_eq!(response.status_code(), 200);
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_books_create() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client
        .post(
            "/api/books/",
            &json!({"title": "Integration Test Book", "source_language": "en", "target_language": "ja"}),
            "json",
        )
        .await
        .unwrap();

    assert_eq!(response.status_code(), 201, "create book should return 201");
    let body = response.json_value().unwrap();
    assert_eq!(body["title"], "Integration Test Book");
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_books_create_invalid() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client
        .post(
            "/api/books/",
            &json!({"title": "", "source_language": "en", "target_language": "ja"}),
            "json",
        )
        .await
        .unwrap();

    assert_eq!(response.status_code(), 400, "empty title should return 400");
    shutdown_test_server(handle).await;
}

// ---------------------------------------------------------------------------
// Learning API
// ---------------------------------------------------------------------------

#[tokio::test]
async fn test_learning_sessions_list() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client.get("/api/learning/sessions/").await.unwrap();
    assert_eq!(response.status_code(), 200);
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_learning_session_create() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client
        .post(
            "/api/learning/sessions/",
            &json!({"book_id": "550e8400-e29b-41d4-a716-446655440000", "session_type": "page_by_page"}),
            "json",
        )
        .await
        .unwrap();

    assert_eq!(response.status_code(), 201);
    shutdown_test_server(handle).await;
}

// ---------------------------------------------------------------------------
// Review API
// ---------------------------------------------------------------------------

#[tokio::test]
async fn test_review_vocabulary_list() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client.get("/api/review/vocabulary/").await.unwrap();
    assert_eq!(response.status_code(), 200);
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_review_queue() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client.get("/api/review/queue/").await.unwrap();
    assert_eq!(response.status_code(), 200);
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_review_stats() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client.get("/api/review/stats/").await.unwrap();
    assert_eq!(response.status_code(), 200);
    shutdown_test_server(handle).await;
}

// ---------------------------------------------------------------------------
// TTS API
// ---------------------------------------------------------------------------

#[tokio::test]
async fn test_tts_languages() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client.get("/api/tts/languages/").await.unwrap();
    assert_eq!(response.status_code(), 200);

    let body = response.json_value().unwrap();
    assert!(body.is_array() || body.is_object());
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_tts_synthesize() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client
        .post(
            "/api/tts/synthesize/",
            &json!({"text": "Hello, integration test!", "language": "en"}),
            "json",
        )
        .await
        .unwrap();

    assert_eq!(response.status_code(), 200);
    let body = response.json_value().unwrap();
    assert!(
        body["audio_base64"].is_string(),
        "response should contain audio_base64"
    );
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_tts_synthesize_empty_text() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client
        .post(
            "/api/tts/synthesize/",
            &json!({"text": "", "language": "en"}),
            "json",
        )
        .await
        .unwrap();

    assert_eq!(response.status_code(), 400);
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_tts_history() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client.get("/api/tts/history/").await.unwrap();
    assert_eq!(response.status_code(), 200);
    shutdown_test_server(handle).await;
}

// ---------------------------------------------------------------------------
// Teacher Mode API
// ---------------------------------------------------------------------------

#[tokio::test]
async fn test_teacher_sessions_list() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client.get("/api/teacher/sessions/").await.unwrap();
    assert_eq!(response.status_code(), 200);
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_teacher_start_lesson() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client
        .post(
            "/api/teacher/start/",
            &json!({"book_id": "00000000-0000-0000-0000-000000000001", "start_page": 1, "end_page": 10}),
            "json",
        )
        .await
        .unwrap();

    assert_eq!(response.status_code(), 201);
    shutdown_test_server(handle).await;
}

#[tokio::test]
async fn test_teacher_status_not_found() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    // No active session → should return 404
    let response = client.get("/api/teacher/status/").await.unwrap();
    assert_eq!(response.status_code(), 404);
    shutdown_test_server(handle).await;
}

// ---------------------------------------------------------------------------
// CORS middleware
// ---------------------------------------------------------------------------

#[tokio::test]
async fn test_cors_headers_present() {
    let (url, handle) = spawn_test_server(Arc::new(build_test_app())).await;
    let client = APIClient::with_base_url(&url);

    let response = client.get("/api/books/").await.unwrap();
    assert_eq!(response.status_code(), 200);

    // CorsMiddleware::permissive() adds Access-Control-Allow-Origin
    let headers = response.headers();
    assert!(
        headers.contains_key("access-control-allow-origin"),
        "CORS header should be present"
    );
    shutdown_test_server(handle).await;
}
