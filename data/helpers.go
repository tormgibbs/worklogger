package data

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"
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

type Session struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Duration  string `json:"duration"`
	Status    string `json:"status"`
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

	todayQuery := `
		SELECT COALESCE(SUM(
			(strftime('%s', MIN(COALESCE(end_time, CURRENT_TIMESTAMP), DATETIME('now', 'start of day', '+1 day', 'localtime')))
			- strftime('%s', MAX(start_time, DATETIME('now', 'start of day', 'localtime')))) / 3600.0
		), 0)
		FROM task_session_intervals
		WHERE
			start_time < DATETIME('now', 'start of day', '+1 day', 'localtime') AND
			COALESCE(end_time, CURRENT_TIMESTAMP) > DATETIME('now', 'start of day', 'localtime')
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.QueryRowContext(ctx, todayQuery).Scan(&today)
	if err != nil {
		return 0, 0, err
	}

	yesterdayQuery := `
		SELECT COALESCE(SUM(
			(strftime('%s', MIN(end_time, DATETIME('now', 'start of day')))
			- strftime('%s', MAX(start_time, DATETIME('now', '-1 day', 'start of day')))) / 3600.0
		), 0)
		FROM task_session_intervals
		WHERE
			start_time < DATETIME('now', 'start of day') AND
			end_time > DATETIME('now', '-1 day', 'start of day')
	`
	err = db.QueryRowContext(ctx, yesterdayQuery).Scan(&yesterday)
	if err != nil {
		return 0, 0, err
	}

	return math.Round(today), calculateChange(today, yesterday), nil
}

func GetWeekHours(db *sql.DB) (float64, float64, error) {
	var currentWeek float64
	var previousWeek float64

	currentWeekQuery := `
		SELECT COALESCE(SUM(
				(strftime('%s', COALESCE(end_time, DATETIME('now'))) 
				- strftime('%s', start_time)) / 3600.0
		), 0)
		FROM task_session_intervals
		WHERE 
				start_time < DATETIME('now', 'weekday 1', 'start of day') AND
				(end_time IS NULL OR end_time >= DATETIME('now', 'weekday 1', '-7 days', 'start of day'));
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.QueryRowContext(ctx, currentWeekQuery).Scan(&currentWeek)
	if err != nil {
		return 0, 0, err
	}

	previousWeekQuery := `
		SELECT COALESCE(SUM(
				(strftime('%s', COALESCE(end_time, DATETIME('now', 'weekday 1', '-7 days', 'start of day'))) 
				- strftime('%s', start_time)) / 3600.0
		), 0)
		FROM task_session_intervals
		WHERE 
				start_time < DATETIME('now', 'weekday 1', '-7 days', 'start of day') AND
				(end_time IS NULL OR end_time >= DATETIME('now', 'weekday 1', '-14 days', 'start of day'));
	`
	err = db.QueryRowContext(ctx, previousWeekQuery).Scan(&previousWeek)
	if err != nil {
		return 0, 0, err
	}

	return math.Round(currentWeek), calculateChange(currentWeek, previousWeek), nil
}

func GetTodaySessions(db *sql.DB) (int, float64, error) {
	var currentDay int
	var previousDay int

	query := `
		SELECT COUNT(*)
		FROM task_sessions
		WHERE
			started_at < DATETIME('now', '+1 day', 'start of day') AND
			COALESCE(ended_at, CURRENT_TIMESTAMP) >= DATETIME('now', 'start of day')
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
		WHERE
			started_at < DATETIME('now', 'start of day') AND
			COALESCE(ended_at, CURRENT_TIMESTAMP) >= DATETIME('now', '-1 day', 'start of day')
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
			SELECT
				(strftime('%s', MIN(end_time, DATETIME('now', '+1 day', 'start of day', 'localtime')))
				- strftime('%s', MAX(start_time, DATETIME('now', 'start of day', 'localtime')))
				) AS duration_sec
			FROM task_session_intervals
			WHERE
				start_time < DATETIME('now', '+1 day', 'start of day', 'localtime') AND
				end_time > DATETIME('now', 'start of day', 'localtime')
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
			SELECT
				(strftime('%s', MIN(end_time, DATETIME('now', 'start of day', 'localtime')))
				- strftime('%s', MAX(start_time, DATETIME('now', '-1 day', 'start of day', 'localtime')))
				) AS duration_sec
			FROM task_session_intervals
			WHERE
				start_time < DATETIME('now', 'start of day', 'localtime') AND
				end_time > DATETIME('now', '-1 day', 'start of day', 'localtime')
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

	change := ((current - previous) / previous) * 100
	return math.Round(change)
}

func setErr(dst *error, err error) {
	if *dst == nil {
		*dst = err
	}
}

func GetDailyStats(db *sql.DB) ([]*DailyStat, error) {
	query := `
		WITH RECURSIVE days(day) AS (
				SELECT DATE('now', '-6 days')
				UNION ALL
				SELECT DATE(day, '+1 day')
				FROM days
				WHERE day < DATE('now')
		),
		intervals AS (
				SELECT 
						DATE(day) AS period,
						session_id,
						(strftime('%s', MIN(COALESCE(end_time, CURRENT_TIMESTAMP), DATETIME(day, '+1 day'))) -
						strftime('%s', MAX(start_time, DATETIME(day)))) / 3600.0 AS duration
				FROM task_session_intervals
				JOIN days
				ON start_time < DATETIME(days.day, '+1 day') AND COALESCE(end_time, CURRENT_TIMESTAMP) > days.day
		)
		SELECT 
				period,
				COUNT(DISTINCT session_id) AS sessions,
				SUM(duration) AS hours
		FROM intervals
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

	stats := make([]*DailyStat, 0)

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
		WITH RECURSIVE weeks(start_date) AS (
			SELECT DATE('now', 'weekday 1', '-21 days')  -- Start from 3 weeks ago Monday
			UNION ALL
			SELECT DATE(start_date, '+7 days')
			FROM weeks
			WHERE start_date < DATE('now', 'weekday 1')
		),
		intervals AS (
			SELECT 
				weeks.start_date AS week_start,
				session_id,
				(strftime('%s', MIN(COALESCE(end_time, CURRENT_TIMESTAMP), DATETIME(weeks.start_date, '+7 days'))) -
				 strftime('%s', MAX(start_time, DATETIME(weeks.start_date)))) / 3600.0 AS duration
			FROM task_session_intervals
			JOIN weeks
				ON start_time < DATETIME(weeks.start_date, '+7 days')
				AND COALESCE(end_time, CURRENT_TIMESTAMP) > weeks.start_date
		)
		SELECT 
			week_start,
			COUNT(DISTINCT session_id) AS sessions,
			ROUND(SUM(duration), 2) AS hours
		FROM intervals
		GROUP BY week_start
		ORDER BY week_start;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	stats := make([]*WeeklyStat, 0)

	for rows.Next() {
		var ws WeeklyStat
		if err := rows.Scan(&ws.Start, &ws.Sessions, &ws.Hours); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		stats = append(stats, &ws)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return stats, nil
}

func GetMonthlyStats(db *sql.DB) ([]*MonthlyStat, error) {
	query := `
		WITH RECURSIVE months(month_start) AS (
			SELECT DATE('now', 'start of month', '-2 months', 'localtime')
			UNION ALL
			SELECT DATE(month_start, '+1 month')
			FROM months
			WHERE month_start < DATE('now', 'start of month', 'localtime')
		)
		SELECT
			month_start,
			COUNT(DISTINCT tsi.session_id) AS sessions,
			COALESCE(SUM(
				CAST(
					(
						strftime('%s', MIN(COALESCE(tsi.end_time, CURRENT_TIMESTAMP), DATE(month_start, '+1 month')))
						-
						strftime('%s', MAX(tsi.start_time, month_start))
					) AS REAL
				) / 3600.0
			), 0) AS hours
		FROM months
		LEFT JOIN task_session_intervals tsi ON
			tsi.start_time < DATE(month_start, '+1 month') AND
			COALESCE(tsi.end_time, CURRENT_TIMESTAMP) > month_start
		GROUP BY month_start
		ORDER BY month_start;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	stats := make([]*MonthlyStat, 0)

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

func GetSessions(db *sql.DB) ([]*Session, error) {
	query := `
		SELECT 
			ts.id,
			t.description,
			MIN(ti.start_time) AS start_time,
			MAX(ti.end_time) AS last_interval_end,
			ts.ended_at,
			EXISTS (
				SELECT 1 FROM task_session_intervals ti2 
				WHERE ti2.session_id = ts.id AND ti2.end_time IS NULL
			) AS has_active_interval,
			GROUP_CONCAT(
				(CASE 
					WHEN ti.end_time IS NOT NULL 
					THEN (JULIANDAY(ti.end_time) - JULIANDAY(ti.start_time)) * 86400
					ELSE (JULIANDAY(CURRENT_TIMESTAMP) - JULIANDAY(ti.start_time)) * 86400
				END)
			)
		FROM task_sessions ts
		JOIN tasks t ON ts.task_id = t.id
		JOIN task_session_intervals ti ON ti.session_id = ts.id
		GROUP BY ts.id
		ORDER BY start_time DESC
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessions := make([]*Session, 0)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			s         Session
			startTime string
			endTime   sql.NullString
			endedAt   sql.NullString
			hasActive bool
			durations string
		)

		err := rows.Scan(&s.ID, &s.Task, &startTime, &endTime, &endedAt, &hasActive, &durations)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		totalSeconds := 0
		for _, part := range strings.Split(durations, ",") {
			sec, _ := strconv.ParseFloat(part, 64)
			totalSeconds += int(sec)
		}

		hours := totalSeconds / 3600
		minutes := (totalSeconds % 3600) / 60

		if hours > 0 {
			s.Duration = fmt.Sprintf("%dh %dm", hours, minutes)
		} else {
			s.Duration = fmt.Sprintf("%dm", minutes)
		}

		s.StartTime = startTime
		s.EndTime = endTime.String

		// Determine session status
		if endedAt.Valid {
			s.Status = "ended"
		} else if hasActive {
			s.Status = "in_progress"
		} else {
			s.Status = "paused"
		}

		sessions = append(sessions, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return sessions, nil

}
