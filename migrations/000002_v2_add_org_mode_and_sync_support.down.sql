DROP INDEX IF EXISTS idx_task_sessions_mode;
DROP INDEX IF EXISTS idx_task_sessions_synced;
DROP INDEX IF EXISTS idx_task_sessions_started_at;
DROP INDEX IF EXISTS idx_session_tags_session_id;
DROP INDEX IF EXISTS idx_session_tags_tag;
DROP INDEX IF EXISTS idx_session_kpis_session_id;
DROP INDEX IF EXISTS idx_session_kpis_kpi;


DROP TABLE IF EXISTS session_tags;
DROP TABLE IF EXISTS session_kpis;


-- PRAGMA foreign_keys=off;

-- CREATE TABLE task_sessions_new (
--   id INTEGER PRIMARY KEY AUTOINCREMENT,
--   task_id INTEGER NOT NULL,
--   started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   ended_at DATETIME,
--   FOREIGN KEY (task_id) REFERENCES tasks(id)
-- );

-- INSERT INTO task_sessions_new (id, task_id, started_at, ended_at)
-- SELECT id, task_id, started_at, ended_at
-- FROM task_sessions;

-- DROP TABLE task_sessions;

-- ALTER TABLE task_sessions_new RENAME TO task_sessions;

-- PRAGMA foreign_keys=on;
