# HaiLanGo - Claude Code Project Memory

このファイルはClaude Codeがプロジェクトのコンテキストを理解するためのメモリです。

## プロジェクト概要

@README.md を参照してプロジェクトの全体像を把握してください。

### 重要なドキュメント
- 要件定義書: @docs/requirements_definition.md
- UI/UX設計書: @docs/ui_ux_design_document.md
- 教師モード技術仕様書: @docs/teacher_mode_technical_spec.md
- モック構築戦略: @docs/mocking_strategy.md
- API統合提案書: @docs/api_integration_proposal.md ⭐ NEW
- 機能実装RD: @docs/featureRDs/ （各機能の詳細仕様）

## 技術スタック

### バックエンド
- **言語**: Go
- **データベース**: PostgreSQL , Redis
- **コンテナ**: Podman（優先）または Docker

### フロントエンド
- **Web**: pnpm使用, Next.js 14+, TypeScript, React, TailwindCSS, ShadCN/UI
- **Mobile**: Flutter 3.0+
- **リンター/フォーマッター**: Biome.js (Web), `dart format` (Mobile)
- **テスト**: Vitest（単体・統合テスト）, Playwright（E2Eテスト）
- **CI/CD**: GitHub Actions

### インフラ
- **開発環境**: Podman Compose
- **本番環境**: AWS / GCP / Cloudflare（将来）
- **IaC**必須: Terraform

## コーディング規約

### Go（バックエンド）

#### 基本ルール
- `gofmt`でフォーマット（保存時に自動実行）
- `golangci-lint`を使用してリント
- エラーハンドリングは必ず行う（`if err != nil`を省略しない）
- コメントは英語で記述

#### ディレクトリ構造
```
backend/
├── cmd/                    # エントリーポイント
│   └── server/
│       └── main.go
├── internal/               # 内部パッケージ（他プロジェクトから使用不可）
│   ├── api/                # APIハンドラー
│   │   ├── handler/        # HTTPハンドラー
│   │   ├── middleware/     # ミドルウェア
│   │   └── router/         # ルーティング
│   ├── service/            # ビジネスロジック
│   ├── repository/         # データアクセス層
│   └── models/             # データモデル
└── pkg/                    # 外部パッケージ（再利用可能）
```

#### 命名規則
- パッケージ名: 小文字、単語区切りなし（例: `userservice`, `bookrepository`）
- インターフェース名: 動詞 + "er" または 名詞 + "Service" （例: `BookReader`, `UserService`）
- 構造体名: PascalCase（例: `UserProfile`, `BookMetadata`）
- メソッド/関数名: PascalCase（外部公開）、camelCase（内部のみ）

#### よく使うコマンド
```bash
# 依存関係のインストール
go mod download

# ビルド
go build -o bin/server cmd/server/main.go

# 実行
go run cmd/server/main.go

# テスト
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# リント
golangci-lint run

# フォーマット
gofmt -w .
```

### TypeScript / Next.js（Webフロントエンド）

#### 基本ルール
- Biome.jsでフォーマットとリント
- `use client`または`use server`を明示的に使用
- TypeScript strictモード有効
- コンポーネントは関数コンポーネントのみ（Classコンポーネント禁止）
- localhost以外でもVPN構成内からの接続もできるように`0.0.0.0`を許可すること

#### ディレクトリ構造
```
frontend/web/
├── app/                    # Next.js App Router
│   ├── (auth)/             # 認証グループ
│   ├── (main)/             # メインアプリ
│   └── api/                # API Routes
├── components/             # Reactコンポーネント
│   ├── ui/                 # ShadCN/UIコンポーネント
│   └── features/           # 機能別コンポーネント
├── lib/                    # ユーティリティ
├── hooks/                  # カスタムフック
├── types/                  # TypeScript型定義
└── public/                 # 静的ファイル
```

#### 命名規則
- コンポーネント: PascalCase（例: `BookCard.tsx`, `UserProfile.tsx`）
- フック: `use` + PascalCase（例: `useBookData.ts`, `useAuth.ts`）
- ユーティリティ: camelCase（例: `formatDate.ts`, `apiClient.ts`）
- 型定義: PascalCase（例: `User`, `BookMetadata`）

#### よく使うコマンド
```bash
# 依存関係のインストール
pnpm install

# 開発サーバー起動
pnpm run dev

# ビルド
pnpm run build

# 本番サーバー起動
pnpm start

# リント
pnpm run lint

# フォーマット
pnpm run format

# 型チェック
pnpm run type-check
```

### Flutter（モバイルアプリ）

#### 基本ルール
- `dart format`でフォーマット
- `flutter analyze`でリント
- Riverpodを状態管理に使用
- Material Design 3準拠

#### ディレクトリ構造
```
frontend/mobile/
├── lib/
│   ├── main.dart           # エントリーポイント
│   ├── app.dart            # アプリルート
│   ├── features/           # 機能別モジュール
│   │   ├── auth/
│   │   ├── book/
│   │   └── learning/
│   ├── core/               # 共通機能
│   │   ├── providers/      # Riverpodプロバイダー
│   │   ├── models/         # データモデル
│   │   └── services/       # サービス層
│   └── shared/             # 共有ウィジェット
└── test/                   # テスト
```

#### 命名規則
- ファイル名: snake_case（例: `book_card.dart`, `user_profile.dart`）
- クラス名: PascalCase（例: `BookCard`, `UserProfile`）
- 変数・関数: camelCase（例: `bookTitle`, `fetchUserData()`）
- 定数: lowerCamelCase（例: `apiBaseUrl`, `maxRetryCount`）

#### よく使うコマンド
```bash
# 依存関係のインストール
flutter pub get

# 実行
flutter run

# ビルド（Android）
flutter build apk

# ビルド（iOS）
flutter build ios

# テスト
flutter test

# リント
flutter analyze

# フォーマット
dart format .
```

## データベース

### PostgreSQL

#### 接続情報（開発環境）
```bash
Host: localhost
Port: 5432
Database: HaiLanGo_dev
User: HaiLanGo
Password: （.envファイル参照）
```

#### よく使うコマンド
```bash
# コンテナ内のPostgreSQLに接続
podman exec -it HaiLanGo-postgres psql -U HaiLanGo -d HaiLanGo_dev

# マイグレーション実行
go run cmd/migrate/main.go up

# マイグレーションロールバック
go run cmd/migrate/main.go down
```

### Redis

#### 接続情報（開発環境）
```bash
Host: localhost
Port: 6379
Password: （.envファイル参照）
```

#### よく使うコマンド
```bash
# Redisに接続
podman exec -it HaiLanGo-redis redis-cli

# キャッシュクリア
podman exec -it HaiLanGo-redis redis-cli FLUSHALL
```

## 環境変数

### 環境変数の設定

`.env.example` ファイルをコピーして `.env` を作成してください：

```bash
cp .env.example .env
```

### 重要な設定

#### APIキーなしでも開発可能

**APIキーがなくても開発・テストは可能です！**

```bash
# .envファイルに追加
USE_MOCK_APIS=true
```

この設定により、すべての外部API呼び出しが自動的にモックに置き換わります。
詳細は [モック構築戦略](docs/mocking_strategy.md) を参照してください。

#### 必須の環境変数（最小構成）

```.env
# アプリケーション環境
APP_ENV=development

# サーバーポート（ポート競合時は変更）
BACKEND_PORT=8080
FRONTEND_PORT=3000

# データベース
DATABASE_URL=postgresql://HaiLanGo:password@localhost:5432/HaiLanGo_dev
REDIS_URL=redis://localhost:6379

# JWT（開発用の簡単な値でも可）
JWT_SECRET=dev-secret-key-change-in-production

# フロントエンド
NEXT_PUBLIC_API_URL=http://localhost:8080

# モック使用（APIキーなしで開発する場合）
USE_MOCK_APIS=true
```

#### オプションの環境変数（実APIを使用する場合）

```.env
# Google Cloud APIs
GOOGLE_CLOUD_VISION_API_KEY=your_key_here
GOOGLE_CLOUD_TTS_API_KEY=your_key_here
GOOGLE_CLOUD_STT_API_KEY=your_key_here

# Azure Computer Vision
AZURE_COMPUTER_VISION_ENDPOINT=https://your-resource.cognitiveservices.azure.com/
AZURE_COMPUTER_VISION_API_KEY=your_key_here

# OpenAI
OPENAI_API_KEY=your_key_here

# Stripe
STRIPE_SECRET_KEY=sk_test_your_key_here
STRIPE_PUBLISHABLE_KEY=pk_test_your_key_here

# その他のAPIキー
# 詳細は .env.example を参照
```

**注意**: `.env` ファイルはGitにコミットしないでください（`.gitignore`に含まれています）

## 開発ワークフロー

### ブランチ戦略
- `main`: 本番環境
- `develop`: 開発環境
- `feature/*`: 新機能開発
- `bugfix/*`: バグ修正
- `hotfix/*`: 緊急修正

### コミットメッセージ
Conventional Commits形式を使用：

```
feat: 新機能追加
fix: バグ修正
docs: ドキュメント変更
style: コードフォーマット（動作に影響なし）
refactor: リファクタリング
test: テスト追加・修正
chore: ビルドプロセスやツール変更
```

例:
```bash
feat(auth): Google OAuth認証を追加
fix(ocr): ロシア語の認識精度を改善
docs: README.mdにクイックスタートを追加
```

### プルリクエスト
1. フィーチャーブランチを作成
2. 変更をコミット
3. `develop`ブランチに対してPR作成
4. レビュー後にマージ

## テスト

### テスト戦略

**TDD原則を徹底**: すべての機能はテストファーストで実装します。

**モックシステム**: APIキーなしでもテスト可能です。テスト実行時は自動的にモックが使用されます。

詳細は [モック構築戦略](docs/mocking_strategy.md) を参照してください。

### バックエンド（Go）
```bash
# すべてのテストを実行（自動的にモック使用）
go test ./...

# カバレッジ付き
go test -cover ./...

# 特定のパッケージのみ
go test ./internal/service/...

# ベンチマーク
go test -bench=. ./...

# 実APIを使用したテスト（APIキー必要）
USE_MOCK_APIS=false go test ./...
```

### フロントエンド（Next.js）
```bash
# Vitest単体・統合テスト
pnpm run test:unit

# E2Eテスト（Playwright）
pnpm run test:e2e

# すべてのテストを実行
pnpm test

# カバレッジ
pnpm run test:coverage

# BiomeJSリント・フォーマット
pnpm run lint
pnpm run format
```

### モバイル（Flutter）
```bash
# すべてのテストを実行
flutter test

# カバレッジ付き
flutter test --coverage

# 統合テスト
flutter test integration_test/
```

## デバッグ

### バックエンド
```bash
# Delveを使用
dlv debug cmd/server/main.go

# またはVS Codeのデバッガーを使用
```

### フロントエンド
- Chrome DevToolsを使用
- React DevToolsでコンポーネントツリーを確認

## よくある問題と解決方法

### Podmanコンテナが起動しない
```bash
# コンテナを停止して再起動
podman-compose down
podman-compose up -d

# ログを確認
podman-compose logs -f
```

### データベース接続エラー
```bash
# PostgreSQLコンテナが起動しているか確認
podman ps | grep postgres

# 接続テスト
podman exec -it HaiLanGo-postgres pg_isready
```

### ポートが既に使用されている
```bash
# 使用中のポートを確認
lsof -i :8080  # バックエンド
lsof -i :3000  # フロントエンド

# プロセスを終了
kill -9 <PID>
```

## セキュリティ

### 重要な注意事項
- **APIキーを絶対にコミットしない**（`.env`ファイルは`.gitignore`に含まれている）
- **ユーザーデータは必ずE2E暗号化**
- **SQLインジェクション対策**：プレースホルダーを使用
- **XSS対策**：入力値のサニタイズ

### 脆弱性スキャン
```bash
# Go依存関係のスキャン
go list -json -m all | nancy sleuth

# pnpmパッケージのスキャン
pnpm audit

# 自動修正
pnpm audit fix
```

## パフォーマンス最適化

### バックエンド
- データベースクエリはインデックスを活用
- N+1問題に注意
- Redisでキャッシュを活用

### フロントエンド
- 画像は必ずWebP形式に変換
- Next.jsの`next/image`を使用
- コンポーネントの遅延読み込み

## AI APIの使用

### モックシステム

**APIキーなしでも開発・テスト可能**: `USE_MOCK_APIS=true` を設定すると、すべての外部APIが自動的にモックに置き換わります。

```bash
# 開発時にモックを使用
USE_MOCK_APIS=true go run cmd/server/main.go

# テスト時にモックを使用（デフォルト）
go test ./...
```

詳細は [モック構築戦略](docs/mocking_strategy.md) を参照してください。

### コスト管理
- **キャッシュを最大限活用**（Redis + CDN）
- **レート制限を設定**（ユーザーごと、APIごと）
- **バッチ処理**で効率化
- **開発・テスト時はモックを使用**してコスト削減

### エラーハンドリング
```go
// リトライロジックの実装
func callAIAPI(ctx context.Context, request APIRequest) (*APIResponse, error) {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        resp, err := api.Call(ctx, request)
        if err == nil {
            return resp, nil
        }
        if i < maxRetries-1 {
            time.Sleep(time.Second * time.Duration(math.Pow(2, float64(i))))
        }
    }
    return nil, errors.New("max retries exceeded")
}
```

### APIクライアントの実装パターン

```go
// インターフェース定義
type OCRClient interface {
    ProcessImage(ctx context.Context, imageData []byte) (*OCRResult, error)
}

// ファクトリー関数（環境変数に基づいて切り替え）
func NewOCRClient() OCRClient {
    if os.Getenv("USE_MOCK_APIS") == "true" {
        return NewMockOCRClient()
    }
    return NewGoogleVisionClient(os.Getenv("GOOGLE_CLOUD_VISION_API_KEY"))
}
```

## プロジェクト固有の規約

### OCR処理
- OCR結果は必ずRedisにキャッシュ（TTL: 7日）
- ユーザーによる手動修正を許可
- 複数のOCR APIを試して最高精度を選択

### 音声データ
- TTS音声はCDNにキャッシュ
- オフライン用音声はZIP形式でパッケージング
- ダウンロードはWi-Fi接続時を推奨

### 教師モード
- バックグラウンド再生はMedia Session API使用（Web）
- audio_serviceパッケージ使用（Flutter）
- プレイリスト生成は非同期で行う

## 参考リンク

### 公式ドキュメント
- [Go公式](https://golang.org/doc/)
- [Next.js公式](https://nextjs.org/docs)
- [Flutter公式](https://flutter.dev/docs)
- [PostgreSQL公式](https://www.postgresql.org/docs/)
- [Redis公式](https://redis.io/documentation)

### AI APIs
- [Google Cloud Vision API](https://cloud.google.com/vision/docs)
- [Google Cloud TTS](https://cloud.google.com/text-to-speech/docs)
- [OpenAI API](https://platform.openai.com/docs)
- [Anthropic API](https://docs.anthropic.com)


---

このファイルはプロジェクトの進化に応じて更新してください。
質問や提案があれば、GitHub Issuesまたはプロジェクトチャットで共有してください。