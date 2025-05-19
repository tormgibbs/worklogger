package data

import (
	"database/sql"
	"time"
)

type TaskModel struct {
	DB *sql.DB
}

type Task struct {
	ID          int
	Description string
	CreatedAt   time.Time
}

func (m TaskModel) CreateTx(tx *sql.Tx, task *Task) error {
	query := `
		INSERT INTO tasks (description)
		VALUES (?)
		RETURNING id, created_at
	`
	return tx.QueryRow(query, task.Description).Scan(&task.ID, &task.CreatedAt)
}


