package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/helper"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

const defaultMemory = 32 << 20 // 32 MB

func User(app core.App, c pyrin.Context) (*database.User, error) {
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
			// TODO(patrik): Should we handle error here?
			return nil, err
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
	UNKNOWN_ARTIST_ID = "unknown"
	UNKNOWN_ARTIST_NAME = "UNKNOWN"
)

func EnsureUnknownArtistExists(ctx context.Context, db *database.Database, workDir types.WorkDir) error {
	_, err := db.GetArtistById(ctx, UNKNOWN_ARTIST_ID)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			log.Info("Creating 'unknown' artist")
			_, err = helper.CreateArtist(ctx, db, workDir, database.CreateArtistParams{
				Id:   UNKNOWN_ARTIST_ID,
				Name: UNKNOWN_ARTIST_NAME,
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
