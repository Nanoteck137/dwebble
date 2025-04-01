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

type Album struct {
	RowId int `db:"rowid"`

	Id        string         `db:"id"`
	Name      string         `db:"name"`
	OtherName sql.NullString `db:"other_name"`

	ArtistId string `db:"artist_id"`

	CoverArt sql.NullString `db:"cover_art"`
	Year     sql.NullInt64  `db:"year"`

	ArtistName      string         `db:"artist_name"`
	ArtistOtherName sql.NullString `db:"artist_other_name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`

	FeaturingArtists FeaturingArtists `db:"featuring_artists"`
}

func AlbumQuery() *goqu.SelectDataset {
	tags := dialect.From("albums_tags").
		Select(
			goqu.I("albums_tags.album_id").As("album_id"),
			goqu.Func("group_concat", goqu.I("tags.slug"), ",").As("tags"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(goqu.I("albums_tags.tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(goqu.I("albums_tags.album_id"))

	query := dialect.From("albums").
		Select(
			"albums.rowid",

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

			goqu.I("tags.tags").As("tags"),

			goqu.I("featuring_artists.artists").As("featuring_artists"),
		).
		Prepared(true).
		Join(
			goqu.I("artists"),
			goqu.On(
				goqu.I("albums.artist_id").Eq(goqu.I("artists.id")),
			),
		).
		LeftJoin(
			tags.As("tags"),
			goqu.On(goqu.I("albums.id").Eq(goqu.I("tags.album_id"))),
		).
		LeftJoin(
			FeaturingArtistsQuery("albums_featuring_artists", "album_id").As("featuring_artists"),
			goqu.On(goqu.I("albums.id").Eq(goqu.I("featuring_artists.id"))),
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

func (db *Database) GetAlbumsPaged(ctx context.Context, opts FetchOptions) ([]Album, types.Page, error) {
	query := AlbumQuery()

	var err error

	a := adapter.AlbumResolverAdapter{}
	resolver := filter.New(&a)

	query, err = applyFilter(query, resolver, opts.Filter)
	if err != nil {
		return nil, types.Page{}, err
	}

	query, err = applySort(query, resolver, opts.Sort)
	if err != nil {
		return nil, types.Page{}, err
	}

	countQuery := query.
		Select(goqu.COUNT("albums.id"))

	if opts.PerPage > 0 {
		query = query.
			Limit(uint(opts.PerPage)).
			Offset(uint(opts.Page * opts.PerPage))
	}

	var totalItems int
	err = db.Get(&totalItems, countQuery)
	if err != nil {
		return nil, types.Page{}, err
	}

	totalPages := utils.TotalPages(opts.PerPage, totalItems)
	page := types.Page{
		Page:       opts.Page,
		PerPage:    opts.PerPage,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	var items []Album
	err = db.Select(&items, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
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

func (db *Database) DeleteAlbum(ctx context.Context, id string) error {
	query := dialect.Delete("albums").
		Prepared(true).
		Where(goqu.I("albums.id").Eq(id))

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

func (db *Database) ChangeAllAlbumArtist(ctx context.Context, artistId, newArtistId string) error {
	record := goqu.Record{
		"artist_id": newArtistId,
		"updated":   time.Now().UnixMilli(),
	}
	query := goqu.Update("albums").
		Prepared(true).
		Set(record).
		Where(goqu.I("albums.artist_id").Eq(artistId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
// TODO(patrik): Rename to AddAlbumTag, same with track
func (db *Database) AddTagToAlbum(ctx context.Context, tagSlug, albumId string) error {
	ds := dialect.Insert("albums_tags").
		Prepared(true).
		Rows(goqu.Record{
			"album_id": albumId,
			"tag_slug": tagSlug,
		})

	_, err := db.Exec(ctx, ds)
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

// TODO(patrik): Generalize
// TODO(patrik): Rename to RemoveAllAlbumTags, same with track
func (db *Database) RemoveAllTagsFromAlbum(ctx context.Context, albumId string) error {
	query := dialect.Delete("albums_tags").
		Prepared(true).
		Where(goqu.I("album_id").Eq(albumId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db *Database) RemoveAllAlbumFeaturingArtists(ctx context.Context, albumId string) error {
	query := dialect.Delete("albums_featuring_artists").
		Prepared(true).
		Where(goqu.I("albums_featuring_artists.album_id").Eq(albumId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db *Database) AddFeaturingArtistToAlbum(ctx context.Context, albumId, artistId string) error {
	query := dialect.Insert("albums_featuring_artists").
		Prepared(true).
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

func (db *Database) RemoveFeaturingArtistFromAlbum(ctx context.Context, albumId, artistId string) error {
	query := goqu.Delete("albums_featuring_artists").
		Prepared(true).
		Where(
			goqu.And(
				goqu.I("albums_featuring_artists.album_id").Eq(albumId),
				goqu.I("albums_featuring_artists.artist_id").Eq(artistId),
			),
		)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
