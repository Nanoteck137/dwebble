package core

import (
	"context"
	"errors"
	"os"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	db       *database.Database
	config   *config.Config
	dbConfig *database.Config
}

func (app *BaseApp) DB() *database.Database {
	return app.db
}

func (app *BaseApp) Config() *config.Config {
	return app.config
}

func (app *BaseApp) DBConfig() *database.Config {
	return app.dbConfig
}

func (app *BaseApp) WorkDir() types.WorkDir {
	return app.config.WorkDir()
}

func (app *BaseApp) Bootstrap() error {
	var err error

	workDir := app.config.WorkDir()

	dirs := []string{
		workDir.Artists(),
		workDir.Albums(),
		workDir.Tracks(),
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	app.db, err = database.Open(workDir)
	if err != nil {
		return err
	}

	err = database.RunMigrateUp(app.db)
	if err != nil {
		return err
	}

	app.dbConfig, err = app.db.GetConfig(context.Background())
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			db, tx, err := app.DB().Begin()
			if err != nil {
				return err
			}
			defer tx.Rollback()

			ctx := context.Background()
			username := app.config.Username
			password := app.config.InitialPassword

			user, err := db.CreateUser(ctx, username, password)
			if err != nil {
				return err
			}

			conf, err := db.CreateConfig(ctx, user.Id)
			if err != nil {
				return err
			}

			err = tx.Commit()
			if err != nil {
				return err
			}

			app.dbConfig = &conf
		} else {
			return err
		}
	}

	ctx := context.TODO()
	db := app.db

	_, err = db.Conn.Exec(`
		DROP TABLE IF EXISTS test;
	`)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec(`
		CREATE VIRTUAL TABLE test USING fts5(name, artist, album);
	`)
	if err != nil {
		return err
	}

	tracks, err := db.GetAllTracks(ctx, "", "", true)
	if err != nil {
		return err
	}

	for _, track := range tracks {
		_, err := db.Conn.Exec(`
			INSERT INTO test(name, artist, album) VALUES (?, ?, ?)
		`, track.Name, track.ArtistName, track.AlbumName)
		if err != nil {
			return err
		}
	}

	rows, err := db.Conn.Query("SELECT * FROM test WHERE test MATCH 'pantera from hell' ORDER BY rank")
	if err != nil {
		return err
	}
	cols, _ := rows.Columns()

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		pretty.Println(m)
	}

	return nil
}

func NewBaseApp(config *config.Config) *BaseApp {
	return &BaseApp{
		config: config,
	}
}
