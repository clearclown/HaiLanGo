//! Application Routes

use leptos::*;
use leptos_router::*;

use crate::pages::*;

/// Application router with all page routes
#[component]
pub fn AppRouter() -> impl IntoView {
    view! {
        <Router>
            <Routes>
                <Route path="/" view=DashboardPage />
                <Route path="/login" view=LoginPage />
                <Route path="/register" view=RegisterPage />
                <Route path="/callback/:provider" view=OAuthCallbackPage />
                <Route path="/books" view=BooksPage />
                <Route path="/learn/:id" view=LearningSessionPage />
                <Route path="/teacher" view=TeacherModePage />
                <Route path="/review" view=ReviewPage />
            </Routes>
        </Router>
    }
}
