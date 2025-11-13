package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHashPassword はパスワードのハッシュ化をテスト
func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "正常なパスワード",
			password: "TestPassword123!",
			wantErr:  false,
		},
		{
			name:     "長いパスワード",
			password: "ThisIsAVeryLongPasswordWithManyCharacters123!@#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "短いパスワード",
			password: "Test123",
			wantErr:  false,
		},
		{
			name:     "空のパスワード",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, hash)
				// ハッシュは元のパスワードと異なることを確認
				assert.NotEqual(t, tt.password, hash)
				// ハッシュは60文字であることを確認（bcryptの仕様）
				assert.Len(t, hash, 60)
				// ハッシュが$2a$で始まることを確認（bcryptの識別子）
				assert.Contains(t, hash, "$2a$")
			}
		})
	}
}

// TestVerifyPassword はパスワードの検証をテスト
func TestVerifyPassword(t *testing.T) {
	password := "TestPassword123!"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "正しいパスワード",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "間違ったパスワード",
			password: "WrongPassword",
			hash:     hash,
			want:     false,
		},
		{
			name:     "空のパスワード",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "大文字小文字が違う",
			password: "testpassword123!",
			hash:     hash,
			want:     false,
		},
		{
			name:     "空のハッシュ",
			password: password,
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VerifyPassword(tt.password, tt.hash)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestHashPasswordDeterminism は同じパスワードでも異なるハッシュが生成されることをテスト
func TestHashPasswordDeterminism(t *testing.T) {
	password := "TestPassword123!"

	hash1, err1 := HashPassword(password)
	require.NoError(t, err1)

	hash2, err2 := HashPassword(password)
	require.NoError(t, err2)

	// 同じパスワードでも異なるハッシュが生成される（ソルト使用）
	assert.NotEqual(t, hash1, hash2)

	// しかし、両方のハッシュで検証は成功する
	assert.True(t, VerifyPassword(password, hash1))
	assert.True(t, VerifyPassword(password, hash2))
}

// TestValidatePasswordStrength はパスワード強度の検証をテスト
func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "強力なパスワード",
			password: "TestPassword123!",
			wantErr:  false,
		},
		{
			name:     "大文字・小文字・数字・記号すべて含む",
			password: "Abc123!@#",
			wantErr:  false,
		},
		{
			name:     "短すぎるパスワード（8文字未満）",
			password: "Test12!",
			wantErr:  true,
			errMsg:   "パスワードは最低8文字必要です",
		},
		{
			name:     "大文字のみ",
			password: "TESTPASSWORD",
			wantErr:  true,
			errMsg:   "パスワードは大文字、小文字、数字、記号のうち3種類以上を含む必要があります",
		},
		{
			name:     "小文字のみ",
			password: "testpassword",
			wantErr:  true,
			errMsg:   "パスワードは大文字、小文字、数字、記号のうち3種類以上を含む必要があります",
		},
		{
			name:     "大文字と小文字のみ",
			password: "TestPassword",
			wantErr:  true,
			errMsg:   "パスワードは大文字、小文字、数字、記号のうち3種類以上を含む必要があります",
		},
		{
			name:     "大文字と数字のみ",
			password: "TESTPASSWORD123",
			wantErr:  true,
			errMsg:   "パスワードは大文字、小文字、数字、記号のうち3種類以上を含む必要があります",
		},
		{
			name:     "空のパスワード",
			password: "",
			wantErr:  true,
			errMsg:   "パスワードは最低8文字必要です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// BenchmarkHashPassword はHashPasswordのベンチマーク
func BenchmarkHashPassword(b *testing.B) {
	password := "TestPassword123!"
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

// BenchmarkVerifyPassword はVerifyPasswordのベンチマーク
func BenchmarkVerifyPassword(b *testing.B) {
	password := "TestPassword123!"
	hash, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = VerifyPassword(password, hash)
	}
}
