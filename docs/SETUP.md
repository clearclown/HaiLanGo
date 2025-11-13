# HaiLanGo - セットアップガイド

## 📋 目次

1. [前提条件](#前提条件)
2. [環境変数の設定](#環境変数の設定)
3. [Podman環境でのセットアップ](#podman環境でのセットアップ)
4. [Docker環境でのセットアップ](#docker環境でのセットアップ)
5. [開発モード](#開発モード)
6. [トラブルシューティング](#トラブルシューティング)
7. [よくある質問](#よくある質問)

---

## 前提条件

### 必須ソフトウェア

#### Podman環境（推奨）
- **Podman** 4.0+
- **podman-compose** 1.0+
- **Go** 1.21+
- **Node.js** 18+
- **pnpm** 8+

#### Docker環境（代替）
- **Docker** 20.10+
- **Docker Compose** 2.0+
- **Go** 1.21+
- **Node.js** 18+
- **pnpm** 8+

### インストール確認

```bash
# Podmanのバージョン確認
podman --version
podman-compose --version

# または Docker
docker --version
docker-compose --version

# Go
go version

# Node.js & pnpm
node --version
pnpm --version
```

---

## 環境変数の設定

### 1. 環境変数ファイルの作成

```bash
cd /path/to/HaiLanGo
cp .env.example .env
```

### 2. 基本設定（必須）

`.env` ファイルを編集します：

```bash
# アプリケーション環境
APP_ENV=development
APP_NAME=HaiLanGo
DEBUG=true

# サーバーポート
BACKEND_PORT=8080
FRONTEND_PORT=3000
SERVER_HOST=0.0.0.0

# データベース（デフォルトで問題なし）
POSTGRES_DB=HaiLanGo_dev
POSTGRES_USER=HaiLanGo
POSTGRES_PASSWORD=password
POSTGRES_PORT=5432

# Redis（デフォルトで問題なし）
REDIS_PORT=6379

# JWT（開発環境用）
JWT_SECRET=dev-secret-key-change-in-production

# フロントエンド
NEXT_PUBLIC_API_URL=http://localhost:8080

# モックAPI使用（APIキーなしで開発可能）
USE_MOCK_APIS=true
```

### 3. 本番API使用（オプション）

実際のAPIキーを使用する場合は、以下を設定します：

```bash
# モックを無効化
USE_MOCK_APIS=false

# Google Cloud APIs
GOOGLE_CLOUD_VISION_API_KEY=your_key_here
GOOGLE_CLOUD_TTS_API_KEY=your_key_here
GOOGLE_CLOUD_STT_API_KEY=your_key_here

# OpenAI
OPENAI_API_KEY=your_key_here

# Stripe
STRIPE_SECRET_KEY=sk_test_your_key_here
STRIPE_PUBLISHABLE_KEY=pk_test_your_key_here
```

**注意**: APIキーがなくても `USE_MOCK_APIS=true` で開発・テスト可能です。詳細は [mocking_strategy.md](mocking_strategy.md) を参照してください。

---

## Podman環境でのセットアップ

### 1. リポジトリのクローン

```bash
git clone https://github.com/clearclown/HaiLanGo.git
cd HaiLanGo
```

### 2. 環境変数の設定

```bash
cp .env.example .env
# .envファイルを編集（上記参照）
```

### 3. Podmanコンテナの起動

```bash
# すべてのサービスをビルド・起動
podman-compose up -d

# ログを確認
podman-compose logs -f
```

### 4. サービスの確認

```bash
# すべてのコンテナが起動しているか確認
podman-compose ps
```

出力例：
```
NAME                  IMAGE                       STATUS       PORTS
hailango-postgres     postgres:15-alpine          Up           0.0.0.0:5432->5432/tcp
hailango-redis        redis:7-alpine              Up           0.0.0.0:6379->6379/tcp
hailango-backend      localhost/backend:latest    Up           0.0.0.0:8080->8080/tcp
hailango-frontend     localhost/frontend:latest   Up           0.0.0.0:3000->3000/tcp
```

### 5. アクセス確認

- **フロントエンド**: http://localhost:3000
- **バックエンドAPI**: http://localhost:8080
- **ヘルスチェック**: http://localhost:8080/health

### 6. 停止とクリーンアップ

```bash
# サービスを停止
podman-compose down

# ボリュームも削除（データベースのデータも削除されます）
podman-compose down -v

# イメージも削除
podman-compose down --rmi all
```

---

## Docker環境でのセットアップ

Dockerを使用する場合も同様の手順です：

```bash
# 1. クローンと設定
git clone https://github.com/clearclown/HaiLanGo.git
cd HaiLanGo
cp .env.example .env

# 2. Docker Composeで起動
docker-compose up -d

# 3. ログ確認
docker-compose logs -f

# 4. 停止
docker-compose down
```

---

## 開発モード

コンテナを使わずにローカルで開発する場合：

### バックエンド開発

```bash
cd backend

# 依存関係のインストール
go mod download

# データベース・Redis起動（Podman）
podman-compose up -d postgres redis

# 開発サーバー起動（ホットリロード）
go run cmd/server/main.go

# または air を使用（ホットリロード）
air
```

### フロントエンド開発

```bash
cd frontend/web

# 依存関係のインストール
pnpm install

# 開発サーバー起動
pnpm dev
```

ブラウザで http://localhost:3000 を開く

### データベース接続

```bash
# PostgreSQLに接続
podman exec -it hailango-postgres psql -U HaiLanGo -d HaiLanGo_dev

# Redisに接続
podman exec -it hailango-redis redis-cli
```

---

## トラブルシューティング

### ポートが既に使用されている

```bash
# 使用中のポートを確認
lsof -i :8080  # バックエンド
lsof -i :3000  # フロントエンド
lsof -i :5432  # PostgreSQL
lsof -i :6379  # Redis

# プロセスを終了
kill -9 <PID>

# または .env でポートを変更
BACKEND_PORT=8081
FRONTEND_PORT=3001
```

### Podmanコンテナが起動しない

```bash
# コンテナのログを確認
podman-compose logs <service-name>

# 例: バックエンドのログ
podman-compose logs backend

# コンテナを再起動
podman-compose restart <service-name>

# すべてのコンテナを再ビルド
podman-compose up -d --build --force-recreate
```

### データベース接続エラー

```bash
# PostgreSQLが起動しているか確認
podman exec -it hailango-postgres pg_isready

# データベースが存在するか確認
podman exec -it hailango-postgres psql -U HaiLanGo -l

# 接続テスト
podman exec -it hailango-postgres psql -U HaiLanGo -d HaiLanGo_dev -c "SELECT 1;"
```

### フロントエンドのビルドエラー

```bash
# node_modules を削除して再インストール
cd frontend/web
rm -rf node_modules .next
pnpm install
pnpm build

# Dockerキャッシュをクリアして再ビルド
podman-compose build --no-cache frontend
```

### 権限エラー（Podman）

```bash
# rootlessモードで実行しているか確認
podman info | grep rootless

# ボリュームの権限を確認
podman volume inspect hailango_postgres_data

# 必要に応じて sudo を使用
sudo podman-compose up -d
```

---

## よくある質問

### Q1. APIキーがなくても開発できますか？

**A**: はい、可能です。`USE_MOCK_APIS=true` を設定することで、すべての外部API呼び出しが自動的にモックに置き換わります。詳細は [mocking_strategy.md](mocking_strategy.md) を参照してください。

### Q2. PodmanとDockerのどちらを使うべきですか？

**A**: 両方とも動作しますが、以下の理由でPodmanを推奨します：
- Rootlessで実行可能（セキュリティ向上）
- Dockerデーモン不要（軽量）
- Docker互換のコマンド

### Q3. データベースのデータを永続化するには？

**A**: Podman/Docker Composeではボリュームが自動的に作成され、データは永続化されます。ボリュームを削除しない限り、コンテナを停止・再起動してもデータは保持されます。

```bash
# ボリュームを削除せずに停止
podman-compose down

# ボリュームも削除（データが消える）
podman-compose down -v
```

### Q4. 本番環境へのデプロイ方法は？

**A**: 現在のセットアップは開発環境用です。本番環境へのデプロイについては、以下を参照してください：
- AWS/GCP/Cloudflareへのデプロイガイド（準備中）
- Kubernetes設定（`infra/k8s/` - 準備中）
- Terraform設定（`infra/terraform/` - 準備中）

### Q5. ホットリロードは動作しますか？

**A**: はい、以下の方法でホットリロードが可能です：
- **バックエンド**: `air` を使用（Goファイルの変更を自動検知）
- **フロントエンド**: `pnpm dev` で自動リロード

コンテナ内でもボリュームマウントにより、ホストのファイル変更が反映されます。

### Q6. テストの実行方法は？

```bash
# バックエンドテスト
cd backend
go test ./...

# フロントエンドテスト
cd frontend/web
pnpm test

# E2Eテスト
cd frontend/web
pnpm test:e2e
```

### Q7. ログの確認方法は？

```bash
# すべてのサービスのログ
podman-compose logs -f

# 特定のサービスのみ
podman-compose logs -f backend
podman-compose logs -f frontend

# 最新100行のみ表示
podman-compose logs --tail=100 backend
```

---

## 参考リンク

- [README.md](../README.md) - プロジェクト概要
- [CLAUDE.md](../CLAUDE.md) - 開発者向けドキュメント
- [mocking_strategy.md](mocking_strategy.md) - モックAPI戦略
- [requirements_definition.md](requirements_definition.md) - 要件定義書

---

## サポート

質問や問題がある場合は、以下を利用してください：

- **GitHub Issues**: https://github.com/clearclown/HaiLanGo/issues
- **Email**: support@HaiLanGo.com
- **Discord**: https://discord.gg/HaiLanGo

---

**最終更新**: 2025-11-14
