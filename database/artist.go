package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

type Artist struct {
	Id      string
	Name    string
	Picture string
	Path    string
}

func (db *Database) GetAllArtists(ctx context.Context) ([]Artist, error) {
	ds := dialect.From("artists").Select("id", "name", "picture", "path")

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
		return Artist{}, err
	}

	return item, nil
}