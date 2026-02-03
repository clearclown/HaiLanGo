//! Page Routes
//!
//! Application page components and routing.

use reinhardt::pages::prelude::*;

use crate::pages::components::{Alert, Button, Card, Input, Spinner};
use crate::pages::layouts::{AuthLayout, MainLayout};

/// Home page - Landing/Dashboard
#[component]
pub fn HomePage() -> impl IntoView {
    view! {
        <MainLayout>
            <div class="dashboard">
                <h1>"Welcome to HaiLanGo"</h1>
                <div class="stats-grid">
                    <Card title="Words Learned".to_string()>
                        <p class="stat-number">"0"</p>
                    </Card>
                    <Card title="Books".to_string()>
                        <p class="stat-number">"0"</p>
                    </Card>
                    <Card title="Study Streak".to_string()>
                        <p class="stat-number">"0 days"</p>
                    </Card>
                </div>
                <div class="quick-actions">
                    <Button text="Start Learning".to_string() variant="primary" />
                    <Button text="Upload Book".to_string() variant="secondary" />
                </div>
            </div>
        </MainLayout>
    }
}

/// Login page
#[component]
pub fn LoginPage() -> impl IntoView {
    let email = create_rw_signal(String::new());
    let password = create_rw_signal(String::new());
    let error = create_rw_signal::<Option<String>>(None);
    let loading = create_rw_signal(false);

    let on_submit = move |ev: web_sys::SubmitEvent| {
        ev.prevent_default();
        loading.set(true);
        error.set(None);

        // TODO: Call login API
        spawn_local(async move {
            // Simulate API call
            loading.set(false);
        });
    };

    view! {
        <AuthLayout>
            <form class="auth-form" on:submit=on_submit>
                <h2>"Sign In"</h2>

                {move || error.get().map(|msg| view! {
                    <Alert message=msg variant="danger" />
                })}

                <Input
                    label="Email".to_string()
                    name="email".to_string()
                    input_type="email"
                    placeholder=Some("your@email.com".to_string())
                    value=Some(email)
                    required=true
                />

                <Input
                    label="Password".to_string()
                    name="password".to_string()
                    input_type="password"
                    placeholder=Some("Enter password".to_string())
                    value=Some(password)
                    required=true
                />

                <Button
                    text={if loading.get() { "Signing in...".to_string() } else { "Sign In".to_string() }}
                    variant="primary"
                    disabled=loading.get()
                />

                <p class="auth-link">
                    "Don't have an account? "
                    <a href="/register">"Sign up"</a>
                </p>
            </form>
        </AuthLayout>
    }
}

/// Register page
#[component]
pub fn RegisterPage() -> impl IntoView {
    let email = create_rw_signal(String::new());
    let password = create_rw_signal(String::new());
    let display_name = create_rw_signal(String::new());
    let error = create_rw_signal::<Option<String>>(None);
    let loading = create_rw_signal(false);

    let on_submit = move |ev: web_sys::SubmitEvent| {
        ev.prevent_default();
        loading.set(true);
        error.set(None);

        // TODO: Call register API
        spawn_local(async move {
            loading.set(false);
        });
    };

    view! {
        <AuthLayout>
            <form class="auth-form" on:submit=on_submit>
                <h2>"Create Account"</h2>

                {move || error.get().map(|msg| view! {
                    <Alert message=msg variant="danger" />
                })}

                <Input
                    label="Display Name".to_string()
                    name="display_name".to_string()
                    placeholder=Some("Your name".to_string())
                    value=Some(display_name)
                    required=true
                />

                <Input
                    label="Email".to_string()
                    name="email".to_string()
                    input_type="email"
                    placeholder=Some("your@email.com".to_string())
                    value=Some(email)
                    required=true
                />

                <Input
                    label="Password".to_string()
                    name="password".to_string()
                    input_type="password"
                    placeholder=Some("Min 8 characters".to_string())
                    value=Some(password)
                    required=true
                />

                <Button
                    text={if loading.get() { "Creating...".to_string() } else { "Create Account".to_string() }}
                    variant="primary"
                    disabled=loading.get()
                />

                <p class="auth-link">
                    "Already have an account? "
                    <a href="/login">"Sign in"</a>
                </p>
            </form>
        </AuthLayout>
    }
}

/// Books list page
#[component]
pub fn BooksPage() -> impl IntoView {
    let books = create_rw_signal::<Vec<BookItem>>(vec![]);
    let loading = create_rw_signal(true);

    // Fetch books on mount
    create_effect(move |_| {
        spawn_local(async move {
            // TODO: Fetch from API
            loading.set(false);
        });
    });

    view! {
        <MainLayout>
            <div class="books-page">
                <div class="page-header">
                    <h1>"My Books"</h1>
                    <Button text="Upload New".to_string() variant="primary" />
                </div>

                {move || if loading.get() {
                    view! { <Spinner /> }.into_view()
                } else if books.get().is_empty() {
                    view! {
                        <div class="empty-state">
                            <p>"No books yet. Upload your first book to start learning!"</p>
                        </div>
                    }.into_view()
                } else {
                    view! {
                        <div class="books-grid">
                            {books.get().into_iter().map(|book| view! {
                                <Card title=book.title>
                                    <p>{book.author}</p>
                                    <p>{format!("{} pages", book.total_pages)}</p>
                                </Card>
                            }).collect_view()}
                        </div>
                    }.into_view()
                }}
            </div>
        </MainLayout>
    }
}

/// Book item for display
#[derive(Clone)]
struct BookItem {
    title: String,
    author: String,
    total_pages: u32,
}

/// Review page
#[component]
pub fn ReviewPage() -> impl IntoView {
    let queue_count = create_rw_signal(0);
    let current_word = create_rw_signal::<Option<VocabItem>>(None);
    let show_answer = create_rw_signal(false);

    view! {
        <MainLayout>
            <div class="review-page">
                <h1>"Review Session"</h1>

                <div class="review-stats">
                    <p>{format!("Cards remaining: {}", queue_count.get())}</p>
                </div>

                {move || if let Some(word) = current_word.get() {
                    view! {
                        <Card>
                            <div class="flashcard">
                                <h2 class="word">{word.word}</h2>

                                {if show_answer.get() {
                                    view! {
                                        <div class="answer">
                                            <p class="reading">{word.reading}</p>
                                            <p class="meaning">{word.meaning}</p>
                                        </div>
                                        <div class="rating-buttons">
                                            <Button text="Again".to_string() variant="danger" />
                                            <Button text="Hard".to_string() variant="secondary" />
                                            <Button text="Good".to_string() variant="primary" />
                                            <Button text="Easy".to_string() variant="primary" />
                                        </div>
                                    }.into_view()
                                } else {
                                    view! {
                                        <Button
                                            text="Show Answer".to_string()
                                            variant="primary"
                                            on_click=Some(Box::new(move || show_answer.set(true)))
                                        />
                                    }.into_view()
                                }}
                            </div>
                        </Card>
                    }.into_view()
                } else {
                    view! {
                        <div class="empty-state">
                            <p>"No cards to review. Great job!"</p>
                        </div>
                    }.into_view()
                }}
            </div>
        </MainLayout>
    }
}

/// Vocabulary item for review
#[derive(Clone)]
struct VocabItem {
    word: String,
    reading: String,
    meaning: String,
}

/// Application router
#[component]
pub fn AppRouter() -> impl IntoView {
    view! {
        <Router>
            <Routes>
                <Route path="/" view=HomePage />
                <Route path="/login" view=LoginPage />
                <Route path="/register" view=RegisterPage />
                <Route path="/books" view=BooksPage />
                <Route path="/review" view=ReviewPage />
            </Routes>
        </Router>
    }
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_routes_module_exists() {
        assert!(true);
    }
}
