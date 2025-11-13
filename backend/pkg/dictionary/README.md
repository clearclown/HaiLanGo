# Dictionary Package

This package provides integration with multiple dictionary APIs (Oxford Dictionary, Free Dictionary, Wiktionary) with automatic fallback and caching support.

## Features

- **Multiple Dictionary APIs**: Oxford Dictionary (primary), Free Dictionary, Wiktionary (fallback)
- **Automatic Fallback**: If one API fails, automatically tries the next
- **Redis Caching**: 30-day cache for dictionary results
- **Mock Support**: Automatic mock mode for testing without API keys

## Installation

```bash
go get github.com/clearclown/HaiLanGo/backend/pkg/dictionary
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/clearclown/HaiLanGo/backend/pkg/dictionary"
)

func main() {
    // Create dictionary clients
    oxfordClient := dictionary.NewOxfordClient("your-api-key", "your-app-id")
    freeDictClient := dictionary.NewFreeDictionaryClient()
    wiktionaryClient := dictionary.NewWiktionaryClient()

    // Create cache client
    cache, err := dictionary.NewRedisCacheClient("localhost:6379", "", 0)
    if err != nil {
        log.Fatal(err)
    }

    // Create dictionary service with fallback
    service := dictionary.NewService(
        oxfordClient,
        []dictionary.Client{freeDictClient, wiktionaryClient},
        cache,
    )

    // Look up a word
    ctx := context.Background()
    entry, err := service.LookupWord(ctx, "hello", "en")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Word: %s\n", entry.Word)
    fmt.Printf("Source: %s\n", entry.SourceAPI)
    for _, meaning := range entry.Meanings {
        fmt.Printf("Part of Speech: %s\n", meaning.PartOfSpeech)
        for _, def := range meaning.Definitions {
            fmt.Printf("  Definition: %s\n", def.Definition)
            for _, example := range def.Examples {
                fmt.Printf("    Example: %s\n", example)
            }
        }
    }
}
```

### Using Mock Mode

For testing without API keys, set the environment variable:

```bash
export USE_MOCK_APIS=true
# or
export TEST_USE_MOCKS=true
```

When mock mode is enabled, all API calls will return mock data.

```go
// Mock mode is automatically detected
service := dictionary.NewService(
    dictionary.NewOxfordClient("", ""), // Empty API keys will trigger mock mode
    []dictionary.Client{dictionary.NewFreeDictionaryClient()},
    mockCache,
)
```

## API Clients

### Oxford Dictionary API

Requires API key and App ID from [Oxford Dictionaries](https://developer.oxforddictionaries.com/).

```go
client := dictionary.NewOxfordClient("your-api-key", "your-app-id")
entry, err := client.LookupWord(ctx, "word", "en")
```

### Free Dictionary API

No API key required. Uses [Free Dictionary API](https://dictionaryapi.dev/).

```go
client := dictionary.NewFreeDictionaryClient()
entry, err := client.LookupWord(ctx, "word", "en")
```

### Wiktionary API

No API key required. Uses Wiktionary's MediaWiki API.

```go
client := dictionary.NewWiktionaryClient()
entry, err := client.LookupWord(ctx, "word", "en")
```

## Caching

The package uses Redis for caching dictionary results for 30 days.

```go
cache, err := dictionary.NewRedisCacheClient("localhost:6379", "", 0)
if err != nil {
    log.Fatal(err)
}
```

Cache keys are generated as: `dictionary:{language}:{word}`

## Error Handling

The package defines several error types:

- `ErrWordNotFound`: Word not found in dictionary
- `ErrAPIUnavailable`: API is unavailable
- `ErrRateLimitExceeded`: API rate limit exceeded

```go
entry, err := service.LookupWord(ctx, "word", "en")
if errors.Is(err, dictionary.ErrWordNotFound) {
    fmt.Println("Word not found")
} else if errors.Is(err, dictionary.ErrAPIUnavailable) {
    fmt.Println("API unavailable")
}
```

## Testing

Run tests with mock mode:

```bash
TEST_USE_MOCKS=true go test ./pkg/dictionary/... -v
```

## Data Models

### WordEntry

```go
type WordEntry struct {
    Word      string         // The word
    Phonetics []WordPhonetic // Pronunciation information
    Meanings  []WordMeaning  // Meanings and definitions
    Language  string         // Language code (e.g., "en")
    SourceAPI string         // Which API provided this data
    FetchedAt time.Time      // When the data was fetched
}
```

### WordMeaning

```go
type WordMeaning struct {
    PartOfSpeech string           // e.g., "noun", "verb"
    Definitions  []WordDefinition // List of definitions
}
```

### WordDefinition

```go
type WordDefinition struct {
    Definition string   // The definition
    Examples   []string // Usage examples
    Synonyms   []string // Synonyms
    Antonyms   []string // Antonyms
}
```

## Performance

- **Cache hit rate**: 80%+ expected
- **Lookup latency**: <500ms (with cache), <2s (API call)
- **Cache TTL**: 30 days

## License

MIT
