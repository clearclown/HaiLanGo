//! Reusable UI Components
//!
//! Common components used across the application.

use reinhardt::pages::prelude::*;

/// Button component with variants
#[component]
pub fn Button(
    #[prop] text: String,
    #[prop(default = "primary")] variant: &'static str,
    #[prop(optional)] on_click: Option<Box<dyn Fn()>>,
    #[prop(default = false)] disabled: bool,
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
            disabled=disabled
            on:click=move |_| {
                if let Some(ref handler) = on_click {
                    handler();
                }
            }
        >
            {text}
        </button>
    }
}

/// Input field component
#[component]
pub fn Input(
    #[prop] label: String,
    #[prop] name: String,
    #[prop(default = "text")] input_type: &'static str,
    #[prop(optional)] placeholder: Option<String>,
    #[prop(optional)] value: Option<RwSignal<String>>,
    #[prop(default = false)] required: bool,
) -> impl IntoView {
    let input_value = value.unwrap_or_else(|| create_rw_signal(String::new()));

    view! {
        <div class="form-group">
            <label for=name.clone()>{label}</label>
            <input
                type=input_type
                id=name.clone()
                name=name
                class="form-control"
                placeholder=placeholder.unwrap_or_default()
                required=required
                prop:value=move || input_value.get()
                on:input=move |ev| {
                    input_value.set(event_target_value(&ev));
                }
            />
        </div>
    }
}

/// Card component for content containers
#[component]
pub fn Card(#[prop(optional)] title: Option<String>, children: Children) -> impl IntoView {
    view! {
        <div class="card">
            {title.map(|t| view! { <div class="card-header"><h3>{t}</h3></div> })}
            <div class="card-body">
                {children()}
            </div>
        </div>
    }
}

/// Loading spinner component
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

/// Alert component for messages
#[component]
pub fn Alert(
    #[prop] message: String,
    #[prop(default = "info")] variant: &'static str,
) -> impl IntoView {
    let class = format!("alert alert-{}", variant);

    view! {
        <div class=class role="alert">
            {message}
        </div>
    }
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_component_module_exists() {
        // Components are tested via integration tests
        assert!(true);
    }
}
