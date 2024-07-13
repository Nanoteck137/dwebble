package handlers

import (
	"fmt"
	"io"
	"time"

	"github.com/faceair/jio"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
)

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

// TODO(patrik): Check confirmPassword
func (h *Handlers) HandlePostSignup(c echo.Context) error {
	body, err := Body[types.PostAuthSignupBody](c, types.PostAuthSignupBodySchema)
	if err != nil {
		return err
	}

	user, err := h.db.CreateUser(c.Request().Context(), body.Username, body.Password)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.PostAuthSignup{
		Id:       user.Id,
		Username: user.Username,
	}))
}

func (h *Handlers) HandlePostSignin(c echo.Context) error {
	body, err := Body[types.PostAuthSigninBody](c, types.PostAuthSigninBodySchema)
	if err != nil {
		return err
	}

	user, err := h.db.GetUserByUsername(c.Request().Context(), body.Username)
	if err != nil {
		return err
	}

	if user.Password != body.Password {
		return types.ErrIncorrectCreds
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id,
		"iat":    time.Now().Unix(),
		// "exp":    time.Now().Add(1000 * time.Second).Unix(),
	})

	tokenString, err := token.SignedString(([]byte)(config.Current.JwtSecret))
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.PostAuthSignin{
		Token: tokenString,
	}))
}

func (h *Handlers) User(c echo.Context) (*database.User, error) {
	authHeader := c.Request().Header.Get("Authorization")
	tokenString, err := utils.ParseAuthHeader(authHeader)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Current.JwtSecret), nil
	})

	if err != nil {
		return nil, types.ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := h.jwtValidator.Validate(token.Claims); err != nil {
			return nil, types.ErrInvalidToken
		}

		userId := claims["userId"].(string)
		user, err := h.db.GetUserById(c.Request().Context(), userId)
		if err != nil {
			return nil, types.ErrInvalidToken
		}

		return &user, nil
	}

	return nil, types.ErrInvalidToken
}

func (h *Handlers) HandleGetMe(c echo.Context) error {
	user, err := h.User(c)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.GetAuthMe{
		Id:       user.Id,
		Username: user.Username,
	}))
}

func (h *Handlers) InstallAuthHandlers(group Group) {
	group.POST(
		"Signup", "/auth/signup", 
		h.HandlePostSignup, 
		types.PostAuthSignup{}, types.PostAuthSignupBody{},
	)

	group.POST(
		"Signin", "/auth/signin", 
		h.HandlePostSignin, 
		types.PostAuthSignin{}, types.PostAuthSigninBody{},
	)

	group.GET(
		"GetMe", "/auth/me", 
		h.HandleGetMe, 
		types.GetAuthMe{}, nil,
	)
}
