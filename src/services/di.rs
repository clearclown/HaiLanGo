//! Service registry for dependency injection
//!
//! Provides a centralized container for application-wide service providers.
//! This lightweight registry bridges Reinhardt's Handler-based routing with
//! the trait-based external service abstractions (TTS, OCR, STT).
//!
//! In later phases this will integrate with Reinhardt DI (`reinhardt_di::Injectable`)
//! when ViewSets are adopted. For now it provides Arc-based service sharing.

use std::sync::Arc;

use crate::services::{
    MockOcrProvider, MockSttProvider, MockTtsProvider, OcrProvider, SttProvider, TtsProvider,
    create_ocr_provider, create_stt_provider, create_tts_provider,
};

/// Application-wide service registry.
///
/// Holds shared, thread-safe handles to external service providers.
/// Clone is cheap — all fields are `Arc`.
#[derive(Clone)]
pub struct ServiceRegistry {
    pub tts: Arc<dyn TtsProvider>,
    pub ocr: Arc<dyn OcrProvider>,
    pub stt: Arc<dyn SttProvider>,
}

impl ServiceRegistry {
    /// Create a registry with mock providers suitable for testing.
    pub fn mock() -> Self {
        Self {
            tts: Arc::new(MockTtsProvider::new()),
            ocr: Arc::new(MockOcrProvider),
            stt: Arc::new(MockSttProvider::new()),
        }
    }

    /// Create a registry with auto-detected real providers from environment variables.
    ///
    /// Falls back to mock providers when API keys are not configured.
    pub fn from_env() -> Self {
        Self {
            tts: Arc::from(create_tts_provider()),
            ocr: Arc::from(create_ocr_provider()),
            stt: Arc::from(create_stt_provider()),
        }
    }
}

impl Default for ServiceRegistry {
    fn default() -> Self {
        Self::mock()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_mock_registry_creation() {
        let registry = ServiceRegistry::mock();
        // Verify providers are cloneable (they will be used via trait methods in handlers)
        let _ = registry.tts.clone();
        let _ = registry.ocr.clone();
        let _ = registry.stt.clone();
    }

    #[test]
    fn test_default_is_mock() {
        let _registry = ServiceRegistry::default();
    }

    #[test]
    fn test_clone_shares_arcs() {
        let registry = ServiceRegistry::mock();
        let cloned = registry.clone();
        // Both point to the same Arc'd providers
        assert!(Arc::ptr_eq(&registry.tts, &cloned.tts));
        assert!(Arc::ptr_eq(&registry.ocr, &cloned.ocr));
        assert!(Arc::ptr_eq(&registry.stt, &cloned.stt));
    }
}
