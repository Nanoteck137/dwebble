package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/dwebble/tools/utils"
)

type Playlist struct {
	Id      string         `db:"id"`
	Name    string         `db:"name"`
	Picture sql.NullString `db:"picture"`

	OwnerId string `db:"owner_id"`

	Created string `db:"created"`
	Updated string `db:"updated"`
}

type PlaylistItem struct {
	PlaylistId string `db:"playlist_id"`
	TrackId    string `db:"track_id"`
}

func PlaylistQuery() *goqu.SelectDataset {
	query := dialect.From("playlists").
		Select(
			"playlists.id",
			"playlists.name",
			"playlists.picture",

			"playlists.owner_id",

			"playlists.created",
			"playlists.updated",
		).
		Prepared(true)

	return query
}

func (db *Database) GetPlaylistsByUser(ctx context.Context, userId string) ([]Playlist, error) {
	query := PlaylistQuery().
		Where(goqu.I("playlists.owner_id").Eq(userId))

	var items []Playlist
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetPlaylistById(ctx context.Context, id string) (Playlist, error) {
	query := PlaylistQuery().
		Where(goqu.I("playlists.id").Eq(id))

	var item Playlist
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Playlist{}, ErrItemNotFound
		}

		return Playlist{}, err
	}

	return item, nil
}

func (db *Database) GetPlaylistItems(ctx context.Context, playlistId string) ([]PlaylistItem, error) {
	query := dialect.From("playlist_items").
		Select(
			"playlist_items.playlist_id",
			"playlist_items.track_id",
		).
		Where(goqu.I("playlist_id").Eq(playlistId))

	var items []PlaylistItem
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetPlaylistTracks(ctx context.Context, playlistId string) ([]Track, error) {
	tracks := TrackQuery().As("tracks")

	query := dialect.From("playlist_items").
		Select("tracks.*").
		Join(tracks, goqu.On(goqu.I("tracks.id").Eq(goqu.I("playlist_items.track_id")))).
		Where(goqu.I("playlist_items.playlist_id").Eq(playlistId)).
		Order(goqu.I("tracks.name").Asc())

	var items []Track
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type CreatePlaylistParams struct {
	Id      string
	Name    string
	Picture sql.NullString

	OwnerId string

	Created int64
	Updated int64
}

func (db *Database) CreatePlaylist(ctx context.Context, params CreatePlaylistParams) (Playlist, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateId()
	}

	query := dialect.Insert("playlists").
		Rows(goqu.Record{
			"id":      id,
			"name":    params.Name,
			"picture": params.Picture,

			"owner_id": params.OwnerId,

			"created": created,
			"updated": updated,
		}).
		Returning(
			"playlists.id",
			"playlists.name",
			"playlists.picture",

			"playlists.owner_id",

			"playlists.created",
			"playlists.updated",
		).
		Prepared(true)

	var item Playlist
	err := db.Get(&item, query)
	if err != nil {
		return Playlist{}, err
	}

	return item, nil
}

func (db *Database) AddItemToPlaylist(ctx context.Context, playlistId, trackId string) error {
	query := goqu.Insert("playlist_items").
		Rows(goqu.Record{
			"playlist_id": playlistId,
			"track_id":    trackId,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		var e sqlite3.Error
		if errors.As(err, &e) {
			if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (db *Database) RemovePlaylistItem(ctx context.Context, playlistId, trackId string) error {
	query := goqu.Delete("playlist_items").
		Where(goqu.And(
			goqu.I("playlist_items.playlist_id").Eq(playlistId),
			goqu.I("playlist_items.track_id").Eq(trackId),
		))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
