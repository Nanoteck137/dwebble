package apis

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

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

	featuringArtists := ConvertDBFeaturingArtists(album.FeaturingArtists)
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
		CoverArt: ConvertAlbumCoverURL(c, album.Id, album.CoverArt),
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
	Page   types.Page `json:"page"`
	Albums []Album    `json:"albums"`
}

type GetAlbumById struct {
	Album
}

type GetAlbumTracks struct {
	Tracks []Track `json:"tracks"`
}

type GetAlbumTracksDetails struct {
	Tracks []TrackDetails `json:"tracks"`
}

type EditAlbumBody struct {
	Name             *string   `json:"name,omitempty"`
	OtherName        *string   `json:"otherName,omitempty"`
	ArtistId         *string   `json:"artistId,omitempty"`
	Year             *int64    `json:"year,omitempty"`
	Tags             *[]string `json:"tags,omitempty"`
	FeaturingArtists *[]string `json:"featuringArtists,omitempty"`
}

func (b *EditAlbumBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
	b.OtherName = transform.StringPtr(b.OtherName)
	b.Tags = DiscardEntriesStringArrayPtr(b.Tags)
	b.FeaturingArtists = DiscardEntriesStringArrayPtr(b.FeaturingArtists)
}

func (b EditAlbumBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
		// validate.Field(&b.OtherName, validate.Required.When(b.OtherName != nil)),
		validate.Field(&b.ArtistId, validate.Required.When(b.ArtistId != nil)),
		validate.Field(&b.Year, validate.Min(0)),
	)
}

type CreateAlbum struct {
	AlbumId string `json:"albumId"`
}

type CreateAlbumBody struct {
	Name      string `json:"name"`
	OtherName string `json:"otherName"`

	ArtistId string `json:"artistId"`

	Year int64 `json:"year"`

	Tags             []string `json:"tags"`
	FeaturingArtists []string `json:"featuringArtists"`
}

func (b *CreateAlbumBody) Transform() {
	b.Name = transform.String(b.Name)
	b.OtherName = transform.String(b.OtherName)
	b.ArtistId = transform.String(b.ArtistId)

	b.Tags = StringArray(b.Tags)
	b.Tags = DiscardEntriesStringArray(b.Tags)
	b.FeaturingArtists = StringArray(b.FeaturingArtists)
	b.FeaturingArtists = DiscardEntriesStringArray(b.FeaturingArtists)
}

func (b CreateAlbumBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
		validate.Field(&b.ArtistId, validate.Required),
	)
}

type GetAlbumTracksForPlay struct {
	Tracks []MusicTrack `json:"tracks"`
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
				opts := getPageOptions(q)

				albums, pageInfo, err := app.DB().GetAlbumsPaged(c.Request().Context(), opts)
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
					Page:   pageInfo,
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

		pyrin.ApiHandler{
			Name:         "GetAlbumTracksDetails",
			Method:       http.MethodGet,
			Path:         "/albums/:id/tracks/details",
			ResponseType: GetAlbumTracksDetails{},
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

				res := GetAlbumTracksDetails{
					Tracks: make([]TrackDetails, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrackToDetails(c, track)
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

				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

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

				{
					album, err := app.DB().GetAlbumById(ctx, album.Id)
					if err != nil {
						return nil, err
					}

					err = app.DB().UpdateSearchAlbum(ctx, album)
					if err != nil {
						return nil, err
					}
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

				_, err := User(app, c, RequireAdmin)
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

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				tracks, err := db.GetTracksByAlbum(ctx, album.Id)
				if err != nil {
					return nil, err
				}

				for _, track := range tracks {
					err = DeleteTrack(ctx, db, app.WorkDir(), track)
					if err != nil {
						return nil, err
					}
				}

				err = DeleteAlbum(ctx, db, app.WorkDir(), album)
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
				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

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

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				album, err := helper.CreateAlbum(ctx, db, app.WorkDir(), database.CreateAlbumParams{
					Name: body.Name,
					OtherName: sql.NullString{
						String: body.OtherName,
						Valid:  body.OtherName != "",
					},
					ArtistId: artist.Id,
					Year: sql.NullInt64{
						Int64: body.Year,
						Valid: body.Year != 0,
					},
				})
				if err != nil {
					return nil, err
				}

				for _, tag := range body.Tags {
					// TODO(patrik): Move slugify to transform function
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

				for _, artistId := range body.FeaturingArtists {
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
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				err = tx.Commit()
				if err != nil {
					return nil, err
				}

				{
					album, err := app.DB().GetAlbumById(ctx, album.Id)
					if err != nil {
						return nil, err
					}

					err = app.DB().InsertAlbumToSearch(ctx, album)
					if err != nil {
						return nil, err
					}
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

				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

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
				err = utils.CreateResizedImage(dst, i, 128, 128)
				if err != nil {
					return nil, err
				}

				i = path.Join(albumDir, "cover-256.png")
				err = utils.CreateResizedImage(dst, i, 256, 256)
				if err != nil {
					return nil, err
				}

				i = path.Join(albumDir, "cover-512.png")
				err = utils.CreateResizedImage(dst, i, 512, 512)
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
