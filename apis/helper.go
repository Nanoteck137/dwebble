package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

type UserCheckFunc func(user *database.User) error

func RequireAdmin(user *database.User) error {
	if user.Role != types.RoleSuperUser && user.Role != types.RoleAdmin {
		return InvalidAuth("user requires 'super_user' or 'admin' role")
	}

	return nil
}

func User(app core.App, c pyrin.Context, checks ...UserCheckFunc) (*database.User, error) {
	user, err := getUser(app, c)
	if err != nil {
		return nil, err
	}

	for _, check := range checks {
		err := check(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func getUser(app core.App, c pyrin.Context) (*database.User, error) {
	apiTokenHeader := c.Request().Header.Get("X-Api-Token")
	if apiTokenHeader != "" {
		ctx := context.TODO()
		token, err := app.DB().GetApiTokenById(ctx, apiTokenHeader)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				return nil, InvalidAuth("invalid api token")
			}

			return nil, err
		}

		user, err := app.DB().GetUserById(c.Request().Context(), token.UserId)
		if err != nil {
			return nil, InvalidAuth("invalid api token")
		}

		return &user, nil
	}

	authHeader := c.Request().Header.Get("Authorization")
	tokenString := utils.ParseAuthHeader(authHeader)
	if tokenString == "" {
		return nil, InvalidAuth("invalid authorization header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(app.Config().JwtSecret), nil
	})

	if err != nil {
		// TODO(patrik): Handle error better
		return nil, InvalidAuth("invalid authorization token")
	}

	jwtValidator := jwt.NewValidator(jwt.WithIssuedAt())

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := jwtValidator.Validate(token.Claims); err != nil {
			return nil, InvalidAuth("invalid authorization token")
		}

		userId := claims["userId"].(string)
		user, err := app.DB().GetUserById(c.Request().Context(), userId)
		if err != nil {
			return nil, InvalidAuth("invalid authorization token")
		}

		return &user, nil
	}

	return nil, InvalidAuth("invalid authorization token")
}

func ConvertSqlNullString(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}

	return nil
}

func ConvertSqlNullInt64(value sql.NullInt64) *int64 {
	if value.Valid {
		return &value.Int64
	}

	return nil
}

const (
	UNKNOWN_ARTIST_ID   = "unknown"
	UNKNOWN_ARTIST_NAME = "UNKNOWN"
)

func EnsureUnknownArtistExists(ctx context.Context, db *database.Database, workDir types.WorkDir) error {
	_, err := db.GetArtistById(ctx, UNKNOWN_ARTIST_ID)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			log.Info("Creating 'unknown' artist")
			_, err := db.CreateArtist(ctx, database.CreateArtistParams{
				Id:   UNKNOWN_ARTIST_ID,
				Name: UNKNOWN_ARTIST_NAME,
				Slug: utils.Slug(UNKNOWN_ARTIST_NAME),
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

const (
	DefaultArtistPictureName = "default/default_artist.png"
	DefaultAlbumCoverArtName = "default/default_album.png"
)

func ConvertURL(c pyrin.Context, path string) string {
	host := c.Request().Host

	scheme := "http"

	h := c.Request().Header.Get("X-Forwarded-Proto")
	if h != "" {
		scheme = h
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

func ConvertImageURL(c pyrin.Context, albumId string, val sql.NullString, def string) string {
	coverArt := def
	if val.Valid && val.String != "" {
		coverArt = val.String
	}

	return ConvertURL(c, "/files/albums/images/"+albumId+"/"+coverArt)
}

func ConvertArtistPicture(c pyrin.Context, artistId string, val sql.NullString) types.Images {
	if val.Valid && val.String != "" {
		first := "/files/artists/" + artistId + "/"
		return types.Images{
			Original: ConvertURL(c, first+val.String),
			Small:    ConvertURL(c, first+"picture-128.png"),
			Medium:   ConvertURL(c, first+"picture-256.png"),
			Large:    ConvertURL(c, first+"picture-512.png"),
		}
	}

	url := ConvertURL(c, "/files/images/"+DefaultArtistPictureName)
	return types.Images{
		Original: url,
		Small:    url,
		Medium:   url,
		Large:    url,
	}
}

func ConvertAlbumCoverURL(c pyrin.Context, albumId string, val sql.NullString) types.Images {
	if val.Valid && val.String != "" {
		coverArt := val.String
		originalExt := path.Ext(coverArt)
		return types.Images{
			// TODO(patrik): Move to const (cover-128.png, cover-256.png, cover-512.png)
			Original: ConvertURL(c, "/files/albums/images/"+albumId+"/"+"original"+originalExt),
			Small:    ConvertURL(c, "/files/albums/images/"+albumId+"/"+"128.png"),
			Medium:   ConvertURL(c, "/files/albums/images/"+albumId+"/"+"256.png"),
			Large:    ConvertURL(c, "/files/albums/images/"+albumId+"/"+"512.png"),
		}
	}

	url := ConvertURL(c, "/files/images/"+DefaultAlbumCoverArtName)
	return types.Images{
		Original: url,
		Small:    url,
		Medium:   url,
		Large:    url,
	}
}
