package cmd

import (
	"bufio"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/tormgibbs/worklogger/auth"
)

var migrationFS embed.FS

func SetMigrationFiles(fs embed.FS) {
	migrationFS = fs
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the local database and setup WorkLogger environment",
	Long:  `Initialize the worklogger environment, setting up database and applying migrations. Run this first before using other commands.`,
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

		d, err := iofs.New(migrationFS, "migrations")
		if err != nil {
			log.Fatalf("Failed to create migration source from embedded files: %v", err)
		}
		defer d.Close()

		m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)
		if err != nil {
			log.Fatalf("Failed to create migrate instance: %v", err)
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration failed: %v", err)
		}

		setupGitHubOAuth()

		fmt.Println("Migrations ran successfully!")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  • Run `worklogger setup-hook` to enable automatic commit tracking")
		fmt.Println("  • OR run `worklogger sync` to import existing commit history manually")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func setupGitHubOAuth() {
	fmt.Println("=== GitHub OAuth Setup ===")
	fmt.Println()

	// Check if already configured in keyring
	clientID, err1 := auth.GetToken("github_client_id")
	clientSecret, err2 := auth.GetToken("github_client_secret")

	if err1 == nil && err2 == nil && clientID != "" && clientSecret != "" {
		fmt.Println("✓ GitHub OAuth credentials already configured in keyring")
		return
	}

	fmt.Println("To use GitHub integration features, you need to create a GitHub OAuth App:")
	fmt.Println()
	fmt.Println("1. Go to: https://github.com/settings/applications/new")
	fmt.Println("2. Fill in:")
	fmt.Println("   - Application name: WorkLogger (or any name you prefer)")
	fmt.Println("   - Homepage URL: https://github.com/tormgibbs/worklogger")
	fmt.Println("   - Authorization callback URL: http://localhost:8080/callback")
	fmt.Println("3. Click 'Register application'")
	fmt.Println("4. Copy the Client ID and generate a Client Secret")
	fmt.Println()

	// Interactive setup
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your GitHub Client ID (or press Enter to skip): ")
	newClientID, _ := reader.ReadString('\n')
	newClientID = strings.TrimSpace(newClientID)

	if newClientID == "" {
		fmt.Println("Skipping GitHub OAuth setup. You can configure it later with:")
		fmt.Println("  worklogger setup-github")
		return
	}

	fmt.Print("Enter your GitHub Client Secret: ")
	newClientSecret, _ := reader.ReadString('\n')
	newClientSecret = strings.TrimSpace(newClientSecret)

	// Save to keyring
	if err := auth.SetToken("github_client_id", newClientID); err != nil {
		fmt.Printf("Error saving Client ID to keyring: %v\n", err)
		return
	}

	if err := auth.SetToken("github_client_secret", newClientSecret); err != nil {
		fmt.Printf("Error saving Client Secret to keyring: %v\n", err)
		// Clean up the client ID if secret failed
		auth.DeleteToken("github_client_id")
		return
	}

	// Also save redirect URI for consistency
	auth.SetToken("github_redirect_uri", "http://localhost:3000/callback")

	fmt.Println("✓ GitHub OAuth credentials saved securely to system keyring")
	fmt.Println("✓ GitHub integration is ready!")
}
