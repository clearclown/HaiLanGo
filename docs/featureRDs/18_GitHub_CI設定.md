# 機能実装: GitHub CI設定

## 要件

### 機能要件

#### CI/CDパイプライン
1. **テスト実行**
   - バックエンド（Go）テスト
   - フロントエンド（Vitest）単体・統合テスト
   - フロントエンド（Playwright）E2Eテスト

2. **リント・フォーマット**
   - Go: golangci-lint
   - フロントエンド: BiomeJS

3. **ビルド**
   - バックエンドビルド
   - フロントエンドビルド

4. **デプロイ**
   - テスト成功時の自動デプロイ（オプション）

### 非機能要件

- **パフォーマンス**:
  - CI実行時間: 10分以内
  - 並列実行で高速化

- **セキュリティ**:
  - シークレット管理
  - 依存関係の脆弱性スキャン

## 実装指示

### Step 1: テスト設計

CI設定ファイルを作成：

1. **バックエンドCI**
   - Goテスト実行
   - golangci-lint実行
   - ビルド確認

2. **フロントエンドCI**
   - Vitest実行
   - Playwright実行
   - BiomeJS実行
   - ビルド確認

### Step 2: 実装

#### 実装ファイル構造

```
.github/
└── workflows/
    ├── backend.yml                  # バックエンドCI
    ├── frontend.yml                 # フロントエンドCI
    └── e2e.yml                      # E2Eテスト（必要に応じて分離）
```

#### バックエンドCI設定

```yaml
# .github/workflows/backend.yml
name: Backend Tests

on:
  push:
    paths:
      - 'backend/**'
      - '.github/workflows/backend.yml'
  pull_request:
    paths:
      - 'backend/**'

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:7
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}

      - name: Install dependencies
        working-directory: ./backend
        run: go mod download

      - name: Run tests
        working-directory: ./backend
        run: go test -v -race -coverprofile=coverage.out ./...
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/test?sslmode=disable
          REDIS_URL: redis://localhost:6379

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          working-directory: ./backend

      - name: Build
        working-directory: ./backend
        run: go build ./cmd/server
```

#### フロントエンドCI設定

```yaml
# .github/workflows/frontend.yml
name: Frontend Tests

on:
  push:
    paths:
      - 'frontend/web/**'
      - '.github/workflows/frontend.yml'
  pull_request:
    paths:
      - 'frontend/web/**'

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'pnpm'

      - name: Install pnpm
        uses: pnpm/action-setup@v2
        with:
          version: 8

      - name: Install dependencies
        working-directory: ./frontend/web
        run: pnpm install

      - name: Run BiomeJS
        working-directory: ./frontend/web
        run: pnpm run lint

      - name: Run Vitest
        working-directory: ./frontend/web
        run: pnpm run test:unit

      - name: Install Playwright Browsers
        working-directory: ./frontend/web
        run: pnpm exec playwright install --with-deps

      - name: Run Playwright tests
        working-directory: ./frontend/web
        run: pnpm run test:e2e

      - name: Build
        working-directory: ./frontend/web
        run: pnpm run build
        env:
          NEXT_PUBLIC_API_URL: ${{ secrets.NEXT_PUBLIC_API_URL }}
```

## 制約事項

- `.github/workflows/` 配下のみ編集可能
- シークレットはGitHub Secretsで管理

## 完了条件

- [ ] バックエンドCIが通る
- [ ] フロントエンドCIが通る
- [ ] すべてのテストが実行される
- [ ] リント・フォーマットチェックが実行される
- [ ] ビルドが成功する
