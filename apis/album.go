package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/faceair/jio"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/helper"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/validate"
)

func ConvertDBAlbum(c pyrin.Context, album database.Album) types.Album {
	var year *int64
	if album.Year.Valid {
		year = &album.Year.Int64
	}

	return types.Album{
		Id:         album.Id,
		Name:       album.Name,
		CoverArt:   utils.ConvertAlbumCoverURL(c, album.Id, album.CoverArt),
		ArtistId:   album.ArtistId,
		ArtistName: album.ArtistName,
		Year:       year,
		Created:    album.Created,
		Updated:    album.Updated,
	}
}

var _ pyrin.Body = (*PatchAlbumBody)(nil)

type PatchAlbumBody struct {
	Name       *string `json:"name"`
	ArtistId   *string `json:"artistId"`
	ArtistName *string `json:"artistName"`
	Year       *int64  `json:"year"`
}

func (b PatchAlbumBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

var _ pyrin.Body = (*CreateAlbumBody)(nil)

type CreateAlbumBody struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

func (b CreateAlbumBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func (CreateAlbumBody) Schema() jio.Schema {
	panic("unimplemented")
}

type CreateAlbum struct {
	AlbumId string `json:"albumId"`
}

var _ pyrin.Body = (*UploadTracksBody)(nil)

type UploadTracksBody struct {
	ForceExtractNumber bool `json:"forceExtractNumber"`
}

func (b UploadTracksBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func InstallAlbumHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "GetAlbums",
			Path:     "/albums",
			Method:   http.MethodGet,
			DataType: types.GetAlbums{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				filter := q.Get("filter")
				sort := q.Get("sort")

				albums, err := app.DB().GetAllAlbums(c.Request().Context(), filter, sort)
				if err != nil {
					return nil, err
				}

				res := types.GetAlbums{
					Albums: make([]types.Album, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = ConvertDBAlbum(c, album)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "SearchAlbums",
			Path:     "/albums/search",
			Method:   http.MethodGet,
			DataType: types.GetAlbums{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				query := strings.TrimSpace(q.Get("query"))

				var err error
				var albums []database.Album

				if query != "" {
					albums, err = app.DB().SearchAlbums(query)
					if err != nil {
						return nil, err
					}
				}

				res := types.GetAlbums{
					Albums: make([]types.Album, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = ConvertDBAlbum(c, album)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetAlbumById",
			Method:   http.MethodGet,
			Path:     "/albums/:id",
			DataType: types.GetAlbumById{},
			Errors:   []pyrin.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				album, err := app.DB().GetAlbumById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				return types.GetAlbumById{
					Album: ConvertDBAlbum(c, album),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetAlbumTracks",
			Method:   http.MethodGet,
			Path:     "/albums/:id/tracks",
			DataType: types.GetAlbumTracksById{},
			Errors:   []pyrin.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				album, err := app.DB().GetAlbumById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				tracks, err := app.DB().GetTracksByAlbum(c.Request().Context(), album.Id)
				if err != nil {
					return nil, err
				}

				res := types.GetAlbumTracksById{
					Tracks: make([]types.Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},
	)

	group.Register(
		pyrin.ApiHandler{
			Name:     "EditAlbum",
			Method:   http.MethodPatch,
			Path:     "/albums/:id",
			BodyType: PatchAlbumBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				var body PatchAlbumBody
				d := json.NewDecoder(c.Request().Body)
				err := d.Decode(&body)
				if err != nil {
					return nil, err
				}

				album, err := app.DB().GetAlbumById(c.Request().Context(), id)
				if err != nil {
					// TODO(patrik): Handle error
					return nil, err
				}

				var name types.Change[string]
				if body.Name != nil {
					n := strings.TrimSpace(*body.Name)
					name.Value = n
					name.Changed = n != album.Name
				}

				ctx := context.TODO()

				var artistId types.Change[string]
				if body.ArtistId != nil {
					artistId.Value = *body.ArtistId
					artistId.Changed = *body.ArtistId != album.ArtistId
				} else if body.ArtistName != nil {
					artistName := strings.TrimSpace(*body.ArtistName)

					artist, err := app.DB().GetArtistByName(ctx, artistName)
					if err != nil {
						if errors.Is(err, database.ErrItemNotFound) {
							artist, err = helper.CreateArtist(ctx, app.DB(), app.WorkDir(), artistName)
							if err != nil {
								return nil, err
							}
						} else {
							return nil, err
						}
					}

					artistId.Value = artist.Id
					artistId.Changed = artist.Id != album.ArtistId
				}

				var year types.Change[sql.NullInt64]
				if body.Year != nil {
					year.Value = sql.NullInt64{
						Int64: *body.Year,
						Valid: *body.Year != 0,
					}
					year.Changed = *body.Year != album.Year.Int64
				}

				err = app.DB().UpdateAlbum(ctx, album.Id, database.AlbumChanges{
					Name:     name,
					ArtistId: artistId,
					Year:     year,
				})
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeleteAlbum",
			Method: http.MethodDelete,
			Path:   "/albums/:id",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				err = db.RemoveAlbumTracks(c.Request().Context(), id)
				if err != nil {
					return nil, fmt.Errorf("Failed to remove album tracks: %w", err)
				}

				err = db.RemoveAlbum(c.Request().Context(), id)
				if err != nil {
					return nil, fmt.Errorf("Failed to remove album: %w", err)
				}

				err = tx.Commit()
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "CreateAlbum",
			Method:   http.MethodPost,
			Path:     "/albums",
			DataType: CreateAlbum{},
			BodyType: CreateAlbumBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Validate and trim body
				body, err := Body[CreateAlbumBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				artist, err := app.DB().GetArtistByName(ctx, body.Artist)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						artist, err = helper.CreateArtist(ctx, app.DB(), app.WorkDir(), body.Artist)
						if err != nil {
							return nil, err
						}
					} else {
						return nil, err
					}
				}

				album, err := helper.CreateAlbum(ctx, app.DB(), app.WorkDir(), database.CreateAlbumParams{
					Name:     body.Name,
					ArtistId: artist.Id,
				})
				if err != nil {
					return nil, err
				}

				return CreateAlbum{
					AlbumId: album.Id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:        "ChangeAlbumCover",
			Method:      http.MethodPost,
			Path:        "/albums/:id/cover",
			RequireForm: true,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				err := c.Request().ParseMultipartForm(defaultMemory)
				if err != nil {
					return nil, err
				}

				form := c.Request().MultipartForm

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				ctx := context.TODO()
				album, err := db.GetAlbumById(ctx, id)
				if err != nil {
					// TODO(patrik): Handle error
					return nil, err
				}

				albumDir := app.WorkDir().Album(album.Id)

				// TODO(patrik): Check content-type
				// TODO(patrik): Return error if there is no files attached
				coverArt := form.File["cover"]
				if len(coverArt) > 0 {
					f := coverArt[0]

					ext := path.Ext(f.Filename)
					filename := "cover-original" + ext

					dst := path.Join(albumDir, filename)

					file, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
					if err != nil {
						return nil, err
					}
					defer file.Close()

					ff, err := f.Open()
					if err != nil {
						return nil, err
					}
					defer ff.Close()

					_, err = io.Copy(file, ff)
					if err != nil {
						return nil, err
					}

					i := path.Join(albumDir, "cover-128.png")
					err = utils.CreateResizedImage(dst, i, 128)
					if err != nil {
						return nil, err
					}

					i = path.Join(albumDir, "cover-256.png")
					err = utils.CreateResizedImage(dst, i, 256)
					if err != nil {
						return nil, err
					}

					i = path.Join(albumDir, "cover-512.png")
					err = utils.CreateResizedImage(dst, i, 512)
					if err != nil {
						return nil, err
					}

					err = db.UpdateAlbum(ctx, album.Id, database.AlbumChanges{
						CoverArt: types.Change[sql.NullString]{
							Value: sql.NullString{
								String: filename,
								Valid:  true,
							},
							Changed: true,
						},
					})

					if err != nil {
						return nil, err
					}
				}

				err = tx.Commit()
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:        "UploadTracks",
			Method:      http.MethodPost,
			Path:        "/albums/:id/upload",
			BodyType:    UploadTracksBody{},
			RequireForm: true,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				err := c.Request().ParseMultipartForm(defaultMemory)
				if err != nil {
					return nil, err
				}

				form := c.Request().MultipartForm

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				ctx := context.TODO()
				album, err := db.GetAlbumById(ctx, id)
				if err != nil {
					// TODO(patrik): Handle error
					return nil, err
				}

				copyFormFileToTemp := func(file *multipart.FileHeader) (string, error) {
					ext := path.Ext(file.Filename)

					src, err := file.Open()
					if err != nil {
						return "", err
					}

					// TODO(patrik): Copy the file to $trackDir/raw.flac instead
					dst, err := os.CreateTemp("", "track.*"+ext)
					if err != nil {
						return "", err
					}
					defer dst.Close()

					_, err = io.Copy(dst, src)
					if err != nil {
						return "", err
					}

					return dst.Name(), nil
				}

				var body UploadTracksBody
				if len(form.Value["body"]) > 0 {
					bodyStr := form.Value["body"][0]
					err := json.Unmarshal(([]byte)(bodyStr), &body)
					if err != nil {
						return nil, err
					}
				}

				files := form.File["files"]
				for _, f := range files {
					ext := path.Ext(f.Filename)
					name := strings.TrimSuffix(f.Filename, ext)

					filename, err := copyFormFileToTemp(f)
					if err != nil {
						return nil, err
					}
					defer os.Remove(filename)

					data := helper.ImportTrackData{
						AlbumId:            album.Id,
						ArtistId:           album.ArtistId,
						Name:               name,
						Filename:           filename,
						ForceExtractNumber: false,
					}
					_, err = helper.ImportTrack(ctx, app.DB(), app.WorkDir(), data)
					if err != nil {
						return nil, err
					}
				}

				err = tx.Commit()
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
