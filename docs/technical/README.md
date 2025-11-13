# 技術仕様書ディレクトリ

このディレクトリには、HaiLanGoプロジェクトの各機能の詳細な技術仕様とアーキテクチャドキュメントが含まれています。

## 📑 技術仕様書一覧

### コア機能

| ドキュメント | 説明 | 実装状況 |
|-------------|------|----------|
| [teacher_mode.md](teacher_mode.md) | 教師モード（自動学習モード）の技術仕様 | ✅ 完了 |
| [websocket.md](websocket.md) | WebSocketリアルタイム通知の実装詳細 | ✅ 完了 |
| [ocr_implementation.md](ocr_implementation.md) | OCR処理機能の実装サマリー | ✅ 完了 |

## 📋 各ドキュメントの概要

### 1. [教師モード技術仕様](teacher_mode.md)

**概要**: ボタン一つで教師が授業をするように連続自動再生する機能

**主な内容**:
- アーキテクチャ設計
- データモデル（TeacherModeSettings, AudioSegment, PageAudio）
- バックグラウンド再生の実装
- オフライン対応
- APIエンドポイント設計
- パフォーマンス最適化
- コスト試算

**実装場所**:
- Backend: `backend/internal/service/teacher-mode/`
- Frontend: `frontend/web/components/teacher-mode/`

---

### 2. [WebSocket実装仕様](websocket.md)

**概要**: リアルタイム通知システムの実装詳細

**主な内容**:
- WebSocketアーキテクチャ
- 通知タイプ（OCR進捗、TTS進捗、学習進捗等）
- Hub パターンによる接続管理
- クライアント・サーバー間の通信プロトコル
- 再接続ロジック
- テスト戦略

**実装場所**:
- Backend: `backend/internal/api/websocket/`
- Frontend: `frontend/web/hooks/useWebSocket.ts`

---

### 3. [OCR実装サマリー](ocr_implementation.md)

**概要**: OCR処理機能の実装の詳細サマリー

**主な内容**:
- OCRクライアントの設計（インターフェース、プロバイダー）
- キャッシュ層の実装
- OCR処理サービス
- モックシステム
- テスト戦略
- パフォーマンス最適化

**実装場所**:
- Backend: `backend/internal/service/ocr/`, `backend/pkg/ocr/`
- Frontend: `frontend/web/components/ocr-editor/`

---

## 🔗 関連ドキュメント

### プロジェクト全体
- [要件定義書](../requirements_definition.md) - プロジェクトの全体像と機能要件
- [UI/UX設計書](../ui_ux_design_document.md) - 画面設計とワイヤーフレーム
- [実装状況サマリー](../IMPLEMENTATION_STATUS.md) - 全機能の実装進捗

### 実装関連
- [モック構築戦略](../mocking_strategy.md) - APIキーなしでもテスト可能な仕組み
- [API統合提案書](../api_integration_proposal.md) - 統合可能な外部API・ツール

### 機能別詳細
- [機能実装RD](../featureRDs/README.md) - 各機能の要件定義書
- [実装済み機能](../featureRDs/archives/README.md) - アーカイブされた実装完了文書

---

## 📖 ドキュメントの使い方

### 新機能の実装時
1. [要件定義書](../requirements_definition.md)で全体像を把握
2. [機能実装RD](../featureRDs/)で詳細要件を確認
3. 該当する技術仕様書で実装詳細を確認
4. [モック構築戦略](../mocking_strategy.md)でテスト方法を確認

### 既存機能の理解時
1. [実装状況サマリー](../IMPLEMENTATION_STATUS.md)で実装状況を確認
2. 該当する技術仕様書で詳細を確認
3. コードベース内の実装を確認

### バグ修正・機能改善時
1. 該当する技術仕様書でアーキテクチャを理解
2. 実装箇所を特定
3. テスト戦略を確認して修正・改善

---

## 🔄 更新履歴

| 日付 | 内容 |
|------|------|
| 2025-11-13 | technical/ディレクトリ作成、技術仕様書を整理 |

---

**Note**: 各機能の要件定義は `featureRDs/archives/` に、詳細な技術仕様はこのディレクトリに保管されています。
