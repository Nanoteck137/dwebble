package handlers

import (
	"fmt"

	"github.com/dwebble/v2/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type ApiConfig struct {
	queries *database.Queries
}

func New(queries *database.Queries) ApiConfig {
	return ApiConfig{
		queries: queries,
	}
}

func (apiConfig *ApiConfig) HandlerGetAllArtists(c *fiber.Ctx) error {
	artists, err := apiConfig.queries.GetAllArtists(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"artists": artists})
}

func (apiConfig *ApiConfig) HandlerGetArtist(c *fiber.Ctx) error {
	id := c.Params("id")

	artist, err := apiConfig.queries.GetArtist(c.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": fmt.Sprintf("No artist with id: %s", id)})
		} else {
			return err
		}
	}

	return c.JSON(artist)
}

func (apiConfig *ApiConfig) HandlerGetAllAlbums(c *fiber.Ctx) error {
	artists, err := apiConfig.queries.GetAllAlbums(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"albums": artists})
}

func (apiConfig *ApiConfig) HandlerGetAlbum(c *fiber.Ctx) error {
	id := c.Params("id")

	artist, err := apiConfig.queries.GetAlbum(c.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": fmt.Sprintf("No album with id: %s", id)})
		} else {
			return err
		}
	}

	return c.JSON(artist)
}
