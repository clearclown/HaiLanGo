//! STT (Speech-to-Text) API routes

use async_trait::async_trait;
use base64::{Engine as _, engine::general_purpose};
use serde::Deserialize;
use serde_json::json;
use std::sync::{Arc, RwLock};
use uuid::Uuid;

use crate::apps::stt::{
    dto::PronunciationRequest,
    models::{PronunciationAttempt, WordFeedback},
    views::{EvaluateResult, SttViewSet},
};
use crate::services::stt::{MockSttProvider, SttAudioFormat, SttProvider};
use crate::{Handler, Method, Request, Response, Result, Route, path};

/// Shared STT state
#[derive(Clone)]
pub struct SttState {
    pub provider: Arc<dyn SttProvider>,
    pub attempts: Arc<RwLock<Vec<PronunciationAttempt>>>,
    pub feedbacks: Arc<RwLock<Vec<WordFeedback>>>,
}

impl Default for SttState {
    fn default() -> Self {
        Self {
            provider: Arc::new(MockSttProvider::new()),
            attempts: Arc::new(RwLock::new(Vec::new())),
            feedbacks: Arc::new(RwLock::new(Vec::new())),
        }
    }
}

/// Request body for the evaluate endpoint (with base64 audio)
#[derive(Debug, Deserialize)]
struct EvaluateApiRequest {
    expected_text: String,
    language: String,
    audio_base64: String,
    #[serde(default)]
    audio_format: Option<String>,
    page_id: Option<Uuid>,
}

/// Handler for POST /evaluate/
struct EvaluateHandler {
    state: SttState,
}

#[async_trait]
impl Handler for EvaluateHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let user_id = request
            .extensions
            .get::<Uuid>()
            .unwrap_or_else(Uuid::new_v4);

        let req: EvaluateApiRequest = request.json()?;

        // Decode base64 audio
        let audio_data = match general_purpose::STANDARD.decode(&req.audio_base64) {
            Ok(data) => data,
            Err(_) => {
                return Response::bad_request()
                    .with_json(&json!({"error": "Invalid base64 audio data"}));
            }
        };

        // Parse audio format
        let audio_format = match req.audio_format.as_deref() {
            Some("mp3") => SttAudioFormat::Mp3,
            Some("ogg") => SttAudioFormat::Ogg,
            Some("webm") => SttAudioFormat::Webm,
            _ => SttAudioFormat::Wav,
        };

        let pronunciation_req = PronunciationRequest {
            expected_text: req.expected_text,
            language: req.language,
            audio_format,
            page_id: req.page_id,
        };

        let result = SttViewSet::evaluate(
            user_id,
            pronunciation_req,
            audio_data,
            self.state.provider.as_ref(),
        )
        .await;

        match result {
            EvaluateResult::Success(response) => Response::ok().with_json(&response),
            EvaluateResult::InvalidInput(msg) => {
                Response::bad_request().with_json(&json!({"error": msg}))
            }
            EvaluateResult::ServiceError(msg) => {
                Response::internal_server_error().with_json(&json!({"error": msg}))
            }
            EvaluateResult::Unauthorized => {
                Response::unauthorized().with_json(&json!({"error": "Unauthorized"}))
            }
        }
    }
}

/// Request body for the transcribe endpoint
#[derive(Debug, Deserialize)]
struct TranscribeApiRequest {
    language: String,
    audio_base64: String,
    #[serde(default)]
    audio_format: Option<String>,
}

/// Handler for POST /transcribe/
struct TranscribeHandler {
    state: SttState,
}

#[async_trait]
impl Handler for TranscribeHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let req: TranscribeApiRequest = request.json()?;

        let audio_data = match general_purpose::STANDARD.decode(&req.audio_base64) {
            Ok(data) => data,
            Err(_) => {
                return Response::bad_request()
                    .with_json(&json!({"error": "Invalid base64 audio data"}));
            }
        };

        if audio_data.is_empty() {
            return Response::bad_request()
                .with_json(&json!({"error": "Audio data is required"}));
        }

        let audio_format = match req.audio_format.as_deref() {
            Some("mp3") => SttAudioFormat::Mp3,
            Some("ogg") => SttAudioFormat::Ogg,
            Some("webm") => SttAudioFormat::Webm,
            _ => SttAudioFormat::Wav,
        };

        use crate::services::stt::SttRequest;
        let stt_request = SttRequest::new(audio_data, req.language.clone())
            .with_format(audio_format);

        match self.state.provider.transcribe(stt_request).await {
            Ok(response) => Response::ok().with_json(&json!({
                "text": response.text,
                "language": response.language,
                "confidence": response.confidence,
                "duration_ms": response.duration_ms,
            })),
            Err(e) => Response::internal_server_error().with_json(&json!({"error": e.to_string()})),
        }
    }
}

/// Handler for GET/POST /attempts/
struct AttemptsHandler {
    state: SttState,
}

#[async_trait]
impl Handler for AttemptsHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let user_id = request
            .extensions
            .get::<Uuid>()
            .unwrap_or_else(Uuid::new_v4);

        match request.method {
            Method::GET => {
                let attempts = self.state.attempts.read().unwrap();
                let result = SttViewSet::list_attempts(user_id, &attempts);
                Response::ok().with_json(&result)
            }
            _ => Err(crate::Error::MethodNotAllowed("Only GET is allowed".into())),
        }
    }
}

/// Handler for GET /stats/
struct StatsHandler {
    state: SttState,
}

#[async_trait]
impl Handler for StatsHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let user_id = request
            .extensions
            .get::<Uuid>()
            .unwrap_or_else(Uuid::new_v4);

        let attempts = self.state.attempts.read().unwrap();
        let feedbacks = self.state.feedbacks.read().unwrap();
        let stats = SttViewSet::stats(user_id, &attempts, &feedbacks);
        Response::ok().with_json(&stats)
    }
}

/// Create STT routes
pub fn routes() -> Vec<Route> {
    let state = SttState::default();

    vec![
        path(
            "/evaluate/",
            EvaluateHandler {
                state: state.clone(),
            },
        ),
        path(
            "/transcribe/",
            TranscribeHandler {
                state: state.clone(),
            },
        ),
        path(
            "/attempts/",
            AttemptsHandler {
                state: state.clone(),
            },
        ),
        path("/stats/", StatsHandler { state }),
    ]
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::Method;
    use bytes::Bytes;

    fn build_handler() -> (EvaluateHandler, TranscribeHandler, AttemptsHandler, StatsHandler) {
        let state = SttState::default();
        (
            EvaluateHandler {
                state: state.clone(),
            },
            TranscribeHandler {
                state: state.clone(),
            },
            AttemptsHandler {
                state: state.clone(),
            },
            StatsHandler { state },
        )
    }

    fn result_to_response(result: Result<Response>) -> Response {
        match result {
            Ok(r) => r,
            Err(e) => Response::from(e),
        }
    }

    #[tokio::test]
    async fn test_evaluate_success() {
        let (handler, _, _, _) = build_handler();

        // Minimal valid WAV audio encoded as base64
        let audio_b64 = general_purpose::STANDARD.encode(vec![0u8; 100]);
        let body = format!(
            r#"{{"expected_text":"hello world","language":"en","audio_base64":"{}"}}"#,
            audio_b64
        );

        let request = Request::builder()
            .method(Method::POST)
            .uri("/evaluate/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_evaluate_invalid_base64() {
        let (handler, _, _, _) = build_handler();

        let body = r#"{"expected_text":"hello","language":"en","audio_base64":"!!!invalid!!!"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/evaluate/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 400);
    }

    #[tokio::test]
    async fn test_transcribe_success() {
        let (_, handler, _, _) = build_handler();

        let audio_b64 = general_purpose::STANDARD.encode(vec![0u8; 100]);
        let body = format!(
            r#"{{"language":"en","audio_base64":"{}"}}"#,
            audio_b64
        );

        let request = Request::builder()
            .method(Method::POST)
            .uri("/transcribe/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_attempts_list() {
        let (_, _, handler, _) = build_handler();

        let request = Request::builder()
            .method(Method::GET)
            .uri("/attempts/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_stats() {
        let (_, _, _, handler) = build_handler();

        let request = Request::builder()
            .method(Method::GET)
            .uri("/stats/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }
}
