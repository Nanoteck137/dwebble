package apis

import (
	"fmt"
	"io"
	"net/http"

	"github.com/faceair/jio"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/handlers"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type playlistApi struct {
	app core.App
}

// TODO(patrik): 'db' should be core.App
// TODO(patrik): Move
func User(db *database.Database, c echo.Context) (*database.User, error) {
	authHeader := c.Request().Header.Get("Authorization")
	tokenString, err := utils.ParseAuthHeader(authHeader)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return []byte(config.LoadedConfig.JwtSecret), nil
		// TODO(patrik): Move back to config JwtSecret
		return []byte("SECRET"), nil
	})

	if err != nil {
		return nil, types.ErrInvalidToken
	}

	jwtValidator := jwt.NewValidator(jwt.WithIssuedAt())

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := jwtValidator.Validate(token.Claims); err != nil {
			return nil, types.ErrInvalidToken
		}

		userId := claims["userId"].(string)
		user, err := db.GetUserById(c.Request().Context(), userId)
		if err != nil {
			return nil, types.ErrInvalidToken
		}

		return &user, nil
	}

	return nil, types.ErrInvalidToken
}

// TODO(patrik): Move
func Decode(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

// TODO(patrik): Move
func Body[T any](c echo.Context, schema jio.Schema) (T, error) {
	var res T

	j, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return res, err
	}

	if len(j) == 0 {
		return res, types.ErrEmptyBody
	}

	data, err := jio.ValidateJSON(&j, schema)
	if err != nil {
		return res, err
	}

	err = Decode(data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (api *playlistApi) HandleGetPlaylists(c echo.Context) error {
	user, err := User(api.app.DB(), c)
	if err != nil {
		return err
	}

	playlists, err := api.app.DB().GetPlaylistsByUser(c.Request().Context(), user.Id)
	if err != nil {
		return err
	}

	res := types.GetPlaylists{
		Playlists: make([]types.Playlist, len(playlists)),
	}

	for i, playlist := range playlists {
		res.Playlists[i] = types.Playlist{
			Id:   playlist.Id,
			Name: playlist.Name,
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (api *playlistApi) HandlePostPlaylist(c echo.Context) error {
	user, err := User(api.app.DB(), c)
	if err != nil {
		return err
	}

	body, err := Body[types.PostPlaylistBody](c, types.PostPlaylistBodySchema)
	if err != nil {
		return err
	}

	playlist, err := api.app.DB().CreatePlaylist(c.Request().Context(), database.CreatePlaylistParams{
		Name:    body.Name,
		OwnerId: user.Id,
	})
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.PostPlaylist{
		Playlist: types.Playlist{
			Id:   playlist.Id,
			Name: playlist.Name,
		},
	}))
}

func (api *playlistApi) HandleGetPlaylistById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := User(api.app.DB(), c)
	if err != nil {
		return err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return err
	}

	if playlist.OwnerId != user.Id {
		return types.ErrNoPlaylist
	}

	items, err := api.app.DB().GetPlaylistItems(c.Request().Context(), playlist.Id)
	if err != nil {
		return err
	}

	tracks := []types.Track{}
	for _, item := range items {
		// TODO(patrik): Optimize
		track, err := api.app.DB().GetTrackById(c.Request().Context(), item.TrackId)
		if err != nil {
			return err
		}

		tracks = append(tracks, types.Track{
			Id:                track.Id,
			Number:            item.ItemIndex,
			Name:              track.Name,
			CoverArt:          handlers.ConvertTrackCoverURL(c, track.CoverArt),
			Duration:          track.Duration,
			BestQualityFile:   handlers.ConvertURL(c, "/tracks/original/"+track.BestQualityFile),
			MobileQualityFile: handlers.ConvertURL(c, "/tracks/mobile/"+track.MobileQualityFile),
			AlbumId:           track.AlbumId,
			ArtistId:          track.ArtistId,
			AlbumName:         track.AlbumName,
			ArtistName:        track.ArtistName,
		})
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.GetPlaylistById{
		Playlist: types.Playlist{
			Id:   playlist.Id,
			Name: playlist.Name,
		},
		Items: tracks,
	}))
}

func (api *playlistApi) HandlePostPlaylistItemsById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := User(api.app.DB(), c)
	if err != nil {
		return err
	}

	body, err := Body[types.PostPlaylistItemsByIdBody](c, types.PostPlaylistItemsByIdBodySchema)
	if err != nil {
		return err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return err
	}

	if playlist.OwnerId != user.Id {
		return types.ErrNoPlaylist
	}

	err = api.app.DB().AddItemsToPlaylist(c.Request().Context(), playlist.Id, body.Tracks)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(nil))
}

func (api *playlistApi) HandleDeletePlaylistItemsById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := User(api.app.DB(), c)
	if err != nil {
		return err
	}

	body, err := Body[types.DeletePlaylistItemsByIdBody](c, types.DeletePlaylistItemsByIdBodySchema)
	if err != nil {
		return err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return err
	}

	if playlist.OwnerId != user.Id {
		return types.ErrNoPlaylist
	}

	err = api.app.DB().DeleteItemsFromPlaylist(c.Request().Context(), playlist.Id, body.TrackIndices)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(nil))
}

func (api *playlistApi) HandlePostPlaylistsItemsMoveById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := User(api.app.DB(), c)
	if err != nil {
		return err
	}

	body, err := Body[types.PostPlaylistsItemMoveByIdBody](c, types.PostPlaylistsItemMoveByIdBodySchema)
	if err != nil {
		return err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return err
	}

	if playlist.OwnerId != user.Id {
		return types.ErrNoPlaylist
	}

	err = api.app.DB().MovePlaylistItem(c.Request().Context(), playlist.Id, body.TrackId, body.ToIndex)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(nil))
}

func InstallPlaylistHandlers(app core.App, group handlers.Group) {
	api := playlistApi{app: app}

	requireSetup := RequireSetup(app)

	group.Register(
		handlers.Handler{
			Name:        "GetPlaylists",
			Path:        "/playlists",
			Method:      http.MethodGet,
			DataType:    types.GetPlaylists{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetPlaylists,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},

		handlers.Handler{
			Name:        "CreatePlaylist",
			Path:        "/playlists",
			Method:      http.MethodPost,
			DataType:    types.PostPlaylist{},
			BodyType:    types.PostPlaylistBody{},
			HandlerFunc: api.HandlePostPlaylist,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},

		handlers.Handler{
			Name:        "GetPlaylistById",
			Path:        "/playlists/:id",
			Method:      http.MethodGet,
			DataType:    types.GetPlaylistById{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetPlaylistById,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},

		handlers.Handler{
			Name:        "AddItemsToPlaylist",
			Path:        "/playlists/:id/items",
			Method:      http.MethodPost,
			DataType:    nil,
			BodyType:    types.PostPlaylistItemsByIdBody{},
			HandlerFunc: api.HandlePostPlaylistItemsById,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},

		handlers.Handler{
			Name:        "DeletePlaylistItems",
			Path:        "/playlists/:id/items",
			Method:      http.MethodDelete,
			DataType:    nil,
			BodyType:    types.DeletePlaylistItemsByIdBody{},
			HandlerFunc: api.HandleDeletePlaylistItemsById,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},

		handlers.Handler{
			Name:        "MovePlaylistItem",
			Path:        "/playlists/:id/items/move",
			Method:      http.MethodPost,
			DataType:    nil,
			BodyType:    types.PostPlaylistItemsByIdBody{},
			HandlerFunc: api.HandlePostPlaylistsItemsMoveById,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},
	)
}
