package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
)

type Playlist struct {
	Id      string
	Name    string
	OwnerId string
}

func (db *Database) GetPlaylistsByUser(ctx context.Context, userId string) ([]Playlist, error) {
	ds := dialect.From("playlists").
		Select("id", "name", "owner_id").
		Where(goqu.I("owner_id").Eq(userId)).
		Prepared(true)

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Playlist
	for rows.Next() {
		var item Playlist
		err := rows.Scan(&item.Id, &item.Name, &item.OwnerId)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetPlaylistById(ctx context.Context, id string) (Playlist, error) {
	ds := dialect.From("playlists").
		Select("id", "name", "owner_id").
		Where(goqu.I("id").Eq(id)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Playlist{}, err
	}

	var item Playlist
	err = row.Scan(&item.Id, &item.Name, &item.OwnerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Playlist{}, types.ErrNoPlaylist
		}

		return Playlist{}, err
	}

	return item, nil
}

type PlaylistItem struct {
	PlaylistId string
	TrackId    string
	ItemIndex  int
}

func (db *Database) GetPlaylistItems(ctx context.Context, id string) ([]PlaylistItem, error) {
	ds := dialect.From("playlist_items").
		Select("playlist_id", "track_id", "item_index").
		Where(goqu.I("playlist_id").Eq(id)).
		Order(goqu.I("item_index").Asc()).
		Prepared(true)

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PlaylistItem
	for rows.Next() {
		var item PlaylistItem
		err := rows.Scan(&item.PlaylistId, &item.TrackId, &item.ItemIndex)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

type CreatePlaylistParams struct {
	Name    string
	OwnerId string
}

func (db *Database) CreatePlaylist(ctx context.Context, params CreatePlaylistParams) (Playlist, error) {
	ds := dialect.Insert("playlists").
		Rows(goqu.Record{
			"id":       utils.CreateId(),
			"name":     params.Name,
			"owner_id": params.OwnerId,
		}).
		Returning("id", "name", "owner_id").
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Playlist{}, nil
	}

	var item Playlist
	err = row.Scan(&item.Id, &item.Name, &item.OwnerId)
	if err != nil {
		return Playlist{}, nil
	}

	return item, nil
}

func (db *Database) AddItemsToPlaylist(ctx context.Context, playlistId string, trackIds []string) error {
	items, err := db.GetPlaylistItems(ctx, playlistId)
	if err != nil {
		return err
	}

	for i := range items {
		items[i].ItemIndex = i + 1
	}

	maxIndex := 0
	for _, item := range items {
		maxIndex = max(maxIndex, item.ItemIndex)
	}

	nextIndex := maxIndex + 1

	fmt.Printf("nextIndex: %v\n", nextIndex)

	for i, trackId := range trackIds {
		ds := dialect.Insert("playlist_items").
			Rows(goqu.Record{
				"playlist_id": playlistId,
				"track_id":    trackId,
				"item_index":  nextIndex + i,
			}).
			Prepared(true)

		_, err = db.Exec(ctx, ds)
		if err != nil {
			if err, ok := err.(sqlite3.Error); ok {
				if errors.Is(err.ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
					// TODO(patrik): Move and fix
					return types.NewApiError(400, "Track already in playlist")
				}
			}
			return err
		}
	}

	for _, item := range items {
		ds := dialect.Update("playlist_items").
			Set(goqu.Record{
				"item_index": item.ItemIndex,
			}).
			Where(
				goqu.And(
					goqu.I("playlist_id").Eq(item.PlaylistId),
					goqu.I("track_id").Eq(item.TrackId),
				),
			).
			Prepared(true)

		_, err := db.Exec(ctx, ds)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO(patrik): Change from trackIndices to trackIds
func (db *Database) DeleteItemsFromPlaylist(ctx context.Context, playlistId string, trackIndices []int) error {
	for _, trackIndex := range trackIndices {
		ds := dialect.Delete("playlist_items").
			Where(
				goqu.And(
					goqu.I("playlist_id").Eq(playlistId),
					goqu.I("item_index").Eq(trackIndex),
				),
			)

		_, err := db.Exec(ctx, ds)
		if err != nil {
			return nil
		}
	}

	items, err := db.GetPlaylistItems(ctx, playlistId)
	if err != nil {
		return err
	}

	for i := range items {
		items[i].ItemIndex = i + 1
	}

	for _, item := range items {
		ds := dialect.Update("playlist_items").
			Set(goqu.Record{
				"item_index": item.ItemIndex,
			}).
			Where(
				goqu.And(
					goqu.I("playlist_id").Eq(item.PlaylistId),
					goqu.I("track_id").Eq(item.TrackId),
				),
			).
			Prepared(true)

		_, err := db.Exec(ctx, ds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) MovePlaylistItem(ctx context.Context, playlistId string, itemIndex int, toIndex int) error {
	items, err := db.GetPlaylistItems(ctx, playlistId)
	if err != nil {
		return err
	}

	// TODO(patrik): Maybe we should try to find the items inside the items
	// array instead
	itemIndex = itemIndex - 1
	toIndex = toIndex - 1

	item := items[itemIndex]

	// TODO(patrik): This need testing
	length := (toIndex - itemIndex) + 1
	for i := itemIndex; i < length; i++ {
		if i < len(items)-1 {
			items[i] = items[i+1]
		}
	}

	items[toIndex] = item

	for i := range items {
		items[i].ItemIndex = i + 1
	}

	for _, item := range items {
		ds := dialect.Update("playlist_items").
			Set(goqu.Record{
				"item_index": item.ItemIndex,
			}).
			Where(
				goqu.And(
					goqu.I("playlist_id").Eq(item.PlaylistId),
					goqu.I("track_id").Eq(item.TrackId),
				),
			).
			Prepared(true)

		_, err := db.Exec(ctx, ds)
		if err != nil {
			return err
		}
	}

	return nil
}
