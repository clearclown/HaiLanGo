//! Teacher Mode app
//!
//! Provides automated lesson playback and pacing control.
//! Allows users to progress through book pages with TTS narration,
//! configurable speed, page intervals, and repeat counts.
//!
//! Models:
//! - TeacherSession: A lesson playback session with config and progress
//! - PlaybackConfig: Speed, interval, repeat, auto-advance settings
//! - PagePlayback: Per-page playback tracking
//!
//! Views:
//! - TeacherModeViewSet: Start/pause/resume/stop/next/config/status

pub mod dto;
pub mod models;
pub mod views;

pub use dto::*;
pub use models::*;
pub use views::*;
