package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/validate"
)

func ConvertDBTrack(c pyrin.Context, track database.Track) types.Track {
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

	return types.Track{
		Id:               track.Id,
		Name:             track.Name,
		AlbumId:          track.AlbumId,
		ArtistId:         track.ArtistId,
		Number:           number,
		Duration:         duration,
		Year:             year,
		OriginalMediaUrl: utils.ConvertURL(c, fmt.Sprintf("/files/tracks/original/%s/%s", track.AlbumId, track.OriginalFilename)),
		MobileMediaUrl:   utils.ConvertURL(c, fmt.Sprintf("/files/tracks/mobile/%s/%s", track.AlbumId, track.MobileFilename)),
		CoverArt:         utils.ConvertAlbumCoverURL(c, track.AlbumId, track.AlbumCoverArt),
		AlbumName:        track.AlbumName,
		ArtistName:       track.ArtistName,
		Created:          track.Created,
		Updated:          track.Updated,
		Tags:             utils.SplitString(track.Tags.String),
		Available:        track.Available,
	}
}

func ParseQueryBool(s string) bool {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	switch s {
	case "true", "1":
		return true
	// TODO(patrik): Add "0"
	case "false", "":
		return false
	default:
		return false
	}
}

var _ pyrin.Body = (*PatchTrackBody)(nil)

type PatchTrackBody struct {
	Name       *string   `json:"name,omitempty"`
	ArtistId   *string   `json:"artistId,omitempty"`
	ArtistName *string   `json:"artistName,omitempty"`
	Year       *int64    `json:"year,omitempty"`
	Number     *int64    `json:"number,omitempty"`
	Tags       *[]string `json:"tags,omitempty"`
}

// Validate implements pyrin.Body.
func (b PatchTrackBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func InstallTrackHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "GetTracks",
			Method:   http.MethodGet,
			Path:     "/tracks",
			DataType: types.GetTracks{},
			Errors:   []pyrin.ErrorType{ErrTypeInvalidFilter, ErrTypeInvalidSort},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				f := q.Get("filter")
				s := q.Get("sort")
				includeAll := ParseQueryBool(q.Get("includeAll"))

				tracks, err := app.DB().GetAllTracks(c.Request().Context(), f, s, includeAll)
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					if errors.Is(err, database.ErrInvalidSort) {
						return nil, InvalidSort(err)
					}

					return nil, err
				}

				res := types.GetTracks{
					Tracks: make([]types.Track, len(tracks)),
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
			DataType: types.GetTrackById{},
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

				return types.GetTrackById{
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
			BodyType: PatchTrackBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				var body PatchTrackBody
				d := json.NewDecoder(c.Request().Body)
				err := d.Decode(&body)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				track, err := app.DB().GetTrackById(ctx, id)
				if err != nil {
					// TODO(patrik): Handle error
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
							artist, err = app.DB().CreateArtist(ctx, database.CreateArtistParams{
								Name:    artistName,
								Picture: sql.NullString{},
							})

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
					Available: types.Change[bool]{
						Value:   true,
						Changed: true,
					},
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
