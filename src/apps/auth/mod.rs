//! Authentication app
//!
//! Handles user registration, login, JWT token generation and validation,
//! password management, and OAuth integration.
//!
//! Models:
//! - User: Core user account with authentication credentials
//!
//! Services:
//! - Password hashing with Argon2id
//! - Password verification
//!
//! DTOs:
//! - RegisterRequest, LoginRequest for request payloads
//! - AuthResponse, UserResponse, TokenResponse for responses
//!
//! Views:
//! - AuthViewSet: Registration and Login endpoints

pub mod dto;
pub mod models;
pub mod services;
pub mod views;

pub use dto::*;
pub use models::User;
pub use services::{AuthError, hash_password, verify_password};
pub use views::{AuthViewSet, LoginResult, RegisterResult};
