//! Data Transfer Objects for books API

use super::models::BookStatus;
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Request to create a new book
#[derive(Debug, Clone, Deserialize)]
pub struct CreateBookRequest {
    pub title: String,
    pub source_language: String,
    pub target_language: String,
    pub reference_language: Option<String>,
}

/// Response for book creation (upload accepted)
#[derive(Debug, Serialize)]
pub struct UploadAcceptedResponse {
    pub id: Uuid,
    pub title: String,
    pub status: BookStatus,
    pub job_id: Uuid,
}

/// Book details response
#[derive(Debug, Serialize)]
pub struct BookResponse {
    pub id: Uuid,
    pub title: String,
    pub source_language: String,
    pub target_language: String,
    pub total_pages: i32,
    pub status: BookStatus,
    pub progress: Option<BookProgress>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

/// Book processing progress
#[derive(Debug, Serialize)]
pub struct BookProgress {
    pub processed_pages: i32,
    pub total_pages: i32,
    pub percentage: f32,
}

impl BookProgress {
    pub fn new(processed: i32, total: i32) -> Self {
        let percentage = if total > 0 {
            (processed as f32 / total as f32) * 100.0
        } else {
            0.0
        };
        Self {
            processed_pages: processed,
            total_pages: total,
            percentage,
        }
    }
}

/// Page content response
#[derive(Debug, Serialize)]
pub struct PageResponse {
    pub id: Uuid,
    pub page_number: i32,
    pub original_content: Option<String>,
    pub processed_content: Option<String>,
    pub audio_url: Option<String>,
    pub is_processed: bool,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_create_book_request_deserialization() {
        let json = r#"{
            "title": "My Book",
            "source_language": "en",
            "target_language": "ja"
        }"#;

        let request: CreateBookRequest = serde_json::from_str(json).unwrap();
        assert_eq!(request.title, "My Book");
        assert!(request.reference_language.is_none());
    }

    #[test]
    fn test_create_book_request_with_reference() {
        let json = r#"{
            "title": "My Book",
            "source_language": "en",
            "target_language": "ja",
            "reference_language": "zh"
        }"#;

        let request: CreateBookRequest = serde_json::from_str(json).unwrap();
        assert_eq!(request.reference_language, Some("zh".to_string()));
    }

    #[test]
    fn test_book_progress_calculation() {
        let progress = BookProgress::new(5, 10);
        assert_eq!(progress.percentage, 50.0);

        let empty = BookProgress::new(0, 0);
        assert_eq!(empty.percentage, 0.0);
    }

    #[test]
    fn test_book_progress_full() {
        let progress = BookProgress::new(10, 10);
        assert_eq!(progress.percentage, 100.0);
    }

    #[test]
    fn test_upload_response_serialization() {
        let response = UploadAcceptedResponse {
            id: Uuid::new_v4(),
            title: "Test Book".to_string(),
            status: BookStatus::Pending,
            job_id: Uuid::new_v4(),
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("Test Book"));
        assert!(json.contains("pending"));
    }

    #[test]
    fn test_page_response_serialization() {
        let response = PageResponse {
            id: Uuid::new_v4(),
            page_number: 1,
            original_content: Some("text".to_string()),
            processed_content: None,
            audio_url: None,
            is_processed: false,
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("\"page_number\":1"));
    }

    #[test]
    fn test_book_response_serialization() {
        let response = BookResponse {
            id: Uuid::new_v4(),
            title: "Test".to_string(),
            source_language: "en".to_string(),
            target_language: "ja".to_string(),
            total_pages: 50,
            status: BookStatus::Ready,
            progress: Some(BookProgress::new(25, 50)),
            created_at: Utc::now(),
            updated_at: Utc::now(),
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("\"total_pages\":50"));
        assert!(json.contains("\"ready\""));
    }
}
