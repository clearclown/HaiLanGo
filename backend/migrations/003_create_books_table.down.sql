-- 書籍テーブルの削除
DROP TRIGGER IF EXISTS update_books_updated_at ON books;
DROP INDEX IF EXISTS idx_books_created_at;
DROP INDEX IF EXISTS idx_books_ocr_status;
DROP INDEX IF EXISTS idx_books_status;
DROP INDEX IF EXISTS idx_books_user_id;
DROP TABLE IF EXISTS books;
