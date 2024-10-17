package apis

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/api"
)

type Handler struct {
	Name        string
	Method      string
	Path        string
	DataType    any
	BodyType    types.Body
	IsMultiForm bool
	Errors      []api.ErrorType
	HandlerFunc echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

type Group interface {
	Register(handlers ...Handler)
}

func InstallHandlers(app core.App, g pyrin.Group) {
	InstallArtistHandlers(app, g)
	InstallAlbumHandlers(app, g)
	InstallTrackHandlers(app, g)
	// InstallQueueHandlers(app, g)
	// InstallTagHandlers(app, g)
	InstallAuthHandlers(app, g)
	InstallPlaylistHandlers(app, g)
	// InstallSystemHandlers(app, g)
}
