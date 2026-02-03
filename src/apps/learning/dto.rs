//! Data Transfer Objects for learning sessions

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use super::models::{SessionSettings, SessionStatus, SessionType};

/// Request to create a new learning session
#[derive(Debug, Deserialize)]
pub struct CreateSessionRequest {
    pub book_id: Uuid,
    pub session_type: SessionType,
    pub start_page: Option<i32>,
    pub end_page: Option<i32>,
    #[serde(default)]
    pub settings: Option<SessionSettings>,
}

/// Request to start a review session (no book)
#[derive(Debug, Deserialize)]
pub struct CreateReviewSessionRequest {
    #[serde(default)]
    pub settings: Option<SessionSettings>,
}

/// Request to update session progress
#[derive(Debug, Deserialize)]
pub struct UpdateProgressRequest {
    pub page_id: Uuid,
    pub time_spent_seconds: i32,
    pub pronunciation_score: Option<i32>,
    pub comprehension_score: Option<i32>,
}

/// Request to update session status
#[derive(Debug, Deserialize)]
pub struct UpdateSessionStatusRequest {
    pub action: SessionAction,
}

/// Available session actions
#[derive(Debug, Clone, Copy, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum SessionAction {
    Pause,
    Resume,
    Complete,
    Abandon,
    NextPage,
}

/// Learning session response
#[derive(Debug, Serialize)]
pub struct SessionResponse {
    pub id: Uuid,
    pub user_id: Uuid,
    pub book_id: Option<Uuid>,
    pub session_type: SessionType,
    pub current_page: i32,
    pub start_page: Option<i32>,
    pub end_page: Option<i32>,
    pub duration_seconds: i32,
    pub status: SessionStatus,
    pub started_at: DateTime<Utc>,
    pub ended_at: Option<DateTime<Utc>>,
}

/// Learning progress response
#[derive(Debug, Serialize)]
pub struct ProgressResponse {
    pub id: Uuid,
    pub session_id: Uuid,
    pub page_id: Uuid,
    pub time_spent_seconds: i32,
    pub pronunciation_score: Option<i32>,
    pub comprehension_score: Option<i32>,
    pub average_score: Option<f32>,
    pub created_at: DateTime<Utc>,
}

/// Session list response with pagination
#[derive(Debug, Serialize)]
pub struct SessionListResponse {
    pub sessions: Vec<SessionResponse>,
    pub total: usize,
    pub page: usize,
    pub per_page: usize,
}

/// Session statistics response
#[derive(Debug, Serialize)]
pub struct SessionStatsResponse {
    pub total_sessions: usize,
    pub completed_sessions: usize,
    pub total_duration_seconds: i64,
    pub average_score: Option<f32>,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_create_session_request_deserialization() {
        let json = r#"{
            "book_id": "550e8400-e29b-41d4-a716-446655440000",
            "session_type": "page_by_page"
        }"#;

        let request: CreateSessionRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.session_type, SessionType::PageByPage);
        assert!(request.start_page.is_none());
    }

    #[test]
    fn test_create_session_request_with_pages() {
        let json = r#"{
            "book_id": "550e8400-e29b-41d4-a716-446655440000",
            "session_type": "teacher_mode",
            "start_page": 1,
            "end_page": 10
        }"#;

        let request: CreateSessionRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.session_type, SessionType::TeacherMode);
        assert_eq!(request.start_page, Some(1));
        assert_eq!(request.end_page, Some(10));
    }

    #[test]
    fn test_update_progress_request_deserialization() {
        let json = r#"{
            "page_id": "550e8400-e29b-41d4-a716-446655440000",
            "time_spent_seconds": 120,
            "pronunciation_score": 85
        }"#;

        let request: UpdateProgressRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.time_spent_seconds, 120);
        assert_eq!(request.pronunciation_score, Some(85));
        assert!(request.comprehension_score.is_none());
    }

    #[test]
    fn test_session_action_deserialization() {
        let json = r#"{"action": "pause"}"#;
        let request: UpdateSessionStatusRequest = serde_json::from_str(json).unwrap();
        assert!(matches!(request.action, SessionAction::Pause));

        let json = r#"{"action": "complete"}"#;
        let request: UpdateSessionStatusRequest = serde_json::from_str(json).unwrap();
        assert!(matches!(request.action, SessionAction::Complete));
    }

    #[test]
    fn test_session_response_serialization() {
        let now = Utc::now();
        let response = SessionResponse {
            id: Uuid::new_v4(),
            user_id: Uuid::new_v4(),
            book_id: Some(Uuid::new_v4()),
            session_type: SessionType::PageByPage,
            current_page: 5,
            start_page: Some(1),
            end_page: Some(20),
            duration_seconds: 300,
            status: SessionStatus::Active,
            started_at: now,
            ended_at: None,
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("page_by_page"));
        assert!(json.contains("active"));
        assert!(json.contains("\"current_page\":5"));
    }

    #[test]
    fn test_progress_response_serialization() {
        let now = Utc::now();
        let response = ProgressResponse {
            id: Uuid::new_v4(),
            session_id: Uuid::new_v4(),
            page_id: Uuid::new_v4(),
            time_spent_seconds: 60,
            pronunciation_score: Some(90),
            comprehension_score: Some(85),
            average_score: Some(87.5),
            created_at: now,
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("\"time_spent_seconds\":60"));
        assert!(json.contains("\"pronunciation_score\":90"));
        assert!(json.contains("87.5"));
    }

    #[test]
    fn test_session_stats_response_serialization() {
        let response = SessionStatsResponse {
            total_sessions: 10,
            completed_sessions: 8,
            total_duration_seconds: 3600,
            average_score: Some(82.5),
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("\"total_sessions\":10"));
        assert!(json.contains("\"completed_sessions\":8"));
    }
}
