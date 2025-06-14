/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/data"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export your data (CSV, JSON, etc)",
	Long:  `Export your tracked work sessions, stats, and summaries to CSV or other formats`,
	RunE: func(cmd *cobra.Command, args []string) error {
		outFile, _ := cmd.Flags().GetString("out")
		useCSV, _ := cmd.Flags().GetBool("csv")

		if !useCSV {
			fmt.Println("No export type specified. Defaulting to CSV...")
			useCSV = true
		}

		if useCSV {
			return data.ExportToCSV(db, outFile)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().Bool("csv", false, "Export data as CSV")
	exportCmd.Flags().StringP("out", "o", "export.csv", "Output filename for exported CSV")
}
