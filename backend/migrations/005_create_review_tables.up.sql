-- Review Items Table
CREATE TABLE IF NOT EXISTS review_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    book_id UUID NOT NULL,
    page_number INTEGER NOT NULL,
    item_type VARCHAR(50) NOT NULL CHECK (item_type IN ('word', 'phrase', 'pattern', 'grammar')),
    content TEXT NOT NULL,
    translation TEXT NOT NULL,
    context TEXT,

    -- SRS (Spaced Repetition System) fields
    ease_factor DECIMAL(4,2) NOT NULL DEFAULT 2.50,
    interval INTEGER NOT NULL DEFAULT 0, -- days
    repetitions INTEGER NOT NULL DEFAULT 0,
    next_review_date TIMESTAMP NOT NULL DEFAULT NOW(),
    last_reviewed_at TIMESTAMP,

    -- Performance tracking
    correct_count INTEGER NOT NULL DEFAULT 0,
    incorrect_count INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- Review History Table
CREATE TABLE IF NOT EXISTS review_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    item_id UUID NOT NULL,
    score INTEGER NOT NULL CHECK (score >= 0 AND score <= 100),
    time_spent_seconds INTEGER NOT NULL,
    reviewed_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES review_items(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_review_items_user_id ON review_items(user_id);
CREATE INDEX idx_review_items_book_id ON review_items(book_id);
CREATE INDEX idx_review_items_next_review ON review_items(next_review_date);
CREATE INDEX idx_review_history_user_id ON review_history(user_id);
CREATE INDEX idx_review_history_item_id ON review_history(item_id);
CREATE INDEX idx_review_history_reviewed_at ON review_history(reviewed_at);
