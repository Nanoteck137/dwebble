package apis

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

// TODO(patrik): Change name?
type MediaResource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type MediaItem struct {
	Track   MediaResource   `json:"track"`
	Artists []MediaResource `json:"artists"`
	Album   MediaResource   `json:"album"`

	CoverArt types.Images `json:"coverArt"`

	MediaType types.MediaType `json:"mediaType"`
	MediaUrl  *string         `json:"mediaUrl"`
}

type GetMedia struct {
	Items []MediaItem `json:"items"`
}

type GetMediaCommonBody struct {
	MediaType types.MediaType `json:"mediaType,omitempty"`

	Shuffle bool   `json:"shuffle,omitempty"`
	Sort    string `json:"sort,omitempty"`

	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type GetMediaFromPlaylistBody struct {
	GetMediaCommonBody
}

type GetMediaFromTaglistBody struct {
	GetMediaCommonBody
}

type GetMediaFromFilterBody struct {
	GetMediaCommonBody

	Filter string `json:"filter"`
}

type GetMediaFromArtistBody struct {
	GetMediaCommonBody
}

type GetMediaFromAlbumBody struct {
	GetMediaCommonBody
}

type GetMediaFromIdsBody struct {
	GetMediaCommonBody

	TrackIds  []string `json:"trackIds"`
	KeepOrder bool     `json:"keepOrder,omitempty"`
}

func packMediaResult(c pyrin.Context, tracks []database.Track, mediaType types.MediaType, shuffle bool) (GetMedia, error) {
	if shuffle {
		rand.Shuffle(len(tracks), func(i, j int) {
			tracks[i], tracks[j] = tracks[j], tracks[i]
		})
	}

	res := GetMedia{
		Items: make([]MediaItem, len(tracks)),
	}

	for i, track := range tracks {
		artists := make([]MediaResource, len(track.FeaturingArtists)+1)

		artists[0] = MediaResource{
			Id:   track.ArtistId,
			Name: track.ArtistName,
		}

		for i, v := range track.FeaturingArtists {
			artists[i+1] = MediaResource{
				Id:   v.Id,
				Name: v.Name,
			}
		}

		ext := ".unknown"
		if e, ok := track.MediaType.ToExt(); ok {
			ext = e
		}

		mediaUrl := ConvertURL(c, fmt.Sprintf("/files/tracks/%s/track%s", track.Id, ext))

		res.Items[i] = MediaItem{
			Track: MediaResource{
				Id:   track.Id,
				Name: track.Name,
			},
			Artists: artists,
			Album: MediaResource{
				Id:   track.AlbumId,
				Name: track.AlbumName,
			},
			CoverArt:  ConvertAlbumCoverURL(c, track.AlbumId, track.AlbumCoverArt),
			MediaType: mediaType,
			MediaUrl:  &mediaUrl,
		}
	}

	return res, nil
}

func InstallMediaHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetMediaFromPlaylist",
			Method:       http.MethodPost,
			Path:         "/media/playlist/:playlistId",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromPlaylistBody{},
			Errors:       []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("playlistId")

				ctx := context.TODO()

				body, err := pyrin.Body[GetMediaFromPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				subquery := database.PlaylistTrackSubquery(playlist.Id)
				tracks, err := app.DB().GetTracksIn(ctx, subquery, "")
				if err != nil {
					return nil, err
				}

				return packMediaResult(c, tracks, body.MediaType, body.Shuffle)
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromTaglist",
			Method:       http.MethodPost,
			Path:         "/media/taglist/:taglistId",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromTaglistBody{},
			Errors:       []pyrin.ErrorType{ErrTypeTaglistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				taglistId := c.Param("taglistId")

				ctx := context.TODO()

				body, err := pyrin.Body[GetMediaFromTaglistBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				taglist, err := app.DB().GetTaglistById(ctx, taglistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TaglistNotFound()
					}

					return nil, err
				}

				if taglist.OwnerId != user.Id {
					return nil, TaglistNotFound()
				}

				tracks, err := app.DB().GetAllTracks(ctx, taglist.Filter, "")
				if err != nil {
					return nil, err
				}

				return packMediaResult(c, tracks, body.MediaType, body.Shuffle)
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromFilter",
			Method:       http.MethodPost,
			Path:         "/media/filter",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromFilterBody{},
			Errors:       []pyrin.ErrorType{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				body, err := pyrin.Body[GetMediaFromFilterBody](c)
				if err != nil {
					return nil, err
				}

				tracks, err := app.DB().GetAllTracks(ctx, body.Filter, "")
				if err != nil {
					return nil, err
				}

				return packMediaResult(c, tracks, body.MediaType, body.Shuffle)
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromArtist",
			Method:       http.MethodPost,
			Path:         "/media/artist/:artistId",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromArtistBody{},
			Errors:       []pyrin.ErrorType{ErrTypeArtistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				artistId := c.Param("artistId")

				ctx := context.TODO()

				body, err := pyrin.Body[GetMediaFromArtistBody](c)
				if err != nil {
					return nil, err
				}

				artist, err := app.DB().GetArtistById(ctx, artistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ArtistNotFound()
					}

					return nil, err
				}

				subquery := database.ArtistTrackSubquery(artist.Id)
				tracks, err := app.DB().GetTracksIn(ctx, subquery, "")
				if err != nil {
					return nil, err
				}

				return packMediaResult(c, tracks, body.MediaType, body.Shuffle)
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromAlbum",
			Method:       http.MethodPost,
			Path:         "/media/album/:albumId",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromAlbumBody{},
			Errors:       []pyrin.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				albumId := c.Param("albumId")

				ctx := context.TODO()

				body, err := pyrin.Body[GetMediaFromAlbumBody](c)
				if err != nil {
					return nil, err
				}

				album, err := app.DB().GetAlbumById(ctx, albumId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				sort := body.Sort
				if sort == "" {
					sort = "sort=number,name"
				}

				subquery := database.AlbumTrackSubquery(album.Id)
				tracks, err := app.DB().GetTracksIn(ctx, subquery, sort)
				if err != nil {
					return nil, err
				}

				return packMediaResult(c, tracks, body.MediaType, body.Shuffle)
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaFromIds",
			Method:       http.MethodPost,
			Path:         "/media/ids",
			ResponseType: GetMedia{},
			BodyType:     GetMediaFromIdsBody{},
			Errors:       []pyrin.ErrorType{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				body, err := pyrin.Body[GetMediaFromIdsBody](c)
				if err != nil {
					return nil, err
				}

				tracks, err := app.DB().GetTracksIn(ctx, body.TrackIds, "")
				if err != nil {
					return nil, err
				}

				if body.KeepOrder {
					trackMap := make(map[string]database.Track)
					for _, t := range tracks {
						trackMap[t.Id] = t
					}

					tracks = make([]database.Track, 0, len(body.TrackIds))
					for _, v := range body.TrackIds {
						track, exists := trackMap[v]
						if !exists {
							continue
						}

						tracks = append(tracks, track)
					}
				}

				return packMediaResult(c, tracks, body.MediaType, body.Shuffle)
			},
		},
	)
}
