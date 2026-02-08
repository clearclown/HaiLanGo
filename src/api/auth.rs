//! Auth API routes

use axum::{
    Router,
    extract::{Path, Query, State},
    http::StatusCode,
    response::{IntoResponse, Json},
    routing::{get, post},
};
use serde_json::json;
use std::sync::Arc;
use uuid::Uuid;

use crate::apps::auth::{
    dto::{LoginRequest, OAuthCallbackQuery, RegisterRequest},
    oauth::{OAuthProvider, OAuthService},
    views::{AuthViewSet, LoginResult, OAuthLoginResult, RegisterResult},
};

/// OAuth-enabled auth state
#[derive(Clone)]
pub struct AuthState {
    pub oauth_service: Arc<OAuthService>,
}

impl Default for AuthState {
    fn default() -> Self {
        Self {
            oauth_service: Arc::new(OAuthService::from_env()),
        }
    }
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

/// GET /api/auth/oauth/providers - List available OAuth providers
async fn oauth_providers(
    State(state): State<AuthState>,
) -> (StatusCode, Json<serde_json::Value>) {
    let providers: Vec<serde_json::Value> = state
        .oauth_service
        .configured_providers()
        .iter()
        .map(|p| json!({"name": p.as_str(), "configured": true}))
        .collect();
    (StatusCode::OK, Json(json!({"providers": providers})))
}

/// GET /api/auth/oauth/{provider} - Get OAuth authorization URL
async fn oauth_redirect(
    State(state): State<AuthState>,
    Path(provider_name): Path<String>,
) -> (StatusCode, Json<serde_json::Value>) {
    let provider = match OAuthProvider::parse(&provider_name) {
        Some(p) => p,
        None => {
            return (
                StatusCode::BAD_REQUEST,
                Json(json!({"error": format!("Unknown provider: {}", provider_name)})),
            )
        }
    };

    // Generate CSRF state token
    let state_token = Uuid::new_v4().to_string();

    match state
        .oauth_service
        .get_authorization_url(provider, &state_token)
    {
        Ok(url) => (
            StatusCode::OK,
            Json(json!({"auth_url": url, "state": state_token})),
        ),
        Err(e) => (
            StatusCode::BAD_REQUEST,
            Json(json!({"error": e.to_string()})),
        ),
    }
}

/// GET /api/auth/callback/{provider} - Handle OAuth callback
async fn oauth_callback(
    State(state): State<AuthState>,
    Path(provider_name): Path<String>,
    Query(query): Query<OAuthCallbackQuery>,
) -> (StatusCode, Json<serde_json::Value>) {
    let provider = match OAuthProvider::parse(&provider_name) {
        Some(p) => p,
        None => {
            return (
                StatusCode::BAD_REQUEST,
                Json(json!({"error": "Unknown provider"})),
            )
        }
    };

    // TODO: Verify state token against stored value for CSRF protection

    match state
        .oauth_service
        .authenticate(provider, &query.code)
        .await
    {
        Ok(user_info) => match AuthViewSet::oauth_login(user_info) {
            OAuthLoginResult::Success(response) => {
                (StatusCode::OK, Json(json!(response)))
            }
            OAuthLoginResult::ProviderError(msg) => (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(json!({"error": msg})),
            ),
        },
        Err(e) => (
            StatusCode::BAD_REQUEST,
            Json(json!({"error": e.to_string()})),
        ),
    }
}

/// Create auth router
pub fn router() -> Router {
    let state = AuthState::default();

    Router::new()
        .route("/register", post(register))
        .route("/login", post(login))
        .route("/oauth/providers", get(oauth_providers))
        .route("/oauth/{provider}", get(oauth_redirect))
        .route("/callback/{provider}", get(oauth_callback))
        .with_state(state)
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

    #[tokio::test]
    async fn test_oauth_providers_list() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/oauth/providers")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_oauth_redirect_unknown_provider() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/oauth/unknown")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_oauth_redirect_unconfigured_google() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/oauth/google")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        // Google is not configured in test env
        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_oauth_callback_unknown_provider() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/callback/invalid?code=abc&state=xyz")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }
}
