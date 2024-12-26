package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/cmd/dwebble-cli/api"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

type ImportDataAlbum struct {
	Name          string   `toml:"name"`
	OtherName     string   `toml:"other_name"`
	Artists       []string `toml:"artists"`
	Year          int      `toml:"year"`
	Tags          []string `toml:"tags"`
	CoverArtFilename string   `toml:"coverArtFilename"`
}

type ImportDataTrack struct {
	Filename  string   `toml:"filename"`
	Name      string   `toml:"name"`
	OtherName string   `toml:"other_name"`
	Artists   []string `toml:"artists"`
	Number    int      `toml:"number"`
	Year      int      `toml:"year"`
	Tags      []string `toml:"tags"`
}

type ImportData struct {
	UseAlbumTagsForTracks bool              `toml:"use_album_tags_for_tracks"`
	Album                 ImportDataAlbum   `toml:"album"`
	Tracks                []ImportDataTrack `toml:"tracks"`
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

func getImportData(dir string) (ImportData, error) {
	p := path.Join(dir, "import.toml")
	data, err := os.ReadFile(p)
	if err != nil {
		return ImportData{}, err
	}

	var res ImportData

	err = toml.Unmarshal(data, &res)
	if err != nil {
		return ImportData{}, err
	}

	return res, nil
}

type Context struct {
	client  *api.Client
	artists map[string]string
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

var importCmd = &cobra.Command{
	Use: "import",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		data, err := getImportData(dir)
		if err != nil {
			log.Fatal("Failed to get import data", "err", err)
		}

		client := api.New("http://localhost:3000")

		c := Context{
			client:  client,
			artists: map[string]string{},
		}

		artistId, featuringArtists, err := c.GetOrCreateMultipleArtists(data.Album.Artists)
		if err != nil {
			log.Fatal("Failed to get/create album artist", "err", err)
		}

		album, err := client.CreateAlbum(api.CreateAlbumBody{
			Name:             data.Album.Name,
			OtherName:        data.Album.OtherName,
			ArtistId:         artistId,
			Year:             data.Album.Year,
			Tags:             data.Album.Tags,
			FeaturingArtists: featuringArtists,
		}, api.Options{})
		if err != nil {
			log.Fatal("Failed", "err", err)
		}

		setAlbumCoverArt := func() error {
			var b bytes.Buffer
			w := multipart.NewWriter(&b)

			{
				dw, err := w.CreateFormFile("cover", data.Album.CoverArtFilename)
				if err != nil {
					return err
				}

				src, err := os.Open(path.Join(dir, data.Album.CoverArtFilename))
				if err != nil {
					return err
				}

				_, err = io.Copy(dw, src)
				if err != nil {
					return err
				}
			}

			err = w.Close()
			if err != nil {
				return err
			}

			_, err = client.ChangeAlbumCover(album.AlbumId, &b, api.Options{
				Boundary: w.Boundary(),
			})
			if err != nil {
				var apiErr *api.ApiError[any]
				if errors.As(err, &apiErr) {
					pretty.Println(apiErr)
				}

				return err
			}

			return nil

		}

		fmt.Printf("Created album: (%s)\n", album.AlbumId)

		fmt.Printf("Setting album cover art...")
		err = setAlbumCoverArt()
		if err != nil {
			log.Fatal("Failed to set album cover art", "err", err)
		}
		fmt.Printf("done\n")

		uploadTrack := func(track ImportDataTrack) error {
			var b bytes.Buffer
			w := multipart.NewWriter(&b)

			{
				dw, err := w.CreateFormField("body")
				if err != nil {
					return err
				}

				var tags []string
				tags = append(tags, track.Tags...)

				if data.UseAlbumTagsForTracks {
					tags = append(tags, data.Album.Tags...)
				}

				artistId, featuringArtists, err := c.GetOrCreateMultipleArtists(track.Artists)
				if err != nil {
					return err
				}

				body := api.UploadTrackBody{
					Name:             track.Name,
					OtherName:        track.OtherName,
					Number:           track.Number,
					Year:             track.Year,
					AlbumId:          album.AlbumId,
					ArtistId:         artistId,
					Tags:             tags,
					FeaturingArtists: featuringArtists,
				}

				encoder := json.NewEncoder(dw)
				err = encoder.Encode(&body)
				if err != nil {
					return err
				}
			}

			{

				dw, err := w.CreateFormFile("track", track.Filename)
				if err != nil {
					return err
				}

				src, err := os.Open(path.Join(dir, track.Filename))
				if err != nil {
					return err
				}

				_, err = io.Copy(dw, src)
				if err != nil {
					return err
				}
			}

			err = w.Close()
			if err != nil {
				return err
			}

			_, err = client.UploadTrack(&b, api.Options{
				Boundary: w.Boundary(),
			})
			if err != nil {
				var apiErr *api.ApiError[any]
				if errors.As(err, &apiErr) {
					pretty.Println(apiErr)
				}

				return err
			}

			return nil
		}

		for i, track := range data.Tracks {
			fmt.Printf("Uploading track (%02d/%d): %s...", i + 1, len(data.Tracks), track.Name)
			err = uploadTrack(track)
			if err != nil {
				log.Fatal("Failed to upload track", "err", err)
			}
			fmt.Printf("done\n")
		}

	},
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

var importInitCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		var importData ImportData

		var images []string
		var tracks []string

		// TODO(patrik): Search only the top-level directory
		filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			ext := path.Ext(p)

			if utils.IsValidTrackExt(ext) {
				tracks = append(tracks, p)
			}

			if utils.IsValidImageExt(ext) {
				images = append(images, p)
			}

			return nil
		})

		for _, p := range tracks {
			filename := path.Base(p)

			trackInfo, err := getTrackInfo(p)
			if err != nil {
				log.Fatal("Failed to get track info", "err", err)
			}

			if trackInfo.Name == "" {
				trackInfo.Name = strings.TrimSuffix(filename, path.Ext(p))
			}

			if trackInfo.Number == 0 {
				trackInfo.Number = utils.ExtractNumber(filename)
			}

			importData.Tracks = append(importData.Tracks, ImportDataTrack{
				Filename:  filename,
				Name:      trackInfo.Name,
				OtherName: "",
				Artists:   parseArtist(trackInfo.Artist),
				Number:    trackInfo.Number,
				Year:      trackInfo.Year,
				Tags:      []string{},
			})
		}

		if len(tracks) >= 1 {
			p := tracks[0]

			probe, err := utils.ProbeTrack(p)
			if err != nil {
				log.Fatal("Failed to probe track", "err", err)
			}

			if tag, exists := probe.Tags["album"]; exists {
				importData.Album.Name = tag
			}

			if tag, exists := probe.Tags["album_artist"]; exists {
				importData.Album.Artists = parseArtist(tag)
			} else {
				if tag, exists := probe.Tags["artist"]; exists {
					importData.Album.Artists = parseArtist(tag)
				}
			}

			if tag, exists := probe.Tags["date"]; exists {
				match := dateRegex.FindStringSubmatch(tag)
				if len(match) > 0 {
					importData.Album.Year, _ = strconv.Atoi(match[1])
				}
			}
		}

		if len(images) > 0 {
			p := images[0]
			filename := path.Base(p)

			importData.Album.CoverArtFilename = filename
		}

		data, err := toml.Marshal(importData)
		if err != nil {
			log.Fatal("Failed to marshal init data", "err", err)
		}

		out := path.Join(dir, "import.toml")
		err = WriteFile(out, data)
		if err != nil {
			log.Fatal("Failed to write data", "err", err)
		}

		fmt.Println("Wrote:", out)
	},
}

func init() {
	importCmd.Flags().StringP("dir", "d", ".", "Directory to search")
	importInitCmd.Flags().StringP("dir", "d", ".", "Directory to search")

	importCmd.AddCommand(importInitCmd)

	rootCmd.AddCommand(importCmd)
}
