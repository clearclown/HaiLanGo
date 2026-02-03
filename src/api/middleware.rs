//! API Middleware
//!
//! Authentication and authorization middleware for API routes.

use axum::{
    Json,
    body::Body,
    extract::Request,
    http::{StatusCode, header},
    middleware::Next,
    response::{IntoResponse, Response},
};
use serde_json::json;
use uuid::Uuid;

use crate::apps::auth::services::JwtService;

/// Extract user ID from Authorization header
pub async fn auth_middleware(request: Request, next: Next) -> Response {
    let auth_header = request
        .headers()
        .get(header::AUTHORIZATION)
        .and_then(|h| h.to_str().ok());

    let token = match auth_header {
        Some(header) if header.starts_with("Bearer ") => &header[7..],
        _ => {
            return (
                StatusCode::UNAUTHORIZED,
                Json(json!({"error": "Missing or invalid Authorization header"})),
            )
                .into_response();
        }
    };

    let jwt_service = JwtService::from_settings();
    match jwt_service.verify_token(token) {
        Ok(_claims) => next.run(request).await,
        Err(_) => (
            StatusCode::UNAUTHORIZED,
            Json(json!({"error": "Invalid or expired token"})),
        )
            .into_response(),
    }
}

/// Optional auth - extracts user if token present, but doesn't require it
pub async fn optional_auth_middleware(request: Request, next: Next) -> Response {
    // Just pass through - user extraction happens in handlers
    next.run(request).await
}

/// Extract user ID from request (for use in handlers after auth middleware)
pub fn extract_user_id(request: &Request<Body>) -> Option<Uuid> {
    let auth_header = request
        .headers()
        .get(header::AUTHORIZATION)
        .and_then(|h| h.to_str().ok())?;

    let token = auth_header.strip_prefix("Bearer ")?;

    let jwt_service = JwtService::from_settings();
    jwt_service.get_user_id(token).ok()
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::{Router, body::Body, http::Request, routing::get};
    use tower::ServiceExt;

    async fn test_handler() -> &'static str {
        "OK"
    }

    fn create_test_router() -> Router {
        Router::new().route(
            "/protected",
            get(test_handler).layer(axum::middleware::from_fn(auth_middleware)),
        )
    }

    #[tokio::test]
    async fn test_missing_auth_header() {
        let app = create_test_router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/protected")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_invalid_auth_header_format() {
        let app = create_test_router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/protected")
                    .header("Authorization", "Basic invalid")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_invalid_token() {
        let app = create_test_router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/protected")
                    .header("Authorization", "Bearer invalid.token.here")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_valid_token() {
        let app = create_test_router();

        // Generate a valid token
        let jwt_service = JwtService::from_settings();
        let user_id = Uuid::new_v4();
        let tokens = jwt_service
            .generate_tokens(user_id, "test@example.com")
            .unwrap();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/protected")
                    .header("Authorization", format!("Bearer {}", tokens.access_token))
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }
}
