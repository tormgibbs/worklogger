/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/tui"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "View your logged work sessions",
	Long: `Launches an interactive TUI to browse your past work logs,
including tasks, durations, and timestamps.`,
	Run: func(cmd *cobra.Command, args []string) {

		logs, err := models.Logs.GetLogsWithDurations()
		if err != nil {
			cmd.PrintErrf("failed to get logs: %v\n", err)
			return
		}

		if _, err := tui.RunLogUI(logs); err != nil {
			cmd.PrintErrf("Failed to start log viewer: %v\n", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
