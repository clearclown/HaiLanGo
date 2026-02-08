//! Teacher Mode DTOs - Request and response types

use serde::{Deserialize, Serialize};
use uuid::Uuid;

use super::models::{PlaybackConfig, TeacherSessionStatus};

/// Request to start a teacher mode lesson
#[derive(Debug, Deserialize)]
pub struct StartLessonRequest {
    pub book_id: Uuid,
    #[serde(default = "default_start_page")]
    pub start_page: u32,
    pub end_page: u32,
    #[serde(default)]
    pub language: Option<String>,
    #[serde(default)]
    pub speed: Option<f32>,
    #[serde(default)]
    pub page_interval_secs: Option<u32>,
    #[serde(default)]
    pub repeat_count: Option<u32>,
    #[serde(default)]
    pub auto_advance: Option<bool>,
}

fn default_start_page() -> u32 {
    1
}

impl StartLessonRequest {
    /// Build a PlaybackConfig from request fields
    pub fn to_config(&self) -> PlaybackConfig {
        let mut config = PlaybackConfig::default();
        if let Some(lang) = &self.language {
            config.language = lang.clone();
        }
        if let Some(speed) = self.speed {
            config.speed = speed;
        }
        if let Some(interval) = self.page_interval_secs {
            config.page_interval_secs = interval;
        }
        if let Some(repeat) = self.repeat_count {
            config.repeat_count = repeat;
        }
        if let Some(auto) = self.auto_advance {
            config.auto_advance = auto;
        }
        config.validated()
    }
}

/// Request to update playback configuration
#[derive(Debug, Deserialize)]
pub struct UpdateConfigRequest {
    #[serde(default)]
    pub speed: Option<f32>,
    #[serde(default)]
    pub page_interval_secs: Option<u32>,
    #[serde(default)]
    pub repeat_count: Option<u32>,
    #[serde(default)]
    pub auto_advance: Option<bool>,
    #[serde(default)]
    pub language: Option<String>,
}

/// Response for lesson status
#[derive(Debug, Serialize)]
pub struct LessonStatusResponse {
    pub id: Uuid,
    pub book_id: Uuid,
    pub status: TeacherSessionStatus,
    pub current_page: u32,
    pub start_page: u32,
    pub end_page: u32,
    pub total_pages: u32,
    pub pages_completed: u32,
    pub progress_percent: f32,
    pub config: PlaybackConfig,
}

/// Response for session in history list
#[derive(Debug, Serialize)]
pub struct SessionSummaryResponse {
    pub id: Uuid,
    pub book_id: Uuid,
    pub status: TeacherSessionStatus,
    pub start_page: u32,
    pub end_page: u32,
    pub pages_completed: u32,
    pub total_pages: u32,
    pub progress_percent: f32,
    pub created_at: chrono::DateTime<chrono::Utc>,
}

/// Response for start/action operations
#[derive(Debug, Serialize)]
pub struct LessonActionResponse {
    pub id: Uuid,
    pub status: TeacherSessionStatus,
    pub current_page: u32,
    pub message: String,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_start_lesson_request_to_config_defaults() {
        let request = StartLessonRequest {
            book_id: Uuid::new_v4(),
            start_page: 1,
            end_page: 10,
            language: None,
            speed: None,
            page_interval_secs: None,
            repeat_count: None,
            auto_advance: None,
        };

        let config = request.to_config();
        assert_eq!(config.speed, 1.0);
        assert_eq!(config.page_interval_secs, 5);
        assert_eq!(config.repeat_count, 1);
        assert!(config.auto_advance);
        assert_eq!(config.language, "en");
    }

    #[test]
    fn test_start_lesson_request_to_config_custom() {
        let request = StartLessonRequest {
            book_id: Uuid::new_v4(),
            start_page: 1,
            end_page: 10,
            language: Some("ja".to_string()),
            speed: Some(1.5),
            page_interval_secs: Some(10),
            repeat_count: Some(2),
            auto_advance: Some(false),
        };

        let config = request.to_config();
        assert_eq!(config.speed, 1.5);
        assert_eq!(config.page_interval_secs, 10);
        assert_eq!(config.repeat_count, 2);
        assert!(!config.auto_advance);
        assert_eq!(config.language, "ja");
    }

    #[test]
    fn test_start_lesson_request_to_config_clamped() {
        let request = StartLessonRequest {
            book_id: Uuid::new_v4(),
            start_page: 1,
            end_page: 10,
            language: None,
            speed: Some(5.0),
            page_interval_secs: Some(100),
            repeat_count: Some(99),
            auto_advance: None,
        };

        let config = request.to_config();
        assert_eq!(config.speed, 2.0);
        assert_eq!(config.page_interval_secs, 30);
        assert_eq!(config.repeat_count, 3);
    }

    #[test]
    fn test_lesson_status_response_serialization() {
        let resp = LessonStatusResponse {
            id: Uuid::nil(),
            book_id: Uuid::nil(),
            status: TeacherSessionStatus::Playing,
            current_page: 3,
            start_page: 1,
            end_page: 10,
            total_pages: 10,
            pages_completed: 2,
            progress_percent: 20.0,
            config: PlaybackConfig::default(),
        };

        let json = serde_json::to_string(&resp).unwrap();
        assert!(json.contains("\"playing\""));
        assert!(json.contains("\"current_page\":3"));
    }

    #[test]
    fn test_lesson_action_response_serialization() {
        let resp = LessonActionResponse {
            id: Uuid::nil(),
            status: TeacherSessionStatus::Paused,
            current_page: 5,
            message: "Lesson paused".to_string(),
        };

        let json = serde_json::to_string(&resp).unwrap();
        assert!(json.contains("\"paused\""));
        assert!(json.contains("Lesson paused"));
    }

    #[test]
    fn test_start_lesson_request_deserialization() {
        let json = r#"{"book_id":"00000000-0000-0000-0000-000000000000","end_page":10}"#;
        let request: StartLessonRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.start_page, 1); // default
        assert_eq!(request.end_page, 10);
        assert!(request.speed.is_none());
    }

    #[test]
    fn test_update_config_request_deserialization() {
        let json = r#"{"speed":1.5}"#;
        let request: UpdateConfigRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.speed, Some(1.5));
        assert!(request.page_interval_secs.is_none());
        assert!(request.repeat_count.is_none());
    }
}
