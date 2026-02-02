# HaiLanGo

**AI言語学習プラットフォーム** - 自分の教科書をAIで学習可能にする

[![Rust](https://img.shields.io/badge/Rust-2024-orange?style=flat&logo=rust)](https://www.rust-lang.org/)
[![Reinhardt](https://img.shields.io/badge/Reinhardt-Full--Stack-blue?style=flat)](https://github.com/kent8192/reinhardt-web)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## 概要

HaiLanGoは、既存の言語学習本とAI技術を組み合わせ、個人に最適化された能動的な言語学習環境を提供するプラットフォームです。

### 主要な価値提案

- **手持ちの本を活用**: 自分の教科書をアップロードして学習
- **AI教師による24/7の個別指導**: いつでもどこでも学習可能
- **スピーキング重視**: 発音練習・会話練習に特化
- **マイナー言語対応**: ペルシャ語、ヘブライ語など幅広い言語をサポート

## 技術スタック

[Reinhardt](https://github.com/kent8192/reinhardt-web) (Rust Full-Stack Framework) で全て構築:

| コンポーネント | クレート | 説明 |
|---------------|---------|------|
| ORM | `reinhardt-db` | SeaQuery + sqlx (PostgreSQL対応) |
| Frontend | `reinhardt-pages` | WASM + SSR リアクティブUI |
| REST API | `reinhardt-rest` | ViewSets, Serializers |
| Auth | `reinhardt-auth` | JWT, OAuth, Session認証 |
| Admin | `reinhardt-admin` | 自動生成管理画面 |
| WebSocket | `reinhardt-websockets` | リアルタイム通信 |
| Cache | Redis | キャッシュ・セッション管理 |

## 主要機能

1. **書籍のデジタル化** - PDF/画像をOCRでテキスト化
2. **TTS音声合成** - 多言語対応のネイティブ発音
3. **STT発音評価** - リアルタイムスコアリング
4. **教師モード** - ハンズフリー自動学習
5. **間隔反復学習 (SRS)** - SM-2アルゴリズムによる最適な復習
6. **オフライン対応** - PWAによる音声キャッシュ

## 外部API

| API | 用途 |
|-----|------|
| Google Vision / Azure CV | OCR (文字認識) |
| Google TTS / Azure Speech | 音声合成 |
| OpenAI Whisper | 音声認識・発音評価 |
| Anthropic Claude | LLM (解説生成等) |
| Stripe | 決済処理 |

## ドキュメント

- [要件定義書](docs/requirements_definition.md) - プロジェクトの詳細仕様
- [CLAUDE.md](CLAUDE.md) - AI開発アシスタント向けガイドライン
- [llms.txt](llms.txt) - LLM向けプロジェクト概要

## 開発状況

🚧 **開発中** - Reinhardt (Rust) フルスタックフレームワークで構築中

## ライセンス

MIT License
