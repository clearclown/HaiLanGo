-- Dictionary Cache Table
CREATE TABLE IF NOT EXISTS dictionary_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    word VARCHAR(500) NOT NULL,
    language VARCHAR(10) NOT NULL,
    phonetics JSONB,
    meanings JSONB NOT NULL,
    source_api VARCHAR(100),
    fetched_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_accessed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    access_count INTEGER NOT NULL DEFAULT 0,

    -- Unique constraint for word+language
    UNIQUE(word, language)
);

-- Supported Languages Table
CREATE TABLE IF NOT EXISTS dictionary_supported_languages (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Insert default supported languages
INSERT INTO dictionary_supported_languages (code, name) VALUES
    ('en', 'English'),
    ('ja', 'Japanese'),
    ('zh', 'Chinese'),
    ('ru', 'Russian'),
    ('es', 'Spanish'),
    ('fr', 'French'),
    ('de', 'German'),
    ('it', 'Italian'),
    ('pt', 'Portuguese'),
    ('ar', 'Arabic'),
    ('he', 'Hebrew'),
    ('fa', 'Persian')
ON CONFLICT (code) DO NOTHING;

-- Indexes for performance
CREATE INDEX idx_dictionary_cache_word_lang ON dictionary_cache(word, language);
CREATE INDEX idx_dictionary_cache_last_accessed ON dictionary_cache(last_accessed_at);
