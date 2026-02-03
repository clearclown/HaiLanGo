# Phase 1: Requirements Extraction

## Execution Date
2026-02-02

## Source Documents Analyzed
1. `/home/ablaze/Projects/HaiLanGo/docs/requirements_definition.md` - Complete feature requirements
2. `/home/ablaze/Projects/HaiLanGo/docs/architecture/system_architecture.md` - Reinhardt architecture
3. `/home/ablaze/Projects/HaiLanGo/docs/architecture/database_schema.md` - Data models
4. `/home/ablaze/Projects/HaiLanGo/docs/architecture/api_specification.md` - API contracts
5. `/home/ablaze/Projects/HaiLanGo/docs/testing/test_strategy.md` - Testing approach
6. `/home/ablaze/Projects/HaiLanGo/CLAUDE.md` - Project guidelines

## Core Features Extracted

### 1. Authentication & User Management
- **OAuth 2.0** (Google Login) + Email/Password
- **JWT** token-based session management
- User profile management (display_name, native_language, avatar)
- Email verification
- Password reset
- Account deletion

### 2. Book Digitization (OCR)
- **PDF upload** (primary format)
- **Image upload** (PNG, HEIC→PNG, JPG)
- **Multi-language OCR** support (Google Vision API / Azure Computer Vision)
- Complex layout handling (tables, ruby text, vertical text)
- **E2E encryption** for book content
- User-correctable OCR results
- Language configuration:
  - Native language (UI language)
  - Target language (learning language)
  - Reference language (intermediary language)

### 3. Text-to-Speech (TTS)
- **Priority languages**: Japanese, Chinese, English, Russian, Persian, Hebrew, Spanish, French, Portuguese, German, Italian, Turkish
- Speed adjustment (0.5x-2.0x)
- Quality tiers:
  - Free: Standard quality
  - Premium: High quality (Stripe-billed)
- Pre-generated audio caching
- Batch download for offline use

### 4. Speech-to-Text (STT) & Pronunciation Evaluation
- **Word-level evaluation**
- Score display (0-100)
- Specific improvement feedback (conversational style)
- Pronunciation visualization (waveform)
- Integration with OpenAI Whisper / Azure Speech

### 5. Interactive Learning Modes

#### Page-by-Page Mode
- Study one page at a time
- Phrase repetition practice
- Role-play conversation exercises

#### Teacher Mode (Auto-Lesson Playback) ⭐
- **One-button automated lesson playback**
- Auto-advance through all pages
- Audio read-aloud + explanations per page
- **Background playback** (screen off, lock screen controls)
- **Customizable settings**:
  - Playback speed: 0.5x-2.0x
  - Page interval: 0-30 seconds
  - Repeat count: 1-3 times per page
  - Content options:
    - Target language read-aloud (required)
    - Native language translation (optional)
    - Vocabulary explanations (optional)
    - Grammar explanations (optional)
    - Pronunciation practice time (optional)
- **Offline support**:
  - Batch audio download
  - Wi-Fi-only download recommendation
  - Download management UI
- **Use cases**: Commute learning, bedtime listening, multitasking

### 6. Spaced Repetition System (SRS)
- **SM-2 algorithm** implementation
- Vocabulary scheduling
- Due reviews tracking
- Progress statistics
- Weakness analysis (problematic words)

### 7. Progress Tracking
- Completed pages tracking
- Phrase encounter history
- Achievement visualization
- Continuous study day streak
- Study time analytics

### 8. Subscription & Monetization
- **Free plan**:
  - 1 page/day limit
  - 30 min/day usage limit
  - Standard TTS quality
- **Premium plan** (monthly/yearly):
  - Unlimited learning
  - High-quality TTS
  - Offline audio download
  - Priority support
- **Stripe integration** for payments

## Non-Functional Requirements

### Security
- **E2E encryption** for book content (user-derived keys)
- TLS 1.3 for all connections
- PostgreSQL AES-256 encryption for sensitive fields
- Private use only (no external sharing except UGC)
- GDPR/CCPA compliance (future)

### Performance
- OCR: 5-10 seconds per page (external API dependent)
- TTS latency: <1 second to start
- STT feedback: <3 seconds after recording
- Real-time WebSocket streaming

### Scalability
- Initial: Personal use (self-hosted)
- Future: Multi-tenant SaaS
- Horizontal scaling architecture
- Circuit breaker pattern for external APIs

### Privacy
- Books stored until user deletion
- OCR results cached in Redis, persisted in PostgreSQL
- Personal learning data never shared

## Technology Stack (Mandatory)

### Framework: Reinhardt (Rust Full-Stack)
- **reinhardt-db**: ORM (SeaQuery + sqlx)
- **reinhardt-pages**: WASM + SSR frontend
- **reinhardt-rest**: REST API (ViewSets, Serializers)
- **reinhardt-auth**: JWT, Session authentication
- **reinhardt-websockets**: Real-time communication
- **reinhardt-admin**: Auto-generated admin panel
- **reinhardt-i18n**: Internationalization

### Database
- **PostgreSQL**: Primary database (users, books, learning data)
- **Redis**: Cache, session management, rate limiting

### External APIs
- **OCR**: Google Vision API / Azure Computer Vision
- **TTS**: Google Cloud TTS / Azure Speech
- **STT**: OpenAI Whisper / Azure Speech
- **LLM**: Anthropic Claude API
- **Payments**: Stripe

### Deployment
- **Development**: Podman/Docker Compose
- **Future**: Cloudflare Workers/Pages, AWS Lambda, GCP Cloud Run

## Development Phases

### Phase 1: MVP (3-4 months)
- User authentication (OAuth + Email)
- PDF upload + OCR
- TTS basic functionality (5 major languages)
- Simple vocabulary notebook
- Web-only

### Phase 2: Core Features (2-3 months)
- STT + pronunciation evaluation
- Page-by-page learning mode
- SRS algorithm
- Stripe payment integration

### Phase 3: Extended Features (3-4 months)
- Teacher read-aloud mode (offline support)
- Dictionary API integration
- Learning analytics dashboard
- Minor language expansion

### Phase 4: Community (TBD)
- User-generated content
- Blog platform
- Community forum

## Key Success Metrics

### MVP (Phase 1)
- 100 registered accounts
- 7-day average retention
- 90%+ OCR success rate

### Phase 2
- 5% paid conversion rate
- 500 monthly active users
- 3 hours/week average study time

### Phase 3+
- 10,000 monthly active users
- Revenue covers server costs
- 4.5/5.0+ user satisfaction

## Alignment with Architecture Documents

✅ All extracted features align with:
- System architecture component diagram
- Database schema ER diagram
- API specification endpoints
- Test strategy coverage areas

## Next Steps
Proceed to Phase 2: Structure Definition
