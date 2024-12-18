package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type Taglist struct {
	Id     string `db:"id"`
	Name   string `db:"name"`
	Filter string `db:"filter"`

	OwnerId string `db:"owner_id"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

type CreateTaglistParams struct {
	Name    string
	Filter  string
	OwnerId string

	Created int64
	Updated int64
}

func (db *Database) CreateTaglist(ctx context.Context, params CreateTaglistParams) (Taglist, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	query := dialect.Insert("taglists").
		Rows(goqu.Record{
			"id":       utils.CreateId(),
			"name":     params.Name,
			"filter":   params.Filter,
			"owner_id": params.OwnerId,

			"created": created,
			"updated": updated,
		}).
		Returning(
			"taglists.id",
			"taglists.name",
			"taglists.filter",
			"taglists.owner_id",
			"taglists.created",
			"taglists.updated",
		).
		Prepared(true)

	var item Taglist
	err := db.Get(&item, query)
	if err != nil {
		return Taglist{}, nil
	}

	return item, nil
}

func (db *Database) GetTaglistByUser(ctx context.Context, userId string) ([]Taglist, error) {
	query := dialect.From("taglists").
		Select(
			"taglists.id",
			"taglists.name",
			"taglists.filter",
			"taglists.owner_id",
		).
		Prepared(true).
		Where(goqu.I("owner_id").Eq(userId))

	var items []Taglist
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetTaglistById(ctx context.Context, id string) (Taglist, error) {
	query := dialect.From("taglists").
		Select(
			"taglists.id",
			"taglists.name",
			"taglists.filter",
			"taglists.owner_id",
			"taglists.created",
			"taglists.updated",
		).
		Prepared(true).
		Where(goqu.I("taglists.id").Eq(id))

	var item Taglist
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Taglist{}, ErrItemNotFound
		}

		return Taglist{}, err
	}

	return item, nil
}

type TaglistChanges struct {
	Name   types.Change[string]
	Filter types.Change[string]

	Created types.Change[int64]
}

func (db *Database) UpdateTaglist(ctx context.Context, id string, changes TaglistChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "filter", changes.Filter)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("taglists").
		Set(record).
		Where(goqu.I("id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}
