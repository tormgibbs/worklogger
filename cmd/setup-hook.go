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
	Short:   "Install Git post-commit hook for auto-logging",
	Long: `Sets up a Git post-commit hook that automatically logs commit details 
(hash, message, author, and date) into WorkLogger after each commit.`,
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
}
