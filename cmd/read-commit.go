/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/data"
)

var (
	hashFlag, messageFlag, authorFlag, dateFlag string
)

// readCommitCmd represents the readCommit command
var readCommitCmd = &cobra.Command{
	Use:     "read-commit",
	Aliases: []string{"readCommit"},
	Short:   "Record a Git commit into the database",
	Long: `Record a Git commit to your work session.

Use this when integrating with Git hooks or external tools
to store commit history in WorkLogger.

All flags are required:
  --hash      Commit hash
  --message   Commit message
  --author    Author of the commit
  --date      Commit date (ISO 8601 or compatible format)`,
	Run: func(cmd *cobra.Command, args []string) {

		if hashFlag == "" || messageFlag == "" || authorFlag == "" || dateFlag == "" {
			fmt.Println("Missing required commit fields")
			return
		}

		ts, err := models.TaskSessions.Get()
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to check active session: %w", err))
			fmt.Println()
			return
		}

		var sessionID *int // nil by default
		if ts != nil {
			sessionID = &ts.ID // Only set if ts is non-nil
		}

		commit := &data.Commit{
			Hash:      hashFlag,
			SessionID: sessionID,
			Message:   messageFlag,
			Author:    authorFlag,
			Date:      dateFlag,
		}

		if err := models.Commits.Create(commit); err != nil {
			fmt.Printf("Failed to insert commit: %v\n", err)
			return
		}

		fmt.Printf("Commit %s recorded\n", commit.Hash)
	},
}

func init() {
	rootCmd.AddCommand(readCommitCmd)

	readCommitCmd.Flags().StringVar(&hashFlag, "hash", "", "Git commit hash")
	readCommitCmd.Flags().StringVar(&messageFlag, "message", "", "Git commit message")
	readCommitCmd.Flags().StringVar(&authorFlag, "author", "", "Commit author")
	readCommitCmd.Flags().StringVar(&dateFlag, "date", "", "Commit date")
}
