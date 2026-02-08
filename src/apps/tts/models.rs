//! TTS models - Audio generation records and cache metadata

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use crate::services::tts::{AudioFormat, QualityTier};

/// Status of an audio generation job
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(rename_all = "snake_case")]
pub enum GenerationStatus {
    #[default]
    Pending,
    Processing,
    Completed,
    Failed,
}

/// Record of a TTS generation request
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AudioGeneration {
    pub id: Uuid,
    pub user_id: Uuid,
    pub page_id: Option<Uuid>,
    pub text: String,
    pub language: String,
    pub speed: f32,
    pub format: AudioFormat,
    pub quality: QualityTier,
    pub provider: String,
    pub status: GenerationStatus,
    pub duration_ms: Option<u64>,
    pub audio_size_bytes: Option<usize>,
    pub created_at: DateTime<Utc>,
}

impl AudioGeneration {
    pub fn new(user_id: Uuid, text: String, language: String) -> Self {
        Self {
            id: Uuid::new_v4(),
            user_id,
            page_id: None,
            text,
            language,
            speed: 1.0,
            format: AudioFormat::default(),
            quality: QualityTier::default(),
            provider: "mock".to_string(),
            status: GenerationStatus::default(),
            duration_ms: None,
            audio_size_bytes: None,
            created_at: Utc::now(),
        }
    }

    pub fn with_page(mut self, page_id: Uuid) -> Self {
        self.page_id = Some(page_id);
        self
    }

    pub fn with_speed(mut self, speed: f32) -> Self {
        self.speed = speed.clamp(0.5, 2.0);
        self
    }

    pub fn with_quality(mut self, quality: QualityTier) -> Self {
        self.quality = quality;
        self
    }
}

/// Cached audio entry (keyed by text + language + format + quality hash)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AudioCache {
    pub id: Uuid,
    pub cache_key: String,
    pub language: String,
    pub format: AudioFormat,
    pub quality: QualityTier,
    pub audio_size_bytes: usize,
    pub duration_ms: u64,
    pub hit_count: u64,
    pub created_at: DateTime<Utc>,
    pub last_accessed_at: DateTime<Utc>,
}

impl AudioCache {
    /// Generate a cache key from synthesis parameters
    pub fn make_key(text: &str, language: &str, format: AudioFormat, quality: QualityTier) -> String {
        use std::collections::hash_map::DefaultHasher;
        use std::hash::{Hash, Hasher};

        let mut hasher = DefaultHasher::new();
        text.hash(&mut hasher);
        language.hash(&mut hasher);
        format!("{:?}", format).hash(&mut hasher);
        format!("{:?}", quality).hash(&mut hasher);
        format!("tts:{:016x}", hasher.finish())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_audio_generation_creation() {
        let user_id = Uuid::new_v4();
        let audio = AudioGeneration::new(
            user_id,
            "Hello world".to_string(),
            "en".to_string(),
        );

        assert_eq!(audio.user_id, user_id);
        assert_eq!(audio.text, "Hello world");
        assert_eq!(audio.language, "en");
        assert_eq!(audio.speed, 1.0);
        assert_eq!(audio.format, AudioFormat::Mp3);
        assert_eq!(audio.quality, QualityTier::Standard);
        assert_eq!(audio.status, GenerationStatus::Pending);
        assert!(audio.page_id.is_none());
        assert!(audio.duration_ms.is_none());
    }

    #[test]
    fn test_audio_generation_builders() {
        let user_id = Uuid::new_v4();
        let page_id = Uuid::new_v4();
        let audio = AudioGeneration::new(user_id, "Test".to_string(), "ja".to_string())
            .with_page(page_id)
            .with_speed(1.5)
            .with_quality(QualityTier::Premium);

        assert_eq!(audio.page_id, Some(page_id));
        assert_eq!(audio.speed, 1.5);
        assert_eq!(audio.quality, QualityTier::Premium);
    }

    #[test]
    fn test_audio_generation_speed_clamping() {
        let user_id = Uuid::new_v4();
        let audio = AudioGeneration::new(user_id, "T".to_string(), "en".to_string())
            .with_speed(5.0);
        assert_eq!(audio.speed, 2.0);

        let audio2 = AudioGeneration::new(user_id, "T".to_string(), "en".to_string())
            .with_speed(0.1);
        assert_eq!(audio2.speed, 0.5);
    }

    #[test]
    fn test_generation_status_default() {
        assert_eq!(GenerationStatus::default(), GenerationStatus::Pending);
    }

    #[test]
    fn test_generation_status_serialization() {
        let json = serde_json::to_string(&GenerationStatus::Completed).unwrap();
        assert_eq!(json, "\"completed\"");

        let de: GenerationStatus = serde_json::from_str("\"failed\"").unwrap();
        assert_eq!(de, GenerationStatus::Failed);
    }

    #[test]
    fn test_cache_key_deterministic() {
        let key1 = AudioCache::make_key("hello", "en", AudioFormat::Mp3, QualityTier::Standard);
        let key2 = AudioCache::make_key("hello", "en", AudioFormat::Mp3, QualityTier::Standard);
        assert_eq!(key1, key2);
    }

    #[test]
    fn test_cache_key_differs_by_text() {
        let key1 = AudioCache::make_key("hello", "en", AudioFormat::Mp3, QualityTier::Standard);
        let key2 = AudioCache::make_key("world", "en", AudioFormat::Mp3, QualityTier::Standard);
        assert_ne!(key1, key2);
    }

    #[test]
    fn test_cache_key_differs_by_quality() {
        let key1 = AudioCache::make_key("hello", "en", AudioFormat::Mp3, QualityTier::Standard);
        let key2 = AudioCache::make_key("hello", "en", AudioFormat::Mp3, QualityTier::Premium);
        assert_ne!(key1, key2);
    }

    #[test]
    fn test_cache_key_differs_by_format() {
        let key1 = AudioCache::make_key("hello", "en", AudioFormat::Mp3, QualityTier::Standard);
        let key2 = AudioCache::make_key("hello", "en", AudioFormat::Ogg, QualityTier::Standard);
        assert_ne!(key1, key2);
    }

    #[test]
    fn test_cache_key_prefix() {
        let key = AudioCache::make_key("test", "en", AudioFormat::Mp3, QualityTier::Standard);
        assert!(key.starts_with("tts:"));
    }
}
