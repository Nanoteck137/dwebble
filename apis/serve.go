package apis

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/pyrin"
)

// TODO(patrik): Add to pyrin
func fsFile(w http.ResponseWriter, r *http.Request, file string, filesystem fs.FS) error {
	f, err := filesystem.Open(file)
	if err != nil {
		// TODO(patrik): Add NoContentError to pyrin
		return echo.ErrNotFound
	}
	defer f.Close()

	fi, _ := f.Stat()

	ff, ok := f.(io.ReadSeeker)
	if !ok {
		return errors.New("file does not implement io.ReadSeeker")
	}

	http.ServeContent(w, r, fi.Name(), fi.ModTime(), ff)

	return nil
}

func Server(app core.App) (*pyrin.Server, error) {
	// TODO(patrik): Add to pyrin
	// e.RouteNotFound("/*", func(c echo.Context) error {
	// 	return RouteNotFound()
	// })
	//
	// TODO(patrik): Add to pyrin
	// e.Use(echolog.LoggerWithName(dwebble.AppName))
	// e.Use(middleware.Recover())
	// e.Use(middleware.CORS())
	//
	// TODO(patrik): Add to pyrin
	// e.StaticFS("/images/default", assets.DefaultImagesFS)

	s := pyrin.NewServer(&pyrin.ServerConfig{
		RegisterHandlers: func(router pyrin.Router) {
			g := router.Group("/api/v1")
			InstallHandlers(app, g)

			g = router.Group("/files")
			g.Register(
				pyrin.NormalHandler{
					Method:      http.MethodGet,
					Path:        "/albums/images/:albumId/:image",
					Middlewares: []echo.MiddlewareFunc{},
					HandlerFunc: func(c pyrin.Context) error {
						albumId := c.Param("albumId")
						image := c.Param("image")

						p := app.WorkDir().Album(albumId).Images()
						f := os.DirFS(p)

						return fsFile(c.Response(), c.Request(), image, f)
					},
				},
				pyrin.NormalHandler{
					Method:      http.MethodGet,
					Path:        "/files/tracks/mobile/:albumId/:track",
					Middlewares: []echo.MiddlewareFunc{},
					HandlerFunc: func(c pyrin.Context) error {
						albumId := c.Param("albumId")
						track := c.Param("track")

						p := app.WorkDir().Album(albumId).MobileFiles()
						f := os.DirFS(p)

						return fsFile(c.Response(), c.Request(), track, f)
					},
				},
				pyrin.NormalHandler{
					Method:      http.MethodGet,
					Path:        "/files/tracks/original/:albumId/:track",
					Middlewares: []echo.MiddlewareFunc{},
					HandlerFunc: func(c pyrin.Context) error {
						albumId := c.Param("albumId")
						track := c.Param("track")

						p := app.WorkDir().Album(albumId).OriginalFiles()
						f := os.DirFS(p)

						return fsFile(c.Response(), c.Request(), track, f)
					},
				},
			)
		},
	})

	app.DB().Invalidate()

	return s, nil
}
