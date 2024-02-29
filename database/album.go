package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
)

type Album struct {
	Id       string
	Name     string
	CoverArt string
	ArtistId string
	Path     string
}

func (db *Database) GetAllAlbums(ctx context.Context) ([]Album, error) {
	ds := dialect.From("albums").Select("id", "name", "cover_art", "artist_id", "path")

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Album
	for rows.Next() {
		var item Album
		err := rows.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetAlbumsByArtist(ctx context.Context, artistId string) ([]Album, error) {
	ds := dialect.From("albums").
		Select("id", "name", "cover_art", "artist_id", "path").
		Where(goqu.C("artist_id").Eq(artistId))

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Album
	for rows.Next() {
		var item Album
		err := rows.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetAlbumById(ctx context.Context, id string) (Album, error) {
	ds := dialect.From("albums").
		Select("id", "name", "cover_art", "artist_id", "path").
		Where(goqu.C("id").Eq(id)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Album{}, err
	}

	var item Album
	err = row.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Album{}, types.ErrNoAlbum
		}

		return Album{}, err
	}

	return item, nil
}

func (db *Database) GetAlbumByPath(ctx context.Context, path string) (Album, error) {
	ds := dialect.From("albums").
		Select("id", "name", "cover_art", "artist_id", "path").
		Where(goqu.C("path").Eq(path)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Album{}, err
	}

	var item Album
	err = row.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Album{}, types.ErrNoAlbum
		}

		return Album{}, err
	}

	return item, nil
}

type CreateAlbumParams struct {
	Name     string
	CoverArt string
	ArtistId string
	Path     string
}

func (db *Database) CreateAlbum(ctx context.Context, params CreateAlbumParams) (Album, error) {
	ds := dialect.Insert("albums").
		Rows(goqu.Record{
			"id":        utils.CreateId(),
			"name":      params.Name,
			"cover_art": params.CoverArt,
			"artist_id": params.ArtistId,
			"path":      params.Path,
		}).
		Returning("id", "name", "cover_art", "artist_id", "path").
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Album{}, err
	}

	var item Album
	err = row.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
	if err != nil {
		return Album{}, err
	}

	return item, nil
}
