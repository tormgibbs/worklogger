package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/auth"
	"github.com/tormgibbs/worklogger/config"
	"github.com/tormgibbs/worklogger/tui"
)

var githubSignUpFlag bool

var localSignUpFlag bool

// signupCmd represents the signup command
var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create a new account using GitHub or local credentials",
	Long: `Use the signup command to create a new account in the application.

You can choose a signup method by passing a flag:
  -g, --github   Sign up using your GitHub account
  -l, --local    Sign up using local credentials

Examples:
  worklogger signup             # Prompts you to choose a signup method
  worklogger signup --github    # Signs up with GitHub
  worklogger signup -l          # Signs up with local credentials

If no method is provided, an interactive prompt will guide you through the signup process.`,
	Run: func(cmd *cobra.Command, args []string) {
		session, err := auth.LoadSession()
		if err != nil {
			cmd.PrintErrf("Failed to load session: %v\n", err)
			return
		}

		if session.Authenticated {
			cmd.PrintErrf("You're already authenticated")
			return
		}

		if githubSignUpFlag && localSignUpFlag {
			cmd.PrintErrf("Error: You can't use both --github and --local")
			return
		}

		if githubSignUpFlag {
			err := auth.StartGitHubOAuth(
				config.Github.ClientID,
				config.Github.ClientSecret,
			)
			if err != nil {
				cmd.PrintErrf("GitHub OAuth failed: %v\n", err)
				return
			}
			fmt.Println("GitHub Signup Successful")
			return
		}

		if localSignUpFlag {
			if err := auth.LocalSignUp(); err != nil {
				cmd.PrintErrf("Local Signup failed: %v\n", err)
				return
			}
			fmt.Println("Local Signup Successful") // Add this line
			return
		}

		m, err := tui.RunAuthUI("🔐 Choose a signup option")
		if err != nil {
			cmd.PrintErrf("Error running Signup TUI: %v\n", err)
			return
		}

		switch m.Selected {
		case tui.GitHubOAuth:
			if err := auth.StartGitHubOAuth(
				config.Github.ClientID,
				config.Github.ClientSecret,
			); err != nil {
				cmd.PrintErrf("GitHub OAuth failed: %v\n", err)
				return
			}
			fmt.Println("GitHub Signup Successful")
			return
		case tui.LocalAuth:
			if err := auth.LocalSignUp(); err != nil {
				cmd.PrintErrf("Local Signup failed: %v\n", err)
				return
			}
			fmt.Println("Local Signup Successful")
			return
		default:
			fmt.Println("No signup method selected. Exiting.")
		}
	},
}

func init() {
	rootCmd.AddCommand(signupCmd)

	signupCmd.Flags().BoolVarP(&githubSignUpFlag, "github", "g", false, "Sign up using GitHub authentication")

	signupCmd.Flags().BoolVarP(&localSignUpFlag, "local", "l", false, "Sign up using local credentials")
}
