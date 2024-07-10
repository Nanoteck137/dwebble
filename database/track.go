package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/database/filtergen"
	"github.com/nanoteck137/dwebble/log"
	"github.com/nanoteck137/dwebble/sortfilter/filter"
	"github.com/nanoteck137/dwebble/sortfilter/sort"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
)

type TrackResolverAdapter struct {
}

func (a *TrackResolverAdapter) GetDefaultSort() string {
	return "tracks.name"
}

func (a *TrackResolverAdapter) MapNameToId(typ, name string) (string, error) {
	switch typ {
	case "tags":
		for _, t := range tags {
			if t.Name == strings.ToLower(name) {
				return t.Id, nil
			}
		}
	case "genres":
		for _, g := range genres {
			if g.Name == strings.ToLower(name) {
				return g.Id, nil
			}
		}
	}

	return "", fmt.Errorf("Unknown name type: %s (%s)", typ, name)
}

func (a *TrackResolverAdapter) MapName(name string) (filter.Name, error) {
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
	case "track":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "tracks.name",
		}, nil
	case "trackNumber":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "tracks.track_number",
		}, nil
	}

	return filter.Name{}, fmt.Errorf("Unknown name: %s", name)
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
	Id       string
	Number   int
	Name     string
	CoverArt sql.NullString
	Duration int

	Path string

	BestQualityFile   string
	MobileQualityFile string

	AlbumId  string
	ArtistId string

	AlbumName  string
	ArtistName string

	Tags   string
	Genres string
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
			"albums.name",
			"artists.name",
			"tags.tags",
			"genres.genres",
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
		Join(
			tags.As("tags"),
			goqu.On(goqu.I("tracks.id").Eq(goqu.I("tags.track_id"))),
		).
		Join(
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

func (db *Database) GetAllTracks(ctx context.Context, filterStr string, sortStr string) ([]Track, error) {
	query := TrackQuery()

	a := TrackResolverAdapter{}

	if filterStr != "" {
		re, err := FullParseFilter(&a, filterStr)
		if err != nil {
			return nil, err
		}

		query = query.Where(goqu.I("tracks.available").Eq(true), re)
	} else {
		query = query.Where(goqu.I("tracks.available").Eq(true))
	}

	if sortStr != "" {
		e, err := sort.Parse(sortStr)
		if err != nil {
			return nil, err
		}

		r := sort.New(&a)

		re, err := r.Resolve(e)

		ge, err := filtergen.GenerateSort(re)
		if err != nil {
			return nil, err
		}

		query = query.Order(ge...)
	} else {
		query = query.Order(goqu.I("tracks.name").Asc())
	}

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var items []Track
	for rows.Next() {
		item, err := ScanTrack(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetTracksByAlbum(ctx context.Context, albumId string) ([]Track, error) {
	query := TrackQuery().
		Where(goqu.I("tracks.album_id").Eq(albumId))
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var items []Track
	for rows.Next() {
		item, err := ScanTrack(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetTrackById(ctx context.Context, id string) (Track, error) {
	query := TrackQuery().
		Where(goqu.I("tracks.id").Eq(id))

	row, err := db.QueryRow(ctx, query)
	if err != nil {
		return Track{}, err
	}

	item, err := ScanTrack(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return Track{}, types.ErrNoTrack
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
			return Track{}, types.ErrNoTrack
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
			return Track{}, types.ErrNoTrack
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
			return Track{}, types.ErrNoTrack
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

func (db *Database) UpdateTrack(ctx context.Context, id string, changes TrackChanges) error {
	record := goqu.Record{
		"available": changes.Available,
	}

	if changes.Number.Changed {
		record["track_number"] = changes.Number.Value
	}

	if changes.Name.Changed {
		record["name"] = changes.Name.Value
	}

	if changes.CoverArt.Changed {
		record["cover_art"] = changes.CoverArt.Value
	}

	if changes.Duration.Changed {
		record["duration"] = changes.Duration.Value
	}

	if changes.BestQualityFile.Changed {
		record["best_quality_file"] = changes.BestQualityFile.Value
	}

	if changes.MobileQualityFile.Changed {
		record["mobile_quality_file"] = changes.MobileQualityFile.Value
	}

	if changes.ArtistId.Changed {
		record["artist_id"] = changes.ArtistId.Value
	}

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
