package apis

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/api"
)

type albumApi struct {
	app core.App
}

func (api *albumApi) HandleGetAlbums(c echo.Context) error {
	filter := c.QueryParam("filter")
	sort := c.QueryParam("sort")
	includeAll := ParseQueryBool(c.QueryParam("includeAll"))

	albums, err := api.app.DB().GetAllAlbums(c.Request().Context(), filter, sort, includeAll)
	if err != nil {
		return err
	}

	res := types.GetAlbums{
		Albums: make([]types.Album, len(albums)),
	}

	for i, album := range albums {
		res.Albums[i] = types.Album{
			Id:         album.Id,
			Name:       album.Name,
			CoverArt:   utils.ConvertAlbumCoverURL(c, album.CoverArt),
			ArtistId:   album.ArtistId,
			ArtistName: album.ArtistName,
		}
	}

	return c.JSON(200, SuccessResponse(res))
}

func (api *albumApi) HandleGetAlbumById(c echo.Context) error {
	id := c.Param("id")
	album, err := api.app.DB().GetAlbumById(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return AlbumNotFound()
		}

		return err
	}

	return c.JSON(200, SuccessResponse(types.GetAlbumById{
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
		if errors.Is(err, database.ErrItemNotFound) {
			return AlbumNotFound()
		}

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
		res.Tracks[i] = ConvertDBTrack(c, track) 
	}

	return c.JSON(200, SuccessResponse(res))
}

func InstallAlbumHandlers(app core.App, group Group) {
	a := albumApi{app: app}

	group.Register(
		Handler{
			Name:        "GetAlbums",
			Path:        "/albums",
			Method:      http.MethodGet,
			DataType:    types.GetAlbums{},
			BodyType:    nil,
			HandlerFunc: a.HandleGetAlbums,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "GetAlbumById",
			Method:      http.MethodGet,
			Path:        "/albums/:id",
			DataType:    types.GetAlbumById{},
			BodyType:    nil,
			Errors:      []api.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: a.HandleGetAlbumById,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "GetAlbumTracks",
			Method:      http.MethodGet,
			Path:        "/albums/:id/tracks",
			DataType:    types.GetAlbumTracksById{},
			BodyType:    nil,
			Errors:      []api.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: a.HandleGetAlbumTracksById,
			Middlewares: []echo.MiddlewareFunc{},
		},
	)
}
