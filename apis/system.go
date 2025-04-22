package apis

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

type GetSystemInfo struct {
	Version string `json:"version"`
}

func FixMetadata(metadata *library.Metadata) error {
	return nil
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
			Name:   "RefillSearch",
			Path:   "/system/search",
			Method: http.MethodPost,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()
				err = app.DB().RefillSearchTables(ctx)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SyncLibrary",
			Method:       http.MethodPost,
			Path:         "/system/library",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				albums, err := library.FindAlbums("/Volumes/media/test/Ado/unravel (live)/")
				if err != nil {
					return nil, err
				}

				pretty.Println(albums)

				ctx := context.TODO()

				err = EnsureUnknownArtistExists(ctx, app.DB(), app.WorkDir())
				if err != nil {
					return nil, err
				}

				artists := map[string]string{}

				getOrCreateArtist := func(name string) (string, error) {
					slug := utils.Slug(name)

					if artist, exists := artists[slug]; exists {
						return artist, nil
					}

					dbArtist, err := app.DB().GetArtistBySlug(ctx, slug)
					if err != nil {
						if errors.Is(err, database.ErrItemNotFound) {
							dbArtist, err = app.DB().CreateArtist(ctx, database.CreateArtistParams{
								Slug: slug,
								Name: name,
							})
							if err != nil {
								return "", err
							}
						} else {
							return "", err
						}
					}

					artists[slug] = dbArtist.Id
					return dbArtist.Id, nil
				}

				for _, album := range albums {
					err := FixMetadata(&album.Metadata)
					if err != nil {
						return nil, err
					}

					dbAlbum, err := app.DB().GetAlbumByPath(ctx, album.Path)
					if err != nil {
						if errors.Is(err, database.ErrItemNotFound) {
							artist, err := getOrCreateArtist(album.Metadata.Album.Artists[0])
							if err != nil {
								return nil, err
							}

							dbAlbum, err = app.DB().CreateAlbum(ctx, database.CreateAlbumParams{
								Path:     album.Path,
								Name:     album.Metadata.Album.Name,
								ArtistId: artist,
							})

							if err != nil {
								return nil, err
							}
						} else {
							return nil, err
						}
					}

					pretty.Println(dbAlbum)

					for _, track := range album.Metadata.Tracks {
						artist, err := getOrCreateArtist(track.Artists[0])
						if err != nil {
							return nil, err
						}

						dbTrack, err := app.DB().GetTrackByPath(ctx, track.File)
						if err != nil {
							if errors.Is(err, database.ErrItemNotFound) {
								// TODO(patrik): Get the media type
								probeResult, err := utils.ProbeTrack(track.File)
								if err != nil {
									return nil, err
								}

								_, err = app.DB().CreateTrack(ctx, database.CreateTrackParams{
									Path:         track.File,
									ModifiedTime: track.ModifiedTime,
									Filename:     track.File,
									MediaType:    "",
									Name:         track.Name,
									OtherName:    sql.NullString{},
									AlbumId:      dbAlbum.Id,
									ArtistId:     artist,
									Duration:     int64(probeResult.Duration),
									Number: sql.NullInt64{
										Int64: int64(track.Number),
										Valid: track.Number != 0,
									},
									Year: sql.NullInt64{
										Int64: int64(track.Year),
										Valid: track.Year != 0,
									},
								})
								if err != nil {
									return nil, err
								}

								continue
							}
						}

						// TODO(patrik): Check modified time and probe again
						// TODO(patrik): Update track

						changes := database.TrackChanges{}

						if track.ModifiedTime > dbTrack.ModifiedTime {
							probeResult, err := utils.ProbeTrack(track.File)
							if err != nil {
								return nil, err
							}

							dur := int64(probeResult.Duration)
							changes.Duration = types.Change[int64]{
								Value:   dur,
								Changed: dur != dbTrack.Duration,
							}
						}

						changes.Name = types.Change[string]{
							Value:   track.Name,
							Changed: dbTrack.Name != track.Name,
						}

						changes.Year = types.Change[sql.NullInt64]{
							Value:   sql.NullInt64{
								Int64: int64(track.Year),
								Valid: track.Year != 0,
							},
							Changed: dbTrack.Year.Int64 != int64(track.Year),
						}

						pretty.Println(changes)

						err = app.DB().UpdateTrack(ctx, dbTrack.Id, changes)
						if err != nil {
							return nil, err
						}

						pretty.Println(dbTrack)
					}
				}

				return nil, nil
			},
		},
	)
}
