package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/data"
)

var (
	taskFlag string

	modeFlag  string
	tagFlags  []string
	kpiFlags  []string
	notesFlag string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new task session",
	Long: `Begin tracking a new task session with optional tags and KPIs.

This will create a new task and log the start time. If you already
have an active session, you'll need to stop or pause it first.

Example:
  worklogger start --task "Write documentation`,
	Run: func(cmd *cobra.Command, args []string) {
		if taskFlag == "" {
			fmt.Println("⚠️  No task provided. Use --task or -t to specify one.")
			return
		}

		if err := validateModeAndFlags(); err != nil {
			fmt.Printf("⚠️  %s\n", err.Error())
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

		session := &data.TaskSession{
			StartedAt: time.Now(),
			Mode:      getSessionMode(),
			Notes:     notesFlag,
			Synced:    false,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer tx.Rollback()

		if err := models.CreateTask(tx, task, session); err != nil {
			cmd.PrintErr(err)
			return
		}

		if len(tagFlags) > 0 {
			if err := models.SessionTags.Create(tx, session.ID, tagFlags); err != nil {
				cmd.PrintErr(err)
				return
			}
		}

		if len(kpiFlags) > 0 {
			if err := models.SessionKPI.Create(tx, session.ID, kpiFlags); err != nil {
				cmd.PrintErr(err)
				return
			}
		}

		if err := tx.Commit(); err != nil {
			cmd.PrintErr(err)
			return
		}

		formattedTime := task.CreatedAt.Format("2006-01-02 15:04:05")
		fmt.Printf("Session started at %s for task: %s\n", formattedTime, task.Description)

		if session.Mode == "org" {
			fmt.Printf("   Mode: %s\n", session.Mode)
		}
		if len(tagFlags) > 0 {
			fmt.Printf("   Tags: %s\n", strings.Join(tagFlags, ", "))
		}
		if len(kpiFlags) > 0 {
			fmt.Printf("   KPIs: %s\n", strings.Join(kpiFlags, ", "))
		}
		if notesFlag != "" {
			fmt.Printf("   Notes: %s\n", notesFlag)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&taskFlag, "task", "t", "", "Task description to start logging")

	startCmd.Flags().StringVar(&modeFlag, "mode", "", "Session mode: personal or org")
	startCmd.Flags().StringSliceVar(&tagFlags, "tag", nil, "Tags to associate with the session")
	startCmd.Flags().StringSliceVar(&kpiFlags, "kpi", nil, "KPIs to associate with the session (required for org mode)")
	startCmd.Flags().StringVar(&notesFlag, "notes", "", "Notes for this session")
}

func validateModeAndFlags() error {
	mode := getSessionMode()

	if mode == "org" && len(kpiFlags) == 0 {
		return fmt.Errorf("organization mode requires at least one KPI (use --kpi)")
	}

	if mode != "personal" && mode != "org" {
		return fmt.Errorf("mode must be 'personal' or 'org', got: %s", mode)
	}

	return nil
}

func getSessionMode() string {
	if modeFlag != "" {
		return modeFlag
	}

	// TODO: Read from config file (.orcta.yaml)
	// For now, default to personal
	return "personal"
}
