package apis

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type albumApi struct {
	app core.App
}

func (api *albumApi) HandleGetAlbums(c echo.Context) error {
	f := c.QueryParam("filter")

	albums, err := api.app.DB().GetAllAlbums(c.Request().Context(), f)
	if err != nil {
		return err
	}

	res := types.GetAlbums{
		Albums: make([]types.Album, len(albums)),
	}

	for i, album := range albums {
		res.Albums[i] = types.Album{
			Id:       album.Id,
			Name:     album.Name,
			CoverArt: utils.ConvertAlbumCoverURL(c, album.CoverArt),
			ArtistId: album.ArtistId,
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (api *albumApi) HandleGetAlbumById(c echo.Context) error {
	id := c.Param("id")
	album, err := api.app.DB().GetAlbumById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.GetAlbumById{
		Album: types.Album{
			Id:       album.Id,
			Name:     album.Name,
			CoverArt: utils.ConvertAlbumCoverURL(c, album.CoverArt),
			ArtistId: album.ArtistId,
		},
	}))
}

func (api *albumApi) HandleGetAlbumTracksById(c echo.Context) error {
	id := c.Param("id")

	album, err := api.app.DB().GetAlbumById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	tracks, err := api.app.DB().GetTracksByAlbum(c.Request().Context(), album.Id)
	if err != nil {
		return err
	}

	res := types.GetAlbumTracksById{
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

func InstallAlbumHandlers(app core.App, group Group) {
	api := albumApi{app: app}

	requireSetup := RequireSetup(app)

	group.Register(
		Handler{
			Name:        "GetAlbums",
			Path:        "/albums",
			Method:      http.MethodGet,
			DataType:    types.GetAlbums{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetAlbums,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},

		Handler{
			Name:        "GetAlbumById",
			Path:        "/albums/:id",
			Method:      http.MethodGet,
			DataType:    types.GetAlbumById{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetAlbumById,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},

		Handler{
			Name:        "GetAlbumTracks",
			Path:        "/albums/:id/tracks",
			Method:      http.MethodGet,
			DataType:    types.GetAlbumTracksById{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetAlbumTracksById,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},
	)
}
