package apis

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
)

func RequireSetup(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !app.IsSetup() {
				// TODO(patrik): Fix error
				return errors.New("Server not setup")
			}

			return next(c)
		}
	}
}
