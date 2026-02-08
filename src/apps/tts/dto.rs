//! Data Transfer Objects for TTS API

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use crate::services::tts::{AudioFormat, QualityTier};

use super::models::GenerationStatus;

/// Request to synthesize speech
#[derive(Debug, Clone, Deserialize)]
pub struct SynthesizeRequest {
    pub text: String,
    pub language: String,
    #[serde(default = "default_speed")]
    pub speed: f32,
    #[serde(default)]
    pub format: AudioFormat,
    #[serde(default)]
    pub quality: QualityTier,
    pub page_id: Option<Uuid>,
}

fn default_speed() -> f32 {
    1.0
}

/// Metadata response for a completed synthesis (without audio bytes)
#[derive(Debug, Clone, Serialize)]
pub struct SynthesizeMetadataResponse {
    pub id: Uuid,
    pub language: String,
    pub format: AudioFormat,
    pub quality: QualityTier,
    pub provider: String,
    pub duration_ms: u64,
    pub audio_size_bytes: usize,
    pub content_type: String,
}

/// History entry for a past generation
#[derive(Debug, Clone, Serialize)]
pub struct GenerationHistoryResponse {
    pub id: Uuid,
    pub text: String,
    pub language: String,
    pub speed: f32,
    pub format: AudioFormat,
    pub quality: QualityTier,
    pub status: GenerationStatus,
    pub duration_ms: Option<u64>,
    pub created_at: DateTime<Utc>,
}

/// Supported languages response
#[derive(Debug, Clone, Serialize)]
pub struct SupportedLanguagesResponse {
    pub provider: String,
    pub languages: Vec<String>,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_synthesize_request_defaults() {
        let json = r#"{"text":"hello","language":"en"}"#;
        let req: SynthesizeRequest = serde_json::from_str(json).unwrap();

        assert_eq!(req.text, "hello");
        assert_eq!(req.language, "en");
        assert_eq!(req.speed, 1.0);
        assert_eq!(req.format, AudioFormat::Mp3);
        assert_eq!(req.quality, QualityTier::Standard);
        assert!(req.page_id.is_none());
    }

    #[test]
    fn test_synthesize_request_full() {
        let page_id = Uuid::new_v4();
        let json = format!(
            r#"{{"text":"bonjour","language":"fr","speed":1.5,"format":"ogg","quality":"premium","page_id":"{}"}}"#,
            page_id
        );
        let req: SynthesizeRequest = serde_json::from_str(&json).unwrap();

        assert_eq!(req.text, "bonjour");
        assert_eq!(req.language, "fr");
        assert_eq!(req.speed, 1.5);
        assert_eq!(req.format, AudioFormat::Ogg);
        assert_eq!(req.quality, QualityTier::Premium);
        assert_eq!(req.page_id, Some(page_id));
    }

    #[test]
    fn test_metadata_response_serialization() {
        let resp = SynthesizeMetadataResponse {
            id: Uuid::new_v4(),
            language: "en".to_string(),
            format: AudioFormat::Mp3,
            quality: QualityTier::Standard,
            provider: "mock".to_string(),
            duration_ms: 3000,
            audio_size_bytes: 48000,
            content_type: "audio/mpeg".to_string(),
        };

        let json = serde_json::to_string(&resp).unwrap();
        assert!(json.contains("\"duration_ms\":3000"));
        assert!(json.contains("\"mp3\""));
        assert!(json.contains("\"standard\""));
    }

    #[test]
    fn test_history_response_serialization() {
        let resp = GenerationHistoryResponse {
            id: Uuid::new_v4(),
            text: "Hello".to_string(),
            language: "en".to_string(),
            speed: 1.0,
            format: AudioFormat::Mp3,
            quality: QualityTier::Standard,
            status: GenerationStatus::Completed,
            duration_ms: Some(2000),
            created_at: Utc::now(),
        };

        let json = serde_json::to_string(&resp).unwrap();
        assert!(json.contains("\"completed\""));
        assert!(json.contains("\"Hello\""));
    }

    #[test]
    fn test_supported_languages_response() {
        let resp = SupportedLanguagesResponse {
            provider: "google".to_string(),
            languages: vec!["en".to_string(), "ja".to_string()],
        };

        let json = serde_json::to_string(&resp).unwrap();
        assert!(json.contains("\"google\""));
        assert!(json.contains("\"ja\""));
    }
}
