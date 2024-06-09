package handlers

import (
	"strings"

	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
)

func parseQueryParam(param string, item func(action database.FilterAction, name string)) error {
	splits := strings.Split(param, ",")
	for _, s := range splits {
		s = strings.TrimSpace(s)

		if s == "" {
			continue
		}

		switch s[0] {
		case '+':
			item(database.FilterActionAdd, s[1:])
		case '-':
			item(database.FilterActionRemove, s[1:])
		default:
			item(database.FilterActionAdd, s)
		}
	}

	return nil
}

func (h *Handlers) getTagFilter(c echo.Context, filters []database.Filter) ([]database.Filter, error) {
	tagsParam := c.QueryParam("tags")

	tags, err := h.db.GetAllTags(c.Request().Context())
	if err != nil {
		return filters, err
	}

	tagToId := func(tag string) string {
		for _, t := range tags {
			if t.Name == tag {
				return t.Id
			}
		}

		return ""
	}

	var tagsToAdd []string
	var tagsToRemove []string

	parseQueryParam(tagsParam, func(action database.FilterAction, name string) {
		id := tagToId(name)
		if id == "" {
			return
		}

		switch action {
		case database.FilterActionAdd:
			tagsToAdd = append(tagsToAdd, id)
		case database.FilterActionRemove:
			tagsToRemove = append(tagsToRemove, id)
		}
	})

	if tagsToAdd != nil && len(tagsToAdd) > 0 {
		filters = append(filters, database.Filter{
			Action: database.FilterActionAdd,
			Type:   database.FilterTypeTags,
			Ids:    tagsToAdd,
		})
	}

	if tagsToRemove != nil && len(tagsToRemove) > 0 {
		filters = append(filters, database.Filter{
			Action: database.FilterActionRemove,
			Type:   database.FilterTypeTags,
			Ids:    tagsToRemove,
		})
	}

	return filters, nil
}

func (h *Handlers) getGenreFilter(c echo.Context, filters []database.Filter) ([]database.Filter, error) {
	param := c.QueryParam("genres")

	genres, err := h.db.GetAllGenres(c.Request().Context())
	if err != nil {
		return filters, err
	}

	genreToId := func(genre string) string {
		for _, g := range genres {
			if g.Name == genre {
				return g.Id
			}
		}

		return ""
	}

	var genresToAdd []string
	var genresToRemove []string

	parseQueryParam(param, func(action database.FilterAction, name string) {
		id := genreToId(name)
		if id == "" {
			return
		}

		switch action {
		case database.FilterActionAdd:
			genresToAdd = append(genresToAdd, id)
		case database.FilterActionRemove:
			genresToRemove = append(genresToRemove, id)
		}
	})

	if genresToAdd != nil && len(genresToAdd) > 0 {
		filters = append(filters, database.Filter{
			Action: database.FilterActionAdd,
			Type:   database.FilterTypeGenres,
			Ids:    genresToAdd,
		})
	}

	if genresToRemove != nil && len(genresToRemove) > 0 {
		filters = append(filters, database.Filter{
			Action: database.FilterActionRemove,
			Type:   database.FilterTypeGenres,
			Ids:    genresToRemove,
		})
	}

	return filters, nil
}

func (h *Handlers) HandleGetTracks(c echo.Context) error {
	var err error
	var filters []database.Filter

	filters, err = h.getTagFilter(c, filters)
	if err != nil {
		return err
	}

	filters, err = h.getGenreFilter(c, filters)
	if err != nil {
		return err
	}

	pretty.Println(filters)

	tracks, err := h.db.GetAllTracks(c.Request().Context(), false, filters)
	if err != nil {
		return err
	}

	res := types.GetTracks{
		Tracks: make([]types.Track, len(tracks)),
	}

	for i, track := range tracks {
		res.Tracks[i] = types.Track{
			Id:                track.Id,
			Number:            track.Number,
			Name:              track.Name,
			CoverArt:          ConvertTrackCoverURL(c, track.CoverArt),
			Duration:          track.Duration,
			BestQualityFile:   ConvertURL(c, "/tracks/original/"+track.BestQualityFile),
			MobileQualityFile: ConvertURL(c, "/tracks/mobile/"+track.MobileQualityFile),
			AlbumId:           track.AlbumId,
			ArtistId:          track.ArtistId,
			AlbumName:         track.AlbumName,
			ArtistName:        track.ArtistName,
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (h *Handlers) HandleGetTrackById(c echo.Context) error {
	id := c.Param("id")
	track, err := h.db.GetTrackById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.GetTrackById{
		Track: types.Track{
			Id:                track.Id,
			Number:            track.Number,
			Name:              track.Name,
			CoverArt:          ConvertTrackCoverURL(c, track.CoverArt),
			Duration:          track.Duration,
			BestQualityFile:   ConvertURL(c, "/tracks/original/"+track.BestQualityFile),
			MobileQualityFile: ConvertURL(c, "/tracks/mobile/"+track.MobileQualityFile),
			AlbumId:           track.AlbumId,
			ArtistId:          track.ArtistId,
		},
	}))
}

func (h *Handlers) InstallTrackHandlers(group *echo.Group) {
	group.GET("/tracks", h.HandleGetTracks)
	group.GET("/tracks/:id", h.HandleGetTrackById)
}
