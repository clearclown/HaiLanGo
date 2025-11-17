-- Learning Sessions Table
CREATE TABLE IF NOT EXISTS learning_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    book_id UUID NOT NULL,
    page_number INTEGER NOT NULL,
    duration_seconds INTEGER NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- Page Completions Table
CREATE TABLE IF NOT EXISTS page_completions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    book_id UUID NOT NULL,
    page_number INTEGER NOT NULL,
    completed_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Prevent duplicate completions
    UNIQUE(user_id, book_id, page_number),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- Book Progress Table
CREATE TABLE IF NOT EXISTS book_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    book_id UUID NOT NULL,
    total_pages INTEGER NOT NULL DEFAULT 0,
    completed_pages INTEGER NOT NULL DEFAULT 0,
    last_page_number INTEGER,
    total_time_seconds INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- One progress record per user per book
    UNIQUE(user_id, book_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_learning_sessions_user_id ON learning_sessions(user_id);
CREATE INDEX idx_learning_sessions_book_id ON learning_sessions(book_id);
CREATE INDEX idx_learning_sessions_created_at ON learning_sessions(created_at);
CREATE INDEX idx_page_completions_user_id ON page_completions(user_id);
CREATE INDEX idx_page_completions_book_id ON page_completions(book_id);
CREATE INDEX idx_book_progress_user_id ON book_progress(user_id);
CREATE INDEX idx_book_progress_book_id ON book_progress(book_id);
