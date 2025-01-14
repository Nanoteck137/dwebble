package apis

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
)

type MusicTrackArtist struct {
	// TODO(patrik): Change ArtistId, ArtistName to Id, Name?
	ArtistId   string `json:"artistId"`
	ArtistName string `json:"artistName"`
}

type MusicTrackAlbum struct {
	AlbumId   string `json:"albumId"`
	AlbumName string `json:"albumName"`
}

type MusicTrack struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Artists []MusicTrackArtist `json:"artists"`
	Album   MusicTrackAlbum    `json:"album"`

	CoverArt types.Images `json:"coverArt"`
	MediaUrl string       `json:"mediaUrl"`
}

type QueueItem struct {
	Id      string `json:"id"`
	QueueId string `json:"queueId"`

	OrderNumber int    `json:"orderNumber"`
	TrackId     string `json:"trackId"`

	Track MusicTrack `json:"track"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

type Queue struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetQueueItems struct {
	Index int         `json:"index"`
	Items []QueueItem `json:"items"`
}

type UpdateQueueBody struct {
	Name      *string `json:"name,omitempty"`
	ItemIndex *int    `json:"itemIndex,omitempty"`
}

func (b *UpdateQueueBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
}

func (b UpdateQueueBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
	)
}

type AddToQueue struct {
	Shuffle    bool `json:"shuffle,omitempty"`
	AddToFront bool `json:"addToFront,omitempty"`
}

func InstallQueueHandlers(app core.App, group pyrin.Group) {
	// NOTE(patrik):
	// 1. Multiple queues per user
	//   - Default Queue for players
	//   - Shared between players
	// 2. One queue per player
	//   - Players should be able to make an copy of the queue to their own queue
	//   - Create playlists from queue? It can't retain the order but can store the tracks
	//   - Not logged in queue

	// TODO(patrik):
	// - GetQueue
	// - CreateQueueFromPlaylist
	// - CreateQueueFromTaglist
	// - CreateQueueFromIds
	// - CreateQueueFromTracks
	// - CreateQueueFromAlbum
	// - CreateQueueFromArtist
	// - AddFromPlaylist
	//   - Clear Flag
	// - AddFromTaglist
	//   - Clear Flag
	// - AddFromIds
	//   - Clear Flag
	// - AddFromTracks
	//   - Clear Flag
	// - AddFromAlbum
	//   - Clear Flag
	// - AddFromArtist
	//   - Clear Flag
	// - ClearQueue
	// - AddToQueue
	// - RemoveFromQueue
	// - DeleteQueue

	group.Register(
		pyrin.ApiHandler{
			Name:         "GetDefaultQueue",
			Method:       http.MethodGet,
			Path:         "/queue/default/:playerId",
			ResponseType: Queue{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playerId := c.Param("playerId")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				_, err = app.DB().GetPlayerById(ctx, playerId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						err := app.DB().CreatePlayer(ctx, database.CreatePlayerParams{
							Id: playerId,
						})
						if err != nil {
							return nil, err
						}
					} else {
						return nil, err
					}
				}

				// TODO(patrik): Move
				const DEFAULT_QUEUE_NAME = "Default Queue"

				var queue database.Queue

				defaultQueue, err := app.DB().GetDefaultQueue(ctx, playerId, user.Id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						log.Info("Default queue not found")
						queue, err = app.DB().CreateQueue(ctx, database.CreateQueueParams{
							PlayerId: playerId,
							UserId:   user.Id,
							Name:     DEFAULT_QUEUE_NAME,
						})
						if err != nil {
							return nil, err
						}

						err = app.DB().CreateDefaultQueue(ctx, database.CreateDefaultQueueParams{
							PlayerId: playerId,
							UserId:   user.Id,
							QueueId:  queue.Id,
						})
						if err != nil {
							return nil, err
						}
					} else {
						return nil, err
					}
				} else {
					log.Info("Default queue found")
					queue, err = app.DB().GetQueueById(ctx, defaultQueue.QueueId)
					if err != nil {
						return nil, err
					}
				}

				pretty.Println(queue)

				return Queue{
					Id:   queue.Id,
					Name: queue.Name,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "ClearQueue",
			Method: http.MethodPost,
			Path:   "/queue/:id/clear",
			Errors: []pyrin.ErrorType{ErrTypeQueueNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				queue, err := app.DB().GetQueueById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, QueueNotFound()
					}

					return nil, err
				}

				pretty.Println(queue)

				err = app.DB().DeleteAllQueueItems(ctx, queue.Id)
				if err != nil {
					return nil, err
				}

				err = app.DB().UpdateQueue(ctx, queue.Id, database.QueueChanges{
					ItemIndex: types.Change[int]{
						Value:   0,
						Changed: queue.ItemIndex != 0,
					},
				})
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "AddToQueueFromAlbum",
			Method: http.MethodPost,
			Path:   "/queue/:id/add/album/:albumId",
			Errors: []pyrin.ErrorType{ErrTypeQueueNotFound, ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				albumId := c.Param("albumId")

				ctx := context.TODO()

				queue, err := app.DB().GetQueueById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, QueueNotFound()
					}

					return nil, err
				}

				album, err := app.DB().GetAlbumById(ctx, albumId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				pretty.Println(queue)
				pretty.Println(album)

				tracks, err := app.DB().GetTracksByAlbum(ctx, album.Id)
				if err != nil {
					return nil, err
				}

				// TODO(patrik): Get the last item index from queue
				for i, track := range tracks {
					_, err := app.DB().CreateQueueItem(ctx, database.CreateQueueItemParams{
						QueueId:     queue.Id,
						OrderNumber: i,
						TrackId:     track.Id,
					})
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddToQueueFromPlaylist",
			Method:   http.MethodPost,
			Path:     "/queue/:id/add/playlist/:playlistId",
			BodyType: AddToQueue{},
			Errors:   []pyrin.ErrorType{ErrTypeQueueNotFound, ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				playlistId := c.Param("playlistId")

				body, err := pyrin.Body[AddToQueue](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				queue, err := app.DB().GetQueueById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, QueueNotFound()
					}

					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				items, err := app.DB().GetPlaylistTracks(ctx, playlist.Id)
				if err != nil {
					return nil, err
				}

				if body.Shuffle {
					rand.Shuffle(len(items), func(i, j int) {
						items[i], items[j] = items[j], items[i]
					})
				}

				// TODO(patrik): Get the last item index from queue
				for i, item := range items {
					_, err := app.DB().CreateQueueItem(ctx, database.CreateQueueItemParams{
						QueueId:     queue.Id,
						OrderNumber: i,
						TrackId:     item.Id,
					})
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetQueueItems",
			Method:       http.MethodGet,
			Path:         "/queue/:id/items",
			ResponseType: GetQueueItems{},
			Errors:       []pyrin.ErrorType{ErrTypeQueueNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				queue, err := app.DB().GetQueueById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, QueueNotFound()
					}

					return nil, err
				}

				items, err := app.DB().GetQueueItemsForPlay(ctx, queue.Id)
				if err != nil {
					return nil, err
				}

				res := GetQueueItems{
					Index: queue.ItemIndex,
					Items: make([]QueueItem, len(items)),
				}

				for i, item := range items {
					mediaItems := *item.MediaItems.Get()

					// TODO(patrik): Better selection algo
					found := false
					var originalItem database.QueueTrackItemMediaItem
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

					res.Items[i] = QueueItem{
						Id:          item.Id,
						QueueId:     item.QueueId,
						OrderNumber: item.OrderNumber,
						TrackId:     item.TrackId,
						Track: MusicTrack{
							Id:   item.TrackId,
							Name: item.Name,
							Artists: []MusicTrackArtist{
								{
									ArtistId:   item.ArtistId,
									ArtistName: item.ArtistName,
								},
							},
							Album: MusicTrackAlbum{
								AlbumId:   item.AlbumId,
								AlbumName: item.AlbumName,
							},
							CoverArt: utils.ConvertAlbumCoverURL(c, item.AlbumId, item.CoverArt),
							MediaUrl: utils.ConvertURL(c, fmt.Sprintf("/files/tracks/%s/media/%s/%s", item.TrackId, originalItem.Id, originalItem.Filename)),
						},
						Created: item.Created,
						Updated: item.Updated,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "UpdateQueue",
			Method:   http.MethodPatch,
			Path:     "/queue/:id",
			BodyType: UpdateQueueBody{},
			Errors:   []pyrin.ErrorType{ErrTypeQueueNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[UpdateQueueBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				queue, err := app.DB().GetQueueById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, QueueNotFound()
					}

					return nil, err
				}

				changes := database.QueueChanges{}

				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != queue.Name,
					}
				}

				if body.ItemIndex != nil {
					changes.ItemIndex = types.Change[int]{
						Value:   *body.ItemIndex,
						Changed: *body.ItemIndex != queue.ItemIndex,
					}
				}

				err = app.DB().UpdateQueue(ctx, queue.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
