.PHONY: help run build test lint clean

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
