# HaiLanGo Implementation Status

Last updated: 2026-02-24

## Test Suite

| Suite | Count | Status |
|-------|-------|--------|
| Unit tests (`cargo test --lib`) | 322 | ✅ All passing |
| Integration tests (`tests/api_integration.rs`) | 19 | ✅ All passing |
| Database tests (`tests/database_integration.rs`) | 4 | ⏭️ Skipped (require Docker) |
| E2E tests (`e2e/`) | Playwright suite | 🔧 Requires running server |

## API Endpoints

| Module | Endpoint | Status |
|--------|----------|--------|
| Auth | `POST /api/auth/register/` | ✅ |
| Auth | `POST /api/auth/login/` | ✅ |
| Auth | `GET /api/auth/oauth/{provider}/` | ✅ |
| Auth | `GET /api/auth/oauth/{provider}/callback/` | ✅ |
| Auth | `GET /api/auth/providers/` | ✅ |
| Books | `GET/POST /api/books/` | ✅ |
| Books | `GET /api/books/{id}/` | ✅ |
| Learning | `GET/POST /api/learning/sessions/` | ✅ |
| Learning | `GET /api/learning/sessions/{id}/` | ✅ |
| Learning | `GET /api/learning/sessions/{id}/status/` | ✅ |
| Review | `GET/POST /api/review/vocabulary/` | ✅ |
| Review | `GET /api/review/queue/` | ✅ |
| Review | `POST /api/review/record/` | ✅ |
| Review | `GET /api/review/stats/` | ✅ |
| TTS | `POST /api/tts/synthesize/` | ✅ |
| TTS | `GET /api/tts/languages/` | ✅ |
| TTS | `GET /api/tts/history/` | ✅ |
| STT | `POST /api/stt/evaluate/` | ✅ |
| STT | `POST /api/stt/transcribe/` | ✅ |
| STT | `GET /api/stt/attempts/` | ✅ |
| STT | `GET /api/stt/stats/` | ✅ |
| Teacher | `POST /api/teacher/start/` | ✅ |
| Teacher | `POST /api/teacher/pause/` | ✅ |
| Teacher | `POST /api/teacher/resume/` | ✅ |
| Teacher | `POST /api/teacher/stop/` | ✅ |
| Teacher | `POST /api/teacher/next/` | ✅ |
| Teacher | `GET /api/teacher/sessions/` | ✅ |
| Teacher | `GET /api/teacher/sessions/{id}/status/` | ✅ |

## Core Features

| Feature | Implementation | Notes |
|---------|----------------|-------|
| Book digitization (OCR) | Mock + service abstraction | Real: Google Vision / Azure |
| TTS synthesis | Mock + service abstraction | Real: edge-tts / Google / Azure |
| STT pronunciation evaluation | Mock + service abstraction | Real: Whisper / Azure |
| Teacher mode (auto-lesson) | ✅ Fully implemented | SM-2 scheduling |
| SRS (spaced repetition) | ✅ Fully implemented | SM-2 algorithm |
| JWT authentication | ✅ Fully implemented | Access + refresh tokens |
| OAuth (Google/GitHub) | ✅ Fully implemented | PKCE + state validation |
| PWA manifest | ✅ Exists | Service worker: TODO |
| Offline audio caching | Partial | Cache service exists |

## External Service Integration

| Service | Code Layer | Real API |
|---------|-----------|----------|
| OpenAI Whisper (STT) | `src/services/stt.rs` | Configured via `STT_PROVIDER=whisper` |
| edge-tts / Google TTS | `src/services/tts.rs` | Configured via `TTS_PROVIDER` |
| Google Vision (OCR) | `src/services/ocr.rs` | Configured via `OCR_PROVIDER` |
| Anthropic Claude (LLM) | `src/services/llm.rs` | Via `ANTHROPIC_API_KEY` |

## Known Gaps

- Service worker for PWA offline support not yet implemented
- Real API integration tests require valid API keys in `.env`
- Mobile app (Flutter) is a future consideration (Issue #24)
