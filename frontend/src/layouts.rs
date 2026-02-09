//! Page Layouts

use leptos::*;
use leptos_router::A;

use crate::api::ApiClient;

/// Main application layout with navigation
#[component]
pub fn MainLayout(children: Children) -> impl IntoView {
    view! {
        <div class="main-layout">
            <NavBar />
            <main class="main-content">
                {children()}
            </main>
            <Footer />
        </div>
    }
}

/// Navigation bar with auth-aware links
#[component]
pub fn NavBar() -> impl IntoView {
    let is_auth = ApiClient::is_authenticated();

    view! {
        <nav class="nav-bar">
            <A href="/" class="nav-brand">"HaiLanGo"</A>
            <div class="nav-links">
                <A href="/books" class="nav-link">"Books"</A>
                <A href="/review" class="nav-link">"Review"</A>
                <A href="/teacher" class="nav-link">"Teacher"</A>
                {if is_auth {
                    view! {
                        <button
                            class="nav-link nav-logout"
                            on:click=move |_| ApiClient::logout()
                        >
                            "Logout"
                        </button>
                    }.into_view()
                } else {
                    view! {
                        <A href="/login" class="nav-link">"Login"</A>
                    }.into_view()
                }}
            </div>
        </nav>
    }
}

/// Footer
#[component]
pub fn Footer() -> impl IntoView {
    view! {
        <footer class="footer">
            <p>"HaiLanGo - AI-Powered Language Learning"</p>
        </footer>
    }
}

/// Auth layout for login/register
#[component]
pub fn AuthLayout(children: Children) -> impl IntoView {
    view! {
        <div class="auth-layout">
            <div class="auth-card">
                <div class="auth-logo">
                    <h1>"HaiLanGo"</h1>
                    <p>"AI-Powered Language Learning"</p>
                </div>
                {children()}
            </div>
        </div>
    }
}
