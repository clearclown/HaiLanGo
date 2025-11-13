# API統合提案書 - HaiLanGoプロジェクト

## 概要

本ドキュメントは、HaiLanGoプロジェクトに統合可能な外部API、ライブラリ、GitHubリポジトリを包括的に調査し、活用方法を提案するものです。

## 1. 音声・対話API

### 1.1 OpenAI Realtime API / gpt-realtime ⭐ NEW

**概要**: 本番環境対応のリアルタイム音声エージェントAPI

**参考リンク**:
- [OpenAI公式ブログ](https://openai.com/index/introducing-gpt-realtime/)
- [Zenn記事](https://zenn.dev/dxclab/articles/e78dac9d46e17f)
- [Note記事](https://note.com/npaka/n/nad9118c22c48)

**主な機能**:
- **リアルタイム音声対話**: 低遅延の音声エージェント構築
- **gpt-realtimeモデル**: 高品質な音声合成（音質、インテリジェンス、指示追従、Function Calling）
- **画像入力対応**: 会話中に画像を共有可能
- **SIPサポート**: 公衆電話網、PBXシステムへの接続
- **リモートMCPサーバー**: 追加ツールやコンテキストへのアクセス

**HaiLanGoでの活用案**:

#### 1.1.1 リアルタイム発音練習モード
```go
// backend/pkg/realtime/realtime.go
package realtime

// RealtimePronunciationPractice はリアルタイム発音練習を提供
type RealtimePronunciationPractice struct {
    client *openai.RealtimeClient
}

func (r *RealtimePronunciationPractice) StartSession(ctx context.Context, userID string, targetLanguage string) (*Session, error) {
    // gpt-realtimeを使用したリアルタイム発音練習セッション
    // ユーザーの発音をリアルタイムで評価・フィードバック
}
```

**活用シーン**:
- フレーズ練習時のリアルタイム発音評価
- 会話ロールプレイの実現
- 自然な対話形式での学習

**メリット**:
- 低遅延のリアルタイム対話
- 高品質な音声合成（20%コスト削減）
- 多言語対応（日本語、中国語、スペイン語など）
- 文の途中での言語切り替え対応

**価格**: 
- 入力: $32/100万トークン（キャッシュ済み: $0.40）
- 出力: $64/100万トークン

**統合優先度**: ⭐⭐⭐⭐⭐（高）

---

### 1.2 Azure OpenAI Service

**概要**: OpenAIモデルをAzure上で利用

**活用案**:
- エンタープライズ向けのセキュリティ要件対応
- EU Data Residency対応
- 既存Azureインフラとの統合

**統合優先度**: ⭐⭐⭐（中）

---

### 1.3 Deepgram

**概要**: リアルタイム音声認識API

**特徴**:
- 低遅延（<100ms）
- 多言語対応
- カスタムモデル対応

**活用案**:
- STTの代替・補完
- リアルタイム音声認識の高速化

**統合優先度**: ⭐⭐⭐（中）

---

### 1.4 AssemblyAI

**概要**: 音声認識・理解API

**特徴**:
- 自動文字起こし
- 感情分析
- トピック検出

**活用案**:
- 発音評価の補完
- 学習コンテンツの自動分析

**統合優先度**: ⭐⭐（低）

---

## 2. PDF処理・OCR

### 2.1 MarkPDFdown ⭐ NEW

**GitHub**: https://github.com/MarkPDFdown/markpdfdown

**概要**: 大規模言語モデルの視覚認識を活用したPDF→Markdown変換ツール

**特徴**:
- 高品質なPDF変換
- 複雑なレイアウト対応
- Markdown形式での出力

**HaiLanGoでの活用案**:

#### 2.1.1 PDF前処理パイプライン
```go
// backend/pkg/pdf/markpdfdown.go
package pdf

import "github.com/MarkPDFdown/markpdfdown"

// MarkPDFdownProcessor はMarkPDFdownを使用したPDF処理
type MarkPDFdownProcessor struct {
    client *markpdfdown.Client
}

func (m *MarkPDFdownProcessor) ConvertToMarkdown(ctx context.Context, pdfData []byte) (*MarkdownResult, error) {
    // PDFをMarkdownに変換
    // OCR処理の前処理として使用
    // 構造化されたテキストを取得
}
```

**活用シーン**:
- OCR処理の前処理として使用
- 複雑なレイアウトのPDFの構造化
- 表や図表を含むPDFの処理

**メリット**:
- OCR精度の向上
- 構造化されたデータの取得
- Markdown形式での保存・再利用

**統合優先度**: ⭐⭐⭐⭐（高）

---

### 2.2 Unstructured.io

**概要**: ドキュメント構造化API

**特徴**:
- PDF、Word、HTMLなどの構造化
- テーブル抽出
- メタデータ抽出

**活用案**:
- PDF前処理の補完
- 書籍構造の自動解析

**統合優先度**: ⭐⭐⭐（中）

---

### 2.3 Adobe PDF Services API

**概要**: Adobe公式PDF処理API

**特徴**:
- PDF操作（結合、分割、変換）
- OCR機能
- フォーム処理

**活用案**:
- PDF前処理
- 高品質なOCR処理

**統合優先度**: ⭐⭐（低）- コストが高い

---

### 2.4 pdf.js (Mozilla)

**概要**: オープンソースPDFレンダリングライブラリ

**特徴**:
- ブラウザ内でのPDF表示
- テキスト抽出
- 無料・オープンソース

**活用案**:
- フロントエンドでのPDF表示
- クライアントサイドでのテキスト抽出

**統合優先度**: ⭐⭐⭐⭐（高）

---

## 3. 言語学習・辞書API

### 3.1 DeepL API

**概要**: 高品質な翻訳API

**特徴**:
- 自然な翻訳
- 多言語対応
- コンテキスト理解

**活用案**:
- 学習先言語→母国語の翻訳
- フレーズの自然な翻訳
- 例文生成

**統合優先度**: ⭐⭐⭐⭐（高）

---

### 3.2 Google Translate API

**概要**: Google翻訳API

**特徴**:
- 100+言語対応
- 低コスト
- 音声合成統合

**活用案**:
- 翻訳の補完
- 多言語対応の拡大

**統合優先度**: ⭐⭐⭐（中）

---

### 3.3 Lingvanex API

**概要**: 多言語翻訳・音声API

**特徴**:
- 100+言語対応
- 音声合成統合
- 低コスト

**活用案**:
- マイナー言語のサポート
- TTSの補完

**統合優先度**: ⭐⭐⭐（中）

---

### 3.4 Wordnik API

**概要**: 英語辞書API

**特徴**:
- 豊富な語彙情報
- 例文
- 類義語・対義語

**活用案**:
- 英語学習の補完
- 単語帳機能の拡張

**統合優先度**: ⭐⭐（低）

---

## 4. AI・機械学習API

### 4.1 Anthropic Claude API

**概要**: Claude AI API

**特徴**:
- 長文理解
- コード生成
- 多言語対応

**活用案**:
- 学習コンテンツの自動生成
- 解説の自動生成
- 会話パターンの抽出

**統合優先度**: ⭐⭐⭐⭐（高）

---

### 4.2 Google Gemini API

**概要**: Google Gemini AI API

**特徴**:
- マルチモーダル（画像・音声・テキスト）
- 低コスト
- 高速レスポンス

**活用案**:
- OCR処理の補完
- 画像理解
- 学習コンテンツ分析

**統合優先度**: ⭐⭐⭐（中）

---

### 4.3 Cohere API

**概要**: NLP API

**特徴**:
- テキスト分類
- 埋め込み生成
- 要約生成

**活用案**:
- 学習コンテンツの要約
- 類似フレーズ検索
- 単語の関連性分析

**統合優先度**: ⭐⭐（低）

---

## 5. ストレージ・CDN

### 5.1 Cloudflare R2

**概要**: S3互換オブジェクトストレージ

**特徴**:
- S3互換API
- エグレス料金なし
- 低コスト

**活用案**:
- 音声ファイルの保存
- PDFファイルの保存
- CDNとしての利用

**統合優先度**: ⭐⭐⭐⭐（高）

---

### 5.2 Backblaze B2

**概要**: 低コストオブジェクトストレージ

**特徴**:
- S3互換API
- 低コスト
- 高速アクセス

**活用案**:
- 大容量ファイルの保存
- バックアップ

**統合優先度**: ⭐⭐⭐（中）

---

## 6. 認証・セキュリティ

### 6.1 Auth0

**概要**: 認証・認可サービス

**特徴**:
- 多様な認証方式
- ソーシャルログイン
- セキュリティ機能

**活用案**:
- OAuth認証の拡張
- 多要素認証（MFA）
- セキュリティ強化

**統合優先度**: ⭐⭐（低）- 既存実装で十分

---

### 6.2 Clerk

**概要**: 認証・ユーザー管理サービス

**特徴**:
- 簡単な統合
- ユーザー管理UI
- セキュリティ機能

**活用案**:
- 認証システムの簡素化
- ユーザー管理の自動化

**統合優先度**: ⭐⭐（低）

---

## 7. 分析・モニタリング

### 7.1 PostHog

**概要**: オープンソース分析プラットフォーム

**特徴**:
- ユーザー行動分析
- 機能フラグ
- A/Bテスト

**活用案**:
- 学習行動の分析
- 機能の効果測定
- ユーザー体験の改善

**統合優先度**: ⭐⭐⭐（中）

---

### 7.2 Sentry

**概要**: エラー監視・パフォーマンス監視

**特徴**:
- エラー追跡
- パフォーマンス監視
- リリース追跡

**活用案**:
- エラーの早期発見
- パフォーマンス最適化
- ユーザー体験の改善

**統合優先度**: ⭐⭐⭐⭐（高）

---

## 8. 通知・コミュニケーション

### 8.1 SendGrid

**概要**: メール送信サービス

**特徴**:
- 高配信率
- テンプレート機能
- 分析機能

**活用案**:
- パスワードリセット
- 学習リマインダー
- 通知メール

**統合優先度**: ⭐⭐⭐（中）

---

### 8.2 Twilio

**概要**: コミュニケーションAPI

**特徴**:
- SMS送信
- 音声通話
- ビデオ通話

**活用案**:
- SMS通知
- 音声通話での学習サポート（将来）

**統合優先度**: ⭐⭐（低）

---

## 9. GitHubリポジトリ活用

### 9.1 MarkPDFdown

**GitHub**: https://github.com/MarkPDFdown/markpdfdown

**統合方法**:
```go
// backend/go.mod に追加
require github.com/MarkPDFdown/markpdfdown v1.0.0

// 使用例
import "github.com/MarkPDFdown/markpdfdown"

func ProcessPDF(pdfData []byte) (string, error) {
    converter := markpdfdown.NewConverter()
    markdown, err := converter.Convert(pdfData)
    return markdown, err
}
```

**統合優先度**: ⭐⭐⭐⭐（高）

---

### 9.2 その他の有用なリポジトリ

#### 9.2.1 PDF処理
- **pdfcpu**: Go製PDF処理ライブラリ
- **unidoc/unipdf**: Go製PDF処理ライブラリ
- **pdf-lib**: JavaScript製PDF処理ライブラリ

#### 9.2.2 OCR
- **tesseract-ocr**: オープンソースOCRエンジン
- **gocv**: Go用OpenCVバインディング

#### 9.2.3 音声処理
- **whisper.cpp**: C++実装のWhisper
- **vosk**: オフライン音声認識

#### 9.2.4 言語学習
- **anki**: 間隔反復学習アプリ（参考）
- **memrise**: 言語学習アプリ（参考）

---

## 10. 統合優先度マトリックス

| API/ツール | 優先度 | 実装時期 | コスト | 技術的難易度 |
|-----------|--------|---------|--------|------------|
| OpenAI Realtime API | ⭐⭐⭐⭐⭐ | Phase 2 | 中 | 中 |
| MarkPDFdown | ⭐⭐⭐⭐ | Phase 1 | 無料 | 低 |
| DeepL API | ⭐⭐⭐⭐ | Phase 2 | 中 | 低 |
| Cloudflare R2 | ⭐⭐⭐⭐ | Phase 2 | 低 | 低 |
| Sentry | ⭐⭐⭐⭐ | Phase 1 | 低 | 低 |
| Anthropic Claude | ⭐⭐⭐⭐ | Phase 3 | 中 | 低 |
| pdf.js | ⭐⭐⭐⭐ | Phase 1 | 無料 | 低 |
| Deepgram | ⭐⭐⭐ | Phase 3 | 中 | 中 |
| PostHog | ⭐⭐⭐ | Phase 3 | 低 | 中 |

---

## 11. 実装ロードマップ

### Phase 1: MVP（即座に統合可能）
- [ ] **MarkPDFdown**: PDF前処理パイプライン
- [ ] **pdf.js**: フロントエンドPDF表示
- [ ] **Sentry**: エラー監視

### Phase 2: コア機能拡張
- [ ] **OpenAI Realtime API**: リアルタイム発音練習
- [ ] **DeepL API**: 高品質翻訳
- [ ] **Cloudflare R2**: ストレージ最適化

### Phase 3: 高度な機能
- [ ] **Anthropic Claude**: 学習コンテンツ生成
- [ ] **Deepgram**: 高速音声認識
- [ ] **PostHog**: ユーザー行動分析

---

## 12. コスト試算

### 月間1,000ユーザー想定

| API | 月間コスト（概算） |
|-----|------------------|
| OpenAI Realtime API | $50-100 |
| DeepL API | $30-50 |
| Cloudflare R2 | $10-20 |
| Sentry | $26（無料プランあり） |
| MarkPDFdown | $0（オープンソース） |
| **合計** | **$116-196** |

---

## 13. 実装例

### 13.1 OpenAI Realtime API統合

```go
// backend/pkg/realtime/openai_realtime.go
package realtime

import (
    "context"
    "github.com/openai/openai-go"
)

type OpenAIRealtimeClient struct {
    client *openai.Client
}

func NewOpenAIRealtimeClient(apiKey string) *OpenAIRealtimeClient {
    return &OpenAIRealtimeClient{
        client: openai.NewClient().WithAPIKey(apiKey),
    }
}

func (c *OpenAIRealtimeClient) StartPronunciationSession(
    ctx context.Context,
    userID string,
    targetLanguage string,
    nativeLanguage string,
) (*PronunciationSession, error) {
    // gpt-realtimeを使用した発音練習セッション
    session := &PronunciationSession{
        UserID: userID,
        TargetLanguage: targetLanguage,
        NativeLanguage: nativeLanguage,
    }
    
    // リアルタイムAPIセッション開始
    // 発音評価・フィードバックをリアルタイムで提供
    
    return session, nil
}
```

### 13.2 MarkPDFdown統合

```go
// backend/pkg/pdf/markpdfdown.go
package pdf

import (
    "context"
    "github.com/MarkPDFdown/markpdfdown"
)

type MarkPDFdownProcessor struct {
    converter *markpdfdown.Converter
}

func NewMarkPDFdownProcessor() *MarkPDFdownProcessor {
    return &MarkPDFdownProcessor{
        converter: markpdfdown.NewConverter(),
    }
}

func (m *MarkPDFdownProcessor) PreprocessPDF(
    ctx context.Context,
    pdfData []byte,
) (*PreprocessedPDF, error) {
    // PDFをMarkdownに変換
    markdown, err := m.converter.Convert(pdfData)
    if err != nil {
        return nil, err
    }
    
    // 構造化されたデータを返す
    return &PreprocessedPDF{
        Markdown: markdown,
        Structured: true,
    }, nil
}
```

---

## 14. まとめ

### 推奨される統合順序

1. **即座に統合**（Phase 1）:
   - MarkPDFdown（PDF前処理）
   - pdf.js（フロントエンドPDF表示）
   - Sentry（エラー監視）

2. **早期統合**（Phase 2）:
   - OpenAI Realtime API（リアルタイム発音練習）
   - DeepL API（高品質翻訳）
   - Cloudflare R2（ストレージ最適化）

3. **将来統合**（Phase 3）:
   - Anthropic Claude（コンテンツ生成）
   - Deepgram（高速音声認識）
   - PostHog（ユーザー分析）

### 重要なポイント

- **APIキーなしでも開発可能**: モックシステムを活用
- **コスト管理**: キャッシュとレート制限を適切に設定
- **段階的統合**: 優先度に基づいて段階的に統合
- **フォールバック**: 複数のAPIを用意して可用性を確保

---

## 参考リンク

- [OpenAI Realtime API](https://openai.com/index/introducing-gpt-realtime/)
- [MarkPDFdown GitHub](https://github.com/MarkPDFdown/markpdfdown)
- [DeepL API](https://www.deepl.com/docs-api)
- [Cloudflare R2](https://developers.cloudflare.com/r2/)
- [Sentry](https://sentry.io/)

