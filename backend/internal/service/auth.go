package service

import (
	"context"
	"errors"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/pkg/jwt"
	"github.com/clearclown/HaiLanGo/backend/pkg/password"
	"github.com/google/uuid"
)

var (
	// ErrInvalidCredentials は認証情報が無効なエラー
	ErrInvalidCredentials = errors.New("メールアドレスまたはパスワードが正しくありません")
	// ErrWeakPassword はパスワードが弱いエラー
	ErrWeakPassword = errors.New("パスワードが弱すぎます")
	// ErrUserAlreadyExists はユーザーが既に存在するエラー
	ErrUserAlreadyExists = errors.New("このメールアドレスは既に使用されています")
)

// AuthService は認証サービスのインターフェース
type AuthService interface {
	// Register はユーザーを登録する
	Register(ctx context.Context, req *models.CreateUserRequest) (*models.AuthResponse, error)
	// Login はメールアドレスとパスワードでログインする
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	// RefreshToken はリフレッシュトークンで新しいアクセストークンを取得する
	RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error)
	// Logout はログアウトしてリフレッシュトークンを無効化する
	Logout(ctx context.Context, refreshToken string) error
	// GetUserByID はIDでユーザーを取得する
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

// authService はAuthServiceの実装
type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService はAuthServiceの新しいインスタンスを作成
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Register はユーザーを登録する
func (s *authService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.AuthResponse, error) {
	// パスワード強度の検証
	if err := password.ValidatePasswordStrength(req.Password); err != nil {
		return nil, ErrWeakPassword
	}

	// メールアドレスの重複チェック
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// パスワードをハッシュ化
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// ユーザーの作成
	user := &models.User{
		Email:        req.Email,
		PasswordHash: &hashedPassword,
		DisplayName:  req.DisplayName,
		EmailVerified: false,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// JWTトークンの生成
	accessToken, err := jwt.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンの生成
	refreshTokenStr, expiresAt, err := jwt.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンをDBに保存
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: expiresAt,
	}

	if err := s.userRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	// レスポンスの作成
	return &models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    int64(jwt.AccessTokenExpiry.Seconds()),
	}, nil
}

// Login はメールアドレスとパスワードでログインする
func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// ユーザーの取得
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// パスワードの検証
	if user.PasswordHash == nil {
		// OAuth認証のみのユーザー
		return nil, ErrInvalidCredentials
	}

	if !password.VerifyPassword(req.Password, *user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// JWTトークンの生成
	accessToken, err := jwt.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンの生成
	refreshTokenStr, expiresAt, err := jwt.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンをDBに保存
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: expiresAt,
	}

	if err := s.userRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	// レスポンスの作成
	return &models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    int64(jwt.AccessTokenExpiry.Seconds()),
	}, nil
}

// RefreshToken はリフレッシュトークンで新しいアクセストークンを取得する
func (s *authService) RefreshToken(ctx context.Context, refreshTokenStr string) (*models.AuthResponse, error) {
	// リフレッシュトークンの取得
	refreshToken, err := s.userRepo.GetRefreshToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, errors.New("無効なリフレッシュトークンです")
	}

	// トークンの有効期限チェック
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("リフレッシュトークンの有効期限が切れています")
	}

	// ユーザーの取得
	user, err := s.userRepo.GetUserByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	// 新しいアクセストークンの生成
	accessToken, err := jwt.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, err
	}

	// レスポンスの作成
	return &models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    int64(jwt.AccessTokenExpiry.Seconds()),
	}, nil
}

// Logout はログアウトしてリフレッシュトークンを無効化する
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.userRepo.RevokeRefreshToken(ctx, refreshToken)
}

// GetUserByID はIDでユーザーを取得する
func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}
