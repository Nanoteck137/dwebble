package apis

import (
	"net/http"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

type authApi struct {
	app core.App
}

// TODO(patrik): Check confirmPassword
func (api *authApi) HandlePostSignup(c pyrin.Context) (any, error) {
	// TODO(patrik): fix
	// body, err := Body[types.PostAuthSignupBody](c)
	// if err != nil {
	// 	return err
	// }
	//
	// user, err := api.app.DB().CreateUser(c.Request().Context(), body.Username, body.Password)
	// if err != nil {
	// 	return err
	// }
	//
	// return c.JSON(200, SuccessResponse(types.PostAuthSignup{
	// 	Id:       user.Id,
	// 	Username: user.Username,
	// }))
	return nil, nil
}

func (api *authApi) HandlePostSignin(c pyrin.Context) (any, error) {
	// TODO(patrik): fix
	// body, err := Body[types.PostAuthSigninBody](c)
	// if err != nil {
	// 	return err
	// }
	//
	// user, err := api.app.DB().GetUserByUsername(c.Request().Context(), body.Username)
	// if err != nil {
	// 	return err
	// }
	//
	// if user.Password != body.Password {
	// 	// TODO(patrik): Fix error
	// 	return errors.New("Incorrect creds")
	// }
	//
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"userId": user.Id,
	// 	"iat":    time.Now().Unix(),
	// 	// "exp":    time.Now().Add(1000 * time.Second).Unix(),
	// })
	//
	// tokenString, err := token.SignedString(([]byte)(api.app.Config().JwtSecret))
	// if err != nil {
	// 	return err
	// }
	//
	// return c.JSON(200, SuccessResponse(types.PostAuthSignin{
	// 	Token: tokenString,
	// }))
	return nil, nil
}

func (api *authApi) HandleGetMe(c pyrin.Context) (any, error) {
	// TODO(patrik): fix
	// user, err := User(api.app, c)
	// if err != nil {
	// 	return err
	// }
	//
	// isOwner := api.app.DBConfig().OwnerId == user.Id
	//
	// return c.JSON(200, SuccessResponse(types.GetAuthMe{
	// 	Id:       user.Id,
	// 	Username: user.Username,
	// 	IsOwner:  isOwner,
	// }))
	return nil, nil
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
