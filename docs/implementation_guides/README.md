# 実装指示書 - HaiLanGo プロジェクト

## 概要

このディレクトリには、HaiLanGo プロジェクトの未実装機能に関する詳細な実装指示書が含まれています。各指示書は、他のエンジニアが独立して作業できるように、詳細な手順とコード例を提供しています。

## 実装優先度と担当分担の推奨

### 🔴 最優先（すぐに着手すべき）

1. **[01_Review_Page.md](./01_Review_Page.md)** - 復習ページ
   - **目安時間**: 8-12時間
   - **スキル要件**: React, TypeScript, API統合
   - **重要度**: ⭐⭐⭐⭐⭐（ユーザー継続率に直結）
   - **担当推奨**: フロントエンド経験者

2. **[03_Upload_Page.md](./03_Upload_Page.md)** - 書籍アップロードページ
   - **目安時間**: 10-16時間
   - **スキル要件**: React, TypeScript, ファイル処理
   - **重要度**: ⭐⭐⭐⭐⭐（コア機能）
   - **担当推奨**: フロントエンド経験者

### 🟡 重要（早めに着手）

3. **[04_Stats_Dashboard.md](./04_Stats_Dashboard.md)** - 学習統計ダッシュボード
   - **目安時間**: 6-10時間
   - **スキル要件**: React, TypeScript, データ可視化
   - **重要度**: ⭐⭐⭐⭐（ユーザーモチベーション維持）
   - **担当推奨**: フロントエンド経験者

4. **[02_Flutter_Mobile_App_Setup.md](./02_Flutter_Mobile_App_Setup.md)** - Flutter モバイルアプリ初期セットアップ
   - **目安時間**: 4-6時間
   - **スキル要件**: Flutter, Dart
   - **重要度**: ⭐⭐⭐⭐（モバイル対応開始）
   - **担当推奨**: Flutter 経験者

### 🟢 通常（順次着手）

5. **Root ページのリダイレクト実装**
   - **目安時間**: 1-2時間
   - **スキル要件**: Next.js, TypeScript
   - **重要度**: ⭐⭐⭐（UX改善）
   - **内容**: `app/page.tsx` を `app/(home)/page.tsx` にリダイレクト

6. **Flutter 認証画面の実装**
   - **目安時間**: 8-12時間
   - **スキル要件**: Flutter, Dart, API統合
   - **重要度**: ⭐⭐⭐⭐（モバイルコア機能）
   - **前提**: `02_Flutter_Mobile_App_Setup.md` 完了後

7. **Flutter ホーム画面の実装**
   - **目安時間**: 6-10時間
   - **スキル要件**: Flutter, Dart, Riverpod
   - **重要度**: ⭐⭐⭐⭐（モバイルコア機能）
   - **前提**: Flutter 認証画面完了後

## 実装ガイド一覧

| # | ドキュメント | 機能 | 優先度 | 目安時間 | 状態 |
|---|------------|------|--------|---------|------|
| 01 | [Review_Page.md](./01_Review_Page.md) | 復習ページ（SRS） | 🔴 最優先 | 8-12h | ⏳ 未着手 |
| 02 | [Flutter_Mobile_App_Setup.md](./02_Flutter_Mobile_App_Setup.md) | Flutter初期セットアップ | 🟡 重要 | 4-6h | ⏳ 未着手 |
| 03 | [Upload_Page.md](./03_Upload_Page.md) | 書籍アップロード | 🔴 最優先 | 10-16h | ⏳ 未着手 |
| 04 | [Stats_Dashboard.md](./04_Stats_Dashboard.md) | 学習統計ダッシュボード | 🟡 重要 | 6-10h | ⏳ 未着手 |

## 並行作業の推奨

### チーム構成例: 4人

**担当A（フロントエンド経験者）**:
- 01_Review_Page.md（8-12h）
- その後: 04_Stats_Dashboard.md（6-10h）

**担当B（フロントエンド経験者）**:
- 03_Upload_Page.md（10-16h）
- その後: Root ページリダイレクト（1-2h）

**担当C（Flutter 経験者）**:
- 02_Flutter_Mobile_App_Setup.md（4-6h）
- その後: Flutter 認証画面（8-12h）

**担当D（Flutter 経験者）**:
- 02_Flutter_Mobile_App_Setup.md のサポート（最初の2日）
- その後: Flutter ホーム画面の準備・設計

### チーム構成例: 2人

**担当A（フルスタック）**:
- 01_Review_Page.md（8-12h）
- 03_Upload_Page.md（10-16h）
- 04_Stats_Dashboard.md（6-10h）

**担当B（Flutter 経験者）**:
- 02_Flutter_Mobile_App_Setup.md（4-6h）
- Flutter 認証画面（8-12h）
- Flutter ホーム画面（6-10h）

## 実装前の確認事項

### 環境セットアップ

1. **バックエンドの起動確認**:
   ```bash
   cd backend
   go run cmd/server/main.go
   # または
   podman-compose -f ../podman-compose-db-only.yml up -d
   go run cmd/server/main.go
   ```

2. **フロントエンドの起動確認**:
   ```bash
   cd frontend/web
   pnpm install
   pnpm run dev
   ```

3. **Flutter環境の確認**（Flutter担当者のみ）:
   ```bash
   flutter doctor
   ```

### 必要なツール

- **フロントエンド開発者**:
  - Node.js 18+
  - pnpm
  - VS Code + 拡張機能（ESLint, Prettier, TypeScript）

- **Flutter 開発者**:
  - Flutter SDK 3.0+
  - Android Studio / Xcode
  - VS Code + Flutter 拡張機能

- **全員**:
  - Git
  - ブラウザ（Chrome推奨）

## 実装の進め方

### 1. 着手前

1. 担当する実装ガイドを熟読
2. 前提条件をすべて満たしているか確認
3. 疑問点があれば team slack/discord で質問

### 2. 実装中

1. 実装ガイドのStep順に実装
2. 各Stepごとにコミット
3. コミットメッセージ例: `feat(review): add ReviewCard component`
4. 詰まったら30分で質問（時間を無駄にしない）

### 3. 完了後

1. 「完了条件」をすべて確認
2. テスト方法に従ってテスト
3. スクリーンショットを撮影（UI関連の場合）
4. プルリクエストを作成
5. Slack/Discordで完了報告

## プルリクエストのテンプレート

```markdown
## 実装内容

実装ガイド: `01_Review_Page.md`

### 完了したタスク

- [x] 型定義の作成
- [x] API クライアントの拡張
- [x] ReviewCard コンポーネント
- [x] ReviewSession コンポーネント
- [x] Review ページ
- [x] エラーハンドリング

### スクリーンショット

（画像を添付）

### テスト結果

- [x] 復習統計が表示される
- [x] 復習セッションが開始できる
- [x] スコア送信が正常に動作する

### 備考

特になし
```

## トラブルシューティング

### 共通の問題

1. **APIエラー**: バックエンドが起動しているか確認
   ```bash
   curl http://localhost:8080/health
   ```

2. **TypeScriptエラー**: 型定義が正しいか確認
   ```bash
   pnpm run type-check
   ```

3. **依存関係エラー**: node_modules を再インストール
   ```bash
   rm -rf node_modules pnpm-lock.yaml
   pnpm install
   ```

### 質問・サポート

- **Slack/Discord**: #hailango-dev チャンネル
- **GitHub Issues**: バグ報告・機能提案
- **Code Review**: プルリクエストでレビュー依頼

## 参考資料

### プロジェクト全体

- [README.md](../../README.md) - プロジェクト概要
- [CLAUDE.md](../../CLAUDE.md) - 開発ガイド
- [requirements_definition.md](../requirements_definition.md) - 要件定義書
- [ui_ux_design_document.md](../ui_ux_design_document.md) - UI/UX設計書

### 機能別詳細仕様（RD）

- [docs/featureRDs/](../featureRDs/) - 各機能の詳細仕様

### 技術スタック

- [Next.js公式](https://nextjs.org/docs)
- [Flutter公式](https://flutter.dev/docs)
- [Riverpod公式](https://riverpod.dev/)
- [Go公式](https://golang.org/doc/)

## 更新履歴

- 2025-11-14: 初版作成（01-04の実装ガイド）
- 今後、実装ガイドを追加予定

---

**質問・提案があれば、気軽に Slack/Discord で共有してください！**
