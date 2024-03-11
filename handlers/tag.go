package handlers

import (
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
		res.Tags[i] = types.Tag{
			Id:   tag.Id,
			Name: tag.Name,
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (h *Handlers) InstallTagHandlers(group *echo.Group) {
	group.GET("/tags", h.HandleGetTags)
}
