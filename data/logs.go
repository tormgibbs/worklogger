package data

import (
	"database/sql"
	"fmt"
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

func (m LogModel) GetLogs() ([]Log, error) {
	query := `
	SELECT
		ts.id,
		t.description,
		ts.started_at,
		ts.ended_at,
		c.message,
		c.hash,
		c.author,
		c.date
	FROM task_sessions ts
	JOIN tasks t ON t.id = ts.task_id
	LEFT JOIN commits c ON c.task_session_id = ts.id
	ORDER BY DATE(ts.started_at) DESC, ts.started_at ASC, c.id ASC;
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type rowData struct {
		sessionID   int
		description string
		startedAt   time.Time
		endedAt     time.Time
		message     sql.NullString
		hash        sql.NullString
		author      sql.NullString
		commitDate  sql.NullString
	}

	var rowsData []rowData
	for rows.Next() {
		var r rowData
		err := rows.Scan(
			&r.sessionID,
			&r.description,
			&r.startedAt,
			&r.endedAt,
			&r.message,
			&r.hash,
			&r.author,
			&r.commitDate,
		)
		if err != nil {
			return nil, err
		}
		rowsData = append(rowsData, r)
	}

	// Group rows by session
	sessionsMap := make(map[int]*LogSession)
	for _, row := range rowsData {
		session, exists := sessionsMap[row.sessionID]
		if !exists {
			active, paused, total, err := TaskSessionModel(m).GetDurations(row.sessionID)
			if err != nil {
				return nil, fmt.Errorf("getting durations for session %d: %w", row.sessionID, err)
			}

			sessionsMap[row.sessionID] = &LogSession{
				ID:         row.sessionID,
				Task:       row.description,
				StartedAt:  row.startedAt,
				EndedAt:    row.endedAt,
				ActiveTime: active,
				PausedTime: paused,
				TotalTime:  total,
			}
			session = sessionsMap[row.sessionID]
		}

		if row.message.Valid {
			session.Commits = append(session.Commits, LogCommit{
				Message: row.message.String,
				Hash:    row.hash.String,
				Author:  row.author.String,
				Date:    row.commitDate.String,
			})
		}
	}

	// Group sessions by date
	logMap := make(map[string]*Log)
	for _, session := range sessionsMap {
		dateKey := session.StartedAt.Format("Jan 02") // like "May 15"

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

	// Convert map to slice
	var logs []Log
	for _, log := range logMap {
		logs = append(logs, *log)
	}

	// Optional: sort logs by date (desc)
	// You can do this if you care about order

	return logs, nil
}
