ALTER TABLE task_sessions ADD COLUMN mode TEXT DEFAULT 'personal';
ALTER TABLE task_sessions ADD COLUMN notes TEXT;
ALTER TABLE task_sessions ADD COLUMN synced BOOLEAN DEFAULT 0;

CREATE TABLE IF NOT EXISTS session_tags (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id INTEGER NOT NULL,
  tag TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (session_id) REFERENCES task_sessions(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS session_kpis (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id INTEGER NOT NULL,
  kpi TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (session_id) REFERENCES task_sessions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_task_sessions_mode ON task_sessions(mode);
CREATE INDEX IF NOT EXISTS idx_task_sessions_sync_status ON task_sessions(synced);
CREATE INDEX IF NOT EXISTS idx_task_sessions_started_at ON task_sessions(started_at);
CREATE INDEX IF NOT EXISTS idx_session_tags_session_id ON session_tags(session_id);
CREATE INDEX IF NOT EXISTS idx_session_tags_tag ON session_tags(tag);
CREATE INDEX IF NOT EXISTS idx_session_kpis_session_id ON session_kpis(session_id);
CREATE INDEX IF NOT EXISTS idx_session_kpis_kpi ON session_kpis(kpi);