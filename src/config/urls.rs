//! URL routing configuration
//!
//! Defines all API routes and endpoint mappings for the application.
//! Uses Reinhardt's DefaultRouter with mount() for prefix-based grouping.

use crate::api;
use crate::{DefaultRouter, Router};

/// Build the application router with all API routes mounted
pub fn configure_urls() -> DefaultRouter {
    let mut router = DefaultRouter::new();

    router.mount("/api/auth", api::auth::routes(), Some("auth".to_string()));
    router.mount(
        "/api/books",
        api::books::routes(),
        Some("books".to_string()),
    );
    router.mount(
        "/api/learning",
        api::learning::routes(),
        Some("learning".to_string()),
    );
    router.mount(
        "/api/review",
        api::review::routes(),
        Some("review".to_string()),
    );
    router.mount("/api/tts", api::tts::routes(), Some("tts".to_string()));
    router.mount(
        "/api/teacher",
        api::teacher::routes(),
        Some("teacher".to_string()),
    );

    router
}
