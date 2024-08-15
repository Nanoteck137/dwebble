package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/utils"
)

type Genre struct {
	Id          string
	Name        string
	DisplayName string
}

func (db *Database) GetAllGenres(ctx context.Context) ([]Genre, error) {
	ds := dialect.From("genres").
		Select("id", "name", "display_name")

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Genre

	for rows.Next() {
		var item Genre
		err := rows.Scan(&item.Id, &item.Name, &item.DisplayName)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetGenreByName(ctx context.Context, name string) (Genre, error) {
	ds := dialect.From("genres").
		Select("id", "name", "display_name").
		Where(goqu.I("name").Eq(goqu.Func("LOWER", name))).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Genre{}, err
	}

	var item Genre
	err = row.Scan(&item.Id, &item.Name, &item.DisplayName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Genre{}, ErrItemNotFound
		}
		return Genre{}, err
	}

	return item, nil
}

func (db *Database) CreateGenre(ctx context.Context, name string) (Genre, error) {
	ds := dialect.Insert("genres").
		Rows(goqu.Record{
			"id":           utils.CreateId(),
			"name":         goqu.Func("LOWER", name),
			"display_name": name,
		}).
		Returning("id", "name", "display_name").
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Genre{}, err
	}

	var item Genre
	err = row.Scan(&item.Id, &item.Name, &item.DisplayName)
	if err != nil {
		return Genre{}, err
	}

	return item, nil
}

func (db *Database) AddGenreToTrack(ctx context.Context, genreId, trackId string) error {
	ds := dialect.Insert("tracks_to_genres").Rows(goqu.Record{
		"track_id": trackId,
		"genre_id": genreId,
	}).Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveGenreFromTrack(ctx context.Context, genreId, trackId string) error {
	ds := dialect.Delete("tracks_to_genres").
		Where(goqu.And(
			goqu.I("track_id").Eq(trackId),
			goqu.I("genre_id").Eq(genreId),
		)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetTrackGenres(ctx context.Context, trackId string) ([]Genre, error) {
	ds := dialect.From("tracks_to_genres").
		Select("genres.id", "genres.name", "genres.display_name").
		Join(
			goqu.I("genres"),
			goqu.On(
				goqu.I("tracks_to_genres.genre_id").Eq(goqu.I("genres.id")),
			),
		).
		Where(goqu.I("tracks_to_genres.track_id").Eq(trackId)).
		Prepared(true)

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Genre

	for rows.Next() {
		var item Genre
		err := rows.Scan(&item.Id, &item.Name, &item.DisplayName)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
