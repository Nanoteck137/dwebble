package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
)

// HandleGetArtists godoc
//
//	@Summary		Get all artists
//	@Description	Get all artists
//	@Tags			artists
//	@Produce		json
//	@Success		200	{object}	types.ApiResponse[types.ApiGetArtistsData]
//	@Failure		400	{object}	types.ApiError
//	@Failure		500	{object}	types.ApiError
//	@Router			/artists [get]
func (api *ApiConfig) HandleGetArtists(c *fiber.Ctx) error {
	var artists []database.Artist

	name := c.Query("name")
	if name != "" {
		a, err := api.queries.GetArtistByName(c.UserContext(), name)
		if err != nil {
			return err
		}

		artists = a
	} else {
		a, err := api.queries.GetAllArtists(c.UserContext())
		if err != nil {
			return err
		}

		artists = a
	}

	result := types.ApiGetArtistsData{
		Artists: make([]types.ApiArtist, len(artists)),
	}

	for i, artist := range artists {
		result.Artists[i] = types.ApiArtist{
			Id:      artist.ID,
			Name:    artist.Name,
			Picture: ConvertURL(c, "/images/"+artist.Picture),
		}
	}

	return c.JSON(types.NewApiResponse(result))
}

type CreateArtistBody struct {
	Name string `json:"name" form:"name" validate:"required"`
}

// HandlePostArtist godoc
//
//	@Summary		Create new artist
//	@Description	Create new artist
//	@Tags			artists
//	@Accept			mpfd
//	@Produce		json
//	@Param			name	formData	string	true	"Artist name"
//	@Success		200		{object}	types.ApiResponse[types.ApiPostArtistData]
//	@Failure		400		{object}	types.ApiError
//	@Failure		500		{object}	types.ApiError
//	@Router			/artists [post]
func (api *ApiConfig) HandlePostArtist(c *fiber.Ctx) error {
	var body CreateArtistBody
	err := c.BodyParser(&body)
	if err != nil {
		return types.ApiBadRequestError("Failed to create artist: " + err.Error())
	}

	errs := api.validateBody(body)
	if errs != nil {
		return types.ApiBadRequestError("Failed to create artist", errs)
	}

	artist, err := api.queries.CreateArtist(c.UserContext(), database.CreateArtistParams{
		ID:      utils.CreateId(),
		Name:    body.Name,
		Picture: "TODO",
	})

	if err != nil {
		return err
	}

	return c.JSON(types.NewApiResponse(types.ApiPostArtistData{
		Id:      artist.ID,
		Name:    artist.Name,
		Picture: ConvertURL(c, "/images/"+artist.Picture),
	}))
}

// HandleGetArtistById godoc
//
//	@Summary		Get artist by id
//	@Description	Get artist by id
//	@Tags			artists
//	@Produce		json
//	@Param			id	path		string	true	"Artist Id"
//	@Success		200	{object}	types.ApiResponse[types.ApiGetArtistByIdData]
//	@Failure		400	{object}	types.ApiError
//	@Failure		404	{object}	types.ApiError
//	@Failure		500	{object}	types.ApiError
//	@Router			/artists/{id} [get]
func (api *ApiConfig) HandleGetArtistById(c *fiber.Ctx) error {
	id := c.Params("id")

	artist, err := api.queries.GetArtist(c.UserContext(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return types.ApiNotFoundError(fmt.Sprintf("No artist with id: '%s'", id))
		} else {
			return err
		}
	}

	return c.JSON(types.NewApiResponse(types.ApiGetArtistByIdData{
		Id:      artist.ID,
		Name:    artist.Name,
		Picture: ConvertURL(c, "/images/"+artist.Picture),
	}))
}

// HandleGetArtistAlbumsById godoc
//
//	@Summary		Get all albums by artist
//	@Description	Get all albums by artist
//	@Tags			artists
//	@Produce		json
//	@Param			id	path		string	true	"Artist Id"
//	@Success		200	{object}	types.ApiResponse[types.ApiGetArtistAlbumsByIdData]
//	@Failure		400	{object}	types.ApiError
//	@Failure		500	{object}	types.ApiError
//	@Router			/artists/{id}/albums [get]
func (api *ApiConfig) HandleGetArtistAlbumsById(c *fiber.Ctx) error {
	id := c.Params("id")

	var albums []database.Album

	nameQuery := c.Query("name")
	if nameQuery != "" {
		a, err := api.queries.GetAlbumsByArtistAndName(c.UserContext(), database.GetAlbumsByArtistAndNameParams{
			ArtistID: id,
			Name:     nameQuery,
		})

		if err != nil {
			return err
		}

		albums = a
	} else {
		a, err := api.queries.GetAlbumsByArtist(c.UserContext(), id)
		if err != nil {
			return err
		}

		albums = a
	}


	result := types.ApiGetArtistAlbumsByIdData{
		Albums: make([]types.ApiAlbum, len(albums)),
	}

	for i, album := range albums {
		result.Albums[i] = types.ApiAlbum{
			Id:       album.ID,
			Name:     album.Name,
			CoverArt: ConvertURL(c, "/images/"+album.CoverArt),
			ArtistId: album.ArtistID,
		}
	}

	return c.JSON(types.NewApiResponse(result))
}

func InstallArtistHandlers(router *fiber.App, apiConfig *ApiConfig) {
	router.Get("/artists", apiConfig.HandleGetArtists)
	router.Post("/artists", apiConfig.HandlePostArtist)

	router.Get("/artists/:id", apiConfig.HandleGetArtistById)
	router.Get("/artists/:id/albums", apiConfig.HandleGetArtistAlbumsById)
}
