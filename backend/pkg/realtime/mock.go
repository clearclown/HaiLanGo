package realtime

import (
	"context"
	"encoding/base64"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MockClient implements the Client interface for testing
type MockClient struct {
	responses []MockResponse
	delay     time.Duration
}

// MockResponse represents a predefined response for testing
type MockResponse struct {
	Text       string
	Audio      []byte
	Transcript string
}

// NewMockClient creates a new mock client
func NewMockClient() *MockClient {
	return &MockClient{
		responses: defaultMockResponses(),
		delay:     100 * time.Millisecond,
	}
}

// NewMockClientWithResponses creates a mock client with custom responses
func NewMockClientWithResponses(responses []MockResponse) *MockClient {
	return &MockClient{
		responses: responses,
		delay:     100 * time.Millisecond,
	}
}

// Connect creates a mock session
func (c *MockClient) Connect(ctx context.Context, config *SessionConfig) (Session, error) {
	session := &mockSession{
		id:        uuid.New().String(),
		config:    config,
		events:    make(chan *ServerEvent, 100),
		done:      make(chan struct{}),
		responses: c.responses,
		delay:     c.delay,
		closeOnce: sync.Once{},
	}

	// Send session created event
	go func() {
		time.Sleep(50 * time.Millisecond)
		session.events <- &ServerEvent{
			Type:      EventSessionCreated,
			SessionID: session.id,
		}
	}()

	return session, nil
}

// GetName returns the provider name
func (c *MockClient) GetName() string {
	return "mock-realtime"
}

// Close closes the client
func (c *MockClient) Close() error {
	return nil
}

// mockSession implements the Session interface for testing
type mockSession struct {
	id           string
	config       *SessionConfig
	events       chan *ServerEvent
	done         chan struct{}
	responses    []MockResponse
	responseIdx  int
	delay        time.Duration
	closeOnce    sync.Once
	mu           sync.Mutex
}

// ID returns the session identifier
func (s *mockSession) ID() string {
	return s.id
}

// SendAudio simulates sending audio and triggers a mock response
func (s *mockSession) SendAudio(ctx context.Context, audio []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Simulate processing delay
	go func() {
		time.Sleep(s.delay)

		// Send speech detected event
		s.events <- &ServerEvent{
			Type: EventInputAudioBufferSpeechStarted,
		}

		time.Sleep(s.delay)

		// Send speech stopped event
		s.events <- &ServerEvent{
			Type: EventInputAudioBufferSpeechStopped,
		}

		// Get next response
		response := s.getNextResponse()

		// Simulate response creation
		s.events <- &ServerEvent{
			Type: EventResponseCreated,
		}

		// Send text delta
		if response.Text != "" {
			s.events <- &ServerEvent{
				Type: EventResponseTextDelta,
				Delta: &DeltaContent{
					Text: response.Text,
				},
			}
		}

		// Send audio delta (simulated)
		if len(response.Audio) > 0 {
			s.events <- &ServerEvent{
				Type: EventResponseAudioDelta,
				Delta: &DeltaContent{
					Audio: response.Audio,
				},
			}
		}

		// Send transcript
		if response.Transcript != "" {
			s.events <- &ServerEvent{
				Type: EventResponseAudioTranscriptDelta,
				Delta: &DeltaContent{
					Text: response.Transcript,
				},
			}
		}

		// Send response done
		s.events <- &ServerEvent{
			Type: EventResponseDone,
			Response: &ResponseData{
				ID:     uuid.New().String(),
				Status: "completed",
				Usage: &UsageData{
					TotalTokens:       100,
					InputTokens:       30,
					OutputTokens:      70,
					InputAudioTokens:  20,
					OutputAudioTokens: 50,
				},
			},
		}
	}()

	return nil
}

// SendText simulates sending text and triggers a mock response
func (s *mockSession) SendText(ctx context.Context, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	go func() {
		time.Sleep(s.delay)

		response := s.getNextResponse()

		s.events <- &ServerEvent{
			Type: EventResponseCreated,
		}

		if response.Text != "" {
			s.events <- &ServerEvent{
				Type: EventResponseTextDelta,
				Delta: &DeltaContent{
					Text: response.Text,
				},
			}
		}

		s.events <- &ServerEvent{
			Type: EventResponseDone,
			Response: &ResponseData{
				ID:     uuid.New().String(),
				Status: "completed",
			},
		}
	}()

	return nil
}

// SendEvent simulates sending a custom event
func (s *mockSession) SendEvent(ctx context.Context, event *ClientEvent) error {
	// For mock, just acknowledge the event
	return nil
}

// Receive returns the events channel
func (s *mockSession) Receive() <-chan *ServerEvent {
	return s.events
}

// Close closes the mock session
func (s *mockSession) Close() error {
	var err error
	s.closeOnce.Do(func() {
		close(s.done)
		close(s.events)
	})
	return err
}

// getNextResponse returns the next mock response in rotation
func (s *mockSession) getNextResponse() MockResponse {
	if len(s.responses) == 0 {
		return MockResponse{
			Text:       "Hello! I'm here to help you practice.",
			Transcript: "Hello! I'm here to help you practice.",
		}
	}

	response := s.responses[s.responseIdx]
	s.responseIdx = (s.responseIdx + 1) % len(s.responses)
	return response
}

// defaultMockResponses returns a set of default mock responses for language learning
func defaultMockResponses() []MockResponse {
	return []MockResponse{
		{
			Text:       "Great pronunciation! Let's try another phrase.",
			Transcript: "Great pronunciation! Let's try another phrase.",
			Audio:      generateMockAudio("Great pronunciation! Let's try another phrase."),
		},
		{
			Text:       "Almost perfect! Try to emphasize the last syllable a bit more.",
			Transcript: "Almost perfect! Try to emphasize the last syllable a bit more.",
			Audio:      generateMockAudio("Almost perfect! Try to emphasize the last syllable a bit more."),
		},
		{
			Text:       "Excellent! You're making great progress.",
			Transcript: "Excellent! You're making great progress.",
			Audio:      generateMockAudio("Excellent! You're making great progress."),
		},
		{
			Text:       "Good effort! Let me show you the correct pronunciation.",
			Transcript: "Good effort! Let me show you the correct pronunciation.",
			Audio:      generateMockAudio("Good effort! Let me show you the correct pronunciation."),
		},
		{
			Text:       "Perfect! Your accent is improving.",
			Transcript: "Perfect! Your accent is improving.",
			Audio:      generateMockAudio("Perfect! Your accent is improving."),
		},
	}
}

// generateMockAudio generates mock audio data (silence)
func generateMockAudio(text string) []byte {
	// Generate 1 second of silence at 24kHz, 16-bit PCM
	// This is just placeholder data for testing
	samples := 24000 // 1 second at 24kHz
	audio := make([]byte, samples*2) // 16-bit = 2 bytes per sample
	return audio
}

// MockLanguageLearningResponses returns mock responses for language learning scenarios
func MockLanguageLearningResponses(targetLanguage string) []MockResponse {
	switch targetLanguage {
	case "ru", "russian":
		return []MockResponse{
			{
				Text:       "Отлично! Ваше произношение очень хорошее. Great! Your pronunciation is very good.",
				Transcript: "Отлично! Ваше произношение очень хорошее. Great! Your pronunciation is very good.",
			},
			{
				Text:       "Попробуйте ещё раз. Обратите внимание на ударение. Try again. Pay attention to the stress.",
				Transcript: "Попробуйте ещё раз. Обратите внимание на ударение. Try again. Pay attention to the stress.",
			},
		}
	case "zh", "chinese":
		return []MockResponse{
			{
				Text:       "很好！你的声调很准确。Very good! Your tones are accurate.",
				Transcript: "很好！你的声调很准确。Very good! Your tones are accurate.",
			},
			{
				Text:       "再试一次，注意第三声。Try again, pay attention to the third tone.",
				Transcript: "再试一次，注意第三声。Try again, pay attention to the third tone.",
			},
		}
	case "ja", "japanese":
		return []MockResponse{
			{
				Text:       "すばらしい！発音がとても自然です。Wonderful! Your pronunciation is very natural.",
				Transcript: "すばらしい！発音がとても自然です。Wonderful! Your pronunciation is very natural.",
			},
			{
				Text:       "もう一度。長音に気をつけてください。Once more. Pay attention to long vowels.",
				Transcript: "もう一度。長音に気をつけてください。Once more. Pay attention to long vowels.",
			},
		}
	default:
		return defaultMockResponses()
	}
}

// CreateTestAudioBase64 creates base64 encoded test audio
func CreateTestAudioBase64() string {
	audio := generateMockAudio("test")
	return base64.StdEncoding.EncodeToString(audio)
}
