/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// resumeCmd represents the resume command
var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a paused work session",
	Long:  `Resumes a previously paused task session if one exists.`,
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := models.TaskSessions.Get()
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to check active task: %w", err))
			fmt.Println()
			return
		}

		if ts == nil {
			fmt.Println("You don't have an active session. Create a new task!")
			return
		}

		running, err := models.TaskSessionIntervals.HasOpenInterval(ts.ID)
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to check session intervals: %w", err))
			fmt.Println()
			return
		}

		if running {
			fmt.Println("Session is already running. Nothing to resume.")
			return
		}

		tsi, err := models.TaskSessionIntervals.StartNew(ts.ID)
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to resume task: %w", err))
			fmt.Println()
			return
		}

		fmt.Printf("Session resumed at %v\n", tsi.StartTime.Format("2006-01-02 15:04:05"))
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}
