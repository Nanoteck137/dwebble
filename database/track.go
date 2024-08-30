package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go/ast"
	"strings"

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
	case "available":
		return types.Name{
			Kind: types.NameKindNumber,
			Name: "tracks.available",
		}, nil
	}

	return types.Name{}, sort.UnknownName(name)

}

func (a *TrackResolverAdapter) MapNameToId(typ, name string) (string, error) {
	switch typ {
	case "tags":
		for _, t := range tags {
			if t.Name == strings.ToLower(name) {
				return t.Id, nil
			}
		}

		return "", nil
	case "genres":
		for _, g := range genres {
			if g.Name == strings.ToLower(name) {
				return g.Id, nil
			}
		}

		return "", nil
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
			WhereName:  "tag_id",
		}, nil
	case "genres":
		return filter.Table{
			Name:       "tracks_to_genres",
			SelectName: "track_id",
			WhereName:  "genre_id",
		}, nil
	}

	return filter.Table{}, fmt.Errorf("Unknown table type: %s", typ)
}

func (a *TrackResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	switch name {
	case "hasTag":
		return resolver.InTable(name, "tags", args)
	case "hasGenre":
		return resolver.InTable(name, "genres", args)
	}

	return nil, fmt.Errorf("Unknown function name: %s", name)
}

type Track struct {
	Id       string         `db:"id"`
	Number   int            `db:"track_number"`
	Name     string         `db:"name"`
	CoverArt sql.NullString `db:"cover_art"`
	Duration int            `db:"duration"`

	Path string `db:"path"`

	BestQualityFile   string `db:"best_quality_file"`
	MobileQualityFile string `db:"mobile_quality_file"`

	AlbumId  string `db:"album_id"`
	ArtistId string `db:"artist_id"`

	AlbumName  string `db:"album_name"`
	ArtistName string `db:"artist_name"`

	Available bool `db:"available"`

	Tags   sql.NullString `db:"tags"`
	Genres sql.NullString `db:"genres"`
}

type TrackChanges struct {
	Number            types.Change[int]
	Name              types.Change[string]
	CoverArt          types.Change[sql.NullString]
	Duration          types.Change[int]
	BestQualityFile   types.Change[string]
	MobileQualityFile types.Change[string]
	ArtistId          types.Change[string]
	Available         bool
}

func TrackQuery() *goqu.SelectDataset {
	tags := dialect.From("tracks_to_tags").
		Select(
			goqu.I("tracks_to_tags.track_id").As("track_id"),
			goqu.Func("group_concat", goqu.I("tags.display_name"), ",").As("tags"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(goqu.I("tracks_to_tags.tag_id").Eq(goqu.I("tags.id"))),
		).
		GroupBy(goqu.I("tracks_to_tags.track_id"))

	genres := dialect.From("tracks_to_genres").
		Select(
			goqu.I("tracks_to_genres.track_id").As("track_id"),
			goqu.Func("group_concat", goqu.I("genres.display_name"), ",").As("genres"),
		).
		Join(
			goqu.I("genres"),
			goqu.On(goqu.I("tracks_to_genres.genre_id").Eq(goqu.I("genres.id"))),
		).
		GroupBy(goqu.I("tracks_to_genres.track_id"))

	query := dialect.From("tracks").
		Select(
			"tracks.id",
			"tracks.track_number",
			"tracks.name",
			"tracks.cover_art",
			"tracks.duration",
			"tracks.path",
			"tracks.best_quality_file",
			"tracks.mobile_quality_file",
			"tracks.album_id",
			"tracks.artist_id",
			"tracks.available",
			goqu.I("albums.name").As("album_name"),
			goqu.I("artists.name").As("artist_name"),
			goqu.I("tags.tags").As("tags"),
			goqu.I("genres.genres").As("genres"),
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
		).
		LeftJoin(
			genres.As("genres"),
			goqu.On(goqu.I("tracks.id").Eq(goqu.I("genres.track_id"))),
		)

	return query
}

type Scan interface {
	Scan(dest ...any) error
}

func ScanTrack(scanner Scan) (Track, error) {
	var res Track
	err := scanner.Scan(
		&res.Id,
		&res.Number,
		&res.Name,
		&res.CoverArt,
		&res.Duration,
		&res.Path,
		&res.BestQualityFile,
		&res.MobileQualityFile,
		&res.AlbumId,
		&res.ArtistId,
		&res.Available,
		&res.AlbumName,
		&res.ArtistName,
		&res.Tags,
		&res.Genres,
	)
	if err != nil {
		return Track{}, err
	}

	return res, nil
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

		if includeAll {
			query = query.Where(re)
		} else {
			query = query.Where(re, goqu.I("tracks.available").Eq(true))
		}
	} else {
		if !includeAll {
			query = query.Where(goqu.I("tracks.available").Eq(true))
		}
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
		Where(goqu.I("tracks.album_id").Eq(albumId))

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

func (db *Database) GetTrackByPath(ctx context.Context, path string) (Track, error) {
	ds := dialect.From("tracks").
		Select(
			"id",
			"track_number",
			"name",
			"cover_art",
			"duration",
			"path",
			"best_quality_file",
			"mobile_quality_file",
			"album_id",
			"artist_id",
		).
		Where(goqu.C("path").Eq(path)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Track{}, err
	}

	var item Track
	err = row.Scan(
		&item.Id,
		&item.Number,
		&item.Name,
		&item.CoverArt,
		&item.Duration,
		&item.Path,
		&item.BestQualityFile,
		&item.MobileQualityFile,
		&item.AlbumId,
		&item.ArtistId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Track{}, ErrItemNotFound
		}

		return Track{}, err
	}

	return item, nil
}

func (db *Database) GetTrackByName(ctx context.Context, name string) (Track, error) {
	query := TrackQuery().
		Where(goqu.I("tracks.name").Eq(name))

	row, err := db.QueryRow(ctx, query)
	if err != nil {
		return Track{}, err
	}

	item, err := ScanTrack(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Track{}, ErrItemNotFound
		}

		return Track{}, err
	}

	return item, nil
}

func (db *Database) GetTrackByNameAndAlbum(ctx context.Context, name string, albumId string) (Track, error) {
	query := TrackQuery().
		Where(
			goqu.And(
				goqu.I("tracks.name").Eq(name),
				goqu.I("tracks.album_id").Eq(albumId),
			),
		)

	row, err := db.QueryRow(ctx, query)
	if err != nil {
		return Track{}, err
	}

	item, err := ScanTrack(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Track{}, ErrItemNotFound
		}

		return Track{}, err
	}

	return item, nil
}

type CreateTrackParams struct {
	TrackNumber       int
	Name              string
	CoverArt          sql.NullString
	Path              string
	Duration          int
	BestQualityFile   string
	MobileQualityFile string
	AlbumId           string
	ArtistId          string
}

func (db *Database) CreateTrack(ctx context.Context, params CreateTrackParams) (Track, error) {
	ds := dialect.Insert("tracks").Rows(goqu.Record{
		"id":                  utils.CreateId(),
		"track_number":        params.TrackNumber,
		"name":                params.Name,
		"cover_art":           params.CoverArt,
		"duration":            params.Duration,
		"best_quality_file":   params.BestQualityFile,
		"mobile_quality_file": params.MobileQualityFile,
		"album_id":            params.AlbumId,
		"artist_id":           params.ArtistId,
		"path":                params.Path,
		"available":           false,
	}).
		Returning(
			"id",
			"track_number",
			"name",
			"cover_art",
			"duration",
			"path",
			"best_quality_file",
			"mobile_quality_file",
			"album_id",
			"artist_id",
		).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Track{}, err
	}

	var item Track
	err = row.Scan(
		&item.Id,
		&item.Number,
		&item.Name,
		&item.CoverArt,
		&item.Duration,
		&item.Path,
		&item.BestQualityFile,
		&item.MobileQualityFile,
		&item.AlbumId,
		&item.ArtistId,
	)
	if err != nil {
		return Track{}, err
	}

	return item, nil
}


func addToRecord[T any](record goqu.Record, name string, change types.Change[T]) {
	if change.Changed {
		record[name] = change.Value
	}
}

func (db *Database) UpdateTrack(ctx context.Context, id string, changes TrackChanges) error {
	record := goqu.Record{
		"available": changes.Available,
	}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "track_number", changes.Number)
	addToRecord(record, "name", changes.Name)
	addToRecord(record, "cover_art", changes.CoverArt)
	addToRecord(record, "duration", changes.Duration)
	addToRecord(record, "best_quality_file", changes.BestQualityFile)
	addToRecord(record, "mobile_quality_file", changes.MobileQualityFile)
	addToRecord(record, "artist_id", changes.ArtistId)

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

func (db *Database) MarkAllTracksUnavailable(ctx context.Context) error {
	ds := dialect.Update("tracks").Set(goqu.Record{
		"available": false,
	})

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

var tags []Tag
var genres []Genre

func (db *Database) Invalidate() {
	log.Debug("Database.Invalidate")

	ctx := context.Background()

	tags, _ = db.GetAllTags(ctx)
	genres, _ = db.GetAllGenres(ctx)
}

func TrackMapNameToId(typ string, name string) string {
	switch typ {
	case "tags":
		for _, t := range tags {
			if t.Name == strings.ToLower(name) {
				return t.Name
			}
		}
	case "genres":
		for _, g := range genres {
			if g.Name == strings.ToLower(name) {
				return g.Id
			}
		}
	}

	return ""
}
