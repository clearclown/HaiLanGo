package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// DictionaryRepositoryInterface は辞書リポジトリのインターフェース
type DictionaryRepositoryInterface interface {
	// LookupWord は単語を検索
	LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error)

	// BatchLookup は複数の単語を検索
	BatchLookup(ctx context.Context, words []string, language string) ([]*models.WordEntry, error)

	// GetSupportedLanguages はサポートされている言語を取得
	GetSupportedLanguages(ctx context.Context) ([]string, error)

	// CacheWord は単語をキャッシュ
	CacheWord(ctx context.Context, entry *models.WordEntry) error
}

// InMemoryDictionaryRepository はインメモリ辞書リポジトリ
type InMemoryDictionaryRepository struct {
	mu    sync.RWMutex
	cache map[string]*models.WordEntry // word:language -> WordEntry
	languages []string
}

// NewInMemoryDictionaryRepository はインメモリ辞書リポジトリを作成
func NewInMemoryDictionaryRepository() *InMemoryDictionaryRepository {
	repo := &InMemoryDictionaryRepository{
		cache: make(map[string]*models.WordEntry),
		languages: []string{
			"en",    // English
			"ja",    // Japanese
			"zh",    // Chinese
			"ru",    // Russian
			"es",    // Spanish
			"fr",    // French
			"de",    // German
			"it",    // Italian
			"pt",    // Portuguese
			"ar",    // Arabic
			"he",    // Hebrew
			"fa",    // Persian
		},
	}

	// サンプルデータを初期化
	repo.initSampleData()

	return repo
}

func (r *InMemoryDictionaryRepository) initSampleData() {
	// 英語のサンプルデータ
	r.cache["hello:en"] = &models.WordEntry{
		Word: "hello",
		Phonetics: []models.WordPhonetic{
			{
				Text:     "/həˈloʊ/",
				AudioURL: "https://example.com/audio/hello.mp3",
			},
		},
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "interjection",
				Definitions: []models.WordDefinition{
					{
						Definition: "used as a greeting or to begin a phone conversation",
						Examples:   []string{"hello there, Katie!"},
						Synonyms:   []string{"hi", "hey", "greetings"},
					},
				},
			},
			{
				PartOfSpeech: "noun",
				Definitions: []models.WordDefinition{
					{
						Definition: "an utterance of 'hello'; a greeting",
						Examples:   []string{"she was getting polite nods and hellos from people"},
					},
				},
			},
		},
		Language:  "en",
		SourceAPI: "mock",
		FetchedAt: time.Now(),
	}

	r.cache["book:en"] = &models.WordEntry{
		Word: "book",
		Phonetics: []models.WordPhonetic{
			{
				Text:     "/bʊk/",
				AudioURL: "https://example.com/audio/book.mp3",
			},
		},
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "noun",
				Definitions: []models.WordDefinition{
					{
						Definition: "a written or printed work consisting of pages glued or sewn together along one side and bound in covers",
						Examples:   []string{"a book of selected poems"},
						Synonyms:   []string{"volume", "tome", "publication"},
					},
					{
						Definition: "a set of records or accounts",
						Examples:   []string{"the book of the club"},
					},
				},
			},
			{
				PartOfSpeech: "verb",
				Definitions: []models.WordDefinition{
					{
						Definition: "reserve (accommodation, a place, etc.); buy (a ticket) in advance",
						Examples:   []string{"I have booked a table at the restaurant"},
						Synonyms:   []string{"reserve", "make a reservation"},
					},
				},
			},
		},
		Language:  "en",
		SourceAPI: "mock",
		FetchedAt: time.Now(),
	}

	// ロシア語のサンプルデータ
	r.cache["здравствуйте:ru"] = &models.WordEntry{
		Word: "здравствуйте",
		Phonetics: []models.WordPhonetic{
			{
				Text:     "/zdrɐˈstvʊjtʲɪ/",
				AudioURL: "https://example.com/audio/zdravstvuyte.mp3",
			},
		},
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "interjection",
				Definitions: []models.WordDefinition{
					{
						Definition: "formal greeting (hello)",
						Examples:   []string{"Здравствуйте, как дела?"},
						Synonyms:   []string{"привет", "добрый день"},
					},
				},
			},
		},
		Language:  "ru",
		SourceAPI: "mock",
		FetchedAt: time.Now(),
	}

	// 日本語のサンプルデータ
	r.cache["こんにちは:ja"] = &models.WordEntry{
		Word: "こんにちは",
		Phonetics: []models.WordPhonetic{
			{
				Text:     "konnichiwa",
				AudioURL: "https://example.com/audio/konnichiwa.mp3",
			},
		},
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "interjection",
				Definitions: []models.WordDefinition{
					{
						Definition: "greeting used during daytime (hello)",
						Examples:   []string{"こんにちは、元気ですか？"},
						Synonyms:   []string{"おはよう", "こんばんは"},
					},
				},
			},
		},
		Language:  "ja",
		SourceAPI: "mock",
		FetchedAt: time.Now(),
	}
}

func (r *InMemoryDictionaryRepository) LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", word, language)
	entry, exists := r.cache[key]
	if !exists {
		// キャッシュにない場合は、ダミーのエントリを返す
		return r.generateDummyEntry(word, language), nil
	}

	return entry, nil
}

func (r *InMemoryDictionaryRepository) BatchLookup(ctx context.Context, words []string, language string) ([]*models.WordEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries := make([]*models.WordEntry, 0, len(words))
	for _, word := range words {
		key := fmt.Sprintf("%s:%s", word, language)
		entry, exists := r.cache[key]
		if exists {
			entries = append(entries, entry)
		} else {
			// キャッシュにない場合は、ダミーのエントリを生成
			entries = append(entries, r.generateDummyEntry(word, language))
		}
	}

	return entries, nil
}

func (r *InMemoryDictionaryRepository) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.languages, nil
}

func (r *InMemoryDictionaryRepository) CacheWord(ctx context.Context, entry *models.WordEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%s:%s", entry.Word, entry.Language)
	r.cache[key] = entry

	return nil
}

func (r *InMemoryDictionaryRepository) generateDummyEntry(word string, language string) *models.WordEntry {
	return &models.WordEntry{
		Word: word,
		Phonetics: []models.WordPhonetic{
			{
				Text: fmt.Sprintf("/%s/", word),
			},
		},
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "unknown",
				Definitions: []models.WordDefinition{
					{
						Definition: fmt.Sprintf("Definition for '%s' (mock data)", word),
						Examples:   []string{fmt.Sprintf("Example sentence with %s.", word)},
					},
				},
			},
		},
		Language:  language,
		SourceAPI: "mock",
		FetchedAt: time.Now(),
	}
}

// PostgreSQL Implementation

type DictionaryRepositoryPostgres struct {
	db *sql.DB
}

func NewDictionaryRepositoryPostgres(db *sql.DB) DictionaryRepositoryInterface {
	return &DictionaryRepositoryPostgres{db: db}
}

func (r *DictionaryRepositoryPostgres) LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	var entry models.WordEntry
	var phoneticsJSON, meaningsJSON []byte
	var sourceAPI sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT word, language, phonetics, meanings, source_api, fetched_at, last_accessed_at
		FROM dictionary_cache
		WHERE word = $1 AND language = $2
	`, word, language).Scan(&entry.Word, &entry.Language, &phoneticsJSON, &meaningsJSON, &sourceAPI, &entry.FetchedAt, &entry.FetchedAt)

	if err == sql.ErrNoRows {
		// Not in cache, generate dummy entry
		return r.generateDummyEntry(word, language), nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal JSONB fields
	json.Unmarshal(phoneticsJSON, &entry.Phonetics)
	json.Unmarshal(meaningsJSON, &entry.Meanings)
	if sourceAPI.Valid {
		entry.SourceAPI = sourceAPI.String
	}

	// Update access tracking
	r.db.ExecContext(ctx, `
		UPDATE dictionary_cache SET last_accessed_at = NOW(), access_count = access_count + 1
		WHERE word = $1 AND language = $2
	`, word, language)

	return &entry, nil
}

func (r *DictionaryRepositoryPostgres) BatchLookup(ctx context.Context, words []string, language string) ([]*models.WordEntry, error) {
	entries := make([]*models.WordEntry, 0, len(words))

	for _, word := range words {
		entry, err := r.LookupWord(ctx, word, language)
		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *DictionaryRepositoryPostgres) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT code FROM dictionary_supported_languages
		WHERE enabled = true
		ORDER BY code
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	languages := []string{}
	for rows.Next() {
		var code string
		rows.Scan(&code)
		languages = append(languages, code)
	}

	return languages, nil
}

func (r *DictionaryRepositoryPostgres) CacheWord(ctx context.Context, entry *models.WordEntry) error {
	phoneticsJSON, _ := json.Marshal(entry.Phonetics)
	meaningsJSON, _ := json.Marshal(entry.Meanings)

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO dictionary_cache (word, language, phonetics, meanings, source_api, fetched_at, last_accessed_at, access_count)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), 1)
		ON CONFLICT (word, language) DO UPDATE SET
			phonetics = $3,
			meanings = $4,
			source_api = $5,
			fetched_at = NOW(),
			last_accessed_at = NOW(),
			access_count = dictionary_cache.access_count + 1
	`, entry.Word, entry.Language, phoneticsJSON, meaningsJSON, entry.SourceAPI)

	return err
}

func (r *DictionaryRepositoryPostgres) generateDummyEntry(word string, language string) *models.WordEntry {
	return &models.WordEntry{
		Word: word,
		Phonetics: []models.WordPhonetic{
			{
				Text: fmt.Sprintf("/%s/", word),
			},
		},
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "unknown",
				Definitions: []models.WordDefinition{
					{
						Definition: fmt.Sprintf("Definition for '%s' (mock data)", word),
						Examples:   []string{fmt.Sprintf("Example sentence with %s.", word)},
					},
				},
			},
		},
		Language:  language,
		SourceAPI: "mock",
		FetchedAt: time.Now(),
	}
}
