//! API Client for backend communication

use gloo_net::http::Request;
use gloo_storage::{LocalStorage, Storage};
use serde::{Deserialize, Serialize};

use crate::routes::BookItem;

const API_BASE: &str = "/api";
const TOKEN_KEY: &str = "hailango_token";

/// API client for backend communication
pub struct ApiClient;

#[derive(Serialize)]
struct LoginRequest<'a> {
    email: &'a str,
    password: &'a str,
}

#[derive(Serialize)]
struct RegisterRequest<'a> {
    email: &'a str,
    password: &'a str,
    display_name: &'a str,
}

#[derive(Deserialize)]
struct AuthResponse {
    token: String,
}

#[derive(Deserialize)]
struct ErrorResponse {
    error: String,
}

impl ApiClient {
    /// Get stored auth token
    pub fn get_token() -> Option<String> {
        LocalStorage::get(TOKEN_KEY).ok()
    }

    /// Store auth token
    pub fn set_token(token: &str) {
        let _ = LocalStorage::set(TOKEN_KEY, token);
    }

    /// Clear auth token
    pub fn clear_token() {
        LocalStorage::delete(TOKEN_KEY);
    }

    /// Check if authenticated
    pub fn is_authenticated() -> bool {
        Self::get_token().is_some()
    }

    /// Login user
    pub async fn login(email: &str, password: &str) -> Result<(), String> {
        let response = Request::post(&format!("{}/auth/login", API_BASE))
            .json(&LoginRequest { email, password })
            .map_err(|e| e.to_string())?
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() {
            let auth: AuthResponse = response.json().await.map_err(|e| e.to_string())?;
            Self::set_token(&auth.token);
            Ok(())
        } else {
            let err: ErrorResponse = response
                .json()
                .await
                .unwrap_or(ErrorResponse {
                    error: "Login failed".to_string(),
                });
            Err(err.error)
        }
    }

    /// Register new user
    pub async fn register(email: &str, password: &str, display_name: &str) -> Result<(), String> {
        let response = Request::post(&format!("{}/auth/register", API_BASE))
            .json(&RegisterRequest {
                email,
                password,
                display_name,
            })
            .map_err(|e| e.to_string())?
            .send()
            .await
            .map_err(|e| e.to_string())?;

        if response.ok() {
            Ok(())
        } else {
            let err: ErrorResponse = response
                .json()
                .await
                .unwrap_or(ErrorResponse {
                    error: "Registration failed".to_string(),
                });
            Err(err.error)
        }
    }

    /// Fetch user's books
    pub async fn get_books() -> Result<Vec<BookItem>, String> {
        let mut request = Request::get(&format!("{}/books", API_BASE));

        if let Some(token) = Self::get_token() {
            request = request.header("Authorization", &format!("Bearer {}", token));
        }

        let response = request.send().await.map_err(|e| e.to_string())?;

        if response.ok() {
            response.json().await.map_err(|e| e.to_string())
        } else {
            Err("Failed to fetch books".to_string())
        }
    }

    /// Logout
    pub fn logout() {
        Self::clear_token();
        if let Some(window) = web_sys::window() {
            let _ = window.location().set_href("/login");
        }
    }
}
