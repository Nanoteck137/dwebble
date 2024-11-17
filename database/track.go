package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go/ast"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/tools/sort"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type TrackResolverAdapter struct{}

func (a *TrackResolverAdapter) DefaultSort() string {
	return "tracks.name"
}

func (a *TrackResolverAdapter) MapSortName(name string) (types.Name, error) {
	// TODO(patrik): Include all available fields to sort on
	// TODO(patrik): Change names to match the types.Track names
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
	case "track":
		return types.Name{
			Kind: types.NameKindString,
			Name: "tracks.name",
		}, nil
	case "trackNumber":
		return types.Name{
			Kind: types.NameKindNumber,
			Name: "tracks.track_number",
		}, nil
	}

	return types.Name{}, sort.UnknownName(name)

}

func (a *TrackResolverAdapter) MapNameToId(typ, name string) (string, error) {
	switch typ {
	case "tags":
		return utils.Slug(name), nil
	}

	return "", fmt.Errorf("Unknown name type: %s (%s)", typ, name)
}

func (a *TrackResolverAdapter) MapName(name string) (types.Name, error) {
	// TODO(patrik): Include all available fields to filter on
	// TODO(patrik): Change names to match the types.Track names
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
	case "albumId":
		return types.Name{
			Kind: types.NameKindString,
			Name: "albums.id",
		}, nil
	case "track":
		return types.Name{
			Kind: types.NameKindString,
			Name: "tracks.name",
		}, nil
	case "trackNumber":
		return types.Name{
			Kind: types.NameKindNumber,
			Name: "tracks.track_number",
		}, nil
	}

	return types.Name{}, fmt.Errorf("Unknown name: %s", name)
}

func (a *TrackResolverAdapter) ResolveTable(typ string) (filter.Table, error) {
	switch typ {
	case "tags":
		return filter.Table{
			Name:       "tracks_to_tags",
			SelectName: "track_id",
			WhereName:  "tag_slug",
		}, nil
	}

	return filter.Table{}, fmt.Errorf("Unknown table type: %s", typ)
}

func (a *TrackResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	switch name {
	case "hasTag":
		return resolver.InTable(name, "tags", args)
	}

	return nil, fmt.Errorf("Unknown function name: %s", name)
}

type Track struct {
	Id   string `db:"id"`
	Name string `db:"name"`

	AlbumId  string `db:"album_id"`
	ArtistId string `db:"artist_id"`

	Number   sql.NullInt64 `db:"track_number"`
	Duration sql.NullInt64 `db:"duration"`
	Year     sql.NullInt64 `db:"year"`

	ExportName       string `db:"export_name"`
	OriginalFilename string `db:"original_filename"`
	MobileFilename   string `db:"mobile_filename"`

	AlbumName     string         `db:"album_name"`
	AlbumCoverArt sql.NullString `db:"album_cover_art"`
	ArtistName    string         `db:"artist_name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`
}

type TrackChanges struct {
	Name     types.Change[string]
	AlbumId  types.Change[string]
	ArtistId types.Change[string]

	Number   types.Change[sql.NullInt64]
	Duration types.Change[sql.NullInt64]
	Year     types.Change[sql.NullInt64]

	ExportName       types.Change[string]
	OriginalFilename types.Change[string]
	MobileFilename   types.Change[string]

	Created types.Change[int64]
}

func TrackQuery() *goqu.SelectDataset {
	// TODO(patrik): Add Back
	tags := dialect.From("tracks_to_tags").
		Select(
			goqu.I("tracks_to_tags.track_id").As("track_id"),
			goqu.Func("group_concat", goqu.I("tags.name"), ",").As("tags"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(goqu.I("tracks_to_tags.tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(goqu.I("tracks_to_tags.track_id"))

	query := dialect.From("tracks").
		Select(
			"tracks.id",
			"tracks.name",

			"tracks.album_id",
			"tracks.artist_id",

			"tracks.track_number",
			"tracks.duration",
			"tracks.year",

			"tracks.export_name",
			"tracks.original_filename",
			"tracks.mobile_filename",

			"tracks.created",
			"tracks.updated",

			goqu.I("albums.name").As("album_name"),
			goqu.I("albums.cover_art").As("album_cover_art"),
			goqu.I("artists.name").As("artist_name"),

			goqu.I("tags.tags").As("tags"),
		).
		Prepared(true).
		Join(
			goqu.I("albums"),
			goqu.On(goqu.I("tracks.album_id").Eq(goqu.I("albums.id"))),
		).
		Join(
			goqu.I("artists"),
			goqu.On(goqu.I("tracks.artist_id").Eq(goqu.I("artists.id"))),
		).
		LeftJoin(
			tags.As("tags"),
			goqu.On(goqu.I("tracks.id").Eq(goqu.I("tags.track_id"))),
		)

	return query
}

type Scan interface {
	Scan(dest ...any) error
}

var ErrInvalidFilter = errors.New("Invalid filter")
var ErrInvalidSort = errors.New("Invalid sort")

func (db *Database) GetAllTracks(ctx context.Context, filterStr string, sortStr string, includeAll bool) ([]Track, error) {
	query := TrackQuery()

	a := TrackResolverAdapter{}

	if filterStr != "" {
		re, err := fullParseFilter(&a, filterStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidFilter, err)
		}

		query = query.Where(re)
	}

	sortExpr, err := sort.Parse(sortStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidSort, err)
	}

	resolver := sort.New(&TrackResolverAdapter{})

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

	var items []Track
	err = db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetTracksByAlbum(ctx context.Context, albumId string) ([]Track, error) {
	query := TrackQuery().
		Where(goqu.I("tracks.album_id").Eq(albumId)).
		Order(goqu.I("tracks.track_number").Asc().NullsLast(), goqu.I("tracks.name").Asc())

	var items []Track
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetTrackById(ctx context.Context, id string) (Track, error) {
	query := TrackQuery().
		Where(goqu.I("tracks.id").Eq(id))

	var item Track
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Track{}, ErrItemNotFound
		}

		return Track{}, err
	}

	return item, nil
}

// func (db *Database) GetTrackByPath(ctx context.Context, path string) (Track, error) {
// 	ds := dialect.From("tracks").
// 		Select(
// 			"id",
// 			"track_number",
// 			"name",
// 			"cover_art",
// 			"duration",
// 			"path",
// 			"best_quality_file",
// 			"mobile_quality_file",
// 			"album_id",
// 			"artist_id",
// 		).
// 		Where(goqu.C("path").Eq(path)).
// 		Prepared(true)
//
// 	row, err := db.QueryRow(ctx, ds)
// 	if err != nil {
// 		return Track{}, err
// 	}
//
// 	var item Track
// 	err = row.Scan(
// 		&item.Id,
// 		&item.Number,
// 		&item.Name,
// 		&item.CoverArt,
// 		&item.Duration,
// 		&item.Path,
// 		&item.BestQualityFile,
// 		&item.MobileQualityFile,
// 		&item.AlbumId,
// 		&item.ArtistId,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return Track{}, ErrItemNotFound
// 		}
//
// 		return Track{}, err
// 	}
//
// 	return item, nil
// }

// func (db *Database) GetTrackByName(ctx context.Context, name string) (Track, error) {
// 	query := TrackQuery().
// 		Where(goqu.I("tracks.name").Eq(name))
//
// 	row, err := db.QueryRow(ctx, query)
// 	if err != nil {
// 		return Track{}, err
// 	}
//
// 	item, err := ScanTrack(row)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return Track{}, ErrItemNotFound
// 		}
//
// 		return Track{}, err
// 	}
//
// 	return item, nil
// }

func (db *Database) GetTrackByNameAndAlbum(ctx context.Context, name string, albumId string) (Track, error) {
	query := TrackQuery().
		Where(
			goqu.And(
				goqu.I("tracks.name").Eq(name),
				goqu.I("tracks.album_id").Eq(albumId),
			),
		)

	var item Track
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Track{}, ErrItemNotFound
		}

		return Track{}, err
	}

	return item, nil
}

type CreateTrackParams struct {
	Id   string
	Name string

	AlbumId  string
	ArtistId string

	Number   sql.NullInt64
	Duration sql.NullInt64
	Year     sql.NullInt64

	ExportName       string
	OriginalFilename string
	MobileFilename   string

	Created int64
	Updated int64
}

func (db *Database) CreateTrack(ctx context.Context, params CreateTrackParams) (string, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateId()
	}

	ds := dialect.Insert("tracks").Rows(goqu.Record{
		"id":   id,
		"name": params.Name,

		"album_id":  params.AlbumId,
		"artist_id": params.ArtistId,

		"track_number": params.Number,
		"duration":     params.Duration,
		"year":         params.Year,

		"export_name":       params.ExportName,
		"original_filename": params.OriginalFilename,
		"mobile_filename":   params.MobileFilename,

		"created": created,
		"updated": updated,
	}).
		Returning("id").
		Prepared(true)

	var item string

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return "", err
	}

	err = row.Scan(&item)
	if err != nil {
		return "", err
	}

	return item, nil
}

func addToRecord[T any](record goqu.Record, name string, change types.Change[T]) {
	if change.Changed {
		record[name] = change.Value
	}
}

func (db *Database) UpdateTrack(ctx context.Context, id string, changes TrackChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "album_id", changes.AlbumId)
	addToRecord(record, "artist_id", changes.ArtistId)

	addToRecord(record, "track_number", changes.Number)
	addToRecord(record, "duration", changes.Duration)
	addToRecord(record, "year", changes.Year)

	addToRecord(record, "export_name", changes.ExportName)
	addToRecord(record, "original_filename", changes.OriginalFilename)
	addToRecord(record, "mobile_filename", changes.MobileFilename)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("tracks").
		Set(record).
		Where(goqu.I("id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAlbumTracks(ctx context.Context, albumId string) error {
	query := dialect.Delete("tracks").
		Prepared(true).
		Where(goqu.I("tracks.album_id").Eq(albumId))

	// TODO(patrik): Check result?
	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveTrack(ctx context.Context, id string) error {
	query := dialect.Delete("tracks").
		Prepared(true).
		Where(goqu.I("tracks.id").Eq(id))

	// TODO(patrik): Check result?
	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

var tags []Tag

func (db *Database) Invalidate() {
	log.Debug("Database.Invalidate")

	ctx := context.Background()

	tags, _ = db.GetAllTags(ctx)
}

func TrackMapNameToId(typ string, name string) string {
	switch typ {
	case "tags":
		for _, t := range tags {
			if t.Name == strings.ToLower(name) {
				return t.Name
			}
		}
	}

	return ""
}
