package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/types"
)

func (h *Handlers) HandleGetTags(c echo.Context) error {
	tags, err := h.db.GetAllTags(c.Request().Context())
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

func (h *Handlers) InstallTagHandlers(group Group) {
	group.Register(
		Handler{
			Name:        "GetTags",
			Path:        "/tags",
			Method:      http.MethodGet,
			DataType:    types.GetTags{},
			BodyType:    nil,
			HandlerFunc: h.HandleGetTags,
		},
	)
}
