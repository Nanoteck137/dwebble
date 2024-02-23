package handlers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
)

type ApiConfig struct {
	workDir  types.WorkDir
	validate *validator.Validate
	db       *database.Database
}

func New(db *database.Database, workDir types.WorkDir) *ApiConfig {
	var validate = validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &ApiConfig{
		workDir:  workDir,
		validate: validate,
		db:       db,
	}
}

func (api *ApiConfig) validateBody(body any) map[string]string {
	err := api.validate.Struct(body)
	if err != nil {
		type ValidationError struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}

		validationErrs := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			validationErrs[field] = fmt.Sprintf("'%v' not satisfying tags '%v'", field, err.Tag())
		}

		return validationErrs
	}

	return nil
}

func ConvertURL(c echo.Context, path string) string {
	host := c.Request().Host

	return fmt.Sprintf("http://%s%s", host, path)
}

// func (apiConfig *ApiConfig) HandlerCreateQueueFromAlbum(c *fiber.Ctx) error {
// 	body := struct {
// 		AlbumId string `json:"albumId"`
// 	}{}
//
// 	fmt.Println("Hello")
//
// 	err := c.BodyParser(&body)
// 	if err != nil {
// 		return err
// 	}
//
// 	tracks, err := apiConfig.queries.GetTracksByAlbum(c.UserContext(), body.AlbumId)
// 	if err != nil {
// 		return err
// 	}
//
// 	res := struct {
// 		Queue []database.GetTracksByAlbumRow `json:"queue"`
// 	}{
// 		Queue: tracks,
// 	}
//
// 	return c.JSON(res)
// }
