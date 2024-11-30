package apis

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/validate"
)

type Signup struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type SignupBody struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

var usernameRegex = regexp.MustCompile("^[a-zA-Z0-9-]+$")

func (b *SignupBody) Transform() {
	b.Username = strings.TrimSpace(b.Username)
}

func (b SignupBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Username, validate.Required, validate.Length(4, 32), validate.Match(usernameRegex).Error("not valid username")),
		validate.Field(&b.Password, validate.Required, validate.Length(8, 32)),
		validate.Field(&b.PasswordConfirm, validate.Required, validate.By(func(value interface{}) error {
			s, _ := value.(string)

			if s != b.Password {
				return errors.New("password mismatch")
			}

			return nil
		})),
	)
}

type Signin struct {
	Token string `json:"token"`
}

type SigninBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordBody struct {
	CurrentPassword    string `json:"currentPassword"`
	NewPassword        string `json:"newPassword"`
	NewPasswordConfirm string `json:"newPasswordConfirm"`
}

func (b ChangePasswordBody) Validate() error {
	return validate.ValidateStruct(
		&b,
		validate.Field(&b.CurrentPassword, validate.Required),
		validate.Field(&b.NewPassword, validate.Required),
		validate.Field(&b.NewPasswordConfirm, validate.Required, validate.By(func(value interface{}) error {
			s, _ := value.(string)
			if s != b.NewPassword {
				return errors.New("password mismatch")
			}

			return nil
		})),
	)
}

type GetMe struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	IsOwner  bool   `json:"isOwner"`
}

func InstallAuthHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "Signup",
			Path:     "/auth/signup",
			Method:   http.MethodPost,
			DataType: Signup{},
			BodyType: SignupBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[SignupBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				_, err = app.DB().GetUserByUsername(ctx, body.Username)
				if err == nil {
					// TODO(patrik): Better error
					return nil, errors.New("user already exists")
				}

				if !errors.Is(err, database.ErrItemNotFound) {
					return nil, err
				}

				user, err := app.DB().CreateUser(ctx, body.Username, body.Password)
				if err != nil {
					return nil, err
				}

				return Signup{
					Id:       user.Id,
					Username: user.Username,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "Signin",
			Path:     "/auth/signin",
			Method:   http.MethodPost,
			DataType: Signin{},
			BodyType: SigninBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[SigninBody](c)
				if err != nil {
					return nil, err
				}

				user, err := app.DB().GetUserByUsername(c.Request().Context(), body.Username)
				if err != nil {
					// TODO(patrik): Check for ErrItemNotFound
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

				return Signin{
					Token: tokenString,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "ChangePassword",
			Path:     "/auth/password",
			Method:   http.MethodPatch,
			BodyType: ChangePasswordBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				body, err := pyrin.Body[ChangePasswordBody](c)
				if err != nil {
					return nil, err
				}

				// TODO(patrik): Check body.CurrentPassword

				if user.Password != body.CurrentPassword {
					// TODO(patrik): Better error
					return nil, errors.New("Password not matching")
				}

				err = app.DB().UpdateUser(ctx, user.Id, database.UserChanges{
					Password: types.Change[string]{
						Value:   body.NewPassword,
						Changed: true,
					},
				})
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetMe",
			Path:     "/auth/me",
			Method:   http.MethodGet,
			DataType: GetMe{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				isOwner := app.DBConfig().OwnerId == user.Id

				return GetMe{
					Id:       user.Id,
					Username: user.Username,
					IsOwner:  isOwner,
				}, nil
			},
		},
	)
}
