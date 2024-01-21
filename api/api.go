package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/handlers"
	"github.com/nanoteck137/dwebble/types"
)

type Api struct {
}

//	@title		Testing Testing
//	@BasePath	/api/v1
func New(db *pgxpool.Pool) *fiber.App {
	router := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			switch err := err.(type) {
			case types.ApiError:
				return c.Status(err.Status).JSON(err)
			case *fiber.Error:
				return c.Status(err.Code).JSON(types.NewApiError(err.Code, err.Error()))
			}

			return c.Status(http.StatusInternalServerError).JSON(types.NewApiError(
				http.StatusInternalServerError,
				err.Error(),
			))
		},
	})

	queries := database.New(db)
	apiConfig := handlers.New(db, queries, "./work")

	handlers.InstallArtistHandlers(router, &apiConfig)
	handlers.InstallAlbumHandlers(router, &apiConfig)
	handlers.InstallTrackHandlers(router, &apiConfig)

	router.Post("/queue/album", apiConfig.HandlerCreateQueueFromAlbum)

	return router
}
