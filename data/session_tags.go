package data

import (
	"database/sql"
	"time"
)

type SessionTagModel struct {
	DB *sql.DB
}

type SessionTag struct {
	ID        int       `json:"id"`
	SessionID int       `json:"session_id"`
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
}

func (m SessionTagModel) Create(tx *sql.Tx, sessionID int, tags []string) error {
	query := `
		INSERT INTO session_tags (session_id, tag) 
		VALUES (?, ?)
	`
	for _, tag := range tags {
		_, err := tx.Exec(query, sessionID, tag)
		if err != nil {
			return err
		}
	}
	return nil
}
