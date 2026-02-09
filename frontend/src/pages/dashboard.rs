//! Dashboard / Home page

use leptos::*;

use crate::api::{ApiClient, BookItem, ReviewStats};
use crate::components::{Button, Card, EmptyState, Spinner, StatCard};
use crate::layouts::MainLayout;

/// Dashboard with learning stats and recent books
#[component]
pub fn DashboardPage() -> impl IntoView {
    let stats = create_rw_signal::<Option<ReviewStats>>(None);
    let recent_books = create_rw_signal::<Vec<BookItem>>(vec![]);
    let loading = create_rw_signal(true);

    create_effect(move |_| {
        spawn_local(async move {
            let (stats_result, books_result) = (
                ApiClient::get_review_stats().await,
                ApiClient::get_books().await,
            );

            if let Ok(s) = stats_result {
                stats.set(Some(s));
            }
            if let Ok(b) = books_result {
                recent_books.set(b.into_iter().take(4).collect());
            }
            loading.set(false);
        });
    });

    view! {
        <MainLayout>
            <div class="dashboard">
                <div class="dashboard-header">
                    <h1>"Dashboard"</h1>
                    <p class="dashboard-subtitle">"Welcome back! Here's your learning progress."</p>
                </div>

                // Stats grid
                <div class="stats-grid">
                    <StatCard
                        label="Words Learned"
                        value=Signal::derive(move || {
                            stats.get().map(|s| s.total_cards.to_string()).unwrap_or("0".to_string())
                        })
                        icon="\u{1F4DA}".to_string()
                        color="primary"
                    />
                    <StatCard
                        label="Due Today"
                        value=Signal::derive(move || {
                            stats.get().map(|s| s.due_today.to_string()).unwrap_or("0".to_string())
                        })
                        icon="\u{1F4CB}".to_string()
                        color="warning"
                    />
                    <StatCard
                        label="Study Streak"
                        value=Signal::derive(move || {
                            stats.get().map(|s| format!("{} days", s.streak)).unwrap_or("0 days".to_string())
                        })
                        icon="\u{1F525}".to_string()
                        color="success"
                    />
                    <StatCard
                        label="Accuracy"
                        value=Signal::derive(move || {
                            stats.get().map(|s| format!("{:.0}%", s.accuracy * 100.0)).unwrap_or("--".to_string())
                        })
                        icon="\u{1F3AF}".to_string()
                        color="info"
                    />
                </div>

                // Quick actions
                <div class="quick-actions">
                    <a href="/review" class="action-link">
                        <Button text="Review Now".to_string() variant="primary" />
                    </a>
                    <a href="/books" class="action-link">
                        <Button text="My Books".to_string() variant="secondary" />
                    </a>
                    <a href="/teacher" class="action-link">
                        <Button text="Teacher Mode".to_string() variant="secondary" />
                    </a>
                </div>

                // Recent books
                <section class="dashboard-section">
                    <div class="section-header">
                        <h2>"Recent Books"</h2>
                        <a href="/books" class="section-link">"View all"</a>
                    </div>

                    {move || if loading.get() {
                        view! { <Spinner /> }.into_view()
                    } else if recent_books.get().is_empty() {
                        view! {
                            <EmptyState
                                message="No books yet. Upload your first book to start learning!"
                                icon="\u{1F4D6}".to_string()
                            />
                        }.into_view()
                    } else {
                        view! {
                            <div class="books-grid">
                                <For
                                    each=move || recent_books.get()
                                    key=|book| book.id.clone()
                                    children=move |book| {
                                        let book_id = book.id.clone();
                                        let href = format!("/learn/{}", book_id);
                                        view! {
                                            <a href=href class="book-card-link">
                                                <Card>
                                                    <h4 class="book-title">{book.title.clone()}</h4>
                                                    <p class="book-author">{book.author.clone()}</p>
                                                    <div class="book-meta">
                                                        <span>{format!("{} pages", book.total_pages)}</span>
                                                        <span class="book-progress">
                                                            {format!("{:.0}%", book.progress * 100.0)}
                                                        </span>
                                                    </div>
                                                </Card>
                                            </a>
                                        }
                                    }
                                />
                            </div>
                        }.into_view()
                    }}
                </section>
            </div>
        </MainLayout>
    }
}
