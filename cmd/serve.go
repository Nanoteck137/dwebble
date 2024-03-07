package cmd

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/server"
	"github.com/nanoteck137/dwebble/types"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}

		dbUrl := os.Getenv("DB_URL")
		if dbUrl == "" {
			log.Fatal("DB_URL not set")
		}

		conn, err := pgxpool.New(context.Background(), dbUrl)
		if err != nil {
			log.Fatal(err)
		}

		workDirPath := os.Getenv("WORK_DIR")
		if workDirPath == "" {
			log.Println("Warning: WORK_DIR not set, using cwd")
			workDirPath = "./"
		}
		workDir := types.WorkDir(workDirPath)

		libraryDir := os.Getenv("LIBRARY_DIR")
		if libraryDir == "" {
			log.Fatal("LIBRARY_DIR not set")
		}

		db := database.New(conn)
		e := server.New(db, libraryDir, workDir)

		err = e.Start("0.0.0.0:3000")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
