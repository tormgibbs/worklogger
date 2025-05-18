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
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		// fmt.Printf("Commit:\n")
		// fmt.Printf("Hash: %s\n", hashFlag)
		// fmt.Printf("Message: %s\n", messageFlag)
		// fmt.Printf("Author: %s\n", authorFlag)
		// fmt.Printf("Date: %s\n", dateFlag)
	},
}

func init() {
	rootCmd.AddCommand(readCommitCmd)

	readCommitCmd.Flags().StringVar(&hashFlag, "hash", "", "Git commit hash")
	readCommitCmd.Flags().StringVar(&messageFlag, "message", "", "Git commit message")
	readCommitCmd.Flags().StringVar(&authorFlag, "author", "", "Commit author")
	readCommitCmd.Flags().StringVar(&dateFlag, "date", "", "Commit date")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCommitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCommitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
