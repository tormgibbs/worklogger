package data

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type MetricStat struct {
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type SummaryStats struct {
	TodayHours        MetricStat `json:"today_hours"`
	WeekHours         MetricStat `json:"week_hours"`
	SessionsToday     MetricStat `json:"sessions_today"`
	ProductivityScore MetricStat `json:"productivity_score"`
}

type DailyStat struct {
	Date     string  `json:"date"`
	Hours    float64 `json:"hours"`
	Sessions int     `json:"sessions"`
}

type WeeklyStat struct {
	Start    string  `json:"week_start"`
	Hours    float64 `json:"hours"`
	Sessions int     `json:"sessions"`
}

type MonthlyStat struct {
	Month    string  `json:"month"`
	Hours    float64 `json:"hours"`
	Sessions int     `json:"sessions"`
}

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

func CreateTaskAndSession(db *sql.DB, description string) (*Task, *TaskSession, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't start transaction: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Step 1: Insert Task and get back ID + CreatedAt
	query := `
		INSERT INTO tasks (description)
		VALUES (?)
		RETURNING id, description, created_at
	`
	task := &Task{}
	err = tx.QueryRowContext(ctx, query, description).Scan(&task.ID, &task.Description, &task.CreatedAt)

	if err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to insert task: %w", err)
	}

	query = `
		INSERT INTO task_sessions (task_id)
		VALUES (?)
		RETURNING id, task_id, started_at, ended_at
	`

	// Step 2: Insert Task Session and get back full info
	session := &TaskSession{}
	err = tx.QueryRowContext(ctx, query, task.ID).Scan(&session.ID, &session.TaskID, &session.StartedAt, &session.EndedAt)

	if err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to insert task session: %w", err)
	}

	// Finalize the transaction
	if err = tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("commit failed: %w", err)
	}

	return task, session, nil
}

func GetSummaryStats(db *sql.DB) (*SummaryStats, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var stats SummaryStats
	var firstErr error

	wg.Add(4)

	go func() {
		defer wg.Done()
		val, change, err := GetTodayHours(db)
		if err != nil {
			setErr(&firstErr, err)
			return
		}
		mu.Lock()
		stats.TodayHours = MetricStat{Value: val, Change: change}
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		val, change, err := GetWeekHours(db)
		if err != nil {
			setErr(&firstErr, err)
			return
		}
		mu.Lock()
		stats.WeekHours = MetricStat{Value: val, Change: change}
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		val, change, err := GetTodaySessions(db)
		if err != nil {
			setErr(&firstErr, err)
			return
		}
		mu.Lock()
		stats.SessionsToday = MetricStat{Value: float64(val), Change: change}
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		val, change, err := GetProductivityScore(db)
		if err != nil {
			setErr(&firstErr, err)
			return
		}
		mu.Lock()
		stats.ProductivityScore = MetricStat{Value: val, Change: change}
		mu.Unlock()
	}()

	wg.Wait()
	return &stats, firstErr

}

func GetTodayHours(db *sql.DB) (float64, float64, error) {
	var today float64
	var yesterday float64

	query := `
		SELECT COALESCE(SUM((strftime('%s', end_time) - strftime('%s', start_time)) / 3600.0), 0)
		FROM task_session_intervals
		WHERE DATE(start_time) = DATE('now', 'localtime') AND end_time IS NOT NULL
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.QueryRowContext(ctx, query).Scan(&today)
	if err != nil {
		return 0, 0, err
	}

	query = `
		SELECT COALESCE(SUM((strftime('%s', end_time) - strftime('%s', start_time)) / 3600.0), 0)
		FROM task_session_intervals
		WHERE DATE(start_time) = DATE('now', '-1 day', 'localtime') AND end_time IS NOT NULL
	`
	err = db.QueryRowContext(ctx, query).Scan(&yesterday)
	if err != nil {
		return 0, 0, err
	}

	return today, calculateChange(today, yesterday), nil
}

func GetWeekHours(db *sql.DB) (float64, float64, error) {
	var currentWeek float64
	var previousWeek float64

	query := `
		SELECT COALESCE(SUM((strftime('%s', end_time) - strftime('%s', start_time)) / 3600.0), 0)
		FROM task_session_intervals
		WHERE start_time >= DATE('now', 'weekday 1', '-6 days', 'localtime')
		AND end_time IS NOT NULL
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.QueryRowContext(ctx, query).Scan(&currentWeek)
	if err != nil {
		return 0, 0, err
	}

	query = `
		SELECT COALESCE(SUM((strftime('%s', end_time) - strftime('%s', start_time)) / 3600.0), 0)
		FROM task_session_intervals
		WHERE start_time >= DATE('now', 'weekday 1', '-13 days', 'localtime')
		AND start_time < DATE('now', 'weekday 1', '-7 days', 'localtime')
		AND end_time IS NOT NULL
	`
	err = db.QueryRowContext(ctx, query).Scan(&previousWeek)
	if err != nil {
		return 0, 0, err
	}

	return currentWeek, calculateChange(currentWeek, previousWeek), nil
}

func GetTodaySessions(db *sql.DB) (int, float64, error) {
	var currentDay int
	var previousDay int

	query := `
		SELECT COUNT(*)
		FROM task_sessions
		WHERE DATE(started_at) = DATE('now', 'localtime')
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.QueryRowContext(ctx, query).Scan(&currentDay)
	if err != nil {
		return 0, 0, err
	}

	query = `
		SELECT COUNT(*)
		FROM task_sessions
		WHERE DATE(started_at) = DATE('now', '-1 day', 'localtime')
	`
	err = db.QueryRowContext(ctx, query).Scan(&previousDay)
	if err != nil {
		return 0, 0, err
	}

	return currentDay, calculateChange(float64(currentDay), float64(previousDay)), nil
}

func GetProductivityScore(db *sql.DB) (float64, float64, error) {
	var score float64
	var lastScore float64

	query := `
		WITH intervals AS (
			SELECT (strftime('%s', end_time) - strftime('%s', start_time)) AS duration_sec
			FROM task_session_intervals
			WHERE DATE(start_time) = DATE('now', 'localtime') AND end_time IS NOT NULL
		),
		productive AS (
			SELECT SUM(duration_sec) AS total_productive FROM intervals WHERE duration_sec >= 1200
		),
		totals AS (
			SELECT SUM(duration_sec) AS total_time, (SELECT total_productive FROM productive) AS productive_time FROM intervals
		)
		SELECT
			CASE WHEN total_time > 0 THEN ROUND(100.0 * productive_time / total_time, 2) ELSE 0 END
		FROM totals;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.QueryRowContext(ctx, query).Scan(&score)
	if err != nil {
		return 0, 0, err
	}

	query = `
		WITH intervals AS (
			SELECT (strftime('%s', end_time) - strftime('%s', start_time)) AS duration_sec
			FROM task_session_intervals
			WHERE DATE(start_time) = DATE('now', '-1 day', 'localtime') AND end_time IS NOT NULL
		),
		productive AS (
			SELECT SUM(duration_sec) AS total_productive FROM intervals WHERE duration_sec >= 1200
		),
		totals AS (
			SELECT SUM(duration_sec) AS total_time, (SELECT total_productive FROM productive) AS productive_time FROM intervals
		)
		SELECT
			CASE WHEN total_time > 0 THEN ROUND(100.0 * productive_time / total_time, 2) ELSE 0 END
		FROM totals;
	`
	err = db.QueryRowContext(ctx, query).Scan(&lastScore)
	if err != nil {
		return 0, 0, err
	}

	return score, calculateChange(score, lastScore), nil

}

func calculateChange(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return ((current - previous) / previous) * 100
}

func setErr(dst *error, err error) {
	if *dst == nil {
		*dst = err
	}
}

func GetDailyStats(db *sql.DB) ([]*DailyStat, error) {
	query := `
		SELECT 
			DATE(start_time, 'localtime') as period,
			COUNT(DISTINCT session_id) as sessions,
			SUM((strftime('%s', COALESCE(end_time, CURRENT_TIMESTAMP)) - strftime('%s', start_time)) / 3600.0) as hours
		FROM task_session_intervals
		WHERE start_time >= DATE('now', '-6 days')
		GROUP BY period
		ORDER BY period;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var stats []*DailyStat = make([]*DailyStat, 0)

	for rows.Next() {
		var stat DailyStat
		if err := rows.Scan(&stat.Date, &stat.Sessions, &stat.Hours); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		stats = append(stats, &stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return stats, nil
}

func GetWeeklyStats(db *sql.DB) ([]*WeeklyStat, error) {
	query := `
		SELECT 
			strftime('%Y-%W', start_time, 'localtime') as period,
			COUNT(DISTINCT session_id) as sessions,
			SUM((strftime('%s', COALESCE(end_time, CURRENT_TIMESTAMP)) - strftime('%s', start_time)) / 3600.0) as hours
		FROM task_session_intervals
		WHERE start_time >= DATE('now', '-28 days')
		GROUP BY period
		ORDER BY period;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var stats []*WeeklyStat = make([]*WeeklyStat, 0)

	for rows.Next() {
		var ws WeeklyStat
		var period string
		if err := rows.Scan(&period, &ws.Sessions, &ws.Hours); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		// Convert the period 'YYYY-WW' into the date of the week start (Monday)
		weekStart, err := parseWeekStart(period)
		if err != nil {
			return nil, fmt.Errorf("failed to parse week start: %w", err)
		}
		ws.Start = weekStart.Format("2006-01-02")

		stats = append(stats, &ws)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return stats, nil
}

func GetMonthlyStats(db *sql.DB) ([]*MonthlyStat, error) {
	query := `
		SELECT 
			strftime('%Y-%m', start_time, 'localtime') as period,
			COUNT(DISTINCT session_id) as sessions,
			SUM((strftime('%s', COALESCE(end_time, CURRENT_TIMESTAMP)) - strftime('%s', start_time)) / 3600.0) as hours
		FROM task_session_intervals
		WHERE start_time >= DATE('now', '-3 months')
		GROUP BY period
		ORDER BY period;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var stats []*MonthlyStat = make([]*MonthlyStat, 0)

	for rows.Next() {
		var ms MonthlyStat
		if err := rows.Scan(&ms.Month, &ms.Sessions, &ms.Hours); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		stats = append(stats, &ms)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return stats, nil

}

// period format is "YYYY-WW" where WW is week number according to strftime
func parseWeekStart(period string) (time.Time, error) {
	// period example: "2025-23" (year-weeknumber)
	var year, week int
	_, err := fmt.Sscanf(period, "%4d-%2d", &year, &week)
	if err != nil {
		return time.Time{}, err
	}

	// ISO 8601 weeks start on Monday.
	// We get the first Monday of the year, then add (week-1)*7 days.
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.Local)
	// Get the Monday of that week
	weekday := int(jan4.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	monday := jan4.AddDate(0, 0, -(weekday - 1))
	weekStart := monday.AddDate(0, 0, (week-1)*7)

	return weekStart, nil
}
