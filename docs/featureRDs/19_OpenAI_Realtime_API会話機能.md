# 機能実装: OpenAI Realtime API会話機能

実装日: 2026-01-23

## 概要

OpenAI Realtime APIを使用したリアルタイム音声会話機能。運転中などハンズフリーでの言語学習を可能にする。

## 要件

### 機能要件

1. **リアルタイム音声会話**
   - WebSocket経由での音声ストリーミング
   - 低遅延の応答（目標: 1秒以内）
   - 言語学習に特化した会話フロー

2. **対応言語**
   - 学習先言語: ロシア語、中国語、日本語、英語、スペイン語、フランス語、ドイツ語、ペルシャ語、ヘブライ語、トルコ語
   - 母国語: 上記言語をサポート

3. **セッション管理**
   - セッションの開始・停止
   - セッション状態の確認
   - 認証ユーザーのみ利用可能

### 非機能要件

- **パフォーマンス**: 音声応答の遅延は1秒以内
- **セキュリティ**: JWT認証、セッションごとのユーザー検証
- **スケーラビリティ**: 複数同時セッション対応
- **モック対応**: APIキーなしでの開発・テスト対応

## 実装詳細

### バックエンド

#### ファイル構造

```
backend/
├── pkg/
│   └── realtime/
│       ├── realtime.go      # インターフェース定義
│       ├── openai.go        # OpenAI Realtime API実装
│       ├── mock.go          # モック実装
│       ├── factory.go       # ファクトリー関数
│       └── realtime_test.go # テスト
└── internal/
    └── api/
        └── handler/
            └── conversation.go  # HTTPハンドラー
```

#### APIエンドポイント

```
POST /api/v1/conversation/start   # セッション開始
POST /api/v1/conversation/stop    # セッション停止
GET  /api/v1/conversation/status  # セッション状態確認
GET  /api/v1/conversation/ws      # WebSocket接続
```

#### WebSocketメッセージ形式

**クライアント→サーバー**
```json
{
  "type": "audio",
  "audio": "<base64-encoded-audio>"
}
```

**サーバー→クライアント**
```json
{
  "type": "response.audio.delta",
  "audio": "<base64-encoded-audio>",
  "text": "レスポンステキスト"
}
```

### フロントエンド

#### ファイル構造

```
frontend/web/
├── app/
│   └── conversation/
│       └── page.tsx          # 会話ページ
├── lib/
│   └── api/
│       └── client.ts         # conversation API メソッド追加
└── components/
    └── layout/
        ├── Sidebar.tsx       # ナビゲーションリンク追加
        └── BottomNav.tsx     # ナビゲーションリンク追加
```

#### 機能

- 言語選択（学習先言語、母国語）
- セッション開始・停止
- 音声録音・送信
- 音声レベル可視化
- チャットメッセージ表示
- 音声再生

## 環境変数

```bash
# OpenAI Realtime API
OPENAI_API_KEY=your_api_key_here

# モック使用（APIキーなしで開発）
USE_MOCK_APIS=true
```

## テスト

```bash
# バックエンドテスト
cd backend
go test ./pkg/realtime/...

# フロントエンド型チェック
cd frontend/web
pnpm run type-check
```

## 使用方法

1. 言語設定を選択
2. 「Start Conversation」をクリック
3. マイクボタンをクリックして話す
4. AIの応答を聞いて会話を続ける
5. 「End Session」で終了

## 制限事項

- 現在はPCブラウザでの使用を想定
- バックグラウンド再生は未対応（将来実装予定）
- オフライン使用は不可（リアルタイムAPI依存）

## 将来の拡張

- 発音スコアリングとの統合
- 会話履歴の保存
- モバイルアプリ対応
- バックグラウンド再生対応
