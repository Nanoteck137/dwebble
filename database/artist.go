package database

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/dwebble/database/adapter"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/ember"
)

type Artist struct {
	RowId int `db:"rowid"`

	Id   string `db:"id"`
	Slug string `db:"slug"`

	Name      string         `db:"name"`
	OtherName sql.NullString `db:"other_name"`

	Picture sql.NullString `db:"picture"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`
}

func ArtistQuery() *goqu.SelectDataset {
	tags := dialect.From("artists_tags").
		Select(
			goqu.I("artists_tags.artist_id").As("artist_id"),
			goqu.Func("group_concat", goqu.I("tags.slug"), ",").As("tags"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(goqu.I("artists_tags.tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(goqu.I("artists_tags.artist_id"))

	query := dialect.From("artists").
		Select(
			"artists.rowid",

			"artists.id",
			"artists.slug",

			"artists.name",
			"artists.other_name",

			"artists.picture",

			"artists.updated",
			"artists.created",

			goqu.I("tags.tags").As("tags"),
		).
		LeftJoin(
			tags.As("tags"),
			goqu.On(goqu.I("artists.id").Eq(goqu.I("tags.artist_id"))),
		)

	return query
}

func (db DB) GetAllArtists(ctx context.Context, filterStr, sortStr string) ([]Artist, error) {
	query := ArtistQuery()

	var err error

	a := adapter.ArtistResolverAdapter{}
	resolver := filter.New(&a)

	query, err = applyFilter(query, resolver, filterStr)
	if err != nil {
		return nil, err
	}

	query, err = applySort(query, resolver, sortStr)
	if err != nil {
		return nil, err
	}

	return ember.Multiple[Artist](db.db, ctx, query)
}

func (db DB) GetArtistsPaged(ctx context.Context, opts FetchOptions) ([]Artist, types.Page, error) {
	query := ArtistQuery()

	var err error

	a := adapter.ArtistResolverAdapter{}
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
		Select(goqu.COUNT("artists.id"))

	if opts.PerPage > 0 {
		query = query.
			Limit(uint(opts.PerPage)).
			Offset(uint(opts.Page * opts.PerPage))
	}

	totalItems, err := ember.Single[int](db.db, ctx, countQuery)
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

	items, err := ember.Multiple[Artist](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db DB) GetArtistById(ctx context.Context, id string) (Artist, error) {
	query := ArtistQuery().
		Where(goqu.I("artists.id").Eq(id))

	return ember.Single[Artist](db.db, ctx, query)
}

func (db DB) GetArtistBySlug(ctx context.Context, slug string) (Artist, error) {
	query := ArtistQuery().
		Where(goqu.I("artists.slug").Eq(slug))

	return ember.Single[Artist](db.db, ctx, query)
}

func (db DB) GetArtistByName(ctx context.Context, name string) (Artist, error) {
	query := ArtistQuery().
		Where(goqu.Func("LOWER", goqu.I("artists.name")).Eq(strings.ToLower(name)))

	return ember.Single[Artist](db.db, ctx, query)
}

type CreateArtistParams struct {
	Id   string
	Slug string

	Name      string
	OtherName sql.NullString

	Picture sql.NullString

	Created int64
	Updated int64
}

func (db DB) CreateArtist(ctx context.Context, params CreateArtistParams) (Artist, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateArtistId()
	}

	query := dialect.Insert("artists").Rows(goqu.Record{
		"id":   id,
		"slug": params.Slug,

		"name":       params.Name,
		"other_name": params.OtherName,

		"picture": params.Picture,

		"created": created,
		"updated": updated,
	}).
		Returning(
			"artists.id",
			"artists.name",
			"artists.other_name",

			"artists.picture",

			"artists.updated",
			"artists.created",
		)

	return ember.Single[Artist](db.db, ctx, query)
}

type ArtistChanges struct {
	Slug      types.Change[string]

	Name      types.Change[string]

	OtherName types.Change[sql.NullString]
	Picture   types.Change[sql.NullString]

	Created types.Change[int64]
}

func (db DB) UpdateArtist(ctx context.Context, id string, changes ArtistChanges) error {
	record := goqu.Record{}

	addToRecord(record, "slug", changes.Slug)

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "other_name", changes.OtherName)

	addToRecord(record, "picture", changes.Picture)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("artists").
		Set(record).
		Where(goqu.I("artists.id").Eq(id))

	_, err := db.db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteArtist(ctx context.Context, id string) error {
	query := dialect.Delete("artists").
		Where(goqu.I("artists.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) MergeArtist(ctx context.Context, targetId, artistId string) error {
	// TODO(patrik): Changes
	// - albums
	// - albums_featuring_artists
	// - tracks
	// - tracks_featuring_artists
	// - artists_tags (maybe)

	albums, err := db.GetAlbumsByArtist(ctx, artistId)
	if err != nil {
		return err
	}

	for _, album := range albums {
		err = db.UpdateAlbum(ctx, album.Id, AlbumChanges{
			ArtistId: types.Change[string]{
				Value:   targetId,
				Changed: true,
			},
		})
		if err != nil {
			return err
		}
	}

	{
		record := goqu.Record{
			"artist_id": targetId,
		}
		query := goqu.Update("albums_featuring_artists").
			Set(record).
			Where(goqu.I("albums_featuring_artists.artist_id").Eq(artistId))

		_, err = db.db.Exec(ctx, query)
		if err != nil {
			// TODO(patrik): Change
			var sqlErr sqlite3.Error
			if errors.As(err, &sqlErr) {
				if sqlErr.ExtendedCode != sqlite3.ErrConstraintPrimaryKey {
					return err
				}
			} else {
				return err
			}
		}
	}

	{
		record := goqu.Record{
			"artist_id": targetId,
			"updated":   time.Now().UnixMilli(),
		}
		query := goqu.Update("tracks").
			Set(record).
			Where(goqu.I("tracks.artist_id").Eq(artistId))

		_, err = db.db.Exec(ctx, query)
		if err != nil {
			return err
		}
	}

	{
		record := goqu.Record{
			"artist_id": targetId,
		}
		query := goqu.Update("tracks_featuring_artists").
			Set(record).
			Where(goqu.I("tracks_featuring_artists.artist_id").Eq(artistId))

		_, err = db.db.Exec(ctx, query)
		if err != nil {
			// TODO(patrik): Change
			var sqlErr sqlite3.Error
			if errors.As(err, &sqlErr) {
				if sqlErr.ExtendedCode != sqlite3.ErrConstraintPrimaryKey {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

// TODO(patrik): Generalize
// TODO(patrik): Rename to AddArtistTag, same with track
func (db DB) AddTagToArtist(ctx context.Context, tagSlug, artistId string) error {
	ds := dialect.Insert("artists_tags").
		Rows(goqu.Record{
			"artist_id": artistId,
			"tag_slug":  tagSlug,
		})

	_, err := db.db.Exec(ctx, ds)
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
// TODO(patrik): Rename to RemoveAllArtistTags, same with track
func (db DB) RemoveAllTagsFromArtist(ctx context.Context, artistId string) error {
	query := dialect.Delete("artists_tags").
		Where(goqu.I("artist_id").Eq(artistId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
