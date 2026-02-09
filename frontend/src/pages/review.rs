//! SRS Review page with flashcards

use leptos::*;

use crate::api::{ApiClient, ReviewCard, ReviewStats};
use crate::components::{Button, Card, EmptyState, Spinner, StatCard};
use crate::layouts::MainLayout;

/// SRS Review session with flashcards and SM-2 ratings
#[component]
pub fn ReviewPage() -> impl IntoView {
    let stats = create_rw_signal::<Option<ReviewStats>>(None);
    let queue = create_rw_signal::<Vec<ReviewCard>>(vec![]);
    let current_index = create_rw_signal(0usize);
    let show_answer = create_rw_signal(false);
    let loading = create_rw_signal(true);
    let submitting = create_rw_signal(false);
    let session_complete = create_rw_signal(false);
    let cards_reviewed = create_rw_signal(0u32);

    // Load review queue and stats
    create_effect(move |_| {
        spawn_local(async move {
            let (stats_result, queue_result) = (
                ApiClient::get_review_stats().await,
                ApiClient::get_review_queue().await,
            );

            if let Ok(s) = stats_result {
                stats.set(Some(s));
            }
            if let Ok(q) = queue_result {
                if q.is_empty() {
                    session_complete.set(true);
                }
                queue.set(q);
            }
            loading.set(false);
        });
    });

    let current_card = Signal::derive(move || {
        let q = queue.get();
        let idx = current_index.get();
        q.get(idx).cloned()
    });

    let remaining = Signal::derive(move || {
        let total = queue.get().len();
        let done = current_index.get();
        if total > done { total - done } else { 0 }
    });

    let submit_rating = move |rating: u8| {
        submitting.set(true);
        show_answer.set(false);

        if let Some(card) = current_card.get() {
            let card_id = card.id.clone();
            spawn_local(async move {
                let _ = ApiClient::submit_review(&card_id, rating).await;
                cards_reviewed.update(|c| *c += 1);

                current_index.update(|i| *i += 1);
                if current_index.get() >= queue.get().len() {
                    session_complete.set(true);
                }
                submitting.set(false);
            });
        }
    };

    view! {
        <MainLayout>
            <div class="review-page">
                <div class="page-header">
                    <h1>"Review"</h1>
                    {move || if !session_complete.get() && !loading.get() {
                        view! {
                            <span class="remaining-badge">
                                {move || format!("{} cards remaining", remaining.get())}
                            </span>
                        }.into_view()
                    } else {
                        view! {}.into_view()
                    }}
                </div>

                // Stats row
                {move || stats.get().map(|s| view! {
                    <div class="review-stats-row">
                        <StatCard
                            label="Total Cards"
                            value=s.total_cards.to_string()
                            color="primary"
                        />
                        <StatCard
                            label="Due Today"
                            value=s.due_today.to_string()
                            color="warning"
                        />
                        <StatCard
                            label="Streak"
                            value=format!("{} days", s.streak)
                            color="success"
                        />
                        <StatCard
                            label="Accuracy"
                            value=format!("{:.0}%", s.accuracy * 100.0)
                            color="info"
                        />
                    </div>
                })}

                // Main review area
                {move || if loading.get() {
                    view! { <Spinner /> }.into_view()
                } else if session_complete.get() {
                    view! {
                        <Card class="completion-card">
                            <div class="review-complete">
                                <div class="complete-icon">"\u{1F389}"</div>
                                <h2>"Review Complete!"</h2>
                                <p>{move || format!("You reviewed {} cards today.", cards_reviewed.get())}</p>
                                <div class="complete-actions">
                                    <a href="/">
                                        <Button text="Back to Dashboard".to_string() variant="primary" />
                                    </a>
                                </div>
                            </div>
                        </Card>
                    }.into_view()
                } else if let Some(card) = current_card.get() {
                    view! {
                        <Card class="flashcard-card">
                            <div class="flashcard">
                                // Front: word
                                <div class="flashcard-front">
                                    <h2 class="flashcard-word">{card.word.clone()}</h2>
                                    {card.sentence.clone().map(|s| view! {
                                        <p class="flashcard-sentence">{s}</p>
                                    })}
                                </div>

                                // Back: reading + meaning (shown on reveal)
                                {move || if show_answer.get() {
                                    view! {
                                        <div class="flashcard-back">
                                            <div class="flashcard-divider" />
                                            <p class="flashcard-reading">{card.reading.clone()}</p>
                                            <p class="flashcard-meaning">{card.meaning.clone()}</p>
                                        </div>

                                        // Rating buttons (SM-2 inspired)
                                        <div class="rating-buttons">
                                            <button
                                                class="rating-btn rating-again"
                                                disabled=move || submitting.get()
                                                on:click=move |_| submit_rating(1)
                                            >
                                                <span class="rating-label">"Again"</span>
                                                <span class="rating-hint">"< 1 min"</span>
                                            </button>
                                            <button
                                                class="rating-btn rating-hard"
                                                disabled=move || submitting.get()
                                                on:click=move |_| submit_rating(2)
                                            >
                                                <span class="rating-label">"Hard"</span>
                                                <span class="rating-hint">"< 10 min"</span>
                                            </button>
                                            <button
                                                class="rating-btn rating-good"
                                                disabled=move || submitting.get()
                                                on:click=move |_| submit_rating(3)
                                            >
                                                <span class="rating-label">"Good"</span>
                                                <span class="rating-hint">"1 day"</span>
                                            </button>
                                            <button
                                                class="rating-btn rating-easy"
                                                disabled=move || submitting.get()
                                                on:click=move |_| submit_rating(4)
                                            >
                                                <span class="rating-label">"Easy"</span>
                                                <span class="rating-hint">"4 days"</span>
                                            </button>
                                        </div>
                                    }.into_view()
                                } else {
                                    view! {
                                        <div class="show-answer-section">
                                            <Button
                                                text="Show Answer".to_string()
                                                variant="primary"
                                                on_click=Callback::new(move |_| show_answer.set(true))
                                                class="btn-full-width btn-show-answer"
                                            />
                                        </div>
                                    }.into_view()
                                }}
                            </div>
                        </Card>
                    }.into_view()
                } else {
                    view! {
                        <EmptyState
                            message="No cards to review right now. Great job!"
                            icon="\u{2705}".to_string()
                        />
                    }.into_view()
                }}
            </div>
        </MainLayout>
    }
}
