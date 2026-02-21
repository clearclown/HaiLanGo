//! API Router - Connects ViewSets to HTTP endpoints via Reinhardt Handlers
//!
//! Each sub-module exposes `routes() -> Vec<Route>` which are mounted
//! by the DefaultRouter in config/urls.rs.

pub mod auth;
pub mod books;
pub mod learning;
pub mod middleware;
pub mod review;
pub mod teacher;
pub mod tts;
