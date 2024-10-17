package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	pyrinapi "github.com/nanoteck137/pyrin/api"
	"github.com/nanoteck137/pyrin/tools/validate"
)

type trackApi struct {
	app core.App
}

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

func (api *trackApi) HandleGetTracks(c pyrin.Context) (any, error) {
	// TODO(patrik): Fix
	// f := c.QueryParam("filter")
	// s := c.QueryParam("sort")
	// includeAll := ParseQueryBool(c.QueryParam("includeAll"))

	f := ""
	s := ""
	includeAll := false

	// TODO(patrik): Fix
	_ = includeAll

	tracks, err := api.app.DB().GetAllTracks(c.Request().Context(), f, s, includeAll)
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

	// TODO(patrik): Fix
	// for i, track := range tracks {
	// 	res.Tracks[i] = ConvertDBTrack(c, track)
	// }

	return res, nil
}

func (api *trackApi) HandleGetTrackById(c pyrin.Context) (any, error) {
	// id := c.Param("id")

	// track, err := api.app.DB().GetTrackById(c.Request().Context(), id)
	// if err != nil {
	// 	if errors.Is(err, database.ErrItemNotFound) {
	// 		return TrackNotFound()
	// 	}
	//
	// 	return err
	// }

	return types.GetTrackById{
		// TODO(patrik): Fix
		// Track: ConvertDBTrack(c, track),
	}, nil
}

// TODO(patrik): Move the track file to a trash can system
func (api *trackApi) HandleDeleteTrack(c pyrin.Context) (any, error) {
	id := c.Param("id")

	db, tx, err := api.app.DB().Begin()
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
func (*PatchTrackBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func (api *trackApi) HandlePatchTrack(c echo.Context) error {
	id := c.Param("id")

	var body PatchTrackBody
	d := json.NewDecoder(c.Request().Body)
	err := d.Decode(&body)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	track, err := api.app.DB().GetTrackById(ctx, id)
	if err != nil {
		// TODO(patrik): Handle error
		return err
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

		artist, err := api.app.DB().GetArtistByName(ctx, artistName)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				artist, err = api.app.DB().CreateArtist(ctx, database.CreateArtistParams{
					Name:    artistName,
					Picture: sql.NullString{},
				})

				if err != nil {
					return err
				}
			} else {
				return err
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
	err = api.app.DB().UpdateTrack(ctx, track.Id, database.TrackChanges{
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
		return err
	}

	if body.Tags != nil {
		tags := *body.Tags

		err = api.app.DB().RemoveAllTagsFromTrack(ctx, track.Id)
		if err != nil {
			return err
		}

		for _, tag := range tags {
			slug := utils.Slug(tag)

			err := api.app.DB().CreateTag(ctx, slug, tag)
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}

			err = api.app.DB().AddTagToTrack(ctx, slug, track.Id)
			if err != nil {
				return err
			}
		}
	}

	return c.JSON(200, pyrinapi.SuccessResponse(nil))
}

func InstallTrackHandlers(app core.App, group pyrin.Group) {
	a := trackApi{app: app}

	group.Register(
		pyrin.ApiHandler{
			Name:        "GetTracks",
			Method:      http.MethodGet,
			Path:        "/tracks",
			DataType:    types.GetTracks{},
			Errors:      []ErrorType{ErrTypeInvalidFilter, ErrTypeInvalidSort},
			HandlerFunc: a.HandleGetTracks,
		},

		pyrin.ApiHandler{
			Name:        "GetTrackById",
			Method:      http.MethodGet,
			Path:        "/tracks/:id",
			DataType:    types.GetTrackById{},
			Errors:      []ErrorType{ErrTypeTrackNotFound},
			HandlerFunc: a.HandleGetTrackById,
		},

		pyrin.ApiHandler{
			Name:        "RemoveTrack",
			Method:      http.MethodDelete,
			Path:        "/tracks/:id",
			HandlerFunc: a.HandleDeleteTrack,
		},

		// pyrin.ApiHandler{
		// 	Name:        "EditTrack",
		// 	Method:      http.MethodPatch,
		// 	Path:        "/tracks/:id",
		// 	BodyType:    PatchTrackBody{},
		// 	HandlerFunc: a.HandlePatchTrack,
		// },
	)
}
