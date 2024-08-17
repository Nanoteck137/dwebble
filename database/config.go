package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
)

type Config struct {
	OwnerId string
}

func (db *Database) CreateConfig(ctx context.Context, ownerId string) (Config, error) {
	ds := dialect.Insert("config").Rows(goqu.Record{
		"id":       1,
		"owner_id": ownerId,
	}).
		Returning("owner_id").
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)

	var item Config
	err = row.Scan(&item.OwnerId)
	if err != nil {
		return Config{}, err
	}

	return item, nil
}

func (db *Database) GetConfig(ctx context.Context) (*Config, error) {
	ds := dialect.From("config").
		Select("owner_id").
		Where(goqu.I("id").Eq(1))

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return nil, err
	}

	var item Config
	err = row.Scan(&item.OwnerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrItemNotFound
		}

		return nil, err
	}

	return &item, nil
}
