package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nanoteck137/dwebble/cmd/dwebble-cli/api"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/spf13/cobra"
)

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	// TODO(patrik): Return the error
	if err != nil {
		log.Fatal("Failed to open url", "err", err)
	}

}

func TransformStringArray(arr []string) []string {
	for i, v := range arr {
		arr[i] = transform.String(v)
	}

	return arr
}

func WriteFile(name string, data []byte) error {
	// f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	// TODO(patrik): used for testing
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err == nil {
		return err
	}

	return nil
}

type Context struct {
	client  *api.Client
	artists map[string]string
}

func (c *Context) SetAlbumCover(albumId, filename string) error {
	// TODO(patrik): Check cover type

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	dw, err := w.CreateFormFile("cover", path.Base(filename))
	if err != nil {
		return err
	}

	src, err := os.Open(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(dw, src)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	_, err = c.client.ChangeAlbumCover(albumId, &b, api.Options{
		Boundary: w.Boundary(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) UploadTrackSimple(trackFilename string, body api.UploadTrackBody) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	{
		dw, err := w.CreateFormField("body")
		if err != nil {
			return err
		}

		encoder := json.NewEncoder(dw)
		err = encoder.Encode(&body)
		if err != nil {
			return err
		}
	}

	{

		dw, err := w.CreateFormFile("track", path.Base(trackFilename))
		if err != nil {
			return err
		}

		src, err := os.Open(trackFilename)
		if err != nil {
			return err
		}

		_, err = io.Copy(dw, src)
		if err != nil {
			return err
		}
	}

	err := w.Close()
	if err != nil {
		return err
	}

	_, err = c.client.UploadTrack(&b, api.Options{
		Boundary: w.Boundary(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) SearchAlbum(search string) ([]api.Album, error) {
	res, err := c.client.SearchAlbums(api.Options{})
	if err != nil {
		return nil, err
	}

	return res.Albums, nil
}

func (c *Context) FindAlbums(name, artist string) ([]api.Album, error) {
	search, err := c.client.GetAlbums(api.Options{
		QueryParams: map[string]string{
			"filter": fmt.Sprintf("name == \"%s\" && artistName == \"%s\"", name, artist),
		},
	})
	if err != nil {
		return nil, err
	}

	return search.Albums, nil
}

func (c *Context) GetOrCreateMultipleArtists(names []string) (string, []string, error) {
	if len(names) <= 0 {
		// TODO(patrik): Better error?
		return "", nil, errors.New("No artists")
	}

	primaryId, err := c.GetOrCreateArtist(names[0])
	if err != nil {
		return "", nil, err
	}

	rest := names[1:]

	featuringArtists := make([]string, len(rest))
	for i, name := range rest {
		featuringArtists[i], err = c.GetOrCreateArtist(name)
		if err != nil {
			return "", nil, err
		}
	}

	return primaryId, featuringArtists, nil
}

func (c *Context) GetOrCreateArtist(name string) (string, error) {
	if id, exists := c.artists[name]; exists {
		return id, nil
	}

	search, err := c.client.GetArtists(api.Options{
		QueryParams: map[string]string{
			"filter": fmt.Sprintf("name == \"%s\"", name),
		},
	})
	if err != nil {
		return "", err
	}

	if len(search.Artists) <= 0 {
		// TODO(patrik): Create artist

		res, err := c.client.CreateArtist(api.CreateArtistBody{
			Name: name,
		}, api.Options{})
		if err != nil {
			return "", err
		}

		c.artists[name] = res.Id
		return res.Id, nil
	}

	// TODO(patrik): Should we always pick the first one?
	artist := search.Artists[0]
	c.artists[name] = artist.Id
	return artist.Id, nil
}

type TrackInfo struct {
	Name   string
	Artist string
	Number int
	Year   int
}

var dateRegex = regexp.MustCompile(`^([12]\d\d\d)`)

func parseArtist(s string) []string {
	if s == "" {
		return []string{}
	}

	return strings.Split(s, ",")
}

func getTrackInfo(p string) (TrackInfo, error) {

	probe, err := utils.ProbeTrack(p)
	if err != nil {
		return TrackInfo{}, err
	}

	var res TrackInfo

	if tag, exists := probe.Tags["title"]; exists {
		res.Name = tag
	}

	if tag, exists := probe.Tags["artist"]; exists {
		res.Artist = tag
	}

	if tag, exists := probe.Tags["date"]; exists {
		match := dateRegex.FindStringSubmatch(tag)
		if len(match) > 0 {
			res.Year, _ = strconv.Atoi(match[1])
		}
	}

	if tag, exists := probe.Tags["track"]; exists {
		res.Number, _ = strconv.Atoi(tag)
	}

	return res, nil
}

var importCmd = &cobra.Command{
	Use: "import",
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		web, _ := cmd.Flags().GetString("web")
		dir, _ := cmd.Flags().GetString("dir")

		client := api.New(server)

		c := Context{
			client:  client,
			artists: map[string]string{},
		}

		var images []string
		var tracks []string

		entries, err := os.ReadDir(dir)
		if err != nil {
			log.Fatal("Failed to read dir", "err", err)
		}

		for _, e := range entries {
			if e.IsDir() {
				continue
			}

			p := path.Join(dir, e.Name())

			ext := path.Ext(p)

			if utils.IsValidTrackExt(ext) {
				tracks = append(tracks, p)
			}

			if utils.IsValidImageExt(ext) {
				images = append(images, p)
			}
		}

		if len(tracks) <= 0 {
			log.Warn("No tracks found... Quitting")
			return
		}

		p := tracks[0]

		name := "Unknown Album name"
		var artists []string
		year := 0

		probe, err := utils.ProbeTrack(p)
		if err != nil {
			log.Fatal("Failed to probe track", "err", err)
		}

		isSingle := len(tracks) == 1

		if !isSingle {
			if tag, exists := probe.Tags["album"]; exists {
				name = tag
			}
		} else {
			// NOTE(patrik): If we only have one track then we make the
			// album name the same as the track name
			if tag, exists := probe.Tags["title"]; exists {
				name = tag
			}
		}

		if !isSingle {
			if tag, exists := probe.Tags["album_artist"]; exists {
				artists = parseArtist(tag)
			} else {
				if tag, exists := probe.Tags["artist"]; exists {
					artists = parseArtist(tag)
				}
			}
		} else {
			if tag, exists := probe.Tags["artist"]; exists {
				artists = parseArtist(tag)
			}
		}

		if tag, exists := probe.Tags["date"]; exists {
			match := dateRegex.FindStringSubmatch(tag)
			if len(match) > 0 {
				year, _ = strconv.Atoi(match[1])
			}
		}

		// NOTE(patrik): Try to extract the real year if yt-dlp has been used
		if tag, exists := probe.Tags["synopsis"]; exists {
			fmt.Printf("tag: %v\n", tag)

			lines := strings.Split(tag, "\n")

			for _, line := range lines {
				if strings.Contains(line, "Released on:") {
					line = strings.TrimPrefix(line, "Released on:")
					line = strings.TrimSpace(line)

					y := utils.ExtractNumber(line)
					if y != 0 {
						year = y
					}
				}
			}
		}

		// TODO(patrik): Add checks
		artistId, featuringArtists, err := c.GetOrCreateMultipleArtists(artists)
		if err != nil {
			log.Fatal("Failed to get/create album artist", "err", err)
		}

		album, err := client.CreateAlbum(api.CreateAlbumBody{
			Name:             name,
			OtherName:        "",
			ArtistId:         artistId,
			Year:             year,
			Tags:             []string{},
			FeaturingArtists: featuringArtists,
		}, api.Options{})
		if err != nil {
			log.Fatal("Failed", "err", err)
		}

		fmt.Printf("Created album: %s (%s)\n", name, album.AlbumId)

		if len(images) > 0 {
			fmt.Printf("Setting album cover art...")
			err = c.SetAlbumCover(album.AlbumId, images[0])
			if err != nil {
				fmt.Printf("failed\n")
				log.Fatal("Failed to set album cover", "err", err)
			}
			fmt.Printf("done\n")
		}

		for i, p := range tracks {
			filename := path.Base(p)
			fmt.Printf("Uploading track (%02d/%02d): %s ...", i+1, len(tracks), filename)

			trackInfo, err := getTrackInfo(p)
			if err != nil {
				fmt.Printf("failed\n")
				log.Fatal("Failed to get track info", "err", err)
			}

			if trackInfo.Name == "" {
				trackInfo.Name = strings.TrimSuffix(filename, path.Ext(p))
			}

			if trackInfo.Number == 0 {
				trackInfo.Number = utils.ExtractNumber(filename)
			}

			// TODO(patrik): If artist is empty then use album maybe
			artists := parseArtist(trackInfo.Artist)

			artistId, featuringArtists, err := c.GetOrCreateMultipleArtists(artists)
			if err != nil {
				fmt.Printf("failed\n")
				log.Fatal("Failed to get/create track artist", "err", err)
			}

			err = c.UploadTrackSimple(p, api.UploadTrackBody{
				Name:             trackInfo.Name,
				OtherName:        "",
				Number:           trackInfo.Number,
				Year:             trackInfo.Year,
				AlbumId:          album.AlbumId,
				ArtistId:         artistId,
				Tags:             []string{},
				FeaturingArtists: featuringArtists,
			})
			if err != nil {
				fmt.Printf("failed\n")
				log.Fatal("Failed to upload track", "err", err)
			}
			fmt.Printf("done\n")
		}

		confirmed := false
		err = huh.NewConfirm().
			Title("Want to open the album in the browser?").
			Value(&confirmed).
			Run()
		if err != nil {
			log.Fatal("Failed", "err", err)
		}

		if confirmed {
			openbrowser(fmt.Sprintf("%s/albums/%s", web, album.AlbumId))
		}
	},
}

func init() {
	importCmd.Flags().StringP("dir", "d", ".", "Directory to search")

	rootCmd.AddCommand(importCmd)
}
