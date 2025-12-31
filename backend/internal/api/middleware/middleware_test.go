package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/clearclown/HaiLanGo/backend/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIsOriginAllowed_DefaultAllow(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	require.True(t, IsOriginAllowed("http://example.com"))
}

func TestIsOriginAllowed_WildcardAllow(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "*")
	require.True(t, IsOriginAllowed("http://example.com"))
}

func TestIsOriginAllowed_ListAllow(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000, https://example.com")

	require.True(t, IsOriginAllowed("http://localhost:3000"))
	require.True(t, IsOriginAllowed("https://example.com"))
	require.False(t, IsOriginAllowed("https://evil.example"))
}

func TestRequestID_GeneratedWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RequestID())
	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	got := w.Header().Get(requestIDHeader)
	require.NotEmpty(t, got)
	require.True(t, isSafeRequestID(got))
}

func TestRequestID_PassthroughWhenSafe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RequestID())
	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(requestIDHeader, "abc-123_DEF.456:ghi")
	r.ServeHTTP(w, req)

	require.Equal(t, "abc-123_DEF.456:ghi", w.Header().Get(requestIDHeader))
}

func TestRequestID_RegeneratedWhenUnsafe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RequestID())
	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(requestIDHeader, "bad id\ninjected")
	r.ServeHTTP(w, req)

	got := w.Header().Get(requestIDHeader)
	require.NotEmpty(t, got)
	require.NotEqual(t, "bad id\ninjected", got)
	require.True(t, isSafeRequestID(got))
}

func TestAuthRequired_TokenQueryRejectedWhenNotWebSocketHandshake(t *testing.T) {
	gin.SetMode(gin.TestMode)

	require.NoError(t, jwt.GenerateRSAKeys())
	userID := uuid.NewString()
	token, err := jwt.GenerateToken(userID, "test@example.com")
	require.NoError(t, err)

	r := gin.New()
	r.Use(RequestID())
	r.GET("/protected", AuthRequired(), func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected?token="+url.QueryEscape(token), nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthRequired_TokenQueryAllowedForWebSocketHandshake(t *testing.T) {
	gin.SetMode(gin.TestMode)

	require.NoError(t, jwt.GenerateRSAKeys())
	userID := uuid.NewString()
	token, err := jwt.GenerateToken(userID, "test@example.com")
	require.NoError(t, err)

	r := gin.New()
	r.Use(RequestID())
	r.GET("/ws", AuthRequired(), func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ws?token="+url.QueryEscape(token), nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}
