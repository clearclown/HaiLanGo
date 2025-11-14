-- ページテーブルの作成
CREATE TABLE IF NOT EXISTS pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    page_number INTEGER NOT NULL,
    image_url TEXT NOT NULL,
    ocr_text TEXT,
    ocr_confidence DECIMAL(5,4) DEFAULT 0.0,
    detected_lang VARCHAR(10),
    ocr_status VARCHAR(50) DEFAULT 'pending' CHECK (ocr_status IN ('pending', 'processing', 'completed', 'failed')),
    ocr_error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- 制約
    CONSTRAINT page_number_positive CHECK (page_number > 0),
    CONSTRAINT ocr_confidence_range CHECK (ocr_confidence >= 0.0 AND ocr_confidence <= 1.0),
    CONSTRAINT unique_book_page UNIQUE (book_id, page_number)
);

-- インデックス
CREATE INDEX idx_pages_book_id ON pages(book_id);
CREATE INDEX idx_pages_book_page ON pages(book_id, page_number);
CREATE INDEX idx_pages_ocr_status ON pages(ocr_status);
CREATE INDEX idx_pages_created_at ON pages(created_at DESC);

-- updated_atの自動更新トリガー
CREATE TRIGGER update_pages_updated_at BEFORE UPDATE ON pages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- コメント
COMMENT ON TABLE pages IS 'ページ情報を管理するテーブル';
COMMENT ON COLUMN pages.id IS 'ページID（UUID）';
COMMENT ON COLUMN pages.book_id IS '書籍ID（外部キー）';
COMMENT ON COLUMN pages.page_number IS 'ページ番号（1から開始）';
COMMENT ON COLUMN pages.image_url IS 'ページ画像URL';
COMMENT ON COLUMN pages.ocr_text IS 'OCR抽出テキスト';
COMMENT ON COLUMN pages.ocr_confidence IS 'OCR信頼度（0.0〜1.0）';
COMMENT ON COLUMN pages.detected_lang IS '検出された言語（ISO 639-1コード）';
COMMENT ON COLUMN pages.ocr_status IS 'OCR処理の状態（pending, processing, completed, failed）';
COMMENT ON COLUMN pages.ocr_error IS 'OCR処理エラーメッセージ';
COMMENT ON COLUMN pages.created_at IS '作成日時';
COMMENT ON COLUMN pages.updated_at IS '更新日時';
