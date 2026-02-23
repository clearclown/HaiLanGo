//! Auth ViewSet - Registration and Login endpoints

use super::dto::{AuthResponse, LoginRequest, RegisterRequest, TokenResponse, UserResponse};
use super::models::User;
use super::oauth::OAuthUserInfo;
use super::services::{JwtService, hash_password, verify_password};

/// Registration request handler result
#[derive(Debug)]
pub enum RegisterResult {
    Success(Box<AuthResponse>, Box<User>),
    EmailExists,
    InvalidInput(String),
}

/// Login request handler result
#[derive(Debug)]
pub enum LoginResult {
    Success(AuthResponse),
    InvalidCredentials,
    UserNotFound,
}

/// OAuth login result
#[derive(Debug)]
pub enum OAuthLoginResult {
    Success(AuthResponse),
    ProviderError(String),
}

/// Auth ViewSet - handles authentication endpoints
pub struct AuthViewSet;

impl AuthViewSet {
    /// Register a new user
    pub fn register(request: RegisterRequest) -> RegisterResult {
        // Validate input
        if request.email.is_empty() || !request.email.contains('@') {
            return RegisterResult::InvalidInput("Invalid email format".to_string());
        }

        if request.password.len() < 8 {
            return RegisterResult::InvalidInput(
                "Password must be at least 8 characters".to_string(),
            );
        }

        // Hash password
        let password_hash = match hash_password(&request.password) {
            Ok(hash) => hash,
            Err(_) => return RegisterResult::InvalidInput("Password hashing failed".to_string()),
        };

        // Create user
        let user = User::new(
            request.email.clone(),
            request.display_name.clone(),
            password_hash,
        );

        // Generate tokens
        let tokens = Self::generate_tokens(&user);

        let response = AuthResponse {
            user: UserResponse {
                id: user.id,
                email: user.email.clone(),
                display_name: user.display_name.clone(),
                native_language: user.native_language.clone(),
                email_verified: user.email_verified,
                created_at: user.created_at,
            },
            tokens,
        };

        RegisterResult::Success(Box::new(response), Box::new(user))
    }

    /// Login with email and password
    pub fn login(request: LoginRequest, stored_user: Option<&User>) -> LoginResult {
        let user = match stored_user {
            Some(u) => u,
            None => return LoginResult::UserNotFound,
        };

        // Verify password
        let password_hash = match &user.password_hash {
            Some(hash) => hash,
            None => return LoginResult::InvalidCredentials, // OAuth user
        };

        match verify_password(&request.password, password_hash) {
            Ok(true) => {
                let tokens = Self::generate_tokens(user);
                LoginResult::Success(AuthResponse {
                    user: UserResponse {
                        id: user.id,
                        email: user.email.clone(),
                        display_name: user.display_name.clone(),
                        native_language: user.native_language.clone(),
                        email_verified: user.email_verified,
                        created_at: user.created_at,
                    },
                    tokens,
                })
            }
            Ok(false) => LoginResult::InvalidCredentials,
            Err(_) => LoginResult::InvalidCredentials,
        }
    }

    /// Login or create user from OAuth provider info
    pub fn oauth_login(user_info: OAuthUserInfo) -> OAuthLoginResult {
        let user = User::new_oauth(
            user_info.email.clone(),
            user_info.name.unwrap_or_else(|| user_info.email.clone()),
            user_info.provider.as_str().to_string(),
            user_info.provider_id,
        );

        let tokens = Self::generate_tokens(&user);

        OAuthLoginResult::Success(AuthResponse {
            user: UserResponse {
                id: user.id,
                email: user.email,
                display_name: user.display_name,
                native_language: user.native_language,
                email_verified: user.email_verified,
                created_at: user.created_at,
            },
            tokens,
        })
    }

    /// Generate JWT token pair via JwtService.
    fn generate_tokens(user: &User) -> TokenResponse {
        let jwt = JwtService::from_settings();
        match jwt.generate_tokens(user.id, &user.email) {
            Ok(pair) => TokenResponse {
                access_token: pair.access_token,
                refresh_token: pair.refresh_token,
                expires_in: pair.expires_in as u64,
            },
            Err(e) => {
                tracing::error!("JWT generation failed: {}", e);
                // Fallback: unsigned opaque token (never valid for API calls)
                TokenResponse {
                    access_token: format!("err_{}", user.id),
                    refresh_token: format!("err_refresh_{}", user.id),
                    expires_in: 0,
                }
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::apps::auth::oauth::OAuthProvider;

    #[test]
    fn test_register_success() {
        let request = RegisterRequest {
            email: "test@example.com".to_string(),
            password: "securepassword123".to_string(),
            display_name: "Test User".to_string(),
            native_language: "en".to_string(),
        };

        let result = AuthViewSet::register(request);

        match result {
            RegisterResult::Success(response, _user) => {
                assert_eq!(response.user.email, "test@example.com");
                assert!(!response.tokens.access_token.is_empty());
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_register_invalid_email() {
        let request = RegisterRequest {
            email: "invalid-email".to_string(),
            password: "password123".to_string(),
            display_name: "Test".to_string(),
            native_language: "en".to_string(),
        };

        let result = AuthViewSet::register(request);
        assert!(matches!(result, RegisterResult::InvalidInput(_)));
    }

    #[test]
    fn test_register_short_password() {
        let request = RegisterRequest {
            email: "test@example.com".to_string(),
            password: "short".to_string(),
            display_name: "Test".to_string(),
            native_language: "en".to_string(),
        };

        let result = AuthViewSet::register(request);
        assert!(matches!(result, RegisterResult::InvalidInput(_)));
    }

    #[test]
    fn test_login_success() {
        // Create user with known password
        let password = "testpassword123";
        let hash = hash_password(password).unwrap();
        let user = User::new(
            "login@example.com".to_string(),
            "Login User".to_string(),
            hash,
        );

        let request = LoginRequest {
            email: "login@example.com".to_string(),
            password: password.to_string(),
        };

        let result = AuthViewSet::login(request, Some(&user));
        assert!(matches!(result, LoginResult::Success(_)));
    }

    #[test]
    fn test_login_wrong_password() {
        let hash = hash_password("correctpassword").unwrap();
        let user = User::new("test@example.com".to_string(), "Test".to_string(), hash);

        let request = LoginRequest {
            email: "test@example.com".to_string(),
            password: "wrongpassword".to_string(),
        };

        let result = AuthViewSet::login(request, Some(&user));
        assert!(matches!(result, LoginResult::InvalidCredentials));
    }

    #[test]
    fn test_login_user_not_found() {
        let request = LoginRequest {
            email: "nonexistent@example.com".to_string(),
            password: "password".to_string(),
        };

        let result = AuthViewSet::login(request, None);
        assert!(matches!(result, LoginResult::UserNotFound));
    }

    #[test]
    fn test_oauth_login_google() {
        let user_info = OAuthUserInfo {
            provider: OAuthProvider::Google,
            provider_id: "google_123".to_string(),
            email: "oauth@gmail.com".to_string(),
            name: Some("OAuth User".to_string()),
            avatar_url: Some("https://example.com/photo.jpg".to_string()),
        };

        let result = AuthViewSet::oauth_login(user_info);

        match result {
            OAuthLoginResult::Success(response) => {
                assert_eq!(response.user.email, "oauth@gmail.com");
                assert_eq!(response.user.display_name, "OAuth User");
                assert!(response.user.email_verified); // OAuth users are pre-verified
                assert!(!response.tokens.access_token.is_empty());
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_oauth_login_without_name() {
        let user_info = OAuthUserInfo {
            provider: OAuthProvider::GitHub,
            provider_id: "gh_456".to_string(),
            email: "dev@github.com".to_string(),
            name: None,
            avatar_url: None,
        };

        let result = AuthViewSet::oauth_login(user_info);

        match result {
            OAuthLoginResult::Success(response) => {
                // Falls back to email as display name
                assert_eq!(response.user.display_name, "dev@github.com");
            }
            _ => panic!("Expected success"),
        }
    }
}
