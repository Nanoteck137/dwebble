package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/core/log"
)

type SearchOptions struct {
	NumItems uint
}

var replacer = strings.NewReplacer(
	`"`, "",
)

func encodeSearchQuery(searchQuery string) (string, bool) {
	searchQuery = replacer.Replace(searchQuery)

	if searchQuery == "" {
		return "", false
	}

	searchQuery = "\"" + searchQuery + "\""

	return searchQuery, true
}

type SearchFunc[T any] func(db *Database, searchQuery string, options SearchOptions) ([]T, error)

type SearchFuncParams struct {
	Table     string
	ItemTable string

	ItemQuery *goqu.SelectDataset
}

func createSearchFunc[T any](params SearchFuncParams) SearchFunc[T] {
	return func(db *Database, searchQuery string, options SearchOptions) ([]T, error) {
		searchQuery, ok := encodeSearchQuery(searchQuery)
		if !ok {
			return nil, nil
		}

		tbl := goqu.T(params.Table)
		targetTbl := goqu.T(params.ItemTable)

		itemQuery := params.ItemQuery

		query := dialect.From(tbl).
			Select(
				targetTbl.Col("*"),
			).
			Prepared(true).
			Join(
				itemQuery.As(targetTbl.GetTable()),
				goqu.On(
					tbl.Col("id").Eq(targetTbl.Col("id")),
				),
			).
			Where(
				goqu.L("? MATCH ?", tbl, searchQuery),
			).
			Order(goqu.I("rank").Asc()).
			Limit(options.NumItems)

		sql, params, _ := query.ToSQL()
		fmt.Printf("sql: %v\n", sql)
		fmt.Printf("params: %v\n", params)

		var res []T
		err := db.Select(&res, query)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

var TestSearchArtist = createSearchFunc[Artist](SearchFuncParams{
	Table:     "artists_search",
	ItemTable: "artists",
	ItemQuery: ArtistQuery(),
})

var TestSearchAlbums = createSearchFunc[Album](SearchFuncParams{
	Table:     "albums_search",
	ItemTable: "albums",
	ItemQuery: AlbumQuery(),
})

var TestSearchTracks = createSearchFunc[Track](SearchFuncParams{
	Table:     "tracks_search",
	ItemTable: "tracks",
	ItemQuery: TrackQuery(),
})

func (db *Database) SearchArtists(searchQuery string, options SearchOptions) ([]Artist, error) {
	// searchQuery, ok := encodeSearchQuery(searchQuery)
	// if !ok {
	// 	return nil, nil
	// }
	//
	// tbl := goqu.T("artists_search")
	// targetTbl := goqu.T("artists")
	//
	// query := dialect.From(tbl).
	// 	Select(
	// 		targetTbl.Col("*"),
	// 	).
	// 	Prepared(true).
	// 	Join(
	// 		ArtistQuery().As(targetTbl.GetTable()),
	// 		goqu.On(
	// 			tbl.Col("id").Eq(targetTbl.Col("id")),
	// 		),
	// 	).
	// 	Where(
	// 		goqu.L("? MATCH ?", tbl, searchQuery),
	// 	).
	// 	Order(goqu.I("rank").Asc()).
	// 	Limit(options.NumItems)
	//
	// sql, params, _ := query.ToSQL()
	// fmt.Printf("sql: %v\n", sql)
	// fmt.Printf("params: %v\n", params)
	//
	// var res []Artist
	// err := db.Select(&res, query)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// return res, nil

	return TestSearchArtist(db, searchQuery, options)
}

func (db *Database) SearchAlbums(searchQuery string) ([]Album, error) {
	// searchQuery, ok := encodeSearchQuery(searchQuery)
	// if !ok {
	// 	return nil, nil
	// }
	//
	// query := dialect.From("albums_search").
	// 	Select("albums.*").
	// 	Prepared(true).
	// 	Join(
	// 		AlbumQuery().As("albums"),
	// 		goqu.On(goqu.I("albums_search.id").Eq(goqu.I("albums.id"))),
	// 	).
	// 	Where(
	// 		goqu.L("? MATCH ?", goqu.T("albums_search"), searchQuery),
	// 	).
	// 	Order(goqu.I("rank").Asc()).
	// 	Limit(10)
	//
	// var res []Album
	// err := db.Select(&res, query)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// return res, nil

	return TestSearchAlbums(db, searchQuery, SearchOptions{
		NumItems: 10,
	})
}

func (db *Database) SearchTracks(searchQuery string) ([]Track, error) {
	// searchQuery, ok := encodeSearchQuery(searchQuery)
	// if !ok {
	// 	return nil, nil
	// }
	//
	// query := dialect.From("tracks_search").
	// 	Select("tracks.*").
	// 	Prepared(true).
	// 	Join(
	// 		TrackQuery().As("tracks"),
	// 		goqu.On(goqu.I("tracks_search.id").Eq(goqu.I("tracks.id"))),
	// 	).
	// 	Where(
	// 		goqu.L("? MATCH ?", goqu.T("tracks_search"), searchQuery),
	// 	).
	// 	Order(goqu.I("rank").Asc()).
	// 	Limit(10)
	//
	// var res []Track
	// err := db.Select(&res, query)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// return res, nil

	return TestSearchTracks(db, searchQuery, SearchOptions{
		NumItems: 5,
	})
}

func (db *Database) InitializeSearch() error {
	_, err := db.Conn.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS artists_search USING fts5(
			id UNINDEXED, 
			name, 
			other_name,
			tags,
			tokenize=trigram
		);

		CREATE VIRTUAL TABLE IF NOT EXISTS albums_search USING fts5(
			id UNINDEXED, 

			name, 
			other_name,

			artist, 
			artist_other_name,

			tags,

			tokenize=trigram
		);

		CREATE VIRTUAL TABLE IF NOT EXISTS tracks_search USING fts5(
			id UNINDEXED, 

			name, 
			other_name,

			artist, 
			artist_other_name,

			album, 
			album_other_name,

			tags, 

			tokenize=trigram
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertArtistToSearch(ctx context.Context, data Artist) error {
	query := dialect.Insert("artists_search").Rows(goqu.Record{
		"rowid": data.RowId,

		"id": data.Id,

		"name":       data.Name,
		"other_name": data.OtherName,

		"tags": data.Tags.String,
	})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteArtistFromSearch(ctx context.Context, data Artist) error {
	query := dialect.Delete("artists_search").Where(
		goqu.I("artists_search.rowid").Eq(data.RowId),
	)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateSearchArtist(ctx context.Context, data Artist) error {
	query := dialect.Update("artists_search").
		Set(goqu.Record{
			"id": data.Id,

			"name":       data.Name,
			"other_name": data.OtherName,

			"tags": data.Tags.String,
		}).
		Where(goqu.I("artists_search.rowid").Eq(data.RowId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertAlbumToSearch(ctx context.Context, data Album) error {
	query := dialect.Insert("albums_search").Rows(goqu.Record{
		"rowid": data.RowId,

		"id": data.Id,

		"name":       data.Name,
		"other_name": data.OtherName,

		"artist":            data.ArtistName,
		"artist_other_name": data.ArtistOtherName,

		"tags": data.Tags.String,
	})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteAlbumFromSearch(ctx context.Context, data Album) error {
	query := dialect.Delete("albums_search").Where(
		goqu.I("albums_search.rowid").Eq(data.RowId),
	)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateSearchAlbum(ctx context.Context, data Album) error {
	query := dialect.Update("albums_search").
		Set(goqu.Record{
			"id": data.Id,

			"name":       data.Name,
			"other_name": data.OtherName,

			"artist":            data.ArtistName,
			"artist_other_name": data.ArtistOtherName,

			"tags": data.Tags.String,
		}).
		Where(goqu.I("albums_search.rowid").Eq(data.RowId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertTrackToSearch(ctx context.Context, data Track) error {
	query := dialect.Insert("tracks_search").Rows(goqu.Record{
		"rowid": data.RowId,

		"id": data.Id,

		"name":       data.Name,
		"other_name": data.OtherName,

		"artist":            data.ArtistName,
		"artist_other_name": data.ArtistOtherName,

		"album":            data.AlbumName,
		"album_other_name": data.AlbumOtherName,

		"tags": data.Tags.String,
	})

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteTrackFromSearch(ctx context.Context, data Track) error {
	query := dialect.Delete("tracks_search").Where(
		goqu.I("tracks_search.rowid").Eq(data.RowId),
	)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateSearchTrack(ctx context.Context, data Track) error {
	log.Debug("Updating track search", "track", data)

	query := dialect.Update("tracks_search").
		Set(goqu.Record{
			"id": data.Id,

			"name":       data.Name,
			"other_name": data.OtherName,

			"artist":            data.ArtistName,
			"artist_other_name": data.ArtistOtherName,

			"album":            data.AlbumName,
			"album_other_name": data.AlbumOtherName,

			"tags": data.Tags.String,
		}).
		Where(goqu.I("tracks_search.rowid").Eq(data.RowId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RefillSearchTables(ctx context.Context) error {
	_, err := db.Conn.Exec(`
		DROP TABLE IF EXISTS artists_search;
		DROP TABLE IF EXISTS albums_search;
		DROP TABLE IF EXISTS tracks_search;
	`)
	if err != nil {
		return err
	}

	err = db.InitializeSearch()
	if err != nil {
		return err
	}

	artists, err := db.GetAllArtists(ctx, "", "")
	if err != nil {
		return err
	}

	for _, artist := range artists {
		err := db.InsertArtistToSearch(ctx, artist)
		if err != nil {
			return err
		}
	}

	albums, err := db.GetAllAlbums(ctx, "", "")
	if err != nil {
		return err
	}

	for _, album := range albums {
		err := db.InsertAlbumToSearch(ctx, album)
		if err != nil {
			return err
		}
	}

	tracks, err := db.GetAllTracks(ctx, "", "")
	if err != nil {
		return err
	}

	for _, track := range tracks {
		err := db.InsertTrackToSearch(ctx, track)
		if err != nil {
			return err
		}
	}

	return nil
}
