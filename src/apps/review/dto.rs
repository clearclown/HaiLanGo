//! Data Transfer Objects for SRS review system

use chrono::{DateTime, NaiveDate, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Request to add a vocabulary word
#[derive(Debug, Clone, Deserialize)]
pub struct CreateVocabularyRequest {
    pub page_id: Uuid,
    pub word: String,
    pub meaning: String,
    pub reading: Option<String>,
    pub part_of_speech: Option<String>,
    pub example_sentence: Option<String>,
}

/// Request to record a review result
#[derive(Debug, Deserialize)]
pub struct RecordReviewRequest {
    pub vocabulary_id: Uuid,
    /// Quality rating: 0-5 (0-2 = fail, 3-5 = pass)
    pub quality: u8,
}

/// Request to bulk record reviews
#[derive(Debug, Deserialize)]
pub struct BulkReviewRequest {
    pub reviews: Vec<RecordReviewRequest>,
}

/// Vocabulary response
#[derive(Debug, Serialize)]
pub struct VocabularyResponse {
    pub id: Uuid,
    pub page_id: Uuid,
    pub word: String,
    pub reading: Option<String>,
    pub meaning: String,
    pub part_of_speech: Option<String>,
    pub example_sentence: Option<String>,
    pub frequency: i32,
    pub created_at: DateTime<Utc>,
}

/// SRS schedule response
#[derive(Debug, Serialize)]
pub struct SrsScheduleResponse {
    pub id: Uuid,
    pub vocabulary_id: Uuid,
    pub next_review_date: NaiveDate,
    pub interval_days: i32,
    pub easiness_factor: f32,
    pub repetitions: i32,
    pub retention_rate: f32,
    pub is_due: bool,
}

/// Review item for study session
#[derive(Debug, Serialize)]
pub struct ReviewItemResponse {
    pub vocabulary: VocabularyResponse,
    pub schedule: SrsScheduleResponse,
}

/// Review queue response
#[derive(Debug, Serialize)]
pub struct ReviewQueueResponse {
    pub items: Vec<ReviewItemResponse>,
    pub due_count: usize,
    pub total_vocabulary: usize,
}

/// Review result after recording
#[derive(Debug, Serialize)]
pub struct ReviewResultResponse {
    pub vocabulary_id: Uuid,
    pub was_correct: bool,
    pub new_interval_days: i32,
    pub next_review_date: NaiveDate,
}

/// Bulk review result
#[derive(Debug, Serialize)]
pub struct BulkReviewResultResponse {
    pub results: Vec<ReviewResultResponse>,
    pub correct_count: usize,
    pub incorrect_count: usize,
}

/// Review statistics response
#[derive(Debug, Serialize)]
pub struct ReviewStatsResponse {
    pub total_vocabulary: usize,
    pub due_today: usize,
    pub overdue: usize,
    pub learned_count: usize,
    pub average_retention_rate: f32,
    pub streak_days: i32,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_create_vocabulary_request_deserialization() {
        let json = r#"{
            "page_id": "550e8400-e29b-41d4-a716-446655440000",
            "word": "hello",
            "meaning": "a greeting"
        }"#;

        let request: CreateVocabularyRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.word, "hello");
        assert_eq!(request.meaning, "a greeting");
        assert!(request.reading.is_none());
    }

    #[test]
    fn test_create_vocabulary_request_full() {
        let json = r#"{
            "page_id": "550e8400-e29b-41d4-a716-446655440000",
            "word": "読む",
            "meaning": "to read",
            "reading": "よむ",
            "part_of_speech": "verb",
            "example_sentence": "本を読む"
        }"#;

        let request: CreateVocabularyRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.word, "読む");
        assert_eq!(request.reading, Some("よむ".to_string()));
        assert_eq!(request.part_of_speech, Some("verb".to_string()));
    }

    #[test]
    fn test_record_review_request_deserialization() {
        let json = r#"{
            "vocabulary_id": "550e8400-e29b-41d4-a716-446655440000",
            "quality": 5
        }"#;

        let request: RecordReviewRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.quality, 5);
    }

    #[test]
    fn test_bulk_review_request_deserialization() {
        let json = r#"{
            "reviews": [
                {"vocabulary_id": "550e8400-e29b-41d4-a716-446655440000", "quality": 5},
                {"vocabulary_id": "550e8400-e29b-41d4-a716-446655440001", "quality": 3}
            ]
        }"#;

        let request: BulkReviewRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.reviews.len(), 2);
        assert_eq!(request.reviews[0].quality, 5);
        assert_eq!(request.reviews[1].quality, 3);
    }

    #[test]
    fn test_vocabulary_response_serialization() {
        let now = Utc::now();
        let response = VocabularyResponse {
            id: Uuid::new_v4(),
            page_id: Uuid::new_v4(),
            word: "hello".to_string(),
            reading: None,
            meaning: "greeting".to_string(),
            part_of_speech: Some("noun".to_string()),
            example_sentence: None,
            frequency: 5,
            created_at: now,
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("hello"));
        assert!(json.contains("greeting"));
        assert!(json.contains("\"frequency\":5"));
    }

    #[test]
    fn test_srs_schedule_response_serialization() {
        let response = SrsScheduleResponse {
            id: Uuid::new_v4(),
            vocabulary_id: Uuid::new_v4(),
            next_review_date: Utc::now().date_naive(),
            interval_days: 6,
            easiness_factor: 2.5,
            repetitions: 3,
            retention_rate: 0.85,
            is_due: true,
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("\"interval_days\":6"));
        assert!(json.contains("\"is_due\":true"));
    }

    #[test]
    fn test_review_queue_response_serialization() {
        let response = ReviewQueueResponse {
            items: vec![],
            due_count: 10,
            total_vocabulary: 50,
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("\"due_count\":10"));
        assert!(json.contains("\"total_vocabulary\":50"));
    }

    #[test]
    fn test_review_result_response_serialization() {
        let response = ReviewResultResponse {
            vocabulary_id: Uuid::new_v4(),
            was_correct: true,
            new_interval_days: 6,
            next_review_date: Utc::now().date_naive(),
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("\"was_correct\":true"));
        assert!(json.contains("\"new_interval_days\":6"));
    }

    #[test]
    fn test_review_stats_response_serialization() {
        let response = ReviewStatsResponse {
            total_vocabulary: 100,
            due_today: 15,
            overdue: 3,
            learned_count: 80,
            average_retention_rate: 0.87,
            streak_days: 7,
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("\"total_vocabulary\":100"));
        assert!(json.contains("\"streak_days\":7"));
    }

    #[test]
    fn test_bulk_review_result_response() {
        let response = BulkReviewResultResponse {
            results: vec![
                ReviewResultResponse {
                    vocabulary_id: Uuid::new_v4(),
                    was_correct: true,
                    new_interval_days: 6,
                    next_review_date: Utc::now().date_naive(),
                },
                ReviewResultResponse {
                    vocabulary_id: Uuid::new_v4(),
                    was_correct: false,
                    new_interval_days: 1,
                    next_review_date: Utc::now().date_naive(),
                },
            ],
            correct_count: 1,
            incorrect_count: 1,
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("\"correct_count\":1"));
        assert!(json.contains("\"incorrect_count\":1"));
    }
}
