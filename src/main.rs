use std::net::SocketAddr;

use axum::{Router, extract::State, http::StatusCode, response::Json, routing::get};
use serde::Serialize;
use tower_http::trace::TraceLayer;

use hailango::api;

/// Application state shared across handlers
#[derive(Clone)]
pub struct AppState {
    pub app_name: String,
    pub version: String,
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
    Json(serde_json::json!({
        "app": "HaiLanGo",
        "version": env!("CARGO_PKG_VERSION"),
        "description": "AI-powered language learning platform",
        "endpoints": {
            "auth": "/api/auth",
            "books": "/api/books",
            "learning": "/api/learning",
            "review": "/api/review"
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

/// Readiness check (for Kubernetes)
async fn ready() -> StatusCode {
    // TODO: Check database and redis connections
    StatusCode::OK
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
        .layer(TraceLayer::new_for_http())
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Initialize logging
    tracing_subscriber::fmt::init();

    tracing::info!("Starting HaiLanGo server...");

    // Create app state
    let state = AppState {
        app_name: "HaiLanGo".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
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
}
