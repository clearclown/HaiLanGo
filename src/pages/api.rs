//! WASM API client for pages
//!
//! Thin HTTP wrapper around the backend REST API, used by the Reinhardt pages
//! (WASM frontend). Mirrors the types defined in `frontend/src/api.rs`.

#[cfg(target_arch = "wasm32")]
use gloo_net::http::Request;
#[cfg(target_arch = "wasm32")]
use gloo_storage::{LocalStorage, Storage};
use serde::{Deserialize, Serialize};

const TOKEN_KEY: &str = "hailango_token";

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
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct TokenInfo {
    pub access_token: String,
    pub refresh_token: String,
    pub expires_in: u64,
}

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

/// Get the stored access token.
#[cfg(target_arch = "wasm32")]
pub fn get_token() -> Option<String> {
    LocalStorage::get(TOKEN_KEY).ok()
}

/// Persist the access token from a successful auth response.
#[cfg(target_arch = "wasm32")]
fn store_tokens(tokens: &TokenInfo) {
    let _ = LocalStorage::set(TOKEN_KEY, &tokens.access_token);
}

/// Remove the stored token (logout).
#[cfg(target_arch = "wasm32")]
pub fn clear_token() {
    LocalStorage::delete(TOKEN_KEY);
}

/// POST /api/auth/login/
#[cfg(target_arch = "wasm32")]
pub async fn login(email: &str, password: &str) -> Result<AuthResponse, String> {
    let body =
        serde_json::to_string(&LoginPayload { email, password }).map_err(|e| e.to_string())?;

    let resp = Request::post("/api/auth/login/")
        .header("Content-Type", "application/json")
        .body(body)
        .map_err(|e| e.to_string())?
        .send()
        .await
        .map_err(|e| e.to_string())?;

    if !resp.ok() {
        return Err(format!("Login failed: HTTP {}", resp.status()));
    }

    let auth: AuthResponse = resp.json().await.map_err(|e| e.to_string())?;
    store_tokens(&auth.tokens);
    Ok(auth)
}

/// POST /api/auth/register/
#[cfg(target_arch = "wasm32")]
pub async fn register(
    email: &str,
    password: &str,
    display_name: &str,
) -> Result<AuthResponse, String> {
    let body = serde_json::to_string(&RegisterPayload {
        email,
        password,
        display_name,
    })
    .map_err(|e| e.to_string())?;

    let resp = Request::post("/api/auth/register/")
        .header("Content-Type", "application/json")
        .body(body)
        .map_err(|e| e.to_string())?
        .send()
        .await
        .map_err(|e| e.to_string())?;

    if !resp.ok() {
        return Err(format!("Registration failed: HTTP {}", resp.status()));
    }

    let auth: AuthResponse = resp.json().await.map_err(|e| e.to_string())?;
    store_tokens(&auth.tokens);
    Ok(auth)
}

/// GET request with Authorization header.
#[cfg(target_arch = "wasm32")]
async fn authed_get(path: &str) -> Result<gloo_net::http::Response, String> {
    let token = get_token().unwrap_or_default();
    Request::get(path)
        .header("Authorization", &format!("Bearer {token}"))
        .send()
        .await
        .map_err(|e| e.to_string())
}

#[derive(Clone, Debug, Deserialize, Default)]
pub struct DashboardStats {
    pub words_learned: u32,
    pub books_count: u32,
    pub streak_days: u32,
}

/// Aggregate dashboard stats from /api/review/stats/ and /api/books/.
#[cfg(target_arch = "wasm32")]
pub async fn fetch_dashboard_stats() -> Result<DashboardStats, String> {
    #[derive(Deserialize)]
    struct ReviewStats {
        learned_count: u32,
        streak_days: i32,
    }

    #[derive(Deserialize)]
    struct BooksResponse {
        books: Vec<serde_json::Value>,
    }

    let stats_resp = authed_get("/api/review/stats/").await?;
    let books_resp = authed_get("/api/books/").await?;

    let stats: ReviewStats = stats_resp.json().await.map_err(|e| e.to_string())?;
    let books: BooksResponse = books_resp.json().await.map_err(|e| e.to_string())?;

    Ok(DashboardStats {
        words_learned: stats.learned_count,
        books_count: books.books.len() as u32,
        streak_days: stats.streak_days.max(0) as u32,
    })
}
