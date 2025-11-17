-- Drop review tables
DROP INDEX IF EXISTS idx_review_history_reviewed_at;
DROP INDEX IF EXISTS idx_review_history_item_id;
DROP INDEX IF EXISTS idx_review_history_user_id;
DROP INDEX IF EXISTS idx_review_items_next_review;
DROP INDEX IF EXISTS idx_review_items_book_id;
DROP INDEX IF EXISTS idx_review_items_user_id;

DROP TABLE IF EXISTS review_history;
DROP TABLE IF EXISTS review_items;
