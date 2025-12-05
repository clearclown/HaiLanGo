# HaiLanGo - 要件定義書 v2

## 概要

このディレクトリには、HaiLanGoプロジェクトの完全な要件定義書が含まれています。

**重要**: このv2ディレクトリがSingle Source of Truth（SSOT）です。
旧ドキュメント（`docs/` 直下）との矛盾がある場合、v2を優先してください。

---

## ドキュメント一覧

| ファイル | 内容 | 対象読者 |
|----------|------|----------|
| [00_PROJECT_VISION.md](./00_PROJECT_VISION.md) | プロジェクトビジョン・コンセプト | 全員 |
| [01_CORE_USECASES.md](./01_CORE_USECASES.md) | コアユースケース定義 | 全員 |
| [02_SYSTEM_ARCHITECTURE.md](./02_SYSTEM_ARCHITECTURE.md) | システムアーキテクチャ | エンジニア |
| [03_DATABASE_SCHEMA.md](./03_DATABASE_SCHEMA.md) | データベース設計 | バックエンド |
| [04_API_SPECIFICATION.md](./04_API_SPECIFICATION.md) | REST API仕様 | フルスタック |
| [05_TEST_STRATEGY.md](./05_TEST_STRATEGY.md) | テスト戦略（Testcontainers） | エンジニア |
| [06_IMPLEMENTATION_ROADMAP.md](./06_IMPLEMENTATION_ROADMAP.md) | 実装ロードマップ | 全員 |

---

## プロジェクトの本質

> **「自分の好きな本で、AIと会話しながら言語を学ぶ」**

- 既存アプリのカリキュラムに縛られない
- 手持ちの教材をそのまま使える
- AIが対応する言語なら何でも学習可能
- 会話形式で楽しく学習

---

## 技術スタック

| 層 | 技術 |
|----|------|
| Frontend (Web) | Next.js 14+, TypeScript, TailwindCSS, ShadCN/UI |
| Frontend (Mobile) | Flutter 3.0+ |
| Backend | Go 1.21+, Chi Router |
| Database | PostgreSQL 15+, Redis 7+ |
| Testing | Testcontainers, Vitest, Playwright |
| External APIs | Google Vision, Google TTS, OpenAI (GPT, Whisper) |

---

## 読む順序

### 初めての人
1. `00_PROJECT_VISION.md` - 何を作るのか理解
2. `01_CORE_USECASES.md` - どう使われるのか理解
3. `06_IMPLEMENTATION_ROADMAP.md` - どう作るのか理解

### 実装する人
1. `02_SYSTEM_ARCHITECTURE.md` - 全体構成を理解
2. `03_DATABASE_SCHEMA.md` - データ構造を理解
3. `04_API_SPECIFICATION.md` - API仕様を確認
4. `05_TEST_STRATEGY.md` - テスト方法を確認
5. `06_IMPLEMENTATION_ROADMAP.md` - タスクを確認

---

## 旧ドキュメントとの関係

```
docs/
├── v2/                          ← このディレクトリ（SSOT）
│   ├── 00_PROJECT_VISION.md
│   ├── 01_CORE_USECASES.md
│   ├── 02_SYSTEM_ARCHITECTURE.md
│   ├── 03_DATABASE_SCHEMA.md
│   ├── 04_API_SPECIFICATION.md
│   ├── 05_TEST_STRATEGY.md
│   └── 06_IMPLEMENTATION_ROADMAP.md
├── requirements_definition.md    ← 旧（参考のみ）
├── ui_ux_design_document.md      ← 旧（参考のみ）
├── featureRDs/                   ← 旧（参考のみ）
└── ...
```

旧ドキュメントは参考として残しますが、矛盾がある場合はv2を優先してください。

---

## 更新履歴

| 日付 | 内容 |
|------|------|
| 2024-12-04 | v2要件定義書作成 |

---

## 質問・フィードバック

要件定義に関する質問や提案は、GitHub Issuesで受け付けています。
