//! HaiLanGo Frontend - WASM Entry Point

use hailango_frontend::AppRouter;
use leptos::*;

fn main() {
    // Initialize panic hook
    console_error_panic_hook::set_once();

    // Configure logging
    _ = console_log::init_with_level(log::Level::Debug);

    log::info!("HaiLanGo frontend starting...");

    // Mount the application
    mount_to_body(|| {
        view! {
            <AppRouter />
        }
    });

    log::info!("HaiLanGo frontend mounted");
}
