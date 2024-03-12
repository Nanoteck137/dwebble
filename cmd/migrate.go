package cmd

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nanoteck137/dwebble/migrations"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
}

var migrateUpCmd = &cobra.Command{
	Use: "up",
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl := getDatabaseUrl()

		db, err := goose.OpenDBWithDriver("pgx", dbUrl)
		if err != nil {
			log.Fatalf("migrate: failed to open database: %v", err)
		}

		err = goose.Up(db, ".")
		if err != nil {
			log.Fatalf("migrate: failed to upgrade migrations: %v", err)
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use: "down",
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl := getDatabaseUrl()

		db, err := goose.OpenDBWithDriver("pgx", dbUrl)
		if err != nil {
			log.Fatalf("migrate: failed to open database: %v", err)
		}

		err = goose.Down(db, ".")
		if err != nil {
			log.Fatalf("migrate: failed to downgrade migrations: %v", err)
		}
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:  "create <MIGRATION_NAME>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl := getDatabaseUrl()

		db, err := goose.OpenDBWithDriver("pgx", dbUrl)
		if err != nil {
			log.Fatalf("migrate: failed to open database: %v", err)
		}

		name := args[0]
		err = goose.Create(db, "./migrations", name, "sql")
		if err != nil {
			log.Fatalf("migrate: failed to create new migration: %v", err)
		}
	},
}

var migrateFixCmd = &cobra.Command{
	Use: "fix",
	Run: func(cmd *cobra.Command, args []string) {
		err := goose.Fix("./migrations")
		if err != nil {
			log.Fatalf("migrate: failed to fix migrations: %v", err)
		}
	},
}

func init() {
	goose.SetBaseFS(migrations.Migrations)

	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd, migrateCreateCmd, migrateFixCmd)

	rootCmd.AddCommand(migrateCmd)
}
