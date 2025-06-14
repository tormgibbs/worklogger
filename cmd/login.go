package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/auth"
	"github.com/tormgibbs/worklogger/config"
	"github.com/tormgibbs/worklogger/tui"
)

var githubLoginFlag bool

var localLoginFlag bool

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log into the app using GitHub or local credentials",
	Long: `Use the login command to authenticate with the application.

You can choose a login method by passing a flag:
  -g, --github   Log in using your GitHub account
  -l, --local    Log in using local credentials

Examples:
  worklogger login             # Prompts you to choose a login method
  worklogger login --github    # Logs in with GitHub
  worklogger login -l          # Logs in with local credentials

If no method is provided, a list of available options will be shown.`,

	Run: func(cmd *cobra.Command, args []string) {

		session, err := auth.LoadSession()
		if err != nil {
			cmd.PrintErrf("Failed to load session: %v\n", err)
			return
		}

		if session.Authenticated {
			fmt.Println("You've been authenticated already")
			return
		}

		if githubLoginFlag && localLoginFlag {
			fmt.Println("Error: You can't use both --github and --local")
			return
		}

		if githubLoginFlag {
			if err := auth.StartGitHubOAuth(
				config.Github.ClientID,
				config.Github.ClientSecret,
				config.Github.RedirectURI,
			); err != nil {
				cmd.PrintErrf("GitHub OAuth failed: %v\n", err)
			}
			return
		}

		if localLoginFlag {
			if err := auth.LocalLogin(); err != nil {
				cmd.PrintErrf("Local Login failed: %v\n", err)
			}
			return
		}

		lm, err := tui.RunAuthUI("🔐 Choose a login option")
		if err != nil {
			cmd.PrintErrf("Error running Login TUI: %v\n", err)
			return
		}

		switch lm.Selected {
		case "GitHub":
			if err := auth.StartGitHubOAuth(
				config.Github.ClientID,
				config.Github.ClientSecret,
				config.Github.RedirectURI,
			); err != nil {
				cmd.PrintErrf("GitHub OAuth failed: %v\n", err)
				return
			}
			fmt.Println("GitHub Login Successful")
		case "Local":
			if err := auth.LocalLogin(); err != nil {
				cmd.PrintErrf("Local Login failed: %v\n", err)
				return
			}
		default:
			fmt.Println("No login method selected. Exiting.")
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().BoolVarP(&githubLoginFlag, "github", "g", false, "Authenticate with Github")

	loginCmd.Flags().BoolVarP(&localLoginFlag, "local", "l", false, "Authenticate locally (email & password)")
}
