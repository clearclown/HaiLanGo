//! External service integrations
//!
//! Provides integrations with external services:
//! - OCR: Optical Character Recognition for book scanning
//! - TTS: Text-to-Speech for audio generation
//! - STT: Speech-to-Text and pronunciation evaluation
//! - Cache: Redis-based caching and session management
//! - WebSocket: Real-time communication for Teacher Mode
//! - i18n: Internationalization support

pub mod cache;
pub mod di;
pub mod i18n;
pub mod ocr;
pub mod stt;
pub mod tts;
pub mod websocket;

// OCR exports
pub use ocr::{
    GoogleVisionOcrProvider, MockOcrProvider, OcrError, OcrProvider, OcrResult, create_ocr_provider,
};

// TTS exports
pub use tts::{
    AudioFormat, ElevenLabsTtsProvider, GoogleCloudTtsProvider, MockTtsProvider, QualityTier,
    TtsError, TtsProvider, TtsRequest, TtsResponse, create_tts_provider, create_tts_provider_by_name,
};

// STT exports
pub use stt::{
    MockSttProvider, PronunciationEvaluator, PronunciationResult, SttAudioFormat, SttError,
    SttProvider, SttRequest, SttResponse, WhisperSttProvider, WordResult, WordScore,
    create_stt_provider,
};

// Cache exports
pub use cache::CacheService;

// WebSocket exports
pub use websocket::{LessonSession, PlaybackState, WsConnectionManager, WsMessage};

// i18n exports
pub use i18n::{Language, get_translations, keys, translate};
