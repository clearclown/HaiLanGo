.PHONY: help run build test lint clean db-up db-down db-logs migrate-up migrate-down migrate-version setup

help: ## ヘルプを表示
	@echo "使用可能なコマンド:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## サーバーを起動
	go run cmd/server/main.go cmd/server/mock_repository.go

build: ## ビルド
	go build -o bin/server cmd/server/main.go cmd/server/mock_repository.go

test: ## テストを実行
	go test -v ./...

test-coverage: ## カバレッジ付きテスト
	go test -v -cover ./...

lint: ## リントを実行
	golangci-lint run || go vet ./...

clean: ## ビルド成果物を削除
	rm -rf bin/
	rm -f coverage.txt

deps: ## 依存関係をインストール
	go mod download
	go mod tidy

db-up: ## PostgreSQL/Redisを起動
	podman-compose up -d
	@echo "✓ データベースコンテナを起動しました"
	@echo "PostgreSQLの準備ができるまで待機中..."
	@sleep 5

db-down: ## PostgreSQL/Redisを停止
	podman-compose down
	@echo "✓ データベースコンテナを停止しました"

db-logs: ## PostgreSQLのログを表示
	podman-compose logs -f postgres

migrate-up: ## マイグレーションを実行
	@echo "マイグレーションを実行中..."
	go run cmd/migrate/main.go -command=up

migrate-down: ## マイグレーションをロールバック
	@echo "⚠  警告: 全テーブルとデータが削除されます!"
	@read -p "本当に実行しますか? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		go run cmd/migrate/main.go -command=down; \
	else \
		echo "キャンセルしました"; \
	fi

migrate-version: ## マイグレーションバージョンを表示
	go run cmd/migrate/main.go -command=version

setup: db-up migrate-up ## データベースセットアップ（起動+マイグレーション）
	@echo "✓ セットアップ完了！データベースの準備ができました"
