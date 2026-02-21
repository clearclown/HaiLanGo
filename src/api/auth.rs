//! Auth API routes

use async_trait::async_trait;
use serde_json::json;
use std::sync::Arc;
use uuid::Uuid;

use crate::apps::auth::{
    dto::{LoginRequest, OAuthCallbackQuery, RegisterRequest},
    oauth::{OAuthProvider, OAuthService},
    views::{AuthViewSet, LoginResult, OAuthLoginResult, RegisterResult},
};
use crate::{Handler, Request, Response, Result, Route, StatusCode, path};

/// OAuth-enabled auth state
#[derive(Clone)]
pub struct AuthState {
    pub oauth_service: Arc<OAuthService>,
}

impl Default for AuthState {
    fn default() -> Self {
        Self {
            oauth_service: Arc::new(OAuthService::from_env()),
        }
    }
}

/// Handler for POST /register
struct RegisterHandler;

#[async_trait]
impl Handler for RegisterHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let req: RegisterRequest = request.json()?;

        match AuthViewSet::register(req) {
            RegisterResult::Success(response) => Response::created().with_json(&response),
            RegisterResult::EmailExists => Response::new(StatusCode::CONFLICT)
                .with_json(&json!({"error": "Email already exists"})),
            RegisterResult::InvalidInput(msg) => {
                Response::bad_request().with_json(&json!({"error": msg}))
            }
        }
    }
}

/// Handler for POST /login
struct LoginHandler;

#[async_trait]
impl Handler for LoginHandler {
    async fn handle(&self, request: Request) -> Result<Response> {
        let req: LoginRequest = request.json()?;

        match AuthViewSet::login(req, None) {
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

        // TODO: Verify state token against stored value for CSRF protection

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
        path("/register/", RegisterHandler),
        path("/login/", LoginHandler),
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

    #[tokio::test]
    async fn test_register_endpoint() {
        let handler = RegisterHandler;

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
        let handler = RegisterHandler;

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
        let handler = LoginHandler;

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
