package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type TaskSessionIntervalModel struct {
	DB *sql.DB
}

type TaskSessionInterval struct {
	ID        int
	SessionID int
	StartTime time.Time
	EndTime   *time.Time
}

func (m TaskSessionIntervalModel) CreateTX(tx *sql.Tx, sessionID int) error {
	query := `
		INSERT INTO task_session_intervals (session_id)
		VALUES (?)
	`
	_, err := tx.Exec(query, sessionID)
	return err
}

func (m TaskSessionIntervalModel) Create(sessionID int) error {
	query := `
		INSERT INTO task_session_intervals (session_id)
		VALUES (?)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, sessionID)
	return err
}

func (m TaskSessionIntervalModel) End(ts *TaskSession) (*TaskSessionInterval, error) {
	query := `
		UPDATE task_session_intervals
		SET end_time = CURRENT_TIMESTAMP
		WHERE session_id = ? AND end_time IS NULL
		RETURNING id, session_id, start_time, end_time
	`
	var tsi TaskSessionInterval

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, ts.ID).Scan(
		&tsi.ID,
		&tsi.SessionID,
		&tsi.StartTime,
		&tsi.EndTime,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil
		default:
			return nil, err
		}
	}
	return &tsi, nil
}

func (m TaskSessionIntervalModel) HasOpenInterval(sessionID int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM task_session_intervals
		WHERE session_id = ? AND end_time IS NULL
	`
	var count int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, sessionID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m TaskSessionIntervalModel) StartNew(sessionID int) (*TaskSessionInterval, error) {
	query := `
		INSERT INTO task_session_intervals (session_id, start_time)
		VALUES (?, CURRENT_TIMESTAMP)
		RETURNING id, session_id, start_time, end_time
	`

	var tsi TaskSessionInterval

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, sessionID).Scan(
		&tsi.ID,
		&tsi.SessionID,
		&tsi.StartTime,
		&tsi.EndTime,
	)

	if err != nil {
		return nil, err
	}

	return &tsi, nil
}

