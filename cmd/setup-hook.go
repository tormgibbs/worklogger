/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// setupHookCmd represents the setupHook command
var setupHookCmd = &cobra.Command{
	Use:     "setup-hook",
	Aliases: []string{"setupHook"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		hookPath := filepath.Join(".git", "hooks", "post-commit")
		hookDir := filepath.Dir(hookPath)

		if _, err := os.Stat(hookDir); os.IsNotExist(err) {
			log.Fatalf("Git hooks directory not found. Please initialize a Git repository using 'git init' before setting up hooks.")
		}

		hookScript := `#!/bin/sh
commit_hash=$(git rev-parse HEAD)
commit_message=$(git log -1 --pretty=%B)
commit_author=$(git log -1 --pretty=%an)
commit_date=$(git log -1 --pretty=%ad)

worklogger record-commit \
  --hash "$commit_hash" \
  --message "$commit_message" \
  --author "$commit_author" \
  --date "$commit_date"
`

		err := os.WriteFile(hookPath, []byte(hookScript), 0755)
		if err != nil {
			log.Fatalf("Failed to write Git hook: %v", err)
		}

		fmt.Println("Git post-commit hook installed successfully!")
		fmt.Println("All future commits will be auto-logged into WorkLogger")

	},
}

func init() {
	rootCmd.AddCommand(setupHookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupHookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupHookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
