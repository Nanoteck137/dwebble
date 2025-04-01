package apis

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/assets"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
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
				return pyrin.ServeFile(c, assets.DefaultImagesFS, image)
			},
		},
		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/albums/images/:albumId/:image",
			HandlerFunc: func(c pyrin.Context) error {
				albumId := c.Param("albumId")
				image := c.Param("image")

				p := app.WorkDir().Album(albumId)
				f := os.DirFS(p)

				return pyrin.ServeFile(c, f, image)
			},
		},
		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/artists/:artistId/:file",
			HandlerFunc: func(c pyrin.Context) error {
				artistId := c.Param("artistId")
				file := c.Param("file")

				p := app.WorkDir().Artist(artistId)
				f := os.DirFS(p)

				return pyrin.ServeFile(c, f, file)
			},
		},
		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/tracks/:trackId/:file",
			HandlerFunc: func(c pyrin.Context) error {
				trackId := c.Param("trackId")
				file := c.Param("file")

				trackDir := app.WorkDir().Track(trackId)
				f := os.DirFS(trackDir)

				return pyrin.ServeFile(c, f, file)
			},
		},
		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/tracks/:trackId/:file",
			HandlerFunc: func(c pyrin.Context) error {
				trackId := c.Param("trackId")
				file := c.Param("file")

				ctx := context.TODO()

				track, err := app.DB().GetTrackById(ctx, trackId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return pyrin.NoContentNotFound()
					}

					return err
				}

				if file == track.Filename {
					trackDir := app.WorkDir().Track(trackId)
					f := os.DirFS(trackDir)

					return pyrin.ServeFile(c, f, file)
				} else {
					// TODO(patrik): Implement transcoding
					// ext := path.Ext(file)
					// fmt.Printf("ext: %v\n", ext)
					// switch ext {
					// case ".mp3":
					// 	return nil
					// case ".opus":
					// 	return nil
					// case ".ogg":
					// 	return nil
					// }

					return pyrin.NoContentNotFound()
				}
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

	return s, nil
}
