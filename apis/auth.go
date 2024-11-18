package apis

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

func InstallAuthHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "Signup",
			Path:     "/auth/signup",
			Method:   http.MethodPost,
			DataType: types.PostAuthSignup{},
			BodyType: types.PostAuthSignupBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Check confirmPassword
				body, err := Body[types.PostAuthSignupBody](c)
				if err != nil {
					return nil, err
				}

				user, err := app.DB().CreateUser(c.Request().Context(), body.Username, body.Password)
				if err != nil {
					return nil, err
				}

				return types.PostAuthSignup{
					Id:       user.Id,
					Username: user.Username,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "Signin",
			Path:     "/auth/signin",
			Method:   http.MethodPost,
			DataType: types.PostAuthSignin{},
			BodyType: types.PostAuthSigninBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := Body[types.PostAuthSigninBody](c)
				if err != nil {
					return nil, err
				}

				user, err := app.DB().GetUserByUsername(c.Request().Context(), body.Username)
				if err != nil {
					return nil, err
				}

				if user.Password != body.Password {
					// TODO(patrik): Fix error
					return nil, errors.New("Incorrect creds")
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"userId": user.Id,
					"iat":    time.Now().Unix(),
					// "exp":    time.Now().Add(1000 * time.Second).Unix(),
				})

				tokenString, err := token.SignedString(([]byte)(app.Config().JwtSecret))
				if err != nil {
					return nil, err
				}

				return types.PostAuthSignin{
					Token: tokenString,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetMe",
			Path:     "/auth/me",
			Method:   http.MethodGet,
			DataType: types.GetAuthMe{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				isOwner := app.DBConfig().OwnerId == user.Id

				return types.GetAuthMe{
					Id:       user.Id,
					Username: user.Username,
					IsOwner:  isOwner,
				}, nil
			},
		},
	)
}
