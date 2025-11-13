# WebSocketリアルタイム通知実装ドキュメント

## 概要

HaiLanGoプロジェクトにWebSocketベースのリアルタイム通知機能を実装しました。この機能により、OCR処理進捗、TTS音声生成進捗、学習進捗更新などをリアルタイムでクライアントに通知できます。

## アーキテクチャ

### バックエンド（Go）

```
backend/
├── internal/
│   ├── models/
│   │   └── notification.go          # 通知データモデル
│   ├── api/
│   │   └── websocket/
│   │       ├── client.go            # WebSocketクライアント
│   │       ├── hub.go               # 接続管理Hub
│   │       ├── handler.go           # HTTPハンドラー
│   │       └── hub_test.go          # テスト
│   └── service/
│       └── notification/
│           ├── service.go           # 通知サービス
│           └── service_test.go      # テスト
└── cmd/
    └── server/
        └── main.go                  # サーバーエントリーポイント
```

### フロントエンド（TypeScript/React）

```
frontend/web/
├── lib/
│   └── types/
│       └── notification.ts          # 型定義
├── hooks/
│   ├── useWebSocket.ts              # WebSocketフック
│   └── useWebSocket.test.ts         # Vitestテスト
└── e2e/
    └── websocket.spec.ts            # Playwrightテスト
```

## 通知タイプ

### 1. OCR処理進捗通知 (`ocr_progress`)

```typescript
interface OCRProgressData {
  book_id: string;
  total_pages: number;
  processed_pages: number;
  current_page: number;
  progress: number; // 0-100
  estimated_time_ms: number;
  status: "processing" | "completed" | "failed";
}
```

### 2. TTS音声生成進捗通知 (`tts_progress`)

```typescript
interface TTSProgressData {
  book_id: string;
  page_number: number;
  total_segments: number;
  processed_segments: number;
  progress: number; // 0-100
  status: "processing" | "completed" | "failed";
}
```

### 3. 学習進捗更新通知 (`learning_update`)

```typescript
interface LearningUpdateData {
  user_id: string;
  book_id: string;
  page_number: number;
  completed_pages: number;
  total_pages: number;
  learned_words: number;
  study_time_ms: number;
}
```

### 4. エラー通知 (`error`)

```typescript
interface ErrorData {
  code: string;
  message: string;
  details?: string;
}
```

### 5. Ping/Pong (`ping`, `pong`)

接続維持のための定期的なping/pongメッセージ。

## 使用方法

### バックエンド

#### 1. Hubの起動

```go
// main.go
hub := websocket.NewHub()
go hub.Run()
```

#### 2. 通知サービスの作成

```go
notificationService := notification.NewService(hub)
```

#### 3. 通知の送信

```go
// OCR進捗通知
err := notificationService.NotifyOCRProgress(userID, &models.OCRProgressData{
    BookID:         "book-123",
    TotalPages:     100,
    ProcessedPages: 50,
    CurrentPage:    50,
    Progress:       50.0,
    Status:         "processing",
})

// TTS進捗通知
err := notificationService.NotifyTTSProgress(userID, &models.TTSProgressData{
    BookID:            "book-123",
    PageNumber:        10,
    TotalSegments:     5,
    ProcessedSegments: 3,
    Progress:          60.0,
    Status:            "processing",
})

// 学習進捗更新
err := notificationService.NotifyLearningUpdate(userID, &models.LearningUpdateData{
    UserID:         userID,
    BookID:         "book-123",
    PageNumber:     20,
    CompletedPages: 20,
    TotalPages:     100,
    LearnedWords:   150,
    StudyTimeMS:    3600000,
})

// エラー通知
err := notificationService.NotifyError(userID, &models.ErrorData{
    Code:    "OCR_FAILED",
    Message: "Failed to process image",
    Details: "OCR service is temporarily unavailable",
})
```

### フロントエンド

#### 1. useWebSocketフックの使用

```typescript
import { useWebSocket } from "@/hooks/useWebSocket";

function MyComponent() {
  const { isConnected, isConnecting, error } = useWebSocket({
    url: "ws://localhost:8080/api/v1/ws",
    userId: "user-123",

    // コールバック関数
    onOCRProgress: (data) => {
      console.log("OCR Progress:", data.progress);
    },

    onTTSProgress: (data) => {
      console.log("TTS Progress:", data.progress);
    },

    onLearningUpdate: (data) => {
      console.log("Learning Update:", data.completed_pages);
    },

    onError: (data) => {
      console.error("Error:", data.message);
    },

    // オプション設定
    reconnectAttempts: 5,
    reconnectInterval: 3000, // 3秒
  });

  return (
    <div>
      {isConnecting && <p>接続中...</p>}
      {isConnected && <p>接続済み</p>}
      {error && <p>エラー: {error.message}</p>}
    </div>
  );
}
```

## APIエンドポイント

### WebSocket接続

```
WS /api/v1/ws?user_id={userId}
```

**パラメータ**:
- `user_id`: ユーザーID（必須）

**接続例**:
```javascript
const ws = new WebSocket("ws://localhost:8080/api/v1/ws?user_id=user-123");
```

### テスト用通知送信エンドポイント

```
GET /api/v1/test/notify?user_id={userId}
```

開発・テスト用のエンドポイント。指定したユーザーにテスト通知を送信します。

## メッセージフォーマット

すべてのWebSocketメッセージは以下のJSON形式です：

```json
{
  "type": "ocr_progress | tts_progress | learning_update | error | ping | pong",
  "data": { ... },
  "timestamp": "2025-11-13T12:00:00Z"
}
```

### 例: OCR進捗通知

```json
{
  "type": "ocr_progress",
  "data": {
    "book_id": "book-123",
    "total_pages": 100,
    "processed_pages": 50,
    "current_page": 50,
    "progress": 50.0,
    "estimated_time_ms": 60000,
    "status": "processing"
  },
  "timestamp": "2025-11-13T12:00:00Z"
}
```

## 接続管理

### 自動再接続

フロントエンドの`useWebSocket`フックは、接続が切断された場合に自動的に再接続を試みます。

- デフォルトの再接続試行回数: 5回
- デフォルトの再接続間隔: 3秒
- 再接続間隔は指数バックオフではなく固定間隔

### Ping/Pong

接続を維持するため、30秒ごとにpingメッセージを送信します。

## テスト

### バックエンドテスト

```bash
# すべてのテストを実行
go test ./...

# WebSocketテストのみ
go test ./internal/api/websocket/... -v

# 通知サービステストのみ
go test ./internal/service/notification/... -v
```

### フロントエンドテスト

```bash
# Vitestユニットテスト
pnpm test:unit

# Playwright E2Eテスト
pnpm test:e2e

# すべてのテスト
pnpm test
```

## パフォーマンス

### タイムアウト設定

- **writeWait**: 10秒 - メッセージ送信のタイムアウト
- **pongWait**: 60秒 - pongメッセージ受信のタイムアウト
- **pingPeriod**: 54秒 - ping送信間隔（pongWaitの9/10）
- **maxMessageSize**: 512バイト - 受信メッセージの最大サイズ

### 同時接続

Hubは複数のクライアント接続を効率的に管理します。同じユーザーが複数のデバイスから接続した場合、すべてのクライアントに通知が送信されます。

## セキュリティ

### 認証

現在の実装では、クエリパラメータ`user_id`でユーザーを識別していますが、本番環境ではJWTトークンベースの認証を実装する必要があります。

**TODO（本番環境）**:
```go
// handler.go
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
    // JWTトークンの検証
    token := r.Header.Get("Authorization")
    userID, err := validateJWT(token)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // WebSocket接続の確立
    // ...
}
```

### CORS

現在、すべてのオリジンを許可していますが、本番環境では適切なオリジンチェックを実装する必要があります。

```go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        // 本番環境では適切なオリジンチェックを実装
        origin := r.Header.Get("Origin")
        return origin == "https://yourdomain.com"
    },
}
```

## トラブルシューティング

### 接続エラー

1. **サーバーが起動しているか確認**
   ```bash
   curl http://localhost:8080/health
   ```

2. **WebSocketエンドポイントが正しいか確認**
   ```
   ws://localhost:8080/api/v1/ws?user_id=test-user
   ```

3. **ファイアウォールがポート8080をブロックしていないか確認**

### 通知が届かない

1. **ユーザーが接続されているか確認**
   ```go
   isConnected := notificationService.IsUserConnected(userID)
   ```

2. **通知送信時のエラーをログで確認**

3. **ブラウザのコンソールでWebSocketメッセージを確認**

## 今後の改善点

1. **認証の強化**: JWTトークンベースの認証実装
2. **スケーリング**: Redis Pub/Subを使用した複数サーバー間の通知配信
3. **メッセージ永続化**: 未配信メッセージの保存と再送
4. **圧縮**: 大きなメッセージの圧縮
5. **メトリクス**: 接続数、メッセージ送信数などの監視

## 参考リンク

- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [WebSocket API (MDN)](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)
- [React Hooks](https://react.dev/reference/react)
