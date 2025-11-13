# 機能実装: 辞書API統合

## 要件

### 機能要件

#### 辞書API統合
1. **対応API**
   - Oxford Dictionary API（優先）
   - Free Dictionary API（フォールバック）
   - Wiktionary API（補助）

2. **単語情報取得**
   - 意味（複数）
   - 発音記号
   - 例文
   - 品詞
   - 類義語・対義語

3. **キャッシュ**
   - 辞書結果のキャッシュ（30日間）
   - Redisで管理
   - API呼び出し回数の削減

### 非機能要件

- **パフォーマンス**:
  - 辞書情報取得: 500ms以内
  - キャッシュヒット率: 80%以上

- **セキュリティ**:
  - APIキーの安全な管理
  - レート制限の遵守

- **拡張性**:
  - 複数APIの切り替え対応
  - 将来的なAPI追加への対応

- **エラーハンドリング**:
  - APIエラー時のフォールバック
   - リトライロジック

## 実装指示

### Step 1: テスト設計

1. **ユニットテスト**
   - 辞書API呼び出し（モック）
   - データパース
   - キャッシュ操作

2. **統合テスト**
   - 辞書情報取得フロー
   - キャッシュの動作
   - フォールバック処理

3. **エッジケースのテスト**
   - 存在しない単語
   - APIエラー
   - レート制限

テストファイルは `backend/pkg/dictionary/dictionary_test.go` に配置。

### Step 2: 実装

#### 実装ファイル構造

```
backend/
├── pkg/
│   └── dictionary/
│       ├── oxford.go                # Oxford Dictionary API
│       ├── free_dictionary.go       # Free Dictionary API
│       ├── wiktionary.go            # Wiktionary API
│       ├── service.go               # 辞書サービス
│       └── dictionary_test.go
└── internal/
    └── service/
        └── dictionary/
            ├── service.go           # 辞書統合サービス
            └── service_test.go
```

#### APIエンドポイント

```
GET /api/v1/dictionary/words/{word}
GET /api/v1/dictionary/words/{word}/details
```

## 制約事項

- 既存の `backend/internal/models/` は変更しない（新規追加のみ）
- `backend/pkg/dictionary/` 配下のみ編集可能

## 完了条件

- [ ] すべてのテストが通る
- [ ] lintエラーがない
- [ ] タイプエラーがない
- [ ] ドキュメントが更新されている
- [ ] 辞書情報取得の動作確認
- [ ] キャッシュの動作確認
