//! HaiLanGo Frontend - WASM Entry Point
//!
//! This binary is compiled to WASM and runs in the browser.
//! When compiled for native targets, it prints a message directing to WASM build.

#[cfg(target_arch = "wasm32")]
fn main() {
    use hailango::pages::routes::AppRouter;
    use reinhardt::pages::prelude::*;

    // Initialize panic hook for better error messages in browser console
    console_error_panic_hook::set_once();

    // Configure logging for browser console
    _ = console_log::init_with_level(log::Level::Debug);

    log::info!("HaiLanGo frontend starting...");

    // Mount the application
    mount_to_body(|| {
        view! {
            <AppRouter />
        }
    });

    log::info!("HaiLanGo frontend mounted successfully");
}

#[cfg(not(target_arch = "wasm32"))]
fn main() {
    eprintln!("This binary is intended to be compiled for WASM.");
    eprintln!("To build the frontend, run:");
    eprintln!("  trunk build --release");
    eprintln!("Or for development:");
    eprintln!("  trunk serve");
    std::process::exit(1);
}
