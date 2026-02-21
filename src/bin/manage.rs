//! Management command entry point
//!
//! Provides CLI interface for administrative tasks:
//!   manage check       — check environment variables
//!   manage db-status   — test database connectivity
//!   manage routes      — list all registered URL routes
//!   manage version     — print application version

use hailango::config::apps::INSTALLED_APPS;
use hailango::config::urls::configure_urls;

/// Print usage and exit.
fn usage() {
    eprintln!(
        "HaiLanGo Management CLI v{}\n\n\
         Usage: manage <command>\n\n\
         Commands:\n\
         \x20  check       Check required environment variables\n\
         \x20  db-status   Attempt to connect to the configured database\n\
         \x20  routes      List all registered API routes\n\
         \x20  apps        List installed apps\n\
         \x20  version     Print application version",
        env!("CARGO_PKG_VERSION")
    );
    std::process::exit(1);
}

/// `manage check` — verify required env vars are set.
fn cmd_check() {
    let required = ["DATABASE_URL", "REDIS_URL", "JWT_SECRET"];

    let mut ok = true;
    for var in &required {
        match std::env::var(var) {
            Ok(v) if !v.is_empty() => {
                println!("[OK]  {var}");
            }
            _ => {
                eprintln!("[MISSING] {var}");
                ok = false;
            }
        }
    }

    let optional = [
        ("GOOGLE_CLOUD_TTS_API_KEY", "TTS provider"),
        ("OPENAI_API_KEY", "STT / AI provider"),
        ("ANTHROPIC_API_KEY", "LLM teacher mode"),
        ("GOOGLE_OAUTH_CLIENT_ID", "Google OAuth"),
        ("GITHUB_OAUTH_CLIENT_ID", "GitHub OAuth"),
    ];

    for (var, desc) in &optional {
        if std::env::var(var).is_ok() {
            println!("[OK]  {var} ({desc})");
        } else {
            println!("[INFO] {var} not set ({desc} disabled)");
        }
    }

    if ok {
        println!("\nAll required variables are set.");
    } else {
        eprintln!("\nSome required variables are missing.");
        std::process::exit(1);
    }
}

/// `manage db-status` — test DB connectivity (async, tokio runtime).
fn cmd_db_status() {
    let url = match std::env::var("DATABASE_URL") {
        Ok(u) => u,
        Err(_) => {
            eprintln!("DATABASE_URL is not set.");
            std::process::exit(1);
        }
    };

    let rt = tokio::runtime::Runtime::new().unwrap();
    rt.block_on(async {
        use sqlx::postgres::PgPoolOptions;

        println!("Connecting to database...");
        match PgPoolOptions::new()
            .max_connections(1)
            .acquire_timeout(std::time::Duration::from_secs(5))
            .connect(&url)
            .await
        {
            Ok(pool) => match sqlx::query("SELECT 1").fetch_one(&pool).await {
                Ok(_) => println!("Database connection: OK"),
                Err(e) => {
                    eprintln!("Database query failed: {e}");
                    std::process::exit(1);
                }
            },
            Err(e) => {
                eprintln!("Database connection failed: {e}");
                std::process::exit(1);
            }
        }
    });
}

/// `manage routes` — print registered API routes.
fn cmd_routes() {
    let router = configure_urls();
    let routes = router.get_routes();

    println!("Registered routes ({} total):", routes.len());
    for route in routes {
        println!("  {}", route.path);
    }
}

/// `manage apps` — list installed apps.
fn cmd_apps() {
    println!("Installed apps ({} total):", INSTALLED_APPS.len());
    for app in INSTALLED_APPS {
        println!("  {} (label: {})", app.name, app.label);
    }
}

fn main() {
    // Load .env if present
    let _ = dotenvy::dotenv();

    let args: Vec<String> = std::env::args().collect();
    let command = args.get(1).map(String::as_str);

    match command {
        Some("check") => cmd_check(),
        Some("db-status") => cmd_db_status(),
        Some("routes") => cmd_routes(),
        Some("apps") => cmd_apps(),
        Some("version") => println!("HaiLanGo v{}", env!("CARGO_PKG_VERSION")),
        _ => usage(),
    }
}
