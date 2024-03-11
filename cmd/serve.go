package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/server"
	"github.com/nanoteck137/dwebble/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getDatabaseUrl() string {
	if !viper.IsSet("database_url") {
		log.Fatal("'database_url' not set")
	}

	return viper.GetString("database_url")
}

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl := getDatabaseUrl()
		conn, err := pgxpool.New(context.Background(), dbUrl)
		if err != nil {
			log.Fatal(err)
		}

		dataDir := viper.GetString("data_dir")
		workDir := types.WorkDir(dataDir)

		if !viper.IsSet("library_dir") {
			log.Fatal("'library_dir' not set")
		}

		libraryDir := viper.GetString("library_dir")

		db := database.New(conn)
		e := server.New(db, libraryDir, workDir)

		listenAddr := viper.GetString("listen_addr")
		fmt.Printf("listenAddr: %v\n", listenAddr)
		err = e.Start(listenAddr)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// serveCmd.Flags().IntP("port", "p", 3000, "Server Port")
	// viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))

	rootCmd.AddCommand(serveCmd)
}
