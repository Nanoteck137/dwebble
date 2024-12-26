package cmd

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"mime/multipart"
// 	"os"
// 	"path"
//
// 	"github.com/kr/pretty"
// 	"github.com/nanoteck137/dwebble/cmd/dwebble-cli/api"
// 	"github.com/spf13/cobra"
// )
//
// type ImportFile struct {
// 	Album    string   `json:"album"`
// 	Artist   string   `json:"artist"`
// 	CoverArt string   `json:"coverArt"`
// 	Tracks   []string `json:"tracks"`
// }
//
// var albumCmd = &cobra.Command{
// 	Use: "album",
// }
//
// var albumCreateCmd = &cobra.Command{
// 	Use:  "create <ALBUM_NAME> <ARTIST_NAME>",
// 	Args: cobra.ExactArgs(2),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		server, _ := cmd.Flags().GetString("server")
// 		client := api.New(server)
//
// 		albumName := args[0]
// 		artistName := args[1]
//
// 		res, err := client.CreateAlbum(api.CreateAlbumBody{
// 			Name:   albumName,
// 			Artist: artistName,
// 		}, api.Options{})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		pretty.Println(res)
//
// 		fmt.Print(res.AlbumId)
// 	},
// }
//
// func createAlbumCoverFormData(w *multipart.Writer, file string) error {
// 	fw, err := w.CreateFormFile("cover", path.Base(file))
// 	if err != nil {
// 		return err
// 	}
//
// 	src, err := os.Open(file)
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = io.Copy(fw, src)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// var albumUploadCoverCmd = &cobra.Command{
// 	Use:  "upload-cover <ALBUM_ID> <COVER_PATH>",
// 	Args: cobra.ExactArgs(2),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		server, _ := cmd.Flags().GetString("server")
// 		client := api.New(server)
//
// 		albumId := args[0]
// 		coverPath := args[1]
//
// 		body := &bytes.Buffer{}
// 		w := multipart.NewWriter(body)
//
// 		err := createAlbumCoverFormData(w, coverPath)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		err = w.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		res, err := client.ChangeAlbumCover(albumId, body, api.Options{
// 			Boundary: w.Boundary(),
// 		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		pretty.Println(res)
// 	},
// }
//
// var albumUploadTracksCmd = &cobra.Command{
// 	Use:  "upload-tracks <ALBUM_ID> <TRACKS...>",
// 	Args: cobra.MinimumNArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		server, _ := cmd.Flags().GetString("server")
// 		client := api.New(server)
//
// 		albumId := args[0]
// 		tracks := args[1:]
//
// 		body := &bytes.Buffer{}
// 		w := multipart.NewWriter(body)
//
// 		bw, err := w.CreateFormField("body")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		e := json.NewEncoder(bw)
// 		err = e.Encode(api.UploadTracksBody{
// 			ForceExtractNumber: true,
// 		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		for _, track := range tracks {
// 			func() {
// 				fw, err := w.CreateFormFile("files", path.Base(track))
// 				if err != nil {
// 					log.Fatal(err)
// 				}
//
// 				src, err := os.Open(track)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				defer src.Close()
//
// 				_, err = io.Copy(fw, src)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 			}()
// 		}
//
// 		err = w.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		res, err := client.UploadTracks(albumId, body, api.Options{
// 			Boundary: w.Boundary(),
// 		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		pretty.Println(res)
// 	},
// }
//
// var albumImportCmd = &cobra.Command{
// 	Use:  "import <IMPORT_FILE>",
// 	Args: cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		server, _ := cmd.Flags().GetString("server")
// 		client := api.New(server)
//
// 		importFile := args[0]
//
// 		d, err := os.ReadFile(importFile)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		var i ImportFile
// 		err = json.Unmarshal(d, &i)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		album, err := client.CreateAlbum(api.CreateAlbumBody{
// 			Name:   i.Album,
// 			Artist: i.Artist,
// 		}, api.Options{})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		body := &bytes.Buffer{}
// 		w := multipart.NewWriter(body)
//
// 		coverPath := path.Join(path.Dir(importFile), i.CoverArt)
// 		err = createAlbumCoverFormData(w, coverPath)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		err = w.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		// TODO(patrik): Pyrin error needs fixing
// 		_, err = client.ChangeAlbumCover(album.AlbumId, body, api.Options{
// 			Boundary: w.Boundary(),
// 		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		{
// 			body := &bytes.Buffer{}
// 			w := multipart.NewWriter(body)
//
// 			bw, err := w.CreateFormField("body")
// 			if err != nil {
// 				log.Fatal(err)
// 			}
//
// 			e := json.NewEncoder(bw)
// 			err = e.Encode(api.UploadTracksBody{
// 				ForceExtractNumber: true,
// 			})
// 			if err != nil {
// 				log.Fatal(err)
// 			}
//
// 			b := path.Dir(importFile)
// 			for _, track := range i.Tracks {
// 				p := path.Join(b, track)
// 				func() {
// 					fw, err := w.CreateFormFile("files", path.Base(p))
// 					if err != nil {
// 						log.Fatal(err)
// 					}
//
// 					src, err := os.Open(p)
// 					if err != nil {
// 						log.Fatal(err)
// 					}
// 					defer src.Close()
//
// 					_, err = io.Copy(fw, src)
// 					if err != nil {
// 						log.Fatal(err)
// 					}
// 				}()
// 			}
//
// 			err = w.Close()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
//
// 			_, err = client.UploadTracks(album.AlbumId, body, api.Options{
// 				Boundary: w.Boundary(),
// 			})
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
// 	},
// }
//
// func init() {
// 	albumCmd.AddCommand(albumCreateCmd)
// 	albumCmd.AddCommand(albumUploadCoverCmd)
// 	albumCmd.AddCommand(albumUploadTracksCmd)
// 	albumCmd.AddCommand(albumImportCmd)
//
// 	rootCmd.AddCommand(albumCmd)
// }
