//! Page Layouts

use leptos::*;

/// Main application layout with navigation
#[component]
pub fn MainLayout(children: Children) -> impl IntoView {
    view! {
        <div class="app-container">
            <NavBar />
            <main class="main-content">
                {children()}
            </main>
            <Footer />
        </div>
    }
}

/// Navigation bar
#[component]
pub fn NavBar() -> impl IntoView {
    view! {
        <nav class="navbar">
            <div class="navbar-brand">
                <a href="/" class="logo">"HaiLanGo"</a>
            </div>
            <div class="navbar-menu">
                <a href="/books" class="nav-link">"Books"</a>
                <a href="/learn" class="nav-link">"Learn"</a>
                <a href="/review" class="nav-link">"Review"</a>
                <a href="/profile" class="nav-link">"Profile"</a>
            </div>
        </nav>
    }
}

/// Footer
#[component]
pub fn Footer() -> impl IntoView {
    view! {
        <footer class="footer">
            <p>"HaiLanGo - AI Language Learning"</p>
            <p class="copyright">"2026"</p>
        </footer>
    }
}

/// Auth layout for login/register
#[component]
pub fn AuthLayout(children: Children) -> impl IntoView {
    view! {
        <div class="auth-container">
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
