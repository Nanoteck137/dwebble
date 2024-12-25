package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/helper"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
)

type Album struct {
	Id   string `json:"id"`
	Name Name   `json:"name"`

	Year *int64 `json:"year"`

	CoverArt types.Images `json:"coverArt"`

	ArtistId   string `json:"artistId"`
	ArtistName Name   `json:"artistName"`

	Tags             []string     `json:"tags"`
	FeaturingArtists []ArtistInfo `json:"featuringArtists"`
	AllArtists       []ArtistInfo `json:"allArtists"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

func ConvertDBAlbum(c pyrin.Context, album database.Album) Album {
	allArtists := make([]ArtistInfo, len(album.FeaturingArtists)+1)

	allArtists[0] = ArtistInfo{
		Id: album.ArtistId,
		Name: Name{
			Default: album.ArtistName,
			Other:   ConvertSqlNullString(album.ArtistOtherName),
		},
	}

	featuringArtists := ConvertDBExtraArtists(album.FeaturingArtists)
	for i, v := range featuringArtists {
		allArtists[i+1] = v
	}

	return Album{
		Id: album.Id,
		Name: Name{
			Default: album.Name,
			Other:   ConvertSqlNullString(album.OtherName),
		},
		Year:     ConvertSqlNullInt64(album.Year),
		CoverArt: utils.ConvertAlbumCoverURL(c, album.Id, album.CoverArt),
		ArtistId: album.ArtistId,
		ArtistName: Name{
			Default: album.ArtistName,
			Other:   ConvertSqlNullString(album.ArtistOtherName),
		},
		Tags:             utils.SplitString(album.Tags.String),
		FeaturingArtists: featuringArtists,
		AllArtists:       allArtists,
		Created:          album.Created,
		Updated:          album.Updated,
	}
}

type GetAlbums struct {
	Albums []Album `json:"albums"`
}

type GetAlbumById struct {
	Album
}

type GetAlbumTracks struct {
	Tracks []Track `json:"tracks"`
}

type EditAlbumBody struct {
	Name             *string   `json:"name,omitempty"`
	OtherName        *string   `json:"otherName,omitempty"`
	ArtistId         *string   `json:"artistId,omitempty"`
	ArtistName       *string   `json:"artistName,omitempty"`
	Year             *int64    `json:"year,omitempty"`
	Tags             *[]string `json:"tags,omitempty"`
	FeaturingArtists *[]string `json:"featuringArtists,omitempty"`
}

func (b *EditAlbumBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
	b.OtherName = transform.StringPtr(b.OtherName)
	b.ArtistName = transform.StringPtr(b.ArtistName)
	b.Tags = DiscardEntriesStringArrayPtr(b.Tags)
	b.FeaturingArtists = DiscardEntriesStringArrayPtr(b.FeaturingArtists)
}

func (b EditAlbumBody) Validate() error {
	checkBoth := validate.By(func(value interface{}) error {
		if b.ArtistName != nil && b.ArtistId != nil {
			return errors.New("both 'artistId' and 'artistName' cannot be set")
		}

		return nil
	})

	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
		// validate.Field(&b.OtherName, validate.Required.When(b.OtherName != nil)),
		validate.Field(&b.ArtistId, checkBoth, validate.Required.When(b.ArtistId != nil)),
		validate.Field(&b.ArtistName, checkBoth, validate.Required.When(b.ArtistName != nil)),
		validate.Field(&b.Year, validate.Min(0)),
	)
}

type CreateAlbum struct {
	AlbumId string `json:"albumId"`
}

// TODO(patrik): Add ArtistId
type CreateAlbumBody struct {
	Name     string `json:"name"`
	ArtistId string `json:"artistId"`
}

func (b *CreateAlbumBody) Transform() {
	b.Name = transform.String(b.Name)
	b.ArtistId = transform.String(b.ArtistId)
}

func (b CreateAlbumBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
		validate.Field(&b.ArtistId, validate.Required),
	)
}

func InstallAlbumHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetAlbums",
			Path:         "/albums",
			Method:       http.MethodGet,
			ResponseType: GetAlbums{},
			Errors:       []pyrin.ErrorType{ErrTypeInvalidFilter, ErrTypeInvalidSort},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				filter := q.Get("filter")
				sort := q.Get("sort")

				albums, err := app.DB().GetAllAlbums(c.Request().Context(), filter, sort)
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					if errors.Is(err, database.ErrInvalidSort) {
						return nil, InvalidSort(err)
					}

					return nil, err
				}

				res := GetAlbums{
					Albums: make([]Album, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = ConvertDBAlbum(c, album)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SearchAlbums",
			Path:         "/albums/search",
			Method:       http.MethodGet,
			ResponseType: GetAlbums{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				query := strings.TrimSpace(q.Get("query"))

				albums, err := app.DB().SearchAlbums(query)
				if err != nil {
					return nil, err
				}

				res := GetAlbums{
					Albums: make([]Album, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = ConvertDBAlbum(c, album)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetAlbumById",
			Method:       http.MethodGet,
			Path:         "/albums/:id",
			ResponseType: GetAlbumById{},
			Errors:       []pyrin.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				album, err := app.DB().GetAlbumById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				return GetAlbumById{
					Album: ConvertDBAlbum(c, album),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetAlbumTracks",
			Method:       http.MethodGet,
			Path:         "/albums/:id/tracks",
			ResponseType: GetAlbumTracks{},
			Errors:       []pyrin.ErrorType{ErrTypeAlbumNotFound},
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

				res := GetAlbumTracks{
					Tracks: make([]Track, len(tracks)),
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
			BodyType: EditAlbumBody{},
			Errors:   []pyrin.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[EditAlbumBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				album, err := app.DB().GetAlbumById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				changes := database.AlbumChanges{}

				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != album.Name,
					}
				}

				if body.OtherName != nil {
					changes.OtherName = types.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.OtherName,
							Valid:  *body.OtherName != "",
						},
						Changed: *body.OtherName != album.OtherName.String,
					}
				}

				if body.ArtistId != nil {
					changes.ArtistId = types.Change[string]{
						Value:   *body.ArtistId,
						Changed: *body.ArtistId != album.ArtistId,
					}
				} else if body.ArtistName != nil {
					artistName := *body.ArtistName

					// TODO(patrik): Move to helper?
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

					changes.ArtistId = types.Change[string]{
						Value:   artist.Id,
						Changed: artist.Id != album.ArtistId,
					}
				}

				if body.Year != nil {
					changes.Year = types.Change[sql.NullInt64]{
						Value: sql.NullInt64{
							Int64: *body.Year,
							Valid: *body.Year != 0,
						},
						Changed: *body.Year != album.Year.Int64,
					}
				}

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				err = db.UpdateAlbum(ctx, album.Id, changes)
				if err != nil {
					return nil, err
				}

				if body.Tags != nil {
					tags := *body.Tags

					err = db.RemoveAllTagsFromAlbum(ctx, album.Id)
					if err != nil {
						return nil, err
					}

					for _, tag := range tags {
						slug := utils.Slug(tag)

						err := db.CreateTag(ctx, slug)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}

						err = db.AddTagToAlbum(ctx, slug, album.Id)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}
				}

				if body.FeaturingArtists != nil {
					featuringArtists := *body.FeaturingArtists

					err := db.RemoveAllAlbumFeaturingArtists(ctx, album.Id)
					if err != nil {
						return nil, err
					}

					for _, artistId := range featuringArtists {
						artist, err := db.GetArtistById(ctx, artistId)
						if err != nil {
							if errors.Is(err, database.ErrItemNotFound) {
								// TODO(patrik): Should we just be
								// silently continuing
								// NOTE(patrik): If we don't find the artist
								// we just silently continue
								continue
							}

							return nil, err
						}

						err = db.AddFeaturingArtistToAlbum(ctx, album.Id, artist.Id)
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

		pyrin.ApiHandler{
			Name:   "DeleteAlbum",
			Method: http.MethodDelete,
			Path:   "/albums/:id",
			Errors: []pyrin.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				album, err := app.DB().GetAlbumById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				tracks, err := db.GetTracksByAlbum(ctx, album.Id)
				if err != nil {
					return nil, err
				}

				// TODO(patrik): This is duplicated code from the tracks api move to helper
				for _, track := range tracks {
					dir := app.WorkDir().Track(track.Id)
					targetName := fmt.Sprintf("track-%s-%d", track.Id, time.Now().UnixMilli())
					target := path.Join(app.WorkDir().Trash(), targetName)

					err = os.Rename(dir, target)
					if err != nil && !os.IsNotExist(err) {
						return nil, err
					}

					err = app.DB().RemoveTrack(ctx, track.Id)
					if err != nil {
						return nil, err
					}
				}

				dir := app.WorkDir().Album(album.Id)
				targetName := fmt.Sprintf("album-%s-%d", album.Id, time.Now().UnixMilli())
				target := path.Join(app.WorkDir().Trash(), targetName)

				err = os.Rename(dir, target)
				if err != nil && !os.IsNotExist(err) {
					return nil, err
				}

				err = db.RemoveAlbum(c.Request().Context(), id)
				if err != nil {
					return nil, err
				}

				err = tx.Commit()
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateAlbum",
			Method:       http.MethodPost,
			Path:         "/albums",
			ResponseType: CreateAlbum{},
			BodyType:     CreateAlbumBody{},
			Errors:       []pyrin.ErrorType{ErrTypeArtistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[CreateAlbumBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				artist, err := app.DB().GetArtistById(ctx, body.ArtistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ArtistNotFound()
					}

					return nil, err
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

		pyrin.FormApiHandler{
			Name:   "ChangeAlbumCover",
			Method: http.MethodPost,
			Path:   "/albums/:id/cover",
			Spec: pyrin.FormSpec{
				Files: map[string]pyrin.FormFileSpec{
					"cover": {
						NumExpected: 1,
					},
				},
			},
			Errors: []pyrin.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				ctx := context.TODO()
				album, err := db.GetAlbumById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				albumDir := app.WorkDir().Album(album.Id)

				files, err := pyrin.FormFiles(c, "cover")
				if err != nil {
					return nil, err
				}

				file := files[0]

				// TODO(patrik): Check content-type
				ext := path.Ext(file.Filename)
				filename := "cover-original" + ext

				dst := path.Join(albumDir, filename)

				dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					return nil, err
				}
				defer dstFile.Close()

				src, err := file.Open()
				if err != nil {
					return nil, err
				}
				defer src.Close()

				_, err = io.Copy(dstFile, src)
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

				err = tx.Commit()
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
