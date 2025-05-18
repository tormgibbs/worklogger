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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
