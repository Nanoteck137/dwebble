package apis

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/api"
)

type trackApi struct {
	app core.App
}

func ConvertDBTrack(c echo.Context, track database.Track) types.Track {
	// /files/tracks/original/:albumId/:track
	// TODO(patrik): Change types.Track to match the new database.Track
	return types.Track{
		Id:                track.Id,
		Name:              track.Name,
		Number:            int(track.Number.Int64),
		CoverArt:          utils.ConvertAlbumCoverURL(c, track.AlbumId, track.AlbumCoverArt),
		Duration:          int(track.Duration.Int64),
		BestQualityFile:   utils.ConvertURL(c, fmt.Sprintf("/files/tracks/original/%s/%s", track.AlbumId, track.OriginalFilename)),
		MobileQualityFile: utils.ConvertURL(c, fmt.Sprintf("/files/tracks/mobile/%s/%s", track.AlbumId, track.MobileFilename)),
		AlbumId:           track.AlbumId,
		ArtistId:          track.ArtistId,
		AlbumName:         track.AlbumName,
		ArtistName:        track.ArtistName,
		Available:         track.Available,
		Tags:              utils.SplitString(track.Tags.String),
		Genres:            utils.SplitString(track.Genres.String),
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

func (api *trackApi) HandleGetTracks(c echo.Context) error {
	f := c.QueryParam("filter")
	s := c.QueryParam("sort")
	includeAll := ParseQueryBool(c.QueryParam("includeAll"))

	_ = includeAll

	tracks, err := api.app.DB().GetAllTracks(c.Request().Context(), f, s, includeAll)
	if err != nil {
		if errors.Is(err, database.ErrInvalidFilter) {
			return InvalidFilter(err)
		}

		if errors.Is(err, database.ErrInvalidSort) {
			return InvalidSort(err)
		}

		return err
	}

	res := types.GetTracks{
		Tracks: make([]types.Track, len(tracks)),
	}

	for i, track := range tracks {
		res.Tracks[i] = ConvertDBTrack(c, track)
	}

	return c.JSON(200, SuccessResponse(res))
}

func (api *trackApi) HandleGetTrackById(c echo.Context) error {
	id := c.Param("id")

	track, err := api.app.DB().GetTrackById(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return TrackNotFound()
		}

		return err
	}

	return c.JSON(200, SuccessResponse(types.GetTrackById{
		Track: ConvertDBTrack(c, track),
	}))
}

// TODO(patrik): Move the track file to a trash can system
func (api *trackApi) HandleDeleteTrack(c echo.Context) error {
	id := c.Param("id")

	db, tx, err := api.app.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = db.RemoveTrack(c.Request().Context(), id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return c.JSON(200, SuccessResponse(nil))
}


func InstallTrackHandlers(app core.App, group Group) {
	a := trackApi{app: app}

	group.Register(
		Handler{
			Name:        "GetTracks",
			Method:      http.MethodGet,
			Path:        "/tracks",
			DataType:    types.GetTracks{},
			BodyType:    nil,
			Errors:      []api.ErrorType{ErrTypeInvalidFilter, ErrTypeInvalidSort},
			HandlerFunc: a.HandleGetTracks,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "GetTrackById",
			Method:      http.MethodGet,
			Path:        "/tracks/:id",
			DataType:    types.GetTrackById{},
			BodyType:    nil,
			Errors:      []api.ErrorType{ErrTypeTrackNotFound},
			HandlerFunc: a.HandleGetTrackById,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "RemoveTrack",
			Method:      http.MethodDelete,
			Path:        "/tracks/:id",
			DataType:    nil,
			BodyType:    nil,
			HandlerFunc: a.HandleDeleteTrack,
		},
	)
}
