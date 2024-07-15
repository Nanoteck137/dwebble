package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go/ast"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/sortfilter/filter"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
)

type AlbumResolverAdapter struct {
}

func (a *AlbumResolverAdapter) GetDefaultSort() string {
	return "albums.name"
}

func (a *AlbumResolverAdapter) MapNameToId(typ, name string) (string, error) {
	return "", fmt.Errorf("Unknown name type: %s (%s)", typ, name)
}

func (a *AlbumResolverAdapter) MapName(name string) (filter.Name, error) {
	switch name {
	case "artist":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "artists.name",
		}, nil
	case "album":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "albums.name",
		}, nil
	}

	return filter.Name{}, fmt.Errorf("Unknown name: %s", name)
}

func (a *AlbumResolverAdapter) ResolveTable(typ string) (filter.Table, error) {
	return filter.Table{}, fmt.Errorf("Unknown table type: %s", typ)
}

func (a *AlbumResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	return nil, fmt.Errorf("Unknown function name: %s", name)
}

type Album struct {
	Id       string
	Name     string
	CoverArt sql.NullString
	ArtistId string
	Path     string
}

func (db *Database) GetAllAlbums(ctx context.Context, filterStr string) ([]Album, error) {
	ds := dialect.From("albums").
		Select(
			"albums.id",
			"albums.name",
			"albums.cover_art",
			"albums.artist_id",
			"albums.path",
		).
		Join(
			goqu.I("artists"),
			goqu.On(
				goqu.I("albums.artist_id").Eq(goqu.I("artists.id")),
			),
		).
		Prepared(true)

	if filterStr != "" {
		a := AlbumResolverAdapter{}
		re, err := FullParseFilter(&a, filterStr)
		if err != nil {
			return nil, err
		}

		ds = ds.Where(goqu.I("albums.available").Eq(true), re)
	} else {
		ds = ds.Where(goqu.I("albums.available").Eq(true))
	}

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Album
	for rows.Next() {
		var item Album
		err := rows.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetAlbumsByArtist(ctx context.Context, artistId string) ([]Album, error) {
	ds := dialect.From("albums").
		Select("id", "name", "cover_art", "artist_id", "path").
		Where(goqu.And(goqu.I("available").Eq(true), goqu.C("artist_id").Eq(artistId)))

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Album
	for rows.Next() {
		var item Album
		err := rows.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetAlbumById(ctx context.Context, id string) (Album, error) {
	ds := dialect.From("albums").
		Select("id", "name", "cover_art", "artist_id", "path").
		Where(goqu.C("id").Eq(id)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Album{}, err
	}

	var item Album
	err = row.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
	if err != nil {
		if err == sql.ErrNoRows {
			return Album{}, types.ErrNoAlbum
		}

		return Album{}, err
	}

	return item, nil
}

func (db *Database) GetAlbumByPath(ctx context.Context, path string) (Album, error) {
	ds := dialect.From("albums").
		Select("id", "name", "cover_art", "artist_id", "path").
		Where(goqu.C("path").Eq(path)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Album{}, err
	}

	var item Album
	err = row.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
	if err != nil {
		if err == sql.ErrNoRows {
			return Album{}, types.ErrNoAlbum
		}

		return Album{}, err
	}

	return item, nil
}

func (db *Database) GetAlbumByName(ctx context.Context, name string) (Album, error) {
	ds := dialect.From("albums").
		Select("id", "name", "cover_art", "artist_id", "path").
		Where(goqu.C("name").Eq(name)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Album{}, err
	}

	var item Album
	err = row.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Album{}, types.ErrNoAlbum
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

	if changes.Name.Changed {
		record["name"] = changes.Name.Value
	}

	if changes.ArtistId.Changed {
		record["artist_id"] = changes.ArtistId.Value
	}

	if changes.CoverArt.Changed {
		record["cover_art"] = changes.CoverArt.Value
	}

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
