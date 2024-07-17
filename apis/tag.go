package apis

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/handlers"
	"github.com/nanoteck137/dwebble/types"
)

type tagApi struct {
	app core.App
}

func (api *tagApi) HandleGetTags(c echo.Context) error {
	tags, err := api.app.DB().GetAllTags(c.Request().Context())
	if err != nil {
		return err
	}

	res := types.GetTags{
		Tags: make([]types.Tag, len(tags)),
	}

	for i, tag := range tags {
		// TODO(patrik): Fix this
		res.Tags[i] = types.Tag{
			Id:   tag.Name,
			Name: tag.Name,
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func InstallTagHandlers(app core.App, group handlers.Group) {
	api := tagApi{app: app}

	requireSetup := RequireSetup(app)

	group.Register(
		handlers.Handler{
			Name:        "GetTags",
			Path:        "/tags",
			Method:      http.MethodGet,
			DataType:    types.GetTags{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetTags,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},
	)
}
