-- Patterns Table
CREATE TABLE IF NOT EXISTS patterns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('greeting', 'question', 'response', 'request', 'confirmation', 'other')),
    pattern TEXT NOT NULL,
    translation TEXT NOT NULL,
    frequency INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- Pattern Examples Table
CREATE TABLE IF NOT EXISTS pattern_examples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pattern_id UUID NOT NULL,
    page_number INTEGER NOT NULL,
    original_text TEXT NOT NULL,
    translated_text TEXT NOT NULL,
    context TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (pattern_id) REFERENCES patterns(id) ON DELETE CASCADE
);

-- Pattern Practice Table
CREATE TABLE IF NOT EXISTS pattern_practices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pattern_id UUID NOT NULL,
    question TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    alternative_answers JSONB,
    difficulty INTEGER NOT NULL DEFAULT 1 CHECK (difficulty >= 1 AND difficulty <= 5),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (pattern_id) REFERENCES patterns(id) ON DELETE CASCADE
);

-- Pattern Progress Table
CREATE TABLE IF NOT EXISTS pattern_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    pattern_id UUID NOT NULL,
    mastery_level INTEGER NOT NULL DEFAULT 0 CHECK (mastery_level >= 0 AND mastery_level <= 100),
    practice_count INTEGER NOT NULL DEFAULT 0,
    correct_count INTEGER NOT NULL DEFAULT 0,
    last_practiced_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- One progress record per user per pattern
    UNIQUE(user_id, pattern_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (pattern_id) REFERENCES patterns(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_patterns_book_id ON patterns(book_id);
CREATE INDEX idx_patterns_type ON patterns(type);
CREATE INDEX idx_patterns_frequency ON patterns(frequency DESC);
CREATE INDEX idx_pattern_examples_pattern_id ON pattern_examples(pattern_id);
CREATE INDEX idx_pattern_practices_pattern_id ON pattern_practices(pattern_id);
CREATE INDEX idx_pattern_practices_difficulty ON pattern_practices(difficulty);
CREATE INDEX idx_pattern_progress_user_id ON pattern_progress(user_id);
CREATE INDEX idx_pattern_progress_pattern_id ON pattern_progress(pattern_id);
CREATE INDEX idx_pattern_progress_mastery ON pattern_progress(mastery_level);
CREATE INDEX idx_pattern_progress_last_practiced ON pattern_progress(last_practiced_at);
