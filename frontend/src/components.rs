//! Reusable UI Components

use leptos::*;

/// Button component with variants
#[component]
pub fn Button(
    #[prop(into)] text: MaybeSignal<String>,
    #[prop(default = "primary")] variant: &'static str,
    #[prop(optional)] on_click: Option<Callback<()>>,
    #[prop(into, default = false.into())] disabled: MaybeSignal<bool>,
) -> impl IntoView {
    let class = match variant {
        "primary" => "btn btn-primary",
        "secondary" => "btn btn-secondary",
        "danger" => "btn btn-danger",
        _ => "btn",
    };

    view! {
        <button
            class=class
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
            <label for=name.clone()>{label}</label>
            <input
                type=input_type
                id=name.clone()
                name=name
                class="form-control"
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
pub fn Card(#[prop(into, optional)] title: Option<String>, children: Children) -> impl IntoView {
    view! {
        <div class="card">
            {title.map(|t| view! { <div class="card-header"><h3>{t}</h3></div> })}
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
        <div class="spinner">
            <div class="spinner-border" role="status">
                <span class="sr-only">"Loading..."</span>
            </div>
        </div>
    }
}

/// Alert component
#[component]
pub fn Alert(#[prop(into)] message: String, #[prop(default = "info")] variant: &'static str) -> impl IntoView {
    let class = format!("alert alert-{}", variant);

    view! {
        <div class=class role="alert">
            {message}
        </div>
    }
}
