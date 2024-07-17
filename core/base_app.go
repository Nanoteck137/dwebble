package core

import (
	"context"
	"os"

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

func (app *BaseApp) IsSetup() bool {
	return app.dbConfig != nil
}

func (app *BaseApp) UpdateDBConfig(conf *database.Config) {
	app.dbConfig = conf
}

func (app *BaseApp) Bootstrap() error {
	var err error

	workDir := app.config.WorkDir()

	app.db, err = database.Open(workDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(workDir.OriginalTracksDir(), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(workDir.MobileTracksDir(), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(workDir.ImagesDir(), 0755)
	if err != nil {
		return err
	}

	app.dbConfig, err = app.db.GetConfig(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func NewBaseApp(config *config.Config) *BaseApp {
	return &BaseApp{
		config: config,
	}
}
