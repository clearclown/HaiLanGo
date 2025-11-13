package pattern

import (
	"context"
	"testing"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

func TestExtractor_ExtractPatterns(t *testing.T) {
	tests := []struct {
		name           string
		pages          []PageText
		minFrequency   int
		expectedCount  int
		expectedTypes  map[models.PatternType]int
	}{
		{
			name: "extract greeting patterns",
			pages: []PageText{
				{PageNumber: 1, Text: "Hello! How are you?", Translation: "こんにちは！元気ですか？"},
				{PageNumber: 2, Text: "Hello! Nice to meet you.", Translation: "こんにちは！はじめまして。"},
				{PageNumber: 3, Text: "Hello! Good morning.", Translation: "こんにちは！おはようございます。"},
			},
			minFrequency:  2,
			expectedCount: 3, // Multiple patterns extracted with minFrequency >= 2
			expectedTypes: map[models.PatternType]int{
				models.PatternTypeGreeting: 2, // "Hello" and "Good morning"
				models.PatternTypeQuestion: 1, // "How are"
			},
		},
		{
			name: "extract question patterns",
			pages: []PageText{
				{PageNumber: 1, Text: "How are you?", Translation: "元気ですか？"},
				{PageNumber: 2, Text: "How is the weather?", Translation: "天気はどうですか？"},
				{PageNumber: 3, Text: "How old are you?", Translation: "何歳ですか？"},
			},
			minFrequency:  2,
			expectedCount: 2, // Multiple "How" patterns extracted
			expectedTypes: map[models.PatternType]int{
				models.PatternTypeQuestion: 2,
			},
		},
		{
			name: "no patterns with high frequency threshold",
			pages: []PageText{
				{PageNumber: 1, Text: "Hello!", Translation: "こんにちは！"},
				{PageNumber: 2, Text: "Goodbye!", Translation: "さようなら！"},
				{PageNumber: 3, Text: "Thank you!", Translation: "ありがとう！"},
			},
			minFrequency:  3,
			expectedCount: 0,
			expectedTypes: map[models.PatternType]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := NewExtractor()
			ctx := context.Background()
			bookID := uuid.New()

			patterns, err := extractor.ExtractPatterns(ctx, bookID, tt.pages, tt.minFrequency)
			if err != nil {
				t.Fatalf("ExtractPatterns() error = %v", err)
			}

			if len(patterns) != tt.expectedCount {
				t.Errorf("ExtractPatterns() got %d patterns, want %d", len(patterns), tt.expectedCount)
			}

			// Check pattern types
			typeCounts := make(map[models.PatternType]int)
			for _, p := range patterns {
				typeCounts[p.Type]++
			}

			for patternType, expectedCount := range tt.expectedTypes {
				if count, ok := typeCounts[patternType]; !ok || count != expectedCount {
					t.Errorf("Pattern type %s: got %d, want %d", patternType, count, expectedCount)
				}
			}
		})
	}
}

func TestExtractor_CalculateFrequency(t *testing.T) {
	tests := []struct {
		name              string
		pattern           string
		texts             []string
		expectedFrequency int
	}{
		{
			name:    "exact match",
			pattern: "Hello",
			texts: []string{
				"Hello world",
				"Hello there",
				"Hello everyone",
			},
			expectedFrequency: 3,
		},
		{
			name:    "case insensitive",
			pattern: "hello",
			texts: []string{
				"Hello world",
				"HELLO there",
				"hello everyone",
			},
			expectedFrequency: 3,
		},
		{
			name:    "no match",
			pattern: "Goodbye",
			texts: []string{
				"Hello world",
				"Hi there",
				"Hey everyone",
			},
			expectedFrequency: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := NewExtractor()
			frequency := extractor.CalculateFrequency(tt.pattern, tt.texts)

			if frequency != tt.expectedFrequency {
				t.Errorf("CalculateFrequency() = %d, want %d", frequency, tt.expectedFrequency)
			}
		})
	}
}

func TestExtractor_GenerateExamples(t *testing.T) {
	tests := []struct {
		name           string
		pattern        models.Pattern
		pages          []PageText
		maxExamples    int
		expectedCount  int
	}{
		{
			name: "generate examples",
			pattern: models.Pattern{
				ID:      uuid.New(),
				Pattern: "Hello",
			},
			pages: []PageText{
				{PageNumber: 1, Text: "Hello world", Translation: "こんにちは世界"},
				{PageNumber: 2, Text: "Hello there", Translation: "こんにちはそこ"},
				{PageNumber: 3, Text: "Hello everyone", Translation: "こんにちはみんな"},
			},
			maxExamples:   2,
			expectedCount: 2,
		},
		{
			name: "no matching examples",
			pattern: models.Pattern{
				ID:      uuid.New(),
				Pattern: "Goodbye",
			},
			pages: []PageText{
				{PageNumber: 1, Text: "Hello world", Translation: "こんにちは世界"},
				{PageNumber: 2, Text: "Hi there", Translation: "やあそこ"},
			},
			maxExamples:   5,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := NewExtractor()
			ctx := context.Background()

			examples, err := extractor.GenerateExamples(ctx, tt.pattern, tt.pages, tt.maxExamples)
			if err != nil {
				t.Fatalf("GenerateExamples() error = %v", err)
			}

			if len(examples) != tt.expectedCount {
				t.Errorf("GenerateExamples() got %d examples, want %d", len(examples), tt.expectedCount)
			}

			// Verify all examples contain the pattern
			for _, example := range examples {
				if example.PatternID != tt.pattern.ID {
					t.Errorf("Example PatternID = %v, want %v", example.PatternID, tt.pattern.ID)
				}
			}
		})
	}
}

func TestExtractor_Performance(t *testing.T) {
	// Test performance requirement: 10 seconds for 100 pages
	extractor := NewExtractor()
	ctx := context.Background()
	bookID := uuid.New()

	// Generate 100 pages of test data
	pages := make([]PageText, 100)
	for i := 0; i < 100; i++ {
		pages[i] = PageText{
			PageNumber:  i + 1,
			Text:        "Hello! How are you? I'm fine, thank you.",
			Translation: "こんにちは！元気ですか？私は元気です、ありがとう。",
		}
	}

	start := time.Now()
	patterns, err := extractor.ExtractPatterns(ctx, bookID, pages, 2)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("ExtractPatterns() error = %v", err)
	}

	if duration > 10*time.Second {
		t.Errorf("ExtractPatterns() took %v, want < 10s", duration)
	}

	t.Logf("Extracted %d patterns in %v", len(patterns), duration)
}
