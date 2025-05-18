package data

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type CommitModel struct {
	DB *sql.DB
}

type Commit struct {
	ID        int
	SessionID *int
	Hash      string
	Message   string
	Author    string
	Date      string
}

// Create inserts a commit into the DB.
func (m CommitModel) Create(c *Commit) error {
	query := `
		INSERT OR IGNORE INTO commits (hash, session_id, message, author, date)
		VALUES (?, ?, ?, ?, ?)
	`
	args := []any{c.Hash, c.SessionID, c.Message, c.Author, c.Date}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)

	return err
}

// GetAllHashes returns a map of all commit hashes stored in the DB.
func (m CommitModel) GetAllHashes() (map[string]bool, error) {
	existing := make(map[string]bool)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT hash FROM commits"

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query commits: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var hash string
		if err := rows.Scan(&hash); err != nil {
			return nil, fmt.Errorf("failed to scan commit hash: %w", err)
		}
		existing[hash] = true
	}

	return existing, nil
}

// FetchGitCommits runs git log and returns a slice of Commit structs.
func FetchGitCommits(sessionID *int) ([]*Commit, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%H|%an|%ad|%s")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run git log: %w", err)
	}

	var commits []*Commit
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 4)
		if len(parts) < 4 {
			continue
		}
		commit := &Commit{
			Hash:      parts[0],
			SessionID: sessionID,
			Author:    parts[1],
			Date:      parts[2],
			Message:   parts[3],
		}
		commits = append(commits, commit)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return commits, nil
}

func (m CommitModel) SyncCommits(sessionID *int) (int, error) {
	existing, err := m.GetAllHashes()
	if err != nil {
		return 0, err
	}

	allCommits, err := FetchGitCommits(sessionID)
	if err != nil {
		return 0, err
	}

	var newCommits []*Commit
	for _, c := range allCommits {
		if !existing[c.Hash] {
			newCommits = append(newCommits, c)
		}
	}

	for _, c := range newCommits {
		if err := m.Create(c); err != nil {
			// Just log and continue to avoid failing the whole sync
			// You can also collect errors if you want
			fmt.Printf("Failed to insert commit %s: %v\n", c.Hash, err)
		}
	}

	return len(newCommits), nil
}
