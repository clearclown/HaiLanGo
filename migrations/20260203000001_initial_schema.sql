-- HaiLanGo Initial Schema Migration
-- Creates all core tables for the language learning platform

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
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
CREATE INDEX idx_users_oauth ON users(oauth_provider, oauth_id) WHERE oauth_provider IS NOT NULL;

-- Books table
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255),
    language VARCHAR(10) NOT NULL,
    cover_url TEXT,
    total_pages INTEGER NOT NULL DEFAULT 0,
    file_type VARCHAR(20) NOT NULL DEFAULT 'pdf',
    file_size_bytes BIGINT,
    ocr_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_books_user ON books(user_id);
CREATE INDEX idx_books_ocr_status ON books(ocr_status);

-- Pages table (OCR results)
CREATE TABLE IF NOT EXISTS pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    page_number INTEGER NOT NULL,
    raw_text TEXT,
    cleaned_text TEXT,
    ocr_confidence FLOAT,
    image_url TEXT,
    audio_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(book_id, page_number)
);

CREATE INDEX idx_pages_book ON pages(book_id);

-- Vocabulary table
CREATE TABLE IF NOT EXISTS vocabulary (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    word VARCHAR(255) NOT NULL,
    reading VARCHAR(255),
    meaning TEXT NOT NULL,
    context TEXT,
    source_book_id UUID REFERENCES books(id) ON DELETE SET NULL,
    source_page INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, word)
);

CREATE INDEX idx_vocabulary_user ON vocabulary(user_id);
CREATE INDEX idx_vocabulary_word ON vocabulary(word);

-- SRS Schedule table (Spaced Repetition)
CREATE TABLE IF NOT EXISTS srs_schedule (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vocabulary_id UUID NOT NULL REFERENCES vocabulary(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    next_review TIMESTAMPTZ NOT NULL,
    interval_days INTEGER NOT NULL DEFAULT 1,
    ease_factor FLOAT NOT NULL DEFAULT 2.5,
    repetitions INTEGER NOT NULL DEFAULT 0,
    last_review TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(vocabulary_id, user_id)
);

CREATE INDEX idx_srs_user_next ON srs_schedule(user_id, next_review);

-- Learning Sessions table
CREATE TABLE IF NOT EXISTS learning_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    session_type VARCHAR(20) NOT NULL DEFAULT 'page_by_page',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    current_page INTEGER NOT NULL DEFAULT 1,
    start_page INTEGER NOT NULL DEFAULT 1,
    end_page INTEGER,
    total_time_seconds INTEGER NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user ON learning_sessions(user_id);
CREATE INDEX idx_sessions_book ON learning_sessions(book_id);
CREATE INDEX idx_sessions_status ON learning_sessions(status);

-- Learning Progress table
CREATE TABLE IF NOT EXISTS learning_progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES learning_sessions(id) ON DELETE CASCADE,
    page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    time_spent_seconds INTEGER NOT NULL DEFAULT 0,
    tts_played BOOLEAN NOT NULL DEFAULT FALSE,
    stt_attempted BOOLEAN NOT NULL DEFAULT FALSE,
    stt_score FLOAT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_progress_session ON learning_progress(session_id);

-- Review History table
CREATE TABLE IF NOT EXISTS review_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    schedule_id UUID NOT NULL REFERENCES srs_schedule(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quality INTEGER NOT NULL CHECK (quality >= 0 AND quality <= 5),
    time_spent_ms INTEGER,
    reviewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_review_history_schedule ON review_history(schedule_id);
CREATE INDEX idx_review_history_user_date ON review_history(user_id, reviewed_at);

-- User Statistics table
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

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_books_updated_at BEFORE UPDATE ON books
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_srs_schedule_updated_at BEFORE UPDATE ON srs_schedule
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_statistics_updated_at BEFORE UPDATE ON user_statistics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
