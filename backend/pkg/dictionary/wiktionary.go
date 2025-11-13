package dictionary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// WiktionaryClient is the client for Wiktionary API
type WiktionaryClient struct {
	baseURL string
	client  *http.Client
}

// NewWiktionaryClient creates a new Wiktionary API client
func NewWiktionaryClient() *WiktionaryClient {
	return &WiktionaryClient{
		baseURL: "https://en.wiktionary.org/w/api.php",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// LookupWord looks up a word in Wiktionary
func (c *WiktionaryClient) LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	// Check if we should use mock
	useMock := os.Getenv("USE_MOCK_APIS") == "true" || os.Getenv("TEST_USE_MOCKS") == "true"
	if useMock {
		return c.mockLookup(word, language)
	}

	// Build URL with query parameters
	params := url.Values{}
	params.Set("action", "query")
	params.Set("format", "json")
	params.Set("titles", word)
	params.Set("prop", "extracts")
	params.Set("explaintext", "1")
	params.Set("exsectionformat", "plain")

	fullURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Make request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", ErrAPIUnavailable)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrWordNotFound
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimitExceeded
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %w", resp.StatusCode, ErrAPIUnavailable)
	}

	// Parse response
	var wiktResp WiktionaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&wiktResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if page exists
	if wiktResp.Query.Pages == nil {
		return nil, ErrWordNotFound
	}

	// Get the first page (should be only one)
	var page *WiktionaryPage
	for _, p := range wiktResp.Query.Pages {
		page = &p
		break
	}

	if page == nil || page.Missing {
		return nil, ErrWordNotFound
	}

	// Convert to WordEntry
	entry := c.convertToWordEntry(page, word, language)
	return entry, nil
}

// GetName returns the name of the dictionary API
func (c *WiktionaryClient) GetName() string {
	return "Wiktionary"
}

// mockLookup provides mock data for testing
func (c *WiktionaryClient) mockLookup(word string, language string) (*models.WordEntry, error) {
	if word == "xyzabc123notaword" {
		return nil, ErrWordNotFound
	}

	return &models.WordEntry{
		Word:      word,
		Language:  language,
		SourceAPI: "wiktionary",
		FetchedAt: time.Now(),
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "interjection",
				Definitions: []models.WordDefinition{
					{
						Definition: "A greeting",
						Synonyms:   []string{"hi", "hey"},
					},
				},
			},
		},
	}, nil
}

// convertToWordEntry converts Wiktionary API response to WordEntry
func (c *WiktionaryClient) convertToWordEntry(page *WiktionaryPage, word string, language string) *models.WordEntry {
	entry := &models.WordEntry{
		Word:      word,
		Language:  language,
		SourceAPI: "wiktionary",
		FetchedAt: time.Now(),
		Phonetics: []models.WordPhonetic{},
		Meanings:  []models.WordMeaning{},
	}

	// Parse the extract text to extract definitions
	// This is a simplified parser - a full implementation would be more complex
	extract := page.Extract
	if extract == "" {
		return entry
	}

	// Split by sections (very simplified)
	sections := strings.Split(extract, "\n\n")

	meaning := models.WordMeaning{
		PartOfSpeech: "general",
		Definitions:  []models.WordDefinition{},
	}

	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section == "" {
			continue
		}

		// Simple heuristic: if it looks like a definition, add it
		if len(section) > 0 && !strings.HasPrefix(section, "==") {
			def := models.WordDefinition{
				Definition: section,
				Examples:   []string{},
			}
			meaning.Definitions = append(meaning.Definitions, def)
		}
	}

	if len(meaning.Definitions) > 0 {
		entry.Meanings = append(entry.Meanings, meaning)
	}

	return entry
}

// WiktionaryResponse represents the response from Wiktionary API
type WiktionaryResponse struct {
	Query WiktionaryQuery `json:"query"`
}

type WiktionaryQuery struct {
	Pages map[string]WiktionaryPage `json:"pages"`
}

type WiktionaryPage struct {
	PageID  int    `json:"pageid"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
	Missing bool   `json:"missing"`
}
