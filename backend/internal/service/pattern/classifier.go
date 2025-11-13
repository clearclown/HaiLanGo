package pattern

import (
	"strings"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// Classifier classifies conversation patterns into types
type Classifier struct {
	greetingKeywords     []string
	questionKeywords     []string
	requestKeywords      []string
	confirmationKeywords []string
	responseKeywords     []string
}

// NewClassifier creates a new pattern classifier
func NewClassifier() *Classifier {
	return &Classifier{
		greetingKeywords: []string{
			"hello", "hi", "hey", "good morning", "good afternoon",
			"good evening", "good night", "greetings", "howdy",
		},
		questionKeywords: []string{
			"how", "what", "where", "when", "why", "who", "which",
			"is", "are", "do", "does", "did", "can", "could",
			"would", "will", "shall", "may", "might",
		},
		requestKeywords: []string{
			"please", "could you", "would you", "can you", "will you",
			"i would like", "i'd like", "may i", "might i",
		},
		confirmationKeywords: []string{
			"yes", "sure", "of course", "certainly", "okay", "ok",
			"alright", "right", "correct", "no", "not at all",
			"i agree", "i understand",
		},
		responseKeywords: []string{
			"thank you", "thanks", "you're welcome", "my pleasure",
			"i'm sorry", "sorry", "excuse me", "pardon", "goodbye",
			"bye", "see you", "farewell",
		},
	}
}

// ClassifyPattern classifies a pattern into a conversation type
func (c *Classifier) ClassifyPattern(pattern string) models.PatternType {
	if c.isGreeting(pattern) {
		return models.PatternTypeGreeting
	}

	// Check request before question since some requests start with question words
	if c.isRequest(pattern) {
		return models.PatternTypeRequest
	}

	if c.isQuestion(pattern) {
		return models.PatternTypeQuestion
	}

	if c.isConfirmation(pattern) {
		return models.PatternTypeConfirmation
	}

	if c.isResponse(pattern) {
		return models.PatternTypeResponse
	}

	return models.PatternTypeOther
}

func (c *Classifier) isGreeting(pattern string) bool {
	patternLower := strings.ToLower(pattern)

	for _, keyword := range c.greetingKeywords {
		if strings.Contains(patternLower, keyword) {
			return true
		}
	}

	return false
}

func (c *Classifier) isQuestion(pattern string) bool {
	patternLower := strings.ToLower(pattern)

	// Check for question mark
	if strings.Contains(pattern, "?") {
		return true
	}

	// Check for question keywords at the beginning
	for _, keyword := range c.questionKeywords {
		if strings.HasPrefix(patternLower, keyword+" ") || patternLower == keyword {
			return true
		}
	}

	return false
}

func (c *Classifier) isRequest(pattern string) bool {
	patternLower := strings.ToLower(pattern)

	for _, keyword := range c.requestKeywords {
		if strings.Contains(patternLower, keyword) {
			return true
		}
	}

	return false
}

func (c *Classifier) isConfirmation(pattern string) bool {
	patternLower := strings.ToLower(pattern)

	for _, keyword := range c.confirmationKeywords {
		if strings.HasPrefix(patternLower, keyword) || patternLower == keyword {
			return true
		}
	}

	return false
}

func (c *Classifier) isResponse(pattern string) bool {
	patternLower := strings.ToLower(pattern)

	for _, keyword := range c.responseKeywords {
		if strings.Contains(patternLower, keyword) {
			return true
		}
	}

	return false
}
