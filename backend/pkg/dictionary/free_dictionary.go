package dictionary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// FreeDictionaryClient is the client for Free Dictionary API
type FreeDictionaryClient struct {
	baseURL string
	client  *http.Client
}

// NewFreeDictionaryClient creates a new Free Dictionary API client
func NewFreeDictionaryClient() *FreeDictionaryClient {
	return &FreeDictionaryClient{
		baseURL: "https://api.dictionaryapi.dev/api/v2/entries",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// LookupWord looks up a word in Free Dictionary
func (c *FreeDictionaryClient) LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	// Check if we should use mock
	useMock := os.Getenv("USE_MOCK_APIS") == "true" || os.Getenv("TEST_USE_MOCKS") == "true"
	if useMock {
		return c.mockLookup(word, language)
	}

	// Build URL
	url := fmt.Sprintf("%s/%s/%s", c.baseURL, language, word)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
	var freeDictResp []FreeDictionaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&freeDictResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(freeDictResp) == 0 {
		return nil, ErrWordNotFound
	}

	// Convert to WordEntry
	entry := c.convertToWordEntry(&freeDictResp[0], word, language)
	return entry, nil
}

// GetName returns the name of the dictionary API
func (c *FreeDictionaryClient) GetName() string {
	return "Free Dictionary"
}

// mockLookup provides mock data for testing
func (c *FreeDictionaryClient) mockLookup(word string, language string) (*models.WordEntry, error) {
	if word == "xyzabc123notaword" {
		return nil, ErrWordNotFound
	}

	return &models.WordEntry{
		Word:      word,
		Language:  language,
		SourceAPI: "free_dictionary",
		FetchedAt: time.Now(),
		Phonetics: []models.WordPhonetic{
			{Text: "/həˈləʊ/"},
		},
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "interjection",
				Definitions: []models.WordDefinition{
					{
						Definition: "used to greet someone",
						Examples:   []string{"Hello there!"},
					},
				},
			},
		},
	}, nil
}

// convertToWordEntry converts Free Dictionary API response to WordEntry
func (c *FreeDictionaryClient) convertToWordEntry(resp *FreeDictionaryResponse, word string, language string) *models.WordEntry {
	entry := &models.WordEntry{
		Word:      word,
		Language:  language,
		SourceAPI: "free_dictionary",
		FetchedAt: time.Now(),
		Phonetics: []models.WordPhonetic{},
		Meanings:  []models.WordMeaning{},
	}

	// Extract phonetics
	for _, phonetic := range resp.Phonetics {
		p := models.WordPhonetic{
			Text: phonetic.Text,
		}
		if phonetic.Audio != "" {
			p.AudioURL = phonetic.Audio
		}
		entry.Phonetics = append(entry.Phonetics, p)
	}

	// Extract meanings
	for _, meaning := range resp.Meanings {
		m := models.WordMeaning{
			PartOfSpeech: meaning.PartOfSpeech,
			Definitions:  []models.WordDefinition{},
		}

		for _, def := range meaning.Definitions {
			d := models.WordDefinition{
				Definition: def.Definition,
				Examples:   []string{},
				Synonyms:   def.Synonyms,
				Antonyms:   def.Antonyms,
			}

			if def.Example != "" {
				d.Examples = append(d.Examples, def.Example)
			}

			m.Definitions = append(m.Definitions, d)
		}

		entry.Meanings = append(entry.Meanings, m)
	}

	return entry
}

// FreeDictionaryResponse represents the response from Free Dictionary API
type FreeDictionaryResponse struct {
	Word      string                   `json:"word"`
	Phonetics []FreeDictionaryPhonetic `json:"phonetics"`
	Meanings  []FreeDictionaryMeaning  `json:"meanings"`
}

type FreeDictionaryPhonetic struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

type FreeDictionaryMeaning struct {
	PartOfSpeech string                     `json:"partOfSpeech"`
	Definitions  []FreeDictionaryDefinition `json:"definitions"`
}

type FreeDictionaryDefinition struct {
	Definition string   `json:"definition"`
	Example    string   `json:"example"`
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
}
