package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

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
		fmt.Printf("dataDir: %v\n", dataDir)

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
