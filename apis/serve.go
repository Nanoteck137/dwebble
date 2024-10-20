package apis

import (
	"net/http"
	"os"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/assets"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/pyrin"
)

func RegisterHandlers(app core.App, router pyrin.Router) {
	g := router.Group("/api/v1")
	InstallHandlers(app, g)

	g = router.Group("/files")
	g.Register(
		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/images/default/:image",
			HandlerFunc: func(c pyrin.Context) error {
				image := c.Param("image")
				return pyrin.ServeFile(c.Response(), c.Request(), assets.DefaultImagesFS, image)
			},
		},
		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/albums/images/:albumId/:image",
			HandlerFunc: func(c pyrin.Context) error {
				albumId := c.Param("albumId")
				image := c.Param("image")

				p := app.WorkDir().Album(albumId).Images()
				f := os.DirFS(p)

				return pyrin.ServeFile(c.Response(), c.Request(), f, image)
			},
		},
		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/tracks/mobile/:albumId/:track",
			HandlerFunc: func(c pyrin.Context) error {
				albumId := c.Param("albumId")
				track := c.Param("track")

				p := app.WorkDir().Album(albumId).MobileFiles()
				f := os.DirFS(p)

				return pyrin.ServeFile(c.Response(), c.Request(), f, track)
			},
		},
		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/tracks/original/:albumId/:track",
			HandlerFunc: func(c pyrin.Context) error {
				albumId := c.Param("albumId")
				track := c.Param("track")

				p := app.WorkDir().Album(albumId).OriginalFiles()
				f := os.DirFS(p)

				return pyrin.ServeFile(c.Response(), c.Request(), f, track)
			},
		},
	)
}

func Server(app core.App) (*pyrin.Server, error) {
	s := pyrin.NewServer(&pyrin.ServerConfig{
		LogName: dwebble.AppName,
		RegisterHandlers: func(router pyrin.Router) {
			RegisterHandlers(app, router)
		},
	})

	app.DB().Invalidate()

	return s, nil
}
