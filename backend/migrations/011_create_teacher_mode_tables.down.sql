DROP INDEX IF EXISTS idx_teacher_mode_playback_last_played;
DROP INDEX IF EXISTS idx_teacher_mode_playback_book_id;
DROP INDEX IF EXISTS idx_teacher_mode_playback_user_id;
DROP INDEX IF EXISTS idx_teacher_mode_playback_user_book;
DROP TABLE IF EXISTS teacher_mode_playback_history;

DROP INDEX IF EXISTS idx_teacher_mode_downloads_book_id;
DROP INDEX IF EXISTS idx_teacher_mode_downloads_user_id;
DROP INDEX IF EXISTS idx_teacher_mode_downloads_user_book;
DROP TABLE IF EXISTS teacher_mode_downloads;
