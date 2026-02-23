//! Auth API routes

use async_trait::async_trait;
use serde_json::json;
use std::collections::{HashMap, HashSet};
use std::sync::{Arc, RwLock};
use uuid::Uuid;

use crate::apps::auth::{
    dto::{LoginRequest, OAuthCallbackQuery, RegisterRequest},
    models::User,
    oauth::{OAuthProvider, OAuthService},
    views::{AuthViewSet, LoginResult, OAuthLoginResult, RegisterResult},
};
use crate::{Handler, Request, Response, Result, Route, StatusCode, path};

/// OAuth-enabled auth state, also serving as the in-memory user store.
#[derive(Clone)]
pub struct AuthState {
    pub oauth_service: Arc<OAuthService>,
    /// Pending OAuth state tokens for CSRF protection.
    /// Populated by OAuthRedirectHandler, consumed by OAuthCallbackHandler.
    pub pending_states: Arc<RwLock<HashSet<String>>>,
    /// In-memory user store keyed by email address.
    /// RegisterHandler writes; LoginHandler reads.
    pub users: Arc<RwLock<HashMap<String, User>>>,
}

impl Default for AuthState {
    fn default() -> Self {
        Self {
            oauth_service: Arc::new(OAuthService::from_env()),
            pending_states: Arc::new(RwLock::new(HashSet::new())),
            users: Arc::new(RwLock::new(HashMap::new())),
        }
    }
}

/// Handler for POST /register
struct RegisterHandler {
    state: AuthState,
}

#[async_trait]
impl Handler for RegisterHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let req: RegisterRequest = request.json()?;
        let email = req.email.clone();

        // Check for duplicate email before registration
        if let Ok(users) = self.state.users.read() {
            if users.contains_key(&email) {
                return Response::new(StatusCode::CONFLICT)
                    .with_json(&json!({"error": "Email already exists"}));
            }
        }

        match AuthViewSet::register(req) {
            RegisterResult::Success(response, user) => {
                // Persist user in the shared store for subsequent logins
                if let Ok(mut users) = self.state.users.write() {
                    users.insert(user.email.clone(), *user);
                }
                Response::created().with_json(&*response)
            }
            RegisterResult::EmailExists => Response::new(StatusCode::CONFLICT)
                .with_json(&json!({"error": "Email already exists"})),
            RegisterResult::InvalidInput(msg) => {
                Response::bad_request().with_json(&json!({"error": msg}))
            }
        }
    }
}

/// Handler for POST /login
struct LoginHandler {
    state: AuthState,
}

#[async_trait]
impl Handler for LoginHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let req: LoginRequest = request.json()?;

        // Look up user by email from the shared store
        let stored_user = self
            .state
            .users
            .read()
            .ok()
            .and_then(|users| users.get(&req.email).cloned());

        match AuthViewSet::login(req, stored_user.as_ref()) {
            LoginResult::Success(response) => Response::ok().with_json(&response),
            LoginResult::InvalidCredentials => {
                Response::unauthorized().with_json(&json!({"error": "Invalid credentials"}))
            }
            LoginResult::UserNotFound => {
                Response::not_found().with_json(&json!({"error": "User not found"}))
            }
        }
    }
}

/// Handler for GET /oauth/providers
struct OAuthProvidersHandler {
    state: AuthState,
}

#[async_trait]
impl Handler for OAuthProvidersHandler {
    async fn handle(&self, _request: Request) -> Result<Response> {
        let providers: Vec<serde_json::Value> = self
            .state
            .oauth_service
            .configured_providers()
            .iter()
            .map(|p| json!({"name": p.as_str(), "configured": true}))
            .collect();
        Response::ok().with_json(&json!({"providers": providers}))
    }
}

/// Handler for GET /oauth/{provider}
struct OAuthRedirectHandler {
    state: AuthState,
}

#[async_trait]
impl Handler for OAuthRedirectHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let provider_name = request
            .path_params
            .get("provider")
            .ok_or_else(|| crate::Error::Validation("Missing provider parameter".into()))?
            .clone();

        let provider = match OAuthProvider::parse(&provider_name) {
            Some(p) => p,
            None => {
                return Response::bad_request()
                    .with_json(&json!({"error": format!("Unknown provider: {}", provider_name)}));
            }
        };

        let state_token = Uuid::new_v4().to_string();

        // Store state token for CSRF verification in the callback
        if let Ok(mut states) = self.state.pending_states.write() {
            states.insert(state_token.clone());
        }

        match self
            .state
            .oauth_service
            .get_authorization_url(provider, &state_token)
        {
            Ok(url) => Response::ok().with_json(&json!({"auth_url": url, "state": state_token})),
            Err(e) => Response::bad_request().with_json(&json!({"error": e.to_string()})),
        }
    }
}

/// Handler for GET /callback/{provider}
struct OAuthCallbackHandler {
    state: AuthState,
}

#[async_trait]
impl Handler for OAuthCallbackHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let provider_name = request
            .path_params
            .get("provider")
            .ok_or_else(|| crate::Error::Validation("Missing provider parameter".into()))?
            .clone();

        let provider = match OAuthProvider::parse(&provider_name) {
            Some(p) => p,
            None => {
                return Response::bad_request().with_json(&json!({"error": "Unknown provider"}));
            }
        };

        let query: OAuthCallbackQuery = request.query_as()?;

        // Verify and consume the state token to prevent CSRF attacks
        let state_valid = self
            .state
            .pending_states
            .write()
            .ok()
            .map(|mut states| states.remove(&query.state))
            .unwrap_or(false);

        if !state_valid {
            return Response::new(StatusCode::FORBIDDEN)
                .with_json(&json!({"error": "Invalid or expired OAuth state parameter"}));
        }

        match self
            .state
            .oauth_service
            .authenticate(provider, &query.code)
            .await
        {
            Ok(user_info) => match AuthViewSet::oauth_login(user_info) {
                OAuthLoginResult::Success(response) => Response::ok().with_json(&response),
                OAuthLoginResult::ProviderError(msg) => {
                    Response::internal_server_error().with_json(&json!({"error": msg}))
                }
            },
            Err(e) => Response::bad_request().with_json(&json!({"error": e.to_string()})),
        }
    }
}

/// Create auth routes
pub fn routes() -> Vec<Route> {
    let state = AuthState::default();

    vec![
        path(
            "/register/",
            RegisterHandler {
                state: state.clone(),
            },
        ),
        path(
            "/login/",
            LoginHandler {
                state: state.clone(),
            },
        ),
        path(
            "/oauth/providers/",
            OAuthProvidersHandler {
                state: state.clone(),
            },
        ),
        path(
            "/oauth/{provider}/",
            OAuthRedirectHandler {
                state: state.clone(),
            },
        ),
        path("/callback/{provider}/", OAuthCallbackHandler { state }),
    ]
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::Method;
    use bytes::Bytes;

    fn result_to_response(result: Result<Response>) -> Response {
        match result {
            Ok(r) => r,
            Err(e) => Response::from(e),
        }
    }

    fn build_auth_handlers() -> (RegisterHandler, LoginHandler) {
        let state = AuthState::default();
        (
            RegisterHandler {
                state: state.clone(),
            },
            LoginHandler { state },
        )
    }

    #[tokio::test]
    async fn test_register_endpoint() {
        let (handler, _) = build_auth_handlers();

        let body =
            r#"{"email":"test@example.com","password":"password123","display_name":"Test User"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/register/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 201);
    }

    #[tokio::test]
    async fn test_register_invalid_email() {
        let (handler, _) = build_auth_handlers();

        let body = r#"{"email":"invalid","password":"password123","display_name":"Test"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/register/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 400);
    }

    #[tokio::test]
    async fn test_login_user_not_found() {
        let (_, handler) = build_auth_handlers();

        let body = r#"{"email":"notfound@example.com","password":"password123"}"#;

        let request = Request::builder()
            .method(Method::POST)
            .uri("/login/")
            .header("content-type", "application/json")
            .body(Bytes::from(body))
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 404);
    }

    #[tokio::test]
    async fn test_login_after_register() {
        let (reg_handler, login_handler) = build_auth_handlers();

        // Register
        let reg_body =
            r#"{"email":"flow@example.com","password":"password123","display_name":"Flow User"}"#;
        let reg_request = Request::builder()
            .method(Method::POST)
            .uri("/register/")
            .header("content-type", "application/json")
            .body(Bytes::from(reg_body))
            .build()
            .unwrap();
        let reg_response = result_to_response(reg_handler.handle(reg_request).await);
        assert_eq!(reg_response.status, 201);

        // Login with same credentials
        let login_body = r#"{"email":"flow@example.com","password":"password123"}"#;
        let login_request = Request::builder()
            .method(Method::POST)
            .uri("/login/")
            .header("content-type", "application/json")
            .body(Bytes::from(login_body))
            .build()
            .unwrap();
        let login_response = result_to_response(login_handler.handle(login_request).await);
        assert_eq!(login_response.status, 200);
    }

    #[tokio::test]
    async fn test_login_wrong_password() {
        let (reg_handler, login_handler) = build_auth_handlers();

        // Register
        let reg_body =
            r#"{"email":"wrong@example.com","password":"correctpassword","display_name":"User"}"#;
        let reg_request = Request::builder()
            .method(Method::POST)
            .uri("/register/")
            .header("content-type", "application/json")
            .body(Bytes::from(reg_body))
            .build()
            .unwrap();
        result_to_response(reg_handler.handle(reg_request).await);

        // Login with wrong password
        let login_body = r#"{"email":"wrong@example.com","password":"wrongpassword"}"#;
        let login_request = Request::builder()
            .method(Method::POST)
            .uri("/login/")
            .header("content-type", "application/json")
            .body(Bytes::from(login_body))
            .build()
            .unwrap();
        let login_response = result_to_response(login_handler.handle(login_request).await);
        assert_eq!(login_response.status, 401);
    }

    #[tokio::test]
    async fn test_oauth_providers_list() {
        let state = AuthState::default();
        let handler = OAuthProvidersHandler { state };

        let request = Request::builder()
            .method(Method::GET)
            .uri("/oauth/providers/")
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_oauth_redirect_unknown_provider() {
        let state = AuthState::default();
        let handler = OAuthRedirectHandler { state };

        let mut params = std::collections::HashMap::new();
        params.insert("provider".to_string(), "unknown".to_string());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/oauth/unknown/")
            .path_params(params)
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 400);
    }

    #[tokio::test]
    async fn test_oauth_redirect_unconfigured_google() {
        let state = AuthState::default();
        let handler = OAuthRedirectHandler { state };

        let mut params = std::collections::HashMap::new();
        params.insert("provider".to_string(), "google".to_string());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/oauth/google/")
            .path_params(params)
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 400);
    }

    #[tokio::test]
    async fn test_oauth_callback_unknown_provider() {
        let state = AuthState::default();
        let handler = OAuthCallbackHandler { state };

        let mut params = std::collections::HashMap::new();
        params.insert("provider".to_string(), "invalid".to_string());

        let request = Request::builder()
            .method(Method::GET)
            .uri("/callback/invalid/?code=abc&state=xyz")
            .path_params(params)
            .body(Bytes::new())
            .build()
            .unwrap();

        let response = result_to_response(handler.handle(request).await);
        assert_eq!(response.status, 400);
    }
}
