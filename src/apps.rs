//! Application modules
//!
//! Contains domain-specific applications for the HaiLanGo platform:
//! - auth: User authentication and authorization
//! - books: Book management and OCR processing
//! - learning: Learning sessions and page interactions
//! - tts: Text-to-Speech services
//! - stt: Speech-to-Text and pronunciation evaluation
//! - review: Spaced Repetition System (SRS)
//! - teacher_mode: Automated lesson playback

pub mod auth;
pub mod books;
pub mod learning;
pub mod review;
pub mod stt;
pub mod teacher_mode;
pub mod tts;
