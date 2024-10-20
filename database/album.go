package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go/ast"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/tools/sort"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type AlbumResolverAdapter struct {
}

func (*AlbumResolverAdapter) DefaultSort() string {
	return "albums.name"
}

func (*AlbumResolverAdapter) MapSortName(name string) (types.Name, error) {
	// TODO(patrik): Include all available fields to sort on
	switch name {
	case "album":
		return types.Name{
			Kind: types.NameKindString,
			Name: "albums.name",
		}, nil
	case "available":
		return types.Name{
			Kind: types.NameKindNumber,
			Name: "albums.available",
		}, nil
	case "created":
		return types.Name{
			Kind: types.NameKindNumber,
			Name: "albums.created",
		}, nil
	case "updated":
		return types.Name{
			Kind: types.NameKindNumber,
			Name: "albums.updated",
		}, nil
	}

	return types.Name{}, sort.UnknownName(name)

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
	Id       string `db:"id"`
	Name     string `db:"name"`
	ArtistId string `db:"artist_id"`

	CoverArt sql.NullString `db:"cover_art"`
	Year     sql.NullInt64  `db:"year"`

	ArtistName string `db:"artist_name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Available bool `db:"available"`
}

func AlbumQuery() *goqu.SelectDataset {
	query := dialect.From("albums").
		Select(
			"albums.id",
			"albums.name",
			"albums.artist_id",

			"albums.cover_art",
			"albums.year",

			"albums.created",
			"albums.updated",

			goqu.I("artists.name").As("artist_name"),

			"albums.available",
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

	sortExpr, err := sort.Parse(sortStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidSort, err)
	}

	resolver := sort.New(&AlbumResolverAdapter{})

	sortExpr, err = resolver.Resolve(sortExpr)
	if err != nil {
		if errors.Is(err, sort.ErrUnknownName) {
			return nil, fmt.Errorf("%w: %w", ErrInvalidSort, err)
		}

		return nil, err
	}

	exprs, err := generateSort(sortExpr)
	if err != nil {
		return nil, err
	}

	query = query.Order(exprs...)

	var items []Album
	err = db.Select(&items, query)
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
	Id       string
	Name     string
	ArtistId string

	CoverArt sql.NullString
	Year     sql.NullInt64

	Created int64
	Updated int64

	Available bool
}

func (db *Database) CreateAlbum(ctx context.Context, params CreateAlbumParams) (Album, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.Slug(params.Name) + "-" + utils.CreateSmallId()
	}

	ds := dialect.Insert("albums").
		Rows(goqu.Record{
			"id":        id,
			"name":      params.Name,
			"artist_id": params.ArtistId,

			"cover_art": params.CoverArt,
			"year":      params.Year,

			"created": created,
			"updated": updated,

			"available": params.Available,
		}).
		Returning("id", "name", "artist_id", "cover_art", "year", "available").
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Album{}, err
	}

	var item Album
	err = row.Scan(&item.Id, &item.Name, &item.ArtistId, &item.CoverArt, &item.Year, &item.Available)
	if err != nil {
		return Album{}, err
	}

	return item, nil
}

func (db *Database) RemoveAlbum(ctx context.Context, id string) error {
	query := dialect.Delete("albums").
		Prepared(true).
		Where(goqu.I("albums.id").Eq(id))

	// TODO(patrik): Check result?
	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type AlbumChanges struct {
	Name     types.Change[string]
	ArtistId types.Change[string]

	CoverArt types.Change[sql.NullString]
	Year     types.Change[sql.NullInt64]

	Created types.Change[int64]

	Available types.Change[bool]
}

func (db *Database) UpdateAlbum(ctx context.Context, id string, changes AlbumChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "artist_id", changes.ArtistId)
	addToRecord(record, "cover_art", changes.CoverArt)
	addToRecord(record, "year", changes.Year)

	addToRecord(record, "created", changes.Created)

	addToRecord(record, "available", changes.Available)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

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
