package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go/ast"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type AlbumResolverAdapter struct {
}

func (a *AlbumResolverAdapter) GetDefaultSort() string {
	return "albums.name"
}

func (a *AlbumResolverAdapter) MapNameToId(typ, name string) (string, error) {
	return "", fmt.Errorf("Unknown name type: %s (%s)", typ, name)
}

func (a *AlbumResolverAdapter) MapName(name string) (types.Name, error) {
	switch name {
	case "artist":
		return types.Name{
			Kind: types.NameKindString,
			Name: "artists.name",
		}, nil
	case "album":
		return types.Name{
			Kind: types.NameKindString,
			Name: "albums.name",
		}, nil
	}

	return types.Name{}, fmt.Errorf("Unknown name: %s", name)
}

func (a *AlbumResolverAdapter) ResolveTable(typ string) (filter.Table, error) {
	return filter.Table{}, fmt.Errorf("Unknown table type: %s", typ)
}

func (a *AlbumResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	return nil, fmt.Errorf("Unknown function name: %s", name)
}

type Album struct {
	Id         string         `db:"id"`
	Name       string         `db:"name"`
	CoverArt   sql.NullString `db:"cover_art"`
	ArtistId   string         `db:"artist_id"`
	Path       string         `db:"path"`
	ArtistName string         `db:"artist_name"`
	Available  bool           `db:"available"`
}

func AlbumQuery() *goqu.SelectDataset {
	query := dialect.From("albums").
		Select(
			"albums.id",
			"albums.name",
			"albums.cover_art",
			"albums.artist_id",
			"albums.path",
			goqu.I("artists.name").As("artist_name"),
		).
		Join(
			goqu.I("artists"),
			goqu.On(
				goqu.I("albums.artist_id").Eq(goqu.I("artists.id")),
			),
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllAlbums(ctx context.Context, filterStr string, sortStr string, includeAll bool) ([]Album, error) {
	query := AlbumQuery()

	if filterStr != "" {
		a := AlbumResolverAdapter{}
		re, err := fullParseFilter(&a, filterStr)
		if err != nil {
			return nil, err
		}

		if includeAll {
			query = query.Where(re)
		} else {
			query = query.Where(goqu.I("albums.available").Eq(true), re)
		}
	} else {
		if !includeAll {
			query = query.Where(goqu.I("albums.available").Eq(true))
		}
	}

	var items []Album
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetAlbumsByArtist(ctx context.Context, artistId string) ([]Album, error) {
	query := AlbumQuery().
		Where(goqu.I("albums.artist_id").Eq(artistId))

	var items []Album
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetAlbumById(ctx context.Context, id string) (Album, error) {
	query := AlbumQuery().
		Where(goqu.I("albums.id").Eq(id))

	var item Album
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Album{}, ErrItemNotFound
		}

		return Album{}, err
	}

	return item, nil
}

func (db *Database) GetAlbumByPath(ctx context.Context, path string) (Album, error) {
	query := AlbumQuery().
		Where(goqu.I("albums.path").Eq(path))

	var item Album
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Album{}, ErrItemNotFound
		}

		return Album{}, err
	}

	return item, nil
}

func (db *Database) GetAlbumByName(ctx context.Context, name string) (Album, error) {
	query := AlbumQuery().
		Where(goqu.I("albums.name").Eq(name))

	var item Album
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Album{}, ErrItemNotFound
		}

		return Album{}, err
	}

	return item, nil
}

type CreateAlbumParams struct {
	Name     string
	CoverArt sql.NullString
	ArtistId string
	Path     string
}

func (db *Database) CreateAlbum(ctx context.Context, params CreateAlbumParams) (Album, error) {
	ds := dialect.Insert("albums").
		Rows(goqu.Record{
			"id":        utils.CreateId(),
			"name":      params.Name,
			"cover_art": params.CoverArt,
			"artist_id": params.ArtistId,
			"path":      params.Path,
			"available": false,
		}).
		Returning("id", "name", "cover_art", "artist_id", "path").
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Album{}, err
	}

	var item Album
	err = row.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
	if err != nil {
		return Album{}, err
	}

	return item, nil
}

type AlbumChanges struct {
	Name      types.Change[string]
	CoverArt  types.Change[sql.NullString]
	ArtistId  types.Change[string]
	Available bool
}

func (db *Database) UpdateAlbum(ctx context.Context, id string, changes AlbumChanges) error {
	record := goqu.Record{
		"available": changes.Available,
	}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "artist_id", changes.ArtistId)
	addToRecord(record, "cover_art", changes.CoverArt)

	ds := dialect.Update("albums").
		Set(record).
		Where(goqu.I("id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) MarkAllAlbumsUnavailable(ctx context.Context) error {
	ds := dialect.Update("albums").Set(goqu.Record{
		"available": false,
	})

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}
