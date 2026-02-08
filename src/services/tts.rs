//! Text-to-Speech service abstraction

use async_trait::async_trait;
use serde::{Deserialize, Serialize};
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
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum AudioFormat {
    #[default]
    Mp3,
    Wav,
    Ogg,
}

impl AudioFormat {
    /// Content-Type header value
    pub fn content_type(&self) -> &'static str {
        match self {
            Self::Mp3 => "audio/mpeg",
            Self::Wav => "audio/wav",
            Self::Ogg => "audio/ogg",
        }
    }

    /// File extension
    pub fn extension(&self) -> &'static str {
        match self {
            Self::Mp3 => "mp3",
            Self::Wav => "wav",
            Self::Ogg => "ogg",
        }
    }
}

/// TTS quality tier
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum QualityTier {
    #[default]
    Standard,
    Premium,
}

/// TTS request configuration
#[derive(Debug, Clone)]
pub struct TtsRequest {
    pub text: String,
    pub language: String,
    pub speed: f32,
    pub format: AudioFormat,
    pub quality: QualityTier,
}

impl TtsRequest {
    pub fn new(text: String, language: String) -> Self {
        Self {
            text,
            language,
            speed: 1.0,
            format: AudioFormat::default(),
            quality: QualityTier::default(),
        }
    }

    pub fn with_speed(mut self, speed: f32) -> Self {
        self.speed = speed.clamp(0.5, 2.0);
        self
    }

    pub fn with_format(mut self, format: AudioFormat) -> Self {
        self.format = format;
        self
    }

    pub fn with_quality(mut self, quality: QualityTier) -> Self {
        self.quality = quality;
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

    /// Provider display name
    fn provider_name(&self) -> &'static str;
}

// ---------------------------------------------------------------------------
// Mock provider
// ---------------------------------------------------------------------------

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

    fn provider_name(&self) -> &'static str {
        "mock"
    }
}

// ---------------------------------------------------------------------------
// Google Cloud TTS provider
// ---------------------------------------------------------------------------

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
    fn get_voice_config(&self, language: &str, quality: QualityTier) -> (String, String) {
        let tier = match quality {
            QualityTier::Standard => "Standard",
            QualityTier::Premium => "Wavenet",
        };
        match language {
            "zh" => ("zh-CN".to_string(), format!("zh-CN-{}-A", tier)),
            "ja" => ("ja-JP".to_string(), format!("ja-JP-{}-A", tier)),
            "ko" => ("ko-KR".to_string(), format!("ko-KR-{}-A", tier)),
            "es" => ("es-ES".to_string(), format!("es-ES-{}-A", tier)),
            "fr" => ("fr-FR".to_string(), format!("fr-FR-{}-A", tier)),
            "de" => ("de-DE".to_string(), format!("de-DE-{}-A", tier)),
            "ru" => ("ru-RU".to_string(), format!("ru-RU-{}-A", tier)),
            "ar" => ("ar-XA".to_string(), format!("ar-XA-{}-A", tier)),
            "he" => ("he-IL".to_string(), format!("he-IL-{}-A", tier)),
            "fa" => ("fa-IR".to_string(), format!("fa-IR-{}-A", tier)),
            "it" => ("it-IT".to_string(), format!("it-IT-{}-A", tier)),
            "pt" => ("pt-BR".to_string(), format!("pt-BR-{}-A", tier)),
            "tr" => ("tr-TR".to_string(), format!("tr-TR-{}-A", tier)),
            _ => ("en-US".to_string(), format!("en-US-{}-A", tier)),
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

        let (language_code, voice_name) =
            self.get_voice_config(&request.language, request.quality);
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

    fn provider_name(&self) -> &'static str {
        "google"
    }
}

// ---------------------------------------------------------------------------
// ElevenLabs TTS provider (premium)
// ---------------------------------------------------------------------------

/// ElevenLabs TTS provider for high-quality voice synthesis
pub struct ElevenLabsTtsProvider {
    api_key: String,
    endpoint: String,
    voice_id: String,
}

impl ElevenLabsTtsProvider {
    /// Create a new ElevenLabs TTS provider with a specific voice
    pub fn new(api_key: &str, voice_id: &str) -> Self {
        Self {
            api_key: api_key.to_string(),
            endpoint: "https://api.elevenlabs.io/v1/text-to-speech".to_string(),
            voice_id: voice_id.to_string(),
        }
    }

    /// Create from environment variables
    pub fn from_env() -> Option<Self> {
        let api_key = std::env::var("ELEVENLABS_API_KEY").ok()?;
        let voice_id = std::env::var("ELEVENLABS_VOICE_ID")
            .unwrap_or_else(|_| "21m00Tcm4TlvDq8ikWAM".to_string()); // Rachel (default)
        Some(Self::new(&api_key, &voice_id))
    }

    /// Map AudioFormat to ElevenLabs output format
    fn output_format(format: AudioFormat) -> &'static str {
        match format {
            AudioFormat::Mp3 => "mp3_44100_128",
            AudioFormat::Ogg => "ogg_vorbis",
            AudioFormat::Wav => "pcm_44100",
        }
    }

    /// Map language code to ElevenLabs model
    fn model_for_language(language: &str) -> &'static str {
        match language {
            // Multilingual v2 supports many languages
            "en" => "eleven_turbo_v2_5",
            _ => "eleven_multilingual_v2",
        }
    }
}

#[async_trait]
impl TtsProvider for ElevenLabsTtsProvider {
    async fn synthesize(&self, request: TtsRequest) -> Result<TtsResponse, TtsError> {
        if request.text.len() > 5000 {
            return Err(TtsError::TextTooLong(request.text.len(), 5000));
        }

        let model_id = Self::model_for_language(&request.language);
        let output_format = Self::output_format(request.format);

        let url = format!(
            "{}/{}?output_format={}",
            self.endpoint, self.voice_id, output_format
        );

        // ElevenLabs speed: stability/similarity_boost affect quality, not speed directly.
        // We approximate speed via the speaking rate range.
        let body = serde_json::json!({
            "text": request.text,
            "model_id": model_id,
            "voice_settings": {
                "stability": 0.5,
                "similarity_boost": 0.75,
                "speed": request.speed
            }
        });

        let client = reqwest::Client::new();
        let response = client
            .post(&url)
            .header("xi-api-key", &self.api_key)
            .json(&body)
            .send()
            .await
            .map_err(|_| TtsError::ServiceUnavailable)?;

        if response.status() == 429 {
            return Err(TtsError::RateLimitExceeded);
        }

        if !response.status().is_success() {
            let status = response.status();
            let body_text: String = response
                .text()
                .await
                .unwrap_or_else(|_| "Unknown error".to_string());
            return Err(TtsError::GenerationFailed(format!(
                "ElevenLabs API returned status {}: {}",
                status, body_text
            )));
        }

        let audio_data = response
            .bytes()
            .await
            .map_err(|e| TtsError::GenerationFailed(e.to_string()))?
            .to_vec();

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
        // ElevenLabs multilingual v2 supports 29 languages
        vec![
            "en", "ja", "zh", "ko", "es", "fr", "de", "ru", "ar", "he", "fa", "it", "pt", "nl",
            "pl", "tr", "vi", "sv", "da", "fi", "no", "uk", "hi", "id", "cs", "ro", "hu", "el",
            "bg",
        ]
    }

    fn provider_name(&self) -> &'static str {
        "elevenlabs"
    }
}

// ---------------------------------------------------------------------------
// Factory
// ---------------------------------------------------------------------------

/// Factory function to create the appropriate TTS provider
pub fn create_tts_provider() -> Box<dyn TtsProvider> {
    // Prefer ElevenLabs for premium quality, fall back to Google, then mock
    if let Some(provider) = ElevenLabsTtsProvider::from_env() {
        tracing::info!("Using ElevenLabs TTS provider");
        Box::new(provider)
    } else if let Some(provider) = GoogleCloudTtsProvider::from_env() {
        tracing::info!("Using Google Cloud TTS provider");
        Box::new(provider)
    } else {
        tracing::warn!("No TTS API key set, using mock TTS provider");
        Box::new(MockTtsProvider::new())
    }
}

/// Create a provider by name (for explicit selection)
pub fn create_tts_provider_by_name(name: &str) -> Option<Box<dyn TtsProvider>> {
    match name {
        "elevenlabs" => ElevenLabsTtsProvider::from_env().map(|p| Box::new(p) as Box<dyn TtsProvider>),
        "google" => GoogleCloudTtsProvider::from_env().map(|p| Box::new(p) as Box<dyn TtsProvider>),
        "mock" => Some(Box::new(MockTtsProvider::new())),
        _ => None,
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
        assert_eq!(request.quality, QualityTier::Standard);
    }

    #[test]
    fn test_tts_request_builder() {
        let request = TtsRequest::new("Test".to_string(), "ja".to_string())
            .with_speed(1.5)
            .with_format(AudioFormat::Ogg)
            .with_quality(QualityTier::Premium);

        assert_eq!(request.speed, 1.5);
        assert_eq!(request.format, AudioFormat::Ogg);
        assert_eq!(request.quality, QualityTier::Premium);
    }

    #[test]
    fn test_tts_request_speed_clamping() {
        let request = TtsRequest::new("Test".to_string(), "en".to_string()).with_speed(3.0);
        assert_eq!(request.speed, 2.0);

        let request2 = TtsRequest::new("Test".to_string(), "en".to_string()).with_speed(0.1);
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

    #[test]
    fn test_audio_format_content_type() {
        assert_eq!(AudioFormat::Mp3.content_type(), "audio/mpeg");
        assert_eq!(AudioFormat::Wav.content_type(), "audio/wav");
        assert_eq!(AudioFormat::Ogg.content_type(), "audio/ogg");
    }

    #[test]
    fn test_audio_format_extension() {
        assert_eq!(AudioFormat::Mp3.extension(), "mp3");
        assert_eq!(AudioFormat::Wav.extension(), "wav");
        assert_eq!(AudioFormat::Ogg.extension(), "ogg");
    }

    #[test]
    fn test_audio_format_serialization() {
        let json = serde_json::to_string(&AudioFormat::Mp3).unwrap();
        assert_eq!(json, "\"mp3\"");

        let deserialized: AudioFormat = serde_json::from_str("\"ogg\"").unwrap();
        assert_eq!(deserialized, AudioFormat::Ogg);
    }

    #[test]
    fn test_quality_tier_serialization() {
        let json = serde_json::to_string(&QualityTier::Premium).unwrap();
        assert_eq!(json, "\"premium\"");

        let deserialized: QualityTier = serde_json::from_str("\"standard\"").unwrap();
        assert_eq!(deserialized, QualityTier::Standard);
    }

    #[test]
    fn test_quality_tier_default() {
        let tier = QualityTier::default();
        assert_eq!(tier, QualityTier::Standard);
    }

    #[test]
    fn test_provider_names() {
        assert_eq!(MockTtsProvider::new().provider_name(), "mock");
    }

    #[test]
    fn test_google_voice_config_quality() {
        let provider = GoogleCloudTtsProvider::new("test-key");
        let (_, voice_std) = provider.get_voice_config("en", QualityTier::Standard);
        assert!(voice_std.contains("Standard"));

        let (_, voice_prem) = provider.get_voice_config("en", QualityTier::Premium);
        assert!(voice_prem.contains("Wavenet"));
    }

    #[test]
    fn test_google_voice_config_languages() {
        let provider = GoogleCloudTtsProvider::new("test-key");

        let (code, _) = provider.get_voice_config("ja", QualityTier::Standard);
        assert_eq!(code, "ja-JP");

        let (code, _) = provider.get_voice_config("zh", QualityTier::Standard);
        assert_eq!(code, "zh-CN");

        let (code, _) = provider.get_voice_config("unknown", QualityTier::Standard);
        assert_eq!(code, "en-US");
    }

    #[test]
    fn test_elevenlabs_output_format() {
        assert_eq!(ElevenLabsTtsProvider::output_format(AudioFormat::Mp3), "mp3_44100_128");
        assert_eq!(ElevenLabsTtsProvider::output_format(AudioFormat::Ogg), "ogg_vorbis");
        assert_eq!(ElevenLabsTtsProvider::output_format(AudioFormat::Wav), "pcm_44100");
    }

    #[test]
    fn test_elevenlabs_model_selection() {
        assert_eq!(ElevenLabsTtsProvider::model_for_language("en"), "eleven_turbo_v2_5");
        assert_eq!(ElevenLabsTtsProvider::model_for_language("ja"), "eleven_multilingual_v2");
        assert_eq!(ElevenLabsTtsProvider::model_for_language("zh"), "eleven_multilingual_v2");
    }

    #[test]
    fn test_elevenlabs_supported_languages() {
        let provider = ElevenLabsTtsProvider::new("test-key", "test-voice");
        assert!(provider.supports_language("en"));
        assert!(provider.supports_language("ja"));
        assert!(provider.supports_language("ko"));
        assert!(!provider.supports_language("xyz"));
    }

    #[test]
    fn test_create_provider_by_name_mock() {
        let provider = create_tts_provider_by_name("mock");
        assert!(provider.is_some());
        assert_eq!(provider.unwrap().provider_name(), "mock");
    }

    #[test]
    fn test_create_provider_by_name_unknown() {
        let provider = create_tts_provider_by_name("unknown");
        assert!(provider.is_none());
    }

    #[test]
    fn test_tts_error_display() {
        let error = TtsError::ServiceUnavailable;
        assert_eq!(error.to_string(), "Service unavailable");

        let error = TtsError::UnsupportedLanguage("xyz".to_string());
        assert_eq!(error.to_string(), "Unsupported language: xyz");

        let error = TtsError::TextTooLong(6000, 5000);
        assert_eq!(error.to_string(), "Text too long: 6000 chars (max: 5000)");

        let error = TtsError::RateLimitExceeded;
        assert_eq!(error.to_string(), "Rate limit exceeded");
    }
}
