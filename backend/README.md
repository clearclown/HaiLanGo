# HaiLanGo Backend - 間隔反復学習（SRS）機能

## 概要

間隔反復学習（Spaced Repetition System）機能のバックエンド実装です。

## 機能

- ✅ SRSアルゴリズム実装
- ✅ 復習項目管理
- ✅ 優先度別復習項目取得（緊急・推奨・余裕あり）
- ✅ 復習完了処理
- ✅ 学習統計
- ✅ RESTful API

## 技術スタック

- **言語**: Go 1.21+
- **フレームワーク**: Gin
- **テスト**: testify
- **データベース**: PostgreSQL（予定）、Redis（予定）

## セットアップ

### 依存関係のインストール

```bash
go mod download
```

### サーバーの起動

```bash
go run cmd/server/main.go
```

サーバーはポート8080で起動します。

### テストの実行

```bash
# すべてのテスト
go test ./... -v

# カバレッジ付き
go test ./... -cover

# 特定のパッケージ
go test ./pkg/srs/... -v
```

## APIエンドポイント

### 復習項目取得
```
GET /api/v1/review/items/:user_id
```

### 復習完了
```
POST /api/v1/review/items/:item_id/complete
Content-Type: application/json

{
  "score": 85,
  "time_spent_sec": 30
}
```

### 統計情報取得
```
GET /api/v1/review/stats/:user_id
```

### ヘルスチェック
```
GET /health
```

## ディレクトリ構造

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # エントリーポイント
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   └── review_handler.go  # APIハンドラー
│   │   └── router/
│   │       └── router.go           # ルーティング
│   ├── service/
│   │   └── srs/
│   │       ├── srs.go              # ビジネスロジック
│   │       └── srs_test.go
│   ├── repository/
│   │   ├── mock_review_item.go     # リポジトリ
│   │   └── review_item_test.go
│   └── models/
│       └── review_item.go          # データモデル
└── pkg/
    └── srs/
        ├── algorithm.go            # SRSアルゴリズム
        └── algorithm_test.go
```

## SRSアルゴリズム

### 基本間隔
- 初回: 1日後
- 2回目: 3日後
- 3回目: 7日後
- 4回目: 14日後
- 5回目: 30日後
- 6回目以降: 60日後

### スコア調整
- 85点以上: 間隔を1.5倍に延長
- 70-84点: 通常の間隔
- 50-69点: 間隔を半分に短縮
- 50点未満: 翌日に再度復習

## テスト

すべてのテストがPASSしています：

- ✅ アルゴリズムテスト（6/6）
- ✅ リポジトリテスト（6/6）
- ✅ サービス層テスト（5/5）

## 開発

### コーディング規約

- コメントは日本語で記述
- `gofmt`でフォーマット
- テストファーストで実装（TDD）

### 新機能の追加

1. テストを書く
2. 実装する
3. テストが通ることを確認
4. リファクタリング

## 今後の予定

- [ ] PostgreSQL実装
- [ ] Redis キャッシュ
- [ ] JWT認証
- [ ] WebSocket通知
- [ ] カスタマイズ設定

## ライセンス

MIT License
