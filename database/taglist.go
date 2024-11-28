package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type Taglist struct {
	Id      string `db:"id"`
	Name    string `db:"name"`
	Filter  string `db:"filter"`
	OwnerId string `db:"owner_id"`
}

type CreateTaglistParams struct {
	Name    string
	Filter  string
	OwnerId string
}

func (db *Database) CreateTaglist(ctx context.Context, params CreateTaglistParams) (Taglist, error) {
	query := dialect.Insert("taglists").
		Rows(goqu.Record{
			"id":       utils.CreateId(),
			"name":     params.Name,
			"filter":   params.Filter,
			"owner_id": params.OwnerId,
		}).
		Returning("id", "name", "filter", "owner_id").
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
		).
		Prepared(true).
		Where(goqu.I("taglists.id").Eq(id))

	var item Taglist
	err := db.Get(&item, query)
	if err != nil {
		return Taglist{}, err
	}

	return item, nil
}

type TaglistChanges struct {
	Name   types.Change[string]
	Filter types.Change[string]
}

func (db *Database) UpdateTaglist(ctx context.Context, id string, changes TaglistChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "filter", changes.Filter)

	if len(record) == 0 {
		return nil
	}

	// record["updated"] = time.Now().UnixMilli()

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
