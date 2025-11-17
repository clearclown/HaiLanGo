-- TTS Audio Cache Table
CREATE TABLE IF NOT EXISTS tts_audio (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text TEXT NOT NULL,
    language VARCHAR(10) NOT NULL,
    voice VARCHAR(100),
    speed DECIMAL(3,2) DEFAULT 1.00,
    audio_url TEXT NOT NULL,
    audio_data BYTEA,
    duration_seconds DECIMAL(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_accessed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    access_count INTEGER NOT NULL DEFAULT 0
);

-- TTS Jobs Table
CREATE TABLE IF NOT EXISTS tts_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    book_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    progress INTEGER NOT NULL DEFAULT 0,
    total_items INTEGER,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- STT Jobs Table
CREATE TABLE IF NOT EXISTS stt_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    book_id UUID NOT NULL,
    audio_url TEXT NOT NULL,
    language VARCHAR(10) NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- STT Results Table
CREATE TABLE IF NOT EXISTS stt_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL,
    transcript TEXT NOT NULL,
    confidence DECIMAL(5,4),
    pronunciation_score INTEGER CHECK (pronunciation_score >= 0 AND pronunciation_score <= 100),
    pronunciation_feedback JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (job_id) REFERENCES stt_jobs(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_tts_audio_text_lang ON tts_audio(text, language);
CREATE INDEX idx_tts_audio_last_accessed ON tts_audio(last_accessed_at);
CREATE INDEX idx_tts_jobs_user_id ON tts_jobs(user_id);
CREATE INDEX idx_tts_jobs_book_id ON tts_jobs(book_id);
CREATE INDEX idx_tts_jobs_status ON tts_jobs(status);
CREATE INDEX idx_stt_jobs_user_id ON stt_jobs(user_id);
CREATE INDEX idx_stt_jobs_book_id ON stt_jobs(book_id);
CREATE INDEX idx_stt_jobs_status ON stt_jobs(status);
CREATE INDEX idx_stt_results_job_id ON stt_results(job_id);
