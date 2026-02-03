//! Review API routes

use axum::{
    Router,
    extract::{Query, State},
    http::StatusCode,
    response::{IntoResponse, Json},
    routing::{get, post},
};
use serde::Deserialize;
use serde_json::json;
use std::sync::{Arc, RwLock};
use uuid::Uuid;

use crate::apps::review::{
    dto::{CreateVocabularyRequest, RecordReviewRequest},
    models::{SrsSchedule, Vocabulary},
    views::{CreateVocabularyResult, RecordReviewResult, ReviewQueueResult, ReviewViewSet},
};

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

/// GET /api/review/vocabulary
async fn list_vocabulary(State(state): State<ReviewState>) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let vocabularies = state.vocabularies.read().unwrap();

    let response = ReviewViewSet::list_vocabularies(&vocabularies, user_id);
    (StatusCode::OK, Json(json!(response)))
}

/// POST /api/review/vocabulary
async fn create_vocabulary(
    State(state): State<ReviewState>,
    Json(request): Json<CreateVocabularyRequest>,
) -> impl IntoResponse {
    let user_id = Uuid::new_v4();

    // Mock page exists check
    let page_exists = true;

    // Check for duplicate
    let vocabularies = state.vocabularies.read().unwrap();
    let word_exists = vocabularies
        .iter()
        .any(|v| v.user_id == user_id && v.word == request.word);
    drop(vocabularies);

    match ReviewViewSet::create_vocabulary(request.clone(), user_id, page_exists, word_exists) {
        CreateVocabularyResult::Success(response) => {
            // Store vocabulary and create SRS schedule
            let vocab = Vocabulary::new(
                request.page_id,
                user_id,
                response.word.clone(),
                response.meaning.clone(),
            );
            let schedule = SrsSchedule::new(user_id, vocab.id);

            state.vocabularies.write().unwrap().push(vocab);
            state.schedules.write().unwrap().push(schedule);

            (StatusCode::CREATED, Json(json!(response)))
        }
        CreateVocabularyResult::PageNotFound => (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "Page not found"})),
        ),
        CreateVocabularyResult::DuplicateWord => (
            StatusCode::CONFLICT,
            Json(json!({"error": "Word already exists"})),
        ),
        CreateVocabularyResult::InvalidInput(msg) => {
            (StatusCode::BAD_REQUEST, Json(json!({"error": msg})))
        }
    }
}

/// GET /api/review/queue
async fn get_review_queue(
    State(state): State<ReviewState>,
    Query(query): Query<QueueQuery>,
) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let vocabularies = state.vocabularies.read().unwrap();
    let schedules = state.schedules.read().unwrap();

    match ReviewViewSet::get_review_queue(&vocabularies, &schedules, user_id, query.limit) {
        ReviewQueueResult::Success(response) => (StatusCode::OK, Json(json!(response))),
        ReviewQueueResult::Empty => (
            StatusCode::OK,
            Json(json!({"items": [], "due_count": 0, "total_vocabulary": 0})),
        ),
    }
}

/// POST /api/review/record
async fn record_review(
    State(state): State<ReviewState>,
    Json(request): Json<RecordReviewRequest>,
) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let mut schedules = state.schedules.write().unwrap();
    let schedule = schedules
        .iter_mut()
        .find(|s| s.vocabulary_id == request.vocabulary_id && s.user_id == user_id);

    match ReviewViewSet::record_review(request, schedule, user_id) {
        RecordReviewResult::Success(response) => (StatusCode::OK, Json(json!(response))),
        RecordReviewResult::VocabularyNotFound => (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "Vocabulary not found"})),
        ),
        RecordReviewResult::ScheduleNotFound => (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "Schedule not found"})),
        ),
        RecordReviewResult::InvalidQuality => (
            StatusCode::BAD_REQUEST,
            Json(json!({"error": "Quality must be 0-5"})),
        ),
    }
}

/// GET /api/review/stats
async fn get_stats(State(state): State<ReviewState>) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let vocabularies = state.vocabularies.read().unwrap();
    let schedules = state.schedules.read().unwrap();

    let response = ReviewViewSet::get_stats(&vocabularies, &schedules, user_id);
    (StatusCode::OK, Json(json!(response)))
}

/// Create review router
pub fn router() -> Router {
    let state = ReviewState::default();

    Router::new()
        .route("/vocabulary", get(list_vocabulary).post(create_vocabulary))
        .route("/queue", get(get_review_queue))
        .route("/record", post(record_review))
        .route("/stats", get(get_stats))
        .with_state(state)
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::{body::Body, http::Request};
    use tower::ServiceExt;

    #[tokio::test]
    async fn test_list_vocabulary_empty() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/vocabulary")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_create_vocabulary() {
        let app = router();

        let body = r#"{"page_id":"550e8400-e29b-41d4-a716-446655440000","word":"hello","meaning":"greeting"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/vocabulary")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::CREATED);
    }

    #[tokio::test]
    async fn test_get_queue() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/queue")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_get_stats() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/stats")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }
}
