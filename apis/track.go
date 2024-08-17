package apis

import (
	"errors"
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
	return types.Track{
		Id:                track.Id,
		Number:            track.Number,
		Name:              track.Name,
		CoverArt:          utils.ConvertTrackCoverURL(c, track.CoverArt),
		Duration:          track.Duration,
		BestQualityFile:   utils.ConvertURL(c, "/tracks/original/"+track.BestQualityFile),
		MobileQualityFile: utils.ConvertURL(c, "/tracks/mobile/"+track.MobileQualityFile),
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
	)
}
