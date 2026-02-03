//! Application settings and configuration
//!
//! Loads configuration from TOML files and environment variables.
//! Uses Reinhardt's settings builder for composable configuration.

use anyhow::Context;
use reinhardt::conf::settings::builder::SettingsBuilder;
use reinhardt::conf::settings::sources::{DefaultSource, EnvSource, TomlFileSource};
use serde::Deserialize;
use std::env;
use std::sync::OnceLock;

/// Global settings instance
static SETTINGS: OnceLock<Settings> = OnceLock::new();

/// Application settings structure
#[derive(Debug, Clone, Deserialize)]
pub struct Settings {
    pub app: AppSettings,
    pub server: ServerSettings,
    pub database: DatabaseSettings,
    pub redis: RedisSettings,
    pub auth: AuthSettings,
    pub security: SecuritySettings,
    pub external_apis: ExternalApiSettings,
    pub storage: StorageSettings,
    pub logging: LoggingSettings,
}

#[derive(Debug, Clone, Deserialize)]
pub struct AppSettings {
    pub name: String,
    pub version: String,
    pub debug: bool,
}

#[derive(Debug, Clone, Deserialize)]
pub struct ServerSettings {
    pub host: String,
    pub port: u16,
    pub workers: usize,
}

#[derive(Debug, Clone, Deserialize)]
pub struct DatabaseSettings {
    pub engine: String,
    pub host: String,
    pub port: u16,
    pub name: String,
    pub user: String,
    pub password: String,
    pub pool_size: u32,
    pub max_connections: u32,
    pub timeout: u32,
}

impl DatabaseSettings {
    /// Generate database URL from settings
    pub fn url(&self) -> String {
        format!(
            "{}://{}:{}@{}:{}/{}",
            self.engine, self.user, self.password, self.host, self.port, self.name
        )
    }
}

#[derive(Debug, Clone, Deserialize)]
pub struct RedisSettings {
    pub host: String,
    pub port: u16,
    pub database: u8,
    pub pool_size: u32,
}

impl RedisSettings {
    /// Generate Redis URL from settings
    pub fn url(&self) -> String {
        format!("redis://{}:{}/{}", self.host, self.port, self.database)
    }
}

#[derive(Debug, Clone, Deserialize)]
pub struct AuthSettings {
    pub jwt_secret: String,
    pub jwt_expiry_hours: u64,
    pub refresh_expiry_days: u64,
    pub session_timeout_minutes: u64,
}

#[derive(Debug, Clone, Deserialize)]
pub struct SecuritySettings {
    pub cors_origins: Vec<String>,
    pub cors_allow_credentials: bool,
    pub rate_limit_requests: u32,
    pub rate_limit_window_seconds: u64,
}

#[derive(Debug, Clone, Deserialize)]
pub struct ExternalApiSettings {
    pub ocr_provider: String,
    pub google_vision_api_key: String,
    pub azure_vision_endpoint: String,
    pub azure_vision_key: String,
    pub tts_provider: String,
    pub google_tts_api_key: String,
    pub azure_speech_endpoint: String,
    pub azure_speech_key: String,
    pub stt_provider: String,
    pub openai_api_key: String,
    pub llm_provider: String,
    pub anthropic_api_key: String,
}

#[derive(Debug, Clone, Deserialize)]
pub struct StorageSettings {
    pub upload_path: String,
    pub max_file_size_mb: u64,
    pub allowed_extensions: Vec<String>,
}

#[derive(Debug, Clone, Deserialize)]
pub struct LoggingSettings {
    pub level: String,
    pub format: String,
}

impl Settings {
    /// Load settings from TOML files and environment
    pub fn load() -> anyhow::Result<Self> {
        let profile = env::var("APP_ENV").unwrap_or_else(|_| "local".to_string());

        let settings: Settings = SettingsBuilder::new()
            .add_source(DefaultSource::new())
            .add_source(TomlFileSource::new("settings/base.toml"))
            .add_source(TomlFileSource::new(format!("settings/{}.toml", profile)))
            .add_source(EnvSource::new().with_prefix("HAILANGO_"))
            .build()
            .context("Failed to build settings")?
            .into_typed()
            .context("Failed to parse settings")?;

        Ok(settings)
    }

    /// Get or initialize global settings
    pub fn get() -> &'static Settings {
        SETTINGS.get_or_init(|| Self::load().expect("Failed to load settings"))
    }

    /// Check if running in development mode
    pub fn is_development(&self) -> bool {
        self.app.debug
    }

    /// Check if running in production mode
    pub fn is_production(&self) -> bool {
        !self.app.debug
    }

    /// Legacy method for backward compatibility
    pub fn from_env() -> anyhow::Result<LegacySettings> {
        let _ = dotenvy::dotenv();

        Ok(LegacySettings {
            app_env: env::var("APP_ENV").unwrap_or_else(|_| "development".to_string()),
            database_url: env::var("DATABASE_URL").unwrap_or_else(|_| {
                "postgresql://postgres:postgres@localhost:5432/hailango".to_string()
            }),
            redis_url: env::var("REDIS_URL")
                .unwrap_or_else(|_| "redis://localhost:6379".to_string()),
            jwt_secret: env::var("JWT_SECRET")
                .unwrap_or_else(|_| "development-secret-key".to_string()),
        })
    }
}

/// Legacy settings structure for backward compatibility
#[derive(Debug, Clone)]
pub struct LegacySettings {
    pub app_env: String,
    pub database_url: String,
    pub redis_url: String,
    pub jwt_secret: String,
}

impl LegacySettings {
    pub fn is_development(&self) -> bool {
        self.app_env == "development"
    }

    pub fn is_testing(&self) -> bool {
        self.app_env == "testing"
    }

    pub fn is_production(&self) -> bool {
        self.app_env == "production"
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_is_development() {
        let settings = LegacySettings {
            app_env: "development".to_string(),
            database_url: String::new(),
            redis_url: String::new(),
            jwt_secret: String::new(),
        };
        assert!(settings.is_development());
        assert!(!settings.is_production());
    }

    #[test]
    fn test_is_production() {
        let settings = LegacySettings {
            app_env: "production".to_string(),
            database_url: String::new(),
            redis_url: String::new(),
            jwt_secret: String::new(),
        };
        assert!(settings.is_production());
        assert!(!settings.is_development());
    }
}
