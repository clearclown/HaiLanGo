-- Drop learning tables
DROP INDEX IF EXISTS idx_book_progress_book_id;
DROP INDEX IF EXISTS idx_book_progress_user_id;
DROP INDEX IF EXISTS idx_page_completions_book_id;
DROP INDEX IF EXISTS idx_page_completions_user_id;
DROP INDEX IF EXISTS idx_learning_sessions_created_at;
DROP INDEX IF EXISTS idx_learning_sessions_book_id;
DROP INDEX IF EXISTS idx_learning_sessions_user_id;

DROP TABLE IF EXISTS book_progress;
DROP TABLE IF EXISTS page_completions;
DROP TABLE IF EXISTS learning_sessions;
