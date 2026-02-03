//! Learning app - Learning sessions and progress tracking
//!
//! Provides session management for different learning modes:
//! - PageByPage: Study one page at a time
//! - TeacherMode: Automated lesson playback with TTS
//! - Review: SRS vocabulary review session

pub mod dto;
pub mod models;
pub mod views;

pub use dto::{
    CreateReviewSessionRequest, CreateSessionRequest, ProgressResponse, SessionAction,
    SessionListResponse, SessionResponse, SessionStatsResponse, UpdateProgressRequest,
    UpdateSessionStatusRequest,
};
pub use models::{LearningProgress, LearningSession, SessionSettings, SessionStatus, SessionType};
pub use views::{
    CreateSessionResult, LearningViewSet, ListSessionResult, UpdateProgressResult,
    UpdateSessionResult,
};
