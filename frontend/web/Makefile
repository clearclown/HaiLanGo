.PHONY: help install dev build test lint format clean

help: ## ヘルプを表示
	@echo "使用可能なコマンド:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## 依存関係をインストール
	pnpm install

dev: ## 開発サーバーを起動
	pnpm run dev

build: ## プロダクションビルド
	pnpm run build

start: ## プロダクションサーバーを起動
	pnpm start

test: ## すべてのテストを実行
	pnpm run test:unit

test-watch: ## テストをウォッチモードで実行
	pnpm run test:watch

test-e2e: ## E2Eテストを実行
	pnpm run test:e2e

test-coverage: ## カバレッジ付きテスト
	pnpm run test:coverage

lint: ## リントを実行
	pnpm run lint

format: ## フォーマットを実行
	pnpm run format

type-check: ## 型チェックを実行
	pnpm run type-check

clean: ## ビルド成果物を削除
	rm -rf .next node_modules out coverage test-results playwright-report
