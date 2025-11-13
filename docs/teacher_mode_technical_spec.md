# 教師モード（自動学習モード）- 技術仕様書

## 1. 概要

### 1.1 機能説明
教師モードは、ボタン一つで本を最初から最後まで自動的に読み進めながら、教師が授業をするように音声読み上げと解説を連続して行う機能です。

### 1.2 主な特徴
- **ハンズフリー学習**：操作不要で連続学習
- **バックグラウンド再生**：画面オフでも継続
- **オフライン対応**：事前ダウンロードで通信不要
- **高度なカスタマイズ**：ユーザーの学習スタイルに合わせた設定

---

## 2. 機能要件

### 2.1 基本フロー

```
開始
  ↓
ページ1を表示
  ↓
学習先言語のテキストを読み上げ（例：Здравствуйте!）
  ↓
[設定により] 母国語訳を読み上げ（例：こんにちは）
  ↓
[設定により] 単語・文法の解説をAI音声で説明
  ↓
[設定により] 発音練習の時間（ユーザーが発話）
  ↓
設定した間隔だけ待機（例：5秒）
  ↓
ページ2へ自動遷移
  ↓
繰り返し（全ページ完了まで）
  ↓
終了
```

### 2.2 ユーザー操作

#### 開始時
```
1. 学習画面で「教師モード」ボタンをタップ
2. 教師モード設定画面が表示（初回のみ詳細設定）
3. 「開始」ボタンで教師モードスタート
```

#### 再生中
```
- 一時停止：画面タップ or ロック画面のボタン
- 再開：もう一度タップ
- 前のページへ：スワイプ右 or コントロールボタン
- 次のページへ：スワイプ左 or コントロールボタン
- 停止：停止ボタン → 通常モードへ戻る
- 設定変更：設定アイコン → その場で反映
```

### 2.3 カスタマイズ可能な設定

#### 音声設定
```yaml
再生速度:
  - 0.5x（超ゆっくり）
  - 0.75x（ゆっくり）
  - 1.0x（標準）  # デフォルト
  - 1.25x（少し速い）
  - 1.5x（速い）
  - 2.0x（超高速）

音質:
  - 標準（無料プラン）
  - 高品質（プレミアムプラン）  # より自然な発音
```

#### タイミング設定
```yaml
ページ間隔:
  - 0秒（即座に次へ）
  - 3秒
  - 5秒  # デフォルト
  - 10秒
  - 15秒
  - 30秒（じっくり考える時間）

フレーズのリピート回数:
  - 1回  # デフォルト
  - 2回（確認用）
  - 3回（完全習得用）
```

#### 学習内容
```yaml
必須:
  - 学習先言語の読み上げ: true（変更不可）

オプション:
  - 母国語訳の読み上げ: true  # デフォルトON
  - 単語の解説: true  # デフォルトON
  - 文法の解説: false  # デフォルトOFF（日常会話重視のため）
  - 発音練習を含む: false  # デフォルトOFF（聞き流し時に邪魔になるため）
  - 例文の追加: false  # デフォルトOFF（時間短縮）
```

---

## 3. 技術実装

### 3.1 システムアーキテクチャ

```
┌─────────────────────────────────────────────┐
│           Frontend (Web/Mobile)             │
│                                             │
│  ┌─────────────────────────────────────┐   │
│  │   教師モードコントローラー           │   │
│  │   - 状態管理（再生/停止/一時停止）   │   │
│  │   - ページ遷移制御                   │   │
│  │   - タイマー管理                     │   │
│  └─────────────────────────────────────┘   │
│              ↓↑                             │
│  ┌─────────────────────────────────────┐   │
│  │   音声再生エンジン                   │   │
│  │   - TTS API呼び出し                  │   │
│  │   - キャッシュ管理                   │   │
│  │   - バックグラウンド再生             │   │
│  └─────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
              ↓↑
┌─────────────────────────────────────────────┐
│              Backend (Go)                   │
│                                             │
│  ┌─────────────────────────────────────┐   │
│  │   教師モード生成API                  │   │
│  │   - OCRテキスト取得                  │   │
│  │   - 辞書API連携（単語解説）          │   │
│  │   - AI解説生成（文法など）           │   │
│  │   - TTS音声生成リクエスト            │   │
│  └─────────────────────────────────────┘   │
│              ↓↑                             │
│  ┌─────────────────────────────────────┐   │
│  │   音声データ管理                     │   │
│  │   - 音声ファイルのキャッシュ         │   │
│  │   - ダウンロードパッケージ生成       │   │
│  │   - ストリーミング配信               │   │
│  └─────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
              ↓↑
┌─────────────────────────────────────────────┐
│            External APIs                    │
│                                             │
│  - Google Cloud TTS / Amazon Polly          │
│  - ElevenLabs (プレミアム高品質)            │
│  - 辞書API (単語解説)                       │
│  - Claude/GPT API (解説生成)                │
└─────────────────────────────────────────────┘
```

### 3.2 データフロー

#### 音声生成フロー
```
1. ユーザーが教師モードを開始
   ↓
2. バックエンドにリクエスト送信
   {
     book_id: "xxx",
     page_start: 1,
     page_end: 150,
     settings: {
       speed: 1.0,
       include_translation: true,
       include_word_explanation: true,
       ...
     }
   }
   ↓
3. バックエンドで各ページの音声スクリプト生成
   ページ1:
   - 学習先言語テキスト: "Здравствуйте!"
   - 母国語訳: "こんにちは"
   - 単語解説: "Здравствуйтеは丁寧な挨拶です..."
   ↓
4. TTS APIで音声生成
   - 各セグメントを個別に生成
   - キャッシュに保存（Redis）
   ↓
5. フロントエンドに音声URL返却
   [
     {page: 1, segments: [
       {type: "phrase", url: "https://..."},
       {type: "translation", url: "https://..."},
       {type: "explanation", url: "https://..."}
     ]},
     ...
   ]
   ↓
6. フロントエンドで順次再生
```

### 3.3 バックグラウンド再生の実装

#### Web版（PWA）
```javascript
// Service Workerでのバックグラウンド処理
// Media Session APIを使用

navigator.mediaSession.metadata = new MediaMetadata({
  title: 'ロシア語入門 - ページ12',
  artist: 'LinguaAI 教師モード',
  album: 'ロシア語入門',
  artwork: [
    { src: 'book-cover.png', sizes: '512x512', type: 'image/png' }
  ]
});

navigator.mediaSession.setActionHandler('play', () => {
  audio.play();
});

navigator.mediaSession.setActionHandler('pause', () => {
  audio.pause();
});

navigator.mediaSession.setActionHandler('previoustrack', () => {
  goToPreviousPage();
});

navigator.mediaSession.setActionHandler('nexttrack', () => {
  goToNextPage();
});
```

#### Mobile版（Flutter）
```dart
// audio_service パッケージを使用

AudioHandler audioHandler = await AudioService.init(
  builder: () => TeacherModeAudioHandler(),
  config: AudioServiceConfig(
    androidNotificationChannelId: 'com.linguaai.teacher_mode',
    androidNotificationChannelName: 'Teacher Mode',
    androidNotificationIcon: 'drawable/ic_notification',
  ),
);

class TeacherModeAudioHandler extends BaseAudioHandler {
  @override
  Future<void> play() async {
    // 再生処理
  }

  @override
  Future<void> pause() async {
    // 一時停止処理
  }

  @override
  Future<void> skipToNext() async {
    // 次のページへ
  }

  @override
  Future<void> skipToPrevious() async {
    // 前のページへ
  }
}
```

### 3.4 オフライン対応

#### ダウンロード処理
```go
// Go バックエンド

func GenerateTeacherModePackage(bookID string, settings TeacherModeSettings) (string, error) {
    // 1. すべてのページのOCRテキストを取得
    pages := GetBookPages(bookID)
    
    // 2. 各ページの音声スクリプトを生成
    var audioFiles []AudioFile
    for _, page := range pages {
        // 学習先言語
        phraseAudio := GenerateTTS(page.PhraseText, page.TargetLanguage, settings)
        audioFiles = append(audioFiles, phraseAudio)
        
        // 母国語訳（設定による）
        if settings.IncludeTranslation {
            translationAudio := GenerateTTS(page.TranslationText, page.NativeLanguage, settings)
            audioFiles = append(audioFiles, translationAudio)
        }
        
        // 単語解説（設定による）
        if settings.IncludeWordExplanation {
            explanationText := GenerateExplanation(page.Words)
            explanationAudio := GenerateTTS(explanationText, page.NativeLanguage, settings)
            audioFiles = append(audioFiles, explanationAudio)
        }
    }
    
    // 3. ZIPファイルにパッケージング
    zipPath := CreateZipPackage(audioFiles, bookID)
    
    // 4. ダウンロードURLを返却
    return uploadToStorage(zipPath), nil
}
```

#### フロントエンドでのダウンロード管理
```typescript
// TypeScript (Next.js)

interface DownloadPackage {
  bookId: string;
  totalSize: number;
  downloadProgress: number;
  status: 'pending' | 'downloading' | 'completed' | 'failed';
}

class TeacherModeDownloadManager {
  async downloadPackage(bookId: string, settings: TeacherModeSettings): Promise<void> {
    // 1. パッケージのメタデータ取得
    const metadata = await api.getTeacherModePackageMetadata(bookId, settings);
    
    // 2. ダウンロード開始
    const response = await fetch(metadata.downloadUrl);
    const reader = response.body.getReader();
    
    // 3. 進捗管理
    let receivedLength = 0;
    const chunks = [];
    
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      
      chunks.push(value);
      receivedLength += value.length;
      
      // 進捗更新
      this.updateProgress(bookId, receivedLength / metadata.totalSize);
    }
    
    // 4. IndexedDBに保存
    const blob = new Blob(chunks);
    await this.saveToIndexedDB(bookId, blob, settings);
  }
  
  async playOffline(bookId: string): Promise<void> {
    // IndexedDBから取得して再生
    const package = await this.loadFromIndexedDB(bookId);
    // 再生処理...
  }
}
```

---

## 4. データモデル

### 4.1 教師モード設定

```typescript
interface TeacherModeSettings {
  speed: 0.5 | 0.75 | 1.0 | 1.25 | 1.5 | 2.0;
  pageInterval: number; // 秒単位（0-30）
  repeatCount: 1 | 2 | 3;
  audioQuality: 'standard' | 'premium';
  content: {
    includeTranslation: boolean;
    includeWordExplanation: boolean;
    includeGrammarExplanation: boolean;
    includePronunciationPractice: boolean;
    includeExampleSentences: boolean;
  };
}
```

### 4.2 音声セグメント

```typescript
interface AudioSegment {
  id: string;
  type: 'phrase' | 'translation' | 'explanation' | 'pause';
  audioUrl: string;
  duration: number; // ミリ秒
  text: string;
  language: string;
}

interface PageAudio {
  pageNumber: number;
  segments: AudioSegment[];
  totalDuration: number;
}

interface TeacherModePlaylist {
  bookId: string;
  pages: PageAudio[];
  settings: TeacherModeSettings;
}
```

### 4.3 再生状態

```typescript
interface PlaybackState {
  status: 'stopped' | 'playing' | 'paused';
  currentPage: number;
  currentSegmentIndex: number;
  elapsedTime: number;
  totalDuration: number;
}
```

---

## 5. データベーススキーマ

### 5.1 PostgreSQL

```sql
-- 教師モードのダウンロード履歴
CREATE TABLE teacher_mode_downloads (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id),
  book_id UUID NOT NULL REFERENCES books(id),
  settings JSONB NOT NULL,
  total_size_bytes BIGINT NOT NULL,
  downloaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMP,
  INDEX idx_user_book (user_id, book_id)
);

-- 教師モードの再生履歴
CREATE TABLE teacher_mode_playback_history (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id),
  book_id UUID NOT NULL REFERENCES books(id),
  last_page INTEGER NOT NULL,
  last_played_at TIMESTAMP NOT NULL DEFAULT NOW(),
  total_play_time_seconds INTEGER DEFAULT 0,
  INDEX idx_user_book (user_id, book_id)
);
```

### 5.2 Redis

```redis
# 音声キャッシュ（TTL: 7日）
teacher_mode:audio:{book_id}:{page}:{segment_type}:{hash} → 音声データURL

# 再生状態（TTL: 24時間）
teacher_mode:state:{user_id}:{book_id} → PlaybackState JSON

# ダウンロード進捗
teacher_mode:download:{user_id}:{book_id} → { progress: 0.45, status: "downloading" }
```

---

## 6. API設計

### 6.1 教師モード生成API

```http
POST /api/v1/books/{book_id}/teacher-mode/generate
Content-Type: application/json
Authorization: Bearer {token}

Request Body:
{
  "settings": {
    "speed": 1.0,
    "pageInterval": 5,
    "repeatCount": 1,
    "audioQuality": "standard",
    "content": {
      "includeTranslation": true,
      "includeWordExplanation": true,
      "includeGrammarExplanation": false,
      "includePronunciationPractice": false
    }
  },
  "pageRange": {
    "start": 1,
    "end": 150
  }
}

Response:
{
  "playlistId": "uuid",
  "totalPages": 150,
  "estimatedDuration": 4500, // 秒
  "pages": [
    {
      "pageNumber": 1,
      "segments": [
        {
          "type": "phrase",
          "audioUrl": "https://cdn.linguaai.com/audio/...",
          "duration": 2000,
          "text": "Здравствуйте!"
        },
        {
          "type": "translation",
          "audioUrl": "https://cdn.linguaai.com/audio/...",
          "duration": 1500,
          "text": "こんにちは"
        }
      ]
    }
  ]
}
```

### 6.2 オフラインパッケージ生成API

```http
POST /api/v1/books/{book_id}/teacher-mode/download-package
Content-Type: application/json
Authorization: Bearer {token}

Request Body:
{
  "settings": { ... }
}

Response:
{
  "packageId": "uuid",
  "downloadUrl": "https://cdn.linguaai.com/packages/...",
  "totalSize": 248000000, // バイト
  "expiresAt": "2025-11-20T12:00:00Z"
}
```

### 6.3 再生状態保存API

```http
PUT /api/v1/books/{book_id}/teacher-mode/playback-state
Content-Type: application/json
Authorization: Bearer {token}

Request Body:
{
  "currentPage": 12,
  "currentSegmentIndex": 2,
  "elapsedTime": 3600
}

Response:
{
  "success": true
}
```

---

## 7. パフォーマンス最適化

### 7.1 音声生成のキャッシュ戦略

```
レベル1: CDN（CloudFlare）
  - 生成済み音声ファイル
  - TTL: 30日

レベル2: Redis
  - 音声URL、メタデータ
  - TTL: 7日

レベル3: PostgreSQL
  - 音声生成履歴
  - 永続化
```

### 7.2 プリロード戦略

```javascript
// 次の3ページ分を事前ロード
class AudioPreloader {
  preloadNext(currentPage: number, playlist: TeacherModePlaylist) {
    const pagesToPreload = [currentPage + 1, currentPage + 2, currentPage + 3];
    
    pagesToPreload.forEach(pageNum => {
      if (pageNum <= playlist.pages.length) {
        const page = playlist.pages[pageNum - 1];
        page.segments.forEach(segment => {
          this.preloadAudio(segment.audioUrl);
        });
      }
    });
  }
  
  private preloadAudio(url: string) {
    const audio = new Audio();
    audio.preload = 'auto';
    audio.src = url;
  }
}
```

### 7.3 バッチ生成

```go
// 夜間バッチで人気の本の音声を事前生成
func PreGeneratePopularBooks() {
    books := GetPopularBooks(limit: 100)
    
    for _, book := range books {
        // 標準設定で事前生成
        GenerateTeacherModePackage(book.ID, DefaultSettings)
        
        // プレミアム設定でも事前生成
        GenerateTeacherModePackage(book.ID, PremiumSettings)
    }
}
```

---

## 8. コスト試算

### 8.1 TTS APIコスト

**Google Cloud TTS料金**
- Standard: $4 / 100万文字
- WaveNet（高品質）: $16 / 100万文字

**1ページあたりの推定文字数**
- 学習先言語のフレーズ: 50文字
- 母国語訳: 50文字
- 単語解説: 200文字
- 合計: 約300文字/ページ

**150ページの本の場合**
- 総文字数: 45,000文字
- Standard: $0.18
- WaveNet: $0.72

**月間1,000ユーザーが各2冊使用する場合**
- 総コスト（Standard）: $360
- 総コスト（WaveNet）: $1,440

### 8.2 ストレージコスト

**S3料金（東京リージョン）**
- ストレージ: $0.025 / GB / 月
- 転送: $0.114 / GB（最初の10TB）

**1冊あたりの音声データサイズ**
- Standard: 約150MB
- WaveNet: 約250MB

**月間1,000ユーザー、各2冊**
- ストレージ: 500GB × $0.025 = $12.5
- 転送（1回ダウンロード）: 500GB × $0.114 = $57

**総月間コスト（Standard）: 約$430**
**総月間コスト（WaveNet）: 約$1,510**

### 8.3 収益モデルとの整合性

**プレミアムプラン月額: $9.99**
- 1,000ユーザーの10%が有料: 100ユーザー
- 月間収益: $999

**結論**
- Standard音質で運営すれば、100人のプレミアムユーザーで収益化可能
- 初期はキャッシュ活用で音声生成を最小化
- 人気の本は事前生成してコスト削減

---

## 9. テスト計画

### 9.1 単体テスト
- [ ] 音声生成ロジック
- [ ] ページ遷移ロジック
- [ ] タイマー制御
- [ ] 設定の保存・読み込み

### 9.2 統合テスト
- [ ] TTS API連携
- [ ] バックグラウンド再生
- [ ] オフライン再生
- [ ] キャッシュ機能

### 9.3 E2Eテスト
- [ ] 教師モードの開始〜終了
- [ ] 一時停止・再開
- [ ] ページスキップ
- [ ] ダウンロード機能

### 9.4 パフォーマンステスト
- [ ] 連続再生の安定性（1時間以上）
- [ ] メモリリーク検証
- [ ] バッテリー消費量

---

## 10. ロードマップ

### Phase 1: 基本機能（1ヶ月）
- [ ] 音声連続再生
- [ ] 基本設定（速度、間隔）
- [ ] ページ自動遷移

### Phase 2: バックグラウンド再生（2週間）
- [ ] Media Session API実装（Web）
- [ ] audio_service実装（Flutter）
- [ ] ロック画面コントロール

### Phase 3: オフライン対応（3週間）
- [ ] ダウンロード機能
- [ ] IndexedDB保存
- [ ] オフライン再生

### Phase 4: 高度な機能（2週間）
- [ ] AI解説生成
- [ ] カスタム設定プリセット
- [ ] 再生統計ダッシュボード

---

## まとめ

教師モードは、LinguaAIの差別化要因となる重要な機能です。
ユーザーが「ながら学習」できることで、学習時間を大幅に増やすことができます。

**重要なポイント**
1. バックグラウンド再生の安定性
2. オフライン対応によるユーザビリティ向上
3. 音声品質とコストのバランス
4. キャッシュ戦略によるコスト削減
