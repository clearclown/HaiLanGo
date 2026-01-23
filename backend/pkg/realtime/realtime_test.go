package realtime

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Use mock for tests
	os.Setenv("TEST_USE_MOCKS", "true")
	code := m.Run()
	os.Exit(code)
}

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.GetName() != "mock-realtime" {
		t.Errorf("GetName() = %v, want %v", client.GetName(), "mock-realtime")
	}
}

func TestMockClient_Connect(t *testing.T) {
	client := NewMockClient()
	ctx := context.Background()

	session, err := client.Connect(ctx, nil)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer session.Close()

	if session.ID() == "" {
		t.Error("Session ID is empty")
	}

	// Wait for session created event
	select {
	case event := <-session.Receive():
		if event.Type != EventSessionCreated {
			t.Errorf("Expected %s event, got %s", EventSessionCreated, event.Type)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for session created event")
	}
}

func TestMockSession_SendText(t *testing.T) {
	client := NewMockClient()
	ctx := context.Background()

	session, err := client.Connect(ctx, nil)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer session.Close()

	// Drain session created event
	<-session.Receive()

	// Send text
	err = session.SendText(ctx, "Hello, how are you?")
	if err != nil {
		t.Fatalf("SendText() error = %v", err)
	}

	// Wait for response
	var gotResponse bool
	timeout := time.After(2 * time.Second)

	for !gotResponse {
		select {
		case event := <-session.Receive():
			if event.Type == EventResponseDone {
				gotResponse = true
			}
		case <-timeout:
			t.Fatal("Timeout waiting for response")
		}
	}
}

func TestMockSession_SendAudio(t *testing.T) {
	client := NewMockClient()
	ctx := context.Background()

	session, err := client.Connect(ctx, nil)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer session.Close()

	// Drain session created event
	<-session.Receive()

	// Send audio (mock data)
	audio := make([]byte, 1024)
	err = session.SendAudio(ctx, audio)
	if err != nil {
		t.Fatalf("SendAudio() error = %v", err)
	}

	// Wait for response events
	var gotSpeechStarted, gotSpeechStopped, gotResponse bool
	timeout := time.After(3 * time.Second)

	for !gotResponse {
		select {
		case event := <-session.Receive():
			switch event.Type {
			case EventInputAudioBufferSpeechStarted:
				gotSpeechStarted = true
			case EventInputAudioBufferSpeechStopped:
				gotSpeechStopped = true
			case EventResponseDone:
				gotResponse = true
			}
		case <-timeout:
			t.Fatal("Timeout waiting for response")
		}
	}

	if !gotSpeechStarted {
		t.Error("Did not receive speech started event")
	}
	if !gotSpeechStopped {
		t.Error("Did not receive speech stopped event")
	}
}

func TestDefaultSessionConfig(t *testing.T) {
	config := DefaultSessionConfig()

	if config.Model != "gpt-4o-realtime-preview" {
		t.Errorf("Model = %v, want %v", config.Model, "gpt-4o-realtime-preview")
	}
	if config.Voice != "alloy" {
		t.Errorf("Voice = %v, want %v", config.Voice, "alloy")
	}
	if config.InputAudioFormat != "pcm16" {
		t.Errorf("InputAudioFormat = %v, want %v", config.InputAudioFormat, "pcm16")
	}
	if config.TurnDetection == nil {
		t.Error("TurnDetection is nil")
	}
	if config.TurnDetection.Type != "server_vad" {
		t.Errorf("TurnDetection.Type = %v, want %v", config.TurnDetection.Type, "server_vad")
	}
}

func TestLanguageLearningConfig(t *testing.T) {
	config := LanguageLearningConfig("russian", "japanese")

	if config.Instructions == "" {
		t.Error("Instructions is empty")
	}
	if config.Temperature != 0.7 {
		t.Errorf("Temperature = %v, want %v", config.Temperature, 0.7)
	}
}

func TestMockLanguageLearningResponses(t *testing.T) {
	tests := []struct {
		language string
		wantLen  int
	}{
		{"ru", 2},
		{"russian", 2},
		{"zh", 2},
		{"chinese", 2},
		{"ja", 2},
		{"japanese", 2},
		{"unknown", 5}, // Default responses
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			responses := MockLanguageLearningResponses(tt.language)
			if len(responses) != tt.wantLen {
				t.Errorf("MockLanguageLearningResponses(%s) returned %d responses, want %d",
					tt.language, len(responses), tt.wantLen)
			}
		})
	}
}

func TestIsUsingMock(t *testing.T) {
	// With TEST_USE_MOCKS set, should be true
	if !IsUsingMock() {
		t.Error("IsUsingMock() = false, want true (TEST_USE_MOCKS is set)")
	}
}

func TestCreateTestAudioBase64(t *testing.T) {
	audio := CreateTestAudioBase64()
	if audio == "" {
		t.Error("CreateTestAudioBase64() returned empty string")
	}
}

func TestSession_Close(t *testing.T) {
	client := NewMockClient()
	ctx := context.Background()

	session, err := client.Connect(ctx, nil)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	err = session.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Should be safe to close again
	err = session.Close()
	if err != nil {
		t.Errorf("Second Close() error = %v", err)
	}
}
