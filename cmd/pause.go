package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause",
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

	// rootCmd.Flags().BoolVarP(&pauseFlag, "pause", "p", false, "Pause Task")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pauseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pauseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
