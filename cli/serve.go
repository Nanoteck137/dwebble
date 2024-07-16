package cli

import (
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/log"
	"github.com/nanoteck137/dwebble/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		// workDir, err := config.LoadedConfig.BootstrapDataDir()
		// if err != nil {
		// 	log.Fatal("Failed to bootstrap data dir", "err", err)
		// }
		//
		// db, err := database.Open(workDir)
		// if err != nil {
		// 	log.Fatal("Failed to open database", "err", err)
		// }

		app := core.NewBaseApp(&config.LoadedConfig)

		err := app.Bootstrap()
		if err != nil {
			log.Fatal("Failed to bootstrap app", "err", err)
		}

		// TODO(patrik): Maybe create a flag to run this on startup
		err = runMigrateUp(app.DB())
		if err != nil {
			log.Fatal("Failed to run migrate up", "err", err)
		}

		e := server.New(app)

		err = e.Start(config.LoadedConfig.ListenAddr)
		if err != nil {
			log.Fatal("Failed to start server", "err", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
