package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateToken はJWTトークンの生成をテスト
func TestGenerateToken(t *testing.T) {
	// テスト用のRSA鍵ペアを生成
	err := GenerateRSAKeys()
	require.NoError(t, err)

	userID := uuid.New().String()
	email := "test@example.com"

	tests := []struct {
		name    string
		userID  string
		email   string
		wantErr bool
	}{
		{
			name:    "正常なトークン生成",
			userID:  userID,
			email:   email,
			wantErr: false,
		},
		{
			name:    "空のユーザーID",
			userID:  "",
			email:   email,
			wantErr: true,
		},
		{
			name:    "空のメールアドレス",
			userID:  userID,
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)
				// JWTトークンは3つのパートに分かれている（header.payload.signature）
				assert.Regexp(t, `^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$`, token)
			}
		})
	}
}

// TestVerifyToken はJWTトークンの検証をテスト
func TestVerifyToken(t *testing.T) {
	// テスト用のRSA鍵ペアを生成
	err := GenerateRSAKeys()
	require.NoError(t, err)

	userID := uuid.New().String()
	email := "test@example.com"
	token, err := GenerateToken(userID, email)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "正しいトークン",
			token:   token,
			wantErr: false,
		},
		{
			name:    "空のトークン",
			token:   "",
			wantErr: true,
		},
		{
			name:    "無効なトークン",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "改ざんされたトークン",
			token:   token + "tampered",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := VerifyToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, userID, claims.UserID)
				assert.Equal(t, email, claims.Email)
				assert.False(t, claims.ExpiresAt.Before(time.Now()))
			}
		})
	}
}

// TestGenerateRefreshToken はリフレッシュトークンの生成をテスト
func TestGenerateRefreshToken(t *testing.T) {
	userID := uuid.New().String()

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "正常なリフレッシュトークン生成",
			userID:  userID,
			wantErr: false,
		},
		{
			name:    "空のユーザーID",
			userID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, expiresAt, err := GenerateRefreshToken(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.False(t, expiresAt.Before(time.Now()))
				// リフレッシュトークンは32文字以上
				assert.GreaterOrEqual(t, len(token), 32)
			}
		})
	}
}

// TestTokenExpiration はトークンの有効期限をテスト
func TestTokenExpiration(t *testing.T) {
	// テスト用のRSA鍵ペアを生成
	err := GenerateRSAKeys()
	require.NoError(t, err)

	// 有効期限を1秒に設定（テスト用）
	originalExpiry := AccessTokenExpiry
	AccessTokenExpiry = 1 * time.Second
	defer func() { AccessTokenExpiry = originalExpiry }()

	userID := uuid.New().String()
	email := "test@example.com"
	token, err := GenerateToken(userID, email)
	require.NoError(t, err)

	// すぐに検証（成功するはず）
	claims, err := VerifyToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)

	// 2秒待機（トークンが期限切れになる）
	time.Sleep(2 * time.Second)

	// 期限切れトークンの検証（失敗するはず）
	claims, err = VerifyToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestTokenClaims はトークンのクレーム情報をテスト
func TestTokenClaims(t *testing.T) {
	// テスト用のRSA鍵ペアを生成
	err := GenerateRSAKeys()
	require.NoError(t, err)

	userID := uuid.New().String()
	email := "test@example.com"
	token, err := GenerateToken(userID, email)
	require.NoError(t, err)

	claims, err := VerifyToken(token)
	require.NoError(t, err)

	// クレーム情報の検証
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.ExpiresAt)
	assert.True(t, claims.ExpiresAt.After(claims.IssuedAt.Time))
}

// TestRefreshTokenUniqueness はリフレッシュトークンの一意性をテスト
func TestRefreshTokenUniqueness(t *testing.T) {
	userID := uuid.New().String()

	token1, _, err1 := GenerateRefreshToken(userID)
	require.NoError(t, err1)

	token2, _, err2 := GenerateRefreshToken(userID)
	require.NoError(t, err2)

	// 同じユーザーIDでも異なるトークンが生成される
	assert.NotEqual(t, token1, token2)
}

// BenchmarkGenerateToken はGenerateTokenのベンチマーク
func BenchmarkGenerateToken(b *testing.B) {
	_ = GenerateRSAKeys()
	userID := uuid.New().String()
	email := "test@example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateToken(userID, email)
	}
}

// BenchmarkVerifyToken はVerifyTokenのベンチマーク
func BenchmarkVerifyToken(b *testing.B) {
	_ = GenerateRSAKeys()
	userID := uuid.New().String()
	email := "test@example.com"
	token, _ := GenerateToken(userID, email)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = VerifyToken(token)
	}
}

// BenchmarkGenerateRefreshToken はGenerateRefreshTokenのベンチマーク
func BenchmarkGenerateRefreshToken(b *testing.B) {
	userID := uuid.New().String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = GenerateRefreshToken(userID)
	}
}
