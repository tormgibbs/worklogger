package data

import (
	"database/sql"
	"time"
)

type LogModel struct {
	DB *sql.DB
}

type Log struct {
	Date     string
	Sessions []LogSession
}

type LogSession struct {
	ID         int
	Task       string
	StartedAt  time.Time
	EndedAt    time.Time
	ActiveTime time.Duration
	PausedTime time.Duration
	TotalTime  time.Duration
	Commits    []LogCommit
}

type LogCommit struct {
	Message string
	Hash    string
	Author  string
	Date    string
}

func (m LogModel) GetLogsWithDurations() ([]Log, error) {
	query := `
	WITH interval_totals AS (
		SELECT
			session_id,
			SUM(CAST(strftime('%s', end_time) - strftime('%s', start_time) AS INTEGER)) AS active_seconds
		FROM task_session_intervals
		WHERE end_time IS NOT NULL
		GROUP BY session_id
	),
	session_logs AS (
		SELECT
			ts.id AS session_id,
			t.description AS task_description,
			ts.started_at,
			ts.ended_at,
			CAST(strftime('%s', COALESCE(ts.ended_at, CURRENT_TIMESTAMP)) - strftime('%s', ts.started_at) AS INTEGER) AS total_seconds,
			IFNULL(it.active_seconds, 0) AS active_seconds,
			c.message AS commit_message,
			c.hash AS commit_hash,
			c.author AS commit_author,
			c.date AS commit_date
		FROM task_sessions ts
		JOIN tasks t ON t.id = ts.task_id
		LEFT JOIN interval_totals it ON it.session_id = ts.id
		LEFT JOIN commits c ON c.session_id = ts.id
	),
	orphan_commits AS (
		SELECT
			NULL AS session_id,
			'[Unassociated]' AS task_description,
			NULL AS started_at,
			NULL AS ended_at,
			0 AS total_seconds,
			0 AS active_seconds,
			c.message AS commit_message,
			c.hash AS commit_hash,
			c.author AS commit_author,
			c.date AS commit_date
		FROM commits c
		WHERE c.session_id IS NULL
	)
	SELECT *
	FROM (
		SELECT * FROM session_logs
		UNION ALL
		SELECT * FROM orphan_commits
	)
	ORDER BY 
		started_at IS NULL,       -- push NULLs (orphans) to the bottom
		date(started_at) DESC,
		started_at ASC;
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type rowData struct {
		SessionID     sql.NullInt64
		Task          string
		StartedAt     NullTime
		EndedAt       NullTime
		TotalSeconds  sql.NullInt64
		ActiveSeconds int64
		Message       sql.NullString
		Hash          sql.NullString
		Author        sql.NullString
		CommitDate    sql.NullString
	}

	rowsData := []rowData{}
	for rows.Next() {
		var r rowData
		if err := rows.Scan(
			&r.SessionID,
			&r.Task,
			&r.StartedAt,
			&r.EndedAt,
			&r.TotalSeconds,
			&r.ActiveSeconds,
			&r.Message,
			&r.Hash,
			&r.Author,
			&r.CommitDate,
		); err != nil {
			return nil, err
		}
		rowsData = append(rowsData, r)
	}

	sessionsMap := make(map[int]*LogSession)
	const orphanKey = -1
	var orphanSession *LogSession

	for _, row := range rowsData {
		if !row.SessionID.Valid {
			if orphanSession == nil {
				orphanSession = &LogSession{
					ID:        orphanKey,
					Task:      "[Unassociated]",
					StartedAt: time.Time{},
					EndedAt:   time.Time{},
				}
			}
			if row.Message.Valid {
				orphanSession.Commits = append(orphanSession.Commits, LogCommit{
					Message: row.Message.String,
					Hash:    row.Hash.String,
					Author:  row.Author.String,
					Date:    row.CommitDate.String,
				})
			}
			continue
		}

		sessionID := int(row.SessionID.Int64)
		session, exists := sessionsMap[sessionID]
		if !exists {
			total := time.Duration(row.TotalSeconds.Int64) * time.Second
			active := time.Duration(row.ActiveSeconds) * time.Second
			paused := total - active
			if paused < 0 {
				paused = 0
			}

			startedAt := time.Time{}
			if row.StartedAt.Valid {
				startedAt = row.StartedAt.Time
			}
			endedAt := time.Time{}
			if row.EndedAt.Valid {
				endedAt = row.EndedAt.Time
			}

			sessionsMap[sessionID] = &LogSession{
				ID:         sessionID,
				Task:       row.Task,
				StartedAt:  startedAt,
				EndedAt:    endedAt,
				TotalTime:  total,
				ActiveTime: active,
				PausedTime: paused,
			}
			session = sessionsMap[sessionID]
		}

		if row.Message.Valid {
			session.Commits = append(session.Commits, LogCommit{
				Message: row.Message.String,
				Hash:    row.Hash.String,
				Author:  row.Author.String,
				Date:    row.CommitDate.String,
			})
		}
	}

	// Group by date
	logMap := make(map[string]*Log)

	for _, session := range sessionsMap {
		dateKey := session.StartedAt.Format("Jan 02")
		logDay, exists := logMap[dateKey]
		if !exists {
			logMap[dateKey] = &Log{
				Date:     dateKey,
				Sessions: []LogSession{*session},
			}
		} else {
			logDay.Sessions = append(logDay.Sessions, *session)
		}
	}

	// Add orphan commits to their own group
	if orphanSession != nil && len(orphanSession.Commits) > 0 {
		logMap["Unassociated"] = &Log{
			Date:     "Unassociated",
			Sessions: []LogSession{*orphanSession},
		}
	}

	var logs []Log
	for _, log := range logMap {
		logs = append(logs, *log)
	}

	return logs, nil
}

