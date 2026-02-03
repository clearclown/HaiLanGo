//! Data Transfer Objects for authentication

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Request to register a new user
#[derive(Debug, Deserialize)]
pub struct RegisterRequest {
    pub email: String,
    pub password: String,
    pub display_name: String,
    #[serde(default = "default_language")]
    pub native_language: String,
}

fn default_language() -> String {
    "en".to_string()
}

/// Request to login
#[derive(Debug, Deserialize)]
pub struct LoginRequest {
    pub email: String,
    pub password: String,
}

/// Response containing user info and tokens
#[derive(Debug, Serialize)]
pub struct AuthResponse {
    pub user: UserResponse,
    pub tokens: TokenResponse,
}

/// User information for responses
#[derive(Debug, Serialize)]
pub struct UserResponse {
    pub id: Uuid,
    pub email: String,
    pub display_name: String,
    pub native_language: String,
    pub email_verified: bool,
    pub created_at: DateTime<Utc>,
}

/// JWT token response
#[derive(Debug, Serialize)]
pub struct TokenResponse {
    pub access_token: String,
    pub refresh_token: String,
    pub expires_in: u64,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_register_request_deserialization() {
        let json = r#"{
            "email": "test@example.com",
            "password": "password123",
            "display_name": "Test User"
        }"#;

        let request: RegisterRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.email, "test@example.com");
        assert_eq!(request.native_language, "en"); // default
    }

    #[test]
    fn test_register_request_with_language() {
        let json = r#"{
            "email": "test@example.com",
            "password": "password123",
            "display_name": "Test User",
            "native_language": "ja"
        }"#;

        let request: RegisterRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.native_language, "ja");
    }

    #[test]
    fn test_login_request_deserialization() {
        let json = r#"{
            "email": "test@example.com",
            "password": "password123"
        }"#;

        let request: LoginRequest = serde_json::from_str(json).unwrap();

        assert_eq!(request.email, "test@example.com");
        assert_eq!(request.password, "password123");
    }

    #[test]
    fn test_token_response_serialization() {
        let response = TokenResponse {
            access_token: "access.token.here".to_string(),
            refresh_token: "refresh.token.here".to_string(),
            expires_in: 3600,
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("access_token"));
        assert!(json.contains("3600"));
    }

    #[test]
    fn test_user_response_serialization() {
        let now = Utc::now();
        let response = UserResponse {
            id: Uuid::new_v4(),
            email: "test@example.com".to_string(),
            display_name: "Test User".to_string(),
            native_language: "en".to_string(),
            email_verified: false,
            created_at: now,
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("test@example.com"));
        assert!(json.contains("Test User"));
    }

    #[test]
    fn test_auth_response_serialization() {
        let now = Utc::now();
        let response = AuthResponse {
            user: UserResponse {
                id: Uuid::new_v4(),
                email: "test@example.com".to_string(),
                display_name: "Test User".to_string(),
                native_language: "en".to_string(),
                email_verified: false,
                created_at: now,
            },
            tokens: TokenResponse {
                access_token: "token".to_string(),
                refresh_token: "refresh".to_string(),
                expires_in: 3600,
            },
        };

        let json = serde_json::to_string(&response).unwrap();
        assert!(json.contains("access_token"));
        assert!(json.contains("user"));
    }
}
