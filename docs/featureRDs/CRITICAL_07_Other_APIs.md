# CRITICAL_07: その他のAPI実装（P2機能）

**優先度**: P2（中優先度）
**担当者**: 未割当
**見積時間**: 各4-8時間
**依存**: OCR/TTS/STT実装完了後に着手

---

## ⚠️ PM指示

**現状**: 拡張機能のバックエンドAPIが全て未実装。
**期限**: P0/P1完了後、2週間以内に実装完了すること。
**重要**: すべてモックから始め、段階的に実APIと統合。

---

## 📋 API一覧

### 1. Dictionary API（辞書統合）

#### GET /api/v1/dictionary/words/:word
**説明**: 単語の詳細情報を取得

**Response** (200 OK):
```json
{
  "word": "Здравствуйте",
  "language": "ru",
  "translation": "こんにちは",
  "phonetic": "/zdrɐˈstvʊjtʲɪ/",
  "part_of_speech": "interjection",
  "definitions": [
    {
      "definition": "Formal greeting",
      "example": "Здравствуйте, как дела?"
    }
  ],
  "frequency": "common",
  "related_words": ["привет", "здорово"]
}
```

#### POST /api/v1/dictionary/lookup/batch
**説明**: 複数の単語を一括検索

**Request**:
```json
{
  "words": ["Здравствуйте", "дела", "как"],
  "language": "ru",
  "translation_language": "ja"
}
```

---

### 2. Pattern API（会話パターン抽出）

#### GET /api/v1/patterns/books/:bookId/analyze
**説明**: 書籍から会話パターンを抽出

**Response** (200 OK):
```json
{
  "book_id": "550e8400",
  "patterns": [
    {
      "id": "pattern-1",
      "pattern": "Как [noun]?",
      "translation": "[名詞]はどう？",
      "frequency": 15,
      "examples": [
        {"text": "Как дела?", "translation": "元気？"},
        {"text": "Как жизнь?", "translation": "人生はどう？"}
      ],
      "difficulty": "beginner"
    }
  ],
  "total_patterns": 50
}
```

#### POST /api/v1/patterns/practice
**説明**: パターン練習セッションを開始

**Request**:
```json
{
  "pattern_id": "pattern-1",
  "practice_mode": "fill_in_blank"
}
```

---

### 3. Teacher Mode API（教師モード自動学習）

#### POST /api/v1/teacher-mode/books/:bookId/start
**説明**: 教師モードセッションを開始

**Request**:
```json
{
  "settings": {
    "speed": 1.0,
    "page_interval": 5,
    "repeat_count": 1,
    "include_translation": true,
    "include_explanation": true,
    "include_pronunciation_practice": false
  },
  "start_page": 1,
  "end_page": 150
}
```

**Response** (200 OK):
```json
{
  "session_id": "session-uuid",
  "playlist": [
    {
      "page_number": 1,
      "audio_url": "/storage/teacher-mode/session-uuid/page-1.mp3",
      "duration": 45
    }
  ],
  "total_duration": 6750,
  "estimated_completion": "2025-11-15T12:00:00Z"
}
```

#### GET /api/v1/teacher-mode/sessions/:sessionId/status
**説明**: セッション進捗を取得

**Response** (200 OK):
```json
{
  "session_id": "session-uuid",
  "status": "in_progress",
  "current_page": 12,
  "total_pages": 150,
  "elapsed_time": 540,
  "remaining_time": 6210
}
```

#### POST /api/v1/teacher-mode/books/:bookId/download
**説明**: オフライン用音声を一括ダウンロード準備

**Response** (200 OK):
```json
{
  "download_id": "download-uuid",
  "status": "preparing",
  "estimated_size": 250000000,
  "total_files": 150
}
```

---

### 4. Payment API（Stripe決済統合）

#### POST /api/v1/payment/create-checkout-session
**説明**: Stripe決済セッションを作成

**Request**:
```json
{
  "plan": "premium_monthly",
  "success_url": "https://hailango.com/payment/success",
  "cancel_url": "https://hailango.com/payment/cancel"
}
```

**Response** (200 OK):
```json
{
  "session_id": "cs_test_...",
  "url": "https://checkout.stripe.com/pay/cs_test_..."
}
```

#### GET /api/v1/payment/subscription
**説明**: 現在のサブスクリプション情報を取得

**Response** (200 OK):
```json
{
  "subscription_id": "sub_...",
  "plan": "premium_monthly",
  "status": "active",
  "current_period_start": "2025-11-01T00:00:00Z",
  "current_period_end": "2025-12-01T00:00:00Z",
  "cancel_at_period_end": false
}
```

#### POST /api/v1/payment/cancel
**説明**: サブスクリプションをキャンセル

---

### 5. WebSocket API（リアルタイム通知）

#### WebSocket接続エンドポイント
```
ws://localhost:8080/api/v1/ws
```

**接続時の認証**:
```json
{
  "type": "auth",
  "token": "JWT_TOKEN"
}
```

**サーバーからのメッセージ例**:
```json
{
  "type": "ocr_completed",
  "data": {
    "job_id": "ocr-job-uuid",
    "book_id": "550e8400",
    "page_number": 12,
    "status": "completed"
  }
}

{
  "type": "tts_generated",
  "data": {
    "page_id": "page-uuid",
    "audio_url": "/storage/..."
  }
}

{
  "type": "progress_update",
  "data": {
    "book_id": "550e8400",
    "completed_pages": 46
  }
}
```

---

## 🏗️ 実装優先順位

1. **Dictionary API** (4時間)
   - モック実装
   - Oxford/Wiktionary統合は後回し

2. **Pattern API** (6時間)
   - 基本的なパターンマッチング
   - AI抽出は後回し

3. **Teacher Mode API** (8時間)
   - プレイリスト生成
   - オフラインダウンロード

4. **Payment API** (6時間)
   - Stripe Test Mode
   - Webhook処理

5. **WebSocket API** (6時間)
   - 接続管理
   - イベント配信

---

## 📝 実装チェックリスト

### 各API共通
- [ ] ハンドラー作成
- [ ] モックリポジトリ実装
- [ ] ルーター登録
- [ ] テスト作成
- [ ] ドキュメント更新

### 動作確認
- [ ] すべてのエンドポイントが200/401を返す
- [ ] モックデータが正常に返される
- [ ] フロントエンドと統合テスト成功

---

## ✅ 完了条件

- [ ] 5つのAPIすべてが実装され、ルーターに登録されている
- [ ] モック実装が動作している
- [ ] すべてのテストがパスする
- [ ] フロントエンドからのリクエストが正常に処理される

---

**期限**: P0/P1完了後、2週間以内
**備考**: 実API統合は段階的に進める（モックファーストアプローチ）
