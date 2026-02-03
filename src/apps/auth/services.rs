//! Authentication services
//!
//! Provides password hashing, JWT token generation and validation.

use argon2::{
    Argon2,
    password_hash::{PasswordHash, PasswordHasher, PasswordVerifier, SaltString, rand_core::OsRng},
};
use chrono::{Duration, Utc};
use jsonwebtoken::{DecodingKey, EncodingKey, Header, Validation, decode, encode};
use serde::{Deserialize, Serialize};
use thiserror::Error;
use uuid::Uuid;

#[derive(Error, Debug)]
pub enum AuthError {
    #[error("Password hashing failed")]
    HashingError,
    #[error("Invalid password")]
    InvalidPassword,
    #[error("User not found")]
    UserNotFound,
    #[error("Email already exists")]
    EmailExists,
    #[error("Token generation failed")]
    TokenGenerationError,
    #[error("Invalid token")]
    InvalidToken,
    #[error("Token expired")]
    TokenExpired,
}

/// Hash a password using Argon2id
pub fn hash_password(password: &str) -> Result<String, AuthError> {
    let salt = SaltString::generate(&mut OsRng);
    let argon2 = Argon2::default();

    argon2
        .hash_password(password.as_bytes(), &salt)
        .map(|hash| hash.to_string())
        .map_err(|_| AuthError::HashingError)
}

/// Verify a password against a hash
pub fn verify_password(password: &str, hash: &str) -> Result<bool, AuthError> {
    let parsed_hash = PasswordHash::new(hash).map_err(|_| AuthError::InvalidPassword)?;

    Ok(Argon2::default()
        .verify_password(password.as_bytes(), &parsed_hash)
        .is_ok())
}

/// JWT token pair (access + refresh)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TokenPair {
    pub access_token: String,
    pub refresh_token: String,
    pub token_type: String,
    pub expires_in: i64,
}

/// JWT claims structure
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Claims {
    pub sub: String, // User ID
    pub email: String,
    pub exp: i64,           // Expiration time
    pub iat: i64,           // Issued at
    pub token_type: String, // "access" or "refresh"
}

/// JWT service for token management
pub struct JwtService {
    secret: String,
    access_expiry_hours: i64,
    refresh_expiry_days: i64,
}

impl JwtService {
    /// Create a new JWT service
    pub fn new(secret: &str, access_expiry_hours: i64, refresh_expiry_days: i64) -> Self {
        Self {
            secret: secret.to_string(),
            access_expiry_hours,
            refresh_expiry_days,
        }
    }

    /// Create from default settings
    pub fn from_settings() -> Self {
        // Use default values for now, will be loaded from settings in production
        Self::new("development-secret-key", 24, 30)
    }

    /// Generate a token pair for a user
    pub fn generate_tokens(&self, user_id: Uuid, email: &str) -> Result<TokenPair, AuthError> {
        let now = Utc::now();
        let access_exp = now + Duration::hours(self.access_expiry_hours);
        let refresh_exp = now + Duration::days(self.refresh_expiry_days);

        // Generate access token
        let access_claims = Claims {
            sub: user_id.to_string(),
            email: email.to_string(),
            exp: access_exp.timestamp(),
            iat: now.timestamp(),
            token_type: "access".to_string(),
        };

        let access_token = encode(
            &Header::default(),
            &access_claims,
            &EncodingKey::from_secret(self.secret.as_bytes()),
        )
        .map_err(|_| AuthError::TokenGenerationError)?;

        // Generate refresh token
        let refresh_claims = Claims {
            sub: user_id.to_string(),
            email: email.to_string(),
            exp: refresh_exp.timestamp(),
            iat: now.timestamp(),
            token_type: "refresh".to_string(),
        };

        let refresh_token = encode(
            &Header::default(),
            &refresh_claims,
            &EncodingKey::from_secret(self.secret.as_bytes()),
        )
        .map_err(|_| AuthError::TokenGenerationError)?;

        Ok(TokenPair {
            access_token,
            refresh_token,
            token_type: "Bearer".to_string(),
            expires_in: self.access_expiry_hours * 3600,
        })
    }

    /// Verify and decode an access token
    pub fn verify_token(&self, token: &str) -> Result<Claims, AuthError> {
        let token_data = decode::<Claims>(
            token,
            &DecodingKey::from_secret(self.secret.as_bytes()),
            &Validation::default(),
        )
        .map_err(|e| {
            if e.to_string().contains("ExpiredSignature") {
                AuthError::TokenExpired
            } else {
                AuthError::InvalidToken
            }
        })?;

        Ok(token_data.claims)
    }

    /// Extract user ID from token
    pub fn get_user_id(&self, token: &str) -> Result<Uuid, AuthError> {
        let claims = self.verify_token(token)?;
        Uuid::parse_str(&claims.sub).map_err(|_| AuthError::InvalidToken)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_password_hashing() {
        let password = "secure_password_123";
        let hash = hash_password(password).expect("Hashing should succeed");

        assert!(!hash.is_empty());
        assert!(hash.starts_with("$argon2"));
    }

    #[test]
    fn test_password_verification_success() {
        let password = "my_secret_password";
        let hash = hash_password(password).unwrap();

        let result = verify_password(password, &hash).unwrap();
        assert!(result);
    }

    #[test]
    fn test_password_verification_failure() {
        let password = "correct_password";
        let hash = hash_password(password).unwrap();

        let result = verify_password("wrong_password", &hash).unwrap();
        assert!(!result);
    }

    #[test]
    fn test_jwt_token_generation() {
        let jwt_service = JwtService::from_settings();
        let user_id = Uuid::new_v4();
        let email = "test@example.com";

        let tokens = jwt_service.generate_tokens(user_id, email).unwrap();

        assert!(!tokens.access_token.is_empty());
        assert!(!tokens.refresh_token.is_empty());
        assert_eq!(tokens.token_type, "Bearer");
    }

    #[test]
    fn test_jwt_token_verification() {
        let jwt_service = JwtService::from_settings();
        let user_id = Uuid::new_v4();
        let email = "test@example.com";

        let tokens = jwt_service.generate_tokens(user_id, email).unwrap();
        let claims = jwt_service.verify_token(&tokens.access_token).unwrap();

        assert_eq!(claims.sub, user_id.to_string());
        assert_eq!(claims.email, email);
        assert_eq!(claims.token_type, "access");
    }

    #[test]
    fn test_get_user_id_from_token() {
        let jwt_service = JwtService::from_settings();
        let user_id = Uuid::new_v4();
        let email = "test@example.com";

        let tokens = jwt_service.generate_tokens(user_id, email).unwrap();
        let extracted_id = jwt_service.get_user_id(&tokens.access_token).unwrap();

        assert_eq!(extracted_id, user_id);
    }

    #[test]
    fn test_invalid_token() {
        let jwt_service = JwtService::from_settings();

        let result = jwt_service.verify_token("invalid.token.here");
        assert!(matches!(result, Err(AuthError::InvalidToken)));
    }
}
