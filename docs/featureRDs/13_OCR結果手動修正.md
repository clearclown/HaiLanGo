# 機能実装: OCR結果手動修正

## 要件

### 機能要件

#### OCR結果の手動修正機能
1. **修正機能**
   - OCR結果のテキストを編集可能
   - ページ単位での修正
   - フレーズ単位での修正
   - 変更履歴の保存

2. **UI機能**
   - インライン編集
   - 差分表示（OCR結果 vs 修正後）
   - 一括修正（複数ページ）
   - 修正の取り消し・やり直し

3. **検証機能**
   - 修正内容のバリデーション
   - 言語検証
   - 文字数制限チェック

### 非機能要件

- **パフォーマンス**:
  - 修正の保存: 500ms以内
  - 差分計算: 100ms以内

- **セキュリティ**:
  - ユーザー認証必須
  - 修正履歴の監査ログ

- **拡張性**:
  - 将来的なAI補助修正への対応

## 実装指示

### Step 1: テスト設計

1. **ユニットテスト（Vitest）**
   - テキスト編集ロジック
   - 差分計算
   - バリデーション

2. **統合テスト（Vitest）**
   - 修正保存フロー
   - 変更履歴の記録

3. **E2Eテスト（Playwright）**
   - ブラウザでのテキスト編集
   - 修正の保存
   - 変更履歴の表示

テストファイルは `frontend/web/components/ocr-editor/` に配置。

### Step 2: 実装

#### 実装ファイル構造

```
frontend/web/
├── components/
│   └── ocr-editor/
│       ├── OCRTextEditor.tsx        # テキストエディタ
│       ├── DiffViewer.tsx           # 差分表示
│       └── OCRTextEditor.test.tsx   # Vitestテスト
├── e2e/
│   └── ocr-editor.spec.ts          # Playwrightテスト
└── services/
    └── ocrApi.ts                    # OCR API呼び出し

backend/
├── internal/
│   ├── api/
│   │   └── ocr/
│   │       ├── handler.go           # HTTPハンドラー
│   │       └── handler_test.go
│   └── service/
│       └── ocr/
│           ├── editor.go            # 修正サービス
│           └── editor_test.go
```

#### APIエンドポイント

```
PUT  /api/v1/books/{book_id}/pages/{page_id}/ocr-text
GET  /api/v1/books/{book_id}/pages/{page_id}/ocr-history
```

## 制約事項

- 既存の `backend/internal/models/` は変更しない（新規追加のみ）
- `frontend/web/components/ocr-editor/` 配下のみ編集可能

## 完了条件

- [ ] すべてのテストが通る（Vitest + Playwright）
- [ ] BiomeJSエラーがない
- [ ] タイプエラーがない
- [ ] GitHub CIが通る
- [ ] ドキュメントが更新されている
