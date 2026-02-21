//! Books API routes

use async_trait::async_trait;
use serde_json::json;
use std::sync::{Arc, RwLock};
use uuid::Uuid;

use crate::{Handler, Method, Request, Response, Result, Route, path};
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

/// Handler for GET/POST /api/books/
struct BooksListHandler {
    state: BooksState,
}

#[async_trait]
impl Handler for BooksListHandler {
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

impl BooksListHandler {
    fn list(&self) -> Result<Response> {
        let user_id = Uuid::new_v4();
        let books = self.state.books.read().unwrap();
        let response = BooksViewSet::list(user_id, &books);
        Response::ok().with_json(&response)
    }

    fn create(&self, request: Request) -> Result<Response> {
        let req: CreateBookRequest = request.json()?;
        let user_id = Uuid::new_v4();

        match BooksViewSet::create(user_id, req.clone()) {
            CreateBookResult::Success(response) => {
                let book = Book::new(
                    user_id,
                    req.title,
                    req.source_language,
                    req.target_language,
                );
                self.state.books.write().unwrap().push(book);
                Response::created().with_json(&response)
            }
            CreateBookResult::InvalidInput(msg) => {
                Response::bad_request().with_json(&json!({"error": msg}))
            }
            CreateBookResult::Unauthorized => {
                Response::unauthorized().with_json(&json!({"error": "Unauthorized"}))
            }
        }
    }
}

/// Handler for GET /api/books/{id}/
struct BooksDetailHandler {
    state: BooksState,
}

#[async_trait]
impl Handler for BooksDetailHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        match request.method {
            Method::GET => self.retrieve(request),
            _ => Err(crate::Error::MethodNotAllowed(
                "Only GET is allowed".into(),
            )),
        }
    }
}

impl BooksDetailHandler {
    fn retrieve(&self, request: Request) -> Result<Response> {
        let id: Uuid = request
            .path_params
            .get("id")
            .ok_or_else(|| crate::Error::Validation("Missing id parameter".into()))?
            .parse()
            .map_err(|_| crate::Error::Validation("Invalid UUID".into()))?;

        let user_id = Uuid::new_v4();
        let books = self.state.books.read().unwrap();
        let book = books.iter().find(|b| b.id == id);

        match BooksViewSet::retrieve(user_id, book) {
            GetBookResult::Success(response) => Response::ok().with_json(&response),
            GetBookResult::NotFound => {
                Response::not_found().with_json(&json!({"error": "Book not found"}))
            }
            GetBookResult::Unauthorized => {
                Response::forbidden().with_json(&json!({"error": "Access denied"}))
            }
        }
    }
}

/// Create books routes
pub fn routes() -> Vec<Route> {
    let state = BooksState::default();

    vec![
        path("/", BooksListHandler { state: state.clone() }),
        path("/{id}/", BooksDetailHandler { state }),
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

    fn build_handler() -> (BooksListHandler, BooksDetailHandler) {
        let state = BooksState::default();
        (
            BooksListHandler { state: state.clone() },
            BooksDetailHandler { state },
        )
    }

    #[tokio::test]
    async fn test_list_books_empty() {
        let (handler, _) = build_handler();

        let request = Request::builder()
            .method(Method::GET)
            .uri("/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_create_book() {
        let (handler, _) = build_handler();

        let body = r#"{"title":"Test Book","source_language":"en","target_language":"ja"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 201);
    }

    #[tokio::test]
    async fn test_create_book_invalid() {
        let (handler, _) = build_handler();

        let body = r#"{"title":"","source_language":"en","target_language":"ja"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 400);
    }
}
