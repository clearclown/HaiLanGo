//! Learning ViewSet - Session management endpoints

use uuid::Uuid;

use super::dto::{
    CreateReviewSessionRequest, CreateSessionRequest, ProgressResponse, SessionAction,
    SessionResponse, UpdateProgressRequest, UpdateSessionStatusRequest,
};
use super::models::{LearningProgress, LearningSession, SessionStatus};

/// Result for session creation
#[derive(Debug)]
pub enum CreateSessionResult {
    Success(SessionResponse),
    BookNotFound,
    InvalidPageRange(String),
}

/// Result for session update
#[derive(Debug)]
pub enum UpdateSessionResult {
    Success(SessionResponse),
    SessionNotFound,
    InvalidAction(String),
}

/// Result for progress update
#[derive(Debug)]
pub enum UpdateProgressResult {
    Success(ProgressResponse),
    SessionNotFound,
    PageNotFound,
    InvalidInput(String),
}

/// Result for session list
#[derive(Debug)]
pub enum ListSessionResult {
    Success(Vec<SessionResponse>),
    Unauthorized,
}

/// Learning ViewSet - handles learning session endpoints
pub struct LearningViewSet;

impl LearningViewSet {
    /// Create a new learning session
    pub fn create(
        request: CreateSessionRequest,
        user_id: Uuid,
        book_exists: bool,
    ) -> CreateSessionResult {
        if !book_exists {
            return CreateSessionResult::BookNotFound;
        }

        // Validate page range
        if let (Some(start), Some(end)) = (request.start_page, request.end_page) {
            if start > end || start < 1 {
                return CreateSessionResult::InvalidPageRange(
                    "Start page must be >= 1 and <= end page".to_string(),
                );
            }
        }

        // Create session
        let mut session = LearningSession::new(user_id, request.book_id, request.session_type);

        // Apply page range if provided
        if let (Some(start), Some(end)) = (request.start_page, request.end_page) {
            session = session.with_pages(start, end);
        }

        // Apply settings if provided
        if let Some(settings) = request.settings {
            session.settings = settings;
        }

        CreateSessionResult::Success(Self::session_to_response(&session))
    }

    /// Create a review session (no book required)
    pub fn create_review(
        request: CreateReviewSessionRequest,
        user_id: Uuid,
    ) -> CreateSessionResult {
        let mut session = LearningSession::new_review(user_id);

        if let Some(settings) = request.settings {
            session.settings = settings;
        }

        CreateSessionResult::Success(Self::session_to_response(&session))
    }

    /// Update session status (pause/resume/complete/abandon)
    pub fn update_status(
        request: UpdateSessionStatusRequest,
        session: Option<&mut LearningSession>,
        user_id: Uuid,
    ) -> UpdateSessionResult {
        let session = match session {
            Some(s) => s,
            None => return UpdateSessionResult::SessionNotFound,
        };

        // Verify ownership
        if session.user_id != user_id {
            return UpdateSessionResult::SessionNotFound;
        }

        match request.action {
            SessionAction::Pause => {
                if session.status != SessionStatus::Active {
                    return UpdateSessionResult::InvalidAction(
                        "Can only pause active sessions".to_string(),
                    );
                }
                session.pause();
            }
            SessionAction::Resume => {
                if session.status != SessionStatus::Paused {
                    return UpdateSessionResult::InvalidAction(
                        "Can only resume paused sessions".to_string(),
                    );
                }
                session.resume();
            }
            SessionAction::Complete => {
                if session.status == SessionStatus::Completed
                    || session.status == SessionStatus::Abandoned
                {
                    return UpdateSessionResult::InvalidAction("Session already ended".to_string());
                }
                session.complete();
            }
            SessionAction::Abandon => {
                if session.status == SessionStatus::Completed
                    || session.status == SessionStatus::Abandoned
                {
                    return UpdateSessionResult::InvalidAction("Session already ended".to_string());
                }
                session.abandon();
            }
            SessionAction::NextPage => {
                if session.status != SessionStatus::Active {
                    return UpdateSessionResult::InvalidAction(
                        "Can only advance pages in active sessions".to_string(),
                    );
                }
                if !session.next_page() {
                    return UpdateSessionResult::InvalidAction(
                        "Already at the last page".to_string(),
                    );
                }
            }
        }

        UpdateSessionResult::Success(Self::session_to_response(session))
    }

    /// Record progress for a page
    pub fn record_progress(
        request: UpdateProgressRequest,
        session: Option<&LearningSession>,
        user_id: Uuid,
        page_exists: bool,
    ) -> UpdateProgressResult {
        let session = match session {
            Some(s) => s,
            None => return UpdateProgressResult::SessionNotFound,
        };

        // Verify ownership
        if session.user_id != user_id {
            return UpdateProgressResult::SessionNotFound;
        }

        if !page_exists {
            return UpdateProgressResult::PageNotFound;
        }

        // Validate scores
        if let Some(score) = request.pronunciation_score {
            if !(0..=100).contains(&score) {
                return UpdateProgressResult::InvalidInput(
                    "Pronunciation score must be 0-100".to_string(),
                );
            }
        }

        if let Some(score) = request.comprehension_score {
            if !(0..=100).contains(&score) {
                return UpdateProgressResult::InvalidInput(
                    "Comprehension score must be 0-100".to_string(),
                );
            }
        }

        // Create progress record
        let mut progress = LearningProgress::new(session.id, request.page_id, user_id);
        progress.add_time(request.time_spent_seconds);

        if let Some(score) = request.pronunciation_score {
            progress.set_pronunciation_score(score);
        }

        progress.comprehension_score = request.comprehension_score;

        UpdateProgressResult::Success(Self::progress_to_response(&progress))
    }

    /// Get a session by ID
    pub fn retrieve(session: Option<&LearningSession>, user_id: Uuid) -> Option<SessionResponse> {
        session
            .filter(|s| s.user_id == user_id)
            .map(Self::session_to_response)
    }

    /// List sessions for a user
    pub fn list(sessions: &[LearningSession], user_id: Uuid) -> ListSessionResult {
        let user_sessions: Vec<_> = sessions
            .iter()
            .filter(|s| s.user_id == user_id)
            .map(Self::session_to_response)
            .collect();

        ListSessionResult::Success(user_sessions)
    }

    /// Convert session to response DTO
    fn session_to_response(session: &LearningSession) -> SessionResponse {
        SessionResponse {
            id: session.id,
            user_id: session.user_id,
            book_id: session.book_id,
            session_type: session.session_type,
            current_page: session.current_page,
            start_page: session.start_page,
            end_page: session.end_page,
            duration_seconds: session.duration_seconds,
            status: session.status,
            started_at: session.started_at,
            ended_at: session.ended_at,
        }
    }

    /// Convert progress to response DTO
    fn progress_to_response(progress: &LearningProgress) -> ProgressResponse {
        ProgressResponse {
            id: progress.id,
            session_id: progress.session_id,
            page_id: progress.page_id,
            time_spent_seconds: progress.time_spent_seconds,
            pronunciation_score: progress.pronunciation_score,
            comprehension_score: progress.comprehension_score,
            average_score: progress.average_score(),
            created_at: progress.created_at,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::apps::learning::models::{SessionSettings, SessionType};

    #[test]
    fn test_create_session_success() {
        let user_id = Uuid::new_v4();
        let request = CreateSessionRequest {
            book_id: Uuid::new_v4(),
            session_type: SessionType::PageByPage,
            start_page: None,
            end_page: None,
            settings: None,
        };

        let result = LearningViewSet::create(request, user_id, true);

        match result {
            CreateSessionResult::Success(response) => {
                assert_eq!(response.user_id, user_id);
                assert_eq!(response.session_type, SessionType::PageByPage);
                assert_eq!(response.status, SessionStatus::Active);
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_create_session_with_pages() {
        let user_id = Uuid::new_v4();
        let request = CreateSessionRequest {
            book_id: Uuid::new_v4(),
            session_type: SessionType::TeacherMode,
            start_page: Some(5),
            end_page: Some(15),
            settings: None,
        };

        let result = LearningViewSet::create(request, user_id, true);

        match result {
            CreateSessionResult::Success(response) => {
                assert_eq!(response.start_page, Some(5));
                assert_eq!(response.end_page, Some(15));
                assert_eq!(response.current_page, 5);
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_create_session_book_not_found() {
        let request = CreateSessionRequest {
            book_id: Uuid::new_v4(),
            session_type: SessionType::PageByPage,
            start_page: None,
            end_page: None,
            settings: None,
        };

        let result = LearningViewSet::create(request, Uuid::new_v4(), false);
        assert!(matches!(result, CreateSessionResult::BookNotFound));
    }

    #[test]
    fn test_create_session_invalid_page_range() {
        let request = CreateSessionRequest {
            book_id: Uuid::new_v4(),
            session_type: SessionType::PageByPage,
            start_page: Some(10),
            end_page: Some(5), // end < start
            settings: None,
        };

        let result = LearningViewSet::create(request, Uuid::new_v4(), true);
        assert!(matches!(result, CreateSessionResult::InvalidPageRange(_)));
    }

    #[test]
    fn test_create_review_session() {
        let user_id = Uuid::new_v4();
        let request = CreateReviewSessionRequest { settings: None };

        let result = LearningViewSet::create_review(request, user_id);

        match result {
            CreateSessionResult::Success(response) => {
                assert_eq!(response.session_type, SessionType::Review);
                assert!(response.book_id.is_none());
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_update_status_pause() {
        let user_id = Uuid::new_v4();
        let mut session = LearningSession::new(user_id, Uuid::new_v4(), SessionType::PageByPage);

        let request = UpdateSessionStatusRequest {
            action: SessionAction::Pause,
        };
        let result = LearningViewSet::update_status(request, Some(&mut session), user_id);

        match result {
            UpdateSessionResult::Success(response) => {
                assert_eq!(response.status, SessionStatus::Paused);
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_update_status_resume() {
        let user_id = Uuid::new_v4();
        let mut session = LearningSession::new(user_id, Uuid::new_v4(), SessionType::PageByPage);
        session.pause();

        let request = UpdateSessionStatusRequest {
            action: SessionAction::Resume,
        };
        let result = LearningViewSet::update_status(request, Some(&mut session), user_id);

        match result {
            UpdateSessionResult::Success(response) => {
                assert_eq!(response.status, SessionStatus::Active);
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_update_status_complete() {
        let user_id = Uuid::new_v4();
        let mut session = LearningSession::new(user_id, Uuid::new_v4(), SessionType::PageByPage);

        let request = UpdateSessionStatusRequest {
            action: SessionAction::Complete,
        };
        let result = LearningViewSet::update_status(request, Some(&mut session), user_id);

        match result {
            UpdateSessionResult::Success(response) => {
                assert_eq!(response.status, SessionStatus::Completed);
                assert!(response.ended_at.is_some());
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_update_status_invalid_action() {
        let user_id = Uuid::new_v4();
        let mut session = LearningSession::new(user_id, Uuid::new_v4(), SessionType::PageByPage);
        session.pause();

        // Try to pause an already paused session
        let request = UpdateSessionStatusRequest {
            action: SessionAction::Pause,
        };
        let result = LearningViewSet::update_status(request, Some(&mut session), user_id);

        assert!(matches!(result, UpdateSessionResult::InvalidAction(_)));
    }

    #[test]
    fn test_update_status_next_page() {
        let user_id = Uuid::new_v4();
        let mut session =
            LearningSession::new(user_id, Uuid::new_v4(), SessionType::PageByPage).with_pages(1, 5);

        let request = UpdateSessionStatusRequest {
            action: SessionAction::NextPage,
        };
        let result = LearningViewSet::update_status(request, Some(&mut session), user_id);

        match result {
            UpdateSessionResult::Success(response) => {
                assert_eq!(response.current_page, 2);
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_record_progress_success() {
        let user_id = Uuid::new_v4();
        let session = LearningSession::new(user_id, Uuid::new_v4(), SessionType::PageByPage);

        let request = UpdateProgressRequest {
            page_id: Uuid::new_v4(),
            time_spent_seconds: 120,
            pronunciation_score: Some(85),
            comprehension_score: Some(90),
        };

        let result = LearningViewSet::record_progress(request, Some(&session), user_id, true);

        match result {
            UpdateProgressResult::Success(response) => {
                assert_eq!(response.time_spent_seconds, 120);
                assert_eq!(response.pronunciation_score, Some(85));
                assert_eq!(response.average_score, Some(87.5));
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_record_progress_invalid_score() {
        let user_id = Uuid::new_v4();
        let session = LearningSession::new(user_id, Uuid::new_v4(), SessionType::PageByPage);

        let request = UpdateProgressRequest {
            page_id: Uuid::new_v4(),
            time_spent_seconds: 60,
            pronunciation_score: Some(150), // Invalid
            comprehension_score: None,
        };

        let result = LearningViewSet::record_progress(request, Some(&session), user_id, true);
        assert!(matches!(result, UpdateProgressResult::InvalidInput(_)));
    }

    #[test]
    fn test_retrieve_session() {
        let user_id = Uuid::new_v4();
        let session = LearningSession::new(user_id, Uuid::new_v4(), SessionType::PageByPage);

        let result = LearningViewSet::retrieve(Some(&session), user_id);
        assert!(result.is_some());

        // Wrong user
        let result = LearningViewSet::retrieve(Some(&session), Uuid::new_v4());
        assert!(result.is_none());
    }

    #[test]
    fn test_list_sessions() {
        let user_id = Uuid::new_v4();
        let other_user_id = Uuid::new_v4();

        let sessions = vec![
            LearningSession::new(user_id, Uuid::new_v4(), SessionType::PageByPage),
            LearningSession::new(user_id, Uuid::new_v4(), SessionType::TeacherMode),
            LearningSession::new(other_user_id, Uuid::new_v4(), SessionType::Review),
        ];

        let result = LearningViewSet::list(&sessions, user_id);

        match result {
            ListSessionResult::Success(user_sessions) => {
                assert_eq!(user_sessions.len(), 2);
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_create_session_with_custom_settings() {
        let user_id = Uuid::new_v4();
        let settings = SessionSettings {
            tts_speed: 0.8,
            page_interval: 10,
            repeat_count: 2,
            include_translation: true,
            include_vocabulary: true,
            include_grammar: true,
        };

        let request = CreateSessionRequest {
            book_id: Uuid::new_v4(),
            session_type: SessionType::TeacherMode,
            start_page: None,
            end_page: None,
            settings: Some(settings),
        };

        let result = LearningViewSet::create(request, user_id, true);
        assert!(matches!(result, CreateSessionResult::Success(_)));
    }
}
