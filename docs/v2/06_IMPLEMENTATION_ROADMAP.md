# 実装ロードマップ

## 概要

HaiLanGoの実装を段階的に進めるためのロードマップ。
各フェーズで明確な成果物と検証基準を定義する。

---

## 実装順序の原則

1. **データ層から上へ**: Repository → Service → Handler
2. **テストファースト**: 実装前にテストを書く
3. **依存関係を尊重**: 依存先を先に実装
4. **動作確認を頻繁に**: 小さな単位で動作確認

---

## フェーズ0: 基盤整備

### 目標
- Testcontainersによるテスト環境
- DBマイグレーション基盤
- CI/CD基盤

### タスク

| # | タスク | 成果物 | 検証方法 |
|---|--------|--------|----------|
| 0.1 | Testcontainersセットアップ | `tests/testhelpers/` | `go test ./tests/testhelpers/... -v` |
| 0.2 | マイグレーション基盤 | `migrations/`, `cmd/migrate/` | マイグレーション実行成功 |
| 0.3 | GitHub Actions設定 | `.github/workflows/` | CI緑 |
| 0.4 | 開発環境Docker Compose | `docker-compose.yml` | `docker-compose up` 成功 |

### 完了基準
- [ ] `go test ./...` が成功
- [ ] マイグレーションが実行できる
- [ ] CIが緑になる

---

## フェーズ1: 認証基盤

### 目標
- ユーザー登録・ログイン
- JWT認証

### 依存関係
```
users テーブル
    ↓
UserRepository
    ↓
AuthService
    ↓
AuthHandler
    ↓
AuthMiddleware
```

### タスク

| # | タスク | 成果物 | 検証方法 |
|---|--------|--------|----------|
| 1.1 | users マイグレーション | `migrations/000001_*.sql` | テーブル作成確認 |
| 1.2 | User モデル | `internal/models/user.go` | - |
| 1.3 | UserRepository | `internal/repository/postgres/user.go` | 統合テスト |
| 1.4 | AuthService | `internal/service/auth/` | 単体テスト |
| 1.5 | AuthHandler | `internal/api/handler/auth.go` | HTTPテスト |
| 1.6 | AuthMiddleware | `internal/api/middleware/auth.go` | 統合テスト |
| 1.7 | Router統合 | `internal/api/router/router.go` | E2Eテスト |

### 完了基準
- [ ] POST /api/v1/auth/register でユーザー作成
- [ ] POST /api/v1/auth/login でJWT取得
- [ ] GET /api/v1/auth/me で認証情報取得
- [ ] 認証なしリクエストが401
- [ ] 統合テストがパス

---

## フェーズ2: 教材管理

### 目標
- 教材CRUD
- ファイルアップロード（チャンク対応）

### 依存関係
```
フェーズ1完了
    ↓
books, pages テーブル
    ↓
BookRepository, PageRepository
    ↓
BookService, UploadService
    ↓
BookHandler, UploadHandler
```

### タスク

| # | タスク | 成果物 | 検証方法 |
|---|--------|--------|----------|
| 2.1 | books マイグレーション | `migrations/000002_*.sql` | - |
| 2.2 | pages マイグレーション | `migrations/000003_*.sql` | - |
| 2.3 | Book, Page モデル | `internal/models/` | - |
| 2.4 | BookRepository | `internal/repository/postgres/book.go` | 統合テスト |
| 2.5 | PageRepository | `internal/repository/postgres/page.go` | 統合テスト |
| 2.6 | BookService | `internal/service/book/` | 単体テスト |
| 2.7 | UploadService（チャンク） | `internal/service/upload/` | 統合テスト |
| 2.8 | BookHandler | `internal/api/handler/book.go` | HTTPテスト |
| 2.9 | UploadHandler | `internal/api/handler/upload.go` | 統合テスト |

### 完了基準
- [ ] POST /api/v1/books で教材作成
- [ ] GET /api/v1/books で一覧取得
- [ ] チャンクアップロードでPDFアップロード
- [ ] アップロード後にpagesレコードが作成される
- [ ] 統合テストがパス

---

## フェーズ3: OCR処理

### 目標
- OCR処理（Google Vision API）
- WebSocket進捗通知
- OCR結果の手動修正

### 依存関係
```
フェーズ2完了
    ↓
Google Vision APIクライアント
    ↓
OCRService
    ↓
OCRHandler
    ↓
WebSocket通知
```

### タスク

| # | タスク | 成果物 | 検証方法 |
|---|--------|--------|----------|
| 3.1 | Vision APIクライアント | `internal/external/google/vision.go` | モック対応 |
| 3.2 | OCRService | `internal/service/ocr/` | 単体テスト |
| 3.3 | OCRHandler | `internal/api/handler/ocr.go` | HTTPテスト |
| 3.4 | WebSocketハンドラー | `internal/api/websocket/` | 統合テスト |
| 3.5 | OCR結果修正API | `internal/api/handler/page.go` | HTTPテスト |
| 3.6 | バックグラウンドジョブ | `internal/service/ocr/worker.go` | 統合テスト |

### 完了基準
- [ ] 画像アップロード後にOCR処理が開始
- [ ] WebSocketで進捗通知
- [ ] OCR完了後にocr_textが保存
- [ ] PUT /api/v1/books/:id/pages/:num/ocr で修正可能
- [ ] モックAPIでテスト可能

---

## フェーズ4: TTS/STT

### 目標
- 音声合成（TTS）
- 発音認識・評価（STT）

### 依存関係
```
フェーズ3完了
    ↓
TTS APIクライアント（Google/ElevenLabs）
STT APIクライアント（OpenAI Whisper）
    ↓
TTSService, STTService
    ↓
TTSHandler, STTHandler
```

### タスク

| # | タスク | 成果物 | 検証方法 |
|---|--------|--------|----------|
| 4.1 | Google TTS クライアント | `internal/external/google/tts.go` | モック対応 |
| 4.2 | Whisper クライアント | `internal/external/openai/whisper.go` | モック対応 |
| 4.3 | TTSService | `internal/service/tts/` | 単体テスト |
| 4.4 | STTService（評価ロジック含む） | `internal/service/stt/` | 単体テスト |
| 4.5 | TTSHandler | `internal/api/handler/tts.go` | HTTPテスト |
| 4.6 | STTHandler | `internal/api/handler/stt.go` | 統合テスト |
| 4.7 | 音声キャッシュ（Redis） | `internal/service/cache/` | 統合テスト |

### 完了基準
- [ ] POST /api/v1/tts/synthesize で音声生成
- [ ] 同じテキストはキャッシュから返す
- [ ] POST /api/v1/stt/evaluate で発音評価
- [ ] スコア（0-100）と詳細フィードバックを返す
- [ ] モックAPIでテスト可能

---

## フェーズ5: AI会話

### 目標
- 会話形式の学習
- ストリーミングレスポンス

### 依存関係
```
フェーズ4完了
    ↓
OpenAI GPT クライアント
    ↓
AIChatService
    ↓
AIChatHandler（SSE対応）
```

### タスク

| # | タスク | 成果物 | 検証方法 |
|---|--------|--------|----------|
| 5.1 | OpenAI クライアント | `internal/external/openai/chat.go` | モック対応 |
| 5.2 | 会話履歴テーブル | `migrations/000009_*.sql` | - |
| 5.3 | AIChatService | `internal/service/ai/chat.go` | 単体テスト |
| 5.4 | プロンプトテンプレート | `internal/service/ai/prompts/` | - |
| 5.5 | AIChatHandler（SSE） | `internal/api/handler/ai_chat.go` | 統合テスト |

### 完了基準
- [ ] POST /api/v1/ai/chat でAI会話
- [ ] SSEでストリーミングレスポンス
- [ ] 会話履歴が保存される
- [ ] ページコンテキストを考慮した会話

---

## フェーズ6: SRS（間隔反復学習）

### 目標
- SM-2アルゴリズム実装
- 復習対象取得
- 復習結果記録

### 依存関係
```
フェーズ5完了
    ↓
vocabularies, srs_items テーブル
    ↓
VocabularyRepository, SRSRepository
    ↓
SRSService
    ↓
ReviewHandler
```

### タスク

| # | タスク | 成果物 | 検証方法 |
|---|--------|--------|----------|
| 6.1 | vocabularies マイグレーション | `migrations/000004_*.sql` | - |
| 6.2 | srs_items マイグレーション | `migrations/000005_*.sql` | - |
| 6.3 | Vocabulary, SRSItem モデル | `internal/models/` | - |
| 6.4 | VocabularyRepository | `internal/repository/postgres/vocabulary.go` | 統合テスト |
| 6.5 | SRSRepository | `internal/repository/postgres/srs.go` | 統合テスト |
| 6.6 | SM-2アルゴリズム | `internal/service/srs/sm2.go` | 単体テスト |
| 6.7 | SRSService | `internal/service/srs/` | 単体テスト |
| 6.8 | ReviewHandler | `internal/api/handler/review.go` | HTTPテスト |

### 完了基準
- [ ] GET /api/v1/review/due で復習対象取得
- [ ] POST /api/v1/review/submit で結果送信
- [ ] SM-2アルゴリズムで次回復習日計算
- [ ] 優先度（urgent/recommended/optional）が正しく計算

---

## フェーズ7: 学習統計

### 目標
- 学習セッション記録
- 統計集計
- ダッシュボードAPI

### 依存関係
```
フェーズ6完了
    ↓
learning_sessions, learning_logs, user_stats テーブル
    ↓
StatsRepository
    ↓
StatsService
    ↓
StatsHandler
```

### タスク

| # | タスク | 成果物 | 検証方法 |
|---|--------|--------|----------|
| 7.1 | learning_sessions マイグレーション | `migrations/000006_*.sql` | - |
| 7.2 | learning_logs マイグレーション | `migrations/000007_*.sql` | - |
| 7.3 | user_stats マイグレーション | `migrations/000008_*.sql` | - |
| 7.4 | StatsRepository | `internal/repository/postgres/stats.go` | 統合テスト |
| 7.5 | StatsService | `internal/service/stats/` | 単体テスト |
| 7.6 | StatsHandler | `internal/api/handler/stats.go` | HTTPテスト |
| 7.7 | 統計集計バッチ | `internal/service/stats/aggregator.go` | 統合テスト |

### 完了基準
- [ ] GET /api/v1/stats/overview で概要取得
- [ ] GET /api/v1/stats/history で履歴取得
- [ ] 連続学習日数が正しく計算
- [ ] 学習時間が正しく集計

---

## フェーズ8: フロントエンド

### 目標
- Next.js フロントエンド実装
- 全ユースケースの実現

### タスク

| # | タスク | 成果物 | 検証方法 |
|---|--------|--------|----------|
| 8.1 | 認証ページ | `app/(auth)/` | E2Eテスト |
| 8.2 | ホーム画面 | `app/(main)/page.tsx` | E2Eテスト |
| 8.3 | 教材一覧・詳細 | `app/(main)/books/` | E2Eテスト |
| 8.4 | アップロード画面 | `app/(main)/upload/` | E2Eテスト |
| 8.5 | 学習画面 | `app/(main)/learn/` | E2Eテスト |
| 8.6 | 復習画面 | `app/(main)/review/` | E2Eテスト |
| 8.7 | 統計画面 | `app/(main)/stats/` | E2Eテスト |
| 8.8 | 設定画面 | `app/(main)/settings/` | E2Eテスト |

### 完了基準
- [ ] 全ページが動作
- [ ] E2Eテストがパス
- [ ] レスポンシブ対応

---

## 全体タイムライン（目安）

```
フェーズ0: 基盤整備      [████] 1週間
フェーズ1: 認証基盤      [████] 1週間
フェーズ2: 教材管理      [████████] 2週間
フェーズ3: OCR処理       [████████] 2週間
フェーズ4: TTS/STT       [████████] 2週間
フェーズ5: AI会話        [████████] 2週間
フェーズ6: SRS           [████] 1週間
フェーズ7: 学習統計      [████] 1週間
フェーズ8: フロントエンド [████████████] 3週間
                        ─────────────────
                        合計: 約15週間
```

**注意**: この見積もりは目安であり、実際の進捗に応じて調整する。

---

## チェックリスト

### 各フェーズ開始前
- [ ] 前フェーズの完了基準をすべて満たしている
- [ ] テストがすべてパスしている
- [ ] コードレビュー完了（該当する場合）

### 各フェーズ完了時
- [ ] 完了基準をすべて満たしている
- [ ] 統合テストがパス
- [ ] ドキュメントを更新
- [ ] CIが緑

---

## 次のアクション

1. フェーズ0のタスク0.1から開始
2. Testcontainersのセットアップを完了
3. マイグレーション基盤を構築
4. CIを設定
