//! Learning API routes

use axum::{
    Router,
    extract::{Path, State},
    http::StatusCode,
    response::{IntoResponse, Json},
    routing::{get, patch, post},
};
use serde_json::json;
use std::sync::{Arc, RwLock};
use uuid::Uuid;

use crate::apps::learning::{
    dto::{CreateSessionRequest, UpdateProgressRequest, UpdateSessionStatusRequest},
    models::LearningSession,
    views::{
        CreateSessionResult, LearningViewSet, ListSessionResult, UpdateProgressResult,
        UpdateSessionResult,
    },
};

/// Simulated session store
#[derive(Clone, Default)]
pub struct LearningState {
    pub sessions: Arc<RwLock<Vec<LearningSession>>>,
}

/// GET /api/learning/sessions
async fn list_sessions(State(state): State<LearningState>) -> impl IntoResponse {
    let user_id = Uuid::new_v4(); // Mock user (from JWT in production)
    let sessions = state.sessions.read().unwrap();

    match LearningViewSet::list(&sessions, user_id) {
        ListSessionResult::Success(response) => (StatusCode::OK, Json(json!(response))),
        ListSessionResult::Unauthorized => (
            StatusCode::UNAUTHORIZED,
            Json(json!({"error": "Unauthorized"})),
        ),
    }
}

/// POST /api/learning/sessions
async fn create_session(
    State(state): State<LearningState>,
    Json(request): Json<CreateSessionRequest>,
) -> impl IntoResponse {
    let user_id = Uuid::new_v4();

    // Mock book exists check (in production, query database)
    let book_exists = true;

    match LearningViewSet::create(request, user_id, book_exists) {
        CreateSessionResult::Success(response) => {
            // Store session
            let session =
                LearningSession::new(user_id, response.book_id.unwrap(), response.session_type);
            state.sessions.write().unwrap().push(session);
            (StatusCode::CREATED, Json(json!(response)))
        }
        CreateSessionResult::BookNotFound => (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "Book not found"})),
        ),
        CreateSessionResult::InvalidPageRange(msg) => {
            (StatusCode::BAD_REQUEST, Json(json!({"error": msg})))
        }
    }
}

/// GET /api/learning/sessions/:id
async fn get_session(
    State(state): State<LearningState>,
    Path(id): Path<Uuid>,
) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let sessions = state.sessions.read().unwrap();
    let session = sessions.iter().find(|s| s.id == id);

    match LearningViewSet::retrieve(session, user_id) {
        Some(response) => (StatusCode::OK, Json(json!(response))),
        None => (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "Session not found"})),
        ),
    }
}

/// PATCH /api/learning/sessions/:id/status
async fn update_session_status(
    State(state): State<LearningState>,
    Path(id): Path<Uuid>,
    Json(request): Json<UpdateSessionStatusRequest>,
) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let mut sessions = state.sessions.write().unwrap();
    let session = sessions.iter_mut().find(|s| s.id == id);

    match LearningViewSet::update_status(request, session, user_id) {
        UpdateSessionResult::Success(response) => (StatusCode::OK, Json(json!(response))),
        UpdateSessionResult::SessionNotFound => (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "Session not found"})),
        ),
        UpdateSessionResult::InvalidAction(msg) => {
            (StatusCode::BAD_REQUEST, Json(json!({"error": msg})))
        }
    }
}

/// POST /api/learning/sessions/:id/progress
async fn record_progress(
    State(state): State<LearningState>,
    Path(id): Path<Uuid>,
    Json(request): Json<UpdateProgressRequest>,
) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let sessions = state.sessions.read().unwrap();
    let session = sessions.iter().find(|s| s.id == id);

    // Mock page exists check
    let page_exists = true;

    match LearningViewSet::record_progress(request, session, user_id, page_exists) {
        UpdateProgressResult::Success(response) => (StatusCode::CREATED, Json(json!(response))),
        UpdateProgressResult::SessionNotFound => (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "Session not found"})),
        ),
        UpdateProgressResult::PageNotFound => (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "Page not found"})),
        ),
        UpdateProgressResult::InvalidInput(msg) => {
            (StatusCode::BAD_REQUEST, Json(json!({"error": msg})))
        }
    }
}

/// Create learning router
pub fn router() -> Router {
    let state = LearningState::default();

    Router::new()
        .route("/sessions", get(list_sessions).post(create_session))
        .route("/sessions/{id}", get(get_session))
        .route("/sessions/{id}/status", patch(update_session_status))
        .route("/sessions/{id}/progress", post(record_progress))
        .with_state(state)
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::{body::Body, http::Request};
    use tower::ServiceExt;

    #[tokio::test]
    async fn test_list_sessions_empty() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/sessions")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_create_session() {
        let app = router();

        let body =
            r#"{"book_id":"550e8400-e29b-41d4-a716-446655440000","session_type":"page_by_page"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/sessions")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::CREATED);
    }
}
