package handler

import (
	"net/http"
	"strings"

	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// DictionaryHandler は辞書APIのハンドラー
type DictionaryHandler struct {
	repo repository.DictionaryRepositoryInterface
}

// NewDictionaryHandler は辞書ハンドラーを作成
func NewDictionaryHandler(repo repository.DictionaryRepositoryInterface) *DictionaryHandler {
	return &DictionaryHandler{
		repo: repo,
	}
}

// RegisterRoutes は辞書APIのルートを登録
func (h *DictionaryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	dictionary := rg.Group("/dictionary")
	{
		dictionary.GET("/words/:word", h.LookupWord)
		dictionary.POST("/batch", h.BatchLookup)
		dictionary.GET("/languages", h.GetSupportedLanguages)
	}
}

// BatchLookupRequest は複数単語検索リクエスト
type BatchLookupRequest struct {
	Words    []string `json:"words" binding:"required"`
	Language string   `json:"language" binding:"required"`
}

// LookupWord は単語を検索
// GET /api/v1/dictionary/words/:word?language=en
func (h *DictionaryHandler) LookupWord(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	word := c.Param("word")
	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Word parameter is required"})
		return
	}

	// 言語パラメータを取得（デフォルト: en）
	language := c.DefaultQuery("language", "en")
	language = strings.ToLower(language)

	entry, err := h.repo.LookupWord(c.Request.Context(), word, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to lookup word"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// BatchLookup は複数の単語を検索
// POST /api/v1/dictionary/batch
func (h *DictionaryHandler) BatchLookup(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req BatchLookupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req.Words) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Words array cannot be empty"})
		return
	}

	if len(req.Words) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot lookup more than 50 words at once"})
		return
	}

	req.Language = strings.ToLower(req.Language)

	entries, err := h.repo.BatchLookup(c.Request.Context(), req.Words, req.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to lookup words"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"entries": entries,
		"count":   len(entries),
	})
}

// GetSupportedLanguages はサポートされている言語を取得
// GET /api/v1/dictionary/languages
func (h *DictionaryHandler) GetSupportedLanguages(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	languages, err := h.repo.GetSupportedLanguages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get supported languages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"languages": languages,
		"count":     len(languages),
	})
}
