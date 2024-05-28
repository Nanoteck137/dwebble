package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"github.com/doug-martin/goqu/v9"
	"github.com/kr/pretty"
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
		items[i].ItemIndex = i
	}

	pretty.Println(items)

	sort.Slice(items, func(i, j int) bool {
		return items[i].ItemIndex > items[j].ItemIndex
	})

	pretty.Println(items)

	maxIndex := 0
	for _, item := range items {
		maxIndex = max(maxIndex, item.ItemIndex)
	}

	nextIndex := maxIndex + 1

	fmt.Printf("nextIndex: %v\n", nextIndex)

	// ds := dialect.Insert("playlist_items").
	// 	Rows(goqu.Record{
	// 		"playlist_id": playlistId,
	// 		"track_id":    trackIds[0],
	// 		"item_index":  0,
	// 	}).
	// 	Prepared(true)
	//
	// _, err := db.Exec(ctx, ds)
	// if err != nil {
	// 	return err
	// }

	return nil
}
