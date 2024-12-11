package cli

import (
	"context"

	"github.com/nanoteck137/dwebble/apis"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		app := core.NewBaseApp(&config.LoadedConfig)

		err := app.Bootstrap()
		if err != nil {
			log.Fatal("Failed to bootstrap app", "err", err)
		}

		{
			ctx := context.TODO()
			artist, err := app.DB().CreateArtist(ctx, database.CreateArtistParams{
				Name: "Test Artist",
			})
			if err != nil {
				log.Fatal("Failed", "err", err)
			}

			album, err := app.DB().CreateAlbum(ctx, database.CreateAlbumParams{
				Name:      "Test Album",
				ArtistId:  artist.Id,
			})
			if err != nil {
				log.Fatal("Failed", "err", err)
			}

			extraArtist, err := app.DB().CreateArtist(ctx, database.CreateArtistParams{
				Name: "Test Artist",
			})
			if err != nil {
				log.Fatal("Failed", "err", err)
			}

			err = app.DB().AddExtraArtistToAlbum(ctx, album.Id, extraArtist.Id)
			if err != nil {
				log.Fatal("Failed", "err", err)
			}

			err = app.DB().AddExtraArtistToAlbum(ctx, album.Id, extraArtist.Id)
			if err != nil {
				log.Fatal("Failed", "err", err)
			}
		}

		e, err := apis.Server(app)
		if err != nil {
			log.Fatal("Failed to create server", "err", err)
		}

		err = e.Start(app.Config().ListenAddr)
		if err != nil {
			log.Fatal("Failed to start server", "err", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
