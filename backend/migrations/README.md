# データベースマイグレーション

このディレクトリには、HaiLanGoプロジェクトのデータベースマイグレーションファイルが含まれています。

## マイグレーションファイルの命名規則

```
{version}_{description}.{up|down}.sql
```

例: `001_create_subscription_tables.up.sql`

## マイグレーションの実行

### 手動実行（開発環境）

```bash
# PostgreSQLに接続
psql -U HaiLanGo -d HaiLanGo_dev

# マイグレーションを実行
\i backend/migrations/001_create_subscription_tables.up.sql
```

### ロールバック

```bash
# PostgreSQLに接続
psql -U HaiLanGo -d HaiLanGo_dev

# ロールバックを実行
\i backend/migrations/001_create_subscription_tables.down.sql
```

### migrate ツールを使用（推奨）

```bash
# Install migrate tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path backend/migrations -database "postgresql://HaiLanGo:password@localhost:5432/HaiLanGo_dev?sslmode=disable" up

# Rollback
migrate -path backend/migrations -database "postgresql://HaiLanGo:password@localhost:5432/HaiLanGo_dev?sslmode=disable" down
```

## 現在のマイグレーション

### 001_create_subscription_tables

決済・サブスクリプション関連のテーブルを作成します。

**テーブル:**
- `subscription_plans` - サブスクリプションプラン
- `subscriptions` - ユーザーのサブスクリプション
- `payments` - 決済履歴

**デフォルトプラン:**
- Premium Monthly ($9.99/月)
- Premium Yearly ($99.99/年)

## 新しいマイグレーションの追加

1. 次のバージョン番号を使用してファイルを作成
2. `.up.sql` と `.down.sql` の両方を作成
3. テーブル作成時は必ず `IF NOT EXISTS` を使用
4. インデックスの作成を忘れずに
5. 外部キー制約を適切に設定

## 注意事項

- 本番環境では必ずバックアップを取ってからマイグレーションを実行してください
- ロールバックのテストも忘れずに行ってください
- マイグレーションは順序に依存するため、番号順に実行してください
