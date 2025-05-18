/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := models.TaskSessions.Get()
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to check active session: %w", err))
			fmt.Println()
			return
		}

		if ts == nil {
			fmt.Println("No active session to stop.")
			return
		}

		_, err = models.TaskSessionIntervals.End(ts)
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to close running interval: %w", err))
			fmt.Println()
			return
		}

		stoppedSession, err := models.TaskSessions.Stop(ts.ID)
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to stop session: %w", err))
			fmt.Println()
			return
		}

		if stoppedSession == nil {
			fmt.Println("Task session was already stopped.")
			return
		}

		active, paused, total, err := models.TaskSessions.GetDurations(ts.ID)
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to calculate session durations: %w", err))
			fmt.Println()
			return
		}

		formattedTime := stoppedSession.EndedAt.Format("2006-01-02 15:04:05")

		fmt.Printf("Session stopped at %v.\n", formattedTime)
		fmt.Println("Session summary:")
		fmt.Printf("  ⏱️  Total:  %v\n", total)
		fmt.Printf("  🟢 Active: %v\n", active)
		fmt.Printf("  🛑 Paused: %v\n", paused)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
