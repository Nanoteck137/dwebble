package apis

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type trackApi struct {
	app core.App
}

func (api *trackApi) HandleGetTracks(c echo.Context) error {
	f := c.QueryParam("filter")
	s := c.QueryParam("sort")

	tracks, err := api.app.DB().GetAllTracks(c.Request().Context(), f, s)
	if err != nil {
		return err
	}

	res := types.GetTracks{
		Tracks: make([]types.Track, len(tracks)),
	}

	for i, track := range tracks {
		res.Tracks[i] = types.Track{
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
			Tags:              strings.Split(track.Tags.String, ","),
			Genres:            strings.Split(track.Genres.String, ","),
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (api *trackApi) HandleGetTrackById(c echo.Context) error {
	id := c.Param("id")
	track, err := api.app.DB().GetTrackById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.GetTrackById{
		Track: types.Track{
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
			Tags:              strings.Split(track.Tags.String, ","),
			Genres:            strings.Split(track.Genres.String, ","),
		},
	}))
}

func InstallTrackHandlers(app core.App, group Group) {
	api := trackApi{app: app}

	requireSetup := RequireSetup(app)

	group.Register(
		Handler{
			Name:        "GetTracks",
			Path:        "/tracks",
			Method:      http.MethodGet,
			DataType:    types.GetTracks{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetTracks,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},

		Handler{
			Name:        "GetTrackById",
			Path:        "/tracks/:id",
			Method:      http.MethodGet,
			DataType:    types.GetTrackById{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetTrackById,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},
	)
}
