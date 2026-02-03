//! WebSocket Service for Real-time Communication
//!
//! Provides WebSocket support for Teacher Mode - automated lesson playback
//! with real-time audio streaming and progress tracking.

use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::{RwLock, broadcast};
use uuid::Uuid;

/// WebSocket message types
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(tag = "type", content = "data")]
pub enum WsMessage {
    /// Client connected
    Connected {
        session_id: String,
    },
    /// Start lesson playback
    StartLesson {
        book_id: Uuid,
        page: u32,
    },
    /// Pause lesson
    PausLesson,
    /// Resume lesson
    ResumeLesson,
    /// Stop lesson
    StopLesson,
    /// Page content update
    PageContent {
        page: u32,
        text: String,
        audio_url: Option<String>,
    },
    /// TTS audio chunk (base64 encoded)
    AudioChunk {
        chunk: String,
        is_final: bool,
    },
    /// Progress update
    Progress {
        page: u32,
        total_pages: u32,
        time_elapsed_secs: u32,
    },
    /// Error message
    Error {
        code: String,
        message: String,
    },
    /// Ping/Pong for keep-alive
    Ping,
    Pong,
}

/// Lesson playback state
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum PlaybackState {
    Idle,
    Playing,
    Paused,
    Stopped,
}

/// Active lesson session
#[derive(Debug, Clone)]
pub struct LessonSession {
    pub id: Uuid,
    pub user_id: Uuid,
    pub book_id: Uuid,
    pub current_page: u32,
    pub total_pages: u32,
    pub state: PlaybackState,
    pub started_at: chrono::DateTime<chrono::Utc>,
    pub time_elapsed_secs: u32,
}

impl LessonSession {
    pub fn new(user_id: Uuid, book_id: Uuid, total_pages: u32) -> Self {
        Self {
            id: Uuid::new_v4(),
            user_id,
            book_id,
            current_page: 1,
            total_pages,
            state: PlaybackState::Idle,
            started_at: chrono::Utc::now(),
            time_elapsed_secs: 0,
        }
    }
}

/// WebSocket connection manager
pub struct WsConnectionManager {
    /// Active sessions by user ID
    sessions: Arc<RwLock<HashMap<Uuid, LessonSession>>>,
    /// Broadcast channel for messages
    broadcast_tx: broadcast::Sender<(Uuid, WsMessage)>,
}

impl WsConnectionManager {
    pub fn new() -> Self {
        let (broadcast_tx, _) = broadcast::channel(1024);
        Self {
            sessions: Arc::new(RwLock::new(HashMap::new())),
            broadcast_tx,
        }
    }

    /// Get a broadcast receiver for a user
    pub fn subscribe(&self) -> broadcast::Receiver<(Uuid, WsMessage)> {
        self.broadcast_tx.subscribe()
    }

    /// Start a new lesson session
    pub async fn start_lesson(
        &self,
        user_id: Uuid,
        book_id: Uuid,
        total_pages: u32,
    ) -> LessonSession {
        let mut sessions = self.sessions.write().await;
        let session = LessonSession::new(user_id, book_id, total_pages);
        sessions.insert(user_id, session.clone());

        // Broadcast start message
        let _ = self
            .broadcast_tx
            .send((user_id, WsMessage::StartLesson { book_id, page: 1 }));

        session
    }

    /// Get active session for a user
    pub async fn get_session(&self, user_id: Uuid) -> Option<LessonSession> {
        let sessions = self.sessions.read().await;
        sessions.get(&user_id).cloned()
    }

    /// Update session state
    pub async fn update_state(&self, user_id: Uuid, state: PlaybackState) {
        let mut sessions = self.sessions.write().await;
        if let Some(session) = sessions.get_mut(&user_id) {
            session.state = state;
        }
    }

    /// Advance to next page
    pub async fn next_page(&self, user_id: Uuid) -> Option<u32> {
        let mut sessions = self.sessions.write().await;
        if let Some(session) = sessions.get_mut(&user_id) {
            if session.current_page < session.total_pages {
                session.current_page += 1;

                // Broadcast progress
                let _ = self.broadcast_tx.send((
                    user_id,
                    WsMessage::Progress {
                        page: session.current_page,
                        total_pages: session.total_pages,
                        time_elapsed_secs: session.time_elapsed_secs,
                    },
                ));

                return Some(session.current_page);
            }
        }
        None
    }

    /// Send page content to user
    pub async fn send_page_content(
        &self,
        user_id: Uuid,
        page: u32,
        text: String,
        audio_url: Option<String>,
    ) {
        let _ = self.broadcast_tx.send((
            user_id,
            WsMessage::PageContent {
                page,
                text,
                audio_url,
            },
        ));
    }

    /// Send audio chunk for streaming TTS
    pub async fn send_audio_chunk(&self, user_id: Uuid, chunk: String, is_final: bool) {
        let _ = self
            .broadcast_tx
            .send((user_id, WsMessage::AudioChunk { chunk, is_final }));
    }

    /// End lesson session
    pub async fn end_session(&self, user_id: Uuid) {
        let mut sessions = self.sessions.write().await;
        sessions.remove(&user_id);

        let _ = self.broadcast_tx.send((user_id, WsMessage::StopLesson));
    }

    /// Get all active sessions
    pub async fn active_sessions(&self) -> Vec<LessonSession> {
        let sessions = self.sessions.read().await;
        sessions.values().cloned().collect()
    }
}

impl Default for WsConnectionManager {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_create_session() {
        let manager = WsConnectionManager::new();
        let user_id = Uuid::new_v4();
        let book_id = Uuid::new_v4();

        let session = manager.start_lesson(user_id, book_id, 10).await;

        assert_eq!(session.user_id, user_id);
        assert_eq!(session.book_id, book_id);
        assert_eq!(session.total_pages, 10);
        assert_eq!(session.current_page, 1);
        assert_eq!(session.state, PlaybackState::Idle);
    }

    #[tokio::test]
    async fn test_next_page() {
        let manager = WsConnectionManager::new();
        let user_id = Uuid::new_v4();
        let book_id = Uuid::new_v4();

        manager.start_lesson(user_id, book_id, 3).await;

        let page = manager.next_page(user_id).await;
        assert_eq!(page, Some(2));

        let page = manager.next_page(user_id).await;
        assert_eq!(page, Some(3));

        // Can't go past total pages
        let page = manager.next_page(user_id).await;
        assert_eq!(page, None);
    }

    #[tokio::test]
    async fn test_update_state() {
        let manager = WsConnectionManager::new();
        let user_id = Uuid::new_v4();
        let book_id = Uuid::new_v4();

        manager.start_lesson(user_id, book_id, 5).await;
        manager.update_state(user_id, PlaybackState::Playing).await;

        let session = manager.get_session(user_id).await.unwrap();
        assert_eq!(session.state, PlaybackState::Playing);
    }

    #[tokio::test]
    async fn test_end_session() {
        let manager = WsConnectionManager::new();
        let user_id = Uuid::new_v4();
        let book_id = Uuid::new_v4();

        manager.start_lesson(user_id, book_id, 5).await;
        manager.end_session(user_id).await;

        let session = manager.get_session(user_id).await;
        assert!(session.is_none());
    }

    #[test]
    fn test_ws_message_serialization() {
        let msg = WsMessage::Progress {
            page: 5,
            total_pages: 10,
            time_elapsed_secs: 120,
        };

        let json = serde_json::to_string(&msg).unwrap();
        assert!(json.contains("Progress"));
        assert!(json.contains("\"page\":5"));
    }
}
