package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause the current active session",
	Long: `Pause an ongoing task session. 
This ends the current interval but keeps the session open so you can resume it later.`,
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := models.TaskSessions.Get()
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to check active session: %w", err))
			fmt.Println()
			return
		}

		if ts == nil {
			fmt.Println("You don't have an active session. Create a new session!")
			return
		}

		tsi, err := models.TaskSessionIntervals.End(ts)
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to pause session: %w", err))
			fmt.Println()
			return
		}

		if tsi == nil {
			fmt.Println("Session's already paused. Resume it first or stop it!")
			return
		}

		formattedTime := tsi.EndTime.Format("2006-01-02 15:04:05")
		fmt.Printf("Session paused at %v\n", formattedTime)
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}
