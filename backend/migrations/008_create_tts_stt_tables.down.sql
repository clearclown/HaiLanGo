-- Drop TTS/STT tables
DROP INDEX IF EXISTS idx_stt_results_job_id;
DROP INDEX IF EXISTS idx_stt_jobs_status;
DROP INDEX IF EXISTS idx_stt_jobs_book_id;
DROP INDEX IF EXISTS idx_stt_jobs_user_id;
DROP INDEX IF EXISTS idx_tts_jobs_status;
DROP INDEX IF EXISTS idx_tts_jobs_book_id;
DROP INDEX IF EXISTS idx_tts_jobs_user_id;
DROP INDEX IF EXISTS idx_tts_audio_last_accessed;
DROP INDEX IF EXISTS idx_tts_audio_text_lang;

DROP TABLE IF EXISTS stt_results;
DROP TABLE IF EXISTS stt_jobs;
DROP TABLE IF EXISTS tts_jobs;
DROP TABLE IF EXISTS tts_audio;
