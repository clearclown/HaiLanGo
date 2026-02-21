//! Teacher Mode API routes

use async_trait::async_trait;
use serde_json::json;
use std::sync::{Arc, RwLock};
use uuid::Uuid;

use crate::{Handler, Request, Response, Result, Route, StatusCode, path};
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

/// Convert TeacherActionResult into a Response
fn action_response(result: TeacherActionResult) -> Result<Response> {
    match result {
        TeacherActionResult::Started(resp) => Response::created().with_json(&resp),
        TeacherActionResult::Updated(resp) => Response::ok().with_json(&resp),
        TeacherActionResult::Status(resp) => Response::ok().with_json(&resp),
        TeacherActionResult::Sessions(list) => Response::ok().with_json(&list),
        TeacherActionResult::InvalidInput(msg) => {
            Response::bad_request().with_json(&json!({"error": msg}))
        }
        TeacherActionResult::NotFound(msg) => {
            Response::not_found().with_json(&json!({"error": msg}))
        }
        TeacherActionResult::InvalidState(msg) => {
            Response::new(StatusCode::CONFLICT).with_json(&json!({"error": msg}))
        }
    }
}

/// Handler for POST /start/
struct StartLessonHandler {
    state: TeacherState,
}

#[async_trait]
impl Handler for StartLessonHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let req: StartLessonRequest = request.json()?;
        let user_id = Uuid::new_v4();
        let mut sessions = self.state.sessions.write().unwrap();
        let result = TeacherModeViewSet::start_lesson(user_id, req, &mut sessions);
        action_response(result)
    }
}

/// Handler for POST /pause/
struct PauseLessonHandler {
    state: TeacherState,
}

#[async_trait]
impl Handler for PauseLessonHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        let user_id = Uuid::new_v4();
        let mut sessions = self.state.sessions.write().unwrap();
        let result = TeacherModeViewSet::pause(user_id, &mut sessions);
        action_response(result)
    }
}

/// Handler for POST /resume/
struct ResumeLessonHandler {
    state: TeacherState,
}

#[async_trait]
impl Handler for ResumeLessonHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        let user_id = Uuid::new_v4();
        let mut sessions = self.state.sessions.write().unwrap();
        let result = TeacherModeViewSet::resume(user_id, &mut sessions);
        action_response(result)
    }
}

/// Handler for POST /stop/
struct StopLessonHandler {
    state: TeacherState,
}

#[async_trait]
impl Handler for StopLessonHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        let user_id = Uuid::new_v4();
        let mut sessions = self.state.sessions.write().unwrap();
        let result = TeacherModeViewSet::stop(user_id, &mut sessions);
        action_response(result)
    }
}

/// Handler for POST /next/
struct NextPageHandler {
    state: TeacherState,
}

#[async_trait]
impl Handler for NextPageHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        let user_id = Uuid::new_v4();
        let mut sessions = self.state.sessions.write().unwrap();
        let result = TeacherModeViewSet::next_page(user_id, &mut sessions);
        action_response(result)
    }
}

/// Handler for PUT /config/
struct UpdateConfigHandler {
    state: TeacherState,
}

#[async_trait]
impl Handler for UpdateConfigHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let req: UpdateConfigRequest = request.json()?;
        let user_id = Uuid::new_v4();
        let mut sessions = self.state.sessions.write().unwrap();
        let result = TeacherModeViewSet::update_config(user_id, req, &mut sessions);
        action_response(result)
    }
}

/// Handler for GET /status/
struct GetStatusHandler {
    state: TeacherState,
}

#[async_trait]
impl Handler for GetStatusHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        let user_id = Uuid::new_v4();
        let sessions = self.state.sessions.read().unwrap();
        let result = TeacherModeViewSet::get_status(user_id, &sessions);
        action_response(result)
    }
}

/// Handler for GET /sessions/
struct ListSessionsHandler {
    state: TeacherState,
}

#[async_trait]
impl Handler for ListSessionsHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        let user_id = Uuid::new_v4();
        let sessions = self.state.sessions.read().unwrap();
        let result = TeacherModeViewSet::list_sessions(user_id, &sessions);
        action_response(result)
    }
}

/// Create Teacher Mode routes
pub fn routes() -> Vec<Route> {
    let state = TeacherState::default();

    vec![
        path("/start/", StartLessonHandler { state: state.clone() }),
        path("/pause/", PauseLessonHandler { state: state.clone() }),
        path(
            "/resume/",
            ResumeLessonHandler {
                state: state.clone(),
            },
        ),
        path("/stop/", StopLessonHandler { state: state.clone() }),
        path("/next/", NextPageHandler { state: state.clone() }),
        path(
            "/config/",
            UpdateConfigHandler {
                state: state.clone(),
            },
        ),
        path("/status/", GetStatusHandler { state: state.clone() }),
        path("/sessions/", ListSessionsHandler { state }),
    ]
}

#[cfg(test)]
mod tests {
    use super::*;
    use bytes::Bytes;
    use crate::Method;

    fn result_to_response(result: Result<Response>) -> Response {
        match result {
            Ok(r) => r,
            Err(e) => Response::from(e),
        }
    }

    fn start_body(end_page: u32) -> String {
        format!(
            r#"{{"book_id":"00000000-0000-0000-0000-000000000001","start_page":1,"end_page":{}}}"#,
            end_page
        )
    }

    fn make_state() -> TeacherState {
        TeacherState::default()
    }

    #[tokio::test]
    async fn test_start_lesson_success() {
        let handler = StartLessonHandler { state: make_state() };

        let request = Request::builder()
            .method(Method::POST)
            .uri("/start/")
            .header("content-type", "application/json")
            .body(Bytes::from(start_body(10)))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 201);
    }

    #[tokio::test]
    async fn test_start_lesson_invalid_range() {
        let handler = StartLessonHandler { state: make_state() };

        let body =
            r#"{"book_id":"00000000-0000-0000-0000-000000000001","start_page":10,"end_page":5}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/start/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 400);
    }

    #[tokio::test]
    async fn test_start_lesson_with_config() {
        let handler = StartLessonHandler { state: make_state() };

        let body = r#"{"book_id":"00000000-0000-0000-0000-000000000001","start_page":1,"end_page":10,"language":"ja","speed":1.5,"page_interval_secs":10,"repeat_count":2}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/start/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 201);
    }

    #[tokio::test]
    async fn test_get_sessions_empty() {
        let handler = ListSessionsHandler { state: make_state() };

        let request = Request::builder()
            .method(Method::GET)
            .uri("/sessions/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_get_status_not_found() {
        let handler = GetStatusHandler { state: make_state() };

        let request = Request::builder()
            .method(Method::GET)
            .uri("/status/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 404);
    }

    #[tokio::test]
    async fn test_pause_not_found() {
        let handler = PauseLessonHandler { state: make_state() };

        let request = Request::builder()
            .method(Method::POST)
            .uri("/pause/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 404);
    }

    #[tokio::test]
    async fn test_invalid_json() {
        let handler = StartLessonHandler { state: make_state() };

        let request = Request::builder()
            .method(Method::POST)
            .uri("/start/")
            .header("content-type", "application/json")
            .body(Bytes::from("not-json"))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 400);
    }
}
