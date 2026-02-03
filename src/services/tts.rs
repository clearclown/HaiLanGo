//! Text-to-Speech service abstraction

use async_trait::async_trait;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum TtsError {
    #[error("Service unavailable")]
    ServiceUnavailable,
    #[error("Unsupported language: {0}")]
    UnsupportedLanguage(String),
    #[error("Text too long: {0} chars (max: {1})")]
    TextTooLong(usize, usize),
    #[error("Generation failed: {0}")]
    GenerationFailed(String),
    #[error("Rate limit exceeded")]
    RateLimitExceeded,
}

/// Audio format for TTS output
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default)]
pub enum AudioFormat {
    #[default]
    Mp3,
    Wav,
    Ogg,
}

/// TTS request configuration
#[derive(Debug, Clone)]
pub struct TtsRequest {
    pub text: String,
    pub language: String,
    pub speed: f32,
    pub format: AudioFormat,
}

impl TtsRequest {
    pub fn new(text: String, language: String) -> Self {
        Self {
            text,
            language,
            speed: 1.0,
            format: AudioFormat::default(),
        }
    }

    pub fn with_speed(mut self, speed: f32) -> Self {
        self.speed = speed.clamp(0.5, 2.0);
        self
    }
}

/// TTS response with audio data
#[derive(Debug)]
pub struct TtsResponse {
    pub audio_data: Vec<u8>,
    pub format: AudioFormat,
    pub duration_ms: u64,
    pub language: String,
}

/// TTS Provider trait - implement for different providers
#[async_trait]
pub trait TtsProvider: Send + Sync {
    /// Generate speech from text
    async fn synthesize(&self, request: TtsRequest) -> Result<TtsResponse, TtsError>;

    /// Check if language is supported
    fn supports_language(&self, language: &str) -> bool;

    /// Get list of supported languages
    fn supported_languages(&self) -> Vec<&'static str>;
}

/// Mock TTS provider for development and testing
pub struct MockTtsProvider;

impl MockTtsProvider {
    pub fn new() -> Self {
        Self
    }
}

impl Default for MockTtsProvider {
    fn default() -> Self {
        Self::new()
    }
}

#[async_trait]
impl TtsProvider for MockTtsProvider {
    async fn synthesize(&self, request: TtsRequest) -> Result<TtsResponse, TtsError> {
        if !self.supports_language(&request.language) {
            return Err(TtsError::UnsupportedLanguage(request.language));
        }

        // Simulate audio generation (1 byte per character, 100ms per 10 chars)
        let audio_data = vec![0u8; request.text.len()];
        let duration_ms = (request.text.len() as u64 / 10).max(100) * 100;

        Ok(TtsResponse {
            audio_data,
            format: request.format,
            duration_ms,
            language: request.language,
        })
    }

    fn supports_language(&self, language: &str) -> bool {
        self.supported_languages().contains(&language)
    }

    fn supported_languages(&self) -> Vec<&'static str> {
        vec![
            "en", "ja", "zh", "ko", "es", "fr", "de", "ru", "ar", "he", "fa",
        ]
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_tts_request_creation() {
        let request = TtsRequest::new("Hello world".to_string(), "en".to_string());

        assert_eq!(request.text, "Hello world");
        assert_eq!(request.language, "en");
        assert_eq!(request.speed, 1.0);
    }

    #[test]
    fn test_tts_request_speed_clamping() {
        let request = TtsRequest::new("Test".to_string(), "en".to_string()).with_speed(3.0); // Should clamp to 2.0

        assert_eq!(request.speed, 2.0);

        let request2 = TtsRequest::new("Test".to_string(), "en".to_string()).with_speed(0.1); // Should clamp to 0.5

        assert_eq!(request2.speed, 0.5);
    }

    #[tokio::test]
    async fn test_mock_tts_synthesize() {
        let provider = MockTtsProvider::new();
        let request = TtsRequest::new("Hello, how are you?".to_string(), "en".to_string());

        let response = provider.synthesize(request).await.unwrap();

        assert!(!response.audio_data.is_empty());
        assert!(response.duration_ms > 0);
        assert_eq!(response.language, "en");
    }

    #[tokio::test]
    async fn test_mock_tts_unsupported_language() {
        let provider = MockTtsProvider::new();
        let request = TtsRequest::new("Test".to_string(), "xyz".to_string());

        let result = provider.synthesize(request).await;

        assert!(matches!(result, Err(TtsError::UnsupportedLanguage(_))));
    }

    #[test]
    fn test_mock_tts_supported_languages() {
        let provider = MockTtsProvider::new();

        assert!(provider.supports_language("en"));
        assert!(provider.supports_language("ja"));
        assert!(!provider.supports_language("xyz"));
    }

    #[test]
    fn test_audio_format_default() {
        let format = AudioFormat::default();
        assert_eq!(format, AudioFormat::Mp3);
    }
}
