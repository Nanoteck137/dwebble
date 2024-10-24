package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/dwebble/tools/utils"
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
			return Playlist{}, ErrItemNotFound
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

func (db *Database) GetPlaylistItems(ctx context.Context, playlistId string) ([]PlaylistItem, error) {
	ds := dialect.From("playlist_items").
		Select("playlist_id", "track_id", "item_index").
		Where(goqu.I("playlist_id").Eq(playlistId)).
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

func (db *Database) GetPlaylistTracks(ctx context.Context, playlistId string) ([]Track, error) {
	tracks := TrackQuery().As("tracks")

	ds := dialect.From("playlist_items").
		Select("tracks.*").
		Join(tracks, goqu.On(goqu.I("tracks.id").Eq(goqu.I("playlist_items.track_id")))).
		Where(goqu.I("playlist_id").Eq(playlistId)).
		Order(goqu.I("item_index").Asc()).
		Prepared(true)

	var items []Track
	err := db.Select(&items, ds)
	if err != nil {
		return nil, err
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
					// TODO(patrik): Fix
					return errors.New("Track already in playlist")
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

func (db *Database) AddItemsToPlaylistRaw(ctx context.Context, playlistId string, trackIds []string) error {
	for i, trackId := range trackIds {
		ds := dialect.Insert("playlist_items").
			Rows(goqu.Record{
				"playlist_id": playlistId,
				"track_id":    trackId,
				"item_index":  i,
			}).
			Prepared(true)

		_, err := db.Exec(ctx, ds)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO(patrik): Use Begin
func (db *Database) DeleteItemsFromPlaylist(ctx context.Context, playlistId string, trackIds []string) error {
	for _, trackId := range trackIds {
		ds := dialect.Delete("playlist_items").
			Where(
				goqu.And(
					goqu.I("playlist_id").Eq(playlistId),
					goqu.I("track_id").Eq(trackId),
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

	for i, item := range items {
		ds := dialect.Update("playlist_items").
			Set(goqu.Record{
				"item_index": i + 1,
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

// TODO(patrik): Add bounds check for toIndex
func (db *Database) MovePlaylistItem(ctx context.Context, playlistId string, trackId string, toIndex int) error {
	items, err := db.GetPlaylistItems(ctx, playlistId)
	if err != nil {
		return err
	}

	itemIndex := -1

	for i, item := range items {
		if item.TrackId == trackId {
			itemIndex = i
		}
	}

	if itemIndex == -1 {
		// TODO(patrik): Fix error
		return errors.New("Track not in playlist")
	}

	// TODO(patrik): Maybe we should try to find the items inside the items
	// array instead
	toIndex = toIndex - 1

	item := items[itemIndex]

	if toIndex > itemIndex {
		// TODO(patrik): This need testing
		length := (toIndex - itemIndex) + 1
		for i := itemIndex; i < length; i++ {
			if i < len(items)-1 {
				items[i] = items[i+1]
			}
		}

		items[toIndex] = item

	} else {
		// TODO(patrik): This need testing
		length := itemIndex - toIndex
		for i := length + toIndex; i > 0; i-- {
			items[i] = items[i-1]
		}

		items[toIndex] = item
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
