package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

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
				_ = image

				ctx := c.Request().Context()

				album, err := app.DB().GetAlbumById(ctx, albumId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return pyrin.NoContentNotFound()
					}

					return err
				}

				if !album.CoverArt.Valid {
					// TODO(patrik): Fix
					return pyrin.ServeFile(c, assets.DefaultImagesFS, "default_album.png")
				} 

				// p := app.WorkDir().Album(albumId)
				p := path.Dir(album.CoverArt.String)
				f := os.DirFS(p)

				return pyrin.ServeFile(c, f, path.Base(album.CoverArt.String))
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

				fileExt := path.Ext(file)

				ctx := context.TODO()

				track, err := app.DB().GetTrackById(ctx, trackId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return pyrin.NoContentNotFound()
					}

					return err
				}

				trackDir := app.WorkDir().Track(trackId)

				// Return the original file if the filename matches the
				// one stored inside the track
				if fileExt == path.Ext(track.Filename) {
					d := path.Dir(track.Filename)
					filename := path.Base(track.Filename)

					f := os.DirFS(d)
					return pyrin.ServeFile(c, f, filename)
				}

				return errors.New("not implemented")

				// Here we need to start transcoding the original track
				// media to the requested format

				ext := path.Ext(track.Filename)
				name := strings.TrimSuffix(track.Filename, ext)
				fmt.Printf("name: %v\n", name)

				cacheDir := app.WorkDir().Cache()
				trackCache := cacheDir.Track(track.Id)

				// Make sure that the cache directory is setup
				dirs := []string{
					cacheDir.String(),
					cacheDir.Tracks(),
					trackCache,
				}

				for _, dir := range dirs {
					err = os.Mkdir(dir, 0755)
					if err != nil && !os.IsExist(err) {
						return err
					}
				}

				input := path.Join(trackDir, track.Filename)

				// TODO(patrik): Implement transcoding
				ext = path.Ext(file)
				fmt.Printf("ext: %v\n", ext)
				switch ext {
				case ".mp3":
					name = name + ".mp3"

					p := path.Join(trackCache, name)
					fmt.Printf("p: %v\n", p)

					_, err := os.Stat(p)
					if err != nil {
						if os.IsNotExist(err) {
							cmd := exec.Command("ffmpeg", "-i", input, "-b:a", "320k", p)
							err := cmd.Run()
							if err != nil {
								return err
							}
						} else {
							return err
						}
					}

					f := os.DirFS(trackCache)
					return pyrin.ServeFile(c, f, name)
				case ".opus":
					return nil
				case ".ogg":
					return nil
				}

				return pyrin.NoContentNotFound()
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
