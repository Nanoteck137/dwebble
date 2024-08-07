package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type Tag struct {
	Id          string
	Name        string
	DisplayName string
}

func (db *Database) GetAllTags(ctx context.Context) ([]Tag, error) {
	ds := dialect.From("tags").
		Select("id", "name", "display_name")

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Tag

	for rows.Next() {
		var item Tag
		err := rows.Scan(&item.Id, &item.Name, &item.DisplayName)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetTagByName(ctx context.Context, name string) (Tag, error) {
	ds := dialect.From("tags").
		Select("id", "name", "display_name").
		Where(goqu.I("name").Eq(goqu.Func("LOWER", name))).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Tag{}, err
	}

	var item Tag
	err = row.Scan(&item.Id, &item.Name, &item.DisplayName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Tag{}, types.ErrNoTag
		}

		return Tag{}, err
	}

	return item, nil
}

func (db *Database) CreateTag(ctx context.Context, name string) (Tag, error) {
	ds := dialect.Insert("tags").
		Rows(goqu.Record{
			"id":           utils.CreateId(),
			"name":         goqu.Func("LOWER", name),
			"display_name": name,
		}).
		Returning("id", "name", "display_name").
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Tag{}, err
	}

	var item Tag
	err = row.Scan(&item.Id, &item.Name, &item.DisplayName)
	if err != nil {
		return Tag{}, err
	}

	return item, nil
}

func (db *Database) AddTagToTrack(ctx context.Context, tagId, trackId string) error {
	ds := dialect.Insert("tracks_to_tags").Rows(goqu.Record{
		"track_id": trackId,
		"tag_id":   tagId,
	}).Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveTagFromTrack(ctx context.Context, tagId, trackId string) error {
	ds := dialect.Delete("tracks_to_tags").
		Where(goqu.And(
			goqu.I("track_id").Eq(trackId),
			goqu.I("tag_id").Eq(tagId),
		)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetTrackTags(ctx context.Context, trackId string) ([]Tag, error) {
	ds := dialect.From("tracks_to_tags").
		Select("tags.id", "tags.name", "tags.display_name").
		Join(goqu.I("tags"), goqu.On(goqu.I("tracks_to_tags.tag_id").Eq(goqu.I("tags.id")))).
		Where(goqu.I("tracks_to_tags.track_id").Eq(trackId)).
		Prepared(true)

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Tag

	for rows.Next() {
		var item Tag
		err := rows.Scan(&item.Id, &item.Name, &item.DisplayName)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
