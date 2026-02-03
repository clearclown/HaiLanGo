//! Books API routes

use axum::{
    Router,
    extract::{Path, State},
    http::StatusCode,
    response::{IntoResponse, Json},
    routing::get,
};
use serde_json::json;
use std::sync::{Arc, RwLock};
use uuid::Uuid;

use crate::apps::books::{
    dto::CreateBookRequest,
    models::Book,
    views::{BooksViewSet, CreateBookResult, GetBookResult},
};

/// Simulated book store (in production, use database)
#[derive(Clone, Default)]
pub struct BooksState {
    pub books: Arc<RwLock<Vec<Book>>>,
}

/// GET /api/books
async fn list_books(State(state): State<BooksState>) -> impl IntoResponse {
    // Mock user_id (in production, extract from JWT)
    let user_id = Uuid::new_v4();
    let books = state.books.read().unwrap();

    let response = BooksViewSet::list(user_id, &books);
    (StatusCode::OK, Json(json!(response)))
}

/// POST /api/books
async fn create_book(
    State(state): State<BooksState>,
    Json(request): Json<CreateBookRequest>,
) -> impl IntoResponse {
    // Mock user_id (in production, extract from JWT)
    let user_id = Uuid::new_v4();

    match BooksViewSet::create(user_id, request.clone()) {
        CreateBookResult::Success(response) => {
            // Store in memory (in production, save to database)
            let book = Book::new(
                user_id,
                request.title,
                request.source_language,
                request.target_language,
            );
            state.books.write().unwrap().push(book);
            (StatusCode::CREATED, Json(json!(response)))
        }
        CreateBookResult::InvalidInput(msg) => {
            (StatusCode::BAD_REQUEST, Json(json!({"error": msg})))
        }
        CreateBookResult::Unauthorized => (
            StatusCode::UNAUTHORIZED,
            Json(json!({"error": "Unauthorized"})),
        ),
    }
}

/// GET /api/books/:id
async fn get_book(State(state): State<BooksState>, Path(id): Path<Uuid>) -> impl IntoResponse {
    let user_id = Uuid::new_v4();
    let books = state.books.read().unwrap();
    let book = books.iter().find(|b| b.id == id);

    match BooksViewSet::retrieve(user_id, book) {
        GetBookResult::Success(response) => (StatusCode::OK, Json(json!(response))),
        GetBookResult::NotFound => (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "Book not found"})),
        ),
        GetBookResult::Unauthorized => (
            StatusCode::FORBIDDEN,
            Json(json!({"error": "Access denied"})),
        ),
    }
}

/// Create books router
pub fn router() -> Router {
    let state = BooksState::default();

    Router::new()
        .route("/", get(list_books).post(create_book))
        .route("/{id}", get(get_book))
        .with_state(state)
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::{body::Body, http::Request};
    use tower::ServiceExt;

    #[tokio::test]
    async fn test_list_books_empty() {
        let app = router();

        let response = app
            .oneshot(Request::builder().uri("/").body(Body::empty()).unwrap())
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_create_book() {
        let app = router();

        let body = r#"{"title":"Test Book","source_language":"en","target_language":"ja"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::CREATED);
    }

    #[tokio::test]
    async fn test_create_book_invalid() {
        let app = router();

        let body = r#"{"title":"","source_language":"en","target_language":"ja"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }
}
