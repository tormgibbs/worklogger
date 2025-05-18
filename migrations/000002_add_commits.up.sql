CREATE TABLE IF NOT EXISTS commits (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	session_id INTEGER,
	hash TEXT UNIQUE,
	message TEXT,
	author TEXT,
	date TEXT,
	FOREIGN KEY (session_id) REFERENCES task_sessions(id) ON DELETE CASCADE
);