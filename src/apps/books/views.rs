//! Books ViewSet - Book management endpoints

use uuid::Uuid;

use super::dto::{BookProgress, BookResponse, CreateBookRequest, UploadAcceptedResponse};
use super::models::Book;

/// Book creation result
#[derive(Debug)]
pub enum CreateBookResult {
    Success(UploadAcceptedResponse),
    InvalidInput(String),
    Unauthorized,
}

/// Book retrieval result
#[derive(Debug)]
pub enum GetBookResult {
    Success(BookResponse),
    NotFound,
    Unauthorized,
}

/// Books ViewSet
pub struct BooksViewSet;

impl BooksViewSet {
    /// List books for a user
    pub fn list(user_id: Uuid, books: &[Book]) -> Vec<BookResponse> {
        books
            .iter()
            .filter(|b| b.user_id == user_id)
            .map(|b| Self::book_to_response(b, None))
            .collect()
    }

    /// Create a new book
    pub fn create(user_id: Uuid, request: CreateBookRequest) -> CreateBookResult {
        // Validate input
        if request.title.is_empty() {
            return CreateBookResult::InvalidInput("Title is required".to_string());
        }

        if request.source_language.is_empty() || request.target_language.is_empty() {
            return CreateBookResult::InvalidInput("Languages are required".to_string());
        }

        // Create book
        let book = Book::new(
            user_id,
            request.title.clone(),
            request.source_language.clone(),
            request.target_language.clone(),
        );

        // Generate job ID for OCR processing
        let job_id = Uuid::new_v4();

        CreateBookResult::Success(UploadAcceptedResponse {
            id: book.id,
            title: book.title,
            status: book.status,
            job_id,
        })
    }

    /// Get a single book
    pub fn retrieve(user_id: Uuid, book: Option<&Book>) -> GetBookResult {
        match book {
            Some(b) if b.user_id == user_id => {
                GetBookResult::Success(Self::book_to_response(b, None))
            }
            Some(_) => GetBookResult::Unauthorized,
            None => GetBookResult::NotFound,
        }
    }

    /// Convert Book to BookResponse
    fn book_to_response(book: &Book, progress: Option<BookProgress>) -> BookResponse {
        BookResponse {
            id: book.id,
            title: book.title.clone(),
            source_language: book.source_language.clone(),
            target_language: book.target_language.clone(),
            total_pages: book.total_pages,
            status: book.status,
            progress,
            created_at: book.created_at,
            updated_at: book.updated_at,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::super::models::BookStatus;
    use super::*;

    #[test]
    fn test_create_book_success() {
        let user_id = Uuid::new_v4();
        let request = CreateBookRequest {
            title: "My Language Book".to_string(),
            source_language: "en".to_string(),
            target_language: "ja".to_string(),
            reference_language: None,
        };

        let result = BooksViewSet::create(user_id, request);

        match result {
            CreateBookResult::Success(response) => {
                assert_eq!(response.title, "My Language Book");
                assert_eq!(response.status, BookStatus::Pending);
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_create_book_empty_title() {
        let user_id = Uuid::new_v4();
        let request = CreateBookRequest {
            title: "".to_string(),
            source_language: "en".to_string(),
            target_language: "ja".to_string(),
            reference_language: None,
        };

        let result = BooksViewSet::create(user_id, request);
        assert!(matches!(result, CreateBookResult::InvalidInput(_)));
    }

    #[test]
    fn test_list_books() {
        let user_id = Uuid::new_v4();
        let other_user_id = Uuid::new_v4();

        let books = vec![
            Book::new(
                user_id,
                "Book 1".to_string(),
                "en".to_string(),
                "ja".to_string(),
            ),
            Book::new(
                user_id,
                "Book 2".to_string(),
                "en".to_string(),
                "es".to_string(),
            ),
            Book::new(
                other_user_id,
                "Other Book".to_string(),
                "en".to_string(),
                "fr".to_string(),
            ),
        ];

        let result = BooksViewSet::list(user_id, &books);

        assert_eq!(result.len(), 2);
        assert!(result.iter().all(|b| b.title != "Other Book"));
    }

    #[test]
    fn test_retrieve_book_success() {
        let user_id = Uuid::new_v4();
        let book = Book::new(
            user_id,
            "Test Book".to_string(),
            "en".to_string(),
            "ja".to_string(),
        );

        let result = BooksViewSet::retrieve(user_id, Some(&book));
        assert!(matches!(result, GetBookResult::Success(_)));
    }

    #[test]
    fn test_retrieve_book_unauthorized() {
        let owner_id = Uuid::new_v4();
        let other_user_id = Uuid::new_v4();
        let book = Book::new(
            owner_id,
            "Test Book".to_string(),
            "en".to_string(),
            "ja".to_string(),
        );

        let result = BooksViewSet::retrieve(other_user_id, Some(&book));
        assert!(matches!(result, GetBookResult::Unauthorized));
    }

    #[test]
    fn test_retrieve_book_not_found() {
        let user_id = Uuid::new_v4();
        let result = BooksViewSet::retrieve(user_id, None);
        assert!(matches!(result, GetBookResult::NotFound));
    }
}
