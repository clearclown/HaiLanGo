# HaiLanGo - ドキュメント

## 📚 ドキュメント一覧

このディレクトリには、HaiLanGoプロジェクトのすべてのドキュメントが含まれています。

---

## 🚀 クイックスタート

- **[SETUP.md](SETUP.md)** - セットアップガイド（Podman/Docker環境）
  - 環境構築の手順
  - トラブルシューティング
  - よくある質問

---

## 📋 プロジェクト概要

- **[requirements_definition.md](requirements_definition.md)** - 要件定義書
  - プロジェクトのコアコンセプト
  - 機能要件と非機能要件
  - 技術スタック
  - ビジネスモデル
  - 開発ロードマップ

- **[ui_ux_design_document.md](ui_ux_design_document.md)** - UI/UX設計書
  - デザインコンセプト
  - カラーパレットとタイポグラフィ
  - 画面構成・ワイヤーフレーム
  - コンポーネント設計
  - アクセシビリティ

---

## 🛠 技術仕様書

詳細な技術仕様は [technical/](technical/) ディレクトリを参照してください。

### コア機能の技術仕様

- **[technical/teacher_mode.md](technical/teacher_mode.md)** - 教師モード（自動学習モード）
  - リアルタイム音声再生
  - バックグラウンド再生
  - オフライン対応
  - カスタマイズ設定

- **[technical/websocket.md](technical/websocket.md)** - WebSocketリアルタイム通知
  - WebSocket実装詳細
  - クライアント・サーバー間通信
  - 接続管理とエラーハンドリング

- **[technical/ocr_implementation.md](technical/ocr_implementation.md)** - OCR処理実装
  - OCRパイプライン
  - 多言語対応
  - 精度向上戦略

完全なリストは [technical/README.md](technical/README.md) を参照してください。

---

## 📦 機能実装RD（Feature Requirements Documents）

各機能の詳細な実装仕様は [featureRDs/](featureRDs/) ディレクトリを参照してください。

### 実装完了（100%）

**Phase 1-5のすべての機能が実装完了しました！** 🎉

- **Phase 1 (MVP)**: 6機能 ✅ 100%
- **Phase 2 (コア機能)**: 3機能 ✅ 100%
- **Phase 3 (拡張機能)**: 6機能 ✅ 100%
- **Phase 4 (UI/UX)**: 2機能 ✅ 100%
- **Phase 5 (インフラ)**: 1機能 ✅ 100%

詳細は以下を参照：
- **[featureRDs/README.md](featureRDs/README.md)** - 機能一覧とステータス
- **[IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md)** - 実装状況サマリー
- **[featureRDs/archives/](featureRDs/archives/)** - 完了済み機能の詳細

---

## 🧪 開発・テスト

- **[mocking_strategy.md](mocking_strategy.md)** - モック構築戦略
  - APIキーなしでもテスト可能な仕組み
  - モックシステムの実装方法
  - テストでの使用方法
  - CI/CDでの活用

---

## 🔌 API統合

- **[api_integration_proposal.md](api_integration_proposal.md)** - API統合提案書
  - 統合可能な外部API・ツールの調査
  - OpenAI Realtime API
  - MarkPDFdown
  - DeepL API
  - その他の推奨API

---

## 📊 実装状況

- **[IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md)** - 実装状況サマリー
  - 全18機能の実装完了状況
  - Phase別進捗
  - 統計情報（コード行数、テストカバレッジなど）
  - 次のステップ

---

## 📁 ディレクトリ構造

```
docs/
├── README.md                          # このファイル（ドキュメントインデックス）
├── SETUP.md                           # セットアップガイド
├── IMPLEMENTATION_STATUS.md           # 実装状況サマリー
│
├── requirements_definition.md         # 要件定義書
├── ui_ux_design_document.md           # UI/UX設計書
├── mocking_strategy.md                # モック構築戦略
├── api_integration_proposal.md        # API統合提案書
│
├── technical/                         # 技術仕様書ディレクトリ
│   ├── README.md                      # 技術仕様インデックス
│   ├── teacher_mode.md                # 教師モード技術仕様
│   ├── websocket.md                   # WebSocket実装
│   └── ocr_implementation.md          # OCR実装サマリー
│
└── featureRDs/                        # 機能要件ドキュメント
    ├── README.md                      # 機能一覧
    └── archives/                      # 完了済み機能（23ファイル）
        ├── 1_ユーザー認証.md
        ├── 2_書籍アップロード.md
        ├── ...
        └── 18_GitHub_CI設定.md
```

---

## 🔄 ドキュメント更新フロー

1. **新機能の追加**
   - `featureRDs/` に新しいRDを作成
   - `featureRDs/README.md` を更新

2. **機能の完了**
   - 実装完了後、RDを `featureRDs/archives/` に移動
   - `IMPLEMENTATION_STATUS.md` を更新
   - `featureRDs/README.md` の進捗を更新

3. **技術仕様の追加**
   - `technical/` に新しい技術仕様を作成
   - `technical/README.md` にリンクを追加

---

## 📖 ドキュメント執筆ガイドライン

### Markdown形式

すべてのドキュメントはMarkdown（`.md`）形式で記述します。

### 見出し階層

```markdown
# H1 - ドキュメントタイトル（1つのみ）
## H2 - メインセクション
### H3 - サブセクション
#### H4 - 詳細セクション
```

### コードブロック

```markdown
```言語名
コード
```
```

### リンク

```markdown
- 相対リンク: [SETUP.md](SETUP.md)
- 絶対リンク: https://github.com/clearclown/HaiLanGo
```

### 更新履歴

各ドキュメントの末尾に更新履歴を記載します：

```markdown
---

**最終更新**: 2025-11-14
**更新者**: Your Name
**変更内容**: 機能Xの実装完了に伴う更新
```

---

## 🔗 関連リソース

### プロジェクト関連

- **[../README.md](../README.md)** - プロジェクトのREADME（トップレベル）
- **[../CLAUDE.md](../CLAUDE.md)** - Claude Code設定（開発者向け）
- **[../.env.example](../.env.example)** - 環境変数テンプレート

### 外部リンク

- **GitHub Repository**: https://github.com/clearclown/HaiLanGo
- **公式サイト**: https://HaiLanGo.com（準備中）
- **Discord Community**: https://discord.gg/HaiLanGo

---

## 📞 サポート・お問い合わせ

ドキュメントに関する質問や改善提案がある場合：

- **GitHub Issues**: https://github.com/clearclown/HaiLanGo/issues
- **Email**: support@HaiLanGo.com
- **Discord**: https://discord.gg/HaiLanGo

---

## 🤝 コントリビューション

ドキュメントの改善に貢献したい場合：

1. このリポジトリをフォーク
2. ドキュメントを編集
3. プルリクエストを作成

詳細は [CONTRIBUTING.md](../CONTRIBUTING.md)（準備中）を参照してください。

---

## 📝 ライセンス

このプロジェクトのドキュメントは [MIT License](../LICENSE) の下で公開されています。

---

**最終更新**: 2025-11-14
**管理者**: HaiLanGo Development Team
