package apis

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

type Track struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Duration int64  `json:"duration"`
	Number   *int64 `json:"number"`
	Year     *int64 `json:"year"`

	CoverArt types.Images `json:"coverArt"`

	AlbumId   string `json:"albumId"`
	AlbumName string `json:"albumName"`

	Artists []ArtistInfo `json:"artists"`

	Tags []string `json:"tags"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

func ConvertDBTrack(c pyrin.Context, track database.Track) Track {
	artists := make([]ArtistInfo, len(track.FeaturingArtists)+1)

	artists[0] = ArtistInfo{
		Id:   track.ArtistId,
		Name: track.ArtistName,
	}

	for i, v := range track.FeaturingArtists {
		artists[i+1] = ArtistInfo{
			Id:   v.Id,
			Name: v.Name,
		}
	}

	return Track{
		Id:        track.Id,
		Name:      track.Name,
		Duration:  track.Duration,
		Number:    ConvertSqlNullInt64(track.Number),
		Year:      ConvertSqlNullInt64(track.Year),
		CoverArt:  ConvertAlbumCoverURL(c, track.AlbumId, track.AlbumCoverArt),
		AlbumId:   track.AlbumId,
		AlbumName: track.AlbumName,
		Artists:   artists,
		Tags:      utils.SplitString(track.Tags.String),
		Created:   track.Created,
		Updated:   track.Updated,
	}
}

type GetTracks struct {
	Page   types.Page `json:"page"`
	Tracks []Track    `json:"tracks"`
}

type GetTrackById struct {
	Track
}

// TODO(patrik): Move
func getPageOptions(q url.Values) database.FetchOptions {
	perPage := 100
	page := 0

	if s := q.Get("perPage"); s != "" {
		i, _ := strconv.Atoi(s)
		if i > 0 {
			perPage = i
		}
	}

	if s := q.Get("page"); s != "" {
		i, _ := strconv.Atoi(s)
		page = i
	}

	return database.FetchOptions{
		Filter:  q.Get("filter"),
		Sort:    q.Get("sort"),
		PerPage: perPage,
		Page:    page,
	}
}

func InstallTrackHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetTracks",
			Method:       http.MethodGet,
			Path:         "/tracks",
			ResponseType: GetTracks{},
			Errors:       []pyrin.ErrorType{ErrTypeInvalidFilter, ErrTypeInvalidSort},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				opts := getPageOptions(q)

				tracks, p, err := app.DB().GetPagedTracks(c.Request().Context(), opts)
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					if errors.Is(err, database.ErrInvalidSort) {
						return nil, InvalidSort(err)
					}

					return nil, err
				}

				res := GetTracks{
					Page:   p,
					Tracks: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SearchTracks",
			Method:       http.MethodGet,
			Path:         "/tracks/search",
			ResponseType: GetTracks{},
			Errors:       []pyrin.ErrorType{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				query := strings.TrimSpace(q.Get("query"))

				tracks, err := app.DB().SearchTracks(query)
				if err != nil {
					return nil, err
				}

				res := GetTracks{
					Tracks: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetTrackById",
			Method:       http.MethodGet,
			Path:         "/tracks/:id",
			ResponseType: GetTrackById{},
			Errors:       []pyrin.ErrorType{ErrTypeTrackNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				track, err := app.DB().GetTrackById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TrackNotFound()
					}

					return nil, err
				}

				return GetTrackById{
					Track: ConvertDBTrack(c, track),
				}, nil
			},
		},
	)
}
