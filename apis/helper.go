package apis

import (
	"errors"
	"fmt"
	"io"

	"github.com/faceair/jio"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

func User(app core.App, c pyrin.Context) (*database.User, error) {
	authHeader := c.Request().Header.Get("Authorization")
	tokenString := utils.ParseAuthHeader(authHeader)
	if tokenString == "" {
		// TODO(patrik): Fix error
		return nil, errors.New("Invalid auth header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(app.Config().JwtSecret), nil
	})

	if err != nil {
		// TODO(patrik): Fix error
		return nil, errors.New("Invalid token")
	}

	jwtValidator := jwt.NewValidator(jwt.WithIssuedAt())

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := jwtValidator.Validate(token.Claims); err != nil {
			// TODO(patrik): Fix error
			return nil, errors.New("Invalid token")
		}

		userId := claims["userId"].(string)
		user, err := app.DB().GetUserById(c.Request().Context(), userId)
		if err != nil {
			// TODO(patrik): Fix error
			return nil, errors.New("Invalid token")
		}

		return &user, nil
	}

	// TODO(patrik): Fix error
	return nil, errors.New("Invalid token")
}

func Body[T types.Body](c pyrin.Context) (T, error) {
	var res T

	schema := res.Schema()

	j, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return res, err
	}

	if len(j) == 0 {
		// TODO(patrik): Fix error
		return res, errors.New("Invalid body")
	}

	data, err := jio.ValidateJSON(&j, schema)
	if err != nil {
		return res, err
	}

	err = utils.Decode(data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
