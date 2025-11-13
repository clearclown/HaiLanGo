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

// OxfordClient is the client for Oxford Dictionary API
type OxfordClient struct {
	apiKey  string
	appID   string
	baseURL string
	client  *http.Client
}

// NewOxfordClient creates a new Oxford Dictionary API client
func NewOxfordClient(apiKey, appID string) *OxfordClient {
	return &OxfordClient{
		apiKey:  apiKey,
		appID:   appID,
		baseURL: "https://od-api.oxforddictionaries.com/api/v2",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// LookupWord looks up a word in Oxford Dictionary
func (c *OxfordClient) LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	// Check if we should use mock
	useMock := os.Getenv("USE_MOCK_APIS") == "true" || os.Getenv("TEST_USE_MOCKS") == "true"
	if useMock || c.apiKey == "" {
		return c.mockLookup(word, language)
	}

	// Build URL
	url := fmt.Sprintf("%s/entries/%s/%s", c.baseURL, language, word)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("app_id", c.appID)
	req.Header.Set("app_key", c.apiKey)

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
	var oxfordResp OxfordResponse
	if err := json.NewDecoder(resp.Body).Decode(&oxfordResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to WordEntry
	entry := c.convertToWordEntry(&oxfordResp, word, language)
	return entry, nil
}

// GetName returns the name of the dictionary API
func (c *OxfordClient) GetName() string {
	return "Oxford Dictionary"
}

// mockLookup provides mock data for testing
func (c *OxfordClient) mockLookup(word string, language string) (*models.WordEntry, error) {
	if word == "xyzabc123notaword" {
		return nil, ErrWordNotFound
	}

	return &models.WordEntry{
		Word:      word,
		Language:  language,
		SourceAPI: "oxford",
		FetchedAt: time.Now(),
		Phonetics: []models.WordPhonetic{
			{Text: "/həˈloʊ/"},
		},
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "noun",
				Definitions: []models.WordDefinition{
					{
						Definition: "used as a greeting",
						Examples:   []string{"Hello, how are you?"},
					},
				},
			},
		},
	}, nil
}

// convertToWordEntry converts Oxford API response to WordEntry
func (c *OxfordClient) convertToWordEntry(resp *OxfordResponse, word string, language string) *models.WordEntry {
	entry := &models.WordEntry{
		Word:      word,
		Language:  language,
		SourceAPI: "oxford",
		FetchedAt: time.Now(),
		Phonetics: []models.WordPhonetic{},
		Meanings:  []models.WordMeaning{},
	}

	if len(resp.Results) == 0 {
		return entry
	}

	result := resp.Results[0]

	// Extract phonetics
	if len(result.LexicalEntries) > 0 {
		lexEntry := result.LexicalEntries[0]
		if len(lexEntry.Pronunciations) > 0 {
			for _, pron := range lexEntry.Pronunciations {
				phonetic := models.WordPhonetic{
					Text: pron.PhoneticSpelling,
				}
				if pron.AudioFile != "" {
					phonetic.AudioURL = pron.AudioFile
				}
				entry.Phonetics = append(entry.Phonetics, phonetic)
			}
		}
	}

	// Extract meanings
	for _, lexEntry := range result.LexicalEntries {
		meaning := models.WordMeaning{
			PartOfSpeech: lexEntry.LexicalCategory.Text,
			Definitions:  []models.WordDefinition{},
		}

		for _, entry := range lexEntry.Entries {
			for _, sense := range entry.Senses {
				def := models.WordDefinition{
					Definition: "",
					Examples:   []string{},
					Synonyms:   []string{},
					Antonyms:   []string{},
				}

				if len(sense.Definitions) > 0 {
					def.Definition = sense.Definitions[0]
				}

				for _, example := range sense.Examples {
					def.Examples = append(def.Examples, example.Text)
				}

				for _, synonym := range sense.Synonyms {
					def.Synonyms = append(def.Synonyms, synonym.Text)
				}

				for _, antonym := range sense.Antonyms {
					def.Antonyms = append(def.Antonyms, antonym.Text)
				}

				meaning.Definitions = append(meaning.Definitions, def)
			}
		}

		entry.Meanings = append(entry.Meanings, meaning)
	}

	return entry
}

// OxfordResponse represents the response from Oxford Dictionary API
type OxfordResponse struct {
	Results []OxfordResult `json:"results"`
}

type OxfordResult struct {
	LexicalEntries []OxfordLexicalEntry `json:"lexicalEntries"`
}

type OxfordLexicalEntry struct {
	LexicalCategory OxfordCategory        `json:"lexicalCategory"`
	Pronunciations  []OxfordPronunciation `json:"pronunciations"`
	Entries         []OxfordEntry         `json:"entries"`
}

type OxfordCategory struct {
	Text string `json:"text"`
}

type OxfordPronunciation struct {
	PhoneticSpelling string `json:"phoneticSpelling"`
	AudioFile        string `json:"audioFile"`
}

type OxfordEntry struct {
	Senses []OxfordSense `json:"senses"`
}

type OxfordSense struct {
	Definitions []string        `json:"definitions"`
	Examples    []OxfordExample `json:"examples"`
	Synonyms    []OxfordSynonym `json:"synonyms"`
	Antonyms    []OxfordAntonym `json:"antonyms"`
}

type OxfordExample struct {
	Text string `json:"text"`
}

type OxfordSynonym struct {
	Text string `json:"text"`
}

type OxfordAntonym struct {
	Text string `json:"text"`
}
