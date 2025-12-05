# データベーススキーマ設計

## 概要

HaiLanGoのデータベース設計を定義する。
PostgreSQL 15+を使用し、適切なインデックスとリレーションを設計する。

---

## ER図

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│   users     │       │   books     │       │   pages     │
├─────────────┤       ├─────────────┤       ├─────────────┤
│ id (PK)     │──────<│ id (PK)     │──────<│ id (PK)     │
│ email       │       │ user_id(FK) │       │ book_id(FK) │
│ password    │       │ title       │       │ page_number │
│ created_at  │       │ target_lang │       │ image_path  │
│ updated_at  │       │ native_lang │       │ ocr_text    │
└─────────────┘       │ ref_lang    │       │ ocr_status  │
                      │ created_at  │       │ created_at  │
                      └─────────────┘       └─────────────┘
                                                   │
                      ┌─────────────┐              │
                      │ vocabularies│<─────────────┘
                      ├─────────────┤
                      │ id (PK)     │
                      │ page_id(FK) │
                      │ user_id(FK) │
                      │ word        │
                      │ meaning     │
                      │ context     │
                      └─────────────┘
                             │
                      ┌──────┴──────┐
                      │             │
               ┌─────────────┐ ┌─────────────┐
               │ srs_items   │ │learning_logs│
               ├─────────────┤ ├─────────────┤
               │ id (PK)     │ │ id (PK)     │
               │ vocab_id(FK)│ │ user_id(FK) │
               │ user_id(FK) │ │ page_id(FK) │
               │ easiness    │ │ vocab_id(FK)│
               │ interval    │ │ action_type │
               │ repetitions │ │ score       │
               │ next_review │ │ created_at  │
               └─────────────┘ └─────────────┘
```

---

## テーブル定義

### users - ユーザー

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    native_language VARCHAR(10) DEFAULT 'ja',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

| カラム | 型 | 説明 |
|--------|-----|------|
| id | UUID | 主キー |
| email | VARCHAR(255) | メールアドレス（ユニーク） |
| password_hash | VARCHAR(255) | bcryptハッシュ |
| display_name | VARCHAR(100) | 表示名 |
| native_language | VARCHAR(10) | 母国語コード（ISO 639-1） |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

---

### books - 教材

```sql
CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    target_language VARCHAR(10) NOT NULL,
    native_language VARCHAR(10) NOT NULL,
    reference_language VARCHAR(10),
    cover_image_path VARCHAR(500),
    total_pages INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'processing',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_books_user_id ON books(user_id);
CREATE INDEX idx_books_status ON books(status);
```

| カラム | 型 | 説明 |
|--------|-----|------|
| id | UUID | 主キー |
| user_id | UUID | 所有者（外部キー） |
| title | VARCHAR(255) | 教材タイトル |
| target_language | VARCHAR(10) | 学習言語（例: ru, ar, he） |
| native_language | VARCHAR(10) | 母国語 |
| reference_language | VARCHAR(10) | 参照言語（任意） |
| cover_image_path | VARCHAR(500) | 表紙画像パス |
| total_pages | INTEGER | 総ページ数 |
| status | VARCHAR(20) | processing / ready / error |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

---

### pages - ページ

```sql
CREATE TABLE pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    page_number INTEGER NOT NULL,
    image_path VARCHAR(500) NOT NULL,
    ocr_text TEXT,
    ocr_confidence DECIMAL(5,4),
    ocr_status VARCHAR(20) DEFAULT 'pending',
    ocr_processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(book_id, page_number)
);

CREATE INDEX idx_pages_book_id ON pages(book_id);
CREATE INDEX idx_pages_ocr_status ON pages(ocr_status);
```

| カラム | 型 | 説明 |
|--------|-----|------|
| id | UUID | 主キー |
| book_id | UUID | 教材（外部キー） |
| page_number | INTEGER | ページ番号 |
| image_path | VARCHAR(500) | 画像ファイルパス |
| ocr_text | TEXT | OCR抽出テキスト |
| ocr_confidence | DECIMAL | OCR信頼度（0-1） |
| ocr_status | VARCHAR(20) | pending / processing / completed / failed |
| ocr_processed_at | TIMESTAMP | OCR処理日時 |

---

### vocabularies - 単語・フレーズ

```sql
CREATE TABLE vocabularies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    page_id UUID REFERENCES pages(id) ON DELETE SET NULL,
    word VARCHAR(500) NOT NULL,
    reading VARCHAR(500),
    meaning TEXT NOT NULL,
    context TEXT,
    part_of_speech VARCHAR(50),
    language VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_vocabularies_user_id ON vocabularies(user_id);
CREATE INDEX idx_vocabularies_page_id ON vocabularies(page_id);
CREATE INDEX idx_vocabularies_word ON vocabularies(word);
```

| カラム | 型 | 説明 |
|--------|-----|------|
| id | UUID | 主キー |
| user_id | UUID | 所有者 |
| page_id | UUID | 出典ページ（任意） |
| word | VARCHAR(500) | 単語・フレーズ |
| reading | VARCHAR(500) | 読み方（ルビ等） |
| meaning | TEXT | 意味 |
| context | TEXT | 使用例・文脈 |
| part_of_speech | VARCHAR(50) | 品詞 |
| language | VARCHAR(10) | 言語コード |

---

### srs_items - 間隔反復学習データ

```sql
CREATE TABLE srs_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vocabulary_id UUID NOT NULL REFERENCES vocabularies(id) ON DELETE CASCADE,
    easiness_factor DECIMAL(4,2) DEFAULT 2.50,
    interval_days INTEGER DEFAULT 0,
    repetitions INTEGER DEFAULT 0,
    next_review_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_reviewed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(user_id, vocabulary_id)
);

CREATE INDEX idx_srs_items_user_id ON srs_items(user_id);
CREATE INDEX idx_srs_items_next_review ON srs_items(next_review_at);
```

| カラム | 型 | 説明 |
|--------|-----|------|
| id | UUID | 主キー |
| user_id | UUID | ユーザー |
| vocabulary_id | UUID | 対象単語 |
| easiness_factor | DECIMAL | SM-2のEF（2.5開始） |
| interval_days | INTEGER | 現在の復習間隔（日） |
| repetitions | INTEGER | 連続正解回数 |
| next_review_at | TIMESTAMP | 次回復習日時 |
| last_reviewed_at | TIMESTAMP | 最終復習日時 |

---

### learning_sessions - 学習セッション

```sql
CREATE TABLE learning_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID REFERENCES books(id) ON DELETE SET NULL,
    session_type VARCHAR(20) NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE,
    duration_seconds INTEGER,
    pages_studied INTEGER DEFAULT 0,
    words_learned INTEGER DEFAULT 0
);

CREATE INDEX idx_learning_sessions_user_id ON learning_sessions(user_id);
CREATE INDEX idx_learning_sessions_started_at ON learning_sessions(started_at);
```

| カラム | 型 | 説明 |
|--------|-----|------|
| id | UUID | 主キー |
| user_id | UUID | ユーザー |
| book_id | UUID | 教材（任意） |
| session_type | VARCHAR(20) | page_learning / review / teacher_mode |
| started_at | TIMESTAMP | 開始日時 |
| ended_at | TIMESTAMP | 終了日時 |
| duration_seconds | INTEGER | 学習時間（秒） |
| pages_studied | INTEGER | 学習ページ数 |
| words_learned | INTEGER | 学習単語数 |

---

### learning_logs - 学習ログ

```sql
CREATE TABLE learning_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES learning_sessions(id) ON DELETE SET NULL,
    page_id UUID REFERENCES pages(id) ON DELETE SET NULL,
    vocabulary_id UUID REFERENCES vocabularies(id) ON DELETE SET NULL,
    action_type VARCHAR(50) NOT NULL,
    score DECIMAL(5,2),
    feedback TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_learning_logs_user_id ON learning_logs(user_id);
CREATE INDEX idx_learning_logs_created_at ON learning_logs(created_at);
CREATE INDEX idx_learning_logs_action_type ON learning_logs(action_type);
```

| カラム | 型 | 説明 |
|--------|-----|------|
| id | UUID | 主キー |
| user_id | UUID | ユーザー |
| session_id | UUID | セッション |
| page_id | UUID | ページ（任意） |
| vocabulary_id | UUID | 単語（任意） |
| action_type | VARCHAR(50) | tts_listen / stt_pronounce / srs_review / page_complete |
| score | DECIMAL | スコア（0-100） |
| feedback | TEXT | AIフィードバック |

---

### user_stats - ユーザー統計（集計テーブル）

```sql
CREATE TABLE user_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    total_study_time_seconds BIGINT DEFAULT 0,
    total_words_learned INTEGER DEFAULT 0,
    total_pages_completed INTEGER DEFAULT 0,
    current_streak_days INTEGER DEFAULT 0,
    longest_streak_days INTEGER DEFAULT 0,
    last_study_date DATE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_user_stats_user_id ON user_stats(user_id);
```

---

### conversation_logs - AI会話ログ

```sql
CREATE TABLE conversation_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES learning_sessions(id) ON DELETE SET NULL,
    page_id UUID REFERENCES pages(id) ON DELETE SET NULL,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_conversation_logs_session_id ON conversation_logs(session_id);
CREATE INDEX idx_conversation_logs_created_at ON conversation_logs(created_at);
```

| カラム | 型 | 説明 |
|--------|-----|------|
| role | VARCHAR(20) | user / assistant / system |
| content | TEXT | メッセージ内容 |

---

## マイグレーション

### ファイル構成

```
backend/migrations/
├── 000001_create_users.up.sql
├── 000001_create_users.down.sql
├── 000002_create_books.up.sql
├── 000002_create_books.down.sql
├── 000003_create_pages.up.sql
├── 000003_create_pages.down.sql
├── 000004_create_vocabularies.up.sql
├── 000004_create_vocabularies.down.sql
├── 000005_create_srs_items.up.sql
├── 000005_create_srs_items.down.sql
├── 000006_create_learning_sessions.up.sql
├── 000006_create_learning_sessions.down.sql
├── 000007_create_learning_logs.up.sql
├── 000007_create_learning_logs.down.sql
├── 000008_create_user_stats.up.sql
├── 000008_create_user_stats.down.sql
├── 000009_create_conversation_logs.up.sql
└── 000009_create_conversation_logs.down.sql
```

### 実行コマンド

```bash
# マイグレーション実行
go run cmd/migrate/main.go up

# ロールバック
go run cmd/migrate/main.go down 1
```

---

## インデックス戦略

### 主要クエリとインデックス

| クエリ | 使用インデックス |
|--------|-----------------|
| ユーザーの教材一覧 | idx_books_user_id |
| 教材のページ一覧 | idx_pages_book_id |
| 今日の復習対象 | idx_srs_items_next_review |
| ユーザーの学習履歴 | idx_learning_logs_user_id, idx_learning_logs_created_at |
| 単語検索 | idx_vocabularies_word |

---

## Redis キャッシュ設計

### キー設計

```
# セッション
session:{session_id} -> JWT payload (TTL: 24h)

# OCR結果キャッシュ
ocr:{page_id} -> OCR result JSON (TTL: 7d)

# TTS音声キャッシュ
tts:{hash(text+lang+voice)} -> audio URL (TTL: 30d)

# レート制限
rate:{user_id}:{endpoint} -> count (TTL: 1min)

# 学習進捗（リアルタイム用）
progress:{user_id}:{book_id} -> current page (TTL: 1h)
```

---

## 次のドキュメント

- [04_API_SPECIFICATION.md](./04_API_SPECIFICATION.md) - API仕様
