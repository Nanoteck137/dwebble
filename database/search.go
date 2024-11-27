package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
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

func (db *Database) SearchArtists(searchQuery string) ([]Artist, error) {
	query := dialect.From("artists_search").Select(
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

func (db *Database) SearchTracks(searchQuery string) ([]Track, error) {
	query := dialect.From("tracks_search").Select(
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

		"tracks.album_name",
		"tracks.album_cover_art",
		"tracks.artist_name",

		"tracks.tags",
	).
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
		CREATE VIRTUAL TABLE IF NOT EXISTS artists_search USING fts5(id UNINDEXED, name, tokenize=trigram);
		CREATE VIRTUAL TABLE IF NOT EXISTS albums_search USING fts5(id UNINDEXED, name, artist, tokenize=trigram);
		CREATE VIRTUAL TABLE IF NOT EXISTS tracks_search USING fts5(id UNINDEXED, name, artist, album, tags, tokenize=trigram);
	`)
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

	artists, err := db.GetAllArtists(ctx)
	if err != nil {
		return err
	}

	for _, artist := range artists {
		query := dialect.Insert("artists_search").Rows(goqu.Record{
			"id":   artist.Id,
			"name": artist.Name,
		})

		_, err = db.Exec(ctx, query)
		if err != nil {
			return err
		}
	}

	tracks, err := db.GetAllTracks(ctx, "", "", true)
	if err != nil {
		return err
	}

	for _, track := range tracks {
		query := dialect.Insert("tracks_search").Rows(goqu.Record{
			"id":     track.Id,
			"name":   track.Name,
			"artist": track.ArtistName,
			"album":  track.AlbumName,
			"tags":   track.Tags.String,
		})

		_, err = db.Exec(ctx, query)
		if err != nil {
			return err
		}
	}

	return nil
}
