package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"log"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/dwebble/filter"
	"github.com/nanoteck137/dwebble/filter/gen"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
)

type TrackResolverAdapter struct {
}

func (a *TrackResolverAdapter) MapNameToId(typ, name string) (string, error) {
	switch typ {
	case "tags":
		for _, t := range tags {
			if t.Name == strings.ToLower(name) {
				return t.Name, nil
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
}

func (db *Database) GetAllTracks(ctx context.Context, filterStr string) ([]Track, error) {
	ds := dialect.From("tracks").
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
		).
		Prepared(true).
		Join(goqu.I("albums"), goqu.On(goqu.I("tracks.album_id").Eq(goqu.I("albums.id")))).
		Join(goqu.I("artists"), goqu.On(goqu.I("tracks.artist_id").Eq(goqu.I("artists.id")))).
		Order(goqu.I("tracks.name").Asc())

	a := TrackResolverAdapter{}

	if filterStr != "" {
		ast, err := parser.ParseExpr(filterStr)
		if err != nil {
			return nil, err
		}

		r := filter.New(&a)
		e, err := r.Resolve(ast)
		if err != nil {
			return nil, err
		}

		re, err := gen.Generate(e)
		if err != nil {
			return nil, err
		}

		ds = ds.Where(goqu.I("tracks.available").Eq(true), re)
	} else {
		ds = ds.Where(goqu.I("tracks.available").Eq(true))
	}

	// order := "sort=+artist,+track"
	order := ""
	split := strings.Split(order, "=")

	fmt.Printf("split: %v\n", split)

	mode := split[0]
	switch mode {
	case "sort":
		args := strings.Split(split[1], ",")

		var orderExprs []exp.OrderedExpression

		fmt.Printf("args: %v\n", args)
		for _, arg := range args {
			switch arg[0] {
			case '+':
				name, err := a.MapName(arg[1:])
				if err != nil {
					return nil, err
				}

				fmt.Printf("name: %v\n", name)
				orderExprs = append(orderExprs, goqu.I(name.Name).Asc())
			case '-':
			default:
			}
		}

		ds = ds.Order(orderExprs...)
	case "random":
		ds = ds.Order(goqu.Func("RANDOM").Asc())
	default:
		ds = ds.Order(goqu.I("tracks.name").Asc())
	}

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Track
	for rows.Next() {
		var item Track
		err := rows.Scan(
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
			&item.AlbumName,
			&item.ArtistName,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetTracksByAlbum(ctx context.Context, albumId string) ([]Track, error) {
	ds := dialect.From("tracks").
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
		).
		Join(goqu.I("albums"), goqu.On(goqu.I("tracks.album_id").Eq(goqu.I("albums.id")))).
		Join(goqu.I("artists"), goqu.On(goqu.I("tracks.artist_id").Eq(goqu.I("artists.id")))).
		Where(goqu.And(goqu.I("tracks.available").Eq(true), goqu.I("tracks.album_id").Eq(albumId))).
		Order(goqu.I("tracks.track_number").Asc(), goqu.I("tracks.name").Asc()).
		Prepared(true)

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Track
	for rows.Next() {
		var item Track
		err := rows.Scan(
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
			&item.AlbumName,
			&item.ArtistName,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetTrackById(ctx context.Context, id string) (Track, error) {
	ds := dialect.From("tracks").
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
		).
		Join(goqu.I("albums"), goqu.On(goqu.I("tracks.album_id").Eq(goqu.I("albums.id")))).
		Join(goqu.I("artists"), goqu.On(goqu.I("tracks.artist_id").Eq(goqu.I("artists.id")))).
		Where(goqu.I("tracks.id").Eq(id)).
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
		&item.AlbumName,
		&item.ArtistName,
	)
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
	ds := dialect.From("tracks").
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
		).
		Join(goqu.I("albums"), goqu.On(goqu.I("tracks.album_id").Eq(goqu.I("albums.id")))).
		Join(goqu.I("artists"), goqu.On(goqu.I("tracks.artist_id").Eq(goqu.I("artists.id")))).
		Where(goqu.I("tracks.name").Eq(name)).
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
		&item.AlbumName,
		&item.ArtistName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Track{}, types.ErrNoTrack
		}

		return Track{}, err
	}

	return item, nil
}

func (db *Database) GetTrackByNameAndAlbum(ctx context.Context, name string, albumId string) (Track, error) {
	ds := dialect.From("tracks").
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
		).
		Join(goqu.I("albums"), goqu.On(goqu.I("tracks.album_id").Eq(goqu.I("albums.id")))).
		Join(goqu.I("artists"), goqu.On(goqu.I("tracks.artist_id").Eq(goqu.I("artists.id")))).
		Where(
			goqu.And(
				goqu.I("tracks.name").Eq(name),
				goqu.I("tracks.album_id").Eq(albumId),
			),
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
		&item.AlbumName,
		&item.ArtistName,
	)
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

	tag, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	// TODO(patrik): Cleanup
	fmt.Printf("tag: %v\n", tag)

	return nil
}

func (db *Database) MarkAllTracksUnavailable(ctx context.Context) error {
	ds := dialect.Update("tracks").Set(goqu.Record{
		"available": false,
	})

	tag, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	// TODO(patrik): Cleanup
	fmt.Printf("tag: %v\n", tag)

	return nil
}

var tags []Tag
var genres []Genre

func (db *Database) Invalidate() {
	log.Printf("Database.Invalidate")

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
