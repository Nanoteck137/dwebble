package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/nanoteck137/dwebble/library"
	"github.com/spf13/cobra"
)

var libraryCmd = &cobra.Command{
	Use: "library",
}

var libraryPrint = &cobra.Command{
	Use: "print",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World")

		libraryDir := os.Getenv("LIBRARY_DIR")
		if libraryDir == "" {
			log.Fatal("LIBRARY_DIR not set")
		}

		lib, err := library.ReadFromDir(libraryDir)
		if err != nil {
			log.Fatal(err)
		}

		for _, artist := range lib.Artists {
			fmt.Printf("%s\n", artist.Name)
			for _, album := range artist.Albums {
				fmt.Printf(" - %s\n", album.Name)
				for _, track := range album.Tracks {
					fmt.Printf("   - %02v. %s\n", track.Number, track.Name)
				}
			}
		}
	},
}

func init() {
	libraryCmd.AddCommand(libraryPrint)
	rootCmd.AddCommand(libraryCmd)
}
