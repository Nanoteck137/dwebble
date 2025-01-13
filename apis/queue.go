package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/pyrin"
)

type MusicTrackArtist struct {
	ArtistId   string
	ArtistName string
}

type MusicTrackAlbum struct {
	AlbumId   string
	AlbumName string
}

type MusicTrack struct {
	Name string

	Artists []MusicTrackArtist
	Album   MusicTrackAlbum

	CoverArtUrl string
	MediaUrl    string
}

type QueueItem struct {
	Id      string
	QueueId string

	OrderNumber int
	TrackId     string

	Track MusicTrack

	Created int64
	Updated int64
}

type Queue struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetQueueItems struct {
	Items []QueueItem `json:"items"`
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

				pretty.Println(queue)

				items, err := app.DB().GetQueueItemsForPlay(ctx, queue.Id)
				if err != nil {
					return nil, err
				}

				pretty.Println(items)

				// items, err := app.DB().GetAllQueueItemIds(ctx, queue.Id)
				// if err != nil {
				// 	return nil, err
				// }
				//
				// pretty.Println(items)
				//
				// tracks, err := app.DB().GetTracksByIds(ctx, items)
				// if err != nil {
				// 	return nil, err
				// }
				//
				// pretty.Println(tracks)

				res := GetQueueItems{
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

					pretty.Println(originalItem)

					res.Items[i] = QueueItem{
						Id:          item.Id,
						QueueId:     item.QueueId,
						OrderNumber: item.OrderNumber,
						TrackId:     item.TrackId,
						Track: MusicTrack{
							Name:        item.Name,
							Artists:     []MusicTrackArtist{},
							Album:       MusicTrackAlbum{},
							CoverArtUrl: "",
							MediaUrl:    utils.ConvertURL(c, fmt.Sprintf("/files/tracks/%s/media/%s/%s", item.TrackId, originalItem.Id, originalItem.Filename)),
						},
						Created: item.Created,
						Updated: item.Updated,
					}
				}

				return res, nil
			},
		},
	)
}
