//! HaiLanGo - AI Language Learning Platform
//!
//! This library provides the core backend functionality for the HaiLanGo
//! language learning platform, built with Reinhardt (Rust full-stack framework).

// Re-export Reinhardt prelude for convenient access
pub use reinhardt::prelude::*;

// Re-export types not included in prelude but needed by API handlers
pub use reinhardt::Method;
pub use reinhardt::path;

pub mod api;
pub mod apps;
pub mod config;
pub mod services;

// Frontend pages (WASM only)
#[cfg(target_arch = "wasm32")]
pub mod pages;
