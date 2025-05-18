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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resumeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resumeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
