package realtime

import (
	"os"
)

// NewClient creates a new Realtime API client based on environment configuration
// If USE_MOCK_APIS=true or OPENAI_API_KEY is not set, returns a mock client
func NewClient() (Client, error) {
	useMocks := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	if useMocks {
		return NewMockClient(), nil
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// No API key, fallback to mock
		return NewMockClient(), nil
	}

	return NewOpenAIClient(apiKey), nil
}

// NewClientWithLanguage creates a client configured for a specific language pair
func NewClientWithLanguage(targetLanguage, nativeLanguage string) (Client, error) {
	useMocks := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	if useMocks {
		responses := MockLanguageLearningResponses(targetLanguage)
		return NewMockClientWithResponses(responses), nil
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		responses := MockLanguageLearningResponses(targetLanguage)
		return NewMockClientWithResponses(responses), nil
	}

	return NewOpenAIClient(apiKey), nil
}

// MustNewClient creates a new client and panics on error
// Use only during initialization
func MustNewClient() Client {
	client, err := NewClient()
	if err != nil {
		panic("failed to create realtime client: " + err.Error())
	}
	return client
}

// IsUsingMock returns true if the client will use mock implementation
func IsUsingMock() bool {
	useMocks := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	if useMocks {
		return true
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	return apiKey == ""
}
