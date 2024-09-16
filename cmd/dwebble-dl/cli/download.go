package cli

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nanoteck137/dwebble/cmd/dwebble-dl/api"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/spf13/cobra"
)

func DownloadFile(url, dst string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func DownloadTrack(track *api.Track, dst string) (string, error) {
	split := strings.Split(track.MobileMediaUrl, ".")
	ext := split[len(split)-1]

	name := utils.Slug(track.Name) + "-" + track.Id

	p := path.Join(dst, name+"."+ext)

	err := DownloadFile(track.MobileMediaUrl, p)
	if err != nil {
		return "", err
	}

	return p, nil
}

func DownloadTrackCover(track *api.Track, dst string) (string, error) {
	split := strings.Split(track.CoverArt.Medium, ".")
	ext := split[len(split)-1]

	p := path.Join(dst, "cover."+ext)

	err := DownloadFile(track.CoverArt.Medium, p)
	if err != nil {
		return "", err
	}

	return p, nil
}

var downloadCmd = &cobra.Command{
	Use: "download",
}

var downloadFilterCmd = &cobra.Command{
	Use: "filter",
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		filter, _ := cmd.Flags().GetString("filter")
		sort, _ := cmd.Flags().GetString("sort")
		albumName, _ := cmd.Flags().GetString("album")
		output, _ := cmd.Flags().GetString("output")

		client := api.New(server)

		res, err := client.GetTracks(api.Options{
			QueryParams: map[string]string{
				"filter": filter,
				"sort":   sort,
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = os.Stat(output)
		if err == nil {
			val := false
			s := huh.NewConfirm().
				Title(fmt.Sprintf("'%s' is not empty\nThis operation will delete the folder", output)).
				Value(&val)
			err = s.Run()
			if err != nil {
				log.Fatal(err)
			}

			err = os.RemoveAll(output)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = os.MkdirAll(output, 0755)
		if err != nil {
			log.Fatal(err)
		}

		tmpDir := path.Join(output, "tmp")

		err = os.MkdirAll(tmpDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		for i, track := range res.Tracks {
			fmt.Printf("Downloading '%s' ", track.Name)
			p, err := DownloadTrack(&track, output)
			if err != nil {
				log.Fatal(err)
			}

			c, err := DownloadTrackCover(&track, tmpDir)
			if err != nil {
				log.Fatal(err)
			}

			var args []string
			args = append(args, "--title", track.Name)
			args = append(args, "--artist", track.ArtistName)
			args = append(args, "--album", albumName)
			args = append(args, "--number", strconv.FormatInt(int64(i+1), 10))
			args = append(args, "--image", c)
			args = append(args, "--remove")

			args = append(args, p)

			cmd := exec.Command("tagopus", args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Done\n")
		}

		err = os.RemoveAll(tmpDir)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	downloadFilterCmd.Flags().String("server", "", "Server Address")
	downloadFilterCmd.Flags().StringP("filter", "f", "", "Filter to use")
	downloadFilterCmd.Flags().StringP("sort", "s", "", "Sorting to use")
	downloadFilterCmd.Flags().StringP("album", "a", "Unknown", "Album Name")
	downloadFilterCmd.Flags().StringP("output", "o", ".", "Output Directory")

	downloadFilterCmd.MarkFlagRequired("server")

	downloadCmd.AddCommand(downloadFilterCmd)

	rootCmd.AddCommand(downloadCmd)
}
