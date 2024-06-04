package handlers

import (
	"context"

	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
)

var cfg *database.Config

func IsSetup() bool {
	return cfg != nil
}

func (h *Handlers) HandleGetSystemInfo(c echo.Context) error {
	return c.JSON(200, types.NewApiSuccessResponse(types.GetSystemInfo{
		Version: config.Version,
		IsSetup: IsSetup(),
	}))
}

func (h *Handlers) HandlePostSystemSetup(c echo.Context) error {
	if IsSetup() {
		return types.NewApiError(400, "System already setup")
	}

	body, err := Body[types.PostSystemSetupBody](c, types.PostSystemSetupBodySchema)
	if err != nil {
		return err
	}

	db, tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	user, err := db.CreateUser(c.Request().Context(), body.Username, body.Password)
	if err != nil {
		return err
	}

	conf, err := db.CreateConfig(c.Request().Context(), user.Id)
	if err != nil {
		return err
	}

	pretty.Println(body)

	err = tx.Commit()
	if err != nil {
		return err
	}

	cfg = &conf


	return c.JSON(200, types.NewApiSuccessResponse(nil))
}


func (h *Handlers) InstallSystemHandlers(group *echo.Group) {
	group.GET("/system/info", h.HandleGetSystemInfo)
	group.POST("/system/setup", h.HandlePostSystemSetup)
}

func InitializeConfig(db *database.Database) error {
	conf, err := db.GetConfig(context.Background())
	if err != nil {
		return err
	}

	cfg = conf

	return nil
}
