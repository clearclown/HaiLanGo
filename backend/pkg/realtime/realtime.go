// Package realtime provides OpenAI Realtime API integration for voice conversations
package realtime

import (
	"context"
	"io"
)

// Client is the interface for real-time voice conversation services
type Client interface {
	// Connect establishes a WebSocket connection to the Realtime API
	Connect(ctx context.Context, config *SessionConfig) (Session, error)

	// GetName returns the provider name
	GetName() string

	// Close closes the client and releases resources
	Close() error
}

// Session represents an active real-time conversation session
type Session interface {
	// ID returns the session identifier
	ID() string

	// SendAudio sends audio data to the session
	SendAudio(ctx context.Context, audio []byte) error

	// SendText sends a text message to the session
	SendText(ctx context.Context, text string) error

	// SendEvent sends a custom event to the session
	SendEvent(ctx context.Context, event *ClientEvent) error

	// Receive returns a channel that receives server events
	Receive() <-chan *ServerEvent

	// Close closes the session
	Close() error
}

// SessionConfig configures a real-time session
type SessionConfig struct {
	// Model to use (e.g., "gpt-4o-realtime-preview")
	Model string `json:"model"`

	// Voice for TTS output (alloy, echo, fable, onyx, nova, shimmer)
	Voice string `json:"voice"`

	// Instructions (system prompt) for the session
	Instructions string `json:"instructions"`

	// InputAudioFormat (pcm16, g711_ulaw, g711_alaw)
	InputAudioFormat string `json:"input_audio_format"`

	// OutputAudioFormat (pcm16, g711_ulaw, g711_alaw)
	OutputAudioFormat string `json:"output_audio_format"`

	// Temperature for response generation (0.6-1.2)
	Temperature float64 `json:"temperature"`

	// MaxOutputTokens limits response length
	MaxOutputTokens int `json:"max_output_tokens"`

	// TurnDetection configuration
	TurnDetection *TurnDetectionConfig `json:"turn_detection"`

	// Tools available to the model
	Tools []Tool `json:"tools,omitempty"`
}

// TurnDetectionConfig configures voice activity detection
type TurnDetectionConfig struct {
	// Type of turn detection ("server_vad" or nil for manual)
	Type string `json:"type"`

	// Threshold for VAD (0.0-1.0)
	Threshold float64 `json:"threshold"`

	// PrefixPaddingMs adds audio before speech start
	PrefixPaddingMs int `json:"prefix_padding_ms"`

	// SilenceDurationMs before ending turn
	SilenceDurationMs int `json:"silence_duration_ms"`
}

// Tool represents a function the model can call
type Tool struct {
	Type     string      `json:"type"` // "function"
	Function ToolFunction `json:"function"`
}

// ToolFunction describes a callable function
type ToolFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

// ClientEvent represents an event sent to the server
type ClientEvent struct {
	Type      string      `json:"type"`
	EventID   string      `json:"event_id,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
}

// ServerEvent represents an event received from the server
type ServerEvent struct {
	Type      string          `json:"type"`
	EventID   string          `json:"event_id,omitempty"`
	SessionID string          `json:"session_id,omitempty"`
	Error     *ErrorDetails   `json:"error,omitempty"`
	Audio     []byte          `json:"audio,omitempty"`
	Text      string          `json:"text,omitempty"`
	Transcript string         `json:"transcript,omitempty"`
	Delta     *DeltaContent   `json:"delta,omitempty"`
	Response  *ResponseData   `json:"response,omitempty"`
}

// ErrorDetails contains error information
type ErrorDetails struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Param   string `json:"param,omitempty"`
}

// DeltaContent represents incremental content
type DeltaContent struct {
	Audio      []byte `json:"audio,omitempty"`
	Text       string `json:"text,omitempty"`
	Transcript string `json:"transcript,omitempty"`
}

// ResponseData contains response metadata
type ResponseData struct {
	ID            string   `json:"id"`
	Status        string   `json:"status"`
	StatusDetails *StatusDetails `json:"status_details,omitempty"`
	Output        []OutputItem `json:"output,omitempty"`
	Usage         *UsageData   `json:"usage,omitempty"`
}

// StatusDetails contains status information
type StatusDetails struct {
	Type   string `json:"type"`
	Reason string `json:"reason,omitempty"`
	Error  *ErrorDetails `json:"error,omitempty"`
}

// OutputItem represents an output item in the response
type OutputItem struct {
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Status  string    `json:"status"`
	Content []Content `json:"content,omitempty"`
}

// Content represents content within an output item
type Content struct {
	Type       string `json:"type"`
	Text       string `json:"text,omitempty"`
	Audio      []byte `json:"audio,omitempty"`
	Transcript string `json:"transcript,omitempty"`
}

// UsageData contains token usage information
type UsageData struct {
	TotalTokens      int `json:"total_tokens"`
	InputTokens      int `json:"input_tokens"`
	OutputTokens     int `json:"output_tokens"`
	InputAudioTokens int `json:"input_audio_tokens"`
	OutputAudioTokens int `json:"output_audio_tokens"`
}

// AudioReader provides streaming audio data
type AudioReader interface {
	io.Reader
	// Format returns the audio format
	Format() string
	// SampleRate returns the sample rate in Hz
	SampleRate() int
}

// AudioWriter receives streaming audio data
type AudioWriter interface {
	io.Writer
	// Format returns the audio format
	Format() string
	// SampleRate returns the sample rate in Hz
	SampleRate() int
}

// Event types sent by client
const (
	EventSessionUpdate       = "session.update"
	EventInputAudioBufferAppend = "input_audio_buffer.append"
	EventInputAudioBufferCommit = "input_audio_buffer.commit"
	EventInputAudioBufferClear  = "input_audio_buffer.clear"
	EventConversationItemCreate = "conversation.item.create"
	EventResponseCreate         = "response.create"
	EventResponseCancel         = "response.cancel"
)

// Event types received from server
const (
	EventError                     = "error"
	EventSessionCreated            = "session.created"
	EventSessionUpdated            = "session.updated"
	EventConversationCreated       = "conversation.created"
	EventConversationItemCreated   = "conversation.item.created"
	EventInputAudioBufferCommitted = "input_audio_buffer.committed"
	EventInputAudioBufferCleared   = "input_audio_buffer.cleared"
	EventInputAudioBufferSpeechStarted = "input_audio_buffer.speech_started"
	EventInputAudioBufferSpeechStopped = "input_audio_buffer.speech_stopped"
	EventResponseCreated           = "response.created"
	EventResponseDone              = "response.done"
	EventResponseOutputItemAdded   = "response.output_item.added"
	EventResponseOutputItemDone    = "response.output_item.done"
	EventResponseContentPartAdded  = "response.content_part.added"
	EventResponseContentPartDone   = "response.content_part.done"
	EventResponseTextDelta         = "response.text.delta"
	EventResponseTextDone          = "response.text.done"
	EventResponseAudioDelta        = "response.audio.delta"
	EventResponseAudioDone         = "response.audio.done"
	EventResponseAudioTranscriptDelta = "response.audio_transcript.delta"
	EventResponseAudioTranscriptDone  = "response.audio_transcript.done"
	EventRateLimitsUpdated         = "rate_limits.updated"
)

// DefaultSessionConfig returns a sensible default configuration
func DefaultSessionConfig() *SessionConfig {
	return &SessionConfig{
		Model:             "gpt-4o-realtime-preview",
		Voice:             "alloy",
		InputAudioFormat:  "pcm16",
		OutputAudioFormat: "pcm16",
		Temperature:       0.8,
		MaxOutputTokens:   4096,
		TurnDetection: &TurnDetectionConfig{
			Type:              "server_vad",
			Threshold:         0.5,
			PrefixPaddingMs:   300,
			SilenceDurationMs: 500,
		},
	}
}

// LanguageLearningConfig returns configuration optimized for language learning
func LanguageLearningConfig(targetLanguage, nativeLanguage string) *SessionConfig {
	config := DefaultSessionConfig()
	config.Instructions = buildLanguageLearningPrompt(targetLanguage, nativeLanguage)
	config.Temperature = 0.7 // More consistent for learning
	return config
}

// buildLanguageLearningPrompt creates a system prompt for language learning
func buildLanguageLearningPrompt(targetLanguage, nativeLanguage string) string {
	return `You are an expert language tutor helping the user learn ` + targetLanguage + `.
The user's native language is ` + nativeLanguage + `.

Your role:
1. Listen to the user's pronunciation and provide gentle corrections
2. Engage in natural conversation practice
3. When the user makes mistakes, correct them kindly and explain briefly
4. Speak clearly and at an appropriate pace for learning
5. Encourage the user and celebrate their progress
6. Provide pronunciation tips when helpful
7. Use the target language primarily, but explain in the native language when needed

Keep responses conversational and natural. Focus on practical communication skills.
If the user seems stuck, offer helpful prompts or switch to simpler phrases.`
}
