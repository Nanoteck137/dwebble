package apis

import (
	"context"
	"net/http"
	"path"
	"strings"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/pyrin"
)

type GetSystemInfo struct {
	Version string `json:"version"`
}

type ExportArtist struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type ExportTrack struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	AlbumId  string `json:"albumId"`
	ArtistId string `json:"artistId"`

	Duration int64 `json:"duration"`
	Number   int64 `json:"number"`
	Year     int64 `json:"year"`

	OriginalFilename string `json:"originalFilename"`
	MobileFilename   string `json:"mobileFilename"`

	Created int64 `json:"created"`

	Tags []string `json:"tags"`
}

type ExportAlbum struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ArtistId string `json:"artistId"`

	CoverArt string `json:"coverArt"`
	Year     int64  `json:"year"`
}

type Export struct {
	Artists []ExportArtist `json:"artists"`
	Albums  []ExportAlbum  `json:"albums"`
	Tracks  []ExportTrack  `json:"tracks"`
}

func InstallSystemHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetSystemInfo",
			Path:         "/system/info",
			Method:       http.MethodGet,
			ResponseType: GetSystemInfo{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				return GetSystemInfo{
					Version: dwebble.Version,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SystemExport",
			Path:         "/system/export",
			Method:       http.MethodPost,
			ResponseType: Export{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				db := app.DB()

				ctx := context.TODO()

				export := Export{}

				artists, err := db.GetAllArtists(ctx, "", "")
				if err != nil {
					return nil, err
				}

				export.Artists = make([]ExportArtist, len(artists))
				for i, artist := range artists {
					export.Artists[i] = ExportArtist{
						Id:      artist.Id,
						Name:    artist.Name,
						Picture: artist.Picture.String,
					}
				}

				albums, err := db.GetAllAlbums(ctx, "", "")
				if err != nil {
					return nil, err
				}

				export.Albums = make([]ExportAlbum, len(albums))
				for i, album := range albums {
					export.Albums[i] = ExportAlbum{
						Id:       album.Id,
						Name:     album.Name,
						ArtistId: album.ArtistId,
						CoverArt: album.CoverArt.String,
						Year:     album.Year.Int64,
					}
				}

				tracks, err := db.GetAllTracks(ctx, "", "")
				if err != nil {
					return nil, err
				}

				export.Tracks = make([]ExportTrack, len(tracks))
				for i, track := range tracks {
					export.Tracks[i] = ExportTrack{
						Id:               track.Id,
						Name:             track.Name,
						AlbumId:          track.AlbumId,
						ArtistId:         track.ArtistId,
						Duration:         track.Duration,
						Number:           track.Number.Int64,
						Year:             track.Year.Int64,
						OriginalFilename: track.OriginalFilename,
						MobileFilename:   track.MobileFilename,
						Created:          track.Created,
						Tags:             utils.SplitString(track.Tags.String),
					}
				}

				return export, nil
			},
		},

		// TODO(patrik): Rename
		pyrin.ApiHandler{
			Name:   "Process",
			Path:   "/system/process",
			Method: http.MethodPost,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				tracks, err := app.DB().GetAllTracks(c.Request().Context(), "", "")
				if err != nil {
					return nil, err
				}

				// TODO(patrik): Use goroutines
				for _, track := range tracks {
					// TODO(patrik): Fix
					// albumDir := app.WorkDir().Album(track.AlbumId)

					// TODO(patrik): Fix
					// file := path.Join(albumDir.OriginalFiles(), track.OriginalFilename)

					filename := track.MobileFilename
					filename = strings.TrimSuffix(filename, path.Ext(filename))

					// TODO(patrik): Fix
					// _, err := utils.ProcessMobileVersion(file, albumDir.MobileFiles(), filename)
					// if err != nil {
					// 	return nil, err
					// }
				}

				return nil, nil
			},
		},
	)
}
