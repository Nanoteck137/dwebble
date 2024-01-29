package cmd

import (
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/nanoteck137/dwebble/sql/schema"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
}

var upCmd = &cobra.Command{
	Use: "up",
	Run: func(cmd *cobra.Command, args []string) {
		godotenv.Load()

		dbUrl := os.Getenv("DB_URL")
		if dbUrl == "" {
			log.Fatal("DB_URL not set")
		}

		db, err := goose.OpenDBWithDriver("pgx", dbUrl)
		if err != nil {
			log.Fatalf("goose: failed to open DB: %v\n", err)
		}

		goose.Up(db, ".")
	},
}
var downCmd = &cobra.Command{
	Use: "down",
	Run: func(cmd *cobra.Command, args []string) {
		godotenv.Load()

		dbUrl := os.Getenv("DB_URL")
		if dbUrl == "" {
			log.Fatal("DB_URL not set")
		}

		db, err := goose.OpenDBWithDriver("pgx", dbUrl)
		if err != nil {
			log.Fatalf("goose: failed to open DB: %v\n", err)
		}

		goose.Down(db, ".")
	},
}

func init() {
	goose.SetBaseFS(schema.Migrations)

	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)
	rootCmd.AddCommand(migrateCmd)
}
