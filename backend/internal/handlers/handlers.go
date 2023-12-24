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

	albums, err := apiConfig.queries.GetAlbumsByArtist(c.Context(), id)
	if err != nil {
		return err
	}

	res := struct {
		Artist database.Artist `json:"artist"`
		Albums []database.Album `json:"albums"`
	} {
		Artist: artist,
		Albums: albums,
	}

	return c.JSON(res)
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

	album, err := apiConfig.queries.GetAlbum(c.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": fmt.Sprintf("No album with id: %s", id)})
		} else {
			return err
		}
	}

	tracks, err := apiConfig.queries.GetTracksByAlbum(c.Context(), id);

	res := struct {
		Album database.Album `json:"album"`
		Tracks []database.Track `json:"tracks"`
	} {
		Album: album,
		Tracks: tracks,
	}

	return c.JSON(res)
}

func (apiConfig *ApiConfig) HandlerGetAllTracks(c *fiber.Ctx) error {
	artists, err := apiConfig.queries.GetAllTracks(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"tracks": artists})
}

func (apiConfig *ApiConfig) HandlerGetTrack(c *fiber.Ctx) error {
	id := c.Params("id")

	artist, err := apiConfig.queries.GetTrack(c.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": fmt.Sprintf("No track with id: %s", id)})
		} else {
			return err
		}
	}

	return c.JSON(artist)
}
