-- HaiLanGo Initial Schema Migration
-- Creates all core tables for the language learning platform
-- Aligned with Rust model definitions in src/apps/*/models.rs

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- Users (src/apps/auth/models.rs::User)
-- ============================================================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    display_name VARCHAR(100) NOT NULL,
    native_language VARCHAR(10) NOT NULL DEFAULT 'en',
    avatar_url TEXT,
    oauth_provider VARCHAR(50),
    oauth_id VARCHAR(255),
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_oauth ON users(oauth_provider, oauth_id) WHERE oauth_provider IS NOT NULL;

-- ============================================================================
-- Books (src/apps/books/models.rs::Book)
-- ============================================================================
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    source_language VARCHAR(10) NOT NULL,
    target_language VARCHAR(10) NOT NULL,
    reference_language VARCHAR(10),
    total_pages INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    encryption_key_hash VARCHAR(255),
    settings JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_books_user ON books(user_id);
CREATE INDEX idx_books_status ON books(status);

-- ============================================================================
-- Pages (src/apps/books/models.rs::Page)
-- ============================================================================
CREATE TABLE IF NOT EXISTS pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    page_number INTEGER NOT NULL,
    original_content TEXT,
    processed_content TEXT,
    layout_data JSONB,
    audio_url TEXT,
    is_processed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(book_id, page_number)
);

CREATE INDEX idx_pages_book ON pages(book_id);

-- ============================================================================
-- Vocabulary (src/apps/review/models.rs::Vocabulary)
-- ============================================================================
CREATE TABLE IF NOT EXISTS vocabularies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    word VARCHAR(255) NOT NULL,
    reading VARCHAR(255),
    meaning TEXT NOT NULL,
    part_of_speech VARCHAR(50),
    example_sentence TEXT,
    frequency INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, word, page_id)
);

CREATE INDEX idx_vocabularies_user ON vocabularies(user_id);
CREATE INDEX idx_vocabularies_page ON vocabularies(page_id);
CREATE INDEX idx_vocabularies_word ON vocabularies(word);

-- ============================================================================
-- SRS Schedule (src/apps/review/models.rs::SrsSchedule)
-- ============================================================================
CREATE TABLE IF NOT EXISTS srs_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vocabulary_id UUID NOT NULL REFERENCES vocabularies(id) ON DELETE CASCADE,
    next_review_date DATE NOT NULL DEFAULT CURRENT_DATE,
    interval_days INTEGER NOT NULL DEFAULT 1,
    easiness_factor REAL NOT NULL DEFAULT 2.5,
    repetitions INTEGER NOT NULL DEFAULT 0,
    correct_count INTEGER NOT NULL DEFAULT 0,
    incorrect_count INTEGER NOT NULL DEFAULT 0,
    last_reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(vocabulary_id, user_id)
);

CREATE INDEX idx_srs_user_next ON srs_schedules(user_id, next_review_date);

-- ============================================================================
-- Learning Sessions (src/apps/learning/models.rs::LearningSession)
-- ============================================================================
CREATE TABLE IF NOT EXISTS learning_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    session_type VARCHAR(20) NOT NULL DEFAULT 'page_by_page',
    start_page INTEGER,
    end_page INTEGER,
    current_page INTEGER NOT NULL DEFAULT 1,
    duration_seconds INTEGER NOT NULL DEFAULT 0,
    settings JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ
);

CREATE INDEX idx_sessions_user ON learning_sessions(user_id);
CREATE INDEX idx_sessions_book ON learning_sessions(book_id);
CREATE INDEX idx_sessions_status ON learning_sessions(status);

-- ============================================================================
-- Learning Progress (src/apps/learning/models.rs::LearningProgress)
-- ============================================================================
CREATE TABLE IF NOT EXISTS learning_progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES learning_sessions(id) ON DELETE CASCADE,
    page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    time_spent_seconds INTEGER NOT NULL DEFAULT 0,
    pronunciation_score INTEGER,
    comprehension_score INTEGER,
    feedback_data JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_progress_session ON learning_progress(session_id);
CREATE INDEX idx_progress_user ON learning_progress(user_id);

-- ============================================================================
-- Teacher Sessions (src/apps/teacher_mode/models.rs::TeacherSession)
-- ============================================================================
CREATE TABLE IF NOT EXISTS teacher_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    start_page INTEGER NOT NULL,
    end_page INTEGER NOT NULL,
    current_page INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'idle',
    config JSONB NOT NULL DEFAULT '{}',
    pages_completed INTEGER NOT NULL DEFAULT 0,
    total_pages INTEGER NOT NULL DEFAULT 0,
    page_playbacks JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ
);

CREATE INDEX idx_teacher_sessions_user ON teacher_sessions(user_id);
CREATE INDEX idx_teacher_sessions_book ON teacher_sessions(book_id);

-- ============================================================================
-- Audio Generations (src/apps/tts/models.rs::AudioGeneration)
-- ============================================================================
CREATE TABLE IF NOT EXISTS audio_generations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    page_id UUID REFERENCES pages(id) ON DELETE SET NULL,
    text TEXT NOT NULL,
    language VARCHAR(10) NOT NULL,
    speed REAL NOT NULL DEFAULT 1.0,
    format VARCHAR(10) NOT NULL DEFAULT 'mp3',
    quality VARCHAR(20) NOT NULL DEFAULT 'standard',
    provider VARCHAR(50) NOT NULL DEFAULT 'mock',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    duration_ms BIGINT,
    audio_size_bytes INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audio_gen_user ON audio_generations(user_id);
CREATE INDEX idx_audio_gen_page ON audio_generations(page_id);

-- ============================================================================
-- Audio Cache (src/apps/tts/models.rs::AudioCache)
-- ============================================================================
CREATE TABLE IF NOT EXISTS audio_cache (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cache_key VARCHAR(255) NOT NULL UNIQUE,
    language VARCHAR(10) NOT NULL,
    format VARCHAR(10) NOT NULL DEFAULT 'mp3',
    quality VARCHAR(20) NOT NULL DEFAULT 'standard',
    audio_size_bytes INTEGER NOT NULL DEFAULT 0,
    duration_ms BIGINT NOT NULL DEFAULT 0,
    hit_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_accessed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audio_cache_key ON audio_cache(cache_key);

-- ============================================================================
-- Pronunciation Attempts (src/apps/stt/models.rs::PronunciationAttempt)
-- ============================================================================
CREATE TABLE IF NOT EXISTS pronunciation_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    page_id UUID REFERENCES pages(id) ON DELETE SET NULL,
    expected_text TEXT NOT NULL,
    recognized_text TEXT,
    language VARCHAR(10) NOT NULL,
    overall_score SMALLINT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    audio_duration_ms BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pronunciation_user ON pronunciation_attempts(user_id);
CREATE INDEX idx_pronunciation_page ON pronunciation_attempts(page_id);

-- ============================================================================
-- Word Feedback (src/apps/stt/models.rs::WordFeedback)
-- ============================================================================
CREATE TABLE IF NOT EXISTS word_feedback (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attempt_id UUID NOT NULL REFERENCES pronunciation_attempts(id) ON DELETE CASCADE,
    word VARCHAR(255) NOT NULL,
    score SMALLINT NOT NULL DEFAULT 0,
    feedback TEXT,
    start_ms BIGINT,
    end_ms BIGINT
);

CREATE INDEX idx_word_feedback_attempt ON word_feedback(attempt_id);

-- ============================================================================
-- Review History (for tracking individual review events)
-- ============================================================================
CREATE TABLE IF NOT EXISTS review_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    schedule_id UUID NOT NULL REFERENCES srs_schedules(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quality INTEGER NOT NULL CHECK (quality >= 0 AND quality <= 5),
    time_spent_ms INTEGER,
    reviewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_review_history_schedule ON review_history(schedule_id);
CREATE INDEX idx_review_history_user_date ON review_history(user_id, reviewed_at);

-- ============================================================================
-- User Statistics
-- ============================================================================
CREATE TABLE IF NOT EXISTS user_statistics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    total_words_learned INTEGER NOT NULL DEFAULT 0,
    total_review_sessions INTEGER NOT NULL DEFAULT 0,
    total_study_time_minutes INTEGER NOT NULL DEFAULT 0,
    current_streak_days INTEGER NOT NULL DEFAULT 0,
    longest_streak_days INTEGER NOT NULL DEFAULT 0,
    last_activity_date DATE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Triggers for updated_at
-- ============================================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_books_updated_at BEFORE UPDATE ON books
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_statistics_updated_at BEFORE UPDATE ON user_statistics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_audio_cache_accessed BEFORE UPDATE ON audio_cache
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
