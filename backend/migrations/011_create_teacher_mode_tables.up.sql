-- 教師モードのダウンロード履歴テーブル
CREATE TABLE IF NOT EXISTS teacher_mode_downloads (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
  settings JSONB NOT NULL,
  total_size_bytes BIGINT NOT NULL,
  downloaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_teacher_mode_downloads_user_book ON teacher_mode_downloads(user_id, book_id);
CREATE INDEX IF NOT EXISTS idx_teacher_mode_downloads_user_id ON teacher_mode_downloads(user_id);
CREATE INDEX IF NOT EXISTS idx_teacher_mode_downloads_book_id ON teacher_mode_downloads(book_id);

-- 教師モードの再生履歴テーブル
CREATE TABLE IF NOT EXISTS teacher_mode_playback_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
  current_page INTEGER NOT NULL DEFAULT 1,
  current_segment_index INTEGER NOT NULL DEFAULT 0,
  elapsed_time INTEGER NOT NULL DEFAULT 0, -- 秒単位
  total_play_time_seconds INTEGER NOT NULL DEFAULT 0,
  last_played_at TIMESTAMP NOT NULL DEFAULT NOW(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, book_id)
);

CREATE INDEX IF NOT EXISTS idx_teacher_mode_playback_user_book ON teacher_mode_playback_history(user_id, book_id);
CREATE INDEX IF NOT EXISTS idx_teacher_mode_playback_user_id ON teacher_mode_playback_history(user_id);
CREATE INDEX IF NOT EXISTS idx_teacher_mode_playback_book_id ON teacher_mode_playback_history(book_id);
CREATE INDEX IF NOT EXISTS idx_teacher_mode_playback_last_played ON teacher_mode_playback_history(last_played_at DESC);
