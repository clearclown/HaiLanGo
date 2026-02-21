//! Learning API routes

use async_trait::async_trait;
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
use crate::{Handler, Method, Request, Response, Result, Route, path};

/// Simulated session store
#[derive(Clone, Default)]
pub struct LearningState {
    pub sessions: Arc<RwLock<Vec<LearningSession>>>,
}

/// Handler for GET/POST /sessions/
struct SessionsListHandler {
    state: LearningState,
}

#[async_trait]
impl Handler for SessionsListHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        match request.method {
            Method::GET => self.list(),
            Method::POST => self.create(request),
            _ => Err(crate::Error::MethodNotAllowed(
                "Only GET and POST are allowed".into(),
            )),
        }
    }
}

impl SessionsListHandler {
    fn list(&self) -> Result<Response> {
        let user_id = Uuid::new_v4();
        let sessions = self.state.sessions.read().unwrap();

        match LearningViewSet::list(&sessions, user_id) {
            ListSessionResult::Success(response) => Response::ok().with_json(&response),
            ListSessionResult::Unauthorized => {
                Response::unauthorized().with_json(&json!({"error": "Unauthorized"}))
            }
        }
    }

    fn create(&self, request: Request) -> Result<Response> {
        let req: CreateSessionRequest = request.json()?;
        let user_id = Uuid::new_v4();
        let book_exists = true;

        match LearningViewSet::create(req, user_id, book_exists) {
            CreateSessionResult::Success(response) => {
                let session =
                    LearningSession::new(user_id, response.book_id.unwrap(), response.session_type);
                self.state.sessions.write().unwrap().push(session);
                Response::created().with_json(&response)
            }
            CreateSessionResult::BookNotFound => {
                Response::not_found().with_json(&json!({"error": "Book not found"}))
            }
            CreateSessionResult::InvalidPageRange(msg) => {
                Response::bad_request().with_json(&json!({"error": msg}))
            }
        }
    }
}

/// Handler for GET /sessions/{id}/
struct SessionDetailHandler {
    state: LearningState,
}

#[async_trait]
impl Handler for SessionDetailHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let id = parse_uuid_param(&request, "id")?;
        let user_id = Uuid::new_v4();
        let sessions = self.state.sessions.read().unwrap();
        let session = sessions.iter().find(|s| s.id == id);

        match LearningViewSet::retrieve(session, user_id) {
            Some(response) => Response::ok().with_json(&response),
            None => Response::not_found().with_json(&json!({"error": "Session not found"})),
        }
    }
}

/// Handler for PATCH /sessions/{id}/status/
struct SessionStatusHandler {
    state: LearningState,
}

#[async_trait]
impl Handler for SessionStatusHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let id = parse_uuid_param(&request, "id")?;
        let req: UpdateSessionStatusRequest = request.json()?;
        let user_id = Uuid::new_v4();
        let mut sessions = self.state.sessions.write().unwrap();
        let session = sessions.iter_mut().find(|s| s.id == id);

        match LearningViewSet::update_status(req, session, user_id) {
            UpdateSessionResult::Success(response) => Response::ok().with_json(&response),
            UpdateSessionResult::SessionNotFound => {
                Response::not_found().with_json(&json!({"error": "Session not found"}))
            }
            UpdateSessionResult::InvalidAction(msg) => {
                Response::bad_request().with_json(&json!({"error": msg}))
            }
        }
    }
}

/// Handler for POST /sessions/{id}/progress/
struct SessionProgressHandler {
    state: LearningState,
}

#[async_trait]
impl Handler for SessionProgressHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let id = parse_uuid_param(&request, "id")?;
        let req: UpdateProgressRequest = request.json()?;
        let user_id = Uuid::new_v4();
        let sessions = self.state.sessions.read().unwrap();
        let session = sessions.iter().find(|s| s.id == id);
        let page_exists = true;

        match LearningViewSet::record_progress(req, session, user_id, page_exists) {
            UpdateProgressResult::Success(response) => Response::created().with_json(&response),
            UpdateProgressResult::SessionNotFound => {
                Response::not_found().with_json(&json!({"error": "Session not found"}))
            }
            UpdateProgressResult::PageNotFound => {
                Response::not_found().with_json(&json!({"error": "Page not found"}))
            }
            UpdateProgressResult::InvalidInput(msg) => {
                Response::bad_request().with_json(&json!({"error": msg}))
            }
        }
    }
}

/// Parse a UUID from path parameters
fn parse_uuid_param(request: &Request, name: &str) -> std::result::Result<Uuid, crate::Error> {
    request
        .path_params
        .get(name)
        .ok_or_else(|| crate::Error::Validation(format!("Missing {} parameter", name)))?
        .parse()
        .map_err(|_| crate::Error::Validation("Invalid UUID".into()))
}

/// Create learning routes
pub fn routes() -> Vec<Route> {
    let state = LearningState::default();

    vec![
        path(
            "/sessions/",
            SessionsListHandler {
                state: state.clone(),
            },
        ),
        path(
            "/sessions/{id}/",
            SessionDetailHandler {
                state: state.clone(),
            },
        ),
        path(
            "/sessions/{id}/status/",
            SessionStatusHandler {
                state: state.clone(),
            },
        ),
        path("/sessions/{id}/progress/", SessionProgressHandler { state }),
    ]
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::Method;
    use bytes::Bytes;

    fn result_to_response(result: Result<Response>) -> Response {
        match result {
            Ok(r) => r,
            Err(e) => Response::from(e),
        }
    }

    #[tokio::test]
    async fn test_list_sessions_empty() {
        let state = LearningState::default();
        let handler = SessionsListHandler { state };

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
    async fn test_create_session() {
        let state = LearningState::default();
        let handler = SessionsListHandler { state };

        let body =
            r#"{"book_id":"550e8400-e29b-41d4-a716-446655440000","session_type":"page_by_page"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/sessions/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 201);
    }
}
