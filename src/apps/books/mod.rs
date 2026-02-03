//! Books app - Book management and OCR processing

pub mod dto;
pub mod models;
pub mod views;

pub use dto::*;
pub use models::{Book, BookSettings, BookStatus, Page};
pub use views::{BooksViewSet, CreateBookResult, GetBookResult};
