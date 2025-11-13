# 機能実装: 決済統合（Stripe）

## 要件

### 機能要件

#### Stripe決済統合
1. **サブスクリプション管理**
   - プレミアムプラン（$9.99/月）
   - 年額プラン（$99.99/年）
   - サブスクリプションの開始・更新・キャンセル

2. **決済フロー**
   - 決済画面の表示
   - クレジットカード入力
   - 決済処理
   - 決済完了通知

3. **プラン管理**
   - プラン情報の取得
   - 現在のプラン状態
   - プラン変更

### 非機能要件

- **パフォーマンス**:
  - 決済処理: 3秒以内
  - プラン情報取得: 100ms以内

- **セキュリティ**:
  - PCI DSS準拠（Stripe経由）
  - 決済情報の暗号化
  - 不正アクセス防止

- **拡張性**:
  - 将来的なプラン追加への対応
  - 複数通貨対応

- **エラーハンドリング**:
  - 決済エラーの処理
  - リトライロジック
  - ユーザーへのエラー通知

## 実装指示

### Step 1: テスト設計

1. **ユニットテスト**
   - Stripe API呼び出し（モック）
   - サブスクリプション作成
   - プラン情報取得

2. **統合テスト**
   - 決済フロー
   - サブスクリプション更新
   - Webhook処理

3. **エッジケースのテスト**
   - 決済失敗
   - カードエラー
   - Webhook重複処理

テストファイルは `backend/internal/service/payment/payment_test.go` に配置。

### Step 2: 実装

#### 実装ファイル構造

```
backend/
├── internal/
│   ├── service/
│   │   └── payment/
│   │       ├── service.go           # 決済サービス
│   │       ├── stripe.go            # Stripe統合
│   │       └── service_test.go
│   └── api/
│       └── payment/
│           ├── handler.go           # HTTPハンドラー
│           └── webhook.go          # Webhook処理
```

#### APIエンドポイント

```
POST /api/v1/payment/subscribe
GET  /api/v1/payment/plans
POST /api/v1/payment/cancel
POST /api/v1/payment/webhook
```

## 制約事項

- 既存の `backend/internal/models/` は変更しない（新規追加のみ）
- `backend/internal/service/payment/` 配下のみ編集可能

## 完了条件

- [ ] すべてのテストが通る
- [ ] lintエラーがない
- [ ] タイプエラーがない
- [ ] ドキュメントが更新されている
- [ ] 決済フローの動作確認
- [ ] Webhook処理の動作確認
