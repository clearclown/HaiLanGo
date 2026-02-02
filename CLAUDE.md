# CLAUDE.md - HaiLanGo Project Guidelines

## Project Overview

HaiLanGo is an AI-powered language learning platform that transforms existing textbooks into interactive learning experiences. The project combines OCR, TTS, STT, and AI technologies to provide personalized language instruction.

**Repository**: https://github.com/clearclown/HaiLanGo

## Core Technology Stack

### Framework: Reinhardt (Rust Full-Stack)

[Reinhardt](https://github.com/kent8192/reinhardt-web) is a composable full-stack API framework for Rust, inspired by Django/FastAPI.

| Component | Crate | Description |
|-----------|-------|-------------|
| **ORM** | `reinhardt-db` | SeaQuery + sqlx integration, PostgreSQL/MySQL/SQLite |
| **Frontend** | `reinhardt-pages` | WASM + SSR reactive framework (Leptos/Solid.js style) |
| **REST API** | `reinhardt-rest` | ViewSets, Serializers, pagination, filtering |
| **GraphQL** | `reinhardt-graphql` | Schema generation, subscriptions |
| **WebSocket** | `reinhardt-websockets` | Real-time communication |
| **Auth** | `reinhardt-auth` | JWT, Token, Session, Basic authentication |
| **Admin** | `reinhardt-admin` | Django-style auto-generated admin panel |
| **i18n** | `reinhardt-i18n` | Internationalization support |

### Database
- **Primary**: PostgreSQL (via `reinhardt-db` with SeaQuery/sqlx)
- **Cache/Session**: Redis (separate dependency for caching, rate limiting)

### External APIs
- **OCR**: Google Vision API / Azure Computer Vision
- **TTS**: Google Cloud TTS / Azure Speech
- **STT**: OpenAI Whisper / Azure Speech
- **LLM**: Anthropic Claude API
- **Payments**: Stripe

## Project Structure

```
HaiLanGo/
├── src/
│   ├── config/
│   │   ├── settings/        # Environment-specific settings (TOML)
│   │   ├── urls.rs          # URL routing
│   │   └── apps.rs          # App configuration
│   ├── apps/
│   │   ├── auth/            # Authentication & users
│   │   ├── books/           # Book management & OCR
│   │   ├── learning/        # Learning sessions & pages
│   │   ├── tts/             # Text-to-Speech
│   │   ├── stt/             # Speech-to-Text & pronunciation
│   │   ├── review/          # SRS review system
│   │   └── teacher_mode/    # Automated lesson playback
│   └── pages/               # WASM frontend components
├── templates/               # Server-side templates (if needed)
├── static/                  # Static files (CSS, JS, images)
├── migrations/              # Database migrations
├── docs/
│   └── requirements_definition.md
├── Cargo.toml
└── README.md
```

## Module System Requirements

- Use `module.rs` + `module/` directory structure (Rust 2024 Edition)
- Never use deprecated `mod.rs` files
- Use `#[routes]` macro for route registration
- Use `installed_apps!` macro for app discovery

## Code Standards

### Comments & Documentation
- All code comments must be in English
- Minimize `.to_string()` calls; prefer borrowing
- Remove obsolete code immediately without deletion records
- Mark placeholders with `todo!()` or `// TODO:` comments
- Document all `#[allow(...)]` attributes with explanatory comments

### Placeholder Notation
- `todo!()` - features that will be implemented
- `unimplemented!()` - intentionally excluded features (retain permanently)
- `// TODO:` - planning notes
- Delete `todo!()` and `// TODO:` upon implementation

## Reinhardt-Specific Patterns

### Model Definition
```rust
use reinhardt_db::prelude::*;

#[derive(Model)]
#[model(table_name = "books")]
pub struct Book {
    #[pk]
    pub id: Uuid,
    pub title: String,
    pub user_id: Uuid,
    #[auto_now_add]
    pub created_at: DateTime<Utc>,
}
```

### ViewSet (REST API)
```rust
use reinhardt_rest::prelude::*;

#[viewset]
impl BookViewSet {
    type Model = Book;
    type Serializer = BookSerializer;

    #[action(detail = false, methods = ["GET"])]
    async fn list(&self, request: Request) -> Response {
        // ...
    }
}
```

### Dependency Injection
```rust
use reinhardt_di::inject;

#[inject]
async fn get_user(
    auth: AuthService,
    db: DatabaseConnection,
) -> Result<User, Error> {
    // Dependencies automatically resolved
}
```

## Testing Philosophy

Tests must contain meaningful assertions and follow:
- Strict assertions (`assert_eq!`) over loose matching
- Arrange-Act-Assert (AAA) structural pattern
- `#[serial(group_name)]` for tests accessing global state
- Use `reinhardt-test` crate with TestContainers

**Quality Assurance Commands**:
```bash
cargo check --workspace --all-features
cargo build --workspace --all-features
cargo test --workspace --all-features
cargo fmt --check && cargo clippy
```

## File Management

Critical restrictions:
- Never save temporary files to the project directory (use `/tmp`)
- Immediately delete temporary files from `/tmp` when finished
- Immediately remove backup files (`.bak`, `.backup`, `.old`, `~` suffix)

## Git & Release Workflow

**Commits**:
- Require explicit user instruction before committing
- Use Conventional Commits v1.0.0 format
- Split commits by specific intent

**Commit Message Format**:
```
feat: New feature
fix: Bug fix
docs: Documentation changes
style: Code formatting
refactor: Code refactoring
test: Add or modify tests
chore: Build process or tool changes
```

**GitHub Operations**:
- Use GitHub CLI (`gh`) exclusively

## Design Principles

### Let LLM Handle It
- Don't hardcode domain-specific templates
- Pass user's natural language input directly to LLM
- Keep design flexible for new domains

### Keep It Simple
- Only make changes that are directly requested
- Don't over-engineer or add unnecessary abstractions

## Environment Variables

**Required**:
```bash
APP_ENV=development
DATABASE_URL=postgresql://user:password@localhost:5432/hailango
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
```

**External APIs**:
```bash
GOOGLE_CLOUD_VISION_API_KEY=
GOOGLE_CLOUD_TTS_API_KEY=
OPENAI_API_KEY=
STRIPE_SECRET_KEY=
ANTHROPIC_API_KEY=
```

## Key Features to Implement

1. **Book Digitization**: PDF/image upload with OCR
2. **TTS**: Multi-language text-to-speech
3. **STT**: Pronunciation evaluation with scoring
4. **Teacher Mode**: Automated lesson playback
5. **SRS**: Spaced repetition learning algorithm (SM-2)
6. **Offline Support**: PWA with service worker caching

## References

- [Reinhardt Framework](https://github.com/kent8192/reinhardt-web)
- [Reinhardt Docs](https://docs.rs/reinhardt-web)
- [Requirements Definition](docs/requirements_definition.md)
- [SeaQuery Documentation](https://www.sea-ql.org/SeaQuery/)

---

For detailed requirements, consult `docs/requirements_definition.md`.
