package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type TaskSessionModel struct {
	DB *sql.DB
}

type TaskSession struct {
	ID        int
	TaskID    int
	StartedAt time.Time
	EndedAt   *time.Time
}

type DetailedTaskSession struct {
	TaskSession
	Task Task
}

func (m TaskSessionModel) CreateTX(tx *sql.Tx, taskID int) (*TaskSession, error) {
	query := `
		INSERT INTO task_sessions (task_id)
		VALUES (?)
		RETURNING id, task_id, started_at
	`
	var ts TaskSession

	err := tx.QueryRow(query, taskID).Scan(&ts.ID, &ts.TaskID, &ts.StartedAt)
	if err != nil {
		return nil, err
	}

	return &ts, nil
}

func (m TaskSessionModel) GetAllWithTask() ([]*DetailedTaskSession, error) {
	query := `
		SELECT 
			ts.id,
			ts.task_id,
			ts.started_at,
			ts.ended_at,
			t.id,
			t.description,
			t.created_at
		FROM task_sessions ts
		JOIN tasks t ON ts.task_id = t.id
		ORDER BY ts.started_at DESC;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*DetailedTaskSession

	for rows.Next() {
		var ts TaskSession
		var t Task
		var endedAt sql.NullTime

		err := rows.Scan(
			&ts.ID,
			&ts.TaskID,
			&ts.StartedAt,
			&endedAt,
			&t.ID,
			&t.Description,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if endedAt.Valid {
			ts.EndedAt = &endedAt.Time
		}

		sessions = append(sessions, &DetailedTaskSession{
			TaskSession: ts,
			Task:        t,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil

}

func (m TaskSessionModel) GetByID(sessionID int) (*TaskSession, error) {
	query := `
		SELECT 
			id,
			task_id,
			started_at,
			ended_at
		FROM task_sessions
		WHERE id = ?
	`

	var ts TaskSession
	var endedAt sql.NullTime

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, sessionID).Scan(
		&ts.ID,
		&ts.TaskID,
		&ts.StartedAt,
		&endedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound // not found, return nil
		default:
			return nil, err // other DB error
		}
	}

	if endedAt.Valid {
		ts.EndedAt = &endedAt.Time
	} else {
		ts.EndedAt = nil
	}

	return &ts, nil
}


func (m TaskSessionModel) Get() (*TaskSession, error) {
	query := `
		SELECT ts.id
		FROM task_sessions ts
		WHERE ts.ended_at IS NULL
		LIMIT 1
	`
	var ts TaskSession

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query).Scan(&ts.ID)
	if err != nil {
		switch {
		// No active task session found
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil
		default:
			return nil, err
		}
	}

	return &ts, nil
}

func (m TaskSessionModel) Stop(sessionID int) (*TaskSession, error) {
	query := `
		UPDATE task_sessions
		SET ended_at = CURRENT_TIMESTAMP
		WHERE id = ? AND ended_at IS NULL
		RETURNING id, task_id, started_at, ended_at
	`

	var ts TaskSession

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, sessionID).Scan(
		&ts.ID,
		&ts.TaskID,
		&ts.StartedAt,
		&ts.EndedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil
		default:
			return nil, err
		}
	}

	return &ts, nil
}
