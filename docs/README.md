# HaiLanGo Documentation

AI-powered language learning platform that transforms existing textbooks into interactive learning experiences.

## Quick Links

| Document | Description |
|----------|-------------|
| [Requirements Definition](requirements_definition.md) | Complete functional and non-functional requirements |
| [System Architecture](architecture/system_architecture.md) | System design, component diagram, data flow |
| [Database Schema](architecture/database_schema.md) | ER diagram, table definitions, model examples |
| [API Specification](architecture/api_specification.md) | REST API, WebSocket, authentication endpoints |
| [Getting Started](guides/getting_started.md) | Quick start guide for developers |
| [Test Strategy](testing/test_strategy.md) | Testing approach, coverage goals, CI/CD |

## Technology Stack

| Layer | Technology | Description |
|-------|------------|-------------|
| **Framework** | [Reinhardt](https://github.com/kent8192/reinhardt-web) | Rust full-stack composable API framework |
| **ORM** | `reinhardt-db` | SeaQuery + sqlx integration |
| **Frontend** | `reinhardt-pages` | WASM + SSR reactive framework |
| **REST API** | `reinhardt-rest` | ViewSets, Serializers, pagination |
| **Auth** | `reinhardt-auth` | JWT, Token, Session authentication |
| **WebSocket** | `reinhardt-websockets` | Real-time communication |
| **Database** | PostgreSQL | Primary data storage |
| **Cache** | Redis | Session, caching, rate limiting |

## External Services

| Service | Provider | Purpose |
|---------|----------|---------|
| OCR | Google Vision / Azure | Text extraction from images/PDFs |
| TTS | Google Cloud TTS / Azure | Text-to-speech synthesis |
| STT | OpenAI Whisper / Azure | Speech recognition & evaluation |
| LLM | Anthropic Claude | AI tutoring, explanations |
| Payments | Stripe | Subscription management |

## Project Structure

```
HaiLanGo/
├── src/
│   ├── config/           # Settings, routing, app configuration
│   ├── apps/             # Feature modules
│   │   ├── auth/         # Authentication & users
│   │   ├── books/        # Book management & OCR
│   │   ├── learning/     # Learning sessions
│   │   ├── tts/          # Text-to-Speech
│   │   ├── stt/          # Speech-to-Text
│   │   ├── review/       # SRS review system
│   │   └── teacher_mode/ # Automated lesson playback
│   └── pages/            # WASM frontend components
├── migrations/           # Database migrations
├── docs/                 # Documentation (you are here)
└── Cargo.toml
```

## Documentation Index

### Architecture

- **[System Architecture](architecture/system_architecture.md)**
  - Design philosophy
  - System component diagram
  - Reinhardt framework integration
  - Data flow patterns
  - External service integration

- **[Database Schema](architecture/database_schema.md)**
  - Entity-Relationship diagram
  - Table definitions (users, books, pages, vocabulary, srs_schedule)
  - Index strategy
  - reinhardt-db model examples

- **[API Specification](architecture/api_specification.md)**
  - Authentication endpoints
  - Book management API
  - Learning API
  - WebSocket API (Teacher Mode)
  - Rate limiting & quotas

### Guides

- **[Getting Started](guides/getting_started.md)**
  - Prerequisites (Rust, PostgreSQL, Redis, Podman)
  - 5-minute quickstart
  - Configuration files (TOML, .env)
  - Mock mode development
  - Troubleshooting

### Testing

- **[Test Strategy](testing/test_strategy.md)**
  - Test pyramid approach
  - Coverage goals (80%+ for core logic)
  - reinhardt-test usage
  - TestContainers setup
  - Mock API strategy
  - CI/CD integration

## Key Features

1. **Book Digitization** - Upload PDFs/images, extract text via OCR
2. **Text-to-Speech** - Multi-language audio synthesis with speed control
3. **Speech-to-Text** - Pronunciation evaluation with scoring (0-100)
4. **Teacher Mode** - Automated lesson playback with background audio
5. **SRS Review** - Spaced repetition algorithm (SM-2) for vocabulary
6. **Offline Support** - PWA with service worker caching

## Development Commands

```bash
# Check and build
cargo check --workspace --all-features
cargo build --workspace --all-features

# Run tests
cargo test --workspace --all-features

# Code quality
cargo fmt --check
cargo clippy

# Run development server
cargo run --release
```

## Contributing

See [CLAUDE.md](../CLAUDE.md) for coding guidelines and conventions.

## License

This project is proprietary. All rights reserved.
