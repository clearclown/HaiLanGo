//! Review app - Spaced Repetition System (SRS)
//!
//! Implements learning algorithm for intelligent review scheduling.
//! Uses SM-2 algorithm to optimize retention and spacing.
//!
//! Models:
//! - Vocabulary: Extracted words from pages
//! - SrsSchedule: Scheduled review dates using SM-2 algorithm
//!
//! Services:
//! - ReviewScheduler: Calculate and manage review queues
//! - DifficultyAnalyzer: Analyze item difficulty from user responses
//!
//! Views:
//! - ReviewViewSet: Get next review items and record responses

pub mod dto;
pub mod models;
pub mod views;

pub use dto::{
    BulkReviewRequest, BulkReviewResultResponse, CreateVocabularyRequest, RecordReviewRequest,
    ReviewItemResponse, ReviewQueueResponse, ReviewResultResponse, ReviewStatsResponse,
    SrsScheduleResponse, VocabularyResponse,
};
pub use models::{SrsSchedule, Vocabulary};
pub use views::{CreateVocabularyResult, RecordReviewResult, ReviewQueueResult, ReviewViewSet};
