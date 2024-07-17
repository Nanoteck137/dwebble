package cli

import (
	"fmt"

	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/library"
	"github.com/spf13/cobra"
)

var libraryCmd = &cobra.Command{
	Use: "library",
}

var libraryPrint = &cobra.Command{
	Use: "print",
	Run: func(cmd *cobra.Command, args []string) {
		lib, err := library.ReadFromDir(config.LoadedConfig.LibraryDir)
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
		app := core.NewBaseApp(&config.LoadedConfig)

		err := app.Bootstrap()
		if err != nil {
			log.Fatal("Failed to bootstrap app", "err", err)
		}

		// TODO(patrik): Maybe create a flag to run this on startup
		err = runMigrateUp(app.DB())
		if err != nil {
			log.Fatal("Failed to run migration up", "err", err)
		}

		lib, err := library.ReadFromDir(config.LoadedConfig.LibraryDir)
		if err != nil {
			log.Fatal("Failed to read library", "err", err)
		}

		err = lib.Sync(app.WorkDir(), app.DB())
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
