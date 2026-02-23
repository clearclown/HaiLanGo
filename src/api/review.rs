//! Review API routes

use async_trait::async_trait;
use serde::Deserialize;
use serde_json::json;
use std::sync::{Arc, RwLock};
use uuid::Uuid;

use crate::apps::review::{
    dto::{CreateVocabularyRequest, RecordReviewRequest},
    models::{SrsSchedule, Vocabulary},
    views::{CreateVocabularyResult, RecordReviewResult, ReviewQueueResult, ReviewViewSet},
};
use crate::{Handler, Method, Request, Response, Result, Route, StatusCode, path};

/// Simulated review store
#[derive(Clone, Default)]
pub struct ReviewState {
    pub vocabularies: Arc<RwLock<Vec<Vocabulary>>>,
    pub schedules: Arc<RwLock<Vec<SrsSchedule>>>,
}

/// Query parameters for review queue
#[derive(Debug, Deserialize)]
pub struct QueueQuery {
    #[serde(default = "default_limit")]
    pub limit: usize,
}

fn default_limit() -> usize {
    20
}

/// Handler for GET/POST /vocabulary/
struct VocabularyListHandler {
    state: ReviewState,
}

#[async_trait]
impl Handler for VocabularyListHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let user_id = request
            .extensions
            .get::<Uuid>()
            .unwrap_or_else(Uuid::new_v4);
        match request.method {
            Method::GET => self.list(user_id),
            Method::POST => self.create(user_id, request),
            _ => Err(crate::Error::MethodNotAllowed(
                "Only GET and POST are allowed".into(),
            )),
        }
    }
}

impl VocabularyListHandler {
    fn list(&self, user_id: Uuid) -> Result<Response> {
        let vocabularies = self.state.vocabularies.read().unwrap();
        let response = ReviewViewSet::list_vocabularies(&vocabularies, user_id);
        Response::ok().with_json(&response)
    }

    fn create(&self, user_id: Uuid, request: Request) -> Result<Response> {
        let req: CreateVocabularyRequest = request.json()?;
        let page_exists = true;

        let vocabularies = self.state.vocabularies.read().unwrap();
        let word_exists = vocabularies
            .iter()
            .any(|v| v.user_id == user_id && v.word == req.word);
        drop(vocabularies);

        match ReviewViewSet::create_vocabulary(req.clone(), user_id, page_exists, word_exists) {
            CreateVocabularyResult::Success(response) => {
                let vocab = Vocabulary::new(
                    req.page_id,
                    user_id,
                    response.word.clone(),
                    response.meaning.clone(),
                );
                let schedule = SrsSchedule::new(user_id, vocab.id);

                self.state.vocabularies.write().unwrap().push(vocab);
                self.state.schedules.write().unwrap().push(schedule);

                Response::created().with_json(&response)
            }
            CreateVocabularyResult::PageNotFound => {
                Response::not_found().with_json(&json!({"error": "Page not found"}))
            }
            CreateVocabularyResult::DuplicateWord => Response::new(StatusCode::CONFLICT)
                .with_json(&json!({"error": "Word already exists"})),
            CreateVocabularyResult::InvalidInput(msg) => {
                Response::bad_request().with_json(&json!({"error": msg}))
            }
        }
    }
}

/// Handler for GET /queue/
struct ReviewQueueHandler {
    state: ReviewState,
}

#[async_trait]
impl Handler for ReviewQueueHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let user_id = request
            .extensions
            .get::<Uuid>()
            .unwrap_or_else(Uuid::new_v4);
        let query: QueueQuery = request.query_as().unwrap_or(QueueQuery {
            limit: default_limit(),
        });
        let vocabularies = self.state.vocabularies.read().unwrap();
        let schedules = self.state.schedules.read().unwrap();

        match ReviewViewSet::get_review_queue(&vocabularies, &schedules, user_id, query.limit) {
            ReviewQueueResult::Success(response) => Response::ok().with_json(&response),
            ReviewQueueResult::Empty => Response::ok()
                .with_json(&json!({"items": [], "due_count": 0, "total_vocabulary": 0})),
        }
    }
}

/// Handler for POST /record/
struct RecordReviewHandler {
    state: ReviewState,
}

#[async_trait]
impl Handler for RecordReviewHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let user_id = request
            .extensions
            .get::<Uuid>()
            .unwrap_or_else(Uuid::new_v4);
        let req: RecordReviewRequest = request.json()?;
        let mut schedules = self.state.schedules.write().unwrap();
        let schedule = schedules
            .iter_mut()
            .find(|s| s.vocabulary_id == req.vocabulary_id && s.user_id == user_id);

        match ReviewViewSet::record_review(req, schedule, user_id) {
            RecordReviewResult::Success(response) => Response::ok().with_json(&response),
            RecordReviewResult::VocabularyNotFound => {
                Response::not_found().with_json(&json!({"error": "Vocabulary not found"}))
            }
            RecordReviewResult::ScheduleNotFound => {
                Response::not_found().with_json(&json!({"error": "Schedule not found"}))
            }
            RecordReviewResult::InvalidQuality => {
                Response::bad_request().with_json(&json!({"error": "Quality must be 0-5"}))
            }
        }
    }
}

/// Handler for GET /stats/
struct ReviewStatsHandler {
    state: ReviewState,
}

#[async_trait]
impl Handler for ReviewStatsHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let user_id = request
            .extensions
            .get::<Uuid>()
            .unwrap_or_else(Uuid::new_v4);
        let vocabularies = self.state.vocabularies.read().unwrap();
        let schedules = self.state.schedules.read().unwrap();

        let response = ReviewViewSet::get_stats(&vocabularies, &schedules, user_id);
        Response::ok().with_json(&response)
    }
}

/// Create review routes
pub fn routes() -> Vec<Route> {
    let state = ReviewState::default();

    vec![
        path(
            "/vocabulary/",
            VocabularyListHandler {
                state: state.clone(),
            },
        ),
        path(
            "/queue/",
            ReviewQueueHandler {
                state: state.clone(),
            },
        ),
        path(
            "/record/",
            RecordReviewHandler {
                state: state.clone(),
            },
        ),
        path("/stats/", ReviewStatsHandler { state }),
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
    async fn test_list_vocabulary_empty() {
        let state = ReviewState::default();
        let handler = VocabularyListHandler { state };

        let request = Request::builder()
            .method(Method::GET)
            .uri("/vocabulary/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_create_vocabulary() {
        let state = ReviewState::default();
        let handler = VocabularyListHandler { state };

        let body = r#"{"page_id":"550e8400-e29b-41d4-a716-446655440000","word":"hello","meaning":"greeting"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/vocabulary/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 201);
    }

    #[tokio::test]
    async fn test_get_queue() {
        let state = ReviewState::default();
        let handler = ReviewQueueHandler { state };

        let request = Request::builder()
            .method(Method::GET)
            .uri("/queue/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_get_stats() {
        let state = ReviewState::default();
        let handler = ReviewStatsHandler { state };

        let request = Request::builder()
            .method(Method::GET)
            .uri("/stats/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }
}
