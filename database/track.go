package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
)

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

type FilterAction int

const (
	FilterActionAdd FilterAction = iota
	FilterActionRemove
)

type FilterType int

const (
	FilterTypeTags FilterType = iota
	FilterTypeGenres
)

type Filter struct {
	Action FilterAction
	Type   FilterType
	Ids    []string
}

func (db *Database) GetAllTracks(ctx context.Context, random bool, filters []Filter) ([]Track, error) {
	// SELECT * FROM tracks WHERE
	// id IN (SELECT track_id FROM tracks_to_tags WHERE tag_id='ytnqxmqyo4plg5nhvxjezv2e8e9gw4jh')
	// AND id NOT IN (SELECT track_id FROM tracks_to_genres WHERE genre_id='y1vwuhgv7cfuab1pym13ox1smgvbluiw')
	// AND id NOT IN (SELECT track_id FROM tracks_to_genres WHERE genre_id='i83do0db648xt7oli8t9y3zj4ne2ud02')
	// ORDER BY name

	// filter = tags[Anime] + genre[Nerdcore] && tags[Jujutsu Kaisen] &~ genre[Soundtrack]

	type Set struct {
		Not      bool
		TagIds   []string
		GenreIds []string
	}

	sets := []Set{
		{
			Not:      false,
			TagIds:   []string{"ytnqxmqyo4plg5nhvxjezv2e8e9gw4jh"},
			GenreIds: []string{},
		},
		{
			Not:      true,
			TagIds:   []string{},
			GenreIds: []string{"y1vwuhgv7cfuab1pym13ox1smgvbluiw", "i83do0db648xt7oli8t9y3zj4ne2ud02"},
		},
	}

	// sets := []Set{
	// 	{
	// 		Not:      false,
	// 		TagIds:   []string{"ytnqxmqyo4plg5nhvxjezv2e8e9gw4jh"},
	// 		GenreIds: []string{"i83do0db648xt7oli8t9y3zj4ne2ud02"},
	// 	},
	// 	{
	// 		Not:      false,
	// 		TagIds:   []string{"bqgsf98akxhl6nc356ao7kbmw1l8pwp2"},
	// 		GenreIds: []string{},
	// 	},
	// 	{
	// 		Not:      true,
	// 		TagIds:   []string{},
	// 		GenreIds: []string{"y1vwuhgv7cfuab1pym13ox1smgvbluiw"},
	// 	},
	// }

	e := []exp.Expression{
		goqu.I("tracks.available").Eq(true),
	}

	for _, set := range sets {
		var sel *goqu.SelectDataset

		if set.TagIds != nil && len(set.TagIds) > 0 {
			sel = dialect.From("tracks_to_tags").
				Select("track_id").
				Where(goqu.I("tag_id").In(set.TagIds))
		}

		if set.GenreIds != nil && len(set.GenreIds) > 0 {
			s := dialect.From("tracks_to_genres").
				Select("track_id").
				Where(goqu.I("genre_id").In(set.GenreIds))

			if sel == nil {
				sel = s
			} else {
				sel = sel.Union(s)
			}
		}

		if sel != nil {
			if set.Not {
				e = append(e, goqu.L("tracks.id NOT IN ?", sel))
			} else {
				e = append(e, goqu.L("tracks.id IN ?", sel))
			}
		}
	}

	// var addFilter *goqu.SelectDataset
	// var removeFilter *goqu.SelectDataset
	// for _, filter := range filters {
	// 	switch filter.Action {
	// 	case FilterActionAdd:
	// 		var s *goqu.SelectDataset
	// 		switch filter.Type {
	// 		case FilterTypeTags:
	// 			s = dialect.From("tracks_to_tags").
	// 				Select("track_id").
	// 				Where(goqu.I("tag_id").In(filter.Ids))
	// 		case FilterTypeGenres:
	// 			s = dialect.From("tracks_to_genres").
	// 				Select("track_id").
	// 				Where(goqu.I("genre_id").In(filter.Ids))
	// 		}
	//
	// 		if addFilter == nil {
	// 			addFilter = s
	// 		} else {
	// 			addFilter = addFilter.Union(s)
	// 		}
	// 	case FilterActionRemove:
	// 		var s *goqu.SelectDataset
	// 		switch filter.Type {
	// 		case FilterTypeTags:
	// 			s = dialect.From("tracks_to_tags").
	// 				Select("track_id").
	// 				Where(goqu.I("tag_id").In(filter.Ids))
	// 		case FilterTypeGenres:
	// 			s = dialect.From("tracks_to_genres").
	// 				Select("track_id").
	// 				Where(goqu.I("genre_id").In(filter.Ids))
	// 		}
	//
	// 		if removeFilter == nil {
	// 			removeFilter = s
	// 		} else {
	// 			removeFilter = removeFilter.Union(s)
	// 		}
	// 	default:
	// 	}
	//
	// }

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
		Where(e...).
		Order(goqu.I("tracks.name").Asc())
		// goqu.I("tracks.available").Eq(true)
		// goqu.Ex{
		// 	"tracks.id": t,
		// 	"tracks.available": true,
		// }

	// e := []exp.Expression{
	// 	goqu.I("tracks.available").Eq(true),
	// }
	//
	// if addFilter != nil {
	// 	e = append(e, goqu.L("tracks.id IN ?", addFilter))
	// }
	//
	// if removeFilter != nil {
	// 	e = append(e, goqu.L("tracks.id NOT IN ?", removeFilter))
	// }
	//
	// ds = ds.Where(e...)

	if random {
		ds = ds.Order(goqu.Func("RANDOM").Asc())
	}

	ds = ds.Prepared(true)

	sql, params, _ := ds.ToSQL()

	fmt.Printf("sql: %v\n", sql)
	fmt.Printf("params: %v\n", params)

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
		Where(goqu.I("tracks.available").Eq(true)).
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

	fmt.Printf("tag: %v\n", tag)

	return nil
}
