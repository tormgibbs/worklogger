/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// setupGithubCmd represents the setupGithub command
var setupGithubCmd = &cobra.Command{
	Use:     "setup-github",
	Aliases: []string{"setupGithub"},
	Short:   "Configure GitHub OAuth credentials",
	Long:    `Set up GitHub OAuth credentials for GitHub integration features.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupGitHubOAuth()
	},
}

func init() {
	rootCmd.AddCommand(setupGithubCmd)
}
