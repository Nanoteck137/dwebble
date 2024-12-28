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
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/cmd/dwebble-cli/api"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
	"github.com/pelletier/go-toml/v2"
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

type ImportDataAlbum struct {
	Name             string   `toml:"name"`
	OtherName        string   `toml:"other_name"`
	Artists          []string `toml:"artists"`
	Year             int      `toml:"year"`
	Tags             []string `toml:"tags"`
	CoverArtFilename string   `toml:"coverArtFilename"`
}

func (d *ImportDataAlbum) Transform() {
	d.Name = transform.String(d.Name)
	d.OtherName = transform.String(d.OtherName)

	d.Artists = TransformStringArray(d.Artists)
}

func (d ImportDataAlbum) Validate() error {
	return validate.ValidateStruct(&d,
		validate.Field(&d.Name, validate.Required),
		validate.Field(&d.Artists, validate.Required),
	)
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

func (d *ImportDataTrack) Transform() {
	d.Name = transform.String(d.Name)
	d.OtherName = transform.String(d.OtherName)
	d.Artists = TransformStringArray(d.Artists)
}

func (d ImportDataTrack) Validate() error {
	return validate.ValidateStruct(&d,
		validate.Field(&d.Filename, validate.Required),
		validate.Field(&d.Name, validate.Required),
		validate.Field(&d.Artists, validate.Required),
	)
}

type ImportData struct {
	UseAlbumTagsForTracks bool              `toml:"use_album_tags_for_tracks"`
	Album                 ImportDataAlbum   `toml:"album"`
	Tracks                []ImportDataTrack `toml:"tracks"`
}

func (d *ImportData) Transform() {
	d.Album.Transform()

	for i := range d.Tracks {
		d.Tracks[i].Transform()
	}
}

func (d ImportData) Validate() error {
	return validate.ValidateStruct(&d,
		validate.Field(&d.Tracks),
	)
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

func (c *Context) UploadTrack(albumId string, dir string, data ImportData, track ImportDataTrack) error {
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
			AlbumId:          albumId,
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

var importCmd = &cobra.Command{
	Use: "import",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		data, err := getImportData(dir)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Testing")

				huh.NewForm(
					huh.NewGroup(
						huh.NewInput().Title("Album Name"),
						huh.NewInput().Title("Other Name"),
						huh.NewInput().Title("Artists"),
						huh.NewInput().Title("Year"),
					),
					huh.NewGroup(
						huh.NewConfirm(),
					),
				).Run()

				return
			}
			log.Fatal("Failed to get import data", "err", err)
		}

		data.Transform()

		err = data.Validate()
		if err != nil {
			log.Fatal("Validation failed", "err", err)
		}

		// pretty.Println(data)

		// TODO(patrik):
		//  - [x] Find albums
		//    - [ ] Find album by id
		//    - [x] Updated album with values
		//    - [ ] Find matching tracks
		//      - [ ] Update tracks with values
		//    - [ ] Give the user the option to replace all tracks inside
		//      the found album

		client := api.New("http://localhost:3000")

		c := Context{
			client:  client,
			artists: map[string]string{},
		}

		const ActionSearchAlbum = "search_album"
		const ActionAlbumById = "album_by_id"
		const ActionCreateAlbum = "create_album"

		{
			var selected string
			err := huh.NewSelect[string]().
				Title("Select action").
				Options(
					huh.NewOption[string]("Create new Album", ActionCreateAlbum),
					huh.NewOption[string]("Search for Album", ActionSearchAlbum),
					huh.NewOption[string]("Use Album ID", ActionAlbumById),
				).
				Value(&selected).
				Run()
			if err != nil {
				log.Fatal("Failed", "err", err)
			}

			switch selected {
			case ActionSearchAlbum:
				log.Fatal("Not implemented yet")
			case ActionAlbumById:
				log.Fatal("Not implemented yet")
			case ActionCreateAlbum:

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

				albumId := album.AlbumId
				fmt.Printf("Created album: (%s)\n", album.AlbumId)

				// TODO(patrik): check 'data.Album.CoverArtFilename'
				fmt.Println("Album cover set")

				// TODO(patrik): Create user option
				fmt.Printf("Setting album cover art...")
				err = c.SetAlbumCover(albumId, path.Join(dir, data.Album.CoverArtFilename))
				if err != nil {
					log.Fatal("Failed to set album cover art", "err", err)
				}
				fmt.Printf("done\n")

				for i, track := range data.Tracks {
					fmt.Printf("Uploading track (%02d/%02d): %s...", i+1, len(data.Tracks), track.Name)
					err = c.UploadTrack(albumId, dir, data, track)
					if err != nil {
						log.Fatal("Failed to upload track", "err", err)
					}
					fmt.Printf("done\n")
				}
			default:
				log.Fatal("Unimplemented action", "action", selected)
			}
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
		// NOTE(patrik): Default behavior
		importData.UseAlbumTagsForTracks = true

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

		importData.Transform()

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

var importTestCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		client := api.New("http://localhost:3000")

		c := Context{
			client:  client,
			artists: map[string]string{},
		}

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

		pretty.Println(images)
		pretty.Println(tracks)

		if len(tracks) >= 1 {
			p := tracks[0]

			name := "Unknown Album name"
			var artists []string
			year := 0

			probe, err := utils.ProbeTrack(p)
			if err != nil {
				log.Fatal("Failed to probe track", "err", err)
			}

			if tag, exists := probe.Tags["album"]; exists {
				name = tag
			}

			if tag, exists := probe.Tags["album_artist"]; exists {
				artists = parseArtist(tag)
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

			fmt.Printf("Created album: (%s)\n", album.AlbumId)

			if len(images) > 0 {
				err = c.SetAlbumCover(album.AlbumId, images[0])
				if err != nil {
					log.Fatal("Failed to set album cover", "err", err)
				}
			}

			for _, p := range tracks {
				fmt.Printf("p: %v\n", p)

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

				artistId, featuringArtists, err := c.GetOrCreateMultipleArtists(artists)
				if err != nil {
					log.Fatal("Failed to get/create album artist", "err", err)
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
					log.Fatal("Failed to upload track")
				}
			}

			openbrowser(fmt.Sprintf("http://localhost:5173/albums/%s", album.AlbumId))
		}
	},
}

func init() {
	importCmd.Flags().StringP("dir", "d", ".", "Directory to search")
	importInitCmd.Flags().StringP("dir", "d", ".", "Directory to search")
	importTestCmd.Flags().StringP("dir", "d", ".", "Directory to search")

	importCmd.AddCommand(importInitCmd)
	importCmd.AddCommand(importTestCmd)

	rootCmd.AddCommand(importCmd)
}
