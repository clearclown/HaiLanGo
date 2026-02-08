//! Speech-to-Text and Pronunciation app
//!
//! Handles speech recognition and pronunciation evaluation.
//! Records user speech and provides feedback on pronunciation accuracy.
//!
//! Models:
//! - PronunciationAttempt: User's speech recording and evaluation
//! - WordFeedback: Per-word feedback stored alongside an attempt
//! - PronunciationStats: Aggregated pronunciation statistics
//!
//! Services (in `crate::services::stt`):
//! - SttProvider trait: Abstraction for STT backends (Whisper, Azure, Mock)
//! - PronunciationEvaluator: Compares recognized text against expected text
//!
//! Views:
//! - SttViewSet: Evaluate pronunciation, list attempts, get stats

pub mod dto;
pub mod models;
pub mod views;

pub use dto::*;
pub use models::{AttemptStatus, PronunciationAttempt, PronunciationStats, WeakWord, WordFeedback};
pub use views::{EvaluateResult, GetAttemptResult, SttViewSet};
