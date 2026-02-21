//! Book and Page models

use chrono::{DateTime, Utc};
use reinhardt::db::orm::{FieldSelector, Model, Timestamped};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Type-safe field selector for Book model
#[derive(Clone)]
pub struct BookFields;

impl FieldSelector for BookFields {
    fn with_alias(self, _alias: &str) -> Self {
        self
    }
}

/// Type-safe field selector for Page model
#[derive(Clone)]
pub struct PageFields;

impl FieldSelector for PageFields {
    fn with_alias(self, _alias: &str) -> Self {
        self
    }
}

/// Book processing status
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(rename_all = "snake_case")]
pub enum BookStatus {
    #[default]
    Pending,
    Processing,
    Ready,
    Error,
}

/// Book settings for TTS and learning
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct BookSettings {
    pub tts_language: Option<String>,
    pub tts_speed: Option<f32>,
    pub auto_play: Option<bool>,
}

/// Book entity representing an uploaded textbook
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Book {
    pub id: Uuid,
    pub user_id: Uuid,
    pub title: String,
    pub source_language: String,
    pub target_language: String,
    pub reference_language: Option<String>,
    pub total_pages: i32,
    pub status: BookStatus,
    pub encryption_key_hash: Option<String>,
    pub settings: BookSettings,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

impl Book {
    /// Create a new book in pending status
    pub fn new(
        user_id: Uuid,
        title: String,
        source_language: String,
        target_language: String,
    ) -> Self {
        let now = Utc::now();
        Self {
            id: Uuid::new_v4(),
            user_id,
            title,
            source_language,
            target_language,
            reference_language: None,
            total_pages: 0,
            status: BookStatus::Pending,
            encryption_key_hash: None,
            settings: BookSettings::default(),
            created_at: now,
            updated_at: now,
        }
    }

    /// Update book status
    pub fn set_status(&mut self, status: BookStatus) {
        self.status = status;
        self.updated_at = Utc::now();
    }

    /// Update total page count
    pub fn set_total_pages(&mut self, count: i32) {
        self.total_pages = count;
        self.updated_at = Utc::now();
    }
}

/// Page entity representing a single page in a book
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Page {
    pub id: Uuid,
    pub book_id: Uuid,
    pub page_number: i32,
    pub original_content: Option<String>,
    pub processed_content: Option<String>,
    pub layout_data: Option<serde_json::Value>,
    pub audio_url: Option<String>,
    pub is_processed: bool,
    pub created_at: DateTime<Utc>,
}

impl Page {
    /// Create a new unprocessed page
    pub fn new(book_id: Uuid, page_number: i32) -> Self {
        Self {
            id: Uuid::new_v4(),
            book_id,
            page_number,
            original_content: None,
            processed_content: None,
            layout_data: None,
            audio_url: None,
            is_processed: false,
            created_at: Utc::now(),
        }
    }

    /// Set OCR content
    pub fn set_content(&mut self, original: String, processed: Option<String>) {
        self.original_content = Some(original);
        self.processed_content = processed;
        self.is_processed = true;
    }
}

impl Model for Book {
    type PrimaryKey = Uuid;
    type Fields = BookFields;

    fn table_name() -> &'static str {
        "books"
    }

    fn app_label() -> &'static str {
        "books"
    }

    fn new_fields() -> Self::Fields {
        BookFields
    }

    fn primary_key(&self) -> Option<Self::PrimaryKey> {
        Some(self.id)
    }

    fn set_primary_key(&mut self, value: Self::PrimaryKey) {
        self.id = value;
    }
}

impl Timestamped for Book {
    fn created_at(&self) -> DateTime<Utc> {
        self.created_at
    }

    fn updated_at(&self) -> DateTime<Utc> {
        self.updated_at
    }

    fn set_updated_at(&mut self, time: DateTime<Utc>) {
        self.updated_at = time;
    }
}

impl Model for Page {
    type PrimaryKey = Uuid;
    type Fields = PageFields;

    fn table_name() -> &'static str {
        "pages"
    }

    fn app_label() -> &'static str {
        "books"
    }

    fn new_fields() -> Self::Fields {
        PageFields
    }

    fn primary_key(&self) -> Option<Self::PrimaryKey> {
        Some(self.id)
    }

    fn set_primary_key(&mut self, value: Self::PrimaryKey) {
        self.id = value;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_create_book() {
        let user_id = Uuid::new_v4();
        let book = Book::new(
            user_id,
            "Language Book".to_string(),
            "en".to_string(),
            "ja".to_string(),
        );

        assert_eq!(book.user_id, user_id);
        assert_eq!(book.title, "Language Book");
        assert_eq!(book.status, BookStatus::Pending);
        assert_eq!(book.total_pages, 0);
    }

    #[test]
    fn test_book_status_update() {
        let mut book = Book::new(
            Uuid::new_v4(),
            "Test".to_string(),
            "en".to_string(),
            "es".to_string(),
        );

        book.set_status(BookStatus::Processing);
        assert_eq!(book.status, BookStatus::Processing);

        book.set_status(BookStatus::Ready);
        assert_eq!(book.status, BookStatus::Ready);
    }

    #[test]
    fn test_create_page() {
        let book_id = Uuid::new_v4();
        let page = Page::new(book_id, 1);

        assert_eq!(page.book_id, book_id);
        assert_eq!(page.page_number, 1);
        assert!(!page.is_processed);
    }

    #[test]
    fn test_page_content_update() {
        let mut page = Page::new(Uuid::new_v4(), 1);

        page.set_content(
            "Hello world".to_string(),
            Some("Processed: Hello world".to_string()),
        );

        assert!(page.is_processed);
        assert_eq!(page.original_content, Some("Hello world".to_string()));
    }

    #[test]
    fn test_book_status_serialization() {
        let status = BookStatus::Processing;
        let json = serde_json::to_string(&status).unwrap();
        assert_eq!(json, "\"processing\"");
    }

    #[test]
    fn test_book_page_count() {
        let mut book = Book::new(
            Uuid::new_v4(),
            "Test".to_string(),
            "en".to_string(),
            "ja".to_string(),
        );

        assert_eq!(book.total_pages, 0);
        book.set_total_pages(100);
        assert_eq!(book.total_pages, 100);
    }

    #[test]
    fn test_page_optional_fields() {
        let page = Page::new(Uuid::new_v4(), 1);

        assert!(page.original_content.is_none());
        assert!(page.processed_content.is_none());
        assert!(page.audio_url.is_none());
    }

    #[test]
    fn test_book_timestamps() {
        let book = Book::new(
            Uuid::new_v4(),
            "Test".to_string(),
            "en".to_string(),
            "ja".to_string(),
        );

        assert!(book.created_at <= Utc::now());
        assert_eq!(book.created_at, book.updated_at);
    }
}
