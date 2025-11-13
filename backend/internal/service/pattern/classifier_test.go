package pattern

import (
	"testing"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

func TestClassifier_ClassifyPattern(t *testing.T) {
	tests := []struct {
		name         string
		pattern      string
		expectedType models.PatternType
	}{
		{
			name:         "greeting pattern",
			pattern:      "Hello",
			expectedType: models.PatternTypeGreeting,
		},
		{
			name:         "greeting pattern - Good morning",
			pattern:      "Good morning",
			expectedType: models.PatternTypeGreeting,
		},
		{
			name:         "question pattern - How",
			pattern:      "How are you?",
			expectedType: models.PatternTypeQuestion,
		},
		{
			name:         "question pattern - What",
			pattern:      "What is your name?",
			expectedType: models.PatternTypeQuestion,
		},
		{
			name:         "question pattern - Where",
			pattern:      "Where are you from?",
			expectedType: models.PatternTypeQuestion,
		},
		{
			name:         "request pattern - Please",
			pattern:      "Please help me",
			expectedType: models.PatternTypeRequest,
		},
		{
			name:         "request pattern - Could you",
			pattern:      "Could you tell me?",
			expectedType: models.PatternTypeRequest,
		},
		{
			name:         "confirmation pattern - Yes",
			pattern:      "Yes, I agree",
			expectedType: models.PatternTypeConfirmation,
		},
		{
			name:         "confirmation pattern - Sure",
			pattern:      "Sure, no problem",
			expectedType: models.PatternTypeConfirmation,
		},
		{
			name:         "response pattern - Thank you",
			pattern:      "Thank you",
			expectedType: models.PatternTypeResponse,
		},
		{
			name:         "other pattern",
			pattern:      "The weather is nice",
			expectedType: models.PatternTypeOther,
		},
	}

	classifier := NewClassifier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patternType := classifier.ClassifyPattern(tt.pattern)

			if patternType != tt.expectedType {
				t.Errorf("ClassifyPattern(%q) = %v, want %v", tt.pattern, patternType, tt.expectedType)
			}
		})
	}
}

func TestClassifier_IsGreeting(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expected bool
	}{
		{"Hello", "Hello", true},
		{"Hi", "Hi there", true},
		{"Good morning", "Good morning everyone", true},
		{"Good afternoon", "Good afternoon", true},
		{"Good evening", "Good evening", true},
		{"Good night", "Good night", true},
		{"Hey", "Hey you", true},
		{"Greetings", "Greetings", true},
		{"Not a greeting", "How are you?", false},
		{"Not a greeting", "Please help", false},
	}

	classifier := NewClassifier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.isGreeting(tt.pattern)

			if result != tt.expected {
				t.Errorf("isGreeting(%q) = %v, want %v", tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestClassifier_IsQuestion(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expected bool
	}{
		{"How question", "How are you?", true},
		{"What question", "What is this?", true},
		{"Where question", "Where are you?", true},
		{"When question", "When will you come?", true},
		{"Why question", "Why did you do that?", true},
		{"Who question", "Who is he?", true},
		{"Which question", "Which one do you want?", true},
		{"Is question", "Is this correct?", true},
		{"Are question", "Are you okay?", true},
		{"Do question", "Do you like it?", true},
		{"Does question", "Does he know?", true},
		{"Can question", "Can you help?", true},
		{"Could question", "Could you tell me?", true},
		{"Would question", "Would you like?", true},
		{"Not a question", "Hello", false},
		{"Not a question", "Thank you", false},
	}

	classifier := NewClassifier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.isQuestion(tt.pattern)

			if result != tt.expected {
				t.Errorf("isQuestion(%q) = %v, want %v", tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestClassifier_IsRequest(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expected bool
	}{
		{"Please", "Please help me", true},
		{"Could you", "Could you tell me?", true},
		{"Would you", "Would you like to?", true},
		{"Can you", "Can you help?", true},
		{"Will you", "Will you come?", true},
		{"I would like", "I would like to know", true},
		{"I'd like", "I'd like some water", true},
		{"May I", "May I ask?", true},
		{"Not a request", "Hello", false},
		{"Not a request", "How are you?", false},
	}

	classifier := NewClassifier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.isRequest(tt.pattern)

			if result != tt.expected {
				t.Errorf("isRequest(%q) = %v, want %v", tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestClassifier_IsConfirmation(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expected bool
	}{
		{"Yes", "Yes", true},
		{"Yes with text", "Yes, I agree", true},
		{"Sure", "Sure", true},
		{"Of course", "Of course", true},
		{"Certainly", "Certainly", true},
		{"Okay", "Okay", true},
		{"OK", "OK", true},
		{"Alright", "Alright", true},
		{"Right", "Right", true},
		{"Correct", "Correct", true},
		{"No", "No", true},
		{"Not at all", "Not at all", true},
		{"Not a confirmation", "Hello", false},
		{"Not a confirmation", "How are you?", false},
	}

	classifier := NewClassifier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.isConfirmation(tt.pattern)

			if result != tt.expected {
				t.Errorf("isConfirmation(%q) = %v, want %v", tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestClassifier_IsResponse(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expected bool
	}{
		{"Thank you", "Thank you", true},
		{"Thanks", "Thanks", true},
		{"You're welcome", "You're welcome", true},
		{"My pleasure", "My pleasure", true},
		{"I'm sorry", "I'm sorry", true},
		{"Excuse me", "Excuse me", true},
		{"Pardon", "Pardon me", true},
		{"Goodbye", "Goodbye", true},
		{"Bye", "Bye", true},
		{"See you", "See you later", true},
		{"Not a response", "Hello", false},
		{"Not a response", "How are you?", false},
	}

	classifier := NewClassifier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.isResponse(tt.pattern)

			if result != tt.expected {
				t.Errorf("isResponse(%q) = %v, want %v", tt.pattern, result, tt.expected)
			}
		})
	}
}
