package apis

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/types"
)

type Handler struct {
	Name        string
	Method      string
	Path        string
	DataType    any
	BodyType    types.Body
	HandlerFunc echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

type Group interface {
	Register(handlers ...Handler)
}

func InstallHandlers(app core.App, g Group) {
	InstallArtistHandlers(app, g)
	InstallAlbumHandlers(app, g)
	InstallTrackHandlers(app, g)
	InstallSyncHandlers(app, g)
	InstallQueueHandlers(app, g)
	InstallTagHandlers(app, g)
	InstallAuthHandlers(app, g)
	InstallPlaylistHandlers(app, g)
	InstallSystemHandlers(app, g)
}
