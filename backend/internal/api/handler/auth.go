package handler

import (
	"net/http"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler は認証APIのハンドラー
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler はAuthHandlerの新しいインスタンスを作成
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register はユーザー登録エンドポイント
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest

	// リクエストのバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエストの形式が正しくありません",
		})
		return
	}

	// ユーザー登録
	authResp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		// エラーハンドリング
		switch err {
		case service.ErrWeakPassword:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		case service.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ユーザー登録に失敗しました",
			})
		}
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusCreated, authResp)
}

// Login はログインエンドポイント
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// リクエストのバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエストの形式が正しくありません",
		})
		return
	}

	// ログイン
	authResp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		// エラーハンドリング
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ログインに失敗しました",
			})
		}
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, authResp)
}

// RefreshToken はトークンリフレッシュエンドポイント
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest

	// リクエストのバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエストの形式が正しくありません",
		})
		return
	}

	// トークンリフレッシュ
	authResp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, authResp)
}

// Logout はログアウトエンドポイント
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req models.RefreshTokenRequest

	// リクエストのバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエストの形式が正しくありません",
		})
		return
	}

	// ログアウト
	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ログアウトに失敗しました",
		})
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, gin.H{
		"message": "ログアウトしました",
	})
}
