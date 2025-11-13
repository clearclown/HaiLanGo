-- ユーザーテーブルの作成
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),  -- OAuth認証の場合はNULL
    display_name VARCHAR(100) NOT NULL,
    profile_image_url TEXT,
    oauth_provider VARCHAR(50),  -- 'google', 'github', 'apple' など
    oauth_provider_id VARCHAR(255),  -- OAuthプロバイダーのユーザーID
    email_verified BOOLEAN DEFAULT FALSE,
    email_verification_token VARCHAR(255),
    email_verification_expires_at TIMESTAMP,
    password_reset_token VARCHAR(255),
    password_reset_expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,  -- ソフトデリート用

    -- 制約
    CONSTRAINT email_valid CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT oauth_user_check CHECK (
        (oauth_provider IS NOT NULL AND oauth_provider_id IS NOT NULL) OR
        (password_hash IS NOT NULL)
    )
);

-- インデックス
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_oauth ON users(oauth_provider, oauth_provider_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_email_verification_token ON users(email_verification_token) WHERE email_verification_token IS NOT NULL;
CREATE INDEX idx_users_password_reset_token ON users(password_reset_token) WHERE password_reset_token IS NOT NULL;

-- updated_atの自動更新トリガー
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- コメント
COMMENT ON TABLE users IS 'ユーザー情報を管理するテーブル';
COMMENT ON COLUMN users.id IS 'ユーザーID（UUID）';
COMMENT ON COLUMN users.email IS 'メールアドレス（ユニーク）';
COMMENT ON COLUMN users.password_hash IS 'パスワードのハッシュ（bcrypt）、OAuth認証の場合はNULL';
COMMENT ON COLUMN users.display_name IS '表示名';
COMMENT ON COLUMN users.profile_image_url IS 'プロフィール画像URL';
COMMENT ON COLUMN users.oauth_provider IS 'OAuthプロバイダー（google, github, appleなど）';
COMMENT ON COLUMN users.oauth_provider_id IS 'OAuthプロバイダーのユーザーID';
COMMENT ON COLUMN users.email_verified IS 'メール認証済みフラグ';
COMMENT ON COLUMN users.email_verification_token IS 'メール認証トークン';
COMMENT ON COLUMN users.email_verification_expires_at IS 'メール認証トークンの有効期限';
COMMENT ON COLUMN users.password_reset_token IS 'パスワードリセットトークン';
COMMENT ON COLUMN users.password_reset_expires_at IS 'パスワードリセットトークンの有効期限';
COMMENT ON COLUMN users.created_at IS '作成日時';
COMMENT ON COLUMN users.updated_at IS '更新日時';
COMMENT ON COLUMN users.deleted_at IS '削除日時（ソフトデリート）';
