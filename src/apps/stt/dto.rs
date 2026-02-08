//! Data Transfer Objects for STT API

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use super::models::AttemptStatus;
use crate::services::stt::SttAudioFormat;

/// Request to evaluate pronunciation
#[derive(Debug, Clone, Deserialize)]
pub struct PronunciationRequest {
    pub expected_text: String,
    pub language: String,
    #[serde(default)]
    pub audio_format: SttAudioFormat,
    pub page_id: Option<Uuid>,
}

/// Response for a pronunciation evaluation
#[derive(Debug, Clone, Serialize)]
pub struct PronunciationResponse {
    pub attempt_id: Uuid,
    pub overall_score: u8,
    pub recognized_text: String,
    pub expected_text: String,
    pub feedback: String,
    pub word_scores: Vec<WordScoreResponse>,
    pub duration_ms: u64,
}

/// Per-word score in response
#[derive(Debug, Clone, Serialize)]
pub struct WordScoreResponse {
    pub word: String,
    pub score: u8,
    pub feedback: Option<String>,
}

/// Pronunciation attempt summary (for history listing)
#[derive(Debug, Clone, Serialize)]
pub struct AttemptSummaryResponse {
    pub id: Uuid,
    pub expected_text: String,
    pub recognized_text: Option<String>,
    pub overall_score: Option<u8>,
    pub language: String,
    pub status: AttemptStatus,
    pub created_at: DateTime<Utc>,
}

/// Pronunciation stats response
#[derive(Debug, Clone, Serialize)]
pub struct PronunciationStatsResponse {
    pub total_attempts: u64,
    pub average_score: f32,
    pub best_score: u8,
    pub weak_words: Vec<WeakWordResponse>,
}

/// Weak word in stats response
#[derive(Debug, Clone, Serialize)]
pub struct WeakWordResponse {
    pub word: String,
    pub language: String,
    pub average_score: f32,
    pub attempt_count: u32,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_pronunciation_request_deserialization() {
        let json = r#"{
            "expected_text": "hello world",
            "language": "en"
        }"#;

        let request: PronunciationRequest = serde_json::from_str(json).unwrap();
        assert_eq!(request.expected_text, "hello world");
        assert_eq!(request.language, "en");
        assert_eq!(request.audio_format, SttAudioFormat::Wav);
        assert!(request.page_id.is_none());
    }

    #[test]
    fn test_pronunciation_request_with_all_fields() {
        let page_id = Uuid::new_v4();
        let json = format!(
            r#"{{
                "expected_text": "bonjour",
                "language": "fr",
                "audio_format": "mp3",
                "page_id": "{}"
            }}"#,
            page_id
        );

        let request: PronunciationRequest = serde_json::from_str(&json).unwrap();
        assert_eq!(request.expected_text, "bonjour");
        assert_eq!(request.language, "fr");
        assert_eq!(request.audio_format, SttAudioFormat::Mp3);
        assert_eq!(request.page_id, Some(page_id));
    }

    #[test]
    fn test_pronunciation_response_serialization() {
        let response = PronunciationResponse {
            attempt_id: Uuid::new_v4(),
            overall_score: 85,
            recognized_text: "hello world".to_string(),
            expected_text: "hello world".to_string(),
            feedback: "Good pronunciation!".to_string(),
            word_scores: vec![
                WordScoreResponse {
                    word: "hello".to_string(),
                    score: 90,
                    feedback: None,
                },
                WordScoreResponse {
                    word: "world".to_string(),
                    score: 80,
                    feedback: Some("Slightly unclear".to_string()),
                },
            ],
            duration_ms: 1200,
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("\"overall_score\":85"));
        assert!(json.contains("\"hello\""));
        assert!(json.contains("\"world\""));
    }

    #[test]
    fn test_attempt_summary_serialization() {
        let summary = AttemptSummaryResponse {
            id: Uuid::new_v4(),
            expected_text: "test".to_string(),
            recognized_text: Some("test".to_string()),
            overall_score: Some(95),
            language: "en".to_string(),
            status: AttemptStatus::Completed,
            created_at: Utc::now(),
        };

        let json = serde_json::to_string(&summary).unwrap();
        assert!(json.contains("\"completed\""));
        assert!(json.contains("\"overall_score\":95"));
    }

    #[test]
    fn test_stats_response_serialization() {
        let stats = PronunciationStatsResponse {
            total_attempts: 100,
            average_score: 82.5,
            best_score: 98,
            weak_words: vec![WeakWordResponse {
                word: "difficult".to_string(),
                language: "en".to_string(),
                average_score: 45.0,
                attempt_count: 8,
            }],
        };

        let json = serde_json::to_string(&stats).unwrap();
        assert!(json.contains("\"total_attempts\":100"));
        assert!(json.contains("\"difficult\""));
    }
}
