package cmd

import (
	"log"

	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		workDir, err := config.BootstrapDataDir()
		if err != nil {
			log.Fatal(err)
		}

		db, err := database.Open(workDir)
		if err != nil {
			log.Fatal(err)
		}

		// TODO(patrik): Maybe create a flag to run this on startup
		err = runMigrateUp(db)
		if err != nil {
			log.Fatal(err)
		}

		_ = db

		e := server.New(db, config.LibraryDir, workDir)

		err = e.Start(config.ListenAddr)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
