//! Login, Register, and OAuth callback pages

use leptos::*;
use leptos_router::*;

use crate::api::ApiClient;
use crate::components::{Alert, Button, Divider, Input, OAuthButton, Spinner};
use crate::layouts::AuthLayout;

/// Login page with email/password and OAuth
#[component]
pub fn LoginPage() -> impl IntoView {
    let email = create_rw_signal(String::new());
    let password = create_rw_signal(String::new());
    let error = create_rw_signal::<Option<String>>(None);
    let loading = create_rw_signal(false);
    let oauth_loading = create_rw_signal::<Option<String>>(None);

    let on_submit = move |ev: web_sys::SubmitEvent| {
        ev.prevent_default();
        loading.set(true);
        error.set(None);

        let email_val = email.get();
        let password_val = password.get();

        spawn_local(async move {
            match ApiClient::login(&email_val, &password_val).await {
                Ok(_) => {
                    if let Some(window) = web_sys::window() {
                        let _ = window.location().set_href("/");
                    }
                }
                Err(e) => error.set(Some(e)),
            }
            loading.set(false);
        });
    };

    let on_oauth = Callback::new(move |provider: String| {
        oauth_loading.set(Some(provider.clone()));
        error.set(None);

        spawn_local(async move {
            match ApiClient::get_oauth_url(&provider).await {
                Ok(redirect) => {
                    if let Some(window) = web_sys::window() {
                        let _ = window.location().set_href(&redirect.auth_url);
                    }
                }
                Err(e) => {
                    error.set(Some(e));
                    oauth_loading.set(None);
                }
            }
        });
    });

    view! {
        <AuthLayout>
            <form class="auth-form" on:submit=on_submit>
                <h2>"Sign In"</h2>

                {move || error.get().map(|msg| view! {
                    <Alert message=msg variant="danger" />
                })}

                <Input
                    label="Email"
                    name="email"
                    input_type="email"
                    placeholder="your@email.com".to_string()
                    value=email
                    required=true
                />

                <Input
                    label="Password"
                    name="password"
                    input_type="password"
                    placeholder="Enter password".to_string()
                    value=password
                    required=true
                />

                <Button
                    text=Signal::derive(move || {
                        if loading.get() { "Signing in...".to_string() } else { "Sign In".to_string() }
                    })
                    variant="primary"
                    disabled=loading
                    class="btn-full-width"
                />

                <Divider text="or continue with".to_string() />

                <div class="oauth-buttons">
                    <OAuthButton
                        provider="Google"
                        label="Google"
                        on_click=on_oauth
                    />
                    <OAuthButton
                        provider="GitHub"
                        label="GitHub"
                        on_click=on_oauth
                    />
                </div>

                {move || oauth_loading.get().map(|p| view! {
                    <div class="oauth-status">
                        <Spinner />
                        <p>{format!("Redirecting to {}...", p)}</p>
                    </div>
                })}

                <p class="auth-link">
                    "Don't have an account? "
                    <A href="/register">"Sign up"</A>
                </p>
            </form>
        </AuthLayout>
    }
}

/// Register page with email/password and OAuth
#[component]
pub fn RegisterPage() -> impl IntoView {
    let email = create_rw_signal(String::new());
    let password = create_rw_signal(String::new());
    let display_name = create_rw_signal(String::new());
    let error = create_rw_signal::<Option<String>>(None);
    let success = create_rw_signal(false);
    let loading = create_rw_signal(false);

    let on_submit = move |ev: web_sys::SubmitEvent| {
        ev.prevent_default();
        loading.set(true);
        error.set(None);

        let email_val = email.get();
        let password_val = password.get();
        let name_val = display_name.get();

        spawn_local(async move {
            match ApiClient::register(&email_val, &password_val, &name_val).await {
                Ok(_) => {
                    success.set(true);
                    if let Some(window) = web_sys::window() {
                        let _ = window.location().set_href("/");
                    }
                }
                Err(e) => error.set(Some(e)),
            }
            loading.set(false);
        });
    };

    let on_oauth = Callback::new(move |provider: String| {
        error.set(None);
        spawn_local(async move {
            match ApiClient::get_oauth_url(&provider).await {
                Ok(redirect) => {
                    if let Some(window) = web_sys::window() {
                        let _ = window.location().set_href(&redirect.auth_url);
                    }
                }
                Err(e) => error.set(Some(e)),
            }
        });
    });

    view! {
        <AuthLayout>
            <form class="auth-form" on:submit=on_submit>
                <h2>"Create Account"</h2>

                {move || error.get().map(|msg| view! {
                    <Alert message=msg variant="danger" />
                })}
                {move || if success.get() {
                    view! { <Alert message="Account created! Redirecting...".to_string() variant="success" /> }.into_view()
                } else {
                    view! {}.into_view()
                }}

                <Input
                    label="Display Name"
                    name="display_name"
                    placeholder="Your name".to_string()
                    value=display_name
                    required=true
                />

                <Input
                    label="Email"
                    name="email"
                    input_type="email"
                    placeholder="your@email.com".to_string()
                    value=email
                    required=true
                />

                <Input
                    label="Password"
                    name="password"
                    input_type="password"
                    placeholder="Min 8 characters".to_string()
                    value=password
                    required=true
                />

                <Button
                    text=Signal::derive(move || {
                        if loading.get() { "Creating...".to_string() } else { "Create Account".to_string() }
                    })
                    variant="primary"
                    disabled=loading
                    class="btn-full-width"
                />

                <Divider text="or continue with".to_string() />

                <div class="oauth-buttons">
                    <OAuthButton
                        provider="Google"
                        label="Google"
                        on_click=on_oauth
                    />
                    <OAuthButton
                        provider="GitHub"
                        label="GitHub"
                        on_click=on_oauth
                    />
                </div>

                <p class="auth-link">
                    "Already have an account? "
                    <A href="/login">"Sign in"</A>
                </p>
            </form>
        </AuthLayout>
    }
}

/// OAuth callback handler page
#[component]
pub fn OAuthCallbackPage() -> impl IntoView {
    let params = use_params_map();
    let query = use_query_map();
    let error = create_rw_signal::<Option<String>>(None);

    create_effect(move |_| {
        let provider = params.with(|p| p.get("provider").cloned().unwrap_or_default());
        let code = query.with(|q| q.get("code").cloned().unwrap_or_default());
        let state = query.with(|q| q.get("state").cloned().unwrap_or_default());

        if code.is_empty() {
            error.set(Some("No authorization code received".to_string()));
            return;
        }

        spawn_local(async move {
            match ApiClient::oauth_callback(&provider, &code, &state).await {
                Ok(_) => {
                    if let Some(window) = web_sys::window() {
                        let _ = window.location().set_href("/");
                    }
                }
                Err(e) => error.set(Some(e)),
            }
        });
    });

    view! {
        <AuthLayout>
            <div class="auth-form" style="text-align: center">
                {move || if let Some(err) = error.get() {
                    view! {
                        <Alert message=err variant="danger" />
                        <A href="/login">"Back to login"</A>
                    }.into_view()
                } else {
                    view! {
                        <Spinner />
                        <p>"Completing sign in..."</p>
                    }.into_view()
                }}
            </div>
        </AuthLayout>
    }
}
