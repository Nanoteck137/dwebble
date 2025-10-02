package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/ember"
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

func (db DB) CreateTaglist(ctx context.Context, params CreateTaglistParams) (Taglist, error) {
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
		)

	return ember.Single[Taglist](db.db, ctx, query)
}

func (db DB) GetTaglistByUser(ctx context.Context, userId string) ([]Taglist, error) {
	query := dialect.From("taglists").
		Select(
			"taglists.id",
			"taglists.name",
			"taglists.filter",
			"taglists.owner_id",
		).
		Where(goqu.I("owner_id").Eq(userId))

	return ember.Multiple[Taglist](db.db, ctx, query)
}

func (db DB) GetTaglistById(ctx context.Context, id string) (Taglist, error) {
	query := dialect.From("taglists").
		Select(
			"taglists.id",
			"taglists.name",
			"taglists.filter",
			"taglists.owner_id",
			"taglists.created",
			"taglists.updated",
		).
		Where(goqu.I("taglists.id").Eq(id))

	return ember.Single[Taglist](db.db, ctx, query)
}

type TaglistChanges struct {
	Name   types.Change[string]
	Filter types.Change[string]

	Created types.Change[int64]
}

func (db DB) UpdateTaglist(ctx context.Context, id string, changes TaglistChanges) error {
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
		Where(goqu.I("taglists.id").Eq(id))

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteTaglist(ctx context.Context, id string) error {
	query := dialect.Delete("taglists").
		Where(goqu.I("taglists.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
