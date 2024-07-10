package cmd

import (
	"fmt"

	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/log"
	"github.com/spf13/cobra"
)

var libraryCmd = &cobra.Command{
	Use: "library",
}

var libraryPrint = &cobra.Command{
	Use: "print",
	Run: func(cmd *cobra.Command, args []string) {
		lib, err := library.ReadFromDir(config.Current.LibraryDir)
		if err != nil {
			log.Fatal("Failed to read library", "err", err)
		}

		for _, artist := range lib.Artists {
			fmt.Printf("%s\n", artist.Name)
			for _, album := range artist.Albums {
				fmt.Printf(" - %s\n", album.Name)
				for _, track := range album.Tracks {
					fmt.Printf("   - %02v. %s\n", track.Number, track.Name)
					fmt.Printf("     Artist: %s\n", track.Artist.Name)
					// fmt.Printf("     Tags: %s\n", strings.Join(track.Tags, ", "))
					// fmt.Printf("     Duration: %v\n", utils.FormatTime(track.Duration))
				}
			}
		}
	},
}

var librarySync = &cobra.Command{
	Use: "sync",
	Run: func(cmd *cobra.Command, args []string) {
		workDir, err := config.Current.BootstrapDataDir()
		if err != nil {
			log.Fatal("Failed to bootstrap data dir", "err", err)
		}

		db, err := database.Open(workDir)
		if err != nil {
			log.Fatal("Failed to open database", "err", err)
		}

		// TODO(patrik): Maybe create a flag to run this on startup
		err = runMigrateUp(db)
		if err != nil {
			log.Fatal("Failed to run migration up", "err", err)
		}

		lib, err := library.ReadFromDir(config.Current.LibraryDir)
		if err != nil {
			log.Fatal("Failed to read library", "err", err)
		}

		err = lib.Sync(workDir, db)
		if err != nil {
			log.Fatal("Failed to sync library", "err", err)
		}
	},
}

func init() {
	libraryCmd.AddCommand(libraryPrint)
	libraryCmd.AddCommand(librarySync)

	rootCmd.AddCommand(libraryCmd)
}
