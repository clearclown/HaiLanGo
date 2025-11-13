# HaiLanGo Backend

HaiLanGoのバックエンドサービス（Go実装）

## 概要

このディレクトリには、HaiLanGoプロジェクトのバックエンドAPIサービスが含まれています。

## 機能

### STT（Speech-to-Text）発音評価 ✅ 実装済み

音声認識と発音評価機能を提供します。

#### 主な機能

1. **音声認識（STT）**
   - Google Cloud STT対応
   - Whisper API対応（フォールバック）
   - 主要12言語対応
   - モックAPIによる開発・テスト対応

2. **発音評価**
   - 総合スコア計算（0-100点）
   - 正確性スコア（テキスト比較）
   - 流暢性スコア（速度・リズム分析）
   - 発音スコア（音素レベル評価）
   - 単語レベルの詳細評価

3. **フィードバック生成**
   - レベル判定（excellent, good, fair, poor）
   - 良かった点の提示
   - 改善ポイントの提示
   - 具体的なアドバイス

4. **音声処理**
   - ノイズ除去
   - 音量正規化
   - サンプリングレート変換（16kHz推奨）
   - 音声品質検証

## ディレクトリ構造

```
backend/
├── cmd/
│   └── server/              # サーバーエントリーポイント
├── internal/
│   ├── api/                 # APIハンドラー・ルーター
│   ├── models/              # データモデル
│   │   └── stt.go          # STT関連のモデル
│   └── service/
│       ├── stt/            # STTサービス ✅
│       │   ├── service.go          # メインサービス
│       │   ├── evaluation.go       # 発音評価ロジック
│       │   ├── constants.go        # 定数定義
│       │   └── service_test.go     # テスト
│       └── audio/          # 音声処理 ✅
│           ├── processor.go        # 音声処理ロジック
│           └── processor_test.go   # テスト
├── pkg/
│   └── stt/                # STTクライアント ✅
│       ├── interface.go            # インターフェース定義
│       ├── google_stt.go           # Google STTクライアント
│       ├── whisper.go              # Whisperクライアント
│       ├── mock_stt.go             # モッククライアント
│       └── stt_test.go             # テスト
├── go.mod
├── go.sum
└── README.md
```

## セットアップ

### 前提条件

- Go 1.21+
- PostgreSQL 15+（将来）
- Redis 7+（将来）

### 依存関係のインストール

```bash
go mod download
```

### 環境変数

`.env`ファイルを作成し、以下を設定：

```bash
# モックAPIを使用（APIキーなしで開発可能）
USE_MOCK_APIS=true

# 実際のAPIを使用する場合
# USE_MOCK_APIS=false
# GOOGLE_CLOUD_STT_API_KEY=your_key_here
# OPENAI_API_KEY=your_key_here
```

## テスト

### すべてのテストを実行

```bash
# モックを使用（デフォルト）
go test ./...

# カバレッジ付き
go test -cover ./...

# 詳細出力
go test -v ./...
```

### 特定のパッケージのみテスト

```bash
# STTサービスのみ
go test ./internal/service/stt/...

# STTクライアントのみ
go test ./pkg/stt/...

# 音声処理のみ
go test ./internal/service/audio/...
```

### 実APIを使用したテスト

```bash
# APIキーを設定
export GOOGLE_CLOUD_STT_API_KEY=your_key_here
export USE_MOCK_APIS=false

# テスト実行
go test ./...
```

## 開発

### コーディング規約

- `gofmt`でフォーマット（保存時に自動実行）
- `golangci-lint`を使用してリント
- エラーハンドリングは必ず行う
- コメントは日本語で記述
- テストファーストで開発（TDD）

### テストの作成

新しい機能を実装する際は、必ずテストを先に作成します（TDD）：

1. テストファイルを作成（`*_test.go`）
2. テストケースを記述
3. テストが失敗することを確認
4. 実装を行う
5. テストがパスすることを確認
6. リファクタリング

### リント

```bash
# インストール
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 実行
golangci-lint run
```

## API仕様

### STT発音評価API

#### 音声認識

```http
POST /api/v1/stt/recognize
Content-Type: application/json

{
  "audio_data": "base64_encoded_audio",
  "language": "en-US"
}
```

**レスポンス:**

```json
{
  "text": "Hello, world!",
  "language": "en-US",
  "confidence": 0.95,
  "duration": 1.5,
  "words": [
    {
      "word": "Hello",
      "start_time": 0.0,
      "end_time": 0.5,
      "confidence": 0.96
    }
  ]
}
```

#### 発音評価

```http
POST /api/v1/stt/evaluate-pronunciation
Content-Type: application/json

{
  "expected_text": "Hello, world!",
  "audio_data": "base64_encoded_audio",
  "language": "en-US"
}
```

**レスポンス:**

```json
{
  "total_score": 85,
  "accuracy_score": 90,
  "fluency_score": 80,
  "pronunc_score": 85,
  "word_scores": [
    {
      "word": "Hello",
      "score": 95,
      "is_correct": true,
      "expected_word": "Hello",
      "recognized_word": "Hello"
    }
  ],
  "feedback": {
    "level": "good",
    "message": "👍 良好です！もう少しで完璧です。",
    "positive_points": [
      "基本的な発音は正確です",
      "理解しやすい発音です"
    ],
    "improvements": [
      "いくつかの単語の発音を改善できます"
    ],
    "specific_advice": [
      "個々の音素をはっきりと発音しましょう"
    ]
  }
}
```

## パフォーマンス

### テストカバレッジ

- `internal/service/audio`: 71.8%
- `internal/service/stt`: 91.3%
- `pkg/stt`: 72.7%

### ベンチマーク

```bash
go test -bench=. ./...
```

## トラブルシューティング

### テストが失敗する

1. 依存関係を更新：`go mod tidy`
2. 環境変数を確認：`USE_MOCK_APIS=true`
3. キャッシュをクリア：`go clean -testcache`

### APIエラー

1. モックを使用：`USE_MOCK_APIS=true`
2. APIキーを確認
3. ログを確認：`LOG_LEVEL=debug`

## コントリビューション

1. フィーチャーブランチを作成
2. TDDでテストを先に作成
3. 実装
4. テストがすべてパスすることを確認
5. リントエラーがないことを確認
6. プルリクエストを作成

## ライセンス

MIT License

## 関連ドキュメント

- [要件定義書](../docs/requirements_definition.md)
- [UI/UX設計書](../docs/ui_ux_design_document.md)
- [モック構築戦略](../docs/mocking_strategy.md)
- [API統合提案書](../docs/api_integration_proposal.md)
- [機能実装RD: STT発音評価](../docs/featureRDs/5_STT発音評価.md)
