# WebSocket Real-Time Notifications Implementation

## Overview

This document describes the complete WebSocket implementation for the HaiLanGo project, providing real-time notifications for OCR progress, book readiness, review reminders, learning updates, and general notifications.

## Architecture

### Backend (Go)

#### Components

1. **Hub (`internal/websocket/hub.go`)**
   - Manages all active WebSocket connections
   - Handles client registration and deregistration
   - Broadcasts messages to all clients or specific users
   - Implements heartbeat/ping-pong mechanism
   - Thread-safe operations with mutex locks

2. **Message Types (`internal/websocket/message.go`)**
   - Defines all message types and payload structures
   - Helper functions for creating typed messages
   - Includes:
     - `ocr_progress`: OCR processing updates
     - `book_ready`: Book preparation complete
     - `review_reminder`: Review items due
     - `learning_update`: Learning session statistics
     - `notification`: General notifications
     - `error`: Error notifications
     - `connection_established`: Initial connection confirmation

3. **WebSocket Handler (`internal/api/handler/websocket.go`)**
   - Handles WebSocket upgrade from HTTP
   - Validates JWT authentication (from query param or header)
   - Creates and registers clients
   - Provides debug endpoint for connection statistics

4. **Router Integration (`internal/api/router/router.go`)**
   - Initializes WebSocket hub as singleton
   - Registers WebSocket endpoint: `GET /api/v1/ws`
   - Provides stats endpoint: `GET /api/v1/ws/stats`

### Frontend (TypeScript/React)

#### Components

1. **WebSocket Client (`lib/websocket/client.ts`)**
   - Singleton WebSocket client
   - Automatic reconnection with exponential backoff
   - Heartbeat mechanism
   - Type-safe message handlers
   - Connection state management

2. **Type Definitions (`lib/websocket/types.ts`)**
   - TypeScript interfaces matching backend message types
   - Strongly typed payload structures
   - Message handler types

3. **React Hooks (`hooks/useWebSocket.ts`)**
   - `useWebSocket`: Main hook for connection management
   - `useWebSocketSubscription`: Subscribe to specific message types
   - `useWebSocketSubscriptions`: Subscribe to multiple message types
   - Typed hooks for specific messages:
     - `useOCRProgress`
     - `useBookReady`
     - `useReviewReminder`
     - `useLearningUpdate`
     - `useNotification`
     - `useErrorNotification`

4. **Provider Component (`components/WebSocketProvider.tsx`)**
   - Context provider for global WebSocket state
   - Automatic connection on mount
   - Easy integration into app

5. **Example Component (`components/examples/WebSocketExample.tsx`)**
   - Demonstrates usage of WebSocket hooks
   - Shows how to handle different message types

## Usage

### Backend

#### Sending Messages to Users

```go
package service

import (
    "github.com/clearclown/HaiLanGo/backend/internal/websocket"
    "github.com/google/uuid"
)

func (s *OCRService) ProcessPage(bookID uuid.UUID, page int) error {
    // Your OCR processing logic...

    // Send progress notification
    message, err := websocket.NewOCRProgressMessage(
        bookID,
        totalPages,
        processedPages,
        "processing",
        "Processing page...",
    )
    if err != nil {
        return err
    }

    // Send to user
    s.wsHub.SendToUser(userID, message)

    return nil
}
```

#### Broadcasting to All Users

```go
// Send notification to all connected users
message, err := websocket.NewNotificationMessage(
    "System Update",
    "System maintenance scheduled for tonight",
    websocket.NotificationLevelInfo,
)
if err != nil {
    return err
}

s.wsHub.BroadcastToAll(message)
```

### Frontend

#### Using the Provider

```tsx
// app/layout.tsx
import { WebSocketProvider } from '@/components/WebSocketProvider';

export default function RootLayout({ children }) {
  return (
    <html>
      <body>
        <WebSocketProvider autoConnect={true}>
          {children}
        </WebSocketProvider>
      </body>
    </html>
  );
}
```

#### Subscribing to Messages in Components

```tsx
'use client';

import { useOCRProgress, useBookReady } from '@/hooks/useWebSocket';

export function BookUploadComponent() {
  // Handle OCR progress updates
  useOCRProgress((payload) => {
    console.log('OCR Progress:', payload.progress);
    // Update UI with progress
    setProgress(payload.progress);
  });

  // Handle book ready notification
  useBookReady((payload) => {
    console.log('Book Ready:', payload.title);
    // Show success message
    toast.success(`${payload.title} is ready!`);
  });

  return <div>Upload UI...</div>;
}
```

#### Manual Connection Control

```tsx
'use client';

import { useWebSocket } from '@/hooks/useWebSocket';

export function MyComponent() {
  const { connected, connect, disconnect, subscribe } = useWebSocket();

  // Manual connection with token
  const handleConnect = () => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      connect(token);
    }
  };

  // Subscribe to specific message types
  useEffect(() => {
    const unsubscribe = subscribe('notification', (payload) => {
      console.log('Notification:', payload);
    });

    return unsubscribe; // Cleanup on unmount
  }, [subscribe]);

  return (
    <div>
      <p>Status: {connected ? 'Connected' : 'Disconnected'}</p>
      <button onClick={handleConnect}>Connect</button>
      <button onClick={disconnect}>Disconnect</button>
    </div>
  );
}
```

## Authentication

WebSocket connections require JWT authentication:

1. **Query Parameter**: `ws://localhost:8080/api/v1/ws?token=YOUR_JWT_TOKEN`
2. **Authorization Header**: `Authorization: Bearer YOUR_JWT_TOKEN`

The JWT token is validated on connection upgrade and must contain:
- `user_id`: UUID of the authenticated user
- `email`: User's email address

## Connection Lifecycle

### Client-Side

1. **Connection**: Client connects with JWT token
2. **Heartbeat**: Automatic ping every 30 seconds
3. **Message Handling**: Messages dispatched to registered handlers
4. **Reconnection**: Automatic reconnection on disconnect (max 5 attempts with exponential backoff)
5. **Cleanup**: Proper cleanup on component unmount

### Server-Side

1. **Upgrade**: HTTP connection upgraded to WebSocket
2. **Registration**: Client registered in Hub
3. **Connection Established**: Initial message sent to client
4. **Read/Write Pumps**: Goroutines handle bidirectional communication
5. **Heartbeat**: Ping/pong every 54 seconds (90% of 60s timeout)
6. **Deregistration**: Client removed on disconnect

## Message Flow

```
┌─────────────────┐                    ┌──────────────────┐
│   Backend       │                    │    Frontend      │
│   Service       │                    │    Component     │
└────────┬────────┘                    └────────┬─────────┘
         │                                      │
         │  1. Create Message                   │
         │  (e.g., OCRProgress)                 │
         │                                      │
         ▼                                      │
┌─────────────────┐                             │
│   WebSocket     │                             │
│      Hub        │                             │
└────────┬────────┘                             │
         │                                      │
         │  2. Send to User                     │
         │                                      │
         ▼                                      │
┌─────────────────┐                             │
│   WebSocket     │                             │
│    Client       │                             │
└────────┬────────┘                             │
         │                                      │
         │  3. Transmit over WS                 │
         │────────────────────────────────────► │
         │                                      │
         │                                      ▼
         │                             ┌─────────────────┐
         │                             │  WebSocket      │
         │                             │  Client (TS)    │
         │                             └────────┬────────┘
         │                                      │
         │                                      │  4. Dispatch
         │                                      │
         │                                      ▼
         │                             ┌─────────────────┐
         │                             │   Message       │
         │                             │   Handler       │
         │                             └────────┬────────┘
         │                                      │
         │                                      │  5. Update UI
         │                                      │
         │                                      ▼
         │                             ┌─────────────────┐
         │                             │   React         │
         │                             │   Component     │
         │                             └─────────────────┘
```

## Configuration

### Backend

```go
// WebSocket configuration constants
const (
    WriteWait      = 10 * time.Second  // Message send timeout
    PongWait       = 60 * time.Second  // Pong message timeout
    PingPeriod     = 54 * time.Second  // Ping interval
    MaxMessageSize = 512 * 1024        // 512KB max message size
)
```

### Frontend

```typescript
// Default WebSocket configuration
const config = {
  url: 'ws://localhost:8080/api/v1/ws',
  reconnectInterval: 1000,
  maxReconnectAttempts: 5,
  heartbeatInterval: 30000,
};
```

## Error Handling

### Backend

- Invalid token: Returns 401 Unauthorized
- Connection errors: Logged and client disconnected
- Message parsing errors: Logged, connection continues
- Send failures: Client marked for disconnect

### Frontend

- Connection errors: Automatic reconnection
- Parse errors: Logged, message skipped
- Handler errors: Caught and logged, other handlers continue
- Max reconnect attempts: Error notification to user

## Testing

### Backend

```bash
# Build and verify
go build -o /tmp/server ./cmd/server/main.go

# Run server
./cmd/server/main.go

# Check WebSocket stats
curl http://localhost:8080/api/v1/ws/stats
```

### Frontend

```bash
# Type check
pnpm run type-check

# Build
pnpm run build

# Run development server
pnpm run dev
```

### Manual WebSocket Testing

```javascript
// Browser console
const ws = new WebSocket('ws://localhost:8080/api/v1/ws?token=YOUR_TOKEN');

ws.onopen = () => console.log('Connected');
ws.onmessage = (e) => console.log('Message:', JSON.parse(e.data));
ws.onerror = (e) => console.error('Error:', e);
ws.onclose = () => console.log('Disconnected');
```

## Integration Points

### OCR Service

Update `internal/service/ocr/service.go` to send progress notifications:

```go
type OCRService struct {
    // ... existing fields
    wsHub *websocket.Hub
}

func (s *OCRService) ProcessPage(ctx context.Context, bookID uuid.UUID, pageNum int) error {
    // ... processing logic

    // Send progress
    msg, _ := websocket.NewOCRProgressMessage(bookID, totalPages, processedPages, "processing", "")
    s.wsHub.SendToUser(userID, msg)

    return nil
}
```

### Book Service

Update `internal/service/book/service.go` to notify when books are ready:

```go
func (s *BookService) CompleteBook(ctx context.Context, bookID uuid.UUID) error {
    // ... completion logic

    // Send notification
    msg, _ := websocket.NewBookReadyMessage(bookID, book.Title, book.TotalPages)
    s.wsHub.SendToUser(book.UserID, msg)

    return nil
}
```

## Security Considerations

1. **JWT Validation**: All connections require valid JWT
2. **User Isolation**: Messages only sent to authenticated user's connections
3. **Rate Limiting**: Consider adding rate limiting for message sending
4. **Message Size**: 512KB limit prevents memory exhaustion
5. **CORS**: Configure proper CORS for production

## Performance

- **Concurrent Connections**: Hub handles multiple clients efficiently with goroutines
- **Message Batching**: Multiple queued messages sent in single write
- **Memory**: ~1KB per active connection
- **CPU**: Minimal overhead, scales horizontally

## Monitoring

### Metrics to Track

- Active connections count
- Messages sent per minute
- Connection errors
- Reconnection attempts
- Average message latency

### Debug Endpoint

```bash
# Get connection statistics
curl http://localhost:8080/api/v1/ws/stats

# Response:
{
  "connected_users": 42,
  "total_connections": 58
}
```

## Future Enhancements

- [ ] Message persistence for offline users
- [ ] Message acknowledgment system
- [ ] User presence indicators
- [ ] Typing indicators
- [ ] Read receipts
- [ ] Push notifications fallback
- [ ] Metrics dashboard
- [ ] Load balancing for multiple servers
- [ ] Redis pub/sub for distributed systems

## Troubleshooting

### Connection Not Established

1. Check JWT token validity
2. Verify WebSocket URL
3. Check CORS configuration
4. Inspect browser console for errors

### Messages Not Received

1. Verify handler is registered
2. Check message type matches
3. Inspect network tab for WebSocket frames
4. Check server logs for send errors

### Frequent Disconnections

1. Check network stability
2. Verify heartbeat settings
3. Review server resource limits
4. Check for firewall/proxy issues

## References

- [gorilla/websocket Documentation](https://pkg.go.dev/github.com/gorilla/websocket)
- [MDN WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)
- [RFC 6455 - The WebSocket Protocol](https://datatracker.ietf.org/doc/html/rfc6455)
