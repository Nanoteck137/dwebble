package apis

import (
	"context"
	"errors"
	"net/http"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/pyrin"
)

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
			Name:         "GetQueue",
			Method:       http.MethodGet,
			Path:         "/queue/:playerId",
			ResponseType: nil,
			BodyType:     nil,
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
							Id:      playerId,
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
					queue, err = app.DB().GetQueue(ctx, defaultQueue.QueueId)
					if err != nil {
						return nil, err
					}
				}

				pretty.Println(queue)

				return nil, nil
			},
		},
	)
}
