//! Text-to-Speech app
//!
//! Provides audio generation for learning content.
//! Supports multiple languages, voice options, and quality tiers.
//!
//! Models:
//! - AudioGeneration: Record of a TTS generation request
//! - AudioCache: Cached TTS results for reuse
//!
//! Services (in `crate::services::tts`):
//! - TtsProvider trait: Abstraction for TTS backends (Google, ElevenLabs, Mock)
//! - QualityTier: Standard vs Premium voice quality
//!
//! Views:
//! - TtsViewSet: Synthesize speech, list history, get supported languages

pub mod dto;
pub mod models;
pub mod views;

pub use dto::*;
pub use models::{AudioCache, AudioGeneration, GenerationStatus};
pub use views::{SynthesizeResult, TtsViewSet};
