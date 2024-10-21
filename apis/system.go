package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/kr/pretty"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

type ExportArtist struct {
	Id      string
	Name    string
	Picture string
}

type ExportTrack struct {
	Id   string
	Name string

	AlbumId  string
	ArtistId string

	Number   int64
	Duration int64
	Year     int64

	ExportName       string
	OriginalFilename string
	MobileFilename   string

	Created int64

	Tags []string
}

type ExportAlbum struct {
	Id       string
	Name     string
	ArtistId string

	CoverArt string
	Year     int64
}

type Export struct {
	Artists []ExportArtist
	Albums  []ExportAlbum
	Tracks  []ExportTrack
}

func InstallSystemHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "GetSystemInfo",
			Path:     "/system/info",
			Method:   http.MethodGet,
			DataType: types.GetSystemInfo{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				return types.GetSystemInfo{
					Version: dwebble.Version,
					IsSetup: app.IsSetup(),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "SystemExport",
			Path:   "/system/export",
			Method: http.MethodPost,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				db := app.DB()

				ctx := context.TODO()

				export := Export{}

				artists, err := db.GetAllArtists(ctx)
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

				albums, err := db.GetAllAlbums(ctx, "", "", false)
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

				tracks, err := db.GetAllTracks(ctx, "", "", false)
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
						Number:           track.Number.Int64,
						Duration:         track.Duration.Int64,
						Year:             track.Year.Int64,
						ExportName:       track.ExportName,
						OriginalFilename: track.OriginalFilename,
						MobileFilename:   track.MobileFilename,
						Created:          track.Created,
						Tags:             utils.SplitString(track.Tags.String),
					}
				}

				pretty.Println(export)

				d, err := json.MarshalIndent(export, "", "  ")
				if err != nil {
					return nil, err
				}

				err = os.WriteFile(app.WorkDir().ExportFile(), d, 0644)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "SystemImport",
			Path:   "/system/import",
			Method: http.MethodPost,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				d, err := os.ReadFile(app.WorkDir().ExportFile())
				if err != nil {
					return nil, err
				}

				var export Export
				err = json.Unmarshal(d, &export)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				for _, artist := range export.Artists {
					_, err := db.CreateArtist(ctx, database.CreateArtistParams{
						Id:   artist.Id,
						Name: artist.Name,
						Picture: sql.NullString{
							String: artist.Picture,
							Valid:  artist.Picture != "",
						},
					})
					if err != nil {
						var e sqlite3.Error
						if errors.As(err, &e) {
							if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
								continue
							}
						}

						return nil, err
					}
				}

				for _, album := range export.Albums {
					_, err := db.CreateAlbum(ctx, database.CreateAlbumParams{
						Id:       album.Id,
						Name:     album.Name,
						ArtistId: album.ArtistId,
						CoverArt: sql.NullString{
							String: album.CoverArt,
							Valid:  album.CoverArt != "",
						},
						Year: sql.NullInt64{
							Int64: album.Year,
							Valid: album.Year != 0,
						},
						Available: true,
					})
					if err != nil {
						var e sqlite3.Error
						if errors.As(err, &e) {
							if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
								continue
							}
						}

						return nil, err
					}
				}

				for _, track := range export.Tracks {
					_, err := db.CreateTrack(ctx, database.CreateTrackParams{
						Id:       track.Id,
						Name:     track.Name,
						AlbumId:  track.AlbumId,
						ArtistId: track.ArtistId,
						Number: sql.NullInt64{
							Int64: track.Number,
							Valid: track.Number != 0,
						},
						Duration: sql.NullInt64{
							Int64: track.Duration,
							Valid: track.Duration != 0,
						},
						Year: sql.NullInt64{
							Int64: track.Year,
							Valid: track.Year != 0,
						},
						ExportName:       track.ExportName,
						OriginalFilename: track.OriginalFilename,
						MobileFilename:   track.MobileFilename,
						Created:          track.Created,
						Available:        true,
					})
					if err != nil {
						var e sqlite3.Error
						if errors.As(err, &e) {
							if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
								err := db.UpdateTrack(ctx, track.Id, database.TrackChanges{
									Name: types.Change[string]{
										Value:   track.Name,
										Changed: true,
									},
									AlbumId: types.Change[string]{
										Value:   track.AlbumId,
										Changed: true,
									},
									ArtistId: types.Change[string]{
										Value:   track.ArtistId,
										Changed: true,
									},
									Number: types.Change[sql.NullInt64]{
										Value: sql.NullInt64{
											Int64: track.Number,
											Valid: track.Number != 0,
										},
										Changed: true,
									},
									Duration: types.Change[sql.NullInt64]{
										Value: sql.NullInt64{
											Int64: track.Duration,
											Valid: track.Duration != 0,
										},
										Changed: false,
									},
									Year: types.Change[sql.NullInt64]{
										Value: sql.NullInt64{
											Int64: track.Year,
											Valid: track.Year != 0,
										},
										Changed: true,
									},
									ExportName: types.Change[string]{
										Value:   track.ExportName,
										Changed: true,
									},
									OriginalFilename: types.Change[string]{
										Value:   track.OriginalFilename,
										Changed: true,
									},
									MobileFilename: types.Change[string]{
										Value:   track.MobileFilename,
										Changed: true,
									},
									Created: types.Change[int64]{
										Value:   track.Created,
										Changed: true,
									},
									Available: types.Change[bool]{
										Value:   true,
										Changed: true,
									},
								})
								if err != nil {
									return nil, err
								}

								continue
							}
						}

						return nil, err
					}

					for _, tag := range track.Tags {
						slug := utils.Slug(tag)

						err := db.CreateTag(ctx, slug, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}

						err = db.AddTagToTrack(ctx, slug, track.Id)
						if err != nil {
							return nil, err
						}
					}

				}

				err = tx.Commit()
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		// TODO(patrik): Rename
		pyrin.ApiHandler{
			Name:   "Process",
			Path:   "/system/process",
			Method: http.MethodPost,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				tracks, err := app.DB().GetAllTracks(c.Request().Context(), "", "", false)
				if err != nil {
					return nil, err
				}

				// TODO(patrik): Use goroutines
				for _, track := range tracks {
					albumDir := app.WorkDir().Album(track.AlbumId)

					file := path.Join(albumDir.OriginalFiles(), track.OriginalFilename)

					filename := track.MobileFilename
					filename = strings.TrimSuffix(filename, path.Ext(filename))

					_, err := utils.ProcessMobileVersion(file, albumDir.MobileFiles(), filename)
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},
	)
}
