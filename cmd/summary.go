/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/data"
)


func color(change float64) string {
	if change < 0 {
		return red
	}
	return green
}

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show a summary of your work stats",
	Long: `Get an overview of your tracked time, session counts,
and productivity score based on the data in your local database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stats, err := data.GetSummaryStats(db)

		if err != nil {
			return fmt.Errorf("failed to get summary stats: %w", err)
		}

		fmt.Println("📊 Summary Stats:")
		fmt.Printf("• Today's Hours: %.2f hrs (%s%+.2f%%%s)\n", stats.TodayHours.Value, color(stats.TodayHours.Change), stats.TodayHours.Change, reset)
		fmt.Printf("• Week Hours: %.2f hrs (%s%+.2f%%%s)\n", stats.WeekHours.Value, color(stats.WeekHours.Change), stats.WeekHours.Change, reset)
		fmt.Printf("• Sessions Today: %.0f (%s%+.2f%%%s)\n", stats.SessionsToday.Value, color(stats.SessionsToday.Change), stats.SessionsToday.Change, reset)
		fmt.Printf("• Productivity Score: %.2f%% (%s%+.2f%%%s)\n", stats.ProductivityScore.Value, color(stats.ProductivityScore.Change), stats.ProductivityScore.Change, reset)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}
