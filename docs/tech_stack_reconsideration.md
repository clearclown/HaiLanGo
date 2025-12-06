# 技術スタック再検討 - マイナー言語＆ドメイン特化アプローチ

## HaiLanGoの独自価値

### 競合との差別化ポイント

1. **マイナー言語・マイナー言語話者のサポート**
   - クルド語、アムハラ語、チベット語など大手アプリが対応しない言語
   - 日本語話者がペルシャ語を学ぶ、ヘブライ語話者がロシア語を学ぶなどのニッチなペア
   - 少数言語コミュニティの学習ニーズに対応

2. **ドメイン特化学習**
   - 政治用語のみ、SNSスラングのみ、クルアーンのみ
   - ビジネス交渉、医療用語、法律用語など専門分野
   - ユーザーが持っている特定の教材に最適化

### 従来の「12言語サポート」アプローチの問題

| 従来アプローチ | HaiLanGoアプローチ |
|--------------|------------------|
| 事前に対応言語を固定 | **任意の言語**をAPIが許す限り対応 |
| 汎用的な教材を提供 | **ユーザーの教材**を活用 |
| 一般的な語彙・表現 | **特定ドメイン**に特化可能 |
| 主要言語優先 | **マイナー言語**も同等に対応 |

---

## 外部API再検討

### 1. OCR（テキスト認識）

#### 現状
- Google Vision API（プライマリ）
- Azure Computer Vision（セカンダリ）
- Tesseract（オープンソース）

#### マイナー言語観点での評価

| API | マイナー言語対応 | 非ラテン文字 | RTL対応 | コスト |
|-----|----------------|-------------|---------|--------|
| **Google Vision** | ⭐⭐⭐⭐⭐ | 優秀 | ✅ | 中 |
| Azure Vision | ⭐⭐⭐⭐ | 良好 | ✅ | 中 |
| Tesseract | ⭐⭐⭐ | 要追加学習 | ✅ | 無料 |
| **Amazon Textract** | ⭐⭐⭐ | 限定的 | ✅ | 高 |

#### 推奨構成
```
Primary:   Google Vision API（200+言語、最も広範囲）
Secondary: Azure Vision（フォールバック）
Fallback:  Tesseract + 言語パック（オフライン/コスト削減）
```

**Google Visionが最適な理由**:
- クルド語、アムハラ語、チベット語などの稀少言語も対応
- 縦書き日本語、アラビア文字、ヘブライ文字の高精度認識
- 複雑なレイアウト（表、ルビ、段組み）への対応

---

### 2. TTS（音声合成）

#### マイナー言語観点での評価

| API | 対応言語数 | マイナー言語 | 音質 | コスト |
|-----|-----------|-------------|------|--------|
| **Google Cloud TTS** | 50+ | ⭐⭐⭐⭐ | 高 | 中 |
| Amazon Polly | 30+ | ⭐⭐⭐ | 高 | 中 |
| Azure TTS | 140+ | ⭐⭐⭐⭐⭐ | 高 | 中 |
| **ElevenLabs** | 32+ | ⭐⭐⭐ | 最高 | 高 |
| Coqui TTS | 多言語 | ⭐⭐⭐⭐ | 中 | 無料 |

#### 推奨構成

```
Primary:   Azure TTS（140言語、最大カバレッジ）
Secondary: Google Cloud TTS（安定性）
Premium:   ElevenLabs（高品質オプション）
Offline:   Coqui TTS / edge-tts（オープンソース）
```

**Azure TTSへの変更を推奨する理由**:
- 140言語以上（Google: 50+）
- クルド語、アムハラ語など多くのマイナー言語に対応
- Neural Voicesで高品質
- カスタム音声モデル作成可能

#### 新規追加候補: edge-tts

```go
// オープンソースでMicrosoft Edge TTSを利用
// 完全無料、50+言語対応
// https://github.com/rany2/edge-tts
```

---

### 3. STT（音声認識・発音評価）

#### マイナー言語観点での評価

| API | 対応言語数 | 発音評価 | リアルタイム | コスト |
|-----|-----------|---------|-------------|--------|
| **Whisper (OpenAI)** | 99言語 | ❌ | ❌ | 中 |
| **whisper.cpp** | 99言語 | ❌ | ✅ | 無料 |
| Google Cloud STT | 125+ | ✅ | ✅ | 中 |
| Azure Speech | 100+ | ✅ | ✅ | 中 |
| **Deepgram** | 36言語 | ❌ | ✅ | 中 |

#### 推奨構成

```
Recognition:  OpenAI Whisper API / whisper.cpp（99言語、マイナー言語最強）
Pronunciation: Azure Speech（発音評価機能内蔵）
Fallback:     Google Cloud STT（安定性）
```

**Whisperを推奨する理由**:
- **99言語対応**（業界最多）
- クルド語、アムハラ語、チベット語など稀少言語も高精度
- whisper.cppでローカル実行可能（オフライン対応）
- 方言・訛りへの対応が優秀

#### 発音評価のアプローチ変更

従来: Google/Azure STTの発音スコア機能
新規: **Whisper + LLM分析**

```go
// Whisperで音声をテキスト化
transcription := whisper.Transcribe(audioData, targetLanguage)

// LLMで発音評価（ドメイン特化プロンプト）
prompt := fmt.Sprintf(`
言語: %s
期待テキスト: %s
認識テキスト: %s

発音の精度を0-100で評価し、具体的な改善点を指摘してください。
特に、この言語特有の発音ポイントに注目してください。
`, targetLanguage, expectedText, transcription)

evaluation := llm.Analyze(prompt)
```

**利点**:
- マイナー言語でも動作（既存の発音評価APIは主要言語のみ）
- ドメイン特化の評価が可能（宗教用語の正確な発音など）
- カスタマイズ可能

---

### 4. LLM（コンテンツ生成・分析）

#### ドメイン特化観点での評価

| API | 多言語能力 | ドメイン特化 | コスト | 長文対応 |
|-----|-----------|-------------|--------|---------|
| **Claude 3.5 Sonnet** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 中 | 200K tokens |
| GPT-4o | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 高 | 128K tokens |
| **Gemini 1.5 Pro** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | 低 | 1M tokens |
| Llama 3 | ⭐⭐⭐ | ⭐⭐⭐ | 無料 | 8K tokens |

#### 推奨構成

```
Primary:   Claude 3.5 Sonnet（多言語最強、指示追従性高）
Secondary: Gemini 1.5 Pro（長文対応、低コスト）
Fallback:  GPT-4o（安定性）
Local:     Llama 3 / Mistral（オフライン/プライバシー）
```

**Claude 3.5を主力にする理由**:
- マイナー言語でも高品質な生成・分析
- ドメイン特化プロンプトへの追従性が高い
- 200Kトークンで長い教材も一度に処理可能

#### ドメイン特化の実装例

```go
// ドメイン特化プロンプトテンプレート
type DomainPrompt struct {
    Domain      string // "quran", "politics", "sns_slang", etc.
    Language    string
    SourceText  string
}

func (d *DomainPrompt) GenerateExplanation() string {
    templates := map[string]string{
        "quran": `
あなたはクルアーン（コーラン）の専門家です。
言語: %s
テキスト: %s

このアーヤ（節）について：
1. 原語（アラビア語）の正確な発音
2. 文法構造と語彙の解説
3. イスラム教における意味と解釈
4. 日常会話との違い
を説明してください。`,
        "politics": `
あなたは政治用語の専門家です。
言語: %s
テキスト: %s

この政治用語について：
1. 正式な定義
2. 歴史的背景
3. 現代の用法
4. 関連する重要語彙
を説明してください。`,
        "sns_slang": `
あなたはSNSスラング・ネット用語の専門家です。
言語: %s
テキスト: %s

このスラング/ネット用語について：
1. 意味と使用場面
2. 由来・語源
3. 使用上の注意（フォーマル/カジュアル）
4. 類似表現
を説明してください。`,
    }
    return fmt.Sprintf(templates[d.Domain], d.Language, d.SourceText)
}
```

---

### 5. 翻訳API

#### マイナー言語観点での評価

| API | 対応言語数 | マイナー言語品質 | コスト |
|-----|-----------|----------------|--------|
| **Google Translate** | 133言語 | ⭐⭐⭐⭐ | 低 |
| DeepL | 31言語 | ⭐⭐⭐⭐⭐ | 中 |
| Azure Translator | 100+言語 | ⭐⭐⭐⭐ | 中 |
| **LibreTranslate** | 30+言語 | ⭐⭐⭐ | 無料 |

#### 推奨構成

```
Primary:   Google Translate API（133言語、マイナー言語カバレッジ最大）
Quality:   DeepL API（主要言語の高品質翻訳）
Fallback:  Azure Translator
Offline:   LibreTranslate / Argos Translate
```

**DeepLを補助に留める理由**:
- 31言語のみ（マイナー言語非対応）
- 主要言語ペアでは最高品質だが、HaiLanGoの価値提案と一致しない
- 主要言語ユーザー向けのプレミアム機能として位置付け

---

### 6. リアルタイム音声対話

#### 比較評価

| API | 多言語 | 遅延 | コスト | 成熟度 |
|-----|--------|-----|--------|--------|
| **OpenAI Realtime API** | ⭐⭐⭐⭐ | 低 | 高 | ⭐⭐⭐⭐⭐ |
| VibeVoice-Realtime | 英語のみ | 低 | 無料 | ⭐⭐ |
| Deepgram + TTS | ⭐⭐⭐ | 中 | 中 | ⭐⭐⭐⭐ |
| **Whisper + TTS** | ⭐⭐⭐⭐⭐ | 中 | 中 | ⭐⭐⭐⭐ |

#### 推奨構成

```
Premium:  OpenAI Realtime API（プレミアム機能として）
Standard: Whisper + Azure TTS + Claude（パイプライン構成）
Offline:  whisper.cpp + Coqui TTS + Llama（ローカル実行）
```

**VibeVoiceを採用しない理由**:
- **英語のみ**対応 → HaiLanGoの価値提案と完全に矛盾
- マイナー言語サポートという差別化ポイントを活かせない

**パイプライン構成の利点**:
- 各コンポーネントを最適化可能
- マイナー言語それぞれに最適なAPIを選択可能
- コスト最適化が容易

---

## 新しい推奨アーキテクチャ

### API選択マトリックス

```
┌─────────────────────────────────────────────────────────────────┐
│                        言語カテゴリ                              │
├────────────────┬────────────────┬────────────────┬──────────────┤
│     主要言語     │   中程度言語     │   マイナー言語   │   実験的言語  │
│  (en, ja, zh等)  │  (fa, he, tr等) │  (ku, am, bo等) │  (その他)    │
├────────────────┼────────────────┼────────────────┼──────────────┤
│ OCR            │                │                │              │
│ Google Vision  │ Google Vision  │ Google Vision  │ Google Vision│
├────────────────┼────────────────┼────────────────┼──────────────┤
│ TTS            │                │                │              │
│ Azure + ElevenLabs │ Azure TTS  │ Azure TTS     │ edge-tts     │
├────────────────┼────────────────┼────────────────┼──────────────┤
│ STT            │                │                │              │
│ Azure Speech   │ Whisper        │ Whisper       │ Whisper      │
├────────────────┼────────────────┼────────────────┼──────────────┤
│ 発音評価        │                │                │              │
│ Azure Speech   │ Whisper+LLM    │ Whisper+LLM   │ Whisper+LLM  │
├────────────────┼────────────────┼────────────────┼──────────────┤
│ 翻訳           │                │                │              │
│ DeepL          │ Google Translate│ Google Translate│ Google Translate│
├────────────────┼────────────────┼────────────────┼──────────────┤
│ LLM            │                │                │              │
│ Claude 3.5     │ Claude 3.5     │ Claude 3.5    │ Claude 3.5   │
└────────────────┴────────────────┴────────────────┴──────────────┘
```

### 環境変数設計

```bash
# .env.example（更新版）

# ===========================================
# OCR Configuration
# ===========================================
OCR_PROVIDER=google_vision  # google_vision, azure_vision, tesseract
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json

# ===========================================
# TTS Configuration（Azure優先に変更）
# ===========================================
TTS_PROVIDER=azure          # azure, google, elevenlabs, edge
AZURE_SPEECH_KEY=your_key
AZURE_SPEECH_REGION=eastus

# Premium TTS (optional)
ELEVENLABS_API_KEY=your_key

# ===========================================
# STT Configuration（Whisper追加）
# ===========================================
STT_PROVIDER=whisper        # whisper, azure, google
OPENAI_API_KEY=your_key     # for Whisper API

# Pronunciation evaluation method
PRONUNCIATION_EVAL_METHOD=whisper_llm  # whisper_llm, azure_native

# ===========================================
# LLM Configuration（Claude優先に変更）
# ===========================================
LLM_PROVIDER=claude         # claude, openai, gemini
ANTHROPIC_API_KEY=your_key

# Fallback LLM
FALLBACK_LLM_PROVIDER=gemini
GOOGLE_AI_API_KEY=your_key

# ===========================================
# Translation Configuration
# ===========================================
TRANSLATION_PROVIDER=google   # google, deepl, azure
DEEPL_API_KEY=your_key        # Premium translation (optional)

# ===========================================
# Domain-Specific Settings
# ===========================================
ENABLE_DOMAIN_PROMPTS=true
DOMAIN_PROMPT_DIR=./prompts/domains
```

---

## 実装ロードマップ

### Phase 1: 基盤整備（即時）
- [x] Language Registry（任意言語対応）- 既に実装済み
- [ ] Azure TTS統合（140言語対応）
- [ ] Whisper STT統合（99言語対応）
- [ ] Claude LLM統合

### Phase 2: ドメイン特化（2週間）
- [ ] ドメインプロンプトシステム実装
- [ ] Whisper + LLM発音評価パイプライン
- [ ] ドメイン別テンプレートライブラリ

### Phase 3: 最適化（1ヶ月）
- [ ] 言語別API自動ルーティング
- [ ] オフライン対応（whisper.cpp, edge-tts, Llama）
- [ ] コスト最適化ダッシュボード

---

## コスト比較

### 従来構成 vs 新構成（月間1,000ユーザー想定）

| 項目 | 従来構成 | 新構成 | 差分 |
|------|---------|--------|------|
| OCR | $50 (Google) | $50 (Google) | ±0 |
| TTS | $30 (Google) | $40 (Azure) | +$10 |
| STT | $40 (Google) | $25 (Whisper) | -$15 |
| LLM | $0 | $60 (Claude) | +$60 |
| 翻訳 | $20 (DeepL) | $15 (Google) | -$5 |
| **合計** | **$140** | **$190** | **+$50** |

**+$50/月で得られる価値**:
- マイナー言語対応の大幅拡大（30言語 → 99言語+）
- ドメイン特化学習機能
- より柔軟な発音評価
- 高品質なコンテンツ生成

---

## まとめ

### 主要な変更点

1. **TTS**: Google → **Azure**（140言語、マイナー言語カバレッジ拡大）
2. **STT**: Google → **Whisper**（99言語、マイナー言語最強）
3. **発音評価**: Azure Native → **Whisper + LLM**（全言語対応）
4. **LLM**: なし → **Claude 3.5**（ドメイン特化に必須）
5. **翻訳**: DeepL優先 → **Google優先**（133言語カバレッジ）

### この構成がHaiLanGoに最適な理由

1. **マイナー言語サポート最大化**
   - Whisper: 99言語（クルド語、アムハラ語含む）
   - Azure TTS: 140言語
   - Google Translate: 133言語
   - Google Vision OCR: 200+言語

2. **ドメイン特化対応**
   - Claude 3.5でドメイン特化プロンプト
   - Whisper + LLMで任意言語の発音評価

3. **柔軟性**
   - 言語カテゴリ別の最適API自動選択
   - 新言語追加が容易（コード変更不要）

4. **コスト効率**
   - 月+$50で大幅な機能拡張
   - マイナー言語ユーザー獲得による差別化
