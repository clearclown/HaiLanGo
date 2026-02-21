//! User model with Reinhardt ORM patterns

use chrono::{DateTime, Utc};
use reinhardt::db::orm::{FieldSelector, Model, Timestamped};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Type-safe field selector for User model
#[derive(Clone)]
pub struct UserFields;

impl FieldSelector for UserFields {
    fn with_alias(self, _alias: &str) -> Self {
        self
    }
}

/// User account model
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    pub id: Uuid,
    pub email: String,
    pub password_hash: Option<String>,
    pub display_name: String,
    pub native_language: String,
    pub avatar_url: Option<String>,
    pub oauth_provider: Option<String>,
    pub oauth_id: Option<String>,
    pub email_verified: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub last_login_at: Option<DateTime<Utc>>,
}

impl User {
    /// Create a new user with email and password
    pub fn new(email: String, display_name: String, password_hash: String) -> Self {
        let now = Utc::now();
        Self {
            id: Uuid::new_v4(),
            email,
            password_hash: Some(password_hash),
            display_name,
            native_language: "en".to_string(),
            avatar_url: None,
            oauth_provider: None,
            oauth_id: None,
            email_verified: false,
            created_at: now,
            updated_at: now,
            last_login_at: None,
        }
    }

    /// Create a new OAuth user
    pub fn new_oauth(
        email: String,
        display_name: String,
        provider: String,
        oauth_id: String,
    ) -> Self {
        let now = Utc::now();
        Self {
            id: Uuid::new_v4(),
            email,
            password_hash: None,
            display_name,
            native_language: "en".to_string(),
            avatar_url: None,
            oauth_provider: Some(provider),
            oauth_id: Some(oauth_id),
            email_verified: true, // OAuth users are pre-verified
            created_at: now,
            updated_at: now,
            last_login_at: None,
        }
    }
}

impl Model for User {
    type PrimaryKey = Uuid;
    type Fields = UserFields;

    fn table_name() -> &'static str {
        "users"
    }

    fn app_label() -> &'static str {
        "auth"
    }

    fn new_fields() -> Self::Fields {
        UserFields
    }

    fn primary_key(&self) -> Option<Self::PrimaryKey> {
        Some(self.id)
    }

    fn set_primary_key(&mut self, value: Self::PrimaryKey) {
        self.id = value;
    }
}

impl Timestamped for User {
    fn created_at(&self) -> DateTime<Utc> {
        self.created_at
    }

    fn updated_at(&self) -> DateTime<Utc> {
        self.updated_at
    }

    fn set_updated_at(&mut self, time: DateTime<Utc>) {
        self.updated_at = time;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_create_user_with_password() {
        let user = User::new(
            "test@example.com".to_string(),
            "Test User".to_string(),
            "hashed_password".to_string(),
        );

        assert_eq!(user.email, "test@example.com");
        assert_eq!(user.display_name, "Test User");
        assert!(user.password_hash.is_some());
        assert!(!user.email_verified);
    }

    #[test]
    fn test_create_oauth_user() {
        let user = User::new_oauth(
            "oauth@example.com".to_string(),
            "OAuth User".to_string(),
            "google".to_string(),
            "google_123".to_string(),
        );

        assert_eq!(user.oauth_provider, Some("google".to_string()));
        assert!(user.email_verified); // OAuth users are verified
        assert!(user.password_hash.is_none());
    }

    #[test]
    fn test_user_defaults() {
        let user = User::new(
            "defaults@example.com".to_string(),
            "Default User".to_string(),
            "hash".to_string(),
        );

        assert_eq!(user.native_language, "en");
        assert!(user.avatar_url.is_none());
        assert!(user.oauth_provider.is_none());
        assert!(user.oauth_id.is_none());
        assert!(user.last_login_at.is_none());
    }
}
