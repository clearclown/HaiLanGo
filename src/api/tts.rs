//! TTS API routes

use async_trait::async_trait;
use serde_json::json;
use std::sync::{Arc, RwLock};
use uuid::Uuid;

use crate::{Handler, Request, Response, Result, Route, StatusCode, path};
use crate::apps::tts::{
    dto::SynthesizeRequest,
    models::AudioGeneration,
    views::{SynthesizeResult, TtsViewSet},
};
use crate::services::tts::{MockTtsProvider, TtsProvider, create_tts_provider};

/// Shared TTS state
#[derive(Clone)]
pub struct TtsState {
    pub provider: Arc<dyn TtsProvider>,
    pub generations: Arc<RwLock<Vec<AudioGeneration>>>,
}

impl Default for TtsState {
    fn default() -> Self {
        Self {
            provider: Arc::new(MockTtsProvider::new()),
            generations: Arc::new(RwLock::new(Vec::new())),
        }
    }
}

impl TtsState {
    /// Create state with auto-detected provider
    pub fn with_auto_provider() -> Self {
        Self {
            provider: Arc::from(create_tts_provider()),
            generations: Arc::new(RwLock::new(Vec::new())),
        }
    }
}

/// Handler for POST /synthesize/
struct SynthesizeHandler {
    state: TtsState,
}

#[async_trait]
impl Handler for SynthesizeHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let req: SynthesizeRequest = request.json()?;
        let user_id = Uuid::new_v4();

        let result = TtsViewSet::synthesize(user_id, req, self.state.provider.as_ref()).await;

        match result {
            SynthesizeResult::Success {
                metadata,
                audio_data,
                generation,
            } => {
                self.state.generations.write().unwrap().push(*generation);

                use base64::{Engine, engine::general_purpose::STANDARD};
                let audio_b64 = STANDARD.encode(&audio_data);

                Response::ok().with_json(&json!({
                    "metadata": metadata,
                    "audio_base64": audio_b64
                }))
            }
            SynthesizeResult::InvalidInput(msg) => {
                Response::bad_request().with_json(&json!({"error": msg}))
            }
            SynthesizeResult::ServiceError(msg) => Response::new(StatusCode::SERVICE_UNAVAILABLE)
                .with_json(&json!({"error": msg})),
        }
    }
}

/// Handler for GET /history/
struct HistoryHandler {
    state: TtsState,
}

#[async_trait]
impl Handler for HistoryHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        let user_id = Uuid::new_v4();
        let generations = self.state.generations.read().unwrap();
        let history = TtsViewSet::list_history(user_id, &generations);
        Response::ok().with_json(&history)
    }
}

/// Handler for GET /languages/
struct LanguagesHandler {
    state: TtsState,
}

#[async_trait]
impl Handler for LanguagesHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        let response = TtsViewSet::supported_languages(self.state.provider.as_ref());
        Response::ok().with_json(&response)
    }
}

/// Create TTS routes
pub fn routes() -> Vec<Route> {
    let state = TtsState::default();

    vec![
        path(
            "/synthesize/",
            SynthesizeHandler {
                state: state.clone(),
            },
        ),
        path("/history/", HistoryHandler { state: state.clone() }),
        path("/languages/", LanguagesHandler { state }),
    ]
}

#[cfg(test)]
mod tests {
    use super::*;
    use bytes::Bytes;
    use crate::Method;

    fn result_to_response(result: Result<Response>) -> Response {
        match result {
            Ok(r) => r,
            Err(e) => Response::from(e),
        }
    }

    fn make_handler() -> SynthesizeHandler {
        SynthesizeHandler {
            state: TtsState::default(),
        }
    }

    #[tokio::test]
    async fn test_synthesize_success() {
        let handler = make_handler();
        let body = r#"{"text":"Hello world","language":"en"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/synthesize/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_synthesize_with_options() {
        let handler = make_handler();
        let body = r#"{"text":"Bonjour le monde","language":"fr","speed":1.2,"format":"ogg","quality":"premium"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/synthesize/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_synthesize_empty_text() {
        let handler = make_handler();
        let body = r#"{"text":"","language":"en"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/synthesize/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 400);
    }

    #[tokio::test]
    async fn test_synthesize_unsupported_language() {
        let handler = make_handler();
        let body = r#"{"text":"test","language":"xyz"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/synthesize/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 400);
    }

    #[tokio::test]
    async fn test_list_history() {
        let state = TtsState::default();
        let handler = HistoryHandler { state };

        let request = Request::builder()
            .method(Method::GET)
            .uri("/history/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_supported_languages() {
        let state = TtsState::default();
        let handler = LanguagesHandler { state };

        let request = Request::builder()
            .method(Method::GET)
            .uri("/languages/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_synthesize_invalid_json() {
        let handler = make_handler();

        let request = Request::builder()
            .method(Method::POST)
            .uri("/synthesize/")
            .header("content-type", "application/json")
            .body(Bytes::from("not-json"))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        // request.json() returns Error::Serialization which maps to 400
        assert_eq!(response.status, 400);
    }
}
