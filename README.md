<div align="center">

# 📚 HaiLanGo

### AI-Powered Language Learning Platform

**Transform your physical language textbooks into an intelligent, interactive learning experience**

[![Tests](https://github.com/clearclown/HaiLanGo/actions/workflows/test.yml/badge.svg)](https://github.com/clearclown/HaiLanGo/actions/workflows/test.yml)
[![Backend CI](https://github.com/clearclown/HaiLanGo/workflows/Backend%20CI/badge.svg)](https://github.com/clearclown/HaiLanGo/actions/workflows/backend.yml)
[![Frontend CI](https://github.com/clearclown/HaiLanGo/workflows/Frontend%20CI/badge.svg)](https://github.com/clearclown/HaiLanGo/actions/workflows/frontend.yml)

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Next.js](https://img.shields.io/badge/Next.js-14+-000000?style=flat&logo=next.js)](https://nextjs.org)
[![Flutter](https://img.shields.io/badge/Flutter-3.0+-02569B?style=flat&logo=flutter)](https://flutter.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![Storybook](https://img.shields.io/badge/Storybook-10+-FF4785?style=flat&logo=storybook)](https://storybook.js.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=flat&logo=redis)](https://redis.io)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

### 🌍 Supported Languages

| Language | TTS | STT | OCR | | Language | TTS | STT | OCR |
|:--------:|:---:|:---:|:---:|---|:--------:|:---:|:---:|:---:|
| 🇯🇵 Japanese | ✅ | ✅ | ✅ | | 🇪🇸 Spanish | ✅ | ✅ | ✅ |
| 🇬🇧 English | ✅ | ✅ | ✅ | | 🇫🇷 French | ✅ | ✅ | ✅ |
| 🇨🇳 Chinese | ✅ | ✅ | ✅ | | 🇵🇹 Portuguese | ✅ | ✅ | ✅ |
| 🇷🇺 Russian | ✅ | ✅ | ✅ | | 🇩🇪 German | ✅ | ✅ | ✅ |
| 🇮🇷 Persian | ✅ | ✅ | ✅ | | 🇮🇹 Italian | ✅ | ✅ | ✅ |
| 🇮🇱 Hebrew | ✅ | ✅ | ✅ | | 🇹🇷 Turkish | ✅ | ✅ | ✅ |

*Plus many more minor languages with varying accuracy*

---

</div>

## 📸 Screenshots

> **UI Components Available**: View our component library in [Storybook](frontend/web/.storybook/)
>
> Run `cd frontend/web && pnpm storybook` to browse interactive component documentation.

## 💡 What is HaiLanGo?

HaiLanGo is an **AI-powered language learning platform** that breathes new life into your physical language textbooks. Using cutting-edge OCR, TTS, and STT technologies, it transforms static pages into an interactive, personalized learning experience available 24/7.

**Key Features:**
- 📖 **Digitize any language textbook** with AI-OCR (12+ languages)
- 🎧 **AI Teacher Mode**: Automatic continuous playback with background support
- 🗣️ **Pronunciation Evaluation**: Real-time feedback with 0-100 scoring
- 📊 **Spaced Repetition System (SRS)**: Scientifically optimized review scheduling
- 🔒 **Privacy-First**: E2E encryption keeps your data completely private
- 💾 **Database-Free Development**: Full InMemory fallbacks for testing without PostgreSQL
- 🎨 **Modern UI Components**: Storybook-documented design system with ShadCN/UI

## 🎯 Why HaiLanGo?

### The Problem
Traditional language learning apps force you into their curriculum. But what if you already have the perfect textbook that works for you? What if you want to learn a less common language pair that mainstream apps don't support?

### The Solution
HaiLanGo lets you use **ANY language textbook** and enhances it with AI:

✅ **Your Book, Your Pace**: Use textbooks you trust
✅ **AI-Powered Practice**: Get pronunciation feedback anytime
✅ **Automated Learning**: Teacher Mode plays through pages automatically
✅ **Offline Capable**: Download lessons for offline use
✅ **Rare Language Pairs**: Support for Persian↔Japanese, Hebrew↔Chinese, etc.

### Who It's For
- 🎓 **Students** learning languages at school/university
- 💼 **Professionals** preparing for business or travel
- 🌏 **Language Enthusiasts** studying rare language pairs
- 📚 **Self-Learners** who prefer textbooks over apps

## 🚀 Installation

### Prerequisites

```bash
# Required
- Go 1.21+
- Node.js 18+
- pnpm 10+

# Optional (for full features)
- PostgreSQL 15+
- Redis 7+
- Podman or Docker
```

### Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/clearclown/HaiLanGo.git
cd HaiLanGo

# 2. Setup environment (optional - works without API keys!)
cp .env.example .env
# Edit .env to add API keys, or leave USE_MOCK_APIS=true for testing

# 3. Start Backend
cd backend
go mod download
make build
./bin/server

# 4. Start Frontend (in new terminal)
cd frontend/web
pnpm install
pnpm dev

# 5. Open browser
# Navigate to http://localhost:3000
```

### Development with Database (Optional)

```bash
# Start full stack (PostgreSQL, Redis, Backend, Frontend)
docker compose up -d
# Podman: podman-compose up -d

# Run migrations
cd backend
go run cmd/migrate/main.go up
```

**Note**: The application works **fully without a database** using InMemory repositories! Perfect for quick testing and development.

## 🗑️ Uninstall

```bash
# Stop all services
docker compose down
# Podman: podman-compose down

# Remove project directory
cd ..
rm -rf HaiLanGo

# Remove Docker/Podman volumes (optional)
docker volume prune
```

## 📖 Documentation

### Core Documentation
- [📋 Requirements Definition](docs/requirements_definition.md) - Project overview and functional requirements
- [🎨 UI/UX Design Document](docs/ui_ux_design_document.md) - Screen designs and wireframes
- [🎓 Teacher Mode Technical Spec](docs/teacher_mode_technical_spec.md) - Auto-learning mode specifications
- [🧪 Mocking Strategy](docs/mocking_strategy.md) - Test without API keys
- [🔌 API Integration Proposal](docs/api_integration_proposal.md) - External API/tool survey
- [📚 Storybook](frontend/web/.storybook/) - Interactive UI component documentation

### Feature Requirements Documents
Detailed implementation specs for each feature:

| Phase 1 (MVP) | Phase 2 (Core) | Phase 3 (Advanced) |
|:-------------|:---------------|:-------------------|
| [1. User Authentication](docs/featureRDs/1_ユーザー認証.md) | [6. Page-by-Page Learning](docs/featureRDs/6_ページバイページ学習モード.md) | [12. Dictionary API Integration](docs/featureRDs/12_辞書API統合.md) |
| [2. Book Upload](docs/featureRDs/2_書籍アップロード.md) | [7. Teacher Auto-Learning](docs/featureRDs/7_教師モード自動学習.md) | [13. OCR Manual Correction](docs/featureRDs/13_OCR結果手動修正.md) |
| [3. OCR Processing](docs/featureRDs/3_OCR処理.md) | [8. Spaced Repetition (SRS)](docs/featureRDs/8_間隔反復学習SRS.md) | [14. Conversation Patterns](docs/featureRDs/14_会話パターン抽出.md) ✅ |
| [4. TTS Voice Synthesis](docs/featureRDs/4_TTS音声読み上げ.md) | [9. Vocabulary Features](docs/featureRDs/9_単語帳機能.md) | [15. WebSocket Notifications](docs/featureRDs/15_WebSocketリアルタイム通知.md) ✅ |
| [5. STT Pronunciation](docs/featureRDs/5_STT発音評価.md) ✅ | [10. Learning Analytics](docs/featureRDs/10_学習統計ダッシュボード.md) | [16. Home Screen](docs/featureRDs/16_ホーム画面実装.md) |
| [11. Stripe Payment](docs/featureRDs/11_決済統合Stripe.md) ✅ | | [17. Settings Screen](docs/featureRDs/17_設定画面実装.md) |
| | | [18. GitHub CI Setup](docs/featureRDs/18_GitHub_CI設定.md) |

## 🤝 Contributing

We welcome contributions! Here's how to get started:

### Development Workflow

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** your changes: `git commit -m 'feat: Add amazing feature'`
4. **Push** to your branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request against `main`

### Code Style

**Backend (Go)**
```bash
# Format code
gofmt -w .

# Run linter
golangci-lint run

# Run tests
go test ./...
```

**Frontend (TypeScript)**
```bash
# Format & lint with Biome
pnpm run lint
pnpm run format

# Run tests
pnpm test              # Unit & integration (Vitest)
pnpm test:e2e          # E2E tests (Playwright)

# Component development
pnpm storybook         # Start Storybook for component development
```

### Commit Message Format
We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: Add new feature
fix: Bug fix
docs: Documentation changes
style: Code formatting
refactor: Code refactoring
test: Add or modify tests
chore: Build/tool changes
```

### Testing Strategy
- **TDD Approach**: Write tests before implementation
- **Mock System**: Use `USE_MOCK_APIS=true` for testing without API keys
- **CI/CD**: GitHub Actions runs tests automatically
- See [Mocking Strategy](docs/mocking_strategy.md) for details

## 📚 Resources

### Official Links
- [📖 Documentation](docs/) - Complete project documentation
- [🐛 Issue Tracker](https://github.com/clearclown/HaiLanGo/issues) - Report bugs or request features
- [💬 Discussions](https://github.com/clearclown/HaiLanGo/discussions) - Ask questions and share ideas

### Technology Documentation
- [Go Official Docs](https://golang.org/doc/)
- [Next.js Documentation](https://nextjs.org/docs)
- [Flutter Documentation](https://flutter.dev/docs)
- [PostgreSQL Manual](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)

### External APIs
- [Google Cloud Vision](https://cloud.google.com/vision/docs) - OCR
- [Google Cloud TTS](https://cloud.google.com/text-to-speech/docs) - Text-to-Speech
- [Google Cloud STT](https://cloud.google.com/speech-to-text/docs) - Speech-to-Text
- [OpenAI Realtime API](https://platform.openai.com/docs/) - Real-time voice interaction
- [DeepL API](https://www.deepl.com/docs-api) - High-quality translation
- [Stripe API](https://stripe.com/docs/api) - Payment processing

### Inspiration
- **abceed** - UI/UX design inspiration
- **Duolingo** - Gamification ideas
- **Anki** - Spaced repetition algorithm

## 🗓️ Roadmap

### ✅ Completed
- [x] WebSocket real-time notifications
- [x] InMemory repository fallbacks (database-free development)
- [x] STT pronunciation evaluation
- [x] Conversation pattern extraction
- [x] Stripe payment integration
- [x] UI Component Library with Storybook integration
- [x] Dynamic language support system (30+ languages)

### 🚧 Phase 1: MVP (In Progress)
- [ ] User authentication (OAuth + Email)
- [ ] PDF upload + OCR processing
- [ ] TTS basic features (5 major languages)
- [ ] Simple vocabulary features
- [ ] Web version only

### 📋 Phase 2: Core Features
- [ ] Page-by-page learning mode
- [ ] Spaced repetition algorithm
- [ ] Mobile app (Flutter)
- [ ] Full payment integration

### 🔮 Phase 3: Advanced Features
- [ ] Teacher Mode (offline support)
- [ ] Dictionary API integration
- [ ] Learning analytics dashboard
- [ ] Expanded language support

### 🌐 Phase 4: Community (TBD)
- [ ] User-generated content
- [ ] Blog platform
- [ ] Community forum

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/clearclown/HaiLanGo/issues)
- **Email**: support@HaiLanGo.com
- **Discord**: [Community Server](https://discord.gg/HaiLanGo) *(Coming Soon)*

## ⚖️ Legal

### License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2024 HaiLanGo Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

### Third-Party Services

This project uses third-party services that may have their own licenses:
- **Google Cloud APIs**: Subject to Google Cloud Platform Terms of Service
- **OpenAI APIs**: Subject to OpenAI Terms of Use
- **Stripe**: Subject to Stripe Services Agreement
- **DeepL**: Subject to DeepL API Terms

See [API Integration Proposal](docs/api_integration_proposal.md) for full details.

---

<div align="center">

Made with ❤️ by [HaiLanGo Contributors](https://github.com/clearclown/HaiLanGo/graphs/contributors)

⭐ **Star this project if you find it useful!**

[Report Bug](https://github.com/clearclown/HaiLanGo/issues) · [Request Feature](https://github.com/clearclown/HaiLanGo/issues) · [Contribute](CONTRIBUTING.md)

</div>
