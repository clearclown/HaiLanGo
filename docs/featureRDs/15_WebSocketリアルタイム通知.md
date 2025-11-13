# 機能実装: WebSocketリアルタイム通知

## 要件

### 機能要件

#### WebSocket機能
1. **リアルタイム通知**
   - OCR処理進捗の通知
   - 音声生成進捗の通知
   - 学習進捗の更新

2. **接続管理**
   - 接続の確立・切断
   - 再接続機能
   - 接続状態の監視

3. **メッセージ形式**
   - JSON形式
   - イベントタイプ別のメッセージ
   - エラーメッセージ

### 非機能要件

- **パフォーマンス**:
  - メッセージ送信: 100ms以内
  - 接続確立: 1秒以内

- **セキュリティ**:
  - ユーザー認証必須
  - メッセージの暗号化（WSS）

- **拡張性**:
  - 将来的な通知タイプ追加への対応

## 実装指示

### Step 1: テスト設計

1. **ユニットテスト（Go）**
   - WebSocket接続処理
   - メッセージ送信・受信
   - 接続管理

2. **統合テスト（Go + Vitest）**
   - WebSocket接続フロー
   - メッセージの送受信
   - 再接続処理

3. **E2Eテスト（Playwright）**
   - ブラウザでのWebSocket接続
   - リアルタイム通知の受信

テストファイルは `backend/internal/api/websocket/websocket_test.go` に配置。

### Step 2: 実装

#### 実装ファイル構造

```
backend/
├── internal/
│   ├── api/
│   │   └── websocket/
│   │       ├── handler.go           # WebSocketハンドラー
│   │       ├── hub.go               # 接続管理
│   │       └── handler_test.go
│   └── service/
│       └── notification/
│           ├── service.go           # 通知サービス
│           └── service_test.go

frontend/web/
├── hooks/
│   ├── useWebSocket.ts              # WebSocketフック
│   └── useWebSocket.test.ts         # Vitestテスト
└── e2e/
    └── websocket.spec.ts            # Playwrightテスト
```

#### APIエンドポイント

```
WS   /api/v1/ws
```

## 制約事項

- 既存の `backend/internal/models/` は変更しない（新規追加のみ）

## 完了条件

- [ ] すべてのテストが通る
- [ ] BiomeJSエラーがない
- [ ] GitHub CIが通る
- [ ] ドキュメントが更新されている
