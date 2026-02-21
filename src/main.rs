use std::net::SocketAddr;

use async_trait::async_trait;
use serde::Serialize;
use serde_json::json;

use hailango::{
    DefaultRouter, Handler, MiddlewareChain, Request, Response, Result, Router, StatusCode, path,
};

#[cfg(test)]
use bytes::Bytes;
#[cfg(test)]
use hailango::Method;
use hailango::api::middleware::{CorsMiddleware, JwtAuthMiddleware, LoggingMiddleware};
use hailango::config::urls::configure_urls;

/// Application state shared across handlers
#[derive(Clone)]
pub struct AppState {
    pub app_name: String,
    pub version: String,
    pub db_pool: Option<sqlx::PgPool>,
    pub redis_url: Option<String>,
}

/// Health check response
#[derive(Serialize)]
pub struct HealthResponse {
    pub status: String,
    pub app: String,
    pub version: String,
}

/// Root endpoint handler - API info
struct RootHandler;

#[async_trait]
impl Handler for RootHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        Response::ok().with_json(&json!({
            "app": "HaiLanGo",
            "version": env!("CARGO_PKG_VERSION"),
            "description": "AI-powered language learning platform",
            "endpoints": {
                "auth": "/api/auth",
                "books": "/api/books",
                "learning": "/api/learning",
                "review": "/api/review",
                "tts": "/api/tts",
                "teacher": "/api/teacher"
            }
        }))
    }
}

/// Health check handler
struct HealthHandler {
    state: AppState,
}

#[async_trait]
impl Handler for HealthHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        Response::ok().with_json(&HealthResponse {
            status: "healthy".to_string(),
            app: self.state.app_name.clone(),
            version: self.state.version.clone(),
        })
    }
}

/// Check database connectivity
async fn check_db(pool: &sqlx::PgPool) -> bool {
    sqlx::query("SELECT 1").fetch_one(pool).await.is_ok()
}

/// Check Redis connectivity
async fn check_redis(url: &str) -> bool {
    let client = match redis::Client::open(url) {
        Ok(c) => c,
        Err(_) => return false,
    };
    match redis::aio::ConnectionManager::new(client).await {
        Ok(mut conn) => {
            use redis::AsyncCommands;
            let result: std::result::Result<String, _> = conn.get("__healthcheck__").await;
            result.is_ok() || result.is_err()
        }
        Err(_) => false,
    }
}

/// Readiness check handler (for Kubernetes)
struct ReadyHandler {
    state: AppState,
}

#[async_trait]
impl Handler for ReadyHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        let mut checks = serde_json::Map::new();
        let mut any_unhealthy = false;

        // Database check
        match &self.state.db_pool {
            Some(pool) => {
                let healthy = check_db(pool).await;
                checks.insert(
                    "database".to_string(),
                    json!(if healthy { "healthy" } else { "unhealthy" }),
                );
                if !healthy {
                    any_unhealthy = true;
                }
            }
            None => {
                checks.insert("database".to_string(), json!("unconfigured"));
            }
        }

        // Redis check
        match &self.state.redis_url {
            Some(url) => {
                let healthy = check_redis(url).await;
                checks.insert(
                    "redis".to_string(),
                    json!(if healthy { "healthy" } else { "unhealthy" }),
                );
                if !healthy {
                    any_unhealthy = true;
                }
            }
            None => {
                checks.insert("redis".to_string(), json!("unconfigured"));
            }
        }

        let status = if any_unhealthy { "not_ready" } else { "ready" };
        let code = if any_unhealthy {
            StatusCode::SERVICE_UNAVAILABLE
        } else {
            StatusCode::OK
        };

        Response::new(code).with_json(&json!({"status": status, "checks": checks}))
    }
}

/// Create the application router with all routes
pub fn create_app(state: AppState) -> DefaultRouter {
    let mut router = configure_urls();

    // Add root-level routes
    router.add_route(path("/", RootHandler));
    router.add_route(path(
        "/health/",
        HealthHandler {
            state: state.clone(),
        },
    ));
    router.add_route(path("/ready/", ReadyHandler { state }));

    router
}

/// Wrap the router in the Reinhardt MiddlewareChain.
///
/// Middleware executes in insertion order (outermost → innermost):
/// 1. `LoggingMiddleware` — structured request/response logging
/// 2. `CorsMiddleware`    — permissive CORS headers (restrict in production)
/// 3. `JwtAuthMiddleware` — optional JWT validation for `/api/*` routes
pub fn create_middleware_stack(router: DefaultRouter) -> MiddlewareChain {
    use std::sync::Arc;

    MiddlewareChain::new(Arc::new(router))
        .with_middleware(Arc::new(LoggingMiddleware::new()))
        .with_middleware(Arc::new(CorsMiddleware::permissive()))
        .with_middleware(Arc::new(JwtAuthMiddleware::new()))
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Initialize logging
    tracing_subscriber::fmt::init();

    tracing::info!("Starting HaiLanGo server...");

    // Initialize database pool (optional)
    let db_pool = match std::env::var("DATABASE_URL") {
        Ok(url) => {
            use sqlx::postgres::PgPoolOptions;
            match PgPoolOptions::new()
                .max_connections(10)
                .min_connections(2)
                .acquire_timeout(std::time::Duration::from_secs(10))
                .connect(&url)
                .await
            {
                Ok(pool) => {
                    tracing::info!("Database connected");
                    Some(pool)
                }
                Err(e) => {
                    tracing::warn!("Database connection failed: {}", e);
                    None
                }
            }
        }
        Err(_) => {
            tracing::info!("DATABASE_URL not set, running without database");
            None
        }
    };

    // Redis URL (optional, checked at readiness time)
    let redis_url = std::env::var("REDIS_URL").ok();
    if redis_url.is_some() {
        tracing::info!("Redis configured");
    } else {
        tracing::info!("REDIS_URL not set, running without Redis");
    }

    // Create app state
    let state = AppState {
        app_name: "HaiLanGo".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
        db_pool,
        redis_url,
    };

    // Create router and wrap in middleware stack
    let router = create_app(state);
    let app = create_middleware_stack(router);

    // Bind to address and start Reinhardt HTTP server
    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    tracing::info!("Listening on {}", addr);

    use reinhardt::server::HttpServer;
    HttpServer::new(app)
        .listen(addr)
        .await
        .map_err(|e| anyhow::anyhow!("{}", e))?;

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    fn test_state() -> AppState {
        AppState {
            app_name: "HaiLanGo".to_string(),
            version: "0.1.0".to_string(),
            db_pool: None,
            redis_url: None,
        }
    }

    fn result_to_response(result: Result<Response>) -> Response {
        match result {
            Ok(r) => r,
            Err(e) => Response::from(e),
        }
    }

    #[tokio::test]
    async fn test_root_endpoint() {
        let app = create_app(test_state());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_health_endpoint() {
        let app = create_app(test_state());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/health/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_ready_endpoint() {
        let app = create_app(test_state());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/ready/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        // No DB/Redis configured -> unconfigured -> OK
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_api_auth_register() {
        let app = create_app(test_state());

        let body = r#"{"email":"test@example.com","password":"password123","display_name":"Test"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/api/auth/register/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        assert_eq!(response.status, 201);
    }

    #[tokio::test]
    async fn test_api_auth_oauth_providers() {
        let app = create_app(test_state());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/api/auth/oauth/providers/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_api_books_list() {
        let app = create_app(test_state());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/api/books/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_api_learning_sessions() {
        let app = create_app(test_state());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/api/learning/sessions/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_api_review_stats() {
        let app = create_app(test_state());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/api/review/stats/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_api_tts_synthesize() {
        let app = create_app(test_state());

        let body = r#"{"text":"Hello world","language":"en"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/api/tts/synthesize/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_api_teacher_start() {
        let app = create_app(test_state());

        let body =
            r#"{"book_id":"00000000-0000-0000-0000-000000000001","start_page":1,"end_page":10}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/api/teacher/start/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        assert_eq!(response.status, 201);
    }

    #[tokio::test]
    async fn test_api_teacher_sessions() {
        let app = create_app(test_state());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/api/teacher/sessions/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_api_tts_languages() {
        let app = create_app(test_state());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/api/tts/languages/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(app.handle(request).await);
        assert_eq!(response.status, 200);
    }
}
