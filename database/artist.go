package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
)

type Artist struct {
	Id      string
	Name    string
	Picture sql.NullString
	Path    string
}

func (db *Database) GetAllArtists(ctx context.Context) ([]Artist, error) {
	ds := dialect.From("artists").
		Select("id", "name", "picture", "path").
		Where(goqu.I("available").Eq(true))

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Artist
	for rows.Next() {
		var item Artist
		err := rows.Scan(&item.Id, &item.Name, &item.Picture, &item.Path)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetArtistById(ctx context.Context, id string) (Artist, error) {
	ds := dialect.From("artists").
		Select("id", "name", "picture", "path").
		Where(goqu.C("id").Eq(id)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Artist{}, err
	}

	var item Artist
	err = row.Scan(&item.Id, &item.Name, &item.Picture, &item.Path)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Artist{}, types.ErrNoArtist
		}

		return Artist{}, err
	}

	return item, nil
}

func (db *Database) GetArtistByPath(ctx context.Context, path string) (Artist, error) {
	ds := dialect.From("artists").
		Select("id", "name", "picture", "path").
		Where(goqu.C("path").Eq(path)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Artist{}, err
	}

	var item Artist
	err = row.Scan(&item.Id, &item.Name, &item.Picture, &item.Path)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Artist{}, types.ErrNoArtist
		}

		return Artist{}, err
	}

	return item, nil
}

func (db *Database) GetArtistByName(ctx context.Context, name string) (Artist, error) {
	ds := dialect.From("artists").
		Select("id", "name", "picture", "path").
		Where(goqu.C("name").Eq(name)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Artist{}, err
	}

	var item Artist
	err = row.Scan(&item.Id, &item.Name, &item.Picture, &item.Path)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Artist{}, types.ErrNoArtist
		}

		return Artist{}, err
	}

	return item, nil
}

type CreateArtistParams struct {
	Name    string
	Picture sql.NullString
	Path    string
}

func (db *Database) CreateArtist(ctx context.Context, params CreateArtistParams) (Artist, error) {
	ds := dialect.Insert("artists").Rows(goqu.Record{
		"id":        utils.CreateId(),
		"name":      params.Name,
		"picture":   params.Picture,
		"path":      params.Path,
		"available": false,
	}).Returning("id", "name", "picture", "path").Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Artist{}, err
	}

	var item Artist
	err = row.Scan(&item.Id, &item.Name, &item.Picture, &item.Path)
	if err != nil {
		return Artist{}, err
	}

	return item, nil
}

type ArtistChanges struct {
	Name      types.Change[string]
	Picture   types.Change[sql.NullString]
	Available bool
}

func (db *Database) UpdateArtist(ctx context.Context, id string, changes ArtistChanges) error {
	record := goqu.Record{
		"available": changes.Available,
	}

	if changes.Name.Changed {
		record["name"] = changes.Name.Value
	}

	if changes.Picture.Changed {
		record["picture"] = changes.Picture.Value
	}

	ds := dialect.Update("artists").
		Set(record).
		Where(goqu.I("id").Eq(id)).
		Prepared(true)

	tag, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	fmt.Printf("tag: %v\n", tag)

	return nil
}

func (db *Database) MarkAllArtistsUnavailable(ctx context.Context) error {
	ds := dialect.Update("artists").Set(goqu.Record{
		"available": false,
	})

	tag, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	fmt.Printf("tag: %v\n", tag)

	return nil
}
