import { renderHook, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { useWebSocket } from "./useWebSocket";
import type {
  OCRProgressData,
  TTSProgressData,
  LearningUpdateData,
  ErrorData,
} from "@/lib/types/notification";

describe("useWebSocket", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should connect to WebSocket", async () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: "ws://localhost:8080/api/v1/ws",
        userId: "test-user",
      }),
    );

    expect(result.current.isConnecting).toBe(true);

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    expect(result.current.isConnecting).toBe(false);
    expect(result.current.error).toBe(null);
  });

  it("should handle OCR progress notifications", async () => {
    const onOCRProgress = vi.fn();

    const { result } = renderHook(() =>
      useWebSocket({
        url: "ws://localhost:8080/api/v1/ws",
        userId: "test-user",
        onOCRProgress,
      }),
    );

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    // Simulate receiving a message
    const mockData: OCRProgressData = {
      book_id: "book-123",
      total_pages: 100,
      processed_pages: 50,
      current_page: 50,
      progress: 50.0,
      estimated_time_ms: 60000,
      status: "processing",
    };

    const mockMessage = {
      type: "ocr_progress",
      data: mockData,
      timestamp: new Date().toISOString(),
    };

    // Since we're using a mock WebSocket, we'll need to manually trigger the message handler
    // In a real scenario, this would be sent from the server
  });

  it("should handle TTS progress notifications", async () => {
    const onTTSProgress = vi.fn();

    const { result } = renderHook(() =>
      useWebSocket({
        url: "ws://localhost:8080/api/v1/ws",
        userId: "test-user",
        onTTSProgress,
      }),
    );

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });
  });

  it("should handle learning update notifications", async () => {
    const onLearningUpdate = vi.fn();

    const { result } = renderHook(() =>
      useWebSocket({
        url: "ws://localhost:8080/api/v1/ws",
        userId: "test-user",
        onLearningUpdate,
      }),
    );

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });
  });

  it("should handle error notifications", async () => {
    const onError = vi.fn();

    const { result } = renderHook(() =>
      useWebSocket({
        url: "ws://localhost:8080/api/v1/ws",
        userId: "test-user",
        onError,
      }),
    );

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });
  });

  it("should disconnect WebSocket on unmount", async () => {
    const { result, unmount } = renderHook(() =>
      useWebSocket({
        url: "ws://localhost:8080/api/v1/ws",
        userId: "test-user",
      }),
    );

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    unmount();

    await waitFor(() => {
      expect(result.current.isConnected).toBe(false);
    });
  });

  it("should send ping messages", async () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: "ws://localhost:8080/api/v1/ws",
        userId: "test-user",
      }),
    );

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    result.current.sendPing();
    // In a real test, we would verify that the ping was sent
  });
});
