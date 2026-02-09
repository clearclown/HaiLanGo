//! Teacher Mode control page

use leptos::*;

use crate::api::{ApiClient, StartLessonRequest, UpdateTeacherConfig};
use crate::components::{Alert, Badge, Button, Card, Input, ProgressBar, RangeInput};
use crate::layouts::MainLayout;

/// Teacher Mode with playback controls and settings
#[component]
pub fn TeacherModePage() -> impl IntoView {
    // Session state
    let session_id = create_rw_signal::<Option<String>>(None);
    let status = create_rw_signal("idle".to_string());
    let current_page = create_rw_signal(0u32);
    let total_pages_val = create_rw_signal(0u32);
    let pages_completed = create_rw_signal(0u32);

    // Settings
    let speed = create_rw_signal(1.0f32);
    let page_interval = create_rw_signal(5.0f32);
    let repeat_count = create_rw_signal(1.0f32);
    let auto_advance = create_rw_signal(true);

    // Start lesson form
    let book_id = create_rw_signal(String::new());
    let start_page = create_rw_signal("1".to_string());
    let end_page = create_rw_signal("10".to_string());

    let error = create_rw_signal::<Option<String>>(None);
    let loading = create_rw_signal(false);

    // Read book_id from URL query if present
    create_effect(move |_| {
        if let Some(window) = web_sys::window() {
            if let Ok(search) = window.location().search() {
                if let Some(bid) = search.strip_prefix("?book=") {
                    book_id.set(bid.to_string());
                }
            }
        }
    });

    let is_active = Signal::derive(move || {
        let s = status.get();
        s == "playing" || s == "paused"
    });

    let is_playing = Signal::derive(move || status.get() == "playing");

    // Start lesson
    let on_start = move |ev: web_sys::SubmitEvent| {
        ev.prevent_default();
        loading.set(true);
        error.set(None);

        let req = StartLessonRequest {
            book_id: book_id.get(),
            start_page: start_page.get().parse().unwrap_or(1),
            end_page: end_page.get().parse().unwrap_or(10),
            speed: Some(speed.get()),
            page_interval: Some(page_interval.get() as u32),
            repeat_count: Some(repeat_count.get() as u32),
        };

        spawn_local(async move {
            match ApiClient::start_lesson(req).await {
                Ok(resp) => {
                    session_id.set(Some(resp.session_id));
                    status.set(resp.status);
                    current_page.set(resp.current_page);
                    total_pages_val.set(resp.total_pages);
                    pages_completed.set(resp.pages_completed);
                }
                Err(e) => error.set(Some(e)),
            }
            loading.set(false);
        });
    };

    // Playback control action
    let do_action = move |action: &'static str| {
        loading.set(true);
        error.set(None);

        spawn_local(async move {
            match ApiClient::teacher_action(action).await {
                Ok(resp) => {
                    status.set(resp.status);
                    current_page.set(resp.current_page);
                    pages_completed.set(resp.pages_completed);
                }
                Err(e) => error.set(Some(e)),
            }
            loading.set(false);
        });
    };

    // Update config
    let save_config = move |_| {
        let config = UpdateTeacherConfig {
            speed: Some(speed.get()),
            page_interval: Some(page_interval.get() as u32),
            repeat_count: Some(repeat_count.get() as u32),
            auto_advance: Some(auto_advance.get()),
        };

        spawn_local(async move {
            if let Err(e) = ApiClient::update_teacher_config(config).await {
                error.set(Some(e));
            }
        });
    };

    let progress = Signal::derive(move || {
        if total_pages_val.get() == 0 { 0.0f32 }
        else { pages_completed.get() as f32 / total_pages_val.get() as f32 }
    });

    let status_color = Signal::derive(move || {
        match status.get().as_str() {
            "playing" => "success",
            "paused" => "warning",
            "completed" => "info",
            "stopped" => "danger",
            _ => "secondary",
        }
    });

    view! {
        <MainLayout>
            <div class="teacher-page">
                <div class="page-header">
                    <h1>"Teacher Mode"</h1>
                    <Badge
                        text=Signal::derive(move || status.get().to_uppercase()).get()
                        variant=status_color.get()
                    />
                </div>

                {move || error.get().map(|msg| view! {
                    <Alert message=msg variant="danger" />
                })}

                {move || if !is_active.get() && session_id.get().is_none() {
                    // Start lesson form
                    view! {
                        <Card title="Start New Lesson".to_string()>
                            <form class="teacher-form" on:submit=on_start>
                                <Input
                                    label="Book ID"
                                    name="teacher_book_id"
                                    placeholder="Enter book ID".to_string()
                                    value=book_id
                                    required=true
                                />
                                <div class="form-row">
                                    <Input
                                        label="Start Page"
                                        name="start_page"
                                        input_type="number"
                                        value=start_page
                                    />
                                    <Input
                                        label="End Page"
                                        name="end_page"
                                        input_type="number"
                                        value=end_page
                                    />
                                </div>

                                // Settings
                                <div class="teacher-settings">
                                    <h4>"Playback Settings"</h4>
                                    <RangeInput
                                        label="Speed"
                                        value=speed
                                        min=0.5
                                        max=2.0
                                        step=0.1
                                        format_value=|v| format!("{:.1}x", v)
                                    />
                                    <RangeInput
                                        label="Page Interval"
                                        value=page_interval
                                        min=0.0
                                        max=30.0
                                        step=1.0
                                        format_value=|v| format!("{:.0}s", v)
                                    />
                                    <RangeInput
                                        label="Repeat Count"
                                        value=repeat_count
                                        min=1.0
                                        max=3.0
                                        step=1.0
                                        format_value=|v| format!("{:.0}", v)
                                    />

                                    <div class="form-group">
                                        <label class="checkbox-label">
                                            <input
                                                type="checkbox"
                                                checked=move || auto_advance.get()
                                                on:change=move |ev| {
                                                    use wasm_bindgen::JsCast;
                                                    if let Some(input) = ev.target().and_then(|t| t.dyn_into::<web_sys::HtmlInputElement>().ok()) {
                                                        auto_advance.set(input.checked());
                                                    }
                                                }
                                            />
                                            " Auto-advance to next page"
                                        </label>
                                    </div>
                                </div>

                                <Button
                                    text=Signal::derive(move || {
                                        if loading.get() { "Starting...".to_string() }
                                        else { "Start Lesson".to_string() }
                                    })
                                    variant="primary"
                                    disabled=loading
                                    class="btn-full-width"
                                />
                            </form>
                        </Card>
                    }.into_view()
                } else {
                    // Active session controls
                    view! {
                        <Card class="playback-card">
                            // Progress
                            <div class="playback-progress">
                                <ProgressBar
                                    value=progress
                                    label=Signal::derive(move || {
                                        format!("Page {} / {} ({:.0}%)",
                                            current_page.get(), total_pages_val.get(),
                                            progress.get() * 100.0)
                                    })
                                />
                            </div>

                            // Transport controls
                            <div class="transport-controls">
                                <Button
                                    text="\u{23EE}".to_string()
                                    variant="ghost"
                                    on_click=Callback::new(move |_| do_action("next"))
                                />

                                {move || if is_playing.get() {
                                    view! {
                                        <Button
                                            text="\u{23F8} Pause".to_string()
                                            variant="secondary"
                                            on_click=Callback::new(move |_| do_action("pause"))
                                            class="btn-transport"
                                        />
                                    }.into_view()
                                } else {
                                    view! {
                                        <Button
                                            text="\u{25B6} Resume".to_string()
                                            variant="primary"
                                            on_click=Callback::new(move |_| do_action("resume"))
                                            class="btn-transport"
                                        />
                                    }.into_view()
                                }}

                                <Button
                                    text="\u{23F9} Stop".to_string()
                                    variant="danger"
                                    on_click=Callback::new(move |_| {
                                        do_action("stop");
                                        session_id.set(None);
                                    })
                                />

                                <Button
                                    text="\u{23ED}".to_string()
                                    variant="ghost"
                                    on_click=Callback::new(move |_| do_action("next"))
                                />
                            </div>

                            // Live settings
                            <div class="playback-settings">
                                <h4>"Adjust Settings"</h4>
                                <RangeInput
                                    label="Speed"
                                    value=speed
                                    min=0.5
                                    max=2.0
                                    step=0.1
                                    format_value=|v| format!("{:.1}x", v)
                                />
                                <Button
                                    text="Apply".to_string()
                                    variant="secondary"
                                    on_click=Callback::new(save_config)
                                />
                            </div>
                        </Card>
                    }.into_view()
                }}

                // Session info
                {move || if status.get() == "completed" {
                    view! {
                        <Card class="completion-card">
                            <div class="completion-content">
                                <h3>"\u{1F389} Lesson Complete!"</h3>
                                <p>{move || format!("Completed {} pages", pages_completed.get())}</p>
                                <Button
                                    text="Start New Lesson".to_string()
                                    variant="primary"
                                    on_click=Callback::new(move |_| {
                                        session_id.set(None);
                                        status.set("idle".to_string());
                                    })
                                />
                            </div>
                        </Card>
                    }.into_view()
                } else {
                    view! {}.into_view()
                }}
            </div>
        </MainLayout>
    }
}
