//! Learning session and progress models

use chrono::{DateTime, Utc};
use reinhardt::db::orm::{FieldSelector, Model, Timestamped};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Type-safe field selector for LearningSession model
#[derive(Clone)]
pub struct LearningSessionFields;

impl FieldSelector for LearningSessionFields {
    fn with_alias(self, _alias: &str) -> Self {
        self
    }
}

/// Type-safe field selector for LearningProgress model
#[derive(Clone)]
pub struct LearningProgressFields;

impl FieldSelector for LearningProgressFields {
    fn with_alias(self, _alias: &str) -> Self {
        self
    }
}

/// Session type
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum SessionType {
    PageByPage,
    TeacherMode,
    Review,
}

/// Session status
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
#[derive(Default)]
pub enum SessionStatus {
    #[default]
    Active,
    Paused,
    Completed,
    Abandoned,
}

/// Session settings for teacher mode
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SessionSettings {
    pub tts_speed: f32,
    pub page_interval: u32,
    pub repeat_count: u32,
    pub include_translation: bool,
    pub include_vocabulary: bool,
    pub include_grammar: bool,
}

impl Default for SessionSettings {
    fn default() -> Self {
        Self {
            tts_speed: 1.0,
            page_interval: 5,
            repeat_count: 1,
            include_translation: true,
            include_vocabulary: true,
            include_grammar: false,
        }
    }
}

/// Learning session
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LearningSession {
    pub id: Uuid,
    pub user_id: Uuid,
    pub book_id: Option<Uuid>,
    pub session_type: SessionType,
    pub start_page: Option<i32>,
    pub end_page: Option<i32>,
    pub current_page: i32,
    pub duration_seconds: i32,
    pub settings: SessionSettings,
    pub status: SessionStatus,
    pub started_at: DateTime<Utc>,
    pub ended_at: Option<DateTime<Utc>>,
}

impl LearningSession {
    /// Create a new learning session
    pub fn new(user_id: Uuid, book_id: Uuid, session_type: SessionType) -> Self {
        Self {
            id: Uuid::new_v4(),
            user_id,
            book_id: Some(book_id),
            session_type,
            start_page: None,
            end_page: None,
            current_page: 1,
            duration_seconds: 0,
            settings: SessionSettings::default(),
            status: SessionStatus::Active,
            started_at: Utc::now(),
            ended_at: None,
        }
    }

    /// Create a review session (no book)
    pub fn new_review(user_id: Uuid) -> Self {
        Self {
            id: Uuid::new_v4(),
            user_id,
            book_id: None,
            session_type: SessionType::Review,
            start_page: None,
            end_page: None,
            current_page: 0,
            duration_seconds: 0,
            settings: SessionSettings::default(),
            status: SessionStatus::Active,
            started_at: Utc::now(),
            ended_at: None,
        }
    }

    /// Set page range for the session
    pub fn with_pages(mut self, start: i32, end: i32) -> Self {
        self.start_page = Some(start);
        self.end_page = Some(end);
        self.current_page = start;
        self
    }

    /// Pause the session
    pub fn pause(&mut self) {
        if self.status == SessionStatus::Active {
            self.status = SessionStatus::Paused;
        }
    }

    /// Resume the session
    pub fn resume(&mut self) {
        if self.status == SessionStatus::Paused {
            self.status = SessionStatus::Active;
        }
    }

    /// Complete the session
    pub fn complete(&mut self) {
        self.status = SessionStatus::Completed;
        self.ended_at = Some(Utc::now());
    }

    /// Abandon the session
    pub fn abandon(&mut self) {
        self.status = SessionStatus::Abandoned;
        self.ended_at = Some(Utc::now());
    }

    /// Advance to next page
    pub fn next_page(&mut self) -> bool {
        if let Some(end) = self.end_page {
            if self.current_page < end {
                self.current_page += 1;
                return true;
            }
        } else {
            self.current_page += 1;
            return true;
        }
        false
    }

    /// Check if session is finished
    pub fn is_finished(&self) -> bool {
        if let Some(end) = self.end_page {
            self.current_page >= end
        } else {
            false
        }
    }
}

/// Learning progress for a page
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LearningProgress {
    pub id: Uuid,
    pub session_id: Uuid,
    pub page_id: Uuid,
    pub user_id: Uuid,
    pub time_spent_seconds: i32,
    pub pronunciation_score: Option<i32>,
    pub comprehension_score: Option<i32>,
    pub feedback_data: Option<serde_json::Value>,
    pub created_at: DateTime<Utc>,
}

impl LearningProgress {
    pub fn new(session_id: Uuid, page_id: Uuid, user_id: Uuid) -> Self {
        Self {
            id: Uuid::new_v4(),
            session_id,
            page_id,
            user_id,
            time_spent_seconds: 0,
            pronunciation_score: None,
            comprehension_score: None,
            feedback_data: None,
            created_at: Utc::now(),
        }
    }

    /// Add time spent on page
    pub fn add_time(&mut self, seconds: i32) {
        self.time_spent_seconds += seconds;
    }

    /// Set pronunciation score (0-100)
    pub fn set_pronunciation_score(&mut self, score: i32) {
        self.pronunciation_score = Some(score.clamp(0, 100));
    }

    /// Calculate average score
    pub fn average_score(&self) -> Option<f32> {
        match (self.pronunciation_score, self.comprehension_score) {
            (Some(p), Some(c)) => Some((p + c) as f32 / 2.0),
            (Some(p), None) => Some(p as f32),
            (None, Some(c)) => Some(c as f32),
            (None, None) => None,
        }
    }
}

impl Model for LearningSession {
    type PrimaryKey = Uuid;
    type Fields = LearningSessionFields;

    fn table_name() -> &'static str {
        "learning_sessions"
    }

    fn app_label() -> &'static str {
        "learning"
    }

    fn new_fields() -> Self::Fields {
        LearningSessionFields
    }

    fn primary_key(&self) -> Option<Self::PrimaryKey> {
        Some(self.id)
    }

    fn set_primary_key(&mut self, value: Self::PrimaryKey) {
        self.id = value;
    }
}

impl Timestamped for LearningSession {
    fn created_at(&self) -> DateTime<Utc> {
        self.started_at
    }

    fn updated_at(&self) -> DateTime<Utc> {
        self.ended_at.unwrap_or(self.started_at)
    }

    fn set_updated_at(&mut self, _time: DateTime<Utc>) {
        // Session timestamps are managed by state transitions
    }
}

impl Model for LearningProgress {
    type PrimaryKey = Uuid;
    type Fields = LearningProgressFields;

    fn table_name() -> &'static str {
        "learning_progress"
    }

    fn app_label() -> &'static str {
        "learning"
    }

    fn new_fields() -> Self::Fields {
        LearningProgressFields
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
    fn test_create_learning_session() {
        let user_id = Uuid::new_v4();
        let book_id = Uuid::new_v4();
        let session = LearningSession::new(user_id, book_id, SessionType::PageByPage);

        assert_eq!(session.user_id, user_id);
        assert_eq!(session.book_id, Some(book_id));
        assert_eq!(session.session_type, SessionType::PageByPage);
        assert_eq!(session.status, SessionStatus::Active);
    }

    #[test]
    fn test_create_review_session() {
        let user_id = Uuid::new_v4();
        let session = LearningSession::new_review(user_id);

        assert_eq!(session.book_id, None);
        assert_eq!(session.session_type, SessionType::Review);
    }

    #[test]
    fn test_session_with_pages() {
        let session =
            LearningSession::new(Uuid::new_v4(), Uuid::new_v4(), SessionType::TeacherMode)
                .with_pages(5, 15);

        assert_eq!(session.start_page, Some(5));
        assert_eq!(session.end_page, Some(15));
        assert_eq!(session.current_page, 5);
    }

    #[test]
    fn test_session_pause_resume() {
        let mut session =
            LearningSession::new(Uuid::new_v4(), Uuid::new_v4(), SessionType::PageByPage);

        session.pause();
        assert_eq!(session.status, SessionStatus::Paused);

        session.resume();
        assert_eq!(session.status, SessionStatus::Active);
    }

    #[test]
    fn test_session_complete() {
        let mut session =
            LearningSession::new(Uuid::new_v4(), Uuid::new_v4(), SessionType::PageByPage);

        session.complete();
        assert_eq!(session.status, SessionStatus::Completed);
        assert!(session.ended_at.is_some());
    }

    #[test]
    fn test_session_next_page() {
        let mut session =
            LearningSession::new(Uuid::new_v4(), Uuid::new_v4(), SessionType::PageByPage)
                .with_pages(1, 3);

        assert!(session.next_page());
        assert_eq!(session.current_page, 2);

        assert!(session.next_page());
        assert_eq!(session.current_page, 3);

        assert!(!session.next_page()); // Can't go past end
        assert_eq!(session.current_page, 3);
    }

    #[test]
    fn test_session_is_finished() {
        let mut session =
            LearningSession::new(Uuid::new_v4(), Uuid::new_v4(), SessionType::PageByPage)
                .with_pages(1, 2);

        assert!(!session.is_finished());
        session.next_page();
        assert!(session.is_finished());
    }

    #[test]
    fn test_create_learning_progress() {
        let session_id = Uuid::new_v4();
        let page_id = Uuid::new_v4();
        let user_id = Uuid::new_v4();

        let progress = LearningProgress::new(session_id, page_id, user_id);

        assert_eq!(progress.session_id, session_id);
        assert_eq!(progress.time_spent_seconds, 0);
    }

    #[test]
    fn test_progress_add_time() {
        let mut progress = LearningProgress::new(Uuid::new_v4(), Uuid::new_v4(), Uuid::new_v4());

        progress.add_time(30);
        progress.add_time(45);

        assert_eq!(progress.time_spent_seconds, 75);
    }

    #[test]
    fn test_progress_score_clamping() {
        let mut progress = LearningProgress::new(Uuid::new_v4(), Uuid::new_v4(), Uuid::new_v4());

        progress.set_pronunciation_score(150); // Should clamp to 100
        assert_eq!(progress.pronunciation_score, Some(100));

        progress.set_pronunciation_score(-10); // Should clamp to 0
        assert_eq!(progress.pronunciation_score, Some(0));
    }

    #[test]
    fn test_progress_average_score() {
        let mut progress = LearningProgress::new(Uuid::new_v4(), Uuid::new_v4(), Uuid::new_v4());

        assert_eq!(progress.average_score(), None);

        progress.pronunciation_score = Some(80);
        assert_eq!(progress.average_score(), Some(80.0));

        progress.comprehension_score = Some(90);
        assert_eq!(progress.average_score(), Some(85.0));
    }

    #[test]
    fn test_session_settings_default() {
        let settings = SessionSettings::default();

        assert_eq!(settings.tts_speed, 1.0);
        assert_eq!(settings.page_interval, 5);
        assert!(settings.include_translation);
    }

    #[test]
    fn test_session_type_serialization() {
        let session_type = SessionType::TeacherMode;
        let json = serde_json::to_string(&session_type).unwrap();
        assert_eq!(json, "\"teacher_mode\"");
    }
}
