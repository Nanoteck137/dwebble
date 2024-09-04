package apis

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/MadAppGang/httplog/echolog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/assets"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/pyrin/api"
)

type echoGroup struct {
	app core.App

	Prefix string
	Group  *echo.Group
}

func (g *echoGroup) Register(handlers ...Handler) {
	for _, h := range handlers {
		log.Debug("Registering", "method", h.Method, "name", h.Name, "path", g.Prefix+h.Path)
		g.Group.Add(h.Method, h.Path, h.HandlerFunc, h.Middlewares...)
	}
}

func newEchoGroup(app core.App, e *echo.Echo, prefix string, m ...echo.MiddlewareFunc) *echoGroup {
	g := e.Group(prefix, m...)

	return &echoGroup{
		app:    app,
		Prefix: prefix,
		Group:  g,
	}
}

const ErrTypeUnknownError api.ErrorType = "UNKNOWN_ERROR"

func errorHandler(err error, c echo.Context) {
	switch err := err.(type) {
	case *api.Error:
		c.JSON(err.Code, api.Response{
			Success: false,
			Error:   err,
		})
	case *echo.HTTPError:
		c.JSON(err.Code, api.Response{
			Success: false,
			Error: &api.Error{
				Code:    err.Code,
				Type:    ErrTypeUnknownError,
				Message: err.Error(),
			},
		})
	default:
		c.JSON(500, api.Response{
			Success: false,
			Error: &api.Error{
				Code:    500,
				Type:    ErrTypeUnknownError,
				Message: "Internal Server Error",
			},
		})
	}

	log.Error("HTTP API Error", "err", err)
}

func fsFile(c echo.Context, file string, filesystem fs.FS) error {
	f, err := filesystem.Open(file)
	if err != nil {
		return echo.ErrNotFound
	}
	defer f.Close()

	fi, _ := f.Stat()

	ff, ok := f.(io.ReadSeeker)
	if !ok {
		return errors.New("file does not implement io.ReadSeeker")
	}
	http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), ff)

	return nil
}

func Server(app core.App) (*echo.Echo, error) {
	e := echo.New()

	e.RouteNotFound("/*", func(c echo.Context) error {
		return RouteNotFound()
	})

	e.HTTPErrorHandler = errorHandler

	e.Use(echolog.LoggerWithName(dwebble.AppName))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Add(http.MethodGet, "/files/tracks/original/:albumId/:track", func(c echo.Context) error {
		albumId := c.Param("albumId")
		track := c.Param("track")

		p := app.WorkDir().Album(albumId).OriginalFiles()
		f := os.DirFS(p)

		return fsFile(c, track, f)
	})

	e.Add(http.MethodGet, "/files/tracks/mobile/:albumId/:track", func(c echo.Context) error {
		albumId := c.Param("albumId")
		track := c.Param("track")

		p := app.WorkDir().Album(albumId).MobileFiles()
		f := os.DirFS(p)

		return fsFile(c, track, f)
	})

	e.Add(http.MethodGet, "/files/albums/images/:albumId/:image", func(c echo.Context) error {
		albumId := c.Param("albumId")
		image := c.Param("image")

		p := app.WorkDir().Album(albumId).Images()
		f := os.DirFS(p)

		return fsFile(c, image, f)
	})

	e.StaticFS("/images/default", assets.DefaultImagesFS)
	// e.Static("/images", app.WorkDir().ImagesDir())

	g := newEchoGroup(app, e, "/api/v1")
	InstallHandlers(app, g)

	app.DB().Invalidate()

	return e, nil
}
