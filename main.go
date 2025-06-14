package main

import (
	"embed"

	"github.com/tormgibbs/worklogger/cmd"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func main() {
	cmd.SetMigrationFiles(migrationFiles)
	cmd.Execute()
}
