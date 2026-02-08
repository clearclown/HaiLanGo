use std::net::SocketAddr;

use axum::{Router, extract::State, http::StatusCode, response::Json, routing::get};
use serde::Serialize;
use serde_json::json;
use tower_http::trace::TraceLayer;

use hailango::api;

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

/// Root endpoint - API info
async fn root() -> Json<serde_json::Value> {
    Json(json!({
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

/// Health check endpoint
async fn health_check(State(state): State<AppState>) -> Json<HealthResponse> {
    Json(HealthResponse {
        status: "healthy".to_string(),
        app: state.app_name,
        version: state.version,
    })
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
            let result: Result<String, _> = conn.get("__healthcheck__").await;
            // GET on non-existent key returns Nil, which is a redis error, but connection works
            result.is_ok() || result.is_err()
        }
        Err(_) => false,
    }
}

/// Readiness check (for Kubernetes)
async fn ready(State(state): State<AppState>) -> (StatusCode, Json<serde_json::Value>) {
    let mut checks = serde_json::Map::new();
    let mut any_unhealthy = false;

    // Database check
    match &state.db_pool {
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
    match &state.redis_url {
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

    (code, Json(json!({"status": status, "checks": checks})))
}

/// Create the application router
pub fn create_app(state: AppState) -> Router {
    // Create base router with state
    let base = Router::new()
        .route("/", get(root))
        .route("/health", get(health_check))
        .route("/ready", get(ready))
        .with_state(state);

    // Nest API routers (each has its own state)
    base.nest("/api/auth", api::auth::router())
        .nest("/api/books", api::books::router())
        .nest("/api/learning", api::learning::router())
        .nest("/api/review", api::review::router())
        .nest("/api/tts", api::tts::router())
        .nest("/api/teacher", api::teacher::router())
        .layer(TraceLayer::new_for_http())
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

    // Create router
    let app = create_app(state);

    // Bind to address
    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    tracing::info!("Listening on {}", addr);

    // Start server
    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::{
        body::Body,
        http::{Request, StatusCode},
    };
    use tower::ServiceExt;

    fn test_state() -> AppState {
        AppState {
            app_name: "HaiLanGo".to_string(),
            version: "0.1.0".to_string(),
            db_pool: None,
            redis_url: None,
        }
    }

    #[tokio::test]
    async fn test_root_endpoint() {
        let app = create_app(test_state());

        let response = app
            .oneshot(Request::builder().uri("/").body(Body::empty()).unwrap())
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_health_endpoint() {
        let app = create_app(test_state());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/health")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_ready_endpoint() {
        let app = create_app(test_state());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/ready")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        // No DB/Redis configured -> unconfigured -> OK
        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_api_auth_register() {
        let app = create_app(test_state());

        let body = r#"{"email":"test@example.com","password":"password123","display_name":"Test"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/auth/register")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::CREATED);
    }

    #[tokio::test]
    async fn test_api_auth_oauth_providers() {
        let app = create_app(test_state());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/api/auth/oauth/providers")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_api_books_list() {
        let app = create_app(test_state());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/api/books")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_api_learning_sessions() {
        let app = create_app(test_state());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/api/learning/sessions")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_api_review_stats() {
        let app = create_app(test_state());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/api/review/stats")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_api_tts_synthesize() {
        let app = create_app(test_state());

        let body = r#"{"text":"Hello world","language":"en"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/tts/synthesize")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_api_teacher_start() {
        let app = create_app(test_state());

        let body = r#"{"book_id":"00000000-0000-0000-0000-000000000001","start_page":1,"end_page":10}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/teacher/start")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::CREATED);
    }

    #[tokio::test]
    async fn test_api_teacher_sessions() {
        let app = create_app(test_state());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/api/teacher/sessions")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_api_tts_languages() {
        let app = create_app(test_state());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/api/tts/languages")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }
}
