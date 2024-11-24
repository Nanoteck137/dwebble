package cli

import (
	"fmt"
	"log"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/cmd/dwebble-cli/api"
	"github.com/spf13/cobra"
)

var albumCmd = &cobra.Command{
	Use: "album",
}

var albumCreateCmd = &cobra.Command{
	Use: "create <ALBUM_NAME> <ARTIST_NAME>",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		client := api.New(server)

		albumName := args[0]
		artistName := args[1]

		res, err := client.CreateAlbum(api.CreateAlbumBody{
			Name:   albumName,
			Artist: artistName,
		}, api.Options{})
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(res)

		fmt.Print(res.AlbumId)
	},
}

func init() {
	albumCmd.AddCommand(albumCreateCmd)

	rootCmd.AddCommand(albumCmd)
}
