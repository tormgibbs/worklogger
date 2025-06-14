package cmd

import (
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/web/server"
)

// serveCmd represents the serve command
var studioCmd = &cobra.Command{
	Use:   "studio",
	Short: "Start the worklogger web studio",
	Long: `The "studio" command launches the Worklogger web interface 
for logging and viewing your work data via the browser.`,
	Run: func(cmd *cobra.Command, args []string) {
	go func() {
		err := browser.OpenURL(server.Addr)
		if err != nil {
			cmd.PrintErr("Failed to open browser:", err)
			return
		}
	}()

	server.Serve(db)
	},
}

func init() {
	rootCmd.AddCommand(studioCmd)
}
