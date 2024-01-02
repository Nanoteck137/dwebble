package handlers

import (
	"fmt"

	"github.com/nanoteck137/dwebble/v2/internal/database"
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

	for i := range artists {
		a := &artists[i]
		a.Picture = ConvertURL(c, a.Picture)
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

	artist.Picture = ConvertURL(c, artist.Picture)

	albums, err := apiConfig.queries.GetAlbumsByArtist(c.Context(), id)
	if err != nil {
		return err
	}

	for i := range albums {
		a := &albums[i]
		a.CoverArt = ConvertURL(c, a.CoverArt)
	}

	res := struct {
		Artist database.Artist  `json:"artist"`
		Albums []database.Album `json:"albums"`
	}{
		Artist: artist,
		Albums: albums,
	}

	return c.JSON(res)
}

func (apiConfig *ApiConfig) HandlerGetAllAlbums(c *fiber.Ctx) error {
	albums, err := apiConfig.queries.GetAllAlbums(c.Context())
	if err != nil {
		return err
	}

	for i := range albums {
		a := &albums[i]
		a.CoverArt = ConvertURL(c, a.CoverArt)
	}

	return c.JSON(fiber.Map{"albums": albums})
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

	album.CoverArt = ConvertURL(c, album.CoverArt)

	tracks, err := apiConfig.queries.GetTracksByAlbum(c.Context(), id)
	if err != nil {
		return err
	}

	for i := range tracks {
		t := &tracks[i]
		t.CoverArt = ConvertURL(c, t.CoverArt)
		t.Filename = ConvertURL(c, "/tracks/" + t.Filename)
	}

	res := struct {
		Album  database.Album                 `json:"album"`
		Tracks []database.GetTracksByAlbumRow `json:"tracks"`
	}{
		Album:  album,
		Tracks: tracks,
	}

	return c.JSON(res)
}

func ConvertURL(c *fiber.Ctx, path string) string {
	return c.BaseURL() + path
}

func (apiConfig *ApiConfig) HandlerGetAllTracks(c *fiber.Ctx) error {
	tracks, err := apiConfig.queries.GetAllTracks(c.Context())
	if err != nil {
		return err
	}

	for i := range tracks {
		a := &tracks[i]
		a.CoverArt = ConvertURL(c, a.CoverArt)
	}

	return c.JSON(fiber.Map{"tracks": tracks})
}

func (apiConfig *ApiConfig) HandlerGetTrack(c *fiber.Ctx) error {
	id := c.Params("id")

	track, err := apiConfig.queries.GetTrack(c.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": fmt.Sprintf("No track with id: %s", id)})
		} else {
			return err
		}
	}

	track.CoverArt = ConvertURL(c, track.CoverArt)

	return c.JSON(track)
}
