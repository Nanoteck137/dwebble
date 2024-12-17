package database

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/database/adapter"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type Artist struct {
	Id        string         `db:"id"`
	Name      string         `db:"name"`
	OtherName sql.NullString `db:"other_name"`

	Picture sql.NullString `db:"picture"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func ArtistQuery() *goqu.SelectDataset {
	query := dialect.From("artists").
		Select(
			"artists.id",
			"artists.name",
			"artists.other_name",

			"artists.picture",

			"artists.updated",
			"artists.created",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllArtists(ctx context.Context, filterStr, sortStr string) ([]Artist, error) {
	query := ArtistQuery()

	var err error

	a := adapter.ArtistResolverAdapter{}
	resolver := filter.New(&a)

	query, err = applyFilter(query, resolver, filterStr)
	if err != nil {
		return nil, err
	}

	query, err = applySort(query, resolver, sortStr)
	if err != nil {
		return nil, err
	}

	var items []Artist
	err = db.Select(&items, query)
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
		Where(goqu.Func("LOWER", goqu.I("artists.name")).Eq(strings.ToLower(name)))

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
	Id        string
	Name      string
	OtherName sql.NullString

	Picture sql.NullString

	Created int64
	Updated int64
}

func (db *Database) CreateArtist(ctx context.Context, params CreateArtistParams) (Artist, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateArtistId()
	}

	query := dialect.Insert("artists").Rows(goqu.Record{
		"id":         id,
		"name":       params.Name,
		"other_name": params.OtherName,

		"picture": params.Picture,

		"created": created,
		"updated": updated,
	}).
		Returning(
			"artists.id",
			"artists.name",
			"artists.other_name",

			"artists.picture",

			"artists.updated",
			"artists.created",
		).
		Prepared(true)

	var item Artist
	err := db.Get(&item, query)
	if err != nil {
		return Artist{}, err
	}

	return item, nil
}

type ArtistChanges struct {
	Name      types.Change[string]
	OtherName types.Change[sql.NullString]
	Picture   types.Change[sql.NullString]

	Created types.Change[int64]
}

func (db *Database) UpdateArtist(ctx context.Context, id string, changes ArtistChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "other_name", changes.OtherName)

	addToRecord(record, "picture", changes.Picture)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

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
