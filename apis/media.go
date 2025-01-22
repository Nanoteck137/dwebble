package apis

import (
	"context"
	"errors"
	"fmt"
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
	MediaUrl  string          `json:"mediaUrl"`
}

type GetMedia struct {
	Items []MediaItem `json:"items"`
}

type GetMediaCommonBody struct {
	MediaType types.MediaType `json:"type,omitempty"`

	Shuffle bool   `json:"shuffle,omitempty"`
	Sort    string `json:"sort,omitempty"`

	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type GetMediaFromPlaylistBody struct {
	GetMediaCommonBody
}

type GetMediaFromAlbumBody struct {
	GetMediaCommonBody
}

func InstallMediaHandlers(app core.App, group pyrin.Group) {
	// TODO(patrik):
	// - [ ] Playlist
	// - [ ] Taglist
	// - [ ] Artist
	// - [ ] Filter
	// - [ ] Album
	// - [ ] Ids

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
				_ = playlistId

				return nil, nil
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

				album, err := app.DB().GetAlbumById(ctx, albumId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				subquery := database.AlbumTrackSubquery(album.Id)
				tracks, err := app.DB().GetTracksIn(ctx, subquery)
				if err != nil {
					return nil, err
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

					mediaItems := *track.Formats.Get()

					// TODO(patrik): Better selection algo
					found := false
					var originalItem database.TrackFormat
					for _, item := range mediaItems {
						if item.IsOriginal {
							originalItem = item
							found = true
							break
						}
					}

					if !found {
						// TODO(patrik): Better error handling
						return nil, errors.New("Contains bad media items")
					}

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
						MediaType: originalItem.MediaType,
						MediaUrl:  ConvertURL(c, fmt.Sprintf("/files/tracks/%s/media/%s/%s", track.Id, originalItem.Id, originalItem.Filename)),
					}
				}

				return res, nil
			},
		},
	)
}
