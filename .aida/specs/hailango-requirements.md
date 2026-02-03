# HaiLanGo - Requirements Specification

**Project**: HaiLanGo AI Language Learning Platform
**Generated**: 2026-02-02
**Source**: docs/requirements_definition.md
**Framework**: Reinhardt (Rust Full-Stack)

---

## 1. Project Vision

HaiLanGo transforms existing language learning textbooks into interactive, AI-powered learning experiences. The platform combines OCR, TTS, STT, and LLM technologies to provide personalized 24/7 language instruction focused on conversational fluency.

### Core Value Propositions
- Maximum utilization of existing learning materials
- AI-powered personal tutoring available anytime
- Privacy-first architecture with E2E encryption
- Conversation-focused learning (not grammar-heavy)

---

## 2. Functional Requirements

### FR-1: Authentication & User Management

#### FR-1.1: Email/Password Authentication
- **Priority**: P0 (MVP)
- **Description**: Users register and login with email/password
- **Acceptance Criteria**:
  - Email validation (RFC 5322)
  - Password strength: min 8 chars, 1 uppercase, 1 number, 1 special
  - Argon2id password hashing
  - Email verification flow
  - Password reset via email token
- **Reinhardt Components**: `reinhardt-auth`, JWT middleware

#### FR-1.2: OAuth 2.0 (Google Login)
- **Priority**: P0 (MVP)
- **Description**: Users authenticate via Google OAuth
- **Acceptance Criteria**:
  - Google OAuth 2.0 authorization code flow
  - Auto-create user on first OAuth login
  - Link OAuth to existing email account
- **Reinhardt Components**: `reinhardt-auth`, OAuth middleware

#### FR-1.3: Session Management
- **Priority**: P0 (MVP)
- **Description**: JWT-based stateless sessions with refresh tokens
- **Acceptance Criteria**:
  - Access token expiry: 1 hour
  - Refresh token expiry: 30 days
  - Refresh token rotation (invalidate old on use)
  - Logout invalidates refresh token
- **Reinhardt Components**: `reinhardt-auth` JWT + Redis session

#### FR-1.4: User Profile Management
- **Priority**: P0 (MVP)
- **Description**: Users manage profile information
- **Acceptance Criteria**:
  - Update display_name, native_language, avatar_url
  - View account creation date, last login
  - Delete account (cascade delete all user data)
- **API**: `GET/PATCH/DELETE /api/users/me`

---

### FR-2: Book Digitization (OCR)

#### FR-2.1: PDF Upload
- **Priority**: P0 (MVP)
- **Description**: Users upload PDF textbooks for OCR processing
- **Acceptance Criteria**:
  - Max file size: 50MB
  - Support PDF 1.4-2.0 formats
  - Extract pages to individual images
  - Queue OCR job (async processing)
  - Return 202 Accepted with job_id
- **API**: `POST /api/books/upload`
- **External**: Google Vision API / Azure Computer Vision

#### FR-2.2: Image Upload
- **Priority**: P0 (MVP)
- **Description**: Users upload individual page images
- **Acceptance Criteria**:
  - Support PNG, JPG, HEIC (auto-convert HEIC→PNG)
  - Max image size: 10MB
  - Max resolution: 4096x4096
  - Consecutive upload support (multiple pages)
- **API**: `POST /api/books/upload`

#### FR-2.3: OCR Processing
- **Priority**: P0 (MVP)
- **Description**: Extract text from uploaded images with high accuracy
- **Acceptance Criteria**:
  - Support complex layouts (tables, ruby text, vertical text)
  - Multi-language OCR (priority: ja, zh, en, ru, fa, he, es, fr, pt, de, it, tr)
  - Confidence score per page (0-1.0)
  - Flag low-confidence pages (<0.8) for manual review
  - Store original_content and processed_content
- **Database**: `pages` table
- **Cache**: Redis (TTL: 24 hours)

#### FR-2.4: Language Configuration
- **Priority**: P0 (MVP)
- **Description**: Configure language settings per book
- **Acceptance Criteria**:
  - **Native language**: User's UI language (e.g., ja)
  - **Target language**: Language to learn (e.g., ku)
  - **Reference language**: Book's intermediary language (e.g., en)
  - Validate language codes (ISO 639-1)
- **Database**: `books.source_language`, `target_language`, `reference_language`

#### FR-2.5: OCR Manual Correction (Optional)
- **Priority**: P2
- **Description**: Users manually correct OCR errors
- **Acceptance Criteria**:
  - Edit original_content in-place
  - Track correction history
  - Mark page as "user-corrected"
- **API**: `PATCH /api/books/{id}/pages/{page_number}`

#### FR-2.6: E2E Encryption
- **Priority**: P1
- **Description**: Encrypt book content with user-derived keys
- **Acceptance Criteria**:
  - AES-256-GCM encryption
  - User password derives encryption key (PBKDF2)
  - Encrypt original_content, processed_content
  - Key never leaves user device (client-side encryption)
- **Implementation**: `apps/books/encryption.rs`

---

### FR-3: Text-to-Speech (TTS)

#### FR-3.1: Multi-Language TTS
- **Priority**: P0 (MVP - 5 languages), P1 (Full)
- **Description**: Generate natural-sounding audio for text content
- **Priority Languages** (MVP): Japanese, English, Spanish, French, German
- **Full Support**: Chinese, Russian, Persian, Hebrew, Portuguese, Italian, Turkish
- **Acceptance Criteria**:
  - Support all priority languages with native voice
  - Fallback to Google Translate TTS for minor languages
  - Generate audio in MP3 format (192kbps)
- **API**: `GET /api/tts/pages/{page_id}`, `POST /api/tts/synthesize`
- **External**: Google Cloud TTS / Azure Speech

#### FR-3.2: Speed Adjustment
- **Priority**: P0 (MVP)
- **Description**: Adjust playback speed
- **Acceptance Criteria**:
  - Speed range: 0.5x - 2.0x (0.1x increments)
  - Default: 1.0x
  - Persist speed preference per user
- **API**: Query param `?speed=0.8`

#### FR-3.3: Quality Tiers
- **Priority**: P1
- **Description**: Different audio quality for free vs premium
- **Acceptance Criteria**:
  - **Free**: Standard quality (Google Cloud TTS Standard voices)
  - **Premium**: High quality (Google Cloud TTS WaveNet voices)
  - Gate premium quality behind subscription check
- **API**: Query param `?quality=premium`

#### FR-3.4: Audio Caching
- **Priority**: P0 (MVP)
- **Description**: Cache generated audio to reduce API costs
- **Acceptance Criteria**:
  - Cache MP3 files in file storage (local or S3/R2)
  - Store audio_url in pages table
  - TTL: 90 days (regenerate on demand)
- **Database**: `pages.audio_url`

#### FR-3.5: Batch Audio Download (Offline)
- **Priority**: P2
- **Description**: Download all page audio for offline use
- **Acceptance Criteria**:
  - Batch generate audio for page range
  - Return ZIP archive with MP3 files
  - Job status tracking (202 Accepted → polling)
  - Recommend Wi-Fi only download (large files)
- **API**: `POST /api/tts/batch`

---

### FR-4: Speech-to-Text (STT) & Pronunciation Evaluation

#### FR-4.1: Audio Transcription
- **Priority**: P1
- **Description**: Transcribe user's spoken audio
- **Acceptance Criteria**:
  - Accept audio formats: MP3, WAV, WebM
  - Max duration: 30 seconds
  - Return transcription + word timings
  - Support target language only
- **API**: `POST /api/stt/evaluate`
- **External**: OpenAI Whisper / Azure Speech

#### FR-4.2: Pronunciation Scoring (Word-Level)
- **Priority**: P1
- **Description**: Evaluate pronunciation accuracy per word
- **Acceptance Criteria**:
  - Compare transcription to reference text
  - Score per word: 0-100
  - Overall score: weighted average
  - Identify problematic words (<70 score)
- **Algorithm**: Levenshtein distance + phonetic similarity

#### FR-4.3: Detailed Feedback
- **Priority**: P1
- **Description**: Conversational feedback like English teacher
- **Acceptance Criteria**:
  - Use Claude API to generate feedback
  - Focus on intonation, stress, rhythm issues
  - Provide 2-3 specific improvement tips
  - Positive reinforcement for good pronunciation
- **LLM**: Anthropic Claude API

#### FR-4.4: Waveform Visualization
- **Priority**: P2
- **Description**: Visual representation of pronunciation
- **Acceptance Criteria**:
  - Display audio waveform
  - Highlight problematic word regions
  - Show reference vs user waveform comparison
- **Frontend**: `reinhardt-pages` canvas component

---

### FR-5: Learning Modes

#### FR-5.1: Page-by-Page Mode
- **Priority**: P0 (MVP)
- **Description**: Study one page at a time with interactive exercises
- **Acceptance Criteria**:
  - Display page content (original + translation)
  - Play TTS audio on demand
  - Record and evaluate pronunciation
  - Extract vocabulary (auto-generate flashcards)
  - Navigate prev/next pages
  - Track time spent per page
- **API**: `GET /api/books/{id}/pages/{page_number}`

#### FR-5.2: Teacher Mode (Auto-Playback) ⭐
- **Priority**: P2
- **Description**: Automated lesson playback like podcast
- **Acceptance Criteria**:
  - **One-button start**: Begin from current page or page 1
  - **Auto-advance**: Move to next page after interval
  - **Background playback**: Continue with screen off
  - **Lock screen controls**: Play/pause/skip via notification
  - **Customizable settings**:
    - Playback speed: 0.5x - 2.0x
    - Page interval: 0-30 seconds
    - Repeat count: 1-3 times per page
    - Include translation: yes/no
    - Include vocabulary: yes/no
    - Include grammar: yes/no
    - Pronunciation practice pause: yes/no
  - **WebSocket streaming**: Real-time audio + page updates
  - **Commands**: Pause, resume, skip, stop, update_settings
- **API**: `ws://api/ws/teacher/{book_id}`
- **Use Cases**: Commute, bedtime, chores, driving

#### FR-5.3: Session Tracking
- **Priority**: P0 (MVP)
- **Description**: Track learning session progress
- **Acceptance Criteria**:
  - Create session on learning start
  - Record session_type: page_by_page | teacher_mode | review
  - Track start_page, end_page, duration_seconds
  - Mark session status: active | paused | completed | abandoned
  - Calculate total study time per book
- **Database**: `learning_sessions`, `learning_progress`

---

### FR-6: Spaced Repetition System (SRS)

#### FR-6.1: SM-2 Algorithm
- **Priority**: P1
- **Description**: Implement SuperMemo SM-2 for vocabulary review
- **Acceptance Criteria**:
  - Initial interval: 1 day
  - Quality grades: 0-5 (0=blackout, 5=perfect)
  - Grade ≥3: Increase interval
  - Grade <3: Reset to 1 day
  - Easiness factor: 1.3 - 2.5
  - Formula: `interval = previous_interval * easiness_factor`
- **Implementation**: `apps/review/sm2.rs`
- **Database**: `srs_schedules`

#### FR-6.2: Vocabulary Extraction
- **Priority**: P1
- **Description**: Auto-extract vocabulary from pages
- **Acceptance Criteria**:
  - Extract nouns, verbs, adjectives from OCR content
  - Use dictionary API for definitions
  - Store word, reading, meaning, part_of_speech
  - Track frequency across book
- **Database**: `vocabularies`
- **External**: Oxford Dictionary API / Free Dictionary API

#### FR-6.3: Review Due Queue
- **Priority**: P1
- **Description**: Show vocabulary due for review today
- **Acceptance Criteria**:
  - Query srs_schedules WHERE next_review_date <= TODAY
  - Order by next_review_date ASC (oldest first)
  - Limit: 20 items per session
  - Return vocabulary + example sentence
- **API**: `GET /api/review/due?limit=20`

#### FR-6.4: Review Submission
- **Priority**: P1
- **Description**: Submit review result and update schedule
- **Acceptance Criteria**:
  - Accept quality grade (0-5)
  - Update easiness_factor, interval_days, next_review_date
  - Increment correct_count or incorrect_count
  - Update last_reviewed_at timestamp
- **API**: `POST /api/review/{srs_id}/submit`

#### FR-6.5: Manual Vocabulary Addition
- **Priority**: P2
- **Description**: Users manually add words to SRS
- **Acceptance Criteria**:
  - Select word from page content
  - Auto-fetch definition from dictionary API
  - Create vocabulary + srs_schedule
  - Initial interval: 1 day
- **API**: `POST /api/review/add`

---

### FR-7: Progress Tracking & Analytics

#### FR-7.1: Learning Statistics
- **Priority**: P1
- **Description**: Display learning metrics
- **Acceptance Criteria**:
  - Total study time (minutes)
  - Pages completed
  - Current streak (consecutive days)
  - Longest streak
  - Vocabulary learned (SRS items)
  - Average pronunciation score
  - Daily/weekly/monthly breakdown
- **API**: `GET /api/learning/stats?period=week`

#### FR-7.2: Book Progress
- **Priority**: P0 (MVP)
- **Description**: Track progress per book
- **Acceptance Criteria**:
  - Completed pages count
  - Last page read
  - Progress percentage
  - Total study time for book
- **Calculation**: `completed_pages / total_pages * 100`

#### FR-7.3: Weakness Analysis
- **Priority**: P2
- **Description**: Identify problematic words/patterns
- **Acceptance Criteria**:
  - Track words with low pronunciation scores
  - Track SRS items with high incorrect_count
  - Suggest targeted practice
- **Algorithm**: `incorrect_count / (correct_count + incorrect_count) > 0.3`

---

### FR-8: Subscription & Monetization

#### FR-8.1: Free Plan Limits
- **Priority**: P1
- **Description**: Enforce free tier usage limits
- **Acceptance Criteria**:
  - 1 page/day learning limit
  - 30 minutes/day usage limit
  - Standard TTS quality only
  - Reset limits at midnight UTC
- **Implementation**: Redis-based rate limiting

#### FR-8.2: Premium Plan Features
- **Priority**: P1
- **Description**: Unlock premium features via subscription
- **Acceptance Criteria**:
  - Unlimited learning (no page/time limits)
  - High-quality TTS (WaveNet voices)
  - Offline audio download
  - Priority support
  - Plans: monthly ($9.99), yearly ($79.99)
- **External**: Stripe subscription API

#### FR-8.3: Stripe Integration
- **Priority**: P1
- **Description**: Process payments via Stripe
- **Acceptance Criteria**:
  - Create Stripe customer on first subscription
  - Handle subscription lifecycle (active, past_due, canceled)
  - Webhook handling for subscription events
  - Prorated upgrades/downgrades
  - Invoice generation
- **API**: `POST /api/subscriptions/create`, Stripe webhooks
- **Database**: `subscriptions`

#### FR-8.4: Usage Tracking
- **Priority**: P1
- **Description**: Track API usage per user for quota enforcement
- **Acceptance Criteria**:
  - Track OCR pages processed
  - Track TTS minutes generated
  - Track STT evaluations
  - Store daily usage in Redis (TTL: 25 hours)
  - Block requests exceeding quota (429 error)
- **Implementation**: `services/usage_tracker.rs`

---

## 3. Non-Functional Requirements

### NFR-1: Security

#### NFR-1.1: Authentication Security
- JWT secret: 256-bit random key (env var)
- Argon2id password hashing (time cost: 2, memory cost: 19MB)
- TLS 1.3 for all connections
- Rate limiting: 10 login attempts per 15 min per IP

#### NFR-1.2: Data Encryption
- E2E encryption for book content (AES-256-GCM)
- PostgreSQL column encryption for sensitive fields
- User encryption keys never stored on server
- Encryption key derivation: PBKDF2-HMAC-SHA256 (100k iterations)

#### NFR-1.3: Privacy Compliance
- GDPR Article 17: Right to erasure (account deletion)
- CCPA: User data export on request
- Private use only: Books not shared externally
- No telemetry without explicit consent

---

### NFR-2: Performance

#### NFR-2.1: API Response Times
- **Target**: p95 < 200ms for non-external-API endpoints
- **OCR**: 5-10 seconds per page (external API dependent)
- **TTS**: <1 second to start streaming
- **STT**: <3 seconds for evaluation result

#### NFR-2.2: Database Performance
- Index all foreign keys
- Partial indexes for common filters (status = 'active')
- Connection pooling (min: 5, max: 20)
- Query timeout: 5 seconds

#### NFR-2.3: Caching Strategy
- Redis cache for:
  - OCR results (TTL: 24 hours)
  - TTS audio URLs (TTL: 90 days)
  - Session data (TTL: 30 days)
  - Rate limit counters (TTL: 25 hours)

---

### NFR-3: Scalability

#### NFR-3.1: Horizontal Scaling
- Stateless application servers
- Session state in Redis (shared)
- Database connection pooling
- Load balancer with health checks

#### NFR-3.2: External API Resilience
- Circuit breaker pattern (fail after 5 consecutive errors)
- Fallback strategies:
  - OCR: Queue for retry, notify user
  - TTS: Use cached audio or fallback provider
  - STT: Return graceful error, allow retry
- Timeout: 30 seconds per external API call

---

### NFR-4: Observability

#### NFR-4.1: Logging
- Structured logging (JSON format)
- Log levels: ERROR, WARN, INFO, DEBUG
- Log aggregation: ELK stack or Loki (future)
- Sensitive data: Never log passwords, tokens, encryption keys

#### NFR-4.2: Monitoring
- Prometheus metrics:
  - Request rate, latency, error rate
  - Database query duration
  - External API call duration
  - Active WebSocket connections
- Grafana dashboards (future)

#### NFR-4.3: Error Tracking
- Sentry integration for exception tracking
- Group errors by endpoint + error type
- Alert on error rate spike (>1% of requests)

---

## 4. Development Phases & Priorities

### Phase 1: MVP (P0 - 3-4 months)
- [ ] FR-1.1, FR-1.2, FR-1.3, FR-1.4: Auth (OAuth + Email)
- [ ] FR-2.1, FR-2.2, FR-2.3, FR-2.4: OCR (PDF + Image)
- [ ] FR-3.1, FR-3.2: TTS (5 languages, speed adjustment)
- [ ] FR-5.1: Page-by-Page mode
- [ ] FR-7.2: Book progress tracking
- [ ] NFR-1, NFR-2: Security + Performance baseline

### Phase 2: Core Features (P1 - 2-3 months)
- [ ] FR-4: STT + Pronunciation evaluation
- [ ] FR-6: SRS system (SM-2 algorithm)
- [ ] FR-8: Stripe subscriptions
- [ ] FR-7.1: Learning statistics
- [ ] FR-3.3: TTS quality tiers
- [ ] FR-2.6: E2E encryption

### Phase 3: Extended Features (P2 - 3-4 months)
- [ ] FR-5.2: Teacher Mode (WebSocket streaming)
- [ ] FR-3.5: Batch audio download (offline)
- [ ] FR-7.3: Weakness analysis
- [ ] FR-4.4: Waveform visualization
- [ ] FR-6.5: Manual vocabulary addition

### Phase 4: Community (TBD)
- User-generated content
- Blog platform
- Community forum

---

## 5. Success Criteria

### MVP Launch (Phase 1)
- 100 registered accounts
- 90%+ OCR success rate (confidence >0.8)
- 7-day average user retention
- <500ms p95 API latency

### Growth (Phase 2)
- 500 monthly active users
- 5% paid conversion rate
- 3 hours/week average study time
- <1% error rate

### Maturity (Phase 3+)
- 10,000 monthly active users
- Revenue covers infrastructure costs
- 4.5/5.0 user satisfaction score
- Support 15+ languages

---

## 6. Technology Stack Summary

### Reinhardt Framework
- **reinhardt-db**: PostgreSQL ORM (SeaQuery + sqlx)
- **reinhardt-pages**: WASM + SSR frontend
- **reinhardt-rest**: REST API (ViewSets)
- **reinhardt-auth**: JWT + Session authentication
- **reinhardt-websockets**: Teacher Mode streaming
- **reinhardt-admin**: Admin panel

### Infrastructure
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Deployment**: Podman (dev), Cloud (prod)

### External APIs
- **OCR**: Google Vision API / Azure Computer Vision
- **TTS**: Google Cloud TTS / Azure Speech
- **STT**: OpenAI Whisper / Azure Speech
- **LLM**: Anthropic Claude API
- **Payments**: Stripe

---

## References
- [requirements_definition.md](../docs/requirements_definition.md)
- [system_architecture.md](../docs/architecture/system_architecture.md)
- [database_schema.md](../docs/architecture/database_schema.md)
- [api_specification.md](../docs/architecture/api_specification.md)
- [CLAUDE.md](../CLAUDE.md)
