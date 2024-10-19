package apis

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/assets"
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
	s := pyrin.NewServer(&pyrin.ServerConfig{
		LogName: dwebble.AppName,
		RegisterHandlers: func(router pyrin.Router) {
			g := router.Group("/api/v1")
			InstallHandlers(app, g)

			g = router.Group("/files")
			g.Register(
				pyrin.NormalHandler{
					Method:      http.MethodGet,
					Path:        "/images/default/:image",
					HandlerFunc: func(c pyrin.Context) error {
						image := c.Param("image")
						return fsFile(c.Response(), c.Request(), image, assets.DefaultImagesFS)
					},
				},
				pyrin.NormalHandler{
					Method:      http.MethodGet,
					Path:        "/albums/images/:albumId/:image",
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
					Path:        "/tracks/mobile/:albumId/:track",
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
					Path:        "/tracks/original/:albumId/:track",
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
