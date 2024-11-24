package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/cmd/dwebble-cli/api"
	"github.com/spf13/cobra"
)

var albumCmd = &cobra.Command{
	Use: "album",
}

var albumCreateCmd = &cobra.Command{
	Use:  "create <ALBUM_NAME> <ARTIST_NAME>",
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

var albumUploadCoverCmd = &cobra.Command{
	Use:  "upload-cover <ALBUM_ID> <COVER_PATH>",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		client := api.New(server)

		albumId := args[0]
		coverPath := args[1]

		body := &bytes.Buffer{}
		w := multipart.NewWriter(body)

		fw, err := w.CreateFormFile("cover", path.Base(coverPath))
		if err != nil {
			log.Fatal(err)
		}

		src, err := os.Open(coverPath)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(fw, src)
		if err != nil {
			log.Fatal(err)
		}

		err = w.Close()
		if err != nil {
			log.Fatal(err)
		}

		res, err := client.ChangeAlbumCover(albumId, body, api.Options{
			Boundary: w.Boundary(),
		})
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(res)
	},
}

var albumUploadTracksCmd = &cobra.Command{
	Use:  "upload-tracks <ALBUM_ID> <TRACKS...>",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		client := api.New(server)

		albumId := args[0]
		tracks := args[1:]

		_ = client
		_ = albumId

		pretty.Println(tracks)

		body := &bytes.Buffer{}
		w := multipart.NewWriter(body)

		for _, track := range tracks {
			func() {
				fmt.Printf("path.Base(track): %v\n", path.Base(track))
				fw, err := w.CreateFormFile("files", path.Base(track))
				if err != nil {
					log.Fatal(err)
				}

				src, err := os.Open(track)
				if err != nil {
					log.Fatal(err)
				}
				defer src.Close()

				_, err = io.Copy(fw, src)
				if err != nil {
					log.Fatal(err)
				}
			}()
		}

		err := w.Close()
		if err != nil {
			log.Fatal(err)
		}

		res, err := client.UploadTracks(albumId, body, api.Options{
			Boundary: w.Boundary(),
		})
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(res)
	},
}

func init() {
	albumCmd.AddCommand(albumCreateCmd)
	albumCmd.AddCommand(albumUploadCoverCmd)
	albumCmd.AddCommand(albumUploadTracksCmd)

	rootCmd.AddCommand(albumCmd)
}
