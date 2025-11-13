-- リフレッシュトークンテーブルの作成
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMP,  -- トークンが無効化された日時

    -- 制約
    CONSTRAINT token_not_empty CHECK (token != '')
);

-- インデックス
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token) WHERE revoked_at IS NULL;
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at) WHERE revoked_at IS NULL;

-- コメント
COMMENT ON TABLE refresh_tokens IS 'リフレッシュトークンを管理するテーブル';
COMMENT ON COLUMN refresh_tokens.id IS 'トークンID（UUID）';
COMMENT ON COLUMN refresh_tokens.user_id IS 'ユーザーID';
COMMENT ON COLUMN refresh_tokens.token IS 'リフレッシュトークン';
COMMENT ON COLUMN refresh_tokens.expires_at IS '有効期限';
COMMENT ON COLUMN refresh_tokens.created_at IS '作成日時';
COMMENT ON COLUMN refresh_tokens.revoked_at IS '無効化日時';
