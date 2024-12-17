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
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
)

type Name struct {
	Name      string  `json:"name"`
	OtherName *string `json:"otherName"`
}

type Track struct {
	Id   string `json:"id"`
	Name Name   `json:"name"`

	Number   *int64 `json:"number"`
	Duration *int64 `json:"duration"`
	Year     *int64 `json:"year"`

	OriginalMediaUrl string       `json:"originalMediaUrl"`
	MobileMediaUrl   string       `json:"mobileMediaUrl"`
	CoverArt         types.Images `json:"coverArt"`

	AlbumId  string `json:"albumId"`
	ArtistId string `json:"artistId"`

	AlbumName  Name `json:"albumName"`
	ArtistName Name `json:"artistName"`

	Tags         []string      `json:"tags"`
	ExtraArtists []ExtraArtist `json:"extraArtists"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

func ConvertDBTrack(c pyrin.Context, track database.Track) Track {
	return Track{
		Id: track.Id,
		Name: Name{
			Name:      track.Name,
			OtherName: ConvertSqlNullString(track.OtherName),
		},
		Number:           ConvertSqlNullInt64(track.Number),
		Duration:         ConvertSqlNullInt64(track.Duration),
		Year:             ConvertSqlNullInt64(track.Year),
		OriginalMediaUrl: utils.ConvertURL(c, fmt.Sprintf("/files/tracks/%s/%s", track.Id, track.OriginalFilename)),
		MobileMediaUrl:   utils.ConvertURL(c, fmt.Sprintf("/files/tracks/%s/%s", track.Id, track.MobileFilename)),
		CoverArt:         utils.ConvertAlbumCoverURL(c, track.AlbumId, track.AlbumCoverArt),
		AlbumId:          track.AlbumId,
		ArtistId:         track.ArtistId,
		AlbumName: Name{
			Name:      track.AlbumName,
			OtherName: ConvertSqlNullString(track.AlbumOtherName),
		},
		ArtistName: Name{
			Name:      track.ArtistName,
			OtherName: ConvertSqlNullString(track.ArtistOtherName),
		},
		Tags:         utils.SplitString(track.Tags.String),
		ExtraArtists: ConvertDBExtraArtists(track.ExtraArtists),
		Created:      track.Created,
		Updated:      track.Updated,
	}
}

type GetTracks struct {
	Page   types.Page `json:"page"`
	Tracks []Track    `json:"tracks"`
}

type GetTrackById struct {
	Track
}

type EditTrackBody struct {
	Name      *string `json:"name,omitempty"`
	OtherName *string `json:"otherName,omitempty"`
	ArtistId  *string `json:"artistId,omitempty"`
	// TODO(patrik): Remove artistName
	ArtistName   *string   `json:"artistName,omitempty"`
	Year         *int64    `json:"year,omitempty"`
	Number       *int64    `json:"number,omitempty"`
	Tags         *[]string `json:"tags,omitempty"`
	ExtraArtists *[]string `json:"extraArtists,omitempty"`
}

// TODO(patrik): Should this be pyrins default
func DiscardEmptyStringEntries(arr *[]string) *[]string {
	if arr == nil {
		return nil
	}

	var res []string

	v := *arr
	for _, s := range v {
		if s != "" {
			res = append(res, s)
		}
	}

	return &res
}

func (b *EditTrackBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
	b.OtherName = transform.StringPtr(b.OtherName)
	b.Tags = transform.StringArrayPtr(b.Tags)
	b.Tags = DiscardEmptyStringEntries(b.Tags)
	b.ExtraArtists = DiscardEmptyStringEntries(b.ExtraArtists)
}

func (b EditTrackBody) Validate() error {
	checkBothArtist := validate.By(func(value interface{}) error {
		if b.ArtistId != nil && b.ArtistName != nil {
			return errors.New("both 'artistId' and 'artistName' can't be specified remove one of them")
		}

		return nil
	})

	return validate.ValidateStruct(&b,
		validate.Field(&b.Name,
			validate.Required.When(b.Name != nil),
		),
		validate.Field(&b.ArtistId,
			validate.Required.When(b.ArtistId != nil),
			checkBothArtist,
		),
		validate.Field(&b.ArtistName,
			validate.Required.When(b.ArtistName != nil),
			checkBothArtist,
		),
		validate.Field(&b.Year, validate.Min(0)),
		validate.Field(&b.Number, validate.Min(0)),
	)
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

				tracks, err := app.DB().SearchTracks(query)
				if err != nil {
					return nil, err
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

				pretty.Println(track)

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
			Errors:   []pyrin.ErrorType{ErrTypeTrackNotFound, ErrTypeArtistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

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

				changes := database.TrackChanges{}

				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != track.Name,
					}
				}

				if body.OtherName != nil {
					changes.OtherName = types.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.OtherName,
							Valid:  *body.OtherName != "",
						},
						Changed: *body.OtherName != track.OtherName.String,
					}
				}

				if body.ArtistId != nil {
					// TODO(patrik): Should we silently continue if
					// we don't find the artist
					_, err := app.DB().GetArtistById(ctx, *body.ArtistId)
					if err != nil {
						if errors.Is(err, database.ErrItemNotFound) {
							return nil, ArtistNotFound()
						}

						return nil, err
					}

					changes.ArtistId = types.Change[string]{
						Value:   *body.ArtistId,
						Changed: *body.ArtistId != track.ArtistId,
					}
				}

				if body.Year != nil {
					changes.Year = types.Change[sql.NullInt64]{
						Value: sql.NullInt64{
							Int64: *body.Year,
							Valid: *body.Year != 0,
						},
						Changed: *body.Year != track.Year.Int64,
					}
				}

				if body.Number != nil {
					changes.Number = types.Change[sql.NullInt64]{
						Value: sql.NullInt64{
							Int64: *body.Number,
							Valid: *body.Number != 0,
						},
						Changed: *body.Number != track.Number.Int64,
					}
				}

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				err = db.UpdateTrack(ctx, track.Id, changes)
				if err != nil {
					return nil, err
				}

				if body.Tags != nil {
					tags := *body.Tags

					err = db.RemoveAllTagsFromTrack(ctx, track.Id)
					if err != nil {
						return nil, err
					}

					for _, tag := range tags {
						slug := utils.Slug(tag)

						err := db.CreateTag(ctx, slug, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}

						err = db.AddTagToTrack(ctx, slug, track.Id)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}
				}

				if body.ExtraArtists != nil {
					extraArtists := *body.ExtraArtists

					err := db.RemoveAllExtraArtistsFromTrack(ctx, track.Id)
					if err != nil {
						return nil, err
					}

					for _, artistId := range extraArtists {
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

						err = db.AddExtraArtistToTrack(ctx, track.Id, artist.Id)
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
			Name:   "DeleteTrack",
			Method: http.MethodDelete,
			Path:   "/tracks/:id",
			Errors: []pyrin.ErrorType{ErrTypeTrackNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				track, err := app.DB().GetTrackById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TrackNotFound()
					}

					return nil, err
				}

				err = app.DB().RemoveTrack(ctx, track.Id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
