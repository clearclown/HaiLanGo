//! Teacher Mode models - Lesson playback and session management

use chrono::{DateTime, Utc};
use reinhardt::db::orm::{FieldSelector, Model, Timestamped};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Type-safe field selector for TeacherSession model
#[derive(Clone)]
pub struct TeacherSessionFields;

impl FieldSelector for TeacherSessionFields {
    fn with_alias(self, _alias: &str) -> Self {
        self
    }
}

/// Playback status of a teacher mode session
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(rename_all = "snake_case")]
pub enum TeacherSessionStatus {
    #[default]
    Idle,
    Playing,
    Paused,
    Completed,
    Stopped,
}

/// Playback configuration for teacher mode
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PlaybackConfig {
    /// TTS reading speed (0.5 - 2.0)
    pub speed: f32,
    /// Interval between pages in seconds (0 - 30)
    pub page_interval_secs: u32,
    /// Number of times to repeat each page (1 - 3)
    pub repeat_count: u32,
    /// Automatically advance to next page after playback
    pub auto_advance: bool,
    /// Language for TTS
    pub language: String,
}

impl Default for PlaybackConfig {
    fn default() -> Self {
        Self {
            speed: 1.0,
            page_interval_secs: 5,
            repeat_count: 1,
            auto_advance: true,
            language: "en".to_string(),
        }
    }
}

impl PlaybackConfig {
    /// Validate and clamp configuration values
    pub fn validated(mut self) -> Self {
        self.speed = self.speed.clamp(0.5, 2.0);
        self.page_interval_secs = self.page_interval_secs.min(30);
        self.repeat_count = self.repeat_count.clamp(1, 3);
        self
    }
}

/// Track playback state for an individual page
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PagePlayback {
    pub page_number: u32,
    pub current_repeat: u32,
    pub completed: bool,
    pub started_at: Option<DateTime<Utc>>,
    pub completed_at: Option<DateTime<Utc>>,
}

impl PagePlayback {
    pub fn new(page_number: u32) -> Self {
        Self {
            page_number,
            current_repeat: 0,
            completed: false,
            started_at: None,
            completed_at: None,
        }
    }

    /// Start playing this page
    pub fn start(&mut self) {
        if self.started_at.is_none() {
            self.started_at = Some(Utc::now());
        }
        self.current_repeat += 1;
    }

    /// Mark page as completed
    pub fn complete(&mut self) {
        self.completed = true;
        self.completed_at = Some(Utc::now());
    }

    /// Check if all repeats are done
    pub fn repeats_done(&self, max_repeats: u32) -> bool {
        self.current_repeat >= max_repeats
    }
}

/// A teacher mode lesson session
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TeacherSession {
    pub id: Uuid,
    pub user_id: Uuid,
    pub book_id: Uuid,
    pub start_page: u32,
    pub end_page: u32,
    pub current_page: u32,
    pub status: TeacherSessionStatus,
    pub config: PlaybackConfig,
    pub pages_completed: u32,
    pub total_pages: u32,
    pub page_playbacks: Vec<PagePlayback>,
    pub created_at: DateTime<Utc>,
    pub started_at: Option<DateTime<Utc>>,
    pub ended_at: Option<DateTime<Utc>>,
}

impl TeacherSession {
    pub fn new(
        user_id: Uuid,
        book_id: Uuid,
        start_page: u32,
        end_page: u32,
        config: PlaybackConfig,
    ) -> Self {
        let total_pages = end_page.saturating_sub(start_page) + 1;
        let page_playbacks = (start_page..=end_page).map(PagePlayback::new).collect();

        Self {
            id: Uuid::new_v4(),
            user_id,
            book_id,
            start_page,
            end_page,
            current_page: start_page,
            status: TeacherSessionStatus::Idle,
            config,
            pages_completed: 0,
            total_pages,
            page_playbacks,
            created_at: Utc::now(),
            started_at: None,
            ended_at: None,
        }
    }

    /// Start or resume playback
    pub fn play(&mut self) {
        match self.status {
            TeacherSessionStatus::Idle => {
                self.status = TeacherSessionStatus::Playing;
                self.started_at = Some(Utc::now());
            }
            TeacherSessionStatus::Paused => {
                self.status = TeacherSessionStatus::Playing;
            }
            _ => {}
        }
    }

    /// Pause playback
    pub fn pause(&mut self) {
        if self.status == TeacherSessionStatus::Playing {
            self.status = TeacherSessionStatus::Paused;
        }
    }

    /// Stop the session
    pub fn stop(&mut self) {
        if matches!(
            self.status,
            TeacherSessionStatus::Playing | TeacherSessionStatus::Paused
        ) {
            self.status = TeacherSessionStatus::Stopped;
            self.ended_at = Some(Utc::now());
        }
    }

    /// Get the current page's playback state (mutable)
    pub fn current_page_playback_mut(&mut self) -> Option<&mut PagePlayback> {
        let idx = (self.current_page - self.start_page) as usize;
        self.page_playbacks.get_mut(idx)
    }

    /// Get the current page's playback state
    pub fn current_page_playback(&self) -> Option<&PagePlayback> {
        let idx = (self.current_page - self.start_page) as usize;
        self.page_playbacks.get(idx)
    }

    /// Mark the current page as started (increment repeat)
    pub fn start_current_page(&mut self) {
        if let Some(pb) = self.current_page_playback_mut() {
            pb.start();
        }
    }

    /// Complete the current page and optionally advance
    pub fn complete_current_page(&mut self) -> bool {
        let repeat_count = self.config.repeat_count;
        let auto_advance = self.config.auto_advance;

        if let Some(pb) = self.current_page_playback_mut() {
            if pb.repeats_done(repeat_count) {
                pb.complete();
                self.pages_completed += 1;

                if auto_advance {
                    return self.advance_page();
                }
            }
        }
        false
    }

    /// Advance to the next page; returns true if advanced, false if at end
    pub fn advance_page(&mut self) -> bool {
        if self.current_page < self.end_page {
            self.current_page += 1;
            true
        } else {
            self.status = TeacherSessionStatus::Completed;
            self.ended_at = Some(Utc::now());
            false
        }
    }

    /// Update playback configuration mid-session
    pub fn update_config(&mut self, config: PlaybackConfig) {
        self.config = config.validated();
    }

    /// Calculate progress percentage (0.0 - 100.0)
    pub fn progress_percent(&self) -> f32 {
        if self.total_pages == 0 {
            return 0.0;
        }
        (self.pages_completed as f32 / self.total_pages as f32) * 100.0
    }

    /// Check if the session is active (playing or paused)
    pub fn is_active(&self) -> bool {
        matches!(
            self.status,
            TeacherSessionStatus::Playing | TeacherSessionStatus::Paused
        )
    }
}

impl Model for TeacherSession {
    type PrimaryKey = Uuid;
    type Fields = TeacherSessionFields;

    fn table_name() -> &'static str {
        "teacher_sessions"
    }

    fn app_label() -> &'static str {
        "teacher_mode"
    }

    fn new_fields() -> Self::Fields {
        TeacherSessionFields
    }

    fn primary_key(&self) -> Option<Self::PrimaryKey> {
        Some(self.id)
    }

    fn set_primary_key(&mut self, value: Self::PrimaryKey) {
        self.id = value;
    }
}

impl Timestamped for TeacherSession {
    fn created_at(&self) -> DateTime<Utc> {
        self.created_at
    }

    fn updated_at(&self) -> DateTime<Utc> {
        self.ended_at.unwrap_or(self.created_at)
    }

    fn set_updated_at(&mut self, _time: DateTime<Utc>) {
        // Session timestamps are managed by state transitions
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn default_config() -> PlaybackConfig {
        PlaybackConfig::default()
    }

    #[test]
    fn test_playback_config_default() {
        let config = PlaybackConfig::default();
        assert_eq!(config.speed, 1.0);
        assert_eq!(config.page_interval_secs, 5);
        assert_eq!(config.repeat_count, 1);
        assert!(config.auto_advance);
        assert_eq!(config.language, "en");
    }

    #[test]
    fn test_playback_config_validation() {
        let config = PlaybackConfig {
            speed: 5.0,
            page_interval_secs: 60,
            repeat_count: 10,
            auto_advance: true,
            language: "ja".to_string(),
        }
        .validated();

        assert_eq!(config.speed, 2.0);
        assert_eq!(config.page_interval_secs, 30);
        assert_eq!(config.repeat_count, 3);
    }

    #[test]
    fn test_playback_config_validation_lower_bounds() {
        let config = PlaybackConfig {
            speed: 0.1,
            page_interval_secs: 0,
            repeat_count: 0,
            auto_advance: false,
            language: "en".to_string(),
        }
        .validated();

        assert_eq!(config.speed, 0.5);
        assert_eq!(config.page_interval_secs, 0);
        assert_eq!(config.repeat_count, 1);
    }

    #[test]
    fn test_teacher_session_creation() {
        let user_id = Uuid::new_v4();
        let book_id = Uuid::new_v4();
        let session = TeacherSession::new(user_id, book_id, 1, 10, default_config());

        assert_eq!(session.user_id, user_id);
        assert_eq!(session.book_id, book_id);
        assert_eq!(session.start_page, 1);
        assert_eq!(session.end_page, 10);
        assert_eq!(session.current_page, 1);
        assert_eq!(session.total_pages, 10);
        assert_eq!(session.pages_completed, 0);
        assert_eq!(session.status, TeacherSessionStatus::Idle);
        assert_eq!(session.page_playbacks.len(), 10);
    }

    #[test]
    fn test_teacher_session_play() {
        let mut session =
            TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 1, 5, default_config());

        session.play();
        assert_eq!(session.status, TeacherSessionStatus::Playing);
        assert!(session.started_at.is_some());
    }

    #[test]
    fn test_teacher_session_pause_resume() {
        let mut session =
            TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 1, 5, default_config());

        session.play();
        session.pause();
        assert_eq!(session.status, TeacherSessionStatus::Paused);

        session.play();
        assert_eq!(session.status, TeacherSessionStatus::Playing);
    }

    #[test]
    fn test_teacher_session_stop() {
        let mut session =
            TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 1, 5, default_config());

        session.play();
        session.stop();
        assert_eq!(session.status, TeacherSessionStatus::Stopped);
        assert!(session.ended_at.is_some());
    }

    #[test]
    fn test_teacher_session_stop_from_paused() {
        let mut session =
            TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 1, 5, default_config());

        session.play();
        session.pause();
        session.stop();
        assert_eq!(session.status, TeacherSessionStatus::Stopped);
    }

    #[test]
    fn test_teacher_session_cannot_stop_idle() {
        let mut session =
            TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 1, 5, default_config());

        session.stop();
        assert_eq!(session.status, TeacherSessionStatus::Idle);
    }

    #[test]
    fn test_advance_page() {
        let mut session =
            TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 1, 3, default_config());

        assert!(session.advance_page());
        assert_eq!(session.current_page, 2);

        assert!(session.advance_page());
        assert_eq!(session.current_page, 3);

        assert!(!session.advance_page());
        assert_eq!(session.status, TeacherSessionStatus::Completed);
    }

    #[test]
    fn test_page_playback_repeats() {
        let mut pb = PagePlayback::new(1);

        pb.start();
        assert_eq!(pb.current_repeat, 1);
        assert!(!pb.repeats_done(2));

        pb.start();
        assert_eq!(pb.current_repeat, 2);
        assert!(pb.repeats_done(2));
    }

    #[test]
    fn test_complete_current_page_with_repeat() {
        let config = PlaybackConfig {
            repeat_count: 2,
            auto_advance: true,
            ..PlaybackConfig::default()
        };
        let mut session = TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 1, 3, config);

        session.play();
        session.start_current_page(); // repeat 1
        assert!(!session.complete_current_page()); // not enough repeats

        session.start_current_page(); // repeat 2
        assert!(session.complete_current_page()); // advances to page 2
        assert_eq!(session.current_page, 2);
        assert_eq!(session.pages_completed, 1);
    }

    #[test]
    fn test_complete_current_page_no_auto_advance() {
        let config = PlaybackConfig {
            repeat_count: 1,
            auto_advance: false,
            ..PlaybackConfig::default()
        };
        let mut session = TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 1, 3, config);

        session.play();
        session.start_current_page();
        session.complete_current_page();
        assert_eq!(session.current_page, 1); // stays on same page
        assert_eq!(session.pages_completed, 1);
    }

    #[test]
    fn test_progress_percent() {
        let mut session =
            TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 1, 4, default_config());

        assert_eq!(session.progress_percent(), 0.0);

        session.pages_completed = 2;
        assert_eq!(session.progress_percent(), 50.0);

        session.pages_completed = 4;
        assert_eq!(session.progress_percent(), 100.0);
    }

    #[test]
    fn test_update_config() {
        let mut session =
            TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 1, 5, default_config());

        let new_config = PlaybackConfig {
            speed: 1.5,
            page_interval_secs: 10,
            repeat_count: 2,
            auto_advance: false,
            language: "ja".to_string(),
        };

        session.update_config(new_config);
        assert_eq!(session.config.speed, 1.5);
        assert_eq!(session.config.page_interval_secs, 10);
        assert_eq!(session.config.repeat_count, 2);
        assert!(!session.config.auto_advance);
    }

    #[test]
    fn test_is_active() {
        let mut session =
            TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 1, 5, default_config());

        assert!(!session.is_active());

        session.play();
        assert!(session.is_active());

        session.pause();
        assert!(session.is_active());

        session.stop();
        assert!(!session.is_active());
    }

    #[test]
    fn test_session_status_serialization() {
        let json = serde_json::to_string(&TeacherSessionStatus::Playing).unwrap();
        assert_eq!(json, "\"playing\"");

        let de: TeacherSessionStatus = serde_json::from_str("\"paused\"").unwrap();
        assert_eq!(de, TeacherSessionStatus::Paused);
    }

    #[test]
    fn test_single_page_session() {
        let mut session =
            TeacherSession::new(Uuid::new_v4(), Uuid::new_v4(), 5, 5, default_config());

        assert_eq!(session.total_pages, 1);
        assert_eq!(session.page_playbacks.len(), 1);

        session.play();
        session.start_current_page();
        session.complete_current_page();
        assert_eq!(session.status, TeacherSessionStatus::Completed);
    }
}
