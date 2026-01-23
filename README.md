<div align="center">

# HaiLanGo

### AI言語学習プラットフォーム

**自分の教科書をAIで学習可能にする**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Next.js](https://img.shields.io/badge/Next.js-14+-000000?style=flat&logo=next.js)](https://nextjs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

</div>

## 現在の開発状況

> **正直な評価**: このプロジェクトは開発中です。UIとアーキテクチャは整っていますが、
> **実際に動作させるにはAPIキーとデータベースが必要**です。

### 動作状況サマリー

| 機能 | APIキーなし | APIキーあり | 備考 |
|------|:-----------:|:-----------:|------|
| ユーザー認証 | ✅ 動作 | ✅ 動作 | JWT認証、完全実装 |
| ファイルアップロード | ✅ 動作 | ✅ 動作 | PDF/画像をディスクに保存 |
| OCR（文字認識） | ⚠️ モック | ✅ 動作 | Google Vision / Azure 必須 |
| TTS（音声合成） | ⚠️ モック | ✅ 動作 | Azure / Google / Edge-TTS(無料) |
| STT（音声認識） | ⚠️ モック | ✅ 動作 | Whisper / Azure 必須 |
| 単語帳・SRS | ✅ 動作 | ✅ 動作 | InMemoryでも動作 |
| 学習統計 | ✅ 動作 | ✅ 動作 | InMemoryでも動作 |
| データ永続化 | ❌ 再起動で消失 | ✅ 永続化 | PostgreSQL必須 |
| AI会話練習 | ❌ スタブ | ⚠️ 未完成 | OpenAI Realtime API（実装中） |

**結論**: 実際にPDF/写真から学習するには、**Google Vision APIキー**と**PostgreSQL**が最低限必要です。

---

## コンセプト

### 解決したい課題

1. **自分の本で勉強したい** - 既存アプリは強制カリキュラム
2. **スピーキング重視** - 発音練習・会話練習
3. **運転中に勉強したい** - ハンズフリー学習
4. **マイナー言語対応** - ペルシャ語↔日本語など
5. **AIと1対1の会話** - 本に沿った個別指導

### 対応言語

TTS/STTが動作する場合：
- **主要言語**: 日本語、英語、中国語、ロシア語、スペイン語、フランス語、ドイツ語
- **追加言語**: ペルシャ語、ヘブライ語、トルコ語、ポルトガル語、イタリア語
- **Whisper STT**: 99言語対応
- **Azure TTS**: 140+言語対応

---

## クイックスタート

### 必要なもの

```bash
# 必須
- Go 1.21+
- Node.js 18+
- pnpm 10+

# 実際に学習するには必要
- PostgreSQL 15+ （データ永続化）
- Google Cloud Vision APIキー （OCR）
- Azure Speech / OpenAI APIキー （TTS/STT）
```

### 1. モック版（UIの確認のみ）

```bash
# クローン
git clone https://github.com/clearclown/HaiLanGo.git
cd HaiLanGo

# バックエンド起動（モックモード）
cd backend
go mod download
USE_MOCK_APIS=true go run cmd/server/main.go

# フロントエンド起動（別ターミナル）
cd frontend/web
pnpm install
pnpm dev

# http://localhost:3000 でアクセス
```

**注意**: モックモードでは:
- OCRはダミーテキストを返す
- TTS/STTは動作しない
- データは再起動で消失

### 2. 実際に動作させる（APIキー必須）

```bash
# .envファイルを作成
cp .env.example .env
```

**.env** に以下を設定:

```bash
# データベース（必須）
DATABASE_URL=postgresql://user:password@localhost:5432/hailango

# OCR（必須 - どちらか一つ）
GOOGLE_APPLICATION_CREDENTIALS=./credentials.json
# または
GOOGLE_CLOUD_VISION_API_KEY=your_key

# TTS（推奨）
AZURE_SPEECH_KEY=your_key
AZURE_SPEECH_REGION=japaneast
# または Edge-TTS（無料だが品質は劣る）

# STT（推奨）
OPENAI_API_KEY=your_key  # Whisper用

# モックを無効化
USE_MOCK_APIS=false
```

```bash
# データベース起動
docker compose up -d postgres

# マイグレーション実行
cd backend
go run cmd/migrate/main.go up

# サーバー起動
go run cmd/server/main.go
```

---

## アーキテクチャ

```
┌─────────────────┐     ┌─────────────────┐
│   Next.js 14    │────▶│   Go Backend    │
│   (Frontend)    │     │   (API Server)  │
└─────────────────┘     └────────┬────────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        │                        │                        │
        ▼                        ▼                        ▼
┌───────────────┐      ┌─────────────────┐      ┌─────────────────┐
│  PostgreSQL   │      │   外部API       │      │  ファイル保存   │
│  (データ)     │      │  OCR/TTS/STT    │      │  ./storage      │
└───────────────┘      └─────────────────┘      └─────────────────┘
```

### ディレクトリ構造

```
HaiLanGo/
├── backend/
│   ├── cmd/server/          # エントリーポイント
│   ├── internal/
│   │   ├── api/handler/     # HTTPハンドラー
│   │   ├── service/         # ビジネスロジック
│   │   └── repository/      # データアクセス（InMemory + PostgreSQL）
│   └── pkg/
│       ├── ocr/             # OCRクライアント
│       ├── tts/             # TTSクライアント
│       ├── stt/             # STTクライアント
│       └── realtime/        # OpenAI Realtime API（実装中）
├── frontend/web/
│   ├── app/                 # Next.js App Router
│   ├── components/          # UIコンポーネント
│   └── lib/api/             # APIクライアント
└── docs/                    # ドキュメント
```

---

## 実装済み機能

### 完全に動作 ✅

| 機能 | 説明 |
|------|------|
| JWT認証 | 登録、ログイン、リフレッシュ、ログアウト |
| ファイルアップロード | PDF/画像のアップロード、チャンク対応 |
| 単語帳 | CRUD、タグ、CSV出力 |
| SRS（間隔反復） | SM-2アルゴリズム |
| 学習統計 | ダッシュボード、グラフ表示 |
| InMemoryフォールバック | DB不要で開発可能 |

### APIキー必須 ⚠️

| 機能 | 必要なAPI | 説明 |
|------|-----------|------|
| OCR | Google Vision / Azure | PDF/画像からテキスト抽出 |
| TTS | Azure / Google / Edge-TTS | テキスト読み上げ |
| STT | Whisper / Azure | 発音評価 |
| 発音スコアリング | Claude / GPT | LLMによる評価 |

### 未完成 / 実装中 ❌

| 機能 | 状態 |
|------|------|
| AI音声会話（Realtime API） | スタブ実装、動作せず |
| オフラインダウンロード | UIのみ、バックエンド未実装 |
| モバイルアプリ（Flutter） | 未着手 |

---

## APIキーの取得方法

### Google Cloud Vision（OCR用）

1. [Google Cloud Console](https://console.cloud.google.com/) でプロジェクト作成
2. Cloud Vision API を有効化
3. サービスアカウントを作成してJSONキーをダウンロード
4. `GOOGLE_APPLICATION_CREDENTIALS=./credentials.json` を設定

**料金**: 1000リクエスト/月まで無料、以降 $1.50/1000リクエスト

### Azure Speech Services（TTS/STT用）

1. [Azure Portal](https://portal.azure.com/) でSpeechリソース作成
2. キーとリージョンを取得
3. `AZURE_SPEECH_KEY` と `AZURE_SPEECH_REGION` を設定

**料金**: 500,000文字/月まで無料

### OpenAI（Whisper STT / LLM用）

1. [OpenAI Platform](https://platform.openai.com/) でAPIキー作成
2. `OPENAI_API_KEY` を設定

**料金**: Whisper $0.006/分、GPT-4 $0.03/1K tokens

### 無料代替（品質は劣る）

- **Edge-TTS**: 無料、インストール不要、品質は中程度
- **Tesseract OCR**: 無料、精度は低め

---

## 開発

### テスト実行

```bash
# バックエンド
cd backend
go test ./...

# フロントエンド
cd frontend/web
pnpm test
pnpm run type-check
```

### コード規約

- **Go**: `gofmt`、`golangci-lint`
- **TypeScript**: Biome.js

### コミットメッセージ

```
feat: 新機能
fix: バグ修正
docs: ドキュメント
refactor: リファクタリング
test: テスト
```

---

## ロードマップ

### 完了済み
- [x] 認証システム
- [x] ファイルアップロード
- [x] OCR/TTS/STTクライアント（モック + 実装）
- [x] 単語帳・SRS
- [x] 学習統計UI
- [x] InMemoryフォールバック

### 進行中
- [ ] OpenAI Realtime API会話機能（スタブ→実装）
- [ ] E2Eテストの充実
- [ ] 実API統合のテスト

### 将来
- [ ] オフラインダウンロード
- [ ] モバイルアプリ（Flutter）
- [ ] ユーザー生成コンテンツ

---

## 技術スタック

| レイヤー | 技術 |
|---------|------|
| フロントエンド | Next.js 14, TypeScript, TailwindCSS, ShadCN/UI |
| バックエンド | Go, Gin, GORM |
| データベース | PostgreSQL, Redis（オプション） |
| 外部API | Google Vision, Azure Speech, OpenAI Whisper |
| インフラ | Docker/Podman |

---

## 正直な評価

このプロジェクトは:

**良い点**:
- アーキテクチャは堅実（ファクトリーパターン、InMemoryフォールバック）
- UIは整っている
- 拡張性を考慮した設計

**課題**:
- 実際に動作させるにはAPIキーが必須
- 一部機能はスタブのまま
- E2Eテストが不十分
- ドキュメントと実装の乖離があった

**個人プロジェクトとして**:
- 自分用に動作させるなら十分
- 商用レベルにはまだ遠い

---

## ライセンス

MIT License - 詳細は [LICENSE](LICENSE) を参照

---

## 関連ドキュメント

- [要件定義書](docs/requirements_definition.md)
- [UI/UX設計書](docs/ui_ux_design_document.md)
- [実装状況](docs/IMPLEMENTATION_STATUS.md)
- [API統合提案](docs/api_integration_proposal.md)
