package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nanoteck137/dwebble/v2/database"
	"github.com/nanoteck137/dwebble/v2/handlers"
	"github.com/nanoteck137/dwebble/v2/types"
)

type Api struct {
}

func New(db *pgxpool.Pool) *fiber.App {
	router := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			switch err := err.(type) {
			case handlers.ApiError:
				return c.Status(err.Status).JSON(err)
			case types.ApiError:
				return c.Status(err.Status).JSON(err)
			case *fiber.Error:
				return c.Status(err.Code).JSON(handlers.ApiError{
					Status:  err.Code,
					Message: err.Error(),
					Data:    nil,
				})
			}

			return c.Status(http.StatusInternalServerError).JSON(handlers.ApiError{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		},
	})

	queries := database.New(db)
	apiConfig := handlers.New(db, queries, "./work")

	router.Get("/artists", apiConfig.HandlerGetAllArtists)
	router.Post("/artists", apiConfig.HandlerCreateArtist)
	router.Get("/artists/:id", apiConfig.HandlerGetArtist)

	router.Get("/albums", apiConfig.HandlerGetAllAlbums)
	router.Post("/albums", apiConfig.HandlerCreateAlbum)
	router.Get("/albums/:id", apiConfig.HandlerGetAlbum)

	router.Get("/tracks", apiConfig.HandlerGetAllTracks)
	router.Post("/tracks", apiConfig.HandlerCreateTrack)
	router.Get("/tracks/:id", apiConfig.HandlerGetTrack)

	router.Post("/queue/album", apiConfig.HandlerCreateQueueFromAlbum)

	return router
}
