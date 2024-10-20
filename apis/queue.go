package apis

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/pyrin"
)

type queueApi struct {
	app core.App
}

func (api *queueApi) HandlePostQueue(c echo.Context) error {
	// TODO(patrik): Add back
	return c.JSON(200, struct{}{})

	// tracks, err := api.app.DB().GetAllTracks(c.Request().Context(), "", "")
	// if err != nil {
	// 	return err
	// }
	//
	// res := types.PostQueue{
	// 	Tracks: make([]types.Track, len(tracks)),
	// }
	//
	// for i, track := range tracks {
	// 	res.Tracks[i] = types.Track{
	// 		Id:                track.Id,
	// 		Number:            track.Number,
	// 		Name:              track.Name,
	// 		CoverArt:          utils.ConvertTrackCoverURL(c, track.CoverArt),
	// 		Duration:          track.Duration,
	// 		BestQualityFile:   utils.ConvertURL(c, "/tracks/original/"+track.BestQualityFile),
	// 		MobileQualityFile: utils.ConvertURL(c, "/tracks/mobile/"+track.MobileQualityFile),
	// 		AlbumId:           track.AlbumId,
	// 		ArtistId:          track.ArtistId,
	// 		AlbumName:         track.AlbumName,
	// 		ArtistName:        track.ArtistName,
	// 	}
	// }
	//
	// return c.JSON(200, SuccessResponse(res))
}

// TODO(patrik): Add back
func InstallQueueHandlers(app core.App, group pyrin.Group) {
	// api := queueApi{app: app}
	//
	// requireSetup := RequireSetup(app)
	//
	// group.Register(
	// 	Handler{
	// 		Name:        "CreateQueue",
	// 		Path:        "/queue",
	// 		Method:      http.MethodPost,
	// 		DataType:    types.PostQueue{},
	// 		BodyType:    nil,
	// 		HandlerFunc: api.HandlePostQueue,
	// 		Middlewares: []echo.MiddlewareFunc{requireSetup},
	// 	},
	// )
}
