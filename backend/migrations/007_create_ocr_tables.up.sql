-- OCR Jobs Table
CREATE TABLE IF NOT EXISTS ocr_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    book_id UUID NOT NULL,
    page_number INTEGER,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    progress INTEGER NOT NULL DEFAULT 0,
    total_pages INTEGER,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- OCR Results Table
CREATE TABLE IF NOT EXISTS ocr_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL,
    page_number INTEGER NOT NULL,
    original_text TEXT NOT NULL,
    translated_text TEXT,
    confidence DECIMAL(5,4),
    processing_time_ms INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (job_id) REFERENCES ocr_jobs(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_ocr_jobs_user_id ON ocr_jobs(user_id);
CREATE INDEX idx_ocr_jobs_book_id ON ocr_jobs(book_id);
CREATE INDEX idx_ocr_jobs_status ON ocr_jobs(status);
CREATE INDEX idx_ocr_jobs_created_at ON ocr_jobs(created_at);
CREATE INDEX idx_ocr_results_job_id ON ocr_results(job_id);
