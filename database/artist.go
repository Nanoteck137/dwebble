package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type Artist struct {
	Id      string         `db:"id"`
	Name    string         `db:"name"`
	Picture sql.NullString `db:"picture"`
}

func ArtistQuery() *goqu.SelectDataset {
	query := dialect.From("artists").
		Select(
			"artists.id",
			"artists.name",
			"artists.picture",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllArtists(ctx context.Context) ([]Artist, error) {
	query := ArtistQuery()

	var items []Artist
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetArtistById(ctx context.Context, id string) (Artist, error) {
	query := ArtistQuery().
		Where(goqu.I("artists.id").Eq(id))

	var item Artist
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Artist{}, ErrItemNotFound
		}

		return Artist{}, err
	}

	return item, nil
}

func (db *Database) GetArtistByName(ctx context.Context, name string) (Artist, error) {
	query := ArtistQuery().
		Where(goqu.I("artists.name").Eq(name))

	var item Artist
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Artist{}, ErrItemNotFound
		}

		return Artist{}, err
	}

	return item, nil
}

type CreateArtistParams struct {
	Id      string
	Name    string
	Picture sql.NullString
}

func (db *Database) CreateArtist(ctx context.Context, params CreateArtistParams) (Artist, error) {
	id := params.Id
	if id == "" {
		id = utils.CreateArtistId()
	}

	ds := dialect.Insert("artists").Rows(goqu.Record{
		"id":      id,
		"name":    params.Name,
		"picture": params.Picture,
	}).Returning("id", "name", "picture").Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Artist{}, err
	}

	var item Artist
	err = row.Scan(&item.Id, &item.Name, &item.Picture)
	if err != nil {
		return Artist{}, err
	}

	return item, nil
}

type ArtistChanges struct {
	Name    types.Change[string]
	Picture types.Change[sql.NullString]
}

func (db *Database) UpdateArtist(ctx context.Context, id string, changes ArtistChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "picture", changes.Picture)

	if len(record) <= 0 {
		return nil
	}

	ds := dialect.Update("artists").
		Set(record).
		Where(goqu.I("artists.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) MarkAllArtistsUnavailable(ctx context.Context) error {
	ds := dialect.Update("artists").Set(goqu.Record{
		"available": false,
	})

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}
