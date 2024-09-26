package cli

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/kr/pretty"
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

var testCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		// workDir, err := config.LoadedConfig.BootstrapDataDir()
		// if err != nil {
		// 	log.Fatal("Failed to bootstrap data dir", "err", err)
		// }
		//
		// db, err := database.Open(workDir)
		// if err != nil {
		// 	log.Fatal("Failed to open database", "err", err)
		// }

		app := core.NewBaseApp(&config.LoadedConfig)

		err := app.Bootstrap()
		if err != nil {
			log.Fatal("Failed to bootstrap app", "err", err)
		}

		db := app.DB().RawConn
		dbx := sqlx.NewDb(db, "sqlite3")

		tracks := database.TrackQuery().As("tracks")

		dialect := goqu.Dialect("sqlite_returning")

		ds := dialect.From("playlist_items").
			Select(
				"tracks.*",
				goqu.I("playlist_items.item_index").As("track_number"),
			).
			Join(tracks, goqu.On(goqu.I("tracks.id").Eq(goqu.I("playlist_items.track_id")))).
			Where(goqu.I("playlist_id").Eq("iodzf6saku2kmylysvwzecavtoxfiruc")).
			Order(goqu.I("item_index").Asc()).
			Prepared(true)

		sql, params, err := ds.ToSQL()
		if err != nil {
			log.Fatal("Failed", "err", err)
		}

		var res []database.Track
		err = dbx.Select(&res, sql, params...)
		if err != nil {
			log.Fatal("Failed", "err", err)
		}

		pretty.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(testCmd)
}
