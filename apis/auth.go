package apis

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

type authApi struct {
	app core.App
}

// TODO(patrik): Check confirmPassword
func (api *authApi) HandlePostSignup(c pyrin.Context) (any, error) {
	body, err := Body[types.PostAuthSignupBody](c)
	if err != nil {
		return nil, err
	}

	user, err := api.app.DB().CreateUser(c.Request().Context(), body.Username, body.Password)
	if err != nil {
		return nil, err
	}

	return types.PostAuthSignup{
		Id:       user.Id,
		Username: user.Username,
	}, nil
}

func (api *authApi) HandlePostSignin(c pyrin.Context) (any, error) {
	body, err := Body[types.PostAuthSigninBody](c)
	if err != nil {
		return nil, err
	}

	user, err := api.app.DB().GetUserByUsername(c.Request().Context(), body.Username)
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

	tokenString, err := token.SignedString(([]byte)(api.app.Config().JwtSecret))
	if err != nil {
		return nil, err
	}

	return types.PostAuthSignin{
		Token: tokenString,
	}, nil
}

func (api *authApi) HandleGetMe(c pyrin.Context) (any, error) {
	user, err := User(api.app, c)
	if err != nil {
		return nil, err
	}

	isOwner := api.app.DBConfig().OwnerId == user.Id

	return types.GetAuthMe{
		Id:       user.Id,
		Username: user.Username,
		IsOwner:  isOwner,
	}, nil
}

func InstallAuthHandlers(app core.App, group pyrin.Group) {
	api := authApi{app: app}

	group.Register(
		pyrin.ApiHandler{
			Name:        "Signup",
			Path:        "/auth/signup",
			Method:      http.MethodPost,
			DataType:    types.PostAuthSignup{},
			BodyType:    types.PostAuthSignupBody{},
			HandlerFunc: api.HandlePostSignup,
		},

		pyrin.ApiHandler{
			Name:        "Signin",
			Path:        "/auth/signin",
			Method:      http.MethodPost,
			DataType:    types.PostAuthSignin{},
			BodyType:    types.PostAuthSigninBody{},
			HandlerFunc: api.HandlePostSignin,
		},

		pyrin.ApiHandler{
			Name:        "GetMe",
			Path:        "/auth/me",
			Method:      http.MethodGet,
			DataType:    types.GetAuthMe{},
			HandlerFunc: api.HandleGetMe,
		},
	)
}
