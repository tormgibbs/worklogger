package cmd

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var MigrationFiles embed.FS

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing WorkLogger...")

		dir := ".worklogger"
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.Mkdir(dir, 0755)
			if err != nil {
				log.Fatalf("Failed to create .worklogger folder: %v", err)
			}
			fmt.Println("Created .worklogger folder")
		}

		dbPath := ".worklogger/db.sqlite"
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("Failed to open DB: %v", err)
		}
		defer db.Close()

		driver, err := sqlite.WithInstance(db, &sqlite.Config{})
		if err != nil {
			log.Fatalf("Failed to create SQLite driver: %v", err)
		}

		d, err := iofs.New(MigrationFiles, "migrations")
		if err != nil {
			log.Fatalf("Failed to create migration source from embedded files: %v", err)
		}

		m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)
		if err != nil {
			log.Fatalf("Failed to create migrate instance: %v", err)
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration failed: %v", err)
		}

		fmt.Println("Migrations ran successfully!")
		fmt.Println("Run `worklogger setup-hook` to enable automatic commit tracking")
		fmt.Println("OR")
		fmt.Println("Run `worklogger sync` to import existing commit history manually")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
