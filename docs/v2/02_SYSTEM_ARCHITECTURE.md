# システムアーキテクチャ

## 概要

HaiLanGoは、ユーザーの教材をAIで強化する言語学習プラットフォームである。
本ドキュメントでは、システム全体のアーキテクチャを定義する。

---

## アーキテクチャ概要図

```
┌─────────────────────────────────────────────────────────────────────┐
│                           クライアント層                              │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐    ┌─────────────────┐                         │
│  │   Web (Next.js)  │    │ Mobile (Flutter) │                        │
│  │   - SSR/SSG      │    │   - iOS/Android  │                        │
│  │   - React        │    │   - カメラ対応    │                        │
│  └────────┬────────┘    └────────┬────────┘                         │
└───────────┼──────────────────────┼──────────────────────────────────┘
            │                      │
            ▼                      ▼
┌─────────────────────────────────────────────────────────────────────┐
│                           API Gateway                                │
│                    (将来: Kong / AWS API Gateway)                    │
└─────────────────────────────────────────────────────────────────────┘
            │
            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       バックエンド層 (Go)                             │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  │
│  │   Auth      │  │   Books     │  │  Learning   │  │   Upload   │  │
│  │   Service   │  │   Service   │  │   Service   │  │   Service  │  │
│  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  │
│  │   OCR       │  │   TTS       │  │   STT       │  │   AI Chat  │  │
│  │   Service   │  │   Service   │  │   Service   │  │   Service  │  │
│  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │
│  │   SRS       │  │   Stats     │  │  Teacher    │                  │
│  │   Service   │  │   Service   │  │   Mode Svc  │                  │
│  └─────────────┘  └─────────────┘  └─────────────┘                  │
└─────────────────────────────────────────────────────────────────────┘
            │
            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         データ層                                      │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │
│  │ PostgreSQL  │  │    Redis    │  │  File Store │                  │
│  │ (メインDB)   │  │  (キャッシュ) │  │  (S3/R2)   │                  │
│  └─────────────┘  └─────────────┘  └─────────────┘                  │
└─────────────────────────────────────────────────────────────────────┘
            │
            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       外部API層                                       │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  │
│  │ Google      │  │ OpenAI      │  │ Stripe      │  │ DeepL      │  │
│  │ Cloud APIs  │  │ APIs        │  │ API         │  │ API        │  │
│  │ (Vision,TTS)│  │ (GPT,Whisper)│ │ (決済)      │  │ (翻訳)     │  │
│  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
```

---

## レイヤー詳細

### 1. クライアント層

#### Web (Next.js 14+)
```
frontend/web/
├── app/                    # App Router
│   ├── (auth)/             # 認証関連ページ
│   │   ├── login/
│   │   └── register/
│   ├── (main)/             # メインアプリ
│   │   ├── page.tsx        # ホーム（ダッシュボード）
│   │   ├── books/          # 教材一覧・詳細
│   │   ├── learn/          # 学習画面
│   │   ├── review/         # 復習画面
│   │   └── settings/       # 設定
│   └── api/                # API Routes (BFF)
├── components/
│   ├── ui/                 # ShadCN/UI
│   ├── learning/           # 学習関連
│   ├── upload/             # アップロード関連
│   └── common/             # 共通コンポーネント
├── lib/
│   ├── api/                # APIクライアント
│   ├── hooks/              # カスタムフック
│   └── utils/              # ユーティリティ
└── types/                  # TypeScript型定義
```

#### Mobile (Flutter)
```
frontend/mobile/
├── lib/
│   ├── main.dart
│   ├── app.dart
│   ├── features/
│   │   ├── auth/
│   │   ├── books/
│   │   ├── learning/
│   │   └── settings/
│   ├── core/
│   │   ├── api/
│   │   ├── models/
│   │   └── providers/      # Riverpod
│   └── shared/
│       └── widgets/
└── test/
```

### 2. バックエンド層 (Go)

#### ディレクトリ構造
```
backend/
├── cmd/
│   └── server/
│       └── main.go         # エントリーポイント
├── internal/
│   ├── api/
│   │   ├── router/         # ルーティング
│   │   ├── handler/        # HTTPハンドラー
│   │   ├── middleware/     # ミドルウェア
│   │   └── websocket/      # WebSocket
│   ├── service/            # ビジネスロジック
│   │   ├── auth/
│   │   ├── books/
│   │   ├── learning/
│   │   ├── ocr/
│   │   ├── tts/
│   │   ├── stt/
│   │   ├── ai/             # AI会話サービス
│   │   ├── srs/
│   │   └── stats/
│   ├── repository/         # データアクセス
│   │   ├── postgres/
│   │   ├── redis/
│   │   └── interfaces.go
│   ├── models/             # ドメインモデル
│   └── external/           # 外部API クライアント
│       ├── google/
│       ├── openai/
│       ├── stripe/
│       └── deepl/
├── pkg/                    # 再利用可能なパッケージ
│   ├── logger/
│   ├── config/
│   └── errors/
├── migrations/             # DBマイグレーション
└── tests/                  # テスト
    ├── integration/
    └── e2e/
```

#### サービス一覧

| サービス | 責務 | 依存外部API |
|----------|------|-------------|
| AuthService | 認証・認可 | - |
| BooksService | 教材CRUD | - |
| UploadService | ファイルアップロード | S3/R2 |
| OCRService | 画像→テキスト変換 | Google Vision |
| TTSService | テキスト→音声変換 | Google TTS / ElevenLabs |
| STTService | 音声→テキスト + 評価 | OpenAI Whisper |
| AIChatService | 会話形式学習 | OpenAI GPT-4 |
| SRSService | 間隔反復学習 | - |
| StatsService | 学習統計 | - |
| TeacherModeService | 自動再生モード | TTS依存 |

### 3. データ層

#### PostgreSQL
- ユーザー情報
- 教材メタデータ
- 学習進捗
- SRS データ
- 統計データ

#### Redis
- セッション管理
- OCR結果キャッシュ
- TTS音声キャッシュ
- レート制限

#### File Storage (S3/R2)
- アップロードされた教材（PDF/画像）
- 生成された音声ファイル
- オフライン用パッケージ

---

## 通信パターン

### 1. 同期通信 (REST API)

```
Client → Backend → External API → Backend → Client
```

使用場面:
- ユーザー認証
- 教材CRUD
- 設定変更
- 統計取得

### 2. 非同期通信 (WebSocket)

```
Client ←→ WebSocket ←→ Backend
              ↓
         Background Job
              ↓
         External API
```

使用場面:
- OCR処理の進捗通知
- リアルタイム会話
- 教師モードの制御

### 3. ストリーミング (Server-Sent Events / Streaming)

```
Client ← SSE ← Backend ← Streaming ← LLM API
```

使用場面:
- AI会話のレスポンス
- 長いテキストの読み上げ

---

## セキュリティ

### 認証フロー

```
1. ユーザーがログイン
2. バックエンドがJWTを発行
3. クライアントがJWTを保存（httpOnly cookie）
4. 以降のリクエストにJWTを添付
5. バックエンドがJWTを検証
```

### JWT構造

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user_id",
    "email": "user@example.com",
    "iat": 1234567890,
    "exp": 1234567890
  }
}
```

### セキュリティ対策

| 脅威 | 対策 |
|------|------|
| XSS | httpOnly cookie, CSP |
| CSRF | SameSite cookie, CSRF token |
| SQLi | Prepared statements |
| 認証情報漏洩 | bcrypt, E2E暗号化 |
| レート制限 | Redis based rate limiting |

---

## スケーラビリティ

### 現在（個人開発）

```
Single Server
├── Go Backend
├── PostgreSQL
├── Redis
└── File Storage (Local)
```

### 将来（スケールアウト）

```
Load Balancer
├── Go Backend (複数インスタンス)
├── PostgreSQL (Primary + Read Replica)
├── Redis Cluster
└── S3/R2 (CDN経由)
```

---

## 監視・ログ

### ログ

```go
// 構造化ログ（JSON形式）
{
  "timestamp": "2024-01-01T00:00:00Z",
  "level": "INFO",
  "service": "ocr",
  "trace_id": "abc123",
  "message": "OCR processing completed",
  "duration_ms": 1234,
  "page_count": 10
}
```

### メトリクス

| メトリクス | 説明 |
|------------|------|
| request_count | リクエスト数 |
| request_duration | レスポンス時間 |
| error_rate | エラー率 |
| ocr_processing_time | OCR処理時間 |
| tts_generation_time | TTS生成時間 |
| active_users | アクティブユーザー数 |

### 監視ツール

- **ログ**: stdout → (将来: CloudWatch / Loki)
- **メトリクス**: Prometheus（将来）
- **エラー追跡**: Sentry
- **APM**: (将来: Datadog / New Relic)

---

## 環境構成

### 開発環境

```yaml
# docker-compose.yml
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://...
      - REDIS_URL=redis://...
    depends_on:
      - postgres
      - redis

  frontend:
    build: ./frontend/web
    ports:
      - "3000:3000"

  postgres:
    image: postgres:15
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7
    volumes:
      - redis_data:/data
```

### 本番環境（将来）

```
Cloudflare (CDN + WAF)
       ↓
AWS/GCP Load Balancer
       ↓
ECS/Cloud Run (Go Backend)
       ↓
RDS/Cloud SQL (PostgreSQL)
ElastiCache/Memorystore (Redis)
S3/R2 (File Storage)
```

---

## 次のドキュメント

- [03_DATABASE_SCHEMA.md](./03_DATABASE_SCHEMA.md) - データベース設計
