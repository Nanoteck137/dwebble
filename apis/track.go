package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/helper"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

type Track struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	AlbumId  string `json:"albumId"`
	ArtistId string `json:"artistId"`

	Number   *int64 `json:"number"`
	Duration *int64 `json:"duration"`
	Year     *int64 `json:"year"`

	OriginalMediaUrl string       `json:"originalMediaUrl"`
	MobileMediaUrl   string       `json:"mobileMediaUrl"`
	CoverArt         types.Images `json:"coverArt"`

	AlbumName  string `json:"albumName"`
	ArtistName string `json:"artistName"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`

	Tags []string `json:"tags"`
}

func ConvertDBTrack(c pyrin.Context, track database.Track) Track {
	var number *int64
	if track.Number.Valid {
		number = &track.Number.Int64
	}

	var duration *int64
	if track.Duration.Valid {
		duration = &track.Duration.Int64
	}

	var year *int64
	if track.Year.Valid {
		year = &track.Year.Int64
	}

	return Track{
		Id:               track.Id,
		Name:             track.Name,
		AlbumId:          track.AlbumId,
		ArtistId:         track.ArtistId,
		Number:           number,
		Duration:         duration,
		Year:             year,
		OriginalMediaUrl: utils.ConvertURL(c, fmt.Sprintf("/files/tracks/%s/%s", track.Id, track.OriginalFilename)),
		MobileMediaUrl:   utils.ConvertURL(c, fmt.Sprintf("/files/tracks/%s/%s", track.Id, track.MobileFilename)),
		CoverArt:         utils.ConvertAlbumCoverURL(c, track.AlbumId, track.AlbumCoverArt),
		AlbumName:        track.AlbumName,
		ArtistName:       track.ArtistName,
		Created:          track.Created,
		Updated:          track.Updated,
		Tags:             utils.SplitString(track.Tags.String),
	}
}

type GetTracks struct {
	Page   types.Page `json:"page"`
	Tracks []Track    `json:"tracks"`
}

type GetTrackById struct {
	Track
}

// TODO(patrik): Validation and Transform
type EditTrackBody struct {
	Name       *string   `json:"name,omitempty"`
	ArtistId   *string   `json:"artistId,omitempty"`
	ArtistName *string   `json:"artistName,omitempty"`
	Year       *int64    `json:"year,omitempty"`
	Number     *int64    `json:"number,omitempty"`
	Tags       *[]string `json:"tags,omitempty"`
}

func InstallTrackHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "GetTracks",
			Method:   http.MethodGet,
			Path:     "/tracks",
			DataType: GetTracks{},
			Errors:   []pyrin.ErrorType{ErrTypeInvalidFilter, ErrTypeInvalidSort},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				perPage := 100
				page := 0

				// TODO(patrik): Move to util
				if s := q.Get("perPage"); s != "" {
					i, err := strconv.Atoi(s)
					if err != nil {
						return nil, err
					}

					perPage = i
				}

				if s := q.Get("page"); s != "" {
					i, err := strconv.Atoi(s)
					if err != nil {
						return nil, err
					}

					page = i
				}

				opts := database.FetchOption{
					Filter:  q.Get("filter"),
					Sort:    q.Get("sort"),
					PerPage: perPage,
					Page:    page,
				}

				tracks, p, err := app.DB().GetPagedTracks(c.Request().Context(), opts)
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					if errors.Is(err, database.ErrInvalidSort) {
						return nil, InvalidSort(err)
					}

					return nil, err
				}

				res := GetTracks{
					Page:   p,
					Tracks: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "SearchTracks",
			Method:   http.MethodGet,
			Path:     "/tracks/search",
			DataType: GetTracks{},
			Errors:   []pyrin.ErrorType{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				query := strings.TrimSpace(q.Get("query"))

				var err error
				var tracks []database.Track

				if query != "" {
					tracks, err = app.DB().SearchTracks(query)
					if err != nil {
						pretty.Println(err)
						return nil, err
					}
				}

				res := GetTracks{
					Tracks: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetTrackById",
			Method:   http.MethodGet,
			Path:     "/tracks/:id",
			DataType: GetTrackById{},
			Errors:   []pyrin.ErrorType{ErrTypeTrackNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				track, err := app.DB().GetTrackById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TrackNotFound()
					}

					return nil, err
				}

				return GetTrackById{
					Track: ConvertDBTrack(c, track),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "RemoveTrack",
			Method: http.MethodDelete,
			Path:   "/tracks/:id",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Move the track file to a trash can system
				id := c.Param("id")

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				err = db.RemoveTrack(c.Request().Context(), id)
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
			Name:     "EditTrack",
			Method:   http.MethodPatch,
			Path:     "/tracks/:id",
			BodyType: EditTrackBody{},
			Errors:   []pyrin.ErrorType{ErrTypeTrackNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				// TODO(patrik): Validation
				body, err := pyrin.Body[EditTrackBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				track, err := app.DB().GetTrackById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TrackNotFound()
					}

					return nil, err
				}

				var name types.Change[string]
				if body.Name != nil {
					n := strings.TrimSpace(*body.Name)
					name.Value = n
					name.Changed = n != track.Name
				}

				var artistId types.Change[string]
				if body.ArtistId != nil {
					artistId.Value = *body.ArtistId
					artistId.Changed = *body.ArtistId != track.ArtistId
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
					artistId.Changed = artist.Id != track.ArtistId
				}

				var year types.Change[sql.NullInt64]
				if body.Year != nil {
					year.Value = sql.NullInt64{
						Int64: *body.Year,
						Valid: *body.Year != 0,
					}
					year.Changed = *body.Year != track.Year.Int64
				}

				var number types.Change[sql.NullInt64]
				if body.Number != nil {
					number.Value = sql.NullInt64{
						Int64: *body.Number,
						Valid: *body.Number != 0,
					}
					number.Changed = *body.Number != track.Number.Int64
				}

				// TODO(patrik): Use transaction
				err = app.DB().UpdateTrack(ctx, track.Id, database.TrackChanges{
					Name:     name,
					ArtistId: artistId,
					Year:     year,
					Number:   number,
				})
				if err != nil {
					return nil, err
				}

				if body.Tags != nil {
					tags := *body.Tags

					err = app.DB().RemoveAllTagsFromTrack(ctx, track.Id)
					if err != nil {
						return nil, err
					}

					for _, tag := range tags {
						slug := utils.Slug(tag)

						err := app.DB().CreateTag(ctx, slug, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}

						err = app.DB().AddTagToTrack(ctx, slug, track.Id)
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
