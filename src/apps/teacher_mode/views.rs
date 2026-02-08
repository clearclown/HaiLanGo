//! Teacher Mode ViewSet - Lesson playback control

use uuid::Uuid;

use super::dto::{
    LessonActionResponse, LessonStatusResponse, SessionSummaryResponse, StartLessonRequest,
    UpdateConfigRequest,
};
use super::models::{TeacherSession, TeacherSessionStatus};

/// Result of a teacher mode action
#[derive(Debug)]
pub enum TeacherActionResult {
    Started(LessonActionResponse),
    Updated(LessonActionResponse),
    Status(LessonStatusResponse),
    Sessions(Vec<SessionSummaryResponse>),
    InvalidInput(String),
    NotFound(String),
    InvalidState(String),
}

/// Teacher Mode ViewSet
pub struct TeacherModeViewSet;

impl TeacherModeViewSet {
    /// Start a new teacher mode lesson
    pub fn start_lesson(
        user_id: Uuid,
        request: StartLessonRequest,
        sessions: &mut Vec<TeacherSession>,
    ) -> TeacherActionResult {
        // Validate page range
        if request.start_page == 0 {
            return TeacherActionResult::InvalidInput(
                "start_page must be >= 1".to_string(),
            );
        }
        if request.end_page < request.start_page {
            return TeacherActionResult::InvalidInput(
                "end_page must be >= start_page".to_string(),
            );
        }

        // Check for existing active session
        if sessions
            .iter()
            .any(|s| s.user_id == user_id && s.is_active())
        {
            return TeacherActionResult::InvalidState(
                "An active session already exists. Stop it first.".to_string(),
            );
        }

        let config = request.to_config();
        let mut session = TeacherSession::new(
            user_id,
            request.book_id,
            request.start_page,
            request.end_page,
            config,
        );

        session.play();
        session.start_current_page();

        let response = LessonActionResponse {
            id: session.id,
            status: session.status,
            current_page: session.current_page,
            message: format!(
                "Lesson started: pages {}-{}",
                session.start_page, session.end_page
            ),
        };

        sessions.push(session);
        TeacherActionResult::Started(response)
    }

    /// Pause the active session
    pub fn pause(
        user_id: Uuid,
        sessions: &mut [TeacherSession],
    ) -> TeacherActionResult {
        let session = sessions
            .iter_mut()
            .find(|s| s.user_id == user_id && s.status == TeacherSessionStatus::Playing);

        match session {
            Some(s) => {
                s.pause();
                TeacherActionResult::Updated(LessonActionResponse {
                    id: s.id,
                    status: s.status,
                    current_page: s.current_page,
                    message: "Lesson paused".to_string(),
                })
            }
            None => TeacherActionResult::NotFound(
                "No playing session found".to_string(),
            ),
        }
    }

    /// Resume the paused session
    pub fn resume(
        user_id: Uuid,
        sessions: &mut [TeacherSession],
    ) -> TeacherActionResult {
        let session = sessions
            .iter_mut()
            .find(|s| s.user_id == user_id && s.status == TeacherSessionStatus::Paused);

        match session {
            Some(s) => {
                s.play();
                TeacherActionResult::Updated(LessonActionResponse {
                    id: s.id,
                    status: s.status,
                    current_page: s.current_page,
                    message: "Lesson resumed".to_string(),
                })
            }
            None => TeacherActionResult::NotFound(
                "No paused session found".to_string(),
            ),
        }
    }

    /// Stop the active session
    pub fn stop(
        user_id: Uuid,
        sessions: &mut [TeacherSession],
    ) -> TeacherActionResult {
        let session = sessions
            .iter_mut()
            .find(|s| s.user_id == user_id && s.is_active());

        match session {
            Some(s) => {
                s.stop();
                TeacherActionResult::Updated(LessonActionResponse {
                    id: s.id,
                    status: s.status,
                    current_page: s.current_page,
                    message: "Lesson stopped".to_string(),
                })
            }
            None => TeacherActionResult::NotFound(
                "No active session found".to_string(),
            ),
        }
    }

    /// Advance to the next page manually
    pub fn next_page(
        user_id: Uuid,
        sessions: &mut [TeacherSession],
    ) -> TeacherActionResult {
        let session = sessions
            .iter_mut()
            .find(|s| s.user_id == user_id && s.is_active());

        match session {
            Some(s) => {
                // Complete current page forcefully
                if let Some(pb) = s.current_page_playback_mut() {
                    if !pb.completed {
                        pb.complete();
                        s.pages_completed += 1;
                    }
                }

                if s.advance_page() {
                    s.start_current_page();
                    TeacherActionResult::Updated(LessonActionResponse {
                        id: s.id,
                        status: s.status,
                        current_page: s.current_page,
                        message: format!("Advanced to page {}", s.current_page),
                    })
                } else {
                    TeacherActionResult::Updated(LessonActionResponse {
                        id: s.id,
                        status: s.status,
                        current_page: s.current_page,
                        message: "Lesson completed".to_string(),
                    })
                }
            }
            None => TeacherActionResult::NotFound(
                "No active session found".to_string(),
            ),
        }
    }

    /// Update playback config for the active session
    pub fn update_config(
        user_id: Uuid,
        request: UpdateConfigRequest,
        sessions: &mut [TeacherSession],
    ) -> TeacherActionResult {
        let session = sessions
            .iter_mut()
            .find(|s| s.user_id == user_id && s.is_active());

        match session {
            Some(s) => {
                let mut config = s.config.clone();
                if let Some(speed) = request.speed {
                    config.speed = speed;
                }
                if let Some(interval) = request.page_interval_secs {
                    config.page_interval_secs = interval;
                }
                if let Some(repeat) = request.repeat_count {
                    config.repeat_count = repeat;
                }
                if let Some(auto) = request.auto_advance {
                    config.auto_advance = auto;
                }
                if let Some(lang) = request.language {
                    config.language = lang;
                }

                s.update_config(config);

                TeacherActionResult::Updated(LessonActionResponse {
                    id: s.id,
                    status: s.status,
                    current_page: s.current_page,
                    message: "Configuration updated".to_string(),
                })
            }
            None => TeacherActionResult::NotFound(
                "No active session found".to_string(),
            ),
        }
    }

    /// Get the status of the active session
    pub fn get_status(
        user_id: Uuid,
        sessions: &[TeacherSession],
    ) -> TeacherActionResult {
        let session = sessions
            .iter()
            .rev()
            .find(|s| s.user_id == user_id);

        match session {
            Some(s) => TeacherActionResult::Status(LessonStatusResponse {
                id: s.id,
                book_id: s.book_id,
                status: s.status,
                current_page: s.current_page,
                start_page: s.start_page,
                end_page: s.end_page,
                total_pages: s.total_pages,
                pages_completed: s.pages_completed,
                progress_percent: s.progress_percent(),
                config: s.config.clone(),
            }),
            None => TeacherActionResult::NotFound(
                "No session found".to_string(),
            ),
        }
    }

    /// List session history for a user
    pub fn list_sessions(
        user_id: Uuid,
        sessions: &[TeacherSession],
    ) -> TeacherActionResult {
        let user_sessions: Vec<SessionSummaryResponse> = sessions
            .iter()
            .filter(|s| s.user_id == user_id)
            .map(|s| SessionSummaryResponse {
                id: s.id,
                book_id: s.book_id,
                status: s.status,
                start_page: s.start_page,
                end_page: s.end_page,
                pages_completed: s.pages_completed,
                total_pages: s.total_pages,
                progress_percent: s.progress_percent(),
                created_at: s.created_at,
            })
            .collect();

        TeacherActionResult::Sessions(user_sessions)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn make_start_request(book_id: Uuid, start: u32, end: u32) -> StartLessonRequest {
        StartLessonRequest {
            book_id,
            start_page: start,
            end_page: end,
            language: None,
            speed: None,
            page_interval_secs: None,
            repeat_count: None,
            auto_advance: None,
        }
    }

    #[test]
    fn test_start_lesson_success() {
        let user_id = Uuid::new_v4();
        let book_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        let result = TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(book_id, 1, 10),
            &mut sessions,
        );

        match result {
            TeacherActionResult::Started(resp) => {
                assert_eq!(resp.status, TeacherSessionStatus::Playing);
                assert_eq!(resp.current_page, 1);
                assert!(resp.message.contains("1-10"));
            }
            other => panic!("Expected Started, got {:?}", other),
        }
        assert_eq!(sessions.len(), 1);
    }

    #[test]
    fn test_start_lesson_invalid_page_range() {
        let user_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        let result = TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 10, 5),
            &mut sessions,
        );

        assert!(matches!(result, TeacherActionResult::InvalidInput(_)));
    }

    #[test]
    fn test_start_lesson_zero_start_page() {
        let user_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        let result = TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 0, 5),
            &mut sessions,
        );

        assert!(matches!(result, TeacherActionResult::InvalidInput(_)));
    }

    #[test]
    fn test_start_lesson_duplicate_active() {
        let user_id = Uuid::new_v4();
        let book_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(book_id, 1, 10),
            &mut sessions,
        );

        let result = TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(book_id, 1, 5),
            &mut sessions,
        );

        assert!(matches!(result, TeacherActionResult::InvalidState(_)));
    }

    #[test]
    fn test_pause_and_resume() {
        let user_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 1, 10),
            &mut sessions,
        );

        let result = TeacherModeViewSet::pause(user_id, &mut sessions);
        match result {
            TeacherActionResult::Updated(resp) => {
                assert_eq!(resp.status, TeacherSessionStatus::Paused);
            }
            other => panic!("Expected Updated, got {:?}", other),
        }

        let result = TeacherModeViewSet::resume(user_id, &mut sessions);
        match result {
            TeacherActionResult::Updated(resp) => {
                assert_eq!(resp.status, TeacherSessionStatus::Playing);
            }
            other => panic!("Expected Updated, got {:?}", other),
        }
    }

    #[test]
    fn test_pause_no_session() {
        let result = TeacherModeViewSet::pause(Uuid::new_v4(), &mut Vec::new());
        assert!(matches!(result, TeacherActionResult::NotFound(_)));
    }

    #[test]
    fn test_stop_session() {
        let user_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 1, 10),
            &mut sessions,
        );

        let result = TeacherModeViewSet::stop(user_id, &mut sessions);
        match result {
            TeacherActionResult::Updated(resp) => {
                assert_eq!(resp.status, TeacherSessionStatus::Stopped);
            }
            other => panic!("Expected Updated, got {:?}", other),
        }
    }

    #[test]
    fn test_next_page() {
        let user_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 1, 3),
            &mut sessions,
        );

        let result = TeacherModeViewSet::next_page(user_id, &mut sessions);
        match result {
            TeacherActionResult::Updated(resp) => {
                assert_eq!(resp.current_page, 2);
                assert!(resp.message.contains("page 2"));
            }
            other => panic!("Expected Updated, got {:?}", other),
        }
    }

    #[test]
    fn test_next_page_completes_lesson() {
        let user_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 1, 2),
            &mut sessions,
        );

        TeacherModeViewSet::next_page(user_id, &mut sessions);
        let result = TeacherModeViewSet::next_page(user_id, &mut sessions);

        match result {
            TeacherActionResult::Updated(resp) => {
                assert_eq!(resp.status, TeacherSessionStatus::Completed);
                assert!(resp.message.contains("completed"));
            }
            other => panic!("Expected Updated, got {:?}", other),
        }
    }

    #[test]
    fn test_update_config() {
        let user_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 1, 10),
            &mut sessions,
        );

        let update = UpdateConfigRequest {
            speed: Some(1.8),
            page_interval_secs: Some(15),
            repeat_count: None,
            auto_advance: Some(false),
            language: Some("ja".to_string()),
        };

        let result = TeacherModeViewSet::update_config(user_id, update, &mut sessions);
        assert!(matches!(result, TeacherActionResult::Updated(_)));

        // Verify the config was updated
        let session = sessions.first().unwrap();
        assert_eq!(session.config.speed, 1.8);
        assert_eq!(session.config.page_interval_secs, 15);
        assert!(!session.config.auto_advance);
        assert_eq!(session.config.language, "ja");
    }

    #[test]
    fn test_get_status() {
        let user_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 1, 10),
            &mut sessions,
        );

        let result = TeacherModeViewSet::get_status(user_id, &sessions);
        match result {
            TeacherActionResult::Status(resp) => {
                assert_eq!(resp.status, TeacherSessionStatus::Playing);
                assert_eq!(resp.current_page, 1);
                assert_eq!(resp.total_pages, 10);
                assert_eq!(resp.progress_percent, 0.0);
            }
            other => panic!("Expected Status, got {:?}", other),
        }
    }

    #[test]
    fn test_get_status_no_session() {
        let result = TeacherModeViewSet::get_status(Uuid::new_v4(), &[]);
        assert!(matches!(result, TeacherActionResult::NotFound(_)));
    }

    #[test]
    fn test_list_sessions() {
        let user_id = Uuid::new_v4();
        let other_user = Uuid::new_v4();
        let mut sessions = Vec::new();

        // Start and stop one session
        TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 1, 5),
            &mut sessions,
        );
        TeacherModeViewSet::stop(user_id, &mut sessions);

        // Start another session
        TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 1, 10),
            &mut sessions,
        );

        // Other user's session
        TeacherModeViewSet::start_lesson(
            other_user,
            make_start_request(Uuid::new_v4(), 1, 3),
            &mut sessions,
        );

        let result = TeacherModeViewSet::list_sessions(user_id, &sessions);
        match result {
            TeacherActionResult::Sessions(list) => {
                assert_eq!(list.len(), 2);
            }
            other => panic!("Expected Sessions, got {:?}", other),
        }
    }

    #[test]
    fn test_list_sessions_empty() {
        let result = TeacherModeViewSet::list_sessions(Uuid::new_v4(), &[]);
        match result {
            TeacherActionResult::Sessions(list) => {
                assert!(list.is_empty());
            }
            other => panic!("Expected Sessions, got {:?}", other),
        }
    }

    #[test]
    fn test_start_lesson_with_custom_config() {
        let user_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        let request = StartLessonRequest {
            book_id: Uuid::new_v4(),
            start_page: 5,
            end_page: 15,
            language: Some("ja".to_string()),
            speed: Some(0.8),
            page_interval_secs: Some(10),
            repeat_count: Some(2),
            auto_advance: Some(false),
        };

        TeacherModeViewSet::start_lesson(user_id, request, &mut sessions);

        let session = sessions.first().unwrap();
        assert_eq!(session.config.language, "ja");
        assert_eq!(session.config.speed, 0.8);
        assert_eq!(session.config.page_interval_secs, 10);
        assert_eq!(session.config.repeat_count, 2);
        assert!(!session.config.auto_advance);
        assert_eq!(session.start_page, 5);
        assert_eq!(session.end_page, 15);
    }

    #[test]
    fn test_can_start_after_stop() {
        let user_id = Uuid::new_v4();
        let mut sessions = Vec::new();

        TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 1, 5),
            &mut sessions,
        );
        TeacherModeViewSet::stop(user_id, &mut sessions);

        // Should be able to start a new session
        let result = TeacherModeViewSet::start_lesson(
            user_id,
            make_start_request(Uuid::new_v4(), 1, 10),
            &mut sessions,
        );
        assert!(matches!(result, TeacherActionResult::Started(_)));
    }
}
