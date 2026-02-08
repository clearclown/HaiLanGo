//! Teacher Mode API routes

use axum::{
    Router,
    extract::State,
    http::StatusCode,
    response::{IntoResponse, Json},
    routing::{get, post, put},
};
use serde_json::json;
use std::sync::{Arc, RwLock};
use uuid::Uuid;

use crate::apps::teacher_mode::{
    dto::{StartLessonRequest, UpdateConfigRequest},
    models::TeacherSession,
    views::{TeacherActionResult, TeacherModeViewSet},
};

/// Shared Teacher Mode state
#[derive(Clone)]
pub struct TeacherState {
    pub sessions: Arc<RwLock<Vec<TeacherSession>>>,
}

impl Default for TeacherState {
    fn default() -> Self {
        Self {
            sessions: Arc::new(RwLock::new(Vec::new())),
        }
    }
}

/// Helper to convert TeacherActionResult into HTTP response
fn action_response(result: TeacherActionResult) -> impl IntoResponse {
    match result {
        TeacherActionResult::Started(resp) => (StatusCode::CREATED, Json(json!(resp))),
        TeacherActionResult::Updated(resp) => (StatusCode::OK, Json(json!(resp))),
        TeacherActionResult::Status(resp) => (StatusCode::OK, Json(json!(resp))),
        TeacherActionResult::Sessions(list) => (StatusCode::OK, Json(json!(list))),
        TeacherActionResult::InvalidInput(msg) => {
            (StatusCode::BAD_REQUEST, Json(json!({"error": msg})))
        }
        TeacherActionResult::NotFound(msg) => {
            (StatusCode::NOT_FOUND, Json(json!({"error": msg})))
        }
        TeacherActionResult::InvalidState(msg) => {
            (StatusCode::CONFLICT, Json(json!({"error": msg})))
        }
    }
}

/// POST /api/teacher/start
async fn start_lesson(
    State(state): State<TeacherState>,
    Json(request): Json<StartLessonRequest>,
) -> impl IntoResponse {
    // Mock user_id (in production, extract from JWT)
    let user_id = Uuid::new_v4();
    let mut sessions = state.sessions.write().unwrap();
    let result = TeacherModeViewSet::start_lesson(user_id, request, &mut sessions);
    action_response(result)
}

/// POST /api/teacher/pause
async fn pause_lesson(State(state): State<TeacherState>) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let mut sessions = state.sessions.write().unwrap();
    let result = TeacherModeViewSet::pause(user_id, &mut sessions);
    action_response(result)
}

/// POST /api/teacher/resume
async fn resume_lesson(State(state): State<TeacherState>) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let mut sessions = state.sessions.write().unwrap();
    let result = TeacherModeViewSet::resume(user_id, &mut sessions);
    action_response(result)
}

/// POST /api/teacher/stop
async fn stop_lesson(State(state): State<TeacherState>) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let mut sessions = state.sessions.write().unwrap();
    let result = TeacherModeViewSet::stop(user_id, &mut sessions);
    action_response(result)
}

/// POST /api/teacher/next
async fn next_page(State(state): State<TeacherState>) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let mut sessions = state.sessions.write().unwrap();
    let result = TeacherModeViewSet::next_page(user_id, &mut sessions);
    action_response(result)
}

/// PUT /api/teacher/config
async fn update_config(
    State(state): State<TeacherState>,
    Json(request): Json<UpdateConfigRequest>,
) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let mut sessions = state.sessions.write().unwrap();
    let result = TeacherModeViewSet::update_config(user_id, request, &mut sessions);
    action_response(result)
}

/// GET /api/teacher/status
async fn get_status(State(state): State<TeacherState>) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let sessions = state.sessions.read().unwrap();
    let result = TeacherModeViewSet::get_status(user_id, &sessions);
    action_response(result)
}

/// GET /api/teacher/sessions
async fn list_sessions(State(state): State<TeacherState>) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let sessions = state.sessions.read().unwrap();
    let result = TeacherModeViewSet::list_sessions(user_id, &sessions);
    action_response(result)
}

/// Create Teacher Mode router
pub fn router() -> Router {
    let state = TeacherState::default();

    Router::new()
        .route("/start", post(start_lesson))
        .route("/pause", post(pause_lesson))
        .route("/resume", post(resume_lesson))
        .route("/stop", post(stop_lesson))
        .route("/next", post(next_page))
        .route("/config", put(update_config))
        .route("/status", get(get_status))
        .route("/sessions", get(list_sessions))
        .with_state(state)
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::{body::Body, http::Request};
    use tower::ServiceExt;

    fn start_body(end_page: u32) -> String {
        format!(
            r#"{{"book_id":"00000000-0000-0000-0000-000000000001","start_page":1,"end_page":{}}}"#,
            end_page
        )
    }

    #[tokio::test]
    async fn test_start_lesson_success() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/start")
                    .header("content-type", "application/json")
                    .body(Body::from(start_body(10)))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::CREATED);
    }

    #[tokio::test]
    async fn test_start_lesson_invalid_range() {
        let app = router();

        let body = r#"{"book_id":"00000000-0000-0000-0000-000000000001","start_page":10,"end_page":5}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/start")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_start_lesson_with_config() {
        let app = router();

        let body = r#"{"book_id":"00000000-0000-0000-0000-000000000001","start_page":1,"end_page":10,"language":"ja","speed":1.5,"page_interval_secs":10,"repeat_count":2}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/start")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::CREATED);
    }

    #[tokio::test]
    async fn test_get_sessions_empty() {
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
    async fn test_get_status_not_found() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/status")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        // New user_id each time, so no session found
        assert_eq!(response.status(), StatusCode::NOT_FOUND);
    }

    #[tokio::test]
    async fn test_pause_not_found() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/pause")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::NOT_FOUND);
    }

    #[tokio::test]
    async fn test_invalid_json() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/start")
                    .header("content-type", "application/json")
                    .body(Body::from("not-json"))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }
}
