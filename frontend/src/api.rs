//! API Client for backend communication

use gloo_net::http::Request;
use gloo_storage::{LocalStorage, Storage};
use serde::{Deserialize, Serialize};

const API_BASE: &str = "/api";
const TOKEN_KEY: &str = "hailango_token";
const REFRESH_TOKEN_KEY: &str = "hailango_refresh_token";

// ── Auth types ──

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct AuthResponse {
    pub user: UserInfo,
    pub tokens: TokenInfo,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct UserInfo {
    pub id: String,
    pub email: String,
    pub display_name: String,
    #[serde(default)]
    pub native_language: String,
    #[serde(default)]
    pub email_verified: bool,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct TokenInfo {
    pub access_token: String,
    pub refresh_token: String,
    pub expires_in: u64,
}

#[derive(Clone, Debug, Deserialize)]
pub struct OAuthProviderInfo {
    pub name: String,
    pub configured: bool,
}

#[derive(Clone, Debug, Deserialize)]
pub struct OAuthRedirectResponse {
    pub auth_url: String,
    pub state: String,
}

// ── Book types ──

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct BookItem {
    pub id: String,
    pub title: String,
    pub author: String,
    #[serde(default)]
    pub language: String,
    pub total_pages: u32,
    #[serde(default)]
    pub progress: f32,
}

#[derive(Clone, Debug, Deserialize)]
pub struct BookListResponse {
    pub books: Vec<BookItem>,
}

// ── Learning types ──

#[derive(Clone, Debug, Deserialize)]
pub struct PageContent {
    pub page_number: u32,
    pub text: String,
    pub book_title: String,
    pub total_pages: u32,
}

#[derive(Clone, Debug, Deserialize)]
pub struct LearningSession {
    pub id: String,
    pub book_id: String,
    pub book_title: String,
    pub current_page: u32,
    pub total_pages: u32,
}

// ── TTS types ──

#[derive(Clone, Debug, Deserialize)]
pub struct TtsLanguage {
    pub code: String,
    pub name: String,
}

#[derive(Clone, Debug, Serialize)]
pub struct TtsSynthesizeRequest {
    pub text: String,
    pub language: String,
}

// ── Teacher Mode types ──

#[derive(Clone, Debug, Serialize)]
pub struct StartLessonRequest {
    pub book_id: String,
    pub start_page: u32,
    pub end_page: u32,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub speed: Option<f32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub page_interval: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub repeat_count: Option<u32>,
}

#[derive(Clone, Debug, Deserialize)]
pub struct LessonStatusResponse {
    pub session_id: String,
    pub status: String,
    pub current_page: u32,
    pub total_pages: u32,
    #[serde(default)]
    pub pages_completed: u32,
}

#[derive(Clone, Debug, Serialize)]
pub struct UpdateTeacherConfig {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub speed: Option<f32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub page_interval: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub repeat_count: Option<u32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub auto_advance: Option<bool>,
}

// ── Review types ──

#[derive(Clone, Debug, Deserialize)]
pub struct ReviewCard {
    pub id: String,
    pub word: String,
    pub reading: String,
    pub meaning: String,
    #[serde(default)]
    pub sentence: Option<String>,
}

#[derive(Clone, Debug, Deserialize)]
pub struct ReviewStats {
    pub total_cards: u32,
    pub due_today: u32,
    pub streak: u32,
    pub accuracy: f32,
}

#[derive(Clone, Debug, Serialize)]
pub struct SubmitReviewRequest {
    pub card_id: String,
    pub rating: u8,
}

#[derive(Clone, Debug, Deserialize)]
pub struct ReviewSubmitResponse {
    pub next_review: String,
    pub interval_days: u32,
}

// ── Error response ──

#[derive(Deserialize)]
struct ErrorResponse {
    error: String,
}

// ── Request payloads ──

#[derive(Serialize)]
struct LoginPayload<'a> {
    email: &'a str,
    password: &'a str,
}

#[derive(Serialize)]
struct RegisterPayload<'a> {
    email: &'a str,
    password: &'a str,
    display_name: &'a str,
}

// ── API Client ──

pub struct ApiClient;

impl ApiClient {
    // ── Token management ──

    pub fn get_token() -> Option<String> {
        LocalStorage::get(TOKEN_KEY).ok()
    }

    pub fn set_token(token: &str) {
        let _ = LocalStorage::set(TOKEN_KEY, token);
    }

    pub fn set_refresh_token(token: &str) {
        let _ = LocalStorage::set(REFRESH_TOKEN_KEY, token);
    }

    pub fn clear_token() {
        LocalStorage::delete(TOKEN_KEY);
        LocalStorage::delete(REFRESH_TOKEN_KEY);
    }

    pub fn is_authenticated() -> bool {
        Self::get_token().is_some()
    }

    fn auth_header() -> Option<String> {
        Self::get_token().map(|t| format!("Bearer {}", t))
    }

    // ── Auth endpoints ──

    pub async fn login(email: &str, password: &str) -> Result<AuthResponse, String> {
        let response = Request::post(&format!("{}/auth/login", API_BASE))
            .json(&LoginPayload { email, password })
            .map_err(|e| e.to_string())?
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() {
            let auth: AuthResponse = response.json().await.map_err(|e| e.to_string())?;
            Self::set_token(&auth.tokens.access_token);
            Self::set_refresh_token(&auth.tokens.refresh_token);
            Ok(auth)
        } else {
            let err: ErrorResponse = response
                .json()
                .await
                .unwrap_or(ErrorResponse { error: "Login failed".to_string() });
            Err(err.error)
        }
    }

    pub async fn register(email: &str, password: &str, display_name: &str) -> Result<AuthResponse, String> {
        let response = Request::post(&format!("{}/auth/register", API_BASE))
            .json(&RegisterPayload { email, password, display_name })
            .map_err(|e| e.to_string())?
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() || response.status() == 201 {
            let auth: AuthResponse = response.json().await.map_err(|e| e.to_string())?;
            Self::set_token(&auth.tokens.access_token);
            Self::set_refresh_token(&auth.tokens.refresh_token);
            Ok(auth)
        } else {
            let err: ErrorResponse = response
                .json()
                .await
                .unwrap_or(ErrorResponse { error: "Registration failed".to_string() });
            Err(err.error)
        }
    }

    pub fn logout() {
        Self::clear_token();
        if let Some(window) = web_sys::window() {
            let _ = window.location().set_href("/login");
        }
    }

    // ── OAuth endpoints ──

    pub async fn get_oauth_providers() -> Result<Vec<OAuthProviderInfo>, String> {
        let response = Request::get(&format!("{}/auth/oauth/providers", API_BASE))
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() {
            #[derive(Deserialize)]
            struct Wrapper { providers: Vec<OAuthProviderInfo> }
            let wrapper: Wrapper = response.json().await.map_err(|e| e.to_string())?;
            Ok(wrapper.providers)
        } else {
            Ok(vec![])
        }
    }

    pub async fn get_oauth_url(provider: &str) -> Result<OAuthRedirectResponse, String> {
        let response = Request::get(&format!("{}/auth/oauth/{}", API_BASE, provider))
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() {
            response.json().await.map_err(|e| e.to_string())
        } else {
            let err: ErrorResponse = response
                .json()
                .await
                .unwrap_or(ErrorResponse { error: "OAuth error".to_string() });
            Err(err.error)
        }
    }

    pub async fn oauth_callback(provider: &str, code: &str, state: &str) -> Result<AuthResponse, String> {
        let response = Request::get(
            &format!("{}/auth/callback/{}?code={}&state={}", API_BASE, provider, code, state),
        )
        .send()
        .await
        .map_err(|e| e.to_string())?;

        if response.ok() {
            let auth: AuthResponse = response.json().await.map_err(|e| e.to_string())?;
            Self::set_token(&auth.tokens.access_token);
            Self::set_refresh_token(&auth.tokens.refresh_token);
            Ok(auth)
        } else {
            let err: ErrorResponse = response
                .json()
                .await
                .unwrap_or(ErrorResponse { error: "OAuth callback failed".to_string() });
            Err(err.error)
        }
    }

    // ── Book endpoints ──

    pub async fn get_books() -> Result<Vec<BookItem>, String> {
        let mut request = Request::get(&format!("{}/books", API_BASE));
        if let Some(auth) = Self::auth_header() {
            request = request.header("Authorization", &auth);
        }

        let response = request.send().await.map_err(|e| e.to_string())?;

        if response.ok() {
            #[derive(Deserialize)]
            struct Wrapper { books: Vec<BookItem> }
            let wrapper: Wrapper = response.json().await.map_err(|e| e.to_string())?;
            Ok(wrapper.books)
        } else {
            // Fallback: try parsing as direct array
            Ok(vec![])
        }
    }

    pub async fn upload_book(title: &str, author: &str, language: &str) -> Result<BookItem, String> {
        #[derive(Serialize)]
        struct UploadReq<'a> { title: &'a str, author: &'a str, language: &'a str }

        let mut builder = Request::post(&format!("{}/books", API_BASE));
        if let Some(auth) = Self::auth_header() {
            builder = builder.header("Authorization", &auth);
        }

        let response = builder
            .json(&UploadReq { title, author, language })
            .map_err(|e| e.to_string())?
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() || response.status() == 201 {
            response.json().await.map_err(|e| e.to_string())
        } else {
            let err: ErrorResponse = response
                .json()
                .await
                .unwrap_or(ErrorResponse { error: "Upload failed".to_string() });
            Err(err.error)
        }
    }

    // ── Learning endpoints ──

    pub async fn get_learning_sessions() -> Result<Vec<LearningSession>, String> {
        let mut request = Request::get(&format!("{}/learning/sessions", API_BASE));
        if let Some(auth) = Self::auth_header() {
            request = request.header("Authorization", &auth);
        }

        let response = request.send().await.map_err(|e| e.to_string())?;

        if response.ok() {
            #[derive(Deserialize)]
            struct Wrapper { sessions: Vec<LearningSession> }
            let wrapper: Wrapper = response.json().await.map_err(|e| e.to_string())?;
            Ok(wrapper.sessions)
        } else {
            Ok(vec![])
        }
    }

    pub async fn get_page_content(book_id: &str, page: u32) -> Result<PageContent, String> {
        let mut request = Request::get(
            &format!("{}/learning/books/{}/pages/{}", API_BASE, book_id, page),
        );
        if let Some(auth) = Self::auth_header() {
            request = request.header("Authorization", &auth);
        }

        let response = request.send().await.map_err(|e| e.to_string())?;

        if response.ok() {
            response.json().await.map_err(|e| e.to_string())
        } else {
            Err("Failed to load page content".to_string())
        }
    }

    // ── TTS endpoints ──

    pub async fn get_tts_languages() -> Result<Vec<TtsLanguage>, String> {
        let response = Request::get(&format!("{}/tts/languages", API_BASE))
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() {
            #[derive(Deserialize)]
            struct Wrapper { languages: Vec<TtsLanguage> }
            let wrapper: Wrapper = response.json().await.map_err(|e| e.to_string())?;
            Ok(wrapper.languages)
        } else {
            Ok(vec![])
        }
    }

    pub async fn synthesize_speech(text: &str, language: &str) -> Result<String, String> {
        let mut builder = Request::post(&format!("{}/tts/synthesize", API_BASE));
        if let Some(auth) = Self::auth_header() {
            builder = builder.header("Authorization", &auth);
        }

        let response = builder
            .json(&TtsSynthesizeRequest {
                text: text.to_string(),
                language: language.to_string(),
            })
            .map_err(|e| e.to_string())?
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() {
            #[derive(Deserialize)]
            struct Wrapper { audio_url: String }
            let wrapper: Wrapper = response.json().await.map_err(|e| e.to_string())?;
            Ok(wrapper.audio_url)
        } else {
            Err("TTS synthesis failed".to_string())
        }
    }

    // ── Teacher Mode endpoints ──

    pub async fn start_lesson(req: StartLessonRequest) -> Result<LessonStatusResponse, String> {
        let mut builder = Request::post(&format!("{}/teacher/start", API_BASE));
        if let Some(auth) = Self::auth_header() {
            builder = builder.header("Authorization", &auth);
        }

        let response = builder
            .json(&req)
            .map_err(|e| e.to_string())?
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() || response.status() == 201 {
            response.json().await.map_err(|e| e.to_string())
        } else {
            let err: ErrorResponse = response
                .json()
                .await
                .unwrap_or(ErrorResponse { error: "Failed to start lesson".to_string() });
            Err(err.error)
        }
    }

    pub async fn teacher_action(action: &str) -> Result<LessonStatusResponse, String> {
        let mut builder = Request::post(&format!("{}/teacher/{}", API_BASE, action));
        if let Some(auth) = Self::auth_header() {
            builder = builder.header("Authorization", &auth);
        }

        let response = builder
            .json(&serde_json::json!({}))
            .map_err(|e| e.to_string())?
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() {
            response.json().await.map_err(|e| e.to_string())
        } else {
            let err: ErrorResponse = response
                .json()
                .await
                .unwrap_or(ErrorResponse { error: format!("Teacher {} failed", action) });
            Err(err.error)
        }
    }

    pub async fn update_teacher_config(config: UpdateTeacherConfig) -> Result<(), String> {
        let mut builder = Request::put(&format!("{}/teacher/config", API_BASE));
        if let Some(auth) = Self::auth_header() {
            builder = builder.header("Authorization", &auth);
        }

        let response = builder
            .json(&config)
            .map_err(|e| e.to_string())?
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() {
            Ok(())
        } else {
            Err("Config update failed".to_string())
        }
    }

    pub async fn get_teacher_status() -> Result<LessonStatusResponse, String> {
        let mut request = Request::get(&format!("{}/teacher/status", API_BASE));
        if let Some(auth) = Self::auth_header() {
            request = request.header("Authorization", &auth);
        }

        let response = request.send().await.map_err(|e| e.to_string())?;

        if response.ok() {
            response.json().await.map_err(|e| e.to_string())
        } else {
            Err("Failed to get teacher status".to_string())
        }
    }

    // ── Review endpoints ──

    pub async fn get_review_stats() -> Result<ReviewStats, String> {
        let mut request = Request::get(&format!("{}/review/stats", API_BASE));
        if let Some(auth) = Self::auth_header() {
            request = request.header("Authorization", &auth);
        }

        let response = request.send().await.map_err(|e| e.to_string())?;

        if response.ok() {
            response.json().await.map_err(|e| e.to_string())
        } else {
            Ok(ReviewStats { total_cards: 0, due_today: 0, streak: 0, accuracy: 0.0 })
        }
    }

    pub async fn get_review_queue() -> Result<Vec<ReviewCard>, String> {
        let mut request = Request::get(&format!("{}/review/queue", API_BASE));
        if let Some(auth) = Self::auth_header() {
            request = request.header("Authorization", &auth);
        }

        let response = request.send().await.map_err(|e| e.to_string())?;

        if response.ok() {
            #[derive(Deserialize)]
            struct Wrapper { cards: Vec<ReviewCard> }
            let wrapper: Wrapper = response.json().await.map_err(|e| e.to_string())?;
            Ok(wrapper.cards)
        } else {
            Ok(vec![])
        }
    }

    pub async fn submit_review(card_id: &str, rating: u8) -> Result<ReviewSubmitResponse, String> {
        let mut builder = Request::post(&format!("{}/review/submit", API_BASE));
        if let Some(auth) = Self::auth_header() {
            builder = builder.header("Authorization", &auth);
        }

        let response = builder
            .json(&SubmitReviewRequest {
                card_id: card_id.to_string(),
                rating,
            })
            .map_err(|e| e.to_string())?
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() {
            response.json().await.map_err(|e| e.to_string())
        } else {
            Err("Review submit failed".to_string())
        }
    }
}
