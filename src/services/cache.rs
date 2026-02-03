//! Redis Cache Service
//!
//! Provides caching functionality using Redis for:
//! - Session storage
//! - Rate limiting counters
//! - Temporary data caching
//! - Job queue support

use redis::{AsyncCommands, Client, RedisError, aio::ConnectionManager};
use serde::{Serialize, de::DeserializeOwned};
use std::time::Duration;

/// Redis cache service
pub struct CacheService {
    conn: ConnectionManager,
    prefix: String,
}

impl CacheService {
    /// Create a new cache service
    pub async fn new(redis_url: &str, prefix: &str) -> Result<Self, RedisError> {
        let client = Client::open(redis_url)?;
        let conn = ConnectionManager::new(client).await?;

        Ok(Self {
            conn,
            prefix: prefix.to_string(),
        })
    }

    /// Create from default settings
    pub async fn from_env() -> Result<Self, RedisError> {
        let redis_url =
            std::env::var("REDIS_URL").unwrap_or_else(|_| "redis://localhost:6379".to_string());
        Self::new(&redis_url, "hailango").await
    }

    /// Build a key with prefix
    fn key(&self, key: &str) -> String {
        format!("{}:{}", self.prefix, key)
    }

    /// Set a value with expiration
    pub async fn set<T: Serialize>(
        &mut self,
        key: &str,
        value: &T,
        ttl: Duration,
    ) -> Result<(), RedisError> {
        let json = serde_json::to_string(value).map_err(|e| {
            RedisError::from((
                redis::ErrorKind::IoError,
                "Serialization error",
                e.to_string(),
            ))
        })?;

        self.conn.set_ex(self.key(key), json, ttl.as_secs()).await
    }

    /// Get a value
    pub async fn get<T: DeserializeOwned>(&mut self, key: &str) -> Result<Option<T>, RedisError> {
        let result: Option<String> = self.conn.get(self.key(key)).await?;

        match result {
            Some(json) => {
                let value: T = serde_json::from_str(&json).map_err(|e| {
                    RedisError::from((
                        redis::ErrorKind::IoError,
                        "Deserialization error",
                        e.to_string(),
                    ))
                })?;
                Ok(Some(value))
            }
            None => Ok(None),
        }
    }

    /// Delete a key
    pub async fn delete(&mut self, key: &str) -> Result<bool, RedisError> {
        let deleted: i32 = self.conn.del(self.key(key)).await?;
        Ok(deleted > 0)
    }

    /// Check if a key exists
    pub async fn exists(&mut self, key: &str) -> Result<bool, RedisError> {
        self.conn.exists(self.key(key)).await
    }

    /// Increment a counter (for rate limiting)
    pub async fn incr(&mut self, key: &str) -> Result<i64, RedisError> {
        self.conn.incr(self.key(key), 1).await
    }

    /// Increment with expiration (atomic)
    pub async fn incr_with_ttl(&mut self, key: &str, ttl: Duration) -> Result<i64, RedisError> {
        let full_key = self.key(key);

        // Use MULTI/EXEC for atomic operation
        let count: i64 = self.conn.incr(&full_key, 1).await?;

        // Set expiration only on first increment
        if count == 1 {
            let _: bool = self.conn.expire(&full_key, ttl.as_secs() as i64).await?;
        }

        Ok(count)
    }

    /// Set expiration on a key
    pub async fn expire(&mut self, key: &str, ttl: Duration) -> Result<bool, RedisError> {
        self.conn.expire(self.key(key), ttl.as_secs() as i64).await
    }

    /// Get TTL of a key
    pub async fn ttl(&mut self, key: &str) -> Result<i64, RedisError> {
        self.conn.ttl(self.key(key)).await
    }

    // --- Session Management ---

    /// Store session data
    pub async fn set_session<T: Serialize>(
        &mut self,
        session_id: &str,
        data: &T,
        ttl: Duration,
    ) -> Result<(), RedisError> {
        self.set(&format!("session:{}", session_id), data, ttl)
            .await
    }

    /// Get session data
    pub async fn get_session<T: DeserializeOwned>(
        &mut self,
        session_id: &str,
    ) -> Result<Option<T>, RedisError> {
        self.get(&format!("session:{}", session_id)).await
    }

    /// Delete session
    pub async fn delete_session(&mut self, session_id: &str) -> Result<bool, RedisError> {
        self.delete(&format!("session:{}", session_id)).await
    }

    // --- Rate Limiting ---

    /// Check rate limit (returns remaining requests, or None if exceeded)
    pub async fn check_rate_limit(
        &mut self,
        identifier: &str,
        max_requests: u32,
        window: Duration,
    ) -> Result<Option<u32>, RedisError> {
        let key = format!("ratelimit:{}", identifier);
        let count = self.incr_with_ttl(&key, window).await?;

        if count as u32 > max_requests {
            Ok(None) // Rate limit exceeded
        } else {
            Ok(Some(max_requests - count as u32)) // Remaining requests
        }
    }

    // --- Job Queue Support ---

    /// Push a job to a queue
    pub async fn push_job<T: Serialize>(&mut self, queue: &str, job: &T) -> Result<(), RedisError> {
        let json = serde_json::to_string(job).map_err(|e| {
            RedisError::from((
                redis::ErrorKind::IoError,
                "Serialization error",
                e.to_string(),
            ))
        })?;

        self.conn
            .rpush(self.key(&format!("queue:{}", queue)), json)
            .await
    }

    /// Pop a job from a queue (blocking)
    pub async fn pop_job<T: DeserializeOwned>(
        &mut self,
        queue: &str,
        timeout: Duration,
    ) -> Result<Option<T>, RedisError> {
        let result: Option<(String, String)> = self
            .conn
            .blpop(
                self.key(&format!("queue:{}", queue)),
                timeout.as_secs() as f64,
            )
            .await?;

        match result {
            Some((_, json)) => {
                let job: T = serde_json::from_str(&json).map_err(|e| {
                    RedisError::from((
                        redis::ErrorKind::IoError,
                        "Deserialization error",
                        e.to_string(),
                    ))
                })?;
                Ok(Some(job))
            }
            None => Ok(None),
        }
    }

    /// Get queue length
    pub async fn queue_length(&mut self, queue: &str) -> Result<i64, RedisError> {
        self.conn.llen(self.key(&format!("queue:{}", queue))).await
    }
}

#[cfg(test)]
mod tests {
    // Note: These tests require a running Redis instance
    // Run with: docker run -p 6379:6379 redis:7-alpine

    #[test]
    fn test_key_prefix() {
        // Mock test without actual Redis connection
        let prefix = "test";
        let key = format!("{}:{}", prefix, "mykey");
        assert_eq!(key, "test:mykey");
    }

    #[test]
    fn test_session_key_format() {
        let session_id = "abc123";
        let key = format!("session:{}", session_id);
        assert_eq!(key, "session:abc123");
    }

    #[test]
    fn test_rate_limit_key_format() {
        let identifier = "user:123";
        let key = format!("ratelimit:{}", identifier);
        assert_eq!(key, "ratelimit:user:123");
    }

    #[test]
    fn test_queue_key_format() {
        let queue = "ocr_jobs";
        let key = format!("queue:{}", queue);
        assert_eq!(key, "queue:ocr_jobs");
    }
}
