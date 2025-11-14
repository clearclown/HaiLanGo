# CRITICAL_06: OCR/TTS/STT APIs実装（AI統合）

**優先度**: P0（OCR）/ P1（TTS/STT）
**担当者**: 未割当
**見積時間**: OCR 8-12時間、TTS 4-6時間、STT 6-8時間
**ブロッカー**: 書籍デジタル化（OCR）、音声機能（TTS/STT）が完全に欠落

---

## ⚠️ PM指示

**現状**: MVP核心機能のAI統合が全く進んでいない。
**期限**: OCR 96時間、TTS/STT 72時間以内に実装完了すること。
**重要**: 外部APIキーがなくてもモックで動作すること。

---

## 📋 OCR API エンドポイント

### 1. POST /api/v1/ocr/books/:bookId/pages/:pageNumber/process
**説明**: ページ画像をOCR処理

**Request**:
```http
POST /api/v1/ocr/books/550e8400/pages/12/process
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "language": "ru",
  "translation_language": "ja",
  "options": {
    "detect_layout": true,
    "extract_tables": true
  }
}
```

**Response** (200 OK):
```json
{
  "job_id": "ocr-job-uuid",
  "status": "processing",
  "estimated_time": 30
}
```

### 2. GET /api/v1/ocr/jobs/:jobId/status
**説明**: OCR処理ステータス確認

**Response** (200 OK):
```json
{
  "job_id": "ocr-job-uuid",
  "status": "completed",
  "result": {
    "text": "Здравствуйте! Как дела?",
    "translation": "こんにちは！元気ですか？",
    "confidence": 0.95,
    "language": "ru",
    "blocks": [
      {
        "text": "Здравствуйте!",
        "bbox": {"x": 10, "y": 20, "width": 100, "height": 30},
        "confidence": 0.96
      }
    ]
  }
}
```

### 3. PUT /api/v1/ocr/books/:bookId/pages/:pageNumber/edit
**説明**: OCR結果の手動修正

**Request**:
```json
{
  "text": "修正後のテキスト",
  "translation": "修正後の翻訳"
}
```

---

## 📋 TTS API エンドポイント

### 1. POST /api/v1/tts/synthesize
**説明**: テキストを音声に変換

**Request**:
```json
{
  "text": "Здравствуйте!",
  "language": "ru",
  "voice": "female",
  "speed": 1.0,
  "quality": "standard"
}
```

**Response** (200 OK):
```json
{
  "audio_url": "/storage/tts/audio-uuid.mp3",
  "duration": 2.5,
  "format": "mp3",
  "sample_rate": 22050
}
```

### 2. POST /api/v1/tts/books/:bookId/pages/:pageNumber/generate
**説明**: ページ全体の音声生成

**Response** (200 OK):
```json
{
  "page_audio_url": "/storage/books/550e8400/pages/12/audio.mp3",
  "phrases": [
    {
      "phrase_id": "phrase-1",
      "audio_url": "/storage/books/550e8400/pages/12/phrase-1.mp3"
    }
  ]
}
```

---

## 📋 STT API エンドポイント

### 1. POST /api/v1/stt/evaluate
**説明**: 発音評価

**Request** (multipart/form-data):
```
audio: <audio_file>
reference_text: "Здравствуйте!"
language: "ru"
```

**Response** (200 OK):
```json
{
  "transcription": "Здравствуйте",
  "reference": "Здравствуйте!",
  "score": 85,
  "feedback": {
    "accuracy": 88,
    "fluency": 82,
    "pronunciation": 85,
    "suggestions": [
      "Try to emphasize the 'вств' part more clearly"
    ]
  },
  "word_scores": [
    {
      "word": "Здравствуйте",
      "score": 85,
      "phonemes": [
        {"phoneme": "z", "score": 90},
        {"phoneme": "d", "score": 85}
      ]
    }
  ]
}
```

### 2. POST /api/v1/stt/transcribe
**説明**: 音声をテキストに変換（発音評価なし）

**Request** (multipart/form-data):
```
audio: <audio_file>
language: "ru"
```

**Response** (200 OK):
```json
{
  "transcription": "Здравствуйте! Как дела?",
  "language": "ru",
  "confidence": 0.92,
  "alternatives": [
    {"text": "Здравствуйте, как дела?", "confidence": 0.88}
  ]
}
```

---

## 🏗️ 実装アーキテクチャ

### モック実装の優先

**環境変数制御**:
```go
USE_MOCK_APIS=true  // 開発・テスト時
USE_MOCK_APIS=false // 本番環境
```

### ファクトリーパターン

```go
// pkg/ocr/factory.go
func NewOCRClient() OCRClient {
    if os.Getenv("USE_MOCK_APIS") == "true" {
        return NewMockOCRClient()
    }
    return NewGoogleVisionClient(os.Getenv("GOOGLE_CLOUD_VISION_API_KEY"))
}

// pkg/tts/factory.go
func NewTTSClient() TTSClient {
    if os.Getenv("USE_MOCK_APIS") == "true" {
        return NewMockTTSClient()
    }
    return NewGoogleTTSClient(os.Getenv("GOOGLE_CLOUD_TTS_API_KEY"))
}

// pkg/stt/factory.go
func NewSTTClient() STTClient {
    if os.Getenv("USE_MOCK_APIS") == "true" {
        return NewMockSTTClient()
    }
    return NewGoogleSTTClient(os.Getenv("GOOGLE_CLOUD_STT_API_KEY"))
}
```

### モック実装サンプル

```go
// pkg/ocr/mock.go
type MockOCRClient struct{}

func (m *MockOCRClient) ProcessImage(ctx context.Context, imageData []byte, language string) (*OCRResult, error) {
    // サンプルOCR結果を返す
    return &OCRResult{
        Text:       "Здравствуйте! Как дела?",
        Translation: "こんにちは！元気ですか？",
        Confidence: 0.95,
        Language:   language,
        Blocks: []TextBlock{
            {
                Text:       "Здравствуйте!",
                BBox:       BoundingBox{X: 10, Y: 20, Width: 100, Height: 30},
                Confidence: 0.96,
            },
        },
    }, nil
}

// pkg/tts/mock.go
type MockTTSClient struct{}

func (m *MockTTSClient) Synthesize(ctx context.Context, text, language string, options *TTSOptions) (*TTSResult, error) {
    // モック音声ファイルのURLを返す
    return &TTSResult{
        AudioURL:   "/mock/audio/" + uuid.New().String() + ".mp3",
        Duration:   2.5,
        Format:     "mp3",
        SampleRate: 22050,
    }, nil
}

// pkg/stt/mock.go
type MockSTTClient struct{}

func (m *MockSTTClient) Evaluate(ctx context.Context, audioData []byte, referenceText, language string) (*STTEvaluation, error) {
    // ランダムスコアを生成（80-95点）
    score := 80 + rand.Intn(16)

    return &STTEvaluation{
        Transcription: referenceText, // 参照テキストをそのまま返す
        Reference:     referenceText,
        Score:         score,
        Feedback: FeedbackDetail{
            Accuracy:      score - 3,
            Fluency:       score + 2,
            Pronunciation: score,
            Suggestions: []string{
                "Try to emphasize certain parts more clearly",
            },
        },
    }, nil
}
```

---

## 🗃️ データベーススキーマ

### ocr_jobs テーブル
```sql
CREATE TABLE ocr_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    book_id UUID NOT NULL REFERENCES books(id),
    page_number INT NOT NULL,
    status VARCHAR(50) NOT NULL, -- queued, processing, completed, failed
    language VARCHAR(10),
    translation_language VARCHAR(10),
    result JSONB,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);
```

### tts_cache テーブル
```sql
CREATE TABLE tts_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text_hash VARCHAR(64) UNIQUE NOT NULL,
    language VARCHAR(10) NOT NULL,
    voice VARCHAR(50),
    audio_url TEXT NOT NULL,
    duration FLOAT,
    created_at TIMESTAMP DEFAULT NOW(),
    last_accessed_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_tts_cache_hash ON tts_cache(text_hash);
```

### pronunciation_history テーブル
```sql
CREATE TABLE pronunciation_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    phrase_id UUID REFERENCES phrases(id),
    audio_url TEXT,
    transcription TEXT,
    reference_text TEXT,
    score INT,
    accuracy INT,
    fluency INT,
    pronunciation INT,
    feedback JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_pronunciation_history_user ON pronunciation_history(user_id);
CREATE INDEX idx_pronunciation_history_score ON pronunciation_history(user_id, score);
```

---

## 📝 実装チェックリスト

### OCR API
- [ ] `pkg/ocr/` パッケージ作成
- [ ] `internal/api/handler/ocr.go` 作成
- [ ] モック実装（MockOCRClient）
- [ ] Google Vision API実装（後回しOK）
- [ ] OCRジョブキュー実装
- [ ] ルーター登録

### TTS API
- [ ] `pkg/tts/` パッケージ作成
- [ ] `internal/api/handler/tts.go` 作成
- [ ] モック実装（MockTTSClient）
- [ ] Google TTS API実装（後回しOK）
- [ ] TTSキャッシュ実装
- [ ] ルーター登録

### STT API
- [ ] `pkg/stt/` パッケージ作成
- [ ] `internal/api/handler/stt.go` 作成
- [ ] モック実装（MockSTTClient）
- [ ] Google STT API実装（後回しOK）
- [ ] 発音評価ロジック
- [ ] ルーター登録

### テスト
- [ ] OCR APIテスト
- [ ] TTS APIテスト
- [ ] STT APIテスト

---

## ✅ 完了条件

- [ ] すべてのエンドポイントが実装され、ルーターに登録されている
- [ ] モック実装が動作している（`USE_MOCK_APIS=true`）
- [ ] すべてのテストがパスする
- [ ] フロントエンドからのリクエストが200/401を返す

---

**期限**: OCR 96時間、TTS/STT 72時間以内
**次のタスク**: CRITICAL_07（その他のAPI）
