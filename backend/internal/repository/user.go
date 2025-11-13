package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

var (
	// ErrUserNotFound はユーザーが見つからないエラー
	ErrUserNotFound = errors.New("ユーザーが見つかりません")
	// ErrUserAlreadyExists はユーザーが既に存在するエラー
	ErrUserAlreadyExists = errors.New("ユーザーは既に存在します")
	// ErrRefreshTokenNotFound はリフレッシュトークンが見つからないエラー
	ErrRefreshTokenNotFound = errors.New("リフレッシュトークンが見つかりません")
)

// UserRepository はユーザーリポジトリのインターフェース
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByOAuth(ctx context.Context, provider, providerID string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	SetEmailVerified(ctx context.Context, userID uuid.UUID) error
	SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	VerifyPasswordResetToken(ctx context.Context, token string) (*models.User, error)

	CreateRefreshToken(ctx context.Context, rt *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}

// userRepository はUserRepositoryの実装
type userRepository struct {
	db *sql.DB
}

// NewUserRepository はUserRepositoryの新しいインスタンスを作成
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser は新しいユーザーを作成
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, display_name, oauth_provider, oauth_provider_id, profile_image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at, email_verified
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.OAuthProvider,
		user.OAuthProviderID,
		user.ProfileImageURL,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.EmailVerified)

	if err != nil {
		// 重複メールアドレスのエラーチェック
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

// GetUserByEmail はメールアドレスでユーザーを取得
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, display_name, profile_image_url,
		       oauth_provider, oauth_provider_id, email_verified,
		       email_verification_token, email_verification_expires_at,
		       password_reset_token, password_reset_expires_at,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.ProfileImageURL,
		&user.OAuthProvider,
		&user.OAuthProviderID,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationExpiresAt,
		&user.PasswordResetToken,
		&user.PasswordResetExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID はIDでユーザーを取得
func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, display_name, profile_image_url,
		       oauth_provider, oauth_provider_id, email_verified,
		       email_verification_token, email_verification_expires_at,
		       password_reset_token, password_reset_expires_at,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.ProfileImageURL,
		&user.OAuthProvider,
		&user.OAuthProviderID,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationExpiresAt,
		&user.PasswordResetToken,
		&user.PasswordResetExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetUserByOAuth はOAuthプロバイダーとIDでユーザーを取得
func (r *userRepository) GetUserByOAuth(ctx context.Context, provider, providerID string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, display_name, profile_image_url,
		       oauth_provider, oauth_provider_id, email_verified,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE oauth_provider = $1 AND oauth_provider_id = $2 AND deleted_at IS NULL
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, provider, providerID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.ProfileImageURL,
		&user.OAuthProvider,
		&user.OAuthProviderID,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// UpdateUser はユーザー情報を更新
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET display_name = $1, profile_image_url = $2, updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, user.DisplayName, user.ProfileImageURL, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// DeleteUser はユーザーを削除（ソフトデリート）
func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// SetEmailVerified はメール認証済みフラグを設定
func (r *userRepository) SetEmailVerified(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET email_verified = true, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// SetPasswordResetToken はパスワードリセットトークンを設定
func (r *userRepository) SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	query := `
		UPDATE users
		SET password_reset_token = $1, password_reset_expires_at = $2, updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, token, expiresAt, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// VerifyPasswordResetToken はパスワードリセットトークンを検証してユーザーを取得
func (r *userRepository) VerifyPasswordResetToken(ctx context.Context, token string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, display_name, profile_image_url,
		       oauth_provider, oauth_provider_id, email_verified,
		       password_reset_token, password_reset_expires_at,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE password_reset_token = $1
		  AND password_reset_expires_at > NOW()
		  AND deleted_at IS NULL
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.ProfileImageURL,
		&user.OAuthProvider,
		&user.OAuthProviderID,
		&user.EmailVerified,
		&user.PasswordResetToken,
		&user.PasswordResetExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// CreateRefreshToken はリフレッシュトークンを作成
func (r *userRepository) CreateRefreshToken(ctx context.Context, rt *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query, rt.UserID, rt.Token, rt.ExpiresAt).Scan(&rt.ID, &rt.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// GetRefreshToken はリフレッシュトークンを取得
func (r *userRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, revoked_at
		FROM refresh_tokens
		WHERE token = $1 AND revoked_at IS NULL
	`

	rt := &models.RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.Token,
		&rt.ExpiresAt,
		&rt.CreatedAt,
		&rt.RevokedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return rt, nil
}

// RevokeRefreshToken はリフレッシュトークンを無効化
func (r *userRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token = $1 AND revoked_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}

// RevokeAllUserRefreshTokens はユーザーのすべてのリフレッシュトークンを無効化
func (r *userRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
