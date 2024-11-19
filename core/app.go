package core

import (
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
)

// Inspiration from Pocketbase: https://github.com/pocketbase/pocketbase
// File: https://github.com/pocketbase/pocketbase/blob/master/core/app.go
type App interface {
	DB() *database.Database
	Config() *config.Config
	DBConfig() *database.Config

	WorkDir() types.WorkDir

	IsSetup() bool
	UpdateDBConfig(conf *database.Config)

	Bootstrap() error
}
