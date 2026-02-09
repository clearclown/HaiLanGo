//! Reusable UI Components

use leptos::*;

/// Button component with variants
#[component]
pub fn Button(
    #[prop(into)] text: MaybeSignal<String>,
    #[prop(default = "primary")] variant: &'static str,
    #[prop(optional)] on_click: Option<Callback<()>>,
    #[prop(into, default = false.into())] disabled: MaybeSignal<bool>,
    #[prop(default = "")] class: &'static str,
) -> impl IntoView {
    let btn_class = match variant {
        "primary" => format!("btn btn-primary {}", class),
        "secondary" => format!("btn btn-secondary {}", class),
        "danger" => format!("btn btn-danger {}", class),
        "success" => format!("btn btn-success {}", class),
        "ghost" => format!("btn btn-ghost {}", class),
        "icon" => format!("btn btn-icon {}", class),
        _ => format!("btn {}", class),
    };

    view! {
        <button
            class=btn_class
            disabled=move || disabled.get()
            on:click=move |_| {
                if let Some(handler) = on_click {
                    handler.call(());
                }
            }
        >
            {move || text.get()}
        </button>
    }
}

/// Input field component
#[component]
pub fn Input(
    #[prop(into)] label: String,
    #[prop(into)] name: String,
    #[prop(default = "text")] input_type: &'static str,
    #[prop(into, optional)] placeholder: Option<String>,
    #[prop(optional)] value: Option<RwSignal<String>>,
    #[prop(default = false)] required: bool,
) -> impl IntoView {
    let input_value = value.unwrap_or_else(|| create_rw_signal(String::new()));
    let placeholder_str = placeholder.unwrap_or_default();

    view! {
        <div class="form-group">
            <label class="form-label" for=name.clone()>{label}</label>
            <input
                type=input_type
                id=name.clone()
                name=name
                class="form-input"
                placeholder=placeholder_str
                required=required
                prop:value=move || input_value.get()
                on:input=move |ev| {
                    input_value.set(event_target_value(&ev));
                }
            />
        </div>
    }
}

/// Card component
#[component]
pub fn Card(
    #[prop(into, optional)] title: Option<String>,
    #[prop(default = "")] class: &'static str,
    children: Children,
) -> impl IntoView {
    let card_class = format!("card {}", class);

    view! {
        <div class=card_class>
            {title.map(|t| view! { <div class="card-header"><h3 class="card-title">{t}</h3></div> })}
            <div class="card-body">
                {children()}
            </div>
        </div>
    }
}

/// Loading spinner
#[component]
pub fn Spinner() -> impl IntoView {
    view! {
        <div class="spinner-container">
            <div class="spinner"></div>
        </div>
    }
}

/// Alert component
#[component]
pub fn Alert(
    #[prop(into)] message: String,
    #[prop(default = "info")] variant: &'static str,
) -> impl IntoView {
    let class = format!("alert alert-{}", variant);

    view! {
        <div class=class role="alert">
            {message}
        </div>
    }
}

/// Stat card for dashboard
#[component]
pub fn StatCard(
    #[prop(into)] label: String,
    #[prop(into)] value: MaybeSignal<String>,
    #[prop(into, optional)] icon: Option<String>,
    #[prop(default = "primary")] color: &'static str,
) -> impl IntoView {
    let color_class = format!("stat-card stat-card-{}", color);

    view! {
        <div class=color_class>
            {icon.map(|i| view! { <span class="stat-icon">{i}</span> })}
            <div class="stat-value">{move || value.get()}</div>
            <div class="stat-label">{label}</div>
        </div>
    }
}

/// Progress bar
#[component]
pub fn ProgressBar(
    #[prop(into)] value: MaybeSignal<f32>,
    #[prop(default = "primary")] color: &'static str,
    #[prop(into, optional)] label: Option<MaybeSignal<String>>,
) -> impl IntoView {
    let bar_class = format!("progress-fill progress-{}", color);

    view! {
        <div class="progress-bar">
            <div
                class=bar_class
                style=move || format!("width: {}%", (value.get() * 100.0).min(100.0).max(0.0))
            />
            {label.map(|l| view! {
                <span class="progress-label">{move || l.get()}</span>
            })}
        </div>
    }
}

/// Range slider input
#[component]
pub fn RangeInput(
    #[prop(into)] label: String,
    value: RwSignal<f32>,
    #[prop(default = 0.0)] min: f32,
    #[prop(default = 1.0)] max: f32,
    #[prop(default = 0.1)] step: f32,
    #[prop(optional)] format_value: Option<fn(f32) -> String>,
) -> impl IntoView {
    let display_fn = format_value.unwrap_or(|v| format!("{:.1}", v));

    view! {
        <div class="range-group">
            <div class="range-header">
                <label class="form-label">{label}</label>
                <span class="range-value">{move || display_fn(value.get())}</span>
            </div>
            <input
                type="range"
                class="range-input"
                min=min.to_string()
                max=max.to_string()
                step=step.to_string()
                prop:value=move || value.get().to_string()
                on:input=move |ev| {
                    if let Ok(v) = event_target_value(&ev).parse::<f32>() {
                        value.set(v);
                    }
                }
            />
        </div>
    }
}

/// Badge/tag component
#[component]
pub fn Badge(
    #[prop(into)] text: String,
    #[prop(default = "primary")] variant: &'static str,
) -> impl IntoView {
    let class = format!("badge badge-{}", variant);

    view! {
        <span class=class>{text}</span>
    }
}

/// Modal dialog
#[component]
pub fn Modal(
    #[prop(into)] title: String,
    #[prop(into)] open: MaybeSignal<bool>,
    on_close: Callback<()>,
    children: Children,
) -> impl IntoView {
    let content = children();

    view! {
        <div
            class="modal-overlay"
            style=move || if open.get() { "display: flex" } else { "display: none" }
            on:click=move |_| on_close.call(())
        >
            <div class="modal-content" on:click=|ev| ev.stop_propagation()>
                <div class="modal-header">
                    <h3>{title}</h3>
                    <button class="modal-close" on:click=move |_| on_close.call(())>
                        "\u{2715}"
                    </button>
                </div>
                <div class="modal-body">
                    {content}
                </div>
            </div>
        </div>
    }
}

/// OAuth provider button
#[component]
pub fn OAuthButton(
    #[prop(into)] provider: String,
    #[prop(into)] label: String,
    on_click: Callback<String>,
) -> impl IntoView {
    let provider_class = format!("btn btn-oauth btn-oauth-{}", provider.to_lowercase());
    let provider_clone = provider.clone();

    view! {
        <button
            class=provider_class
            on:click=move |_| on_click.call(provider_clone.clone())
        >
            <span class="oauth-icon">{provider_icon(&provider)}</span>
            <span>{label}</span>
        </button>
    }
}

fn provider_icon(provider: &str) -> &'static str {
    match provider.to_lowercase().as_str() {
        "google" => "\u{1F310}",
        "github" => "\u{1F4BB}",
        _ => "\u{1F511}",
    }
}

/// Divider with optional text
#[component]
pub fn Divider(#[prop(into, optional)] text: Option<String>) -> impl IntoView {
    view! {
        <div class="divider">
            <hr />
            {text.map(|t| view! { <span class="divider-text">{t}</span> })}
            {Some(()).map(|_| view! { <hr /> }).filter(|_| true)}
        </div>
    }
}

/// Empty state placeholder
#[component]
pub fn EmptyState(
    #[prop(into)] message: String,
    #[prop(into, optional)] icon: Option<String>,
    #[prop(optional)] action: Option<Children>,
) -> impl IntoView {
    view! {
        <div class="empty-state">
            {icon.map(|i| view! { <div class="empty-icon">{i}</div> })}
            <p class="empty-message">{message}</p>
            {action.map(|a| view! { <div class="empty-action">{a()}</div> })}
        </div>
    }
}
