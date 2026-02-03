//! API Router - Connects ViewSets to HTTP endpoints
//!
//! This module provides the HTTP routing layer that connects
//! the business logic in ViewSets to axum routes.

pub mod auth;
pub mod books;
pub mod learning;
pub mod middleware;
pub mod review;
