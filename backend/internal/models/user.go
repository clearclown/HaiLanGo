package models

import (
	"time"

	"github.com/google/uuid"
)

// User はユーザー情報を表すモデル
type User struct {
	ID                        uuid.UUID  `json:"id" db:"id"`
	Email                     string     `json:"email" db:"email"`
	PasswordHash              *string    `json:"-" db:"password_hash"` // JSONには含めない
	DisplayName               string     `json:"display_name" db:"display_name"`
	ProfileImageURL           *string    `json:"profile_image_url,omitempty" db:"profile_image_url"`
	OAuthProvider             *string    `json:"oauth_provider,omitempty" db:"oauth_provider"`
	OAuthProviderID           *string    `json:"oauth_provider_id,omitempty" db:"oauth_provider_id"`
	EmailVerified             bool       `json:"email_verified" db:"email_verified"`
	EmailVerificationToken    *string    `json:"-" db:"email_verification_token"`
	EmailVerificationExpiresAt *time.Time `json:"-" db:"email_verification_expires_at"`
	PasswordResetToken        *string    `json:"-" db:"password_reset_token"`
	PasswordResetExpiresAt    *time.Time `json:"-" db:"password_reset_expires_at"`
	CreatedAt                 time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                 *time.Time `json:"-" db:"deleted_at"`
}

// RefreshToken はリフレッシュトークンを表すモデル
type RefreshToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}

// CreateUserRequest はユーザー登録リクエスト
type CreateUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=100"`
}

// LoginRequest はログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// OAuthLoginRequest はOAuthログインリクエスト
type OAuthLoginRequest struct {
	Provider string `json:"provider" binding:"required,oneof=google github apple"`
	Code     string `json:"code" binding:"required"`
}

// RefreshTokenRequest はトークンリフレッシュリクエスト
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// PasswordResetRequest はパスワードリセットリクエスト
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetConfirmRequest はパスワードリセット確認リクエスト
type PasswordResetConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// AuthResponse は認証レスポンス
type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // 秒単位
}
