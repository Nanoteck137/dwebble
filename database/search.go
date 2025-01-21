package database

import (
	"context"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/core/log"
)

type ArtistSearch struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TrackSearch struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Tags   string `json:"tags"`
}

var replacer = strings.NewReplacer(
	`"`, "",
)

func (db *Database) SearchArtists(searchQuery string) ([]Artist, error) {
	searchQuery = replacer.Replace(searchQuery)

	if searchQuery == "" {
		return nil, nil
	}

	searchQuery = "\"" + searchQuery + "\""

	query := dialect.From("artists_search").
		Select(
			"artists.id",
			"artists.name",
			"artists.picture",
		).
		Prepared(true).
		Join(
			ArtistQuery().As("artists"),
			goqu.On(goqu.I("artists_search.id").Eq(goqu.I("artists.id"))),
		).
		Where(
			goqu.L("? MATCH ?", goqu.T("artists_search"), searchQuery),
		).
		Order(goqu.I("rank").Asc()).
		Limit(10)

	var res []Artist
	err := db.Select(&res, query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *Database) SearchAlbums(searchQuery string) ([]Album, error) {
	searchQuery = replacer.Replace(searchQuery)

	if searchQuery == "" {
		return nil, nil
	}

	searchQuery = "\"" + searchQuery + "\""

	query := dialect.From("albums_search").Select(
		"albums.id",
		"albums.name",
		"albums.artist_id",

		"albums.cover_art",
		"albums.year",

		"albums.created",
		"albums.updated",

		"artist_name",
	).
		Prepared(true).
		Join(
			AlbumQuery().As("albums"),
			goqu.On(goqu.I("albums_search.id").Eq(goqu.I("albums.id"))),
		).
		Where(
			goqu.L("? MATCH ?", goqu.T("albums_search"), searchQuery),
		).
		Order(goqu.I("rank").Asc()).
		Limit(10)

	var res []Album
	err := db.Select(&res, query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *Database) SearchTracks(searchQuery string) ([]Track, error) {
	searchQuery = replacer.Replace(searchQuery)

	if searchQuery == "" {
		return nil, nil
	}

	searchQuery = "\"" + searchQuery + "\""

	query := dialect.From("tracks_search").
		Select("tracks.*").
		Prepared(true).
		Join(
			TrackQuery().As("tracks"),
			goqu.On(goqu.I("tracks_search.id").Eq(goqu.I("tracks.id"))),
		).
		Where(
			goqu.L("? MATCH ?", goqu.T("tracks_search"), searchQuery),
		).
		Order(goqu.I("rank").Asc()).
		Limit(10)

	var res []Track
	err := db.Select(&res, query)
	if err != nil {
		return nil, err
	}

	return res, nil
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
