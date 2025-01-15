package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
)

type Tag struct {
	Slug string `db:"slug"`
}

func TagQuery() *goqu.SelectDataset {
	query := dialect.From("tags").
		Select(
			"tags.slug",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllTags(ctx context.Context) ([]Tag, error) {
	query := TagQuery()

	var items []Tag
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetTagBySlug(ctx context.Context, slug string) (Tag, error) {
	query := TagQuery().
		Where(goqu.I("tags.slug").Eq(slug))

	var item Tag
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Tag{}, ErrItemNotFound
		}

		return Tag{}, err
	}

	return item, nil
}

func (db *Database) CreateTag(ctx context.Context, slug string) error {
	query := dialect.Insert("tags").
		Rows(goqu.Record{
			"slug": slug,
		}).
		Prepared(true)

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
