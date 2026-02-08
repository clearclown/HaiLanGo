//! TTS API routes

use axum::{
    Router,
    extract::State,
    http::StatusCode,
    response::{IntoResponse, Json},
    routing::{get, post},
};
use serde_json::json;
use std::sync::{Arc, RwLock};
use uuid::Uuid;

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

/// POST /api/tts/synthesize
async fn synthesize(
    State(state): State<TtsState>,
    Json(request): Json<SynthesizeRequest>,
) -> impl IntoResponse {
    // Mock user_id (in production, extract from JWT)
    let user_id = Uuid::new_v4();

    let result = TtsViewSet::synthesize(user_id, request, state.provider.as_ref()).await;

    match result {
        SynthesizeResult::Success {
            metadata,
            audio_data,
            generation,
        } => {
            // Store generation record
            state.generations.write().unwrap().push(*generation);

            // Return metadata + base64 audio in JSON
            // (In production, audio would be streamed or stored in object storage)
            use base64::{Engine, engine::general_purpose::STANDARD};
            let audio_b64 = STANDARD.encode(&audio_data);

            (
                StatusCode::OK,
                Json(json!({
                    "metadata": metadata,
                    "audio_base64": audio_b64
                })),
            )
        }
        SynthesizeResult::InvalidInput(msg) => {
            (StatusCode::BAD_REQUEST, Json(json!({"error": msg})))
        }
        SynthesizeResult::ServiceError(msg) => (
            StatusCode::SERVICE_UNAVAILABLE,
            Json(json!({"error": msg})),
        ),
    }
}

/// GET /api/tts/history
async fn list_history(State(state): State<TtsState>) -> impl IntoResponse {
    // Mock user_id
    let user_id = Uuid::new_v4();
    let generations = state.generations.read().unwrap();
    let history = TtsViewSet::list_history(user_id, &generations);
    (StatusCode::OK, Json(json!(history)))
}

/// GET /api/tts/languages
async fn supported_languages(State(state): State<TtsState>) -> impl IntoResponse {
    let response = TtsViewSet::supported_languages(state.provider.as_ref());
    (StatusCode::OK, Json(json!(response)))
}

/// Create TTS router
pub fn router() -> Router {
    let state = TtsState::default();

    Router::new()
        .route("/synthesize", post(synthesize))
        .route("/history", get(list_history))
        .route("/languages", get(supported_languages))
        .with_state(state)
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::{body::Body, http::Request};
    use tower::ServiceExt;

    #[tokio::test]
    async fn test_synthesize_success() {
        let app = router();

        let body = r#"{"text":"Hello world","language":"en"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/synthesize")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_synthesize_with_options() {
        let app = router();

        let body = r#"{"text":"Bonjour le monde","language":"fr","speed":1.2,"format":"ogg","quality":"premium"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/synthesize")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_synthesize_empty_text() {
        let app = router();

        let body = r#"{"text":"","language":"en"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/synthesize")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_synthesize_unsupported_language() {
        let app = router();

        let body = r#"{"text":"test","language":"xyz"}"#;

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/synthesize")
                    .header("content-type", "application/json")
                    .body(Body::from(body))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_list_history() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/history")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_supported_languages() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/languages")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_synthesize_invalid_json() {
        let app = router();

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/synthesize")
                    .header("content-type", "application/json")
                    .body(Body::from("not-json"))
                    .unwrap(),
            )
            .await
            .unwrap();

        // Axum returns 400 for malformed JSON bodies
        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }
}
