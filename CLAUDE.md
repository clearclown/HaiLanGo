# CLAUDE.md - HaiLanGo Project Guidelines

## Project Overview

HaiLanGo is an AI-powered language learning platform that transforms existing textbooks into interactive learning experiences. The project combines OCR, TTS, STT, and AI technologies to provide personalized language instruction.

**Repository**: https://github.com/clearclown/HaiLanGo

## Core Technology Stack

- **Framework**: [Reinhardt](https://github.com/kent8192/reinhardt-web) - A composable full-stack API framework for Rust
- **Language**: Rust 2024 Edition
- **ORM**: SeaORM with SeaQuery
- **Database**: PostgreSQL + Redis
- **Authentication**: JWT, OAuth 2.0 (Google Login)
- **API**: REST + WebSocket (real-time features) + GraphQL (optional)
- **Templates**: Reinhardt's built-in template engine (Tera/Askama)
- **Frontend**: Server-side rendered HTML with HTMX/Alpine.js for interactivity

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
│   │   ├── settings/        # Environment-specific settings
│   │   ├── urls.rs          # URL routing
│   │   └── apps.rs          # App configuration
│   └── apps/
│       ├── auth/            # Authentication & users
│       ├── books/           # Book management & OCR
│       ├── learning/        # Learning sessions & pages
│       ├── tts/             # Text-to-Speech
│       ├── stt/             # Speech-to-Text & pronunciation
│       ├── review/          # SRS review system
│       └── teacher_mode/    # Automated lesson playback
├── templates/               # HTML templates (Tera/Askama)
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
- See Reinhardt docs for module organization patterns

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
- Never use alternative notations like `FIXME:` for unimplemented features

## Testing Philosophy

Tests must contain meaningful assertions and follow:
- Strict assertions (`assert_eq!`) over loose matching
- Arrange-Act-Assert (AAA) structural pattern
- `#[serial(group_name)]` for tests accessing global state
- Complete cleanup of all test artifacts

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
- Avoid relative paths beyond one level up; prefer absolute paths

## Git & Release Workflow

**Commits**:
- Require explicit user instruction before committing
- Use Conventional Commits v1.0.0 format
- Split commits by specific intent
- Never execute batch commits without confirmation

**Commit Message Format**:
```
feat: New feature
fix: Bug fix
docs: Documentation changes
style: Code formatting (no functional changes)
refactor: Code refactoring
test: Add or modify tests
chore: Build process or tool changes
```

**GitHub Operations**:
- Use GitHub CLI (`gh`) exclusively for all GitHub interactions
- Never use raw `curl` or web browsers when `gh` is available

## Design Principles

### Let LLM Handle It

Modern LLMs are smart enough - avoid over-engineering:

**DON'T**:
- Hardcode domain-specific prompt templates
- Define fixed category lists
- Pre-control what LLM can naturally handle

**DO**:
- Pass user's natural language input directly to LLM
- For UI domain selection, use free input or let LLM generate choices
- Keep design simple and flexible for new domains

### Keep It Simple

- Only make changes that are directly requested
- Don't add features, refactor, or make "improvements" beyond what was asked
- A bug fix doesn't need surrounding code cleaned up
- Three similar lines of code is better than a premature abstraction

## Environment Variables

**Required (Minimum)**:
```bash
APP_ENV=development
DATABASE_URL=postgresql://user:password@localhost:5432/hailango
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
```

**Optional (External APIs)**:
```bash
GOOGLE_CLOUD_VISION_API_KEY=
GOOGLE_CLOUD_TTS_API_KEY=
OPENAI_API_KEY=
STRIPE_SECRET_KEY=
ANTHROPIC_API_KEY=
```

## Container Configuration

- Docker Desktop must be installed and running
- Ensure `DOCKER_HOST` points to Docker sockets, not Podman
- Use TestContainers for integration tests

## Key Features to Implement

1. **Book Digitization**: PDF/image upload with OCR
2. **TTS**: Multi-language text-to-speech
3. **STT**: Pronunciation evaluation with scoring
4. **Teacher Mode**: Automated lesson playback
5. **SRS**: Spaced repetition learning algorithm
6. **Offline Support**: Pre-download audio for offline use (PWA)

## References

- [Reinhardt Framework](https://github.com/kent8192/reinhardt-web)
- [Requirements Definition](docs/requirements_definition.md)
- [SeaORM Documentation](https://www.sea-ql.org/SeaORM/)

---

For detailed requirements, consult `docs/requirements_definition.md`.
