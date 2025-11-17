-- Drop OCR tables
DROP INDEX IF EXISTS idx_ocr_results_job_id;
DROP INDEX IF EXISTS idx_ocr_jobs_created_at;
DROP INDEX IF EXISTS idx_ocr_jobs_status;
DROP INDEX IF EXISTS idx_ocr_jobs_book_id;
DROP INDEX IF EXISTS idx_ocr_jobs_user_id;

DROP TABLE IF EXISTS ocr_results;
DROP TABLE IF EXISTS ocr_jobs;
