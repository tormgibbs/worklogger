package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentSession, err := models.TaskSessions.Get()
		if err != nil {
			cmd.PrintErr(fmt.Errorf("failed to check active session: %w", err))
			fmt.Println()
			return
		}

		// syncing with current active session
		if currentSession != nil {
			commits, err := models.Commits.SyncCommits(&currentSession.ID)
			if err != nil {
				cmd.PrintErr(fmt.Errorf("failed to sync commits: %w", err))
				fmt.Println()
				return
			}

			fmt.Printf("Synced %d new commits\n", commits)
			return
		}

		if createNewSession {
			if currentSession != nil {
				cmd.Printf("Cannot create a new session while another session is active. Finish it first!")
				return
			}

			// new session flow here
			description := taskDescription
			if description == "" {
				fmt.Print("📝 Enter a description for the new task: ")
				_, err = fmt.Scanln(&description)
				if err != nil || description == "" {
					cmd.Println("❌ Task description cannot be empty.")
					return
				}
			}
		}

		if sessionID > 0 {
			// Make sure the session exists
			session, err := models.TaskSessions.GetByID(sessionID)
			if err != nil {
				cmd.PrintErrf("Could not find session ID %d: %v\n", sessionID, err)
				return
			}

			commits, err := models.Commits.SyncCommits(&session.ID)
			if err != nil {
				cmd.PrintErr(fmt.Errorf("failed to sync commits: %w", err))
				fmt.Println()
				return
			}

			fmt.Printf("✅ Synced %d new commits to session #%d\n", commits, session.ID)
			return
		}

		if leaveUnassociated {
			var sessionID *int
			commits, err := models.Commits.SyncCommits(sessionID)
			if err != nil {
				cmd.PrintErr(fmt.Errorf("failed to sync commits: %w", err))
				fmt.Println()
				return
			}
			fmt.Printf("Synced %d new commits\n", commits)
			return
		}

		model, err := tui.RunSyncUI("Syncing prompt")
		if err != nil {
			cmd.PrintErrf("Error running Login TUI: %v\n", err)
			return
		}

		switch model.Selected {
		case tui.SyncOptionExisting:
			// GET SESSIONS
			sessions, err := models.TaskSessions.GetAllWithTask()
			if err != nil {
				cmd.PrintErr(fmt.Errorf("failed to get sessions: %w", err))
				fmt.Println()
				return
			}

			if len(sessions) == 0 {
				cmd.Println("\n⚠️ No sessions found to associate.")
				return
			}

			// Run the TUI to select a session
			selectedSession, err := tui.RunTaskSelectUI("Select a task to associate with these commits", sessions)
			if err != nil {
				cmd.PrintErrf("Error selecting task: %v\n", err)
				return
			}

			commits, err := models.Commits.SyncCommits(&selectedSession.ID)
			if err != nil {
				cmd.PrintErr(fmt.Errorf("failed to sync commits: %w", err))
				fmt.Println()
				return
			}

			fmt.Printf("Synced %d new commits\n", commits)
		case tui.SyncOptionUnassociated:
			var sessionID *int // null session id

			commits, err := models.Commits.SyncCommits(sessionID)
			if err != nil {
				cmd.PrintErr(fmt.Errorf("failed to sync commits: %w", err))
				fmt.Println()
				return
			}

			fmt.Printf("Synced %d new commits\n", commits)
		case tui.SyncOptionCancel:
			fmt.Println("\nCancelling")
		default:
			fmt.Println("\nNo valid option selected")
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
