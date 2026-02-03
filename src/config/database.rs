//! Database configuration and connection pool
//!
//! Uses SQLx for PostgreSQL connection with async support.

use sqlx::{PgPool, postgres::PgPoolOptions};
use std::sync::Arc;
use std::time::Duration;
use tokio::sync::OnceCell;

/// Global database connection pool
static DB_POOL: OnceCell<Arc<PgPool>> = OnceCell::const_new();

/// Database connection configuration
#[derive(Debug, Clone)]
pub struct DbConfig {
    pub url: String,
    pub max_connections: u32,
    pub min_connections: u32,
    pub connect_timeout_secs: u64,
    pub idle_timeout_secs: u64,
}

impl Default for DbConfig {
    fn default() -> Self {
        Self {
            url: "postgresql://postgres:postgres@localhost:5432/hailango".to_string(),
            max_connections: 10,
            min_connections: 2,
            connect_timeout_secs: 30,
            idle_timeout_secs: 600,
        }
    }
}

impl DbConfig {
    /// Create config from environment
    pub fn from_env() -> Self {
        let url = std::env::var("DATABASE_URL").unwrap_or_else(|_| {
            "postgresql://postgres:postgres@localhost:5432/hailango".to_string()
        });

        Self {
            url,
            ..Default::default()
        }
    }
}

/// Initialize database connection pool
pub async fn init_db(config: &DbConfig) -> Result<Arc<PgPool>, sqlx::Error> {
    let pool = PgPoolOptions::new()
        .max_connections(config.max_connections)
        .min_connections(config.min_connections)
        .acquire_timeout(Duration::from_secs(config.connect_timeout_secs))
        .idle_timeout(Duration::from_secs(config.idle_timeout_secs))
        .connect(&config.url)
        .await?;

    Ok(Arc::new(pool))
}

/// Get or initialize the global database connection
pub async fn get_db() -> Result<Arc<PgPool>, sqlx::Error> {
    DB_POOL
        .get_or_try_init(|| async {
            let config = DbConfig::from_env();
            init_db(&config).await
        })
        .await
        .map(Arc::clone)
}

/// Health check for database connection
pub async fn check_db_health(pool: &PgPool) -> bool {
    sqlx::query("SELECT 1").fetch_one(pool).await.is_ok()
}

/// Repository trait for database operations
#[async_trait::async_trait]
pub trait Repository<T, ID> {
    async fn find_by_id(&self, id: ID) -> Result<Option<T>, sqlx::Error>;
    async fn find_all(&self) -> Result<Vec<T>, sqlx::Error>;
    async fn create(&self, entity: &T) -> Result<T, sqlx::Error>;
    async fn update(&self, entity: &T) -> Result<T, sqlx::Error>;
    async fn delete(&self, id: ID) -> Result<bool, sqlx::Error>;
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_config() {
        let config = DbConfig::default();
        assert!(config.url.contains("postgresql://"));
        assert_eq!(config.max_connections, 10);
        assert_eq!(config.min_connections, 2);
    }

    #[test]
    fn test_config_from_env() {
        // SAFETY: This test is run in isolation and modifies environment variables
        unsafe {
            std::env::set_var(
                "DATABASE_URL",
                "postgresql://test:test@localhost:5432/test_db",
            );
        }

        let config = DbConfig::from_env();
        assert!(config.url.contains("test_db"));

        // SAFETY: Clean up the environment variable
        unsafe {
            std::env::remove_var("DATABASE_URL");
        }
    }

    #[test]
    fn test_config_values() {
        let config = DbConfig {
            url: "postgresql://user:pass@host:5432/db".to_string(),
            max_connections: 20,
            min_connections: 5,
            connect_timeout_secs: 60,
            idle_timeout_secs: 300,
        };

        assert_eq!(config.max_connections, 20);
        assert_eq!(config.min_connections, 5);
        assert_eq!(config.connect_timeout_secs, 60);
    }
}
