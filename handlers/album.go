package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/types"
)

func (api *ApiConfig) HandleGetAlbums(c echo.Context) error {
	albums, err := api.db.GetAllAlbums(c.Request().Context())
	if err != nil {
		return err
	}

	res := types.ApiGetAlbumsData{
		Albums: make([]types.ApiAlbum, len(albums)),
	}

	for i, album := range albums {
		res.Albums[i] = types.ApiAlbum{
			Id:       album.Id,
			Name:     album.Name,
			CoverArt: album.CoverArt,
			ArtistId: album.ArtistId,
		}
	}

	return c.JSON(200, types.NewApiResponse(res))
}

func (api *ApiConfig) HandleGetAlbumById(c echo.Context) error {
	id := c.Param("id")
	album, err := api.db.GetAlbumById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiResponse(types.ApiGetAlbumByIdData{
		Id:       album.Id,
		Name:     album.Name,
		CoverArt: album.CoverArt,
		ArtistId: album.ArtistId,
	}))
}

func (api *ApiConfig) HandleGetAlbumTracksById(c echo.Context) error {
	id := c.Param("id")

	album, err := api.db.GetAlbumById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	tracks, err := api.db.GetTracksByAlbum(c.Request().Context(), album.Id)
	if err != nil {
		return err
	}

	res := types.ApiGetAlbumTracksByIdData{
		Tracks: make([]types.ApiTrack, len(tracks)),
	}

	for i, track := range tracks {
		res.Tracks[i] = types.ApiTrack{
			Id:                track.Id,
			Number:            int32(track.Number),
			Name:              track.Name,
			CoverArt:          track.CoverArt,
			BestQualityFile:   ConvertURL(c, "/tracks/original/"+track.BestQualityFile),
			MobileQualityFile: ConvertURL(c, "/tracks/mobile/"+track.MobileQualityFile),
			AlbumId:           track.AlbumId,
			ArtistId:          track.ArtistId,
		}
	}

	return c.JSON(200, types.NewApiResponse(res))
}

func InstallAlbumHandlers(group *echo.Group, apiConfig *ApiConfig) {
	group.GET("/albums", apiConfig.HandleGetAlbums)
	group.GET("/albums/:id", apiConfig.HandleGetAlbumById)
	group.GET("/albums/:id/tracks", apiConfig.HandleGetAlbumTracksById)
}
