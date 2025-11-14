package handler

import (
	"net/http"

	"github.com/clearclown/HaiLanGo/backend/internal/service/dictionary"
	pkgDict "github.com/clearclown/HaiLanGo/backend/pkg/dictionary"
	"github.com/gin-gonic/gin"
)

// DictionaryHandler handles dictionary-related HTTP requests
type DictionaryHandler struct {
	service *dictionary.Service
}

// NewDictionaryHandler creates a new dictionary handler
func NewDictionaryHandler(service *dictionary.Service) *DictionaryHandler {
	return &DictionaryHandler{
		service: service,
	}
}

// LookupWord handles GET /api/v1/dictionary/words/:word
func (h *DictionaryHandler) LookupWord(c *gin.Context) {
	word := c.Param("word")

	// Get language from query parameter (default: en)
	language := c.DefaultQuery("language", "en")

	// Lookup word
	entry, err := h.service.LookupWord(c.Request.Context(), word, language)
	if err != nil {
		if err == pkgDict.ErrWordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// LookupWordDetails handles GET /api/v1/dictionary/words/:word/details
func (h *DictionaryHandler) LookupWordDetails(c *gin.Context) {
	word := c.Param("word")

	// Get language from query parameter (default: en)
	language := c.DefaultQuery("language", "en")

	// Lookup word details
	entry, err := h.service.LookupWordDetails(c.Request.Context(), word, language)
	if err != nil {
		if err == pkgDict.ErrWordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// RegisterRoutes registers dictionary routes
func (h *DictionaryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	dictionary := rg.Group("/dictionary")
	{
		dictionary.GET("/words/:word", h.LookupWord)
		dictionary.GET("/words/:word/details", h.LookupWordDetails)
	}
}
