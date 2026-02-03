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

/// Google Cloud Text-to-Speech provider
pub struct GoogleCloudTtsProvider {
    api_key: String,
    endpoint: String,
}

impl GoogleCloudTtsProvider {
    /// Create a new Google Cloud TTS provider
    pub fn new(api_key: &str) -> Self {
        Self {
            api_key: api_key.to_string(),
            endpoint: "https://texttospeech.googleapis.com/v1/text:synthesize".to_string(),
        }
    }

    /// Create from environment variable
    pub fn from_env() -> Option<Self> {
        std::env::var("GOOGLE_CLOUD_TTS_API_KEY")
            .ok()
            .map(|key| Self::new(&key))
    }

    /// Map language code to Google TTS voice configuration
    fn get_voice_config(&self, language: &str) -> (String, String) {
        match language {
            "zh" => ("zh-CN".to_string(), "zh-CN-Standard-A".to_string()),
            "ja" => ("ja-JP".to_string(), "ja-JP-Standard-A".to_string()),
            "ko" => ("ko-KR".to_string(), "ko-KR-Standard-A".to_string()),
            "es" => ("es-ES".to_string(), "es-ES-Standard-A".to_string()),
            "fr" => ("fr-FR".to_string(), "fr-FR-Standard-A".to_string()),
            "de" => ("de-DE".to_string(), "de-DE-Standard-A".to_string()),
            "ru" => ("ru-RU".to_string(), "ru-RU-Standard-A".to_string()),
            "ar" => ("ar-XA".to_string(), "ar-XA-Standard-A".to_string()),
            _ => ("en-US".to_string(), "en-US-Standard-A".to_string()),
        }
    }

    /// Map AudioFormat to Google's encoding
    fn get_audio_encoding(&self, format: AudioFormat) -> &'static str {
        match format {
            AudioFormat::Mp3 => "MP3",
            AudioFormat::Wav => "LINEAR16",
            AudioFormat::Ogg => "OGG_OPUS",
        }
    }
}

#[async_trait]
impl TtsProvider for GoogleCloudTtsProvider {
    async fn synthesize(&self, request: TtsRequest) -> Result<TtsResponse, TtsError> {
        if request.text.len() > 5000 {
            return Err(TtsError::TextTooLong(request.text.len(), 5000));
        }

        let (language_code, voice_name) = self.get_voice_config(&request.language);
        let audio_encoding = self.get_audio_encoding(request.format);

        let body = serde_json::json!({
            "input": {
                "text": request.text
            },
            "voice": {
                "languageCode": language_code,
                "name": voice_name
            },
            "audioConfig": {
                "audioEncoding": audio_encoding,
                "speakingRate": request.speed
            }
        });

        let client = reqwest::Client::new();
        let url = format!("{}?key={}", self.endpoint, self.api_key);

        let response = client
            .post(&url)
            .json(&body)
            .send()
            .await
            .map_err(|_| TtsError::ServiceUnavailable)?;

        if response.status() == 429 {
            return Err(TtsError::RateLimitExceeded);
        }

        if !response.status().is_success() {
            return Err(TtsError::GenerationFailed(format!(
                "API returned status {}",
                response.status()
            )));
        }

        let result: serde_json::Value = response
            .json()
            .await
            .map_err(|e| TtsError::GenerationFailed(e.to_string()))?;

        let audio_content = result["audioContent"].as_str().ok_or_else(|| {
            TtsError::GenerationFailed("No audio content in response".to_string())
        })?;

        use base64::{Engine, engine::general_purpose::STANDARD};
        let audio_data = STANDARD
            .decode(audio_content)
            .map_err(|e| TtsError::GenerationFailed(e.to_string()))?;

        // Estimate duration based on text length and speed
        let base_duration_ms = (request.text.len() as u64 * 60) / request.speed.max(0.5) as u64;

        Ok(TtsResponse {
            audio_data,
            format: request.format,
            duration_ms: base_duration_ms,
            language: request.language,
        })
    }

    fn supports_language(&self, language: &str) -> bool {
        self.supported_languages().contains(&language)
    }

    fn supported_languages(&self) -> Vec<&'static str> {
        vec![
            "en", "ja", "zh", "ko", "es", "fr", "de", "ru", "ar", "he", "fa", "it", "pt", "nl",
            "pl", "tr", "vi", "th", "id",
        ]
    }
}

/// Factory function to create the appropriate TTS provider
pub fn create_tts_provider() -> Box<dyn TtsProvider> {
    if let Some(provider) = GoogleCloudTtsProvider::from_env() {
        Box::new(provider)
    } else {
        tracing::warn!("GOOGLE_CLOUD_TTS_API_KEY not set, using mock TTS provider");
        Box::new(MockTtsProvider::new())
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
