package realtime

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// OpenAI Realtime API WebSocket endpoint
	openAIRealtimeURL = "wss://api.openai.com/v1/realtime"

	// Default model
	defaultModel = "gpt-4o-realtime-preview"
)

// OpenAIClient implements the Client interface for OpenAI Realtime API
type OpenAIClient struct {
	apiKey string
	model  string
}

// NewOpenAIClient creates a new OpenAI Realtime API client
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
		model:  defaultModel,
	}
}

// NewOpenAIClientFromEnv creates a client using environment variables
func NewOpenAIClientFromEnv() (*OpenAIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}
	return NewOpenAIClient(apiKey), nil
}

// Connect establishes a WebSocket connection to the Realtime API
func (c *OpenAIClient) Connect(ctx context.Context, config *SessionConfig) (Session, error) {
	if config == nil {
		config = DefaultSessionConfig()
	}

	model := config.Model
	if model == "" {
		model = c.model
	}

	url := fmt.Sprintf("%s?model=%s", openAIRealtimeURL, model)

	headers := http.Header{
		"Authorization": []string{"Bearer " + c.apiKey},
		"OpenAI-Beta":   []string{"realtime=v1"},
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 30 * time.Second,
	}

	conn, resp, err := dialer.DialContext(ctx, url, headers)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("failed to connect to OpenAI Realtime API (status %d): %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("failed to connect to OpenAI Realtime API: %w", err)
	}

	session := &openAISession{
		id:         uuid.New().String(),
		conn:       conn,
		config:     config,
		events:     make(chan *ServerEvent, 100),
		done:       make(chan struct{}),
		closeOnce:  sync.Once{},
	}

	// Start reading events in background
	go session.readLoop()

	// Update session configuration
	if err := session.updateSession(ctx, config); err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to configure session: %w", err)
	}

	return session, nil
}

// GetName returns the provider name
func (c *OpenAIClient) GetName() string {
	return "openai-realtime"
}

// Close closes the client
func (c *OpenAIClient) Close() error {
	return nil
}

// openAISession implements the Session interface
type openAISession struct {
	id        string
	conn      *websocket.Conn
	config    *SessionConfig
	events    chan *ServerEvent
	done      chan struct{}
	closeOnce sync.Once
	mu        sync.Mutex
}

// ID returns the session identifier
func (s *openAISession) ID() string {
	return s.id
}

// SendAudio sends audio data to the session
func (s *openAISession) SendAudio(ctx context.Context, audio []byte) error {
	// Encode audio as base64
	encoded := base64.StdEncoding.EncodeToString(audio)

	event := map[string]interface{}{
		"type":     EventInputAudioBufferAppend,
		"audio":    encoded,
	}

	return s.sendJSON(event)
}

// SendText sends a text message to the session
func (s *openAISession) SendText(ctx context.Context, text string) error {
	event := map[string]interface{}{
		"type": EventConversationItemCreate,
		"item": map[string]interface{}{
			"type": "message",
			"role": "user",
			"content": []map[string]interface{}{
				{
					"type": "input_text",
					"text": text,
				},
			},
		},
	}

	if err := s.sendJSON(event); err != nil {
		return err
	}

	// Trigger response generation
	return s.sendJSON(map[string]interface{}{
		"type": EventResponseCreate,
	})
}

// SendEvent sends a custom event to the session
func (s *openAISession) SendEvent(ctx context.Context, event *ClientEvent) error {
	data := map[string]interface{}{
		"type": event.Type,
	}
	if event.EventID != "" {
		data["event_id"] = event.EventID
	}
	if event.Payload != nil {
		// Merge payload into data
		if payloadMap, ok := event.Payload.(map[string]interface{}); ok {
			for k, v := range payloadMap {
				data[k] = v
			}
		}
	}
	return s.sendJSON(data)
}

// Receive returns a channel that receives server events
func (s *openAISession) Receive() <-chan *ServerEvent {
	return s.events
}

// Close closes the session
func (s *openAISession) Close() error {
	var err error
	s.closeOnce.Do(func() {
		close(s.done)
		err = s.conn.Close()
		close(s.events)
	})
	return err
}

// updateSession updates the session configuration
func (s *openAISession) updateSession(ctx context.Context, config *SessionConfig) error {
	session := map[string]interface{}{}

	if config.Voice != "" {
		session["voice"] = config.Voice
	}
	if config.Instructions != "" {
		session["instructions"] = config.Instructions
	}
	if config.InputAudioFormat != "" {
		session["input_audio_format"] = config.InputAudioFormat
	}
	if config.OutputAudioFormat != "" {
		session["output_audio_format"] = config.OutputAudioFormat
	}
	if config.Temperature > 0 {
		session["temperature"] = config.Temperature
	}
	if config.MaxOutputTokens > 0 {
		session["max_response_output_tokens"] = config.MaxOutputTokens
	}
	if config.TurnDetection != nil {
		session["turn_detection"] = map[string]interface{}{
			"type":               config.TurnDetection.Type,
			"threshold":          config.TurnDetection.Threshold,
			"prefix_padding_ms":  config.TurnDetection.PrefixPaddingMs,
			"silence_duration_ms": config.TurnDetection.SilenceDurationMs,
		}
	}
	if len(config.Tools) > 0 {
		session["tools"] = config.Tools
	}

	event := map[string]interface{}{
		"type":    EventSessionUpdate,
		"session": session,
	}

	return s.sendJSON(event)
}

// sendJSON sends a JSON message over the WebSocket
func (s *openAISession) sendJSON(v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.conn.WriteJSON(v)
}

// readLoop continuously reads events from the WebSocket
func (s *openAISession) readLoop() {
	defer func() {
		s.Close()
	}()

	for {
		select {
		case <-s.done:
			return
		default:
		}

		_, message, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.events <- &ServerEvent{
					Type: EventError,
					Error: &ErrorDetails{
						Type:    "connection_error",
						Message: err.Error(),
					},
				}
			}
			return
		}

		event, err := s.parseServerEvent(message)
		if err != nil {
			s.events <- &ServerEvent{
				Type: EventError,
				Error: &ErrorDetails{
					Type:    "parse_error",
					Message: err.Error(),
				},
			}
			continue
		}

		select {
		case s.events <- event:
		case <-s.done:
			return
		default:
			// Channel full, drop oldest event
			select {
			case <-s.events:
				s.events <- event
			default:
			}
		}
	}
}

// parseServerEvent parses a raw JSON message into a ServerEvent
func (s *openAISession) parseServerEvent(data []byte) (*ServerEvent, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse event: %w", err)
	}

	event := &ServerEvent{}

	// Parse type
	if typeData, ok := raw["type"]; ok {
		var eventType string
		if err := json.Unmarshal(typeData, &eventType); err == nil {
			event.Type = eventType
		}
	}

	// Parse event_id
	if idData, ok := raw["event_id"]; ok {
		var eventID string
		if err := json.Unmarshal(idData, &eventID); err == nil {
			event.EventID = eventID
		}
	}

	// Parse error
	if errorData, ok := raw["error"]; ok {
		var errDetails ErrorDetails
		if err := json.Unmarshal(errorData, &errDetails); err == nil {
			event.Error = &errDetails
		}
	}

	// Parse delta (for streaming content)
	if deltaData, ok := raw["delta"]; ok {
		var deltaStr string
		if err := json.Unmarshal(deltaData, &deltaStr); err == nil {
			// Handle audio delta (base64 encoded)
			if event.Type == EventResponseAudioDelta {
				if decoded, err := base64.StdEncoding.DecodeString(deltaStr); err == nil {
					event.Delta = &DeltaContent{Audio: decoded}
				}
			} else if event.Type == EventResponseTextDelta || event.Type == EventResponseAudioTranscriptDelta {
				event.Delta = &DeltaContent{Text: deltaStr}
			}
		}
	}

	// Parse transcript
	if transcriptData, ok := raw["transcript"]; ok {
		var transcript string
		if err := json.Unmarshal(transcriptData, &transcript); err == nil {
			event.Transcript = transcript
		}
	}

	// Parse response
	if responseData, ok := raw["response"]; ok {
		var response ResponseData
		if err := json.Unmarshal(responseData, &response); err == nil {
			event.Response = &response
		}
	}

	return event, nil
}

// CommitAudioBuffer commits the current audio buffer and triggers processing
func (s *openAISession) CommitAudioBuffer(ctx context.Context) error {
	return s.sendJSON(map[string]interface{}{
		"type": EventInputAudioBufferCommit,
	})
}

// ClearAudioBuffer clears the current audio buffer
func (s *openAISession) ClearAudioBuffer(ctx context.Context) error {
	return s.sendJSON(map[string]interface{}{
		"type": EventInputAudioBufferClear,
	})
}

// TriggerResponse triggers a response from the model
func (s *openAISession) TriggerResponse(ctx context.Context) error {
	return s.sendJSON(map[string]interface{}{
		"type": EventResponseCreate,
	})
}

// CancelResponse cancels the current response
func (s *openAISession) CancelResponse(ctx context.Context) error {
	return s.sendJSON(map[string]interface{}{
		"type": EventResponseCancel,
	})
}
