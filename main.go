package main

import (
	"embed"

	"github.com/tormgibbs/worklogger/cmd"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func main() {
	cmd.MigrationFiles = migrationFiles
	cmd.Execute()
}
