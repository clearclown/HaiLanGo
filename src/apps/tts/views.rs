//! TTS ViewSet - Text-to-speech synthesis endpoints

use uuid::Uuid;

use super::dto::{
    GenerationHistoryResponse, SupportedLanguagesResponse, SynthesizeMetadataResponse,
    SynthesizeRequest,
};
use super::models::{AudioGeneration, GenerationStatus};
use crate::services::tts::{TtsProvider, TtsRequest, TtsResponse};

/// Synthesis result (metadata + raw audio bytes returned separately)
#[derive(Debug)]
pub enum SynthesizeResult {
    Success {
        metadata: SynthesizeMetadataResponse,
        audio_data: Vec<u8>,
        generation: Box<AudioGeneration>,
    },
    InvalidInput(String),
    ServiceError(String),
}

/// TTS ViewSet
pub struct TtsViewSet;

impl TtsViewSet {
    /// Synthesize speech from text
    pub async fn synthesize(
        user_id: Uuid,
        request: SynthesizeRequest,
        provider: &dyn TtsProvider,
    ) -> SynthesizeResult {
        // Validate
        if request.text.is_empty() {
            return SynthesizeResult::InvalidInput("Text is required".to_string());
        }
        if request.text.len() > 5000 {
            return SynthesizeResult::InvalidInput(format!(
                "Text too long: {} chars (max 5000)",
                request.text.len()
            ));
        }
        if !provider.supports_language(&request.language) {
            return SynthesizeResult::InvalidInput(format!(
                "Language '{}' is not supported by provider '{}'",
                request.language,
                provider.provider_name()
            ));
        }

        // Create generation record
        let mut generation =
            AudioGeneration::new(user_id, request.text.clone(), request.language.clone())
                .with_speed(request.speed)
                .with_quality(request.quality);
        generation.format = request.format;
        generation.provider = provider.provider_name().to_string();

        if let Some(page_id) = request.page_id {
            generation = generation.with_page(page_id);
        }

        generation.status = GenerationStatus::Processing;

        // Build service-layer request
        let tts_req = TtsRequest::new(request.text, request.language)
            .with_speed(request.speed)
            .with_format(request.format)
            .with_quality(request.quality);

        // Synthesize
        let tts_resp: TtsResponse = match provider.synthesize(tts_req).await {
            Ok(r) => r,
            Err(e) => {
                generation.status = GenerationStatus::Failed;
                return SynthesizeResult::ServiceError(e.to_string());
            }
        };

        generation.status = GenerationStatus::Completed;
        generation.duration_ms = Some(tts_resp.duration_ms);
        generation.audio_size_bytes = Some(tts_resp.audio_data.len());

        let metadata = SynthesizeMetadataResponse {
            id: generation.id,
            language: tts_resp.language,
            format: tts_resp.format,
            quality: generation.quality,
            provider: generation.provider.clone(),
            duration_ms: tts_resp.duration_ms,
            audio_size_bytes: tts_resp.audio_data.len(),
            content_type: tts_resp.format.content_type().to_string(),
        };

        SynthesizeResult::Success {
            metadata,
            audio_data: tts_resp.audio_data,
            generation: Box::new(generation),
        }
    }

    /// List generation history for a user
    pub fn list_history(
        user_id: Uuid,
        generations: &[AudioGeneration],
    ) -> Vec<GenerationHistoryResponse> {
        generations
            .iter()
            .filter(|g| g.user_id == user_id)
            .map(|g| GenerationHistoryResponse {
                id: g.id,
                text: g.text.clone(),
                language: g.language.clone(),
                speed: g.speed,
                format: g.format,
                quality: g.quality,
                status: g.status,
                duration_ms: g.duration_ms,
                created_at: g.created_at,
            })
            .collect()
    }

    /// Get supported languages for a provider
    pub fn supported_languages(provider: &dyn TtsProvider) -> SupportedLanguagesResponse {
        SupportedLanguagesResponse {
            provider: provider.provider_name().to_string(),
            languages: provider
                .supported_languages()
                .into_iter()
                .map(String::from)
                .collect(),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::services::tts::{AudioFormat, MockTtsProvider, QualityTier};

    fn make_request(text: &str, language: &str) -> SynthesizeRequest {
        SynthesizeRequest {
            text: text.to_string(),
            language: language.to_string(),
            speed: 1.0,
            format: AudioFormat::Mp3,
            quality: QualityTier::Standard,
            page_id: None,
        }
    }

    #[tokio::test]
    async fn test_synthesize_success() {
        let provider = MockTtsProvider::new();
        let user_id = Uuid::new_v4();
        let request = make_request("Hello world", "en");

        let result = TtsViewSet::synthesize(user_id, request, &provider).await;

        match result {
            SynthesizeResult::Success {
                metadata,
                audio_data,
                generation,
            } => {
                assert_eq!(metadata.language, "en");
                assert_eq!(metadata.format, AudioFormat::Mp3);
                assert_eq!(metadata.provider, "mock");
                assert!(metadata.duration_ms > 0);
                assert!(!audio_data.is_empty());
                assert_eq!(generation.status, GenerationStatus::Completed);
                assert_eq!(generation.user_id, user_id);
            }
            other => panic!("Expected Success, got {:?}", other),
        }
    }

    #[tokio::test]
    async fn test_synthesize_with_page_id() {
        let provider = MockTtsProvider::new();
        let user_id = Uuid::new_v4();
        let page_id = Uuid::new_v4();
        let request = SynthesizeRequest {
            text: "Test".to_string(),
            language: "en".to_string(),
            speed: 1.0,
            format: AudioFormat::Mp3,
            quality: QualityTier::Standard,
            page_id: Some(page_id),
        };

        let result = TtsViewSet::synthesize(user_id, request, &provider).await;

        match result {
            SynthesizeResult::Success { generation, .. } => {
                assert_eq!(generation.page_id, Some(page_id));
            }
            other => panic!("Expected Success, got {:?}", other),
        }
    }

    #[tokio::test]
    async fn test_synthesize_premium_quality() {
        let provider = MockTtsProvider::new();
        let user_id = Uuid::new_v4();
        let request = SynthesizeRequest {
            text: "Premium test".to_string(),
            language: "ja".to_string(),
            speed: 0.8,
            format: AudioFormat::Ogg,
            quality: QualityTier::Premium,
            page_id: None,
        };

        let result = TtsViewSet::synthesize(user_id, request, &provider).await;

        match result {
            SynthesizeResult::Success { metadata, .. } => {
                assert_eq!(metadata.format, AudioFormat::Ogg);
                assert_eq!(metadata.quality, QualityTier::Premium);
            }
            other => panic!("Expected Success, got {:?}", other),
        }
    }

    #[tokio::test]
    async fn test_synthesize_empty_text() {
        let provider = MockTtsProvider::new();
        let user_id = Uuid::new_v4();
        let request = make_request("", "en");

        let result = TtsViewSet::synthesize(user_id, request, &provider).await;

        assert!(matches!(result, SynthesizeResult::InvalidInput(_)));
    }

    #[tokio::test]
    async fn test_synthesize_text_too_long() {
        let provider = MockTtsProvider::new();
        let user_id = Uuid::new_v4();
        let long_text = "a".repeat(6000);
        let request = make_request(&long_text, "en");

        let result = TtsViewSet::synthesize(user_id, request, &provider).await;

        assert!(matches!(result, SynthesizeResult::InvalidInput(_)));
    }

    #[tokio::test]
    async fn test_synthesize_unsupported_language() {
        let provider = MockTtsProvider::new();
        let user_id = Uuid::new_v4();
        let request = make_request("test", "xyz");

        let result = TtsViewSet::synthesize(user_id, request, &provider).await;

        assert!(matches!(result, SynthesizeResult::InvalidInput(_)));
    }

    #[test]
    fn test_list_history() {
        let user_id = Uuid::new_v4();
        let other_id = Uuid::new_v4();

        let generations = vec![
            AudioGeneration::new(user_id, "Hello".to_string(), "en".to_string()),
            AudioGeneration::new(user_id, "World".to_string(), "en".to_string()),
            AudioGeneration::new(other_id, "Other".to_string(), "fr".to_string()),
        ];

        let history = TtsViewSet::list_history(user_id, &generations);

        assert_eq!(history.len(), 2);
        assert!(history.iter().all(|h| h.text != "Other"));
    }

    #[test]
    fn test_list_history_empty() {
        let user_id = Uuid::new_v4();
        let history = TtsViewSet::list_history(user_id, &[]);
        assert!(history.is_empty());
    }

    #[test]
    fn test_supported_languages() {
        let provider = MockTtsProvider::new();
        let result = TtsViewSet::supported_languages(&provider);

        assert_eq!(result.provider, "mock");
        assert!(result.languages.contains(&"en".to_string()));
        assert!(result.languages.contains(&"ja".to_string()));
    }
}
