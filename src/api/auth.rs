//! Auth API routes

use axum::{
    Router,
    http::StatusCode,
    response::{IntoResponse, Json},
    routing::post,
};
use serde_json::json;

use crate::apps::auth::{
    dto::{LoginRequest, RegisterRequest},
    views::{AuthViewSet, LoginResult, RegisterResult},
};

/// Simulated user store (in production, use database)
#[derive(Clone, Default)]
pub struct AuthState {
    // In production: DatabasePool
}

/// POST /api/auth/register
async fn register(Json(request): Json<RegisterRequest>) -> impl IntoResponse {
    match AuthViewSet::register(request) {
        RegisterResult::Success(response) => (StatusCode::CREATED, Json(json!(response))),
        RegisterResult::EmailExists => (
            StatusCode::CONFLICT,
            Json(json!({"error": "Email already exists"})),
        ),
        RegisterResult::InvalidInput(msg) => (StatusCode::BAD_REQUEST, Json(json!({"error": msg}))),
    }
}

/// POST /api/auth/login
async fn login(Json(request): Json<LoginRequest>) -> impl IntoResponse {
    // In production: lookup user from database
    // For now, return unauthorized as we don't have persistent storage
    match AuthViewSet::login(request, None) {
        LoginResult::Success(response) => (StatusCode::OK, Json(json!(response))),
        LoginResult::InvalidCredentials => (
            StatusCode::UNAUTHORIZED,
            Json(json!({"error": "Invalid credentials"})),
        ),
        LoginResult::UserNotFound => (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "User not found"})),
        ),
    }
}

/// Create auth router
pub fn router() -> Router {
    Router::new()
        .route("/register", post(register))
        .route("/login", post(login))
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::{body::Body, http::Request};
    use tower::ServiceExt;

    #[tokio::test]
    async fn test_register_endpoint() {
        let app = router();

        let body =
            r#"{"email":"test@example.com","password":"password123","display_name":"Test User"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/register")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::CREATED);
    }

    #[tokio::test]
    async fn test_register_invalid_email() {
        let app = router();

        let body = r#"{"email":"invalid","password":"password123","display_name":"Test"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/register")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_login_user_not_found() {
        let app = router();

        let body = r#"{"email":"notfound@example.com","password":"password123"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/login")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::NOT_FOUND);
    }
}
