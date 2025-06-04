CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	description TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS task_sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id INTEGER NOT NULL,
  started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ended_at DATETIME,
  FOREIGN KEY (task_id) REFERENCES tasks(id)
);

CREATE TABLE IF NOT EXISTS task_session_intervals (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id INTEGER NOT NULL,
  start_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  end_time DATETIME,
  FOREIGN KEY (session_id) REFERENCES task_sessions(id)
);

CREATE TABLE IF NOT EXISTS commits (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	session_id INTEGER,
	hash TEXT UNIQUE,
	message TEXT,
	author TEXT,
	date TEXT,
	FOREIGN KEY (session_id) REFERENCES task_sessions(id) ON DELETE CASCADE
);


-- A task created creates a session and starts a session interval...a session can be paused(ie. session interval ends), a session can be resumed (ie a new session interval starts) and a session can be ended (the session interval ends which ends the session)
--
