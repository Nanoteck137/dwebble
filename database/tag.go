package database

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/utils"
)

type Tag struct {
	Id   string
	Name string
}

func (db *Database) GetAllTags(ctx context.Context) ([]Tag, error) {
	ds := dialect.From("tags").
		Select("id", "name")

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Tag

	for rows.Next() {
		var item Tag
		err := rows.Scan(&item.Id, &item.Name)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetTagByName(ctx context.Context, name string) (Tag, error) {
	ds := dialect.From("tags").
		Select("id", "name").
		Where(goqu.I("name").Eq(name)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Tag{}, err
	}

	var item Tag
	err = row.Scan(&item.Id, &item.Name)
	if err != nil {
		return Tag{}, err
	}

	return item, nil
}

func (db *Database) CreateTag(ctx context.Context, name string) (Tag, error) {
	ds := dialect.Insert("tags").
		Rows(goqu.Record{
			"id":   utils.CreateId(),
			"name": name,
		}).
		Returning("id", "name").
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Tag{}, err
	}

	var item Tag
	err = row.Scan(&item.Id, &item.Name)
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

	tag, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	fmt.Printf("tag: %v\n", tag)

	return nil
}

func (db *Database) RemoveTagFromTrack(ctx context.Context, tagId, trackId string) error {
	ds := dialect.Delete("tracks_to_tags").
		Where(goqu.And(
			goqu.I("track_id").Eq(trackId),
			goqu.I("tag_id").Eq(tagId),
		)).
		Prepared(true)

	tag, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	fmt.Printf("tag: %v\n", tag)

	return nil
}

func (db *Database) GetTrackTags(ctx context.Context, trackId string) ([]Tag, error) {
	ds := dialect.From("tracks_to_tags").
		Select("tags.id", "tags.name").
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
		err := rows.Scan(&item.Id, &item.Name)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
