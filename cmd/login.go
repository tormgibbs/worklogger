/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/auth"
	"github.com/tormgibbs/worklogger/tui"
)

var githubAuthFlag bool

var LocalAuthFlag bool

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

		session, _ := auth.LoadSession()
		if session.Authenticated {
			fmt.Println("You've been authenticated already")
			return
		}

		if githubAuthFlag && LocalAuthFlag {
			fmt.Println("Error: You can't use both --github and --local")
			return
		}

		if githubAuthFlag {
			auth.SetToken("github_token", "ghp_mockedGithubToken")
			auth.SaveSession(auth.Session{
				Method:        "github",
				Authenticated: true,
			})

			fmt.Println("Logged in with GitHub")
			return
		}

		if LocalAuthFlag {
			var username, password string

			fmt.Print("Enter your username: ")
			fmt.Scanln(&username)

			fmt.Print("Enter your password: ")
			bytePass, _ := term.ReadPassword(uintptr(os.Stdin.Fd()))
			
			fmt.Println()

			password = string(bytePass)

			if err := auth.VerifyLocalCredentials(username, password); err != nil {
				fmt.Println("Login failed:", err)
				return
			}


			// auth.SetToken("local_token", "mockedLocalToken123")
			auth.SaveSession(auth.Session{
				Method:        "local",
				Authenticated: true,
			})

			fmt.Println("Logged in successfully with local credentials.")
			return
		}

		lm, err := tui.RunLoginUI()
		if err != nil {
			cmd.PrintErrln("Error running Login TUI", err)
			return
		}

		var token string

		switch lm.Selected {
		case "Github":
			token = "ghp_token_from_ui"
			auth.SetToken("github_token", token)
		case "Local Authentication":
			token = "local_token_from_ui"
			auth.SetToken("local_token", token)
		default:
			fmt.Println("No valid login method selected.")
			return
		}

		auth.SaveSession(auth.Session{Method: lm.Selected, Authenticated: true})
		fmt.Printf("Successfully logged in using: %s\n", lm.Selected)

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().BoolVarP(&githubAuthFlag, "github", "g", false, "Authenticate with Github")

	loginCmd.Flags().BoolVarP(&LocalAuthFlag, "local", "l", false, "Authenticate locally(email & password)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
