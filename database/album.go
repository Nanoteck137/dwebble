package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/dwebble/database/adapter"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

// TODO(patrik): Remove
// type AlbumResolverAdapter struct {
// }
//
// func (*AlbumResolverAdapter) DefaultSort() string {
// 	return "albums.name"
// }
//
// func (*AlbumResolverAdapter) MapSortName(name string) (filter.Name, error) {
// 	// TODO(patrik): Include all available fields to sort on
// 	switch name {
// 	case "album":
// 		return filter.Name{
// 			Kind: filter.NameKindString,
// 			Name: "albums.name",
// 		}, nil
// 	case "created":
// 		return filter.Name{
// 			Kind: filter.NameKindNumber,
// 			Name: "albums.created",
// 		}, nil
// 	case "updated":
// 		return filter.Name{
// 			Kind: filter.NameKindNumber,
// 			Name: "albums.updated",
// 		}, nil
// 	}
//
// 	return filter.Name{}, filter.UnknownName(name)
//
// }
//
// func (a *AlbumResolverAdapter) GetDefaultSort() string {
// 	return "albums.name"
// }
//
// func (a *AlbumResolverAdapter) MapNameToId(typ, name string) (string, error) {
// 	return "", fmt.Errorf("Unknown name type: %s (%s)", typ, name)
// }
//
// func (a *AlbumResolverAdapter) MapName(name string) (filter.Name, error) {
// 	switch name {
// 	case "artist":
// 		return filter.Name{
// 			Kind: filter.NameKindString,
// 			Name: "artists.name",
// 		}, nil
// 	case "album":
// 		return filter.Name{
// 			Kind: filter.NameKindString,
// 			Name: "albums.name",
// 		}, nil
// 	}
//
// 	return filter.Name{}, fmt.Errorf("Unknown name: %s", name)
// }
//
// func (a *AlbumResolverAdapter) ResolveTable(typ string) (filter.Table, error) {
// 	return filter.Table{}, fmt.Errorf("Unknown table type: %s", typ)
// }
//
// func (a *AlbumResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
// 	return nil, fmt.Errorf("Unknown function name: %s", name)
// }

type Album struct {
	Id        string         `db:"id"`
	Name      string         `db:"name"`
	OtherName sql.NullString `db:"other_name"`

	ArtistId string `db:"artist_id"`

	CoverArt sql.NullString `db:"cover_art"`
	Year     sql.NullInt64  `db:"year"`

	ArtistName      string `db:"artist_name"`
	ArtistOtherName sql.NullString `db:"artist_other_name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	ExtraArtists ExtraArtists `db:"extra_artists"`
}

func AlbumQuery() *goqu.SelectDataset {
	query := dialect.From("albums").
		Select(
			"albums.id",
			"albums.name",
			"albums.other_name",

			"albums.artist_id",

			"albums.cover_art",
			"albums.year",

			"albums.created",
			"albums.updated",

			goqu.I("artists.name").As("artist_name"),
			goqu.I("artists.other_name").As("artist_other_name"),

			goqu.I("extra_artists.artists").As("extra_artists"),
		).
		Prepared(true).
		Join(
			goqu.I("artists"),
			goqu.On(
				goqu.I("albums.artist_id").Eq(goqu.I("artists.id")),
			),
		).
		LeftJoin(
			ExtraArtistsQuery("albums_extra_artists", "album_id").As("extra_artists"),
			goqu.On(goqu.I("albums.id").Eq(goqu.I("extra_artists.id"))),
		)

	return query
}

func (db *Database) GetAllAlbums(ctx context.Context, filterStr string, sortStr string) ([]Album, error) {
	query := AlbumQuery()

	var err error

	a := adapter.AlbumResolverAdapter{}
	resolver := filter.New(&a)

	query, err = applyFilter(query, resolver, filterStr)
	if err != nil {
		return nil, err
	}

	query, err = applySort(query, resolver, sortStr)
	if err != nil {
		return nil, err
	}

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
	Id        string
	Name      string
	OtherName sql.NullString

	ArtistId string

	CoverArt sql.NullString
	Year     sql.NullInt64

	Created int64
	Updated int64
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
		id = utils.CreateAlbumId()
	}

	query := dialect.Insert("albums").
		Rows(goqu.Record{
			"id":         id,
			"name":       params.Name,
			"other_name": params.OtherName,

			"artist_id": params.ArtistId,

			"cover_art": params.CoverArt,
			"year":      params.Year,

			"created": created,
			"updated": updated,
		}).
		Returning(
			"albums.id",
			"albums.name",
			"albums.other_name",

			"albums.artist_id",

			"albums.cover_art",
			"albums.year",
		).
		Prepared(true)

	var item Album
	err := db.Get(&item, query)
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
	Name      types.Change[string]
	OtherName types.Change[sql.NullString]

	ArtistId types.Change[string]

	CoverArt types.Change[sql.NullString]
	Year     types.Change[sql.NullInt64]

	Created types.Change[int64]
}

func (db *Database) UpdateAlbum(ctx context.Context, id string, changes AlbumChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "other_name", changes.OtherName)

	addToRecord(record, "artist_id", changes.ArtistId)

	addToRecord(record, "cover_art", changes.CoverArt)
	addToRecord(record, "year", changes.Year)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("albums").
		Set(record).
		Where(goqu.I("albums.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db *Database) RemoveAllExtraArtistsFromAlbum(ctx context.Context, trackId string) error {
	query := dialect.Delete("albums_extra_artists").
		Where(goqu.I("albums_extra_artists.album_id").Eq(trackId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db *Database) AddExtraArtistToAlbum(ctx context.Context, albumId, artistId string) error {
	query := dialect.Insert("albums_extra_artists").
		Rows(goqu.Record{
			"album_id":  albumId,
			"artist_id": artistId,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) {
			if sqlErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}
