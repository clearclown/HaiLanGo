package dictionary

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Enable mock mode for all tests
	os.Setenv("TEST_USE_MOCKS", "true")
	code := m.Run()
	os.Exit(code)
}

func TestService_LookupWord(t *testing.T) {
	service, err := NewService()
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		entry, err := service.LookupWord(ctx, "hello", "en")
		require.NoError(t, err)
		require.NotNil(t, entry)

		assert.Equal(t, "hello", entry.Word)
		assert.NotEmpty(t, entry.Meanings)
	})

	t.Run("EmptyWord", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.LookupWord(ctx, "", "en")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "word cannot be empty")
	})

	t.Run("DefaultLanguage", func(t *testing.T) {
		ctx := context.Background()
		entry, err := service.LookupWord(ctx, "hello", "")
		require.NoError(t, err)
		assert.Equal(t, "en", entry.Language)
	})

	t.Run("WordNotFound", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.LookupWord(ctx, "xyzabc123notaword", "en")
		assert.Error(t, err)
	})
}

func TestService_LookupWordDetails(t *testing.T) {
	service, err := NewService()
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		entry, err := service.LookupWordDetails(ctx, "hello", "en")
		require.NoError(t, err)
		require.NotNil(t, entry)

		assert.Equal(t, "hello", entry.Word)
		assert.NotEmpty(t, entry.Meanings)
	})

	t.Run("EmptyWord", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.LookupWordDetails(ctx, "", "en")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "word cannot be empty")
	})

	t.Run("DefaultLanguage", func(t *testing.T) {
		ctx := context.Background()
		entry, err := service.LookupWordDetails(ctx, "hello", "")
		require.NoError(t, err)
		assert.Equal(t, "en", entry.Language)
	})
}
