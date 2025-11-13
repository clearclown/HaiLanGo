# HaiLanGo - AI-Powered Language Learning Platform

<div align="center">

📚 既存の言語学習本 × 🤖 AI技術 = 🚀 個人に最適化された能動的な学習体験

[![Tests](https://github.com/clearclown/HaiLanGo/actions/workflows/test.yml/badge.svg)](https://github.com/clearclown/HaiLanGo/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Next.js](https://img.shields.io/badge/Next.js-14+-000000?style=flat&logo=next.js)](https://nextjs.org)
[![Flutter](https://img.shields.io/badge/Flutter-3.0+-02569B?style=flat&logo=flutter)](https://flutter.dev)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

</div>

## 🎯 プロジェクト概要

HaiLanGoは、手持ちの言語学習本をAI技術で最大限に活用する革新的な学習プラットフォームです。OCR、TTS、STTなどの最新AI技術を組み合わせ、24/7利用可能な「AI教師」があなたの語学学習をサポートします。

### コアコンセプト
- **既存の本を活用**: お気に入りの言語学習本を持っているなら、それが最高の教材
- **AI教師**: 音声読み上げ、発音評価、解説生成で個別指導
- **日常会話レベルを目指す**: 文法よりも会話パターンと単語量を重視
- **プライバシー重視**: E2E暗号化で書籍データは完全にプライベート

## ✨ 主な機能

### 📖 書籍のデジタル化
- **AI-OCR**: 多言語対応（12言語以上）
- **対応フォーマット**: PDF、PNG、JPG、HEIC
- **複雑なレイアウト対応**: 表、ルビ、縦書きもOK

### 🎧 音声機能
- **TTS（読み上げ）**:
  - 主要12言語サポート
  - 速度調整（0.5x〜2.0x）
  - 無料版（標準品質）/ プレミアム版（高品質）
- **STT（音声認識）**:
  - 発音評価（0-100点スコア）
  - 具体的な改善点提示
  - 英会話教室のようなフィードバック

### 🎓 教師モード（自動学習モード）
- **ボタン一つで連続学習**: 全ページを自動で順次再生
- **バックグラウンド再生**: 画面オフでも学習継続
- **カスタマイズ可能**: 速度、間隔、学習内容を自由に設定
- **オフライン対応**: 事前ダウンロードで通信不要
- **ユースケース**: 通勤・通学、家事中、就寝前のリスニング

### 📚 インタラクティブ学習
- **ページバイページモード**: 1ページずつ丁寧に学習
- **フレーズ練習**: リピート、ロールプレイ
- **発音チェック**: リアルタイムフィードバック
- **単語帳自動生成**: 学習中の単語を自動収集

### 📊 学習管理
- **間隔反復学習（SRS）**: 科学的に最適化された復習タイミング
- **進捗トラッキング**: 学習時間、完了ページ、習得単語数
- **ストリーク機能**: 連続学習日数でモチベーション維持
- **統計ダッシュボード**: 詳細な学習分析

## 🌍 対応言語（主要）

| 言語 | TTS | STT | OCR |
|------|-----|-----|-----|
| 日本語 | ✅ | ✅ | ✅ |
| 英語 | ✅ | ✅ | ✅ |
| 中国語 | ✅ | ✅ | ✅ |
| ロシア語 | ✅ | ✅ | ✅ |
| ペルシャ語 | ✅ | ✅ | ✅ |
| ヘブライ語 | ✅ | ✅ | ✅ |
| スペイン語 | ✅ | ✅ | ✅ |
| フランス語 | ✅ | ✅ | ✅ |
| ポルトガル語 | ✅ | ✅ | ✅ |
| ドイツ語 | ✅ | ✅ | ✅ |
| イタリア語 | ✅ | ✅ | ✅ |
| トルコ語 | ✅ | ✅ | ✅ |

その他のマイナー言語もサポートするが正確さはわからない

## 🛠 技術スタック

### バックエンド
- **言語**: Go 1.21+
- **フレームワーク**: Gin / Echo
- **データベース**: PostgreSQL 15+, Redis 7+
- **API**: RESTful + WebSocket

### フロントエンド
- **Web**: Next.js 14+ (TypeScript, React, TailwindCSS, ShadCN/UI)
- **Mobile**: Flutter 3.0+
- **状態管理**: React Context / Zustand (Web), Riverpod (Flutter)

### AI / 外部API
- **OCR**: Google Vision API / Azure Computer Vision / MarkPDFdown ⭐ NEW
- **TTS**: Google Cloud TTS / Amazon Polly / ElevenLabs
- **STT**: Google Cloud STT / Whisper API / OpenAI Realtime API ⭐ NEW
- **リアルタイム対話**: OpenAI Realtime API / gpt-realtime ⭐ NEW
- **翻訳**: DeepL API / Google Translate API
- **辞書**: Oxford Dictionary API / Wiktionary API
- **決済**: Stripe
- **詳細**: [API統合提案書](docs/api_integration_proposal.md) を参照

### インフラ
- **初期**: オンプレミス（Podman / Docker Compose）
- **将来**: AWS / GCP / Cloudflare
- **IaC**: Terraform
- **CI/CD**: GitHub Actions

## 📁 プロジェクト構造

```
HaiLanGo/
├── docs/                          # ドキュメント
│   ├── requirements_definition.md # 要件定義書
│   ├── ui_ux_design_document.md  # UI/UX設計書
│   └── teacher_mode_technical_spec.md # 教師モード技術仕様書
├── backend/                       # バックエンド（Go）
│   ├── cmd/                       # エントリーポイント
│   ├── internal/                  # 内部パッケージ
│   │   ├── api/                   # APIハンドラー
│   │   ├── service/               # ビジネスロジック
│   │   ├── repository/            # データアクセス層
│   │   └── models/                # データモデル
│   ├── pkg/                       # 再利用可能なパッケージ
│   └── go.mod
├── frontend/                      # フロントエンド
│   ├── web/                       # Next.js Webアプリ
│   └── mobile/                    # Flutterモバイルアプリ
├── infra/                         # インフラ設定
│   ├── docker-compose.yml         # ローカル開発環境
│   ├── terraform/                 # クラウドインフラ
│   └── k8s/                       # Kubernetes設定（将来）
├── scripts/                       # 各種スクリプト
├── .env.example                   # 環境変数テンプレート
├── CLAUDE.md                      # Claude Code設定
└── README.md                      # このファイル
```

## 🚀 クイックスタート

### 前提条件
- Go 1.21+
- Node.js 18+
- Podman / Docker
- PostgreSQL 15+
- Redis 7+

### 1. リポジトリのクローン

```bash
git clone https://github.com/clearclown/HaiLanGo.git
cd HaiLanGo
```

### 2. 環境変数の設定

```bash
cp .env.example .env
# .envファイルを編集して必要なAPIキーを設定
#
# 重要: APIキーがなくても開発・テストは可能です！
# USE_MOCK_APIS=true を設定すると自動的にモックが使用されます
# 詳細は docs/mocking_strategy.md を参照してください
```

### 3. 開発環境の起動

```bash
# Podmanを使用
podman-compose up -d

# または Docker Compose
docker-compose up -d
```

### 4. バックエンドの起動

```bash
cd backend
go mod download
go run cmd/server/main.go
```

### 5. フロントエンド（Web）の起動

```bash
cd frontend/web
pnpm install
pnpm run dev
```

ブラウザで http://localhost:3000 を開く

### 6. モバイルアプリの起動（オプション）

```bash
cd frontend/mobile
flutter pub get
flutter run
```

## 📖 ドキュメント

詳細なドキュメントは`docs/`ディレクトリを参照してください：

- [要件定義書](docs/requirements_definition.md) - プロジェクトの全体像と機能要件
- [UI/UX設計書](docs/ui_ux_design_document.md) - 画面設計とワイヤーフレーム
- [教師モード技術仕様書](docs/teacher_mode_technical_spec.md) - 自動学習モードの詳細仕様
- [モック構築戦略](docs/mocking_strategy.md) - APIキーなしでもテスト可能な仕組み
- [API統合提案書](docs/api_integration_proposal.md) - 統合可能な外部API・ツールの包括的調査

### 機能実装RD（Feature Requirements Documents）

各機能の詳細な実装仕様は `docs/featureRDs/` を参照してください：

- [1. ユーザー認証](docs/featureRDs/1_ユーザー認証.md)
- [2. 書籍アップロード](docs/featureRDs/2_書籍アップロード.md)
- [3. OCR処理](docs/featureRDs/3_OCR処理.md)
- [4. TTS音声読み上げ](docs/featureRDs/4_TTS音声読み上げ.md)
- [5. STT発音評価](docs/featureRDs/5_STT発音評価.md)
- [6. ページバイページ学習モード](docs/featureRDs/6_ページバイページ学習モード.md)
- [7. 教師モード自動学習](docs/featureRDs/7_教師モード自動学習.md)
- [8. 間隔反復学習SRS](docs/featureRDs/8_間隔反復学習SRS.md)
- [9. 単語帳機能](docs/featureRDs/9_単語帳機能.md)
- [10. 学習統計ダッシュボード](docs/featureRDs/10_学習統計ダッシュボード.md)
- [11. 決済統合Stripe](docs/featureRDs/11_決済統合Stripe.md) ✅ **実装完了**
- [12. 辞書API統合](docs/featureRDs/12_辞書API統合.md)
- [13. OCR結果手動修正](docs/featureRDs/13_OCR結果手動修正.md)
- [14. 会話パターン抽出](docs/featureRDs/14_会話パターン抽出.md)
- [15. WebSocketリアルタイム通知](docs/featureRDs/15_WebSocketリアルタイム通知.md)
- [16. ホーム画面実装](docs/featureRDs/16_ホーム画面実装.md)
- [17. 設定画面実装](docs/featureRDs/17_設定画面実装.md)
- [18. GitHub CI設定](docs/featureRDs/18_GitHub_CI設定.md)

## 💰 料金プラン

### 無料プラン
- 1日1ページまで学習
- 1日30分まで使用
- 標準品質のTTS
- 基本的な学習統計

### プレミアムプラン（$9.99/月）
- 無制限の学習
- 高品質TTS（ElevenLabs）
- オフライン音声ダウンロード
- 詳細な学習分析
- 優先サポート

## 🗓 開発ロードマップ

### Phase 1: MVP（3-4ヶ月） - 進行中
- [ ] ユーザー認証（OAuth + Email）
- [ ] PDFアップロード + OCR
- [ ] TTS基本機能（主要言語5つ）
- [ ] 簡易な単語帳機能
- [ ] Web版のみ
- [x] **会話パターン自動抽出** ✅

### Phase 2: コア機能（2-3ヶ月）
- [x] STT + 発音評価 ✅ 実装完了
- [ ] ページバイページ学習モード
- [ ] 間隔反復学習アルゴリズム
- [ ] モバイルアプリ（Flutter）
- [ ] Stripe決済統合

### Phase 3: 拡張機能（3-4ヶ月）
- [ ] 教師モード（オフライン対応）
- [ ] 辞書API統合
- [ ] 学習分析ダッシュボード
- [ ] マイナー言語対応拡大

### Phase 4: コミュニティ（時期未定）
- [ ] ユーザー生成コンテンツ機能
- [ ] ブログ投稿プラットフォーム
- [ ] コミュニティフォーラム

## 🤝 コントリビューション

コントリビューションは大歓迎です！以下の手順でお願いします：

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

### コーディング規約
- Go: `gofmt`, `golangci-lint`を使用
- TypeScript: Biome.jsでフォーマット・リント
- テスト: Vitest（単体・統合）、Playwright（E2E）
- コミットメッセージ: Conventional Commits形式

### テスト戦略
- **TDD原則**: すべての機能はテストファーストで実装
- **モックシステム**: APIキーなしでもテスト可能（`USE_MOCK_APIS=true`）
- **CI/CD**: GitHub Actionsで自動テスト実行
- 詳細は [モック構築戦略](docs/mocking_strategy.md) を参照

## 📄 ライセンス

このプロジェクトは MIT License の下でライセンスされています。詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 🙏 謝辞

- **abceed** - UI/UXデザインのインスピレーション
- **Duolingo** - ゲーミフィケーションのアイデア
- **Anki** - 間隔反復学習アルゴリズム

## 📞 サポート・お問い合わせ

- **Issue**: [GitHub Issues](https://github.com/clearclown/HaiLanGo/issues)
- **Email**: support@HaiLanGo.com
- **Discord**: [Community Server](https://discord.gg/HaiLanGo)

---

<div align="center">

Made with ❤️ by [Your Name]

⭐ このプロジェクトが気に入ったらスターをお願いします！

</div>