//! Learning session page - page display, TTS, pronunciation practice

use leptos::*;
use leptos_router::*;

use crate::api::ApiClient;
use crate::components::{Alert, Button, Card, ProgressBar, Spinner};
use crate::layouts::MainLayout;

/// Learning session with page content, TTS playback, and pronunciation practice
#[component]
pub fn LearningSessionPage() -> impl IntoView {
    let params = use_params_map();

    let book_id = Signal::derive(move || {
        params.with(|p| p.get("id").cloned().unwrap_or_default())
    });
    let current_page = create_rw_signal(1u32);
    let total_pages = create_rw_signal(1u32);
    let page_text = create_rw_signal(String::new());
    let book_title = create_rw_signal(String::new());
    let loading = create_rw_signal(true);
    let error = create_rw_signal::<Option<String>>(None);
    let tts_playing = create_rw_signal(false);
    let tts_language = create_rw_signal("zh".to_string());
    let show_pronunciation = create_rw_signal(false);

    // Load page content
    let load_page = move || {
        loading.set(true);
        error.set(None);
        let bid = book_id.get();
        let page = current_page.get();

        spawn_local(async move {
            match ApiClient::get_page_content(&bid, page).await {
                Ok(content) => {
                    page_text.set(content.text);
                    book_title.set(content.book_title);
                    total_pages.set(content.total_pages);
                }
                Err(e) => error.set(Some(e)),
            }
            loading.set(false);
        });
    };

    create_effect(move |_| {
        let _ = book_id.get();
        load_page();
    });

    let on_prev = move |_| {
        if current_page.get() > 1 {
            current_page.update(|p| *p -= 1);
            load_page();
        }
    };

    let on_next = move |_| {
        if current_page.get() < total_pages.get() {
            current_page.update(|p| *p += 1);
            load_page();
        }
    };

    let on_tts = move |_| {
        tts_playing.set(true);
        let text = page_text.get();
        let lang = tts_language.get();

        spawn_local(async move {
            match ApiClient::synthesize_speech(&text, &lang).await {
                Ok(audio_url) => {
                    // Play audio using browser Audio API
                    if let Some(window) = web_sys::window() {
                        if let Some(document) = window.document() {
                            if let Ok(audio) = document.create_element("audio") {
                                let _ = audio.set_attribute("src", &audio_url);
                                let _ = audio.set_attribute("autoplay", "true");
                                // Audio element handles playback
                            }
                        }
                    }
                }
                Err(e) => log::error!("TTS failed: {}", e),
            }
            tts_playing.set(false);
        });
    };

    let progress = Signal::derive(move || {
        if total_pages.get() == 0 {
            0.0f32
        } else {
            current_page.get() as f32 / total_pages.get() as f32
        }
    });

    view! {
        <MainLayout>
            <div class="learning-page">
                // Header with book info
                <div class="learning-header">
                    <a href="/books" class="back-link">"\u{2190} Back to Books"</a>
                    <h1 class="learning-title">{move || book_title.get()}</h1>
                    <div class="learning-progress">
                        <span class="page-indicator">
                            {move || format!("Page {} / {}", current_page.get(), total_pages.get())}
                        </span>
                        <ProgressBar
                            value=progress
                            label=Signal::derive(move || format!("{:.0}%", progress.get() * 100.0))
                        />
                    </div>
                </div>

                {move || error.get().map(|msg| view! {
                    <Alert message=msg variant="danger" />
                })}

                {move || if loading.get() {
                    view! { <Spinner /> }.into_view()
                } else {
                    view! {
                        // Page content
                        <Card class="content-card">
                            <div class="page-content">
                                <pre class="page-text">{move || page_text.get()}</pre>
                            </div>
                        </Card>

                        // Controls bar
                        <div class="learning-controls">
                            // Navigation
                            <div class="nav-controls">
                                <Button
                                    text="\u{2190} Previous".to_string()
                                    variant="secondary"
                                    on_click=Callback::new(on_prev)
                                    disabled=Signal::derive(move || current_page.get() <= 1)
                                />
                                <Button
                                    text="Next \u{2192}".to_string()
                                    variant="secondary"
                                    on_click=Callback::new(on_next)
                                    disabled=Signal::derive(move || current_page.get() >= total_pages.get())
                                />
                            </div>

                            // TTS controls
                            <div class="tts-controls">
                                <select
                                    class="form-input tts-lang-select"
                                    on:change=move |ev| tts_language.set(event_target_value(&ev))
                                >
                                    <option value="zh" selected=true>"Chinese"</option>
                                    <option value="en">"English"</option>
                                    <option value="ja">"Japanese"</option>
                                    <option value="ko">"Korean"</option>
                                </select>
                                <Button
                                    text=Signal::derive(move || {
                                        if tts_playing.get() { "\u{23F9} Playing...".to_string() }
                                        else { "\u{25B6} Listen".to_string() }
                                    })
                                    variant="primary"
                                    on_click=Callback::new(on_tts)
                                    disabled=tts_playing
                                />
                            </div>
                        </div>

                        // Pronunciation practice section
                        <div class="pronunciation-section">
                            <Button
                                text=Signal::derive(move || {
                                    if show_pronunciation.get() { "Hide Practice".to_string() }
                                    else { "\u{1F3A4} Practice Pronunciation".to_string() }
                                })
                                variant="ghost"
                                on_click=Callback::new(move |_| show_pronunciation.update(|v| *v = !*v))
                            />

                            {move || if show_pronunciation.get() {
                                view! {
                                    <Card class="practice-card">
                                        <div class="practice-content">
                                            <p class="practice-instructions">
                                                "Read the text aloud, then compare with the TTS pronunciation."
                                            </p>
                                            <div class="practice-actions">
                                                <Button
                                                    text="\u{1F3A4} Record".to_string()
                                                    variant="danger"
                                                />
                                                <Button
                                                    text="\u{25B6} Listen Again".to_string()
                                                    variant="secondary"
                                                    on_click=Callback::new(on_tts)
                                                />
                                            </div>
                                        </div>
                                    </Card>
                                }.into_view()
                            } else {
                                view! {}.into_view()
                            }}
                        </div>
                    }.into_view()
                }}
            </div>
        </MainLayout>
    }
}
