package repository

import (
	"context"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// InMemoryUserRepository はインメモリのユーザーリポジトリ
type InMemoryUserRepository struct {
	users         map[string]*models.User
	refreshTokens map[string]*models.RefreshToken
	mu            sync.RWMutex
}

// NewInMemoryUserRepository は新しいInMemoryUserRepositoryを作成
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:         make(map[string]*models.User),
		refreshTokens: make(map[string]*models.RefreshToken),
	}
}

// CreateUser は新しいユーザーを作成
func (r *InMemoryUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// メールアドレスの重複チェック
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email && existingUser.DeletedAt == nil {
			return ErrUserAlreadyExists
		}
	}

	// UUIDとタイムスタンプの設定
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// ユーザーを保存
	r.users[user.ID.String()] = user

	return nil
}

// GetUserByEmail はメールアドレスでユーザーを取得
func (r *InMemoryUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email && user.DeletedAt == nil {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

// GetUserByID はIDでユーザーを取得
func (r *InMemoryUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id.String()]
	if !exists || user.DeletedAt != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetUserByOAuth はOAuthプロバイダーとIDでユーザーを取得
func (r *InMemoryUserRepository) GetUserByOAuth(ctx context.Context, provider, providerID string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.OAuthProvider != nil && *user.OAuthProvider == provider &&
			user.OAuthProviderID != nil && *user.OAuthProviderID == providerID &&
			user.DeletedAt == nil {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

// UpdateUser はユーザー情報を更新
func (r *InMemoryUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingUser, exists := r.users[user.ID.String()]
	if !exists || existingUser.DeletedAt != nil {
		return ErrUserNotFound
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID.String()] = user

	return nil
}

// DeleteUser はユーザーを削除（ソフトデリート）
func (r *InMemoryUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id.String()]
	if !exists || user.DeletedAt != nil {
		return ErrUserNotFound
	}

	now := time.Now()
	user.DeletedAt = &now
	user.UpdatedAt = now

	return nil
}

// SetEmailVerified はメール認証済みフラグを設定
func (r *InMemoryUserRepository) SetEmailVerified(ctx context.Context, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID.String()]
	if !exists || user.DeletedAt != nil {
		return ErrUserNotFound
	}

	user.EmailVerified = true
	user.UpdatedAt = time.Now()

	return nil
}

// SetPasswordResetToken はパスワードリセットトークンを設定
func (r *InMemoryUserRepository) SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID.String()]
	if !exists || user.DeletedAt != nil {
		return ErrUserNotFound
	}

	user.PasswordResetToken = &token
	user.PasswordResetExpiresAt = &expiresAt
	user.UpdatedAt = time.Now()

	return nil
}

// VerifyPasswordResetToken はパスワードリセットトークンを検証してユーザーを取得
func (r *InMemoryUserRepository) VerifyPasswordResetToken(ctx context.Context, token string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	now := time.Now()
	for _, user := range r.users {
		if user.PasswordResetToken != nil && *user.PasswordResetToken == token &&
			user.PasswordResetExpiresAt != nil && user.PasswordResetExpiresAt.After(now) &&
			user.DeletedAt == nil {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

// CreateRefreshToken はリフレッシュトークンを作成
func (r *InMemoryUserRepository) CreateRefreshToken(ctx context.Context, rt *models.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// UUIDとタイムスタンプの設定
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	rt.CreatedAt = time.Now()

	// トークンを保存
	r.refreshTokens[rt.Token] = rt

	return nil
}

// GetRefreshToken はリフレッシュトークンを取得
func (r *InMemoryUserRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rt, exists := r.refreshTokens[token]
	if !exists || rt.RevokedAt != nil {
		return nil, ErrRefreshTokenNotFound
	}

	return rt, nil
}

// RevokeRefreshToken はリフレッシュトークンを無効化
func (r *InMemoryUserRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	rt, exists := r.refreshTokens[token]
	if !exists || rt.RevokedAt != nil {
		return ErrRefreshTokenNotFound
	}

	now := time.Now()
	rt.RevokedAt = &now

	return nil
}

// RevokeAllUserRefreshTokens はユーザーのすべてのリフレッシュトークンを無効化
func (r *InMemoryUserRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for _, rt := range r.refreshTokens {
		if rt.UserID == userID && rt.RevokedAt == nil {
			rt.RevokedAt = &now
		}
	}

	return nil
}
