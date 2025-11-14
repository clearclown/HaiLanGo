-- ページテーブルの削除
DROP TRIGGER IF EXISTS update_pages_updated_at ON pages;
DROP INDEX IF EXISTS idx_pages_created_at;
DROP INDEX IF EXISTS idx_pages_ocr_status;
DROP INDEX IF EXISTS idx_pages_book_page;
DROP INDEX IF EXISTS idx_pages_book_id;
DROP TABLE IF EXISTS pages;
