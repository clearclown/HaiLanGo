package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateUser はユーザー作成のテスト
func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	email := "test@example.com"
	passwordHash := "hashed_password"
	displayName := "Test User"

	t.Run("正常なユーザー作成", func(t *testing.T) {
		user := &models.User{
			Email:        email,
			PasswordHash: &passwordHash,
			DisplayName:  displayName,
		}

		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "email_verified"}).
			AddRow(userID, time.Now(), time.Now(), false)

		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(email, passwordHash, displayName, nil, nil, nil).
			WillReturnRows(rows)

		err := repo.CreateUser(ctx, user)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGetUserByEmail はメールアドレスでユーザーを取得するテスト
func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	email := "test@example.com"
	passwordHash := "hashed_password"
	displayName := "Test User"
	createdAt := time.Now()
	updatedAt := time.Now()

	t.Run("正常なユーザー取得", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "display_name", "profile_image_url",
			"oauth_provider", "oauth_provider_id", "email_verified",
			"email_verification_token", "email_verification_expires_at",
			"password_reset_token", "password_reset_expires_at",
			"created_at", "updated_at", "deleted_at",
		}).AddRow(
			userID, email, passwordHash, displayName, nil,
			nil, nil, false,
			nil, nil,
			nil, nil,
			createdAt, updatedAt, nil,
		)

		mock.ExpectQuery(`SELECT .+ FROM users WHERE email = \$1`).
			WithArgs(email).
			WillReturnRows(rows)

		user, err := repo.GetUserByEmail(ctx, email)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, displayName, user.DisplayName)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ユーザーが見つからない", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM users WHERE email = \$1`).
			WithArgs("notfound@example.com").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByEmail(ctx, "notfound@example.com")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGetUserByID はIDでユーザーを取得するテスト
func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	email := "test@example.com"
	passwordHash := "hashed_password"
	displayName := "Test User"
	createdAt := time.Now()
	updatedAt := time.Now()

	t.Run("正常なユーザー取得", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "email", "password_hash", "display_name", "profile_image_url",
			"oauth_provider", "oauth_provider_id", "email_verified",
			"email_verification_token", "email_verification_expires_at",
			"password_reset_token", "password_reset_expires_at",
			"created_at", "updated_at", "deleted_at",
		}).AddRow(
			userID, email, passwordHash, displayName, nil,
			nil, nil, false,
			nil, nil,
			nil, nil,
			createdAt, updatedAt, nil,
		)

		mock.ExpectQuery(`SELECT .+ FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnRows(rows)

		user, err := repo.GetUserByID(ctx, userID)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, email, user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ユーザーが見つからない", func(t *testing.T) {
		notFoundID := uuid.New()
		mock.ExpectQuery(`SELECT .+ FROM users WHERE id = \$1`).
			WithArgs(notFoundID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByID(ctx, notFoundID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestUpdateUser はユーザー更新のテスト
func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	newDisplayName := "Updated User"

	t.Run("正常なユーザー更新", func(t *testing.T) {
		user := &models.User{
			ID:          userID,
			DisplayName: newDisplayName,
		}

		mock.ExpectExec(`UPDATE users SET`).
			WithArgs(newDisplayName, nil, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUser(ctx, user)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ユーザーが見つからない", func(t *testing.T) {
		user := &models.User{
			ID:          uuid.New(),
			DisplayName: newDisplayName,
		}

		mock.ExpectExec(`UPDATE users SET`).
			WithArgs(newDisplayName, nil, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateUser(ctx, user)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestDeleteUser はユーザー削除（ソフトデリート）のテスト
func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	userID := uuid.New()

	t.Run("正常なユーザー削除", func(t *testing.T) {
		mock.ExpectExec(`UPDATE users SET deleted_at`).
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteUser(ctx, userID)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ユーザーが見つからない", func(t *testing.T) {
		notFoundID := uuid.New()
		mock.ExpectExec(`UPDATE users SET deleted_at`).
			WithArgs(notFoundID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteUser(ctx, notFoundID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCreateRefreshToken はリフレッシュトークン作成のテスト
func TestCreateRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	tokenID := uuid.New()
	userID := uuid.New()
	token := "refresh_token_string"
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	t.Run("正常なリフレッシュトークン作成", func(t *testing.T) {
		rt := &models.RefreshToken{
			UserID:    userID,
			Token:     token,
			ExpiresAt: expiresAt,
		}

		rows := sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow(tokenID, time.Now())

		mock.ExpectQuery(`INSERT INTO refresh_tokens`).
			WithArgs(userID, token, expiresAt).
			WillReturnRows(rows)

		err := repo.CreateRefreshToken(ctx, rt)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGetRefreshToken はリフレッシュトークン取得のテスト
func TestGetRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	tokenID := uuid.New()
	userID := uuid.New()
	token := "refresh_token_string"
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	createdAt := time.Now()

	t.Run("正常なリフレッシュトークン取得", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "user_id", "token", "expires_at", "created_at", "revoked_at",
		}).AddRow(
			tokenID, userID, token, expiresAt, createdAt, nil,
		)

		mock.ExpectQuery(`SELECT .+ FROM refresh_tokens WHERE token = \$1`).
			WithArgs(token).
			WillReturnRows(rows)

		rt, err := repo.GetRefreshToken(ctx, token)
		require.NoError(t, err)
		assert.NotNil(t, rt)
		assert.Equal(t, token, rt.Token)
		assert.Equal(t, userID, rt.UserID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("トークンが見つからない", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM refresh_tokens WHERE token = \$1`).
			WithArgs("invalid_token").
			WillReturnError(sql.ErrNoRows)

		rt, err := repo.GetRefreshToken(ctx, "invalid_token")
		assert.Error(t, err)
		assert.Nil(t, rt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestRevokeRefreshToken はリフレッシュトークン無効化のテスト
func TestRevokeRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	token := "refresh_token_string"

	t.Run("正常なトークン無効化", func(t *testing.T) {
		mock.ExpectExec(`UPDATE refresh_tokens SET revoked_at`).
			WithArgs(token).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.RevokeRefreshToken(ctx, token)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("トークンが見つからない", func(t *testing.T) {
		mock.ExpectExec(`UPDATE refresh_tokens SET revoked_at`).
			WithArgs("invalid_token").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.RevokeRefreshToken(ctx, "invalid_token")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
