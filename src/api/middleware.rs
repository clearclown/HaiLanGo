//! Reinhardt Middleware for HaiLanGo API
//!
//! Provides custom middleware using the Reinhardt `Middleware` trait.
//! Built-in middleware (CORS, Logging) is re-exported from the reinhardt crate.

use async_trait::async_trait;
use std::sync::Arc;

use crate::apps::auth::services::JwtService;
use crate::{Handler, Middleware, Request, Response, Result};

// Re-export Reinhardt built-in middleware for use in main.rs
pub use reinhardt::CorsMiddleware;
pub use reinhardt::LoggingMiddleware;

/// JWT Bearer token authentication middleware.
///
/// Validates `Authorization: Bearer <token>` headers on API routes.
/// When a valid token is present, the claims are logged for traceability.
/// Missing or invalid tokens are logged but do not block requests at this stage —
/// per-handler enforcement is added in subsequent phases.
///
/// Only activates for `/api/*` paths (skips health checks and root).
pub struct JwtAuthMiddleware {
    jwt: JwtService,
}

impl JwtAuthMiddleware {
    pub fn new() -> Self {
        Self {
            jwt: JwtService::from_settings(),
        }
    }
}

impl Default for JwtAuthMiddleware {
    fn default() -> Self {
        Self::new()
    }
}

#[async_trait]
impl Middleware for JwtAuthMiddleware {
    async fn process(&self, request: Request, next: Arc<dyn Handler>) -> Result<Response> {
        if let Some(auth_header) = request.headers.get("authorization") {
            if let Ok(auth_str) = auth_header.to_str() {
                if let Some(token) = auth_str.strip_prefix("Bearer ") {
                    match self.jwt.verify_token(token) {
                        Ok(claims) => {
                            tracing::debug!(user_id = %claims.sub, "JWT authenticated");
                        }
                        Err(e) => {
                            tracing::warn!("Invalid JWT token: {}", e);
                        }
                    }
                }
            }
        }

        next.handle(request).await
    }

    /// Only run on API routes — skip health checks and root.
    fn should_continue(&self, request: &Request) -> bool {
        request.uri.path().starts_with("/api/")
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::Method;
    use bytes::Bytes;

    struct EchoHandler;

    #[async_trait]
    impl Handler for EchoHandler {
        async fn handle(&self, _request: Request) -> Result<Response> {
            Response::ok().with_json(&serde_json::json!({"ok": true}))
        }
    }

    fn make_request(uri: &str) -> Request {
        Request::builder()
            .method(Method::GET)
            .uri(uri)
            .body(Bytes::new())
            .build()
            .unwrap()
    }

    #[test]
    fn test_should_continue_api_paths() {
        let mw = JwtAuthMiddleware::new();
        assert!(mw.should_continue(&make_request("/api/books/")));
        assert!(mw.should_continue(&make_request("/api/auth/login/")));
        assert!(mw.should_continue(&make_request("/api/tts/synthesize/")));
    }

    #[test]
    fn test_should_not_continue_non_api_paths() {
        let mw = JwtAuthMiddleware::new();
        assert!(!mw.should_continue(&make_request("/")));
        assert!(!mw.should_continue(&make_request("/health/")));
        assert!(!mw.should_continue(&make_request("/ready/")));
    }

    #[tokio::test]
    async fn test_passthrough_no_auth_header() {
        let mw = JwtAuthMiddleware::new();
        let handler = Arc::new(EchoHandler);
        let request = make_request("/api/books/");
        let response = mw.process(request, handler).await.unwrap();
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_passthrough_invalid_token() {
        let mw = JwtAuthMiddleware::new();
        let handler = Arc::new(EchoHandler);

        let mut request = make_request("/api/books/");
        request.headers.insert(
            "authorization",
            "Bearer invalid.token.here".parse().unwrap(),
        );

        // Invalid token is logged but does not block (permissive phase)
        let response = mw.process(request, handler).await.unwrap();
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_passthrough_valid_token() {
        use uuid::Uuid;
        let jwt = JwtService::from_settings();
        let user_id = Uuid::new_v4();
        let tokens = jwt.generate_tokens(user_id, "test@example.com").unwrap();

        let mw = JwtAuthMiddleware::new();
        let handler = Arc::new(EchoHandler);

        let mut request = make_request("/api/books/");
        let auth_value = format!("Bearer {}", tokens.access_token);
        request
            .headers
            .insert("authorization", auth_value.parse().unwrap());

        let response = mw.process(request, handler).await.unwrap();
        assert_eq!(response.status, 200);
    }
}
