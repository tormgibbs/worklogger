package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func (m Models) CreateTask(tx *sql.Tx, task *Task) error {
	err := m.Tasks.CreateTx(tx, task)
	if err != nil {
		tx.Rollback()
		return err
	}

	session, err := m.TaskSessions.CreateTX(tx, task.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = m.TaskSessionIntervals.CreateTX(tx, session.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m TaskSessionModel) GetDurations(sessionID int) (totalTime, activeTime, pausedTime time.Duration, err error) {
	var (
		startedAt, endedAt time.Time
		activeSeconds      sql.NullInt64
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use errgroup to run two queries concurrently
	g, ctx := errgroup.WithContext(ctx)

	// Fetch session times
	g.Go(func() error {
		query := `
			SELECT started_at, ended_at
			FROM task_sessions
			WHERE id = ?
		`
		return m.DB.QueryRowContext(ctx, query, sessionID).Scan(&startedAt, &endedAt)
	})

	// Fetch total active time from intervals
	g.Go(func() error {
		query := `
			SELECT SUM(strftime('%s', end_time) - strftime('%s', start_time))
			FROM task_session_intervals
			WHERE session_id = ? AND end_time IS NOT NULL
		`
		return m.DB.QueryRowContext(ctx, query, sessionID).Scan(&activeSeconds)
	})

	// Wait for both queries to finish
	if err := g.Wait(); err != nil {
		return 0, 0, 0, err
	}

	// Validation
	if endedAt.IsZero() {
		return 0, 0, 0, fmt.Errorf("session is still active; stop it first to calculate durations")
	}

	totalTime = endedAt.Sub(startedAt)
	activeTime = time.Duration(activeSeconds.Int64) * time.Second
	pausedTime = totalTime - activeTime

	return activeTime, pausedTime, totalTime, nil
}
