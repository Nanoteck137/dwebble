package apis

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/types"
)

func RequireSetup(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !app.IsSetup() {
				return types.NewApiError(400, "Server not setup")
			}

			return next(c)
		}
	}
}
