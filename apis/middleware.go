package apis

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/handlers"
	"github.com/nanoteck137/dwebble/types"
)

func RequireSetup() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO(patrik): Fix
			if !handlers.IsSetup() {
				return types.NewApiError(400, "Server not setup")
			}

			return next(c)
		}
	}
}
