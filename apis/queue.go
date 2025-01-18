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
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
)

const DEFAULT_QUEUE_NAME = "Default Queue"

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

type AddToQueuePlaylistBody struct {
	AddToQueue

	PlaylistId string `json:"playlistId"`
}

type AddToQueueTaglistBody struct {
	AddToQueue

	TaglistId string `json:"taglistId"`
}

type AddToQueueAlbumBody struct {
	AddToQueue

	AlbumId string `json:"albumId"`
}

func getNextIndex(ctx context.Context, db *database.Database, queueId string) (int, error) {
	index, hasIndex, err := db.GetLastQueueItemIndex(ctx, queueId)
	if err != nil {
		return 0, err
	}

	if hasIndex {
		return index + 1, nil
	}

	return 0, nil
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
	// - GetDefaultQueue
	// - ClearQueue
	// - GetQueueItems
	// - RemoveFromQueue
	// - ShuffleQueue
	//
	// # Clear Flag for these
	// - AddFromPlaylist
	// - AddFromTaglist
	// - AddFromIds
	// - AddFromTracks
	// - AddFromAlbum
	// - AddFromArtist

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
							CoverArt: ConvertAlbumCoverURL(c, item.AlbumId, item.CoverArt),
							MediaUrl: ConvertURL(c, fmt.Sprintf("/files/tracks/%s/media/%s/%s", item.TrackId, originalItem.Id, originalItem.Filename)),
						},
						Created: item.Created,
						Updated: item.Updated,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddToQueueFromPlaylist",
			Method:   http.MethodPost,
			Path:     "/queue/:id/add/playlist",
			BodyType: AddToQueuePlaylistBody{},
			Errors:   []pyrin.ErrorType{ErrTypeQueueNotFound, ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[AddToQueuePlaylistBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				queue, err := app.DB().GetQueueById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, QueueNotFound()
					}

					return nil, err
				}

				if queue.UserId != user.Id {
					return nil, QueueNotFound()
				}

				playlist, err := app.DB().GetPlaylistById(ctx, body.PlaylistId)
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

				nextIndex, err := getNextIndex(ctx, app.DB(), queue.Id)
				if err != nil {
					return nil, err
				}

				for i, item := range items {
					_, err := app.DB().CreateQueueItem(ctx, database.CreateQueueItemParams{
						QueueId:     queue.Id,
						OrderNumber: nextIndex + i,
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
			Name:     "AddToQueueFromTaglist",
			Method:   http.MethodPost,
			Path:     "/queue/:id/add/taglist",
			BodyType: AddToQueueTaglistBody{},
			Errors:   []pyrin.ErrorType{ErrTypeQueueNotFound, ErrTypeTaglistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[AddToQueueTaglistBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				queue, err := app.DB().GetQueueById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, QueueNotFound()
					}

					return nil, err
				}

				if queue.UserId != user.Id {
					return nil, QueueNotFound()
				}

				taglist, err := app.DB().GetTaglistById(ctx, body.TaglistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TaglistNotFound()
					}

					return nil, err
				}

				if taglist.OwnerId != user.Id {
					return nil, TaglistNotFound()
				}

				items, err := app.DB().GetAllTracks(ctx, taglist.Filter, "")
				if err != nil {
					return nil, err
				}

				if body.Shuffle {
					rand.Shuffle(len(items), func(i, j int) {
						items[i], items[j] = items[j], items[i]
					})
				}

				nextIndex, err := getNextIndex(ctx, app.DB(), queue.Id)
				if err != nil {
					return nil, err
				}

				for i, item := range items {
					_, err := app.DB().CreateQueueItem(ctx, database.CreateQueueItemParams{
						QueueId:     queue.Id,
						OrderNumber: nextIndex + i,
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
			Name:   "AddToQueueFromAlbum",
			Method: http.MethodPost,
			Path:   "/queue/:id/add/album",
			BodyType: AddToQueueAlbumBody{},
			Errors: []pyrin.ErrorType{ErrTypeQueueNotFound, ErrTypeAlbumNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[AddToQueueAlbumBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				queue, err := app.DB().GetQueueById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, QueueNotFound()
					}

					return nil, err
				}

				if queue.UserId != user.Id {
					return nil, QueueNotFound()
				}

				album, err := app.DB().GetAlbumById(ctx, body.AlbumId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				tracks, err := app.DB().GetTracksByAlbum(ctx, album.Id)
				if err != nil {
					return nil, err
				}

				nextIndex, err := getNextIndex(ctx, app.DB(), queue.Id)
				if err != nil {
					return nil, err
				}

				for i, track := range tracks {
					_, err := app.DB().CreateQueueItem(ctx, database.CreateQueueItemParams{
						QueueId:     queue.Id,
						OrderNumber: nextIndex + i,
						TrackId:     track.Id,
					})
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},	
	)
}
