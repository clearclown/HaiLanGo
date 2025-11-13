-- Learning Sessions Table
CREATE TABLE IF NOT EXISTS learning_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    book_id UUID NOT NULL,
    page_id UUID NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    duration INTEGER NOT NULL, -- in seconds
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_learning_sessions_user_id ON learning_sessions(user_id);
CREATE INDEX idx_learning_sessions_start_time ON learning_sessions(start_time);
CREATE INDEX idx_learning_sessions_book_id ON learning_sessions(book_id);

-- Vocabulary Progress Table
CREATE TABLE IF NOT EXISTS vocabulary_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    word VARCHAR(255) NOT NULL,
    language VARCHAR(10) NOT NULL,
    mastery_level INTEGER NOT NULL DEFAULT 0 CHECK (mastery_level >= 0 AND mastery_level <= 100),
    last_reviewed TIMESTAMP NOT NULL DEFAULT NOW(),
    review_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, word, language)
);

CREATE INDEX idx_vocabulary_progress_user_id ON vocabulary_progress(user_id);
CREATE INDEX idx_vocabulary_progress_mastery_level ON vocabulary_progress(mastery_level);

-- Phrase Progress Table
CREATE TABLE IF NOT EXISTS phrase_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    phrase TEXT NOT NULL,
    language VARCHAR(10) NOT NULL,
    mastery_level INTEGER NOT NULL DEFAULT 0 CHECK (mastery_level >= 0 AND mastery_level <= 100),
    last_reviewed TIMESTAMP NOT NULL DEFAULT NOW(),
    review_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_phrase_progress_user_id ON phrase_progress(user_id);
CREATE INDEX idx_phrase_progress_mastery_level ON phrase_progress(mastery_level);

-- Pronunciation Scores Table
CREATE TABLE IF NOT EXISTS pronunciation_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    text TEXT NOT NULL,
    language VARCHAR(10) NOT NULL,
    score FLOAT NOT NULL CHECK (score >= 0 AND score <= 100),
    accuracy FLOAT NOT NULL CHECK (accuracy >= 0 AND accuracy <= 100),
    fluency FLOAT NOT NULL CHECK (fluency >= 0 AND fluency <= 100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_pronunciation_scores_user_id ON pronunciation_scores(user_id);
CREATE INDEX idx_pronunciation_scores_created_at ON pronunciation_scores(created_at);

-- User Streaks Table (for tracking daily study streaks)
CREATE TABLE IF NOT EXISTS user_streaks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    current_streak INTEGER NOT NULL DEFAULT 0,
    longest_streak INTEGER NOT NULL DEFAULT 0,
    last_study_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_streaks_user_id ON user_streaks(user_id);
