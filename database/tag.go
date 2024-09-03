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
	Name string `db:"name"`
}

func TagQuery() *goqu.SelectDataset {
	query := dialect.From("tags").
		Select(
			"tags.slug",
			"tags.name",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllTags(ctx context.Context) ([]Tag, error) {
	query := dialect.From("tags").
		Select("tags.slug", "tags.name")

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

func (db *Database) CreateTag(ctx context.Context, slug, name string) error {
	query := dialect.Insert("tags").
		Rows(goqu.Record{
			"slug": slug,
			"name": name,
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

func (db *Database) AddTagToTrack(ctx context.Context, tagSlug, trackId string) error {
	ds := dialect.Insert("tracks_to_tags").Rows(goqu.Record{
		"track_id": trackId,
		"tag_slug":   tagSlug,
	}).Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveTagFromTrack(ctx context.Context, tagSlug, trackId string) error {
	ds := dialect.Delete("tracks_to_tags").
		Where(goqu.And(
			goqu.I("track_id").Eq(trackId),
			goqu.I("tag_slug").Eq(tagSlug),
		)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAllTagsFromTrack(ctx context.Context, trackId string) error {
	ds := dialect.Delete("tracks_to_tags").
		Where(goqu.And(
			goqu.I("track_id").Eq(trackId),
		)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetTrackTags(ctx context.Context, trackId string) ([]Tag, error) {
	query := dialect.From("tracks_to_tags").
		Select("tags.id", "tags.name", "tags.display_name").
		Join(
			goqu.I("tags"), 
			goqu.On(goqu.I("tracks_to_tags.tag_id").Eq(goqu.I("tags.id"))),
		).
		Where(goqu.I("tracks_to_tags.track_id").Eq(trackId)).
		Prepared(true)

	var items []Tag
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}
