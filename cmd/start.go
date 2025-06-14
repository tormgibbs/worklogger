package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/data"
)

var taskFlag string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new task session",
	Long: `Begin tracking a new task session.

This will create a new task and log the start time. If you already
have an active session, you'll need to stop or pause it first.

Example:
  worklogger start --task "Write documentation`,
	Run: func(cmd *cobra.Command, args []string) {
		if taskFlag == "" {
			fmt.Println("⚠️  No task provided. Use --task or -t to specify one.")
			return
		}

		ts, err := models.TaskSessions.Get()
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to check active task: %w", err))
			fmt.Println()
			return
		}

		if ts != nil {
			fmt.Println("You already have an active task. Finish that one before starting a new one.")
			return
		}

		description := taskFlag

		task := &data.Task{
			Description: description,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if err := models.CreateTask(tx, task); err != nil {
			cmd.PrintErr(err)
			return
		}

		formattedTime := task.CreatedAt.Format("2006-01-02 15:04:05")
		fmt.Printf("Session started at %s for task: %s\n", formattedTime, task.Description)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&taskFlag, "task", "t", "", "Task description to start logging")

}
