package data

import (
	"database/sql"
	"time"
)

type SessionKPIModel struct {
	DB *sql.DB
}

type SessionKPI struct {
	ID        int       `json:"id"`
	SessionID int       `json:"session_id"`
	KPI       string    `json:"kpi"`
	CreatedAt time.Time `json:"created_at"`
}

func (m SessionKPIModel) Create(tx *sql.Tx, sessionID int, kpis []string) error {
	query := `
		INSERT INTO session_kpis (session_id, kpi) 
		VALUES (?, ?)
	`
	for _, kpi := range kpis {
		_, err := tx.Exec(query, sessionID, kpi)
		if err != nil {
			return err
		}
	}
	return nil
}
