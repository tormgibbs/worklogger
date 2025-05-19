package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/data"
	"github.com/tormgibbs/worklogger/tui"
)

var (
	sessionID         int
	taskDescription   string
	createNewSession  bool
	leaveUnassociated bool
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync Git commits with a worklogger session",
	Long: `Sync local Git commits to a worklogger task session.

Examples:
  worklogger sync --new -d "Fix login bug"
  worklogger sync --existing 12
  worklogger sync --unassociated

If no flag is passed, a prompt will guide you through the sync process.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if createNewSession && taskDescription == "" {
			return fmt.Errorf("when using --new, you must also provide --desc")
		}
		if createNewSession && sessionID > 0 {
			return fmt.Errorf("--new and --existing can't be used together")
		}
		if leaveUnassociated {
			if createNewSession || sessionID > 0 || taskDescription != "" {
				return fmt.Errorf("--unassociated can't be used with --new, --existing, or --desc")
			}
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		currentSession, err := models.TaskSessions.Get()
		if err != nil {
			cmd.PrintErrf("failed to check active session: %v\n", err)
			return
		}

		switch {
		case currentSession != nil:
			handleActiveSessionSync(cmd, currentSession)

		case createNewSession:
			if taskDescription == "" {
				cmd.Println("When using --new, you must also provide --desc (or -d).")
				return
			}
			handleNewSessionSync(cmd)

		case sessionID > 0:
			handleExistingSessionSync(cmd, sessionID)

		case leaveUnassociated:
			handleUnassociatedSync(cmd)

		default:
			runInteractiveSync(cmd)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().IntVarP(&sessionID, "existing", "e", 0, "Associate with existing session")
	syncCmd.Flags().StringVarP(&taskDescription, "desc", "d", "", "Description for new task")
	syncCmd.Flags().BoolVarP(&createNewSession, "new", "n", false, "Create new session")
	syncCmd.Flags().BoolVarP(&leaveUnassociated, "unassociated", "u", false, "Leave commits unassociated")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func handleActiveSessionSync(cmd *cobra.Command, session *data.TaskSession) {
	commits, err := models.Commits.SyncCommits(&session.ID)
	if err != nil {
		cmd.PrintErrf("Failed to sync commits :%v\n", err)
		return
	}
	fmt.Printf("Synced %d new commits\n", commits)
}

func handleNewSessionSync(cmd *cobra.Command) {
	if taskDescription == "" {
		cmd.Println("Task description cannot be empty. Use --desc or -d to provide one.")
		return
	}

	_, newSession, err := data.CreateTaskAndSession(db, taskDescription)
	if err != nil {
		cmd.PrintErrf("Failed to create new session: %v\n", err)
		return
	}

	commits, err := models.Commits.SyncCommits(&newSession.ID)
	if err != nil {
		cmd.PrintErrf("Failed to sync commits: %v\n", err)
		return
	}

	fmt.Printf("🎉 Created new session #%d and synced %d commits\n", newSession.ID, commits)
}

func handleExistingSessionSync(cmd *cobra.Command, id int) {
	session, err := models.TaskSessions.GetByID(id)
	if err != nil {
		cmd.PrintErrf("Could not find session ID %d: %v\n", id, err)
		return
	}

	commits, err := models.Commits.SyncCommits(&session.ID)
	if err != nil {
		cmd.PrintErrf("Failed to sync commits:%v\n", err)
		return
	}

	fmt.Printf("✅ Synced %d new commits to session #%d\n", commits, session.ID)
}

func handleUnassociatedSync(cmd *cobra.Command) {
	var nilSessionID *int
	commits, err := models.Commits.SyncCommits(nilSessionID)
	if err != nil {
		cmd.PrintErrf("Failed to sync commits:%v\n", err)
		return
	}
	fmt.Printf("\n✅ Synced %d unassociated commits\n", commits)
}

func runInteractiveSync(cmd *cobra.Command) {
	model, err := tui.RunSyncUI("Pick a Sync Option")
	if err != nil {
		cmd.PrintErrf("Error running sync TUI:%v\n", err)
		return
	}

	switch model.Selected {
	case tui.SyncOptionExisting:
		sessions, err := models.TaskSessions.GetAllWithTask()
		if err != nil {
			cmd.PrintErrf("Failed to get sessions:%v\n", err)
			return
		}
		if len(sessions) == 0 {
			cmd.Println("\n⚠️  No sessions found to associate.")
			return
		}

		selectedSession, err := tui.RunTaskSelectUI("Select a task to associate with these commits", sessions)
		if err != nil {
			cmd.PrintErrf("Error selecting session:%v\n", err)
			return
		}

		commits, err := models.Commits.SyncCommits(&selectedSession.ID)
		if err != nil {
			cmd.PrintErrf("Failed to sync commits:%v\n", err)
			return
		}
		fmt.Printf("\n✅ Synced %d commits\n", commits)

	case tui.SyncOptionNew:
		desc, err := tui.RunNewTaskUI()
		if err != nil {
			cmd.PrintErrf("Error running new task TUI:%v\n", err)
			return
		}
		if desc == "" {
			cmd.Println("Task description cannot be empty.")
			return
		}

		_, newSession, err := data.CreateTaskAndSession(db, desc)
		if err != nil {
			cmd.PrintErrf("Failed to create new session:%v\n", err)
			return
		}

		commits, err := models.Commits.SyncCommits(&newSession.ID)
		if err != nil {
			cmd.PrintErrf("Failed to sync commits:%v\n", err)
			return
		}

		fmt.Printf("🎉 Created new session #%d and synced %d commits\n", newSession.ID, commits)

	case tui.SyncOptionUnassociated:
		handleUnassociatedSync(cmd)

	case tui.SyncOptionCancel:
		fmt.Println("\nSyncing cancelled")

	default:
		fmt.Println("⚠️  No valid option selected")
	}
}
