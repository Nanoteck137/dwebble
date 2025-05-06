package apis

import (
	"errors"
	"net/http"
	"strings"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

type Album struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Year *int64 `json:"year"`

	CoverArt types.Images `json:"coverArt"`

	Artists []ArtistInfo `json:"artists"`

	Tags []string `json:"tags"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

func ConvertDBAlbum(c pyrin.Context, album database.Album) Album {
	allArtists := make([]ArtistInfo, len(album.FeaturingArtists)+1)

	allArtists[0] = ArtistInfo{
		Id:   album.ArtistId,
		Name: album.ArtistName,
	}

	for i, v := range album.FeaturingArtists {
		allArtists[i+1] = ArtistInfo{
			Id:   v.Id,
			Name: v.Name,
		}
	}
	return Album{
		Id:       album.Id,
		Name:     album.Name,
		Year:     ConvertSqlNullInt64(album.Year),
		CoverArt: ConvertAlbumCoverURL(c, album.Id, album.CoverArt),
		Artists:  allArtists,
		Tags:     utils.SplitString(album.Tags.String),
		Created:  album.Created,
		Updated:  album.Updated,
	}
}

type GetAlbums struct {
	Page   types.Page `json:"page"`
	Albums []Album    `json:"albums"`
}

type GetAlbumById struct {
	Album
}

type GetAlbumTracks struct {
	Tracks []Track `json:"tracks"`
}

func InstallAlbumHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetAlbums",
			Path:         "/albums",
			Method:       http.MethodGet,
			ResponseType: GetAlbums{},
			Errors:       []pyrin.ErrorType{ErrTypeInvalidFilter, ErrTypeInvalidSort},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				opts := getPageOptions(q)

				albums, pageInfo, err := app.DB().GetAlbumsPaged(c.Request().Context(), opts)
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					if errors.Is(err, database.ErrInvalidSort) {
						return nil, InvalidSort(err)
					}

					return nil, err
				}

				res := GetAlbums{
					Page:   pageInfo,
					Albums: make([]Album, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = ConvertDBAlbum(c, album)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SearchAlbums",
			Path:         "/albums/search",
			Method:       http.MethodGet,
			ResponseType: GetAlbums{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				query := strings.TrimSpace(q.Get("query"))

				albums, err := app.DB().SearchAlbums(query)
				if err != nil {
					return nil, err
				}

				res := GetAlbums{
					Albums: make([]Album, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = ConvertDBAlbum(c, album)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetAlbumById",
			Method:       http.MethodGet,
			Path:         "/albums/:id",
			ResponseType: GetAlbumById{},
			Errors:       []pyrin.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				album, err := app.DB().GetAlbumById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				return GetAlbumById{
					Album: ConvertDBAlbum(c, album),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetAlbumTracks",
			Method:       http.MethodGet,
			Path:         "/albums/:id/tracks",
			ResponseType: GetAlbumTracks{},
			Errors:       []pyrin.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				album, err := app.DB().GetAlbumById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				tracks, err := app.DB().GetTracksByAlbum(c.Request().Context(), album.Id)
				if err != nil {
					return nil, err
				}

				res := GetAlbumTracks{
					Tracks: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},
	)
}
