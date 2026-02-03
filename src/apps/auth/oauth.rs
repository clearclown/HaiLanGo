//! OAuth Provider Integration
//!
//! Supports Google and GitHub OAuth2 authentication flows.

use serde::{Deserialize, Serialize};
use thiserror::Error;

/// OAuth provider types
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum OAuthProvider {
    Google,
    GitHub,
}

impl OAuthProvider {
    /// Get the provider name as string
    pub fn as_str(&self) -> &'static str {
        match self {
            OAuthProvider::Google => "google",
            OAuthProvider::GitHub => "github",
        }
    }

    /// Parse from string
    pub fn parse(s: &str) -> Option<Self> {
        match s.to_lowercase().as_str() {
            "google" => Some(OAuthProvider::Google),
            "github" => Some(OAuthProvider::GitHub),
            _ => None,
        }
    }
}

/// OAuth configuration for a provider
#[derive(Debug, Clone)]
pub struct OAuthConfig {
    pub client_id: String,
    pub client_secret: String,
    pub redirect_uri: String,
    pub auth_url: String,
    pub token_url: String,
    pub userinfo_url: String,
    pub scopes: Vec<String>,
}

impl OAuthConfig {
    /// Create Google OAuth configuration
    pub fn google(client_id: &str, client_secret: &str, redirect_uri: &str) -> Self {
        Self {
            client_id: client_id.to_string(),
            client_secret: client_secret.to_string(),
            redirect_uri: redirect_uri.to_string(),
            auth_url: "https://accounts.google.com/o/oauth2/v2/auth".to_string(),
            token_url: "https://oauth2.googleapis.com/token".to_string(),
            userinfo_url: "https://www.googleapis.com/oauth2/v2/userinfo".to_string(),
            scopes: vec![
                "openid".to_string(),
                "email".to_string(),
                "profile".to_string(),
            ],
        }
    }

    /// Create GitHub OAuth configuration
    pub fn github(client_id: &str, client_secret: &str, redirect_uri: &str) -> Self {
        Self {
            client_id: client_id.to_string(),
            client_secret: client_secret.to_string(),
            redirect_uri: redirect_uri.to_string(),
            auth_url: "https://github.com/login/oauth/authorize".to_string(),
            token_url: "https://github.com/login/oauth/access_token".to_string(),
            userinfo_url: "https://api.github.com/user".to_string(),
            scopes: vec!["user:email".to_string()],
        }
    }

    /// Build the authorization URL
    pub fn authorization_url(&self, state: &str) -> String {
        let scopes = self.scopes.join(" ");
        format!(
            "{}?client_id={}&redirect_uri={}&scope={}&state={}&response_type=code",
            self.auth_url,
            urlencoding::encode(&self.client_id),
            urlencoding::encode(&self.redirect_uri),
            urlencoding::encode(&scopes),
            urlencoding::encode(state)
        )
    }
}

/// OAuth errors
#[derive(Debug, Error)]
pub enum OAuthError {
    #[error("Invalid OAuth provider: {0}")]
    InvalidProvider(String),
    #[error("OAuth token exchange failed: {0}")]
    TokenExchangeFailed(String),
    #[error("Failed to fetch user info: {0}")]
    UserInfoFailed(String),
    #[error("Missing required field: {0}")]
    MissingField(String),
    #[error("State mismatch")]
    StateMismatch,
    #[error("HTTP error: {0}")]
    HttpError(String),
}

/// OAuth token response
#[derive(Debug, Deserialize)]
pub struct OAuthTokenResponse {
    pub access_token: String,
    pub token_type: String,
    pub expires_in: Option<u64>,
    pub refresh_token: Option<String>,
    pub scope: Option<String>,
}

/// Google user info response
#[derive(Debug, Deserialize)]
pub struct GoogleUserInfo {
    pub id: String,
    pub email: String,
    pub verified_email: Option<bool>,
    pub name: Option<String>,
    pub picture: Option<String>,
}

/// GitHub user info response
#[derive(Debug, Deserialize)]
pub struct GitHubUserInfo {
    pub id: i64,
    pub login: String,
    pub email: Option<String>,
    pub name: Option<String>,
    pub avatar_url: Option<String>,
}

/// Normalized OAuth user info
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OAuthUserInfo {
    pub provider: OAuthProvider,
    pub provider_id: String,
    pub email: String,
    pub name: Option<String>,
    pub avatar_url: Option<String>,
}

/// OAuth service for handling authentication flows
pub struct OAuthService {
    google_config: Option<OAuthConfig>,
    github_config: Option<OAuthConfig>,
}

impl OAuthService {
    /// Create a new OAuth service
    pub fn new() -> Self {
        Self {
            google_config: None,
            github_config: None,
        }
    }

    /// Configure from environment variables
    pub fn from_env() -> Self {
        let mut service = Self::new();

        // Google OAuth
        if let (Ok(client_id), Ok(client_secret)) = (
            std::env::var("GOOGLE_CLIENT_ID"),
            std::env::var("GOOGLE_CLIENT_SECRET"),
        ) {
            let redirect_uri = std::env::var("GOOGLE_REDIRECT_URI")
                .unwrap_or_else(|_| "http://localhost:8080/auth/callback/google".to_string());
            service.google_config = Some(OAuthConfig::google(
                &client_id,
                &client_secret,
                &redirect_uri,
            ));
        }

        // GitHub OAuth
        if let (Ok(client_id), Ok(client_secret)) = (
            std::env::var("GITHUB_CLIENT_ID"),
            std::env::var("GITHUB_CLIENT_SECRET"),
        ) {
            let redirect_uri = std::env::var("GITHUB_REDIRECT_URI")
                .unwrap_or_else(|_| "http://localhost:8080/auth/callback/github".to_string());
            service.github_config = Some(OAuthConfig::github(
                &client_id,
                &client_secret,
                &redirect_uri,
            ));
        }

        service
    }

    /// Set Google OAuth configuration
    pub fn with_google(mut self, config: OAuthConfig) -> Self {
        self.google_config = Some(config);
        self
    }

    /// Set GitHub OAuth configuration
    pub fn with_github(mut self, config: OAuthConfig) -> Self {
        self.github_config = Some(config);
        self
    }

    /// Get configuration for a provider
    pub fn get_config(&self, provider: OAuthProvider) -> Option<&OAuthConfig> {
        match provider {
            OAuthProvider::Google => self.google_config.as_ref(),
            OAuthProvider::GitHub => self.github_config.as_ref(),
        }
    }

    /// Check if a provider is configured
    pub fn is_configured(&self, provider: OAuthProvider) -> bool {
        self.get_config(provider).is_some()
    }

    /// Get the authorization URL for a provider
    pub fn get_authorization_url(
        &self,
        provider: OAuthProvider,
        state: &str,
    ) -> Result<String, OAuthError> {
        let config = self.get_config(provider).ok_or_else(|| {
            OAuthError::InvalidProvider(format!("{:?} is not configured", provider))
        })?;
        Ok(config.authorization_url(state))
    }

    /// Get list of configured providers
    pub fn configured_providers(&self) -> Vec<OAuthProvider> {
        let mut providers = Vec::new();
        if self.google_config.is_some() {
            providers.push(OAuthProvider::Google);
        }
        if self.github_config.is_some() {
            providers.push(OAuthProvider::GitHub);
        }
        providers
    }
}

impl Default for OAuthService {
    fn default() -> Self {
        Self::new()
    }
}

/// URL encoding helper
mod urlencoding {
    pub fn encode(input: &str) -> String {
        let mut encoded = String::new();
        for byte in input.bytes() {
            match byte {
                b'A'..=b'Z' | b'a'..=b'z' | b'0'..=b'9' | b'-' | b'_' | b'.' | b'~' => {
                    encoded.push(byte as char);
                }
                _ => {
                    encoded.push_str(&format!("%{:02X}", byte));
                }
            }
        }
        encoded
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_oauth_provider_from_str() {
        assert_eq!(
            OAuthProvider::parse("google"),
            Some(OAuthProvider::Google)
        );
        assert_eq!(
            OAuthProvider::parse("GITHUB"),
            Some(OAuthProvider::GitHub)
        );
        assert_eq!(OAuthProvider::parse("invalid"), None);
    }

    #[test]
    fn test_oauth_provider_as_str() {
        assert_eq!(OAuthProvider::Google.as_str(), "google");
        assert_eq!(OAuthProvider::GitHub.as_str(), "github");
    }

    #[test]
    fn test_google_config() {
        let config = OAuthConfig::google("client_id", "secret", "http://localhost/callback");
        assert!(config.auth_url.contains("google"));
        assert!(config.scopes.contains(&"email".to_string()));
    }

    #[test]
    fn test_github_config() {
        let config = OAuthConfig::github("client_id", "secret", "http://localhost/callback");
        assert!(config.auth_url.contains("github"));
        assert!(config.scopes.contains(&"user:email".to_string()));
    }

    #[test]
    fn test_authorization_url() {
        let config = OAuthConfig::google("my_client_id", "secret", "http://localhost/callback");
        let url = config.authorization_url("random_state");

        assert!(url.contains("client_id=my_client_id"));
        assert!(url.contains("state=random_state"));
        assert!(url.contains("response_type=code"));
    }

    #[test]
    fn test_oauth_service_configuration() {
        let service = OAuthService::new()
            .with_google(OAuthConfig::google(
                "g_id",
                "g_secret",
                "http://localhost/g",
            ))
            .with_github(OAuthConfig::github(
                "gh_id",
                "gh_secret",
                "http://localhost/gh",
            ));

        assert!(service.is_configured(OAuthProvider::Google));
        assert!(service.is_configured(OAuthProvider::GitHub));
        assert_eq!(service.configured_providers().len(), 2);
    }

    #[test]
    fn test_get_authorization_url() {
        let service = OAuthService::new().with_google(OAuthConfig::google(
            "id",
            "secret",
            "http://localhost/callback",
        ));

        let url = service.get_authorization_url(OAuthProvider::Google, "state123");
        assert!(url.is_ok());
        assert!(url.unwrap().contains("state=state123"));

        let err = service.get_authorization_url(OAuthProvider::GitHub, "state");
        assert!(err.is_err());
    }

    #[test]
    fn test_url_encoding() {
        assert_eq!(urlencoding::encode("hello"), "hello");
        assert_eq!(urlencoding::encode("hello world"), "hello%20world");
        assert_eq!(urlencoding::encode("test@email.com"), "test%40email.com");
    }
}
