-- Drop Pattern tables
DROP INDEX IF EXISTS idx_pattern_progress_last_practiced;
DROP INDEX IF EXISTS idx_pattern_progress_mastery;
DROP INDEX IF EXISTS idx_pattern_progress_pattern_id;
DROP INDEX IF EXISTS idx_pattern_progress_user_id;
DROP INDEX IF EXISTS idx_pattern_practices_difficulty;
DROP INDEX IF EXISTS idx_pattern_practices_pattern_id;
DROP INDEX IF EXISTS idx_pattern_examples_pattern_id;
DROP INDEX IF EXISTS idx_patterns_frequency;
DROP INDEX IF EXISTS idx_patterns_type;
DROP INDEX IF EXISTS idx_patterns_book_id;

DROP TABLE IF EXISTS pattern_progress;
DROP TABLE IF EXISTS pattern_practices;
DROP TABLE IF EXISTS pattern_examples;
DROP TABLE IF EXISTS patterns;
