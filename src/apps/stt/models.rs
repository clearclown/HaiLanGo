//! STT models - Pronunciation attempt and feedback data

use chrono::{DateTime, Utc};
use reinhardt::db::orm::{FieldSelector, Model};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Type-safe field selector for PronunciationAttempt model
#[derive(Clone)]
pub struct PronunciationAttemptFields;

impl FieldSelector for PronunciationAttemptFields {
    fn with_alias(self, _alias: &str) -> Self {
        self
    }
}

/// Type-safe field selector for WordFeedback model
#[derive(Clone)]
pub struct WordFeedbackFields;

impl FieldSelector for WordFeedbackFields {
    fn with_alias(self, _alias: &str) -> Self {
        self
    }
}

/// Status of a pronunciation attempt
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(rename_all = "snake_case")]
pub enum AttemptStatus {
    #[default]
    Pending,
    Processing,
    Completed,
    Failed,
}

/// A user's pronunciation attempt
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PronunciationAttempt {
    pub id: Uuid,
    pub user_id: Uuid,
    pub page_id: Option<Uuid>,
    pub expected_text: String,
    pub recognized_text: Option<String>,
    pub language: String,
    pub overall_score: Option<u8>,
    pub status: AttemptStatus,
    pub audio_duration_ms: Option<u64>,
    pub created_at: DateTime<Utc>,
}

impl PronunciationAttempt {
    pub fn new(user_id: Uuid, expected_text: String, language: String) -> Self {
        Self {
            id: Uuid::new_v4(),
            user_id,
            page_id: None,
            expected_text,
            recognized_text: None,
            language,
            overall_score: None,
            status: AttemptStatus::default(),
            audio_duration_ms: None,
            created_at: Utc::now(),
        }
    }

    pub fn with_page(mut self, page_id: Uuid) -> Self {
        self.page_id = Some(page_id);
        self
    }
}

impl Model for PronunciationAttempt {
    type PrimaryKey = Uuid;
    type Fields = PronunciationAttemptFields;

    fn table_name() -> &'static str {
        "pronunciation_attempts"
    }

    fn app_label() -> &'static str {
        "stt"
    }

    fn new_fields() -> Self::Fields {
        PronunciationAttemptFields
    }

    fn primary_key(&self) -> Option<Self::PrimaryKey> {
        Some(self.id)
    }

    fn set_primary_key(&mut self, value: Self::PrimaryKey) {
        self.id = value;
    }
}

/// Detailed per-word feedback stored alongside an attempt
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WordFeedback {
    pub id: Uuid,
    pub attempt_id: Uuid,
    pub word: String,
    pub score: u8,
    pub feedback: Option<String>,
    pub start_ms: Option<u64>,
    pub end_ms: Option<u64>,
}

impl Model for WordFeedback {
    type PrimaryKey = Uuid;
    type Fields = WordFeedbackFields;

    fn table_name() -> &'static str {
        "word_feedback"
    }

    fn app_label() -> &'static str {
        "stt"
    }

    fn new_fields() -> Self::Fields {
        WordFeedbackFields
    }

    fn primary_key(&self) -> Option<Self::PrimaryKey> {
        Some(self.id)
    }

    fn set_primary_key(&mut self, value: Self::PrimaryKey) {
        self.id = value;
    }
}

/// Pronunciation statistics for a user
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PronunciationStats {
    pub user_id: Uuid,
    pub total_attempts: u64,
    pub average_score: f32,
    pub best_score: u8,
    pub weak_words: Vec<WeakWord>,
}

/// A word the user consistently struggles with
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WeakWord {
    pub word: String,
    pub language: String,
    pub average_score: f32,
    pub attempt_count: u32,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_pronunciation_attempt_creation() {
        let user_id = Uuid::new_v4();
        let attempt =
            PronunciationAttempt::new(user_id, "hello world".to_string(), "en".to_string());

        assert_eq!(attempt.user_id, user_id);
        assert_eq!(attempt.expected_text, "hello world");
        assert_eq!(attempt.language, "en");
        assert_eq!(attempt.status, AttemptStatus::Pending);
        assert!(attempt.recognized_text.is_none());
        assert!(attempt.overall_score.is_none());
        assert!(attempt.page_id.is_none());
    }

    #[test]
    fn test_pronunciation_attempt_with_page() {
        let user_id = Uuid::new_v4();
        let page_id = Uuid::new_v4();
        let attempt = PronunciationAttempt::new(user_id, "test".to_string(), "en".to_string())
            .with_page(page_id);

        assert_eq!(attempt.page_id, Some(page_id));
    }

    #[test]
    fn test_attempt_status_default() {
        let status = AttemptStatus::default();
        assert_eq!(status, AttemptStatus::Pending);
    }

    #[test]
    fn test_attempt_status_serialization() {
        let status = AttemptStatus::Completed;
        let json = serde_json::to_string(&status).unwrap();
        assert_eq!(json, "\"completed\"");

        let deserialized: AttemptStatus = serde_json::from_str(&json).unwrap();
        assert_eq!(deserialized, AttemptStatus::Completed);
    }

    #[test]
    fn test_word_feedback_creation() {
        let attempt_id = Uuid::new_v4();
        let feedback = WordFeedback {
            id: Uuid::new_v4(),
            attempt_id,
            word: "hello".to_string(),
            score: 85,
            feedback: Some("Good pronunciation".to_string()),
            start_ms: Some(0),
            end_ms: Some(300),
        };

        assert_eq!(feedback.attempt_id, attempt_id);
        assert_eq!(feedback.score, 85);
    }

    #[test]
    fn test_pronunciation_stats() {
        let stats = PronunciationStats {
            user_id: Uuid::new_v4(),
            total_attempts: 50,
            average_score: 78.5,
            best_score: 96,
            weak_words: vec![WeakWord {
                word: "pronunciation".to_string(),
                language: "en".to_string(),
                average_score: 45.0,
                attempt_count: 5,
            }],
        };

        assert_eq!(stats.total_attempts, 50);
        assert_eq!(stats.weak_words.len(), 1);
        assert_eq!(stats.weak_words[0].word, "pronunciation");
    }
}
