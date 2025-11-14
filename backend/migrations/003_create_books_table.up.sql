-- 書籍テーブルの作成
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    target_language VARCHAR(10) NOT NULL,
    native_language VARCHAR(10) NOT NULL,
    reference_language VARCHAR(10),
    cover_image_url TEXT,
    total_pages INTEGER DEFAULT 0,
    processed_pages INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'uploading' CHECK (status IN ('uploading', 'processing', 'ready', 'failed')),
    ocr_status VARCHAR(50) DEFAULT 'pending' CHECK (ocr_status IN ('pending', 'processing', 'completed', 'failed')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- 制約
    CONSTRAINT total_pages_non_negative CHECK (total_pages >= 0),
    CONSTRAINT processed_pages_non_negative CHECK (processed_pages >= 0),
    CONSTRAINT processed_pages_lte_total CHECK (processed_pages <= total_pages)
);

-- インデックス
CREATE INDEX idx_books_user_id ON books(user_id);
CREATE INDEX idx_books_status ON books(status);
CREATE INDEX idx_books_ocr_status ON books(ocr_status);
CREATE INDEX idx_books_created_at ON books(created_at DESC);

-- updated_atの自動更新トリガー
CREATE TRIGGER update_books_updated_at BEFORE UPDATE ON books
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- コメント
COMMENT ON TABLE books IS '書籍情報を管理するテーブル';
COMMENT ON COLUMN books.id IS '書籍ID（UUID）';
COMMENT ON COLUMN books.user_id IS 'ユーザーID（外部キー）';
COMMENT ON COLUMN books.title IS '書籍タイトル';
COMMENT ON COLUMN books.target_language IS '学習先言語（ISO 639-1コード）';
COMMENT ON COLUMN books.native_language IS '母国語（ISO 639-1コード）';
COMMENT ON COLUMN books.reference_language IS '参照言語（本に使用されている言語）';
COMMENT ON COLUMN books.cover_image_url IS '表紙画像URL';
COMMENT ON COLUMN books.total_pages IS '総ページ数';
COMMENT ON COLUMN books.processed_pages IS 'OCR処理済みページ数';
COMMENT ON COLUMN books.status IS '書籍の状態（uploading, processing, ready, failed）';
COMMENT ON COLUMN books.ocr_status IS 'OCR処理の状態（pending, processing, completed, failed）';
COMMENT ON COLUMN books.created_at IS '作成日時';
COMMENT ON COLUMN books.updated_at IS '更新日時';
