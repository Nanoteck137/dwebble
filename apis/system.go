package apis

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
)

type GetSystemInfo struct {
	Version string `json:"version"`
}

// TODO(patrik): Add testing for this
func FixMetadata(metadata *library.Metadata) error {
	album := &metadata.Album

	album.Name = transform.String(album.Name)

	if album.Year == 0 {
		album.Year = metadata.General.Year
	}

	if len(album.Artists) == 0 {
		album.Artists = []string{UNKNOWN_ARTIST_NAME}
	}

	for i := range metadata.Tracks {
		t := &metadata.Tracks[i]

		if t.Year == 0 {
			t.Year = metadata.General.Year
		}

		t.Name = transform.String(t.Name)

		t.Tags = append(t.Tags, metadata.General.Tags...)
		t.Tags = append(t.Tags, metadata.General.TrackTags...)

		if len(t.Artists) == 0 {
			t.Artists = []string{UNKNOWN_ARTIST_NAME}
		}

		for i, tag := range t.Tags {
			t.Tags[i] = utils.Slug(strings.TrimSpace(tag))
		}
	}

	// err := validate.ValidateStruct(&metadata.Album,
	// 	validate.Field(&metadata.Album.Name, validate.Required),
	// 	validate.Field(&metadata.Album.Artists, validate.Length(1, 0)),
	// )
	// if err != nil {
	// 	return err
	// }
	//
	// for _, track := range metadata.Tracks {
	// 	err := validate.ValidateStruct(&track,
	// 		validate.Field(&track.File, validate.Required),
	// 		validate.Field(&track.Name, validate.Required),
	// 		validate.Field(&track.Artists, validate.Length(1, 0)),
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

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
				// TODO(patrik):
				//  - Handle Errors
				//  - Handle Single Album Syncing
				//  - Handle album modified syncing

				// TODO(patrik): Check for duplicated ids
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

					dbAlbum, err := app.DB().GetAlbumById(ctx, album.Metadata.Album.Id)
					if err != nil {
						if errors.Is(err, database.ErrItemNotFound) {
							artist, err := getOrCreateArtist(album.Metadata.Album.Artists[0])
							if err != nil {
								return nil, err
							}

							dbAlbum, err = app.DB().CreateAlbum(ctx, database.CreateAlbumParams{
								Id:       album.Metadata.Album.Id,
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

					for _, track := range album.Metadata.Tracks {
						artist, err := getOrCreateArtist(track.Artists[0])
						if err != nil {
							return nil, err
						}

						dbTrack, err := app.DB().GetTrackById(ctx, track.Id)
						if err != nil {
							if errors.Is(err, database.ErrItemNotFound) {
								// TODO(patrik): Get the media type
								probeResult, err := utils.ProbeTrack(track.File)
								if err != nil {
									return nil, err
								}

								_, err = app.DB().CreateTrack(ctx, database.CreateTrackParams{
									Id:           track.Id,
									Filename:     track.File,
									ModifiedTime: track.ModifiedTime,
									MediaType:    probeResult.MediaType,
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

							changes.MediaType = types.Change[types.MediaType]{
								Value:   probeResult.MediaType,
								Changed: probeResult.MediaType != dbTrack.MediaType,
							}

							changes.ModifiedTime = types.Change[int64]{
								Value:   track.ModifiedTime,
								Changed: track.ModifiedTime != dbTrack.ModifiedTime,
							}
						}

						// TODO(patrik): Implement all the changes here

						changes.Filename = types.Change[string]{
							Value:   track.File,
							Changed: dbTrack.Filename != track.File,
						}

						changes.Name = types.Change[string]{
							Value:   track.Name,
							Changed: dbTrack.Name != track.Name,
						}

						changes.Year = types.Change[sql.NullInt64]{
							Value: sql.NullInt64{
								Int64: int64(track.Year),
								Valid: track.Year != 0,
							},
							Changed: dbTrack.Year.Int64 != int64(track.Year),
						}

						err = app.DB().UpdateTrack(ctx, dbTrack.Id, changes)
						if err != nil {
							return nil, err
						}
					}
				}

				return nil, nil
			},
		},
	)
}
