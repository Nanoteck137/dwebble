package cmd

import (
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/log"
	"github.com/nanoteck137/dwebble/migrations"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
}

func runMigrateUp(db *database.Database) error {
	return goose.Up(db.RawConn, ".")
}

var upCmd = &cobra.Command{
	Use: "up",
	Run: func(cmd *cobra.Command, args []string) {
		workDir, err := config.Current.BootstrapDataDir()

		db, err := database.Open(workDir)
		if err != nil {
			log.Fatal("Failed to open database", "err", err)
		}

		err = runMigrateUp(db)
		if err != nil {
			log.Fatal("Failed to run migrate up", "err", err)
		}
	},
}

var downCmd = &cobra.Command{
	Use: "down",
	Run: func(cmd *cobra.Command, args []string) {
		workDir, err := config.Current.BootstrapDataDir()

		db, err := database.Open(workDir)
		if err != nil {
			log.Fatal("Failed to open database", "err", err)
		}

		err = goose.Down(db.RawConn, ".")
		if err != nil {
			log.Fatal("Failed to run migrate down", "err", err)
		}
	},
}

// TODO(patrik): Move to dev cmd
var createCmd = &cobra.Command{
	Use:  "create <MIGRATION_NAME>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		workDir, err := config.Current.BootstrapDataDir()

		db, err := database.Open(workDir)
		if err != nil {
			log.Fatal("Failed to open database", "err", err)
		}

		err = goose.Create(db.RawConn, "./migrations", name, "sql")
		if err != nil {
			log.Fatal("Failed to create migration", "err", err)
		}
	},
}

// TODO(patrik): Move to dev cmd?
var fixCmd = &cobra.Command{
	Use: "fix",
	Run: func(cmd *cobra.Command, args []string) {
		err := goose.Fix("./migrations")
		if err != nil {
			log.Fatal("Failed to fix migrations", "err", err)
		}
	},
}

func init() {
	// TODO(patrik): Move?
	goose.SetBaseFS(migrations.Migrations)
	goose.SetDialect("sqlite3")

	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)
	migrateCmd.AddCommand(createCmd)
	migrateCmd.AddCommand(fixCmd)

	rootCmd.AddCommand(migrateCmd)
}
