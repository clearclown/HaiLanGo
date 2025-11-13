package pattern

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// PageText represents the text content of a page
type PageText struct {
	PageNumber  int
	Text        string
	Translation string
}

// Extractor handles pattern extraction from book pages
type Extractor struct {
	classifier *Classifier
}

// NewExtractor creates a new pattern extractor
func NewExtractor() *Extractor {
	return &Extractor{
		classifier: NewClassifier(),
	}
}

// ExtractPatterns extracts conversation patterns from pages
func (e *Extractor) ExtractPatterns(ctx context.Context, bookID uuid.UUID, pages []PageText, minFrequency int) ([]models.Pattern, error) {
	// Extract all potential patterns
	patternCandidates := make(map[string]*patternCandidate)

	for _, page := range pages {
		// Split text into sentences
		sentences := e.splitIntoSentences(page.Text)

		for _, sentence := range sentences {
			// Extract phrases from sentence
			phrases := e.extractPhrases(sentence)

			for _, phrase := range phrases {
				normalized := e.normalizePattern(phrase)
				if normalized == "" {
					continue
				}

				if candidate, exists := patternCandidates[normalized]; exists {
					candidate.frequency++
					candidate.examples = append(candidate.examples, page)
				} else {
					patternCandidates[normalized] = &patternCandidate{
						pattern:   normalized,
						frequency: 1,
						examples:  []PageText{page},
					}
				}
			}
		}
	}

	// Filter patterns by minimum frequency and create Pattern objects
	var patterns []models.Pattern
	now := time.Now()

	for _, candidate := range patternCandidates {
		if candidate.frequency >= minFrequency {
			pattern := models.Pattern{
				ID:          uuid.New(),
				BookID:      bookID,
				Type:        e.classifier.ClassifyPattern(candidate.pattern),
				Pattern:     candidate.pattern,
				Translation: e.extractTranslation(candidate.pattern, candidate.examples),
				Frequency:   candidate.frequency,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			patterns = append(patterns, pattern)
		}
	}

	return patterns, nil
}

// CalculateFrequency calculates how many times a pattern appears in texts
func (e *Extractor) CalculateFrequency(pattern string, texts []string) int {
	frequency := 0
	normalized := strings.ToLower(pattern)

	for _, text := range texts {
		textLower := strings.ToLower(text)
		count := strings.Count(textLower, normalized)
		frequency += count
	}

	return frequency
}

// GenerateExamples generates usage examples for a pattern
func (e *Extractor) GenerateExamples(ctx context.Context, pattern models.Pattern, pages []PageText, maxExamples int) ([]models.PatternExample, error) {
	var examples []models.PatternExample
	now := time.Now()

	patternLower := strings.ToLower(pattern.Pattern)

	for _, page := range pages {
		if len(examples) >= maxExamples {
			break
		}

		textLower := strings.ToLower(page.Text)
		if strings.Contains(textLower, patternLower) {
			// Find the sentence containing the pattern
			sentences := e.splitIntoSentences(page.Text)
			for _, sentence := range sentences {
				if strings.Contains(strings.ToLower(sentence), patternLower) {
					example := models.PatternExample{
						ID:             uuid.New(),
						PatternID:      pattern.ID,
						PageNumber:     page.PageNumber,
						OriginalText:   sentence,
						TranslatedText: page.Translation,
						Context:        page.Text,
						CreatedAt:      now,
					}
					examples = append(examples, example)
					break
				}
			}
		}
	}

	return examples, nil
}

// Helper functions

type patternCandidate struct {
	pattern   string
	frequency int
	examples  []PageText
}

func (e *Extractor) splitIntoSentences(text string) []string {
	// Split by common sentence delimiters
	re := regexp.MustCompile(`[.!?]+`)
	sentences := re.Split(text, -1)

	var result []string
	for _, sentence := range sentences {
		trimmed := strings.TrimSpace(sentence)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

func (e *Extractor) extractPhrases(sentence string) []string {
	var phrases []string

	// Extract whole sentence
	trimmed := strings.TrimSpace(sentence)
	if trimmed != "" {
		phrases = append(phrases, trimmed)
	}

	// Extract phrases by splitting at commas
	parts := strings.Split(sentence, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" && trimmed != sentence {
			phrases = append(phrases, trimmed)
		}
	}

	// Extract phrases starting with common words
	words := strings.Fields(sentence)
	if len(words) >= 2 {
		// Extract 2-word phrases
		for i := 0; i < len(words)-1; i++ {
			phrase := words[i] + " " + words[i+1]
			phrases = append(phrases, phrase)
		}

		// Extract 3-word phrases
		for i := 0; i < len(words)-2; i++ {
			phrase := words[i] + " " + words[i+1] + " " + words[i+2]
			phrases = append(phrases, phrase)
		}
	}

	return phrases
}

func (e *Extractor) normalizePattern(pattern string) string {
	// Trim whitespace
	normalized := strings.TrimSpace(pattern)

	// Remove multiple spaces
	re := regexp.MustCompile(`\s+`)
	normalized = re.ReplaceAllString(normalized, " ")

	// Skip very short patterns
	if len(normalized) < 3 {
		return ""
	}

	return normalized
}

func (e *Extractor) extractTranslation(pattern string, examples []PageText) string {
	// Try to find a translation from the examples
	for _, example := range examples {
		if example.Translation != "" {
			return example.Translation
		}
	}
	return ""
}
