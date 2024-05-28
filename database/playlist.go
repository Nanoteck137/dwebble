package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
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
		return Playlist{}, err
	}

	return item, nil
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
