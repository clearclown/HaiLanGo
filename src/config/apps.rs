//! HaiLanGo application configuration
//!
//! Declares the installed apps and their metadata for the Reinhardt framework.
//! Each app corresponds to a module under `src/apps/`.

/// Metadata for an installed app module.
pub struct AppConfig {
    pub name: &'static str,
    pub label: &'static str,
}

/// All apps installed in HaiLanGo.
///
/// This list drives app discovery, migration ordering, and admin registration.
pub const INSTALLED_APPS: &[AppConfig] = &[
    AppConfig {
        name: "hailango.auth",
        label: "auth",
    },
    AppConfig {
        name: "hailango.books",
        label: "books",
    },
    AppConfig {
        name: "hailango.learning",
        label: "learning",
    },
    AppConfig {
        name: "hailango.review",
        label: "review",
    },
    AppConfig {
        name: "hailango.tts",
        label: "tts",
    },
    AppConfig {
        name: "hailango.stt",
        label: "stt",
    },
    AppConfig {
        name: "hailango.teacher_mode",
        label: "teacher_mode",
    },
];
