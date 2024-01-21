package handlers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nanoteck137/dwebble/database"
)

type ApiConfig struct {
	workDir  string
	validate *validator.Validate

	db      *pgxpool.Pool
	queries *database.Queries
}

func New(db *pgxpool.Pool, queries *database.Queries, workDir string) ApiConfig {
	var validate = validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return ApiConfig{
		workDir:  workDir,
		queries:  queries,
		db:       db,
		validate: validate,
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

func ConvertURL(c *fiber.Ctx, path string) string {
	return c.BaseURL() + path
}


func (apiConfig *ApiConfig) HandlerCreateQueueFromAlbum(c *fiber.Ctx) error {
	body := struct {
		AlbumId string `json:"albumId"`
	}{}

	fmt.Println("Hello")

	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	tracks, err := apiConfig.queries.GetTracksByAlbum(c.UserContext(), body.AlbumId)
	if err != nil {
		return err
	}

	res := struct {
		Queue []database.GetTracksByAlbumRow `json:"queue"`
	}{
		Queue: tracks,
	}

	return c.JSON(res)
}
