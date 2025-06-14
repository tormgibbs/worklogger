package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/auth"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of the current session",
	Long: `Logs you out of WorkLogger by deleting your active session.

This will remove your authentication token and clear any
locally stored session data.`,
	Run: func(cmd *cobra.Command, args []string) {
		session, _ := auth.LoadSession()

		if !session.Authenticated {
			fmt.Println("You're not logged in.")
			return
		}

		switch session.Method {
		case "github":
			auth.DeleteToken("github_token")
		case "local":
			auth.DeleteToken("local_token")
		}

		err := auth.DeleteSession()
		if err != nil {
			cmd.PrintErrln("Failed to delete session:", err)
			return
		}

		fmt.Println("Successfully logged out.")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
