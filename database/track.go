package database

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/dwebble/utils"
)

type Track struct {
	Id       string
	Number   int
	Name     string
	CoverArt string

	Path string

	BestQualityFile   string
	MobileQualityFile string

	AlbumId  string
	ArtistId string

	AlbumName  string
	ArtistName string
}

func (db *Database) GetAllTracks(ctx context.Context) ([]Track, error) {
	ds := dialect.From("tracks").
		Select("tracks.id", "tracks.track_number", "tracks.name", "tracks.cover_art", "tracks.path", "tracks.best_quality_file", "tracks.mobile_quality_file", "tracks.album_id", "tracks.artist_id", "albums.name", "artists.name").
		Join(goqu.I("albums"), goqu.On(goqu.I("tracks.album_id").Eq(goqu.I("albums.id")))).
		Join(goqu.I("artists"), goqu.On(goqu.I("tracks.artist_id").Eq(goqu.I("artists.id"))))

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Track
	for rows.Next() {
		var item Track
		err := rows.Scan(&item.Id, &item.Number, &item.Name, &item.CoverArt, &item.Path, &item.BestQualityFile, &item.MobileQualityFile, &item.AlbumId, &item.ArtistId, &item.AlbumName, &item.ArtistName)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (db *Database) GetTrackById(ctx context.Context, id string) (Track, error) {
	ds := dialect.From("tracks").
		Select("tracks.id", "tracks.track_number", "tracks.name", "tracks.cover_art", "tracks.path", "tracks.best_quality_file", "tracks.mobile_quality_file", "tracks.album_id", "tracks.artist_id", "albums.name", "artists.name").
		Join(goqu.I("albums"), goqu.On(goqu.I("tracks.album_id").Eq(goqu.I("albums.id")))).
		Join(goqu.I("artists"), goqu.On(goqu.I("tracks.artist_id").Eq(goqu.I("artists.id")))).
		Where(goqu.I("tracks.id").Eq(id)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Track{}, err
	}

	var item Track
	err = row.Scan(&item.Id, &item.Number, &item.Name, &item.CoverArt, &item.Path, &item.BestQualityFile, &item.MobileQualityFile, &item.AlbumId, &item.ArtistId, &item.AlbumName, &item.ArtistName)
	if err != nil {
		return Track{}, err
	}

	return item, nil
}

func (db *Database) GetTrackByPath(ctx context.Context, path string) (Track, error) {
	ds := dialect.From("tracks").
		Select("id", "track_number", "name", "cover_art", "path", "best_quality_file", "mobile_quality_file", "album_id", "artist_id").
		Where(goqu.C("path").Eq(path)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Track{}, err
	}

	var item Track
	err = row.Scan(&item.Id, &item.Number, &item.Name, &item.CoverArt, &item.Path, &item.BestQualityFile, &item.MobileQualityFile, &item.AlbumId, &item.ArtistId)
	if err != nil {
		return Track{}, err
	}

	return item, nil
}

type CreateTrackParams struct {
	TrackNumber       int
	Name              string
	CoverArt          string
	Path              string
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
		"path":                params.Path,
		"best_quality_file":   params.BestQualityFile,
		"mobile_quality_file": params.MobileQualityFile,
		"album_id":            params.AlbumId,
		"artist_id":           params.ArtistId,
	}).
		Returning("id", "track_number", "name", "cover_art", "path", "best_quality_file", "mobile_quality_file", "album_id", "artist_id").
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Track{}, err
	}

	var item Track
	err = row.Scan(&item.Id, &item.Number, &item.Name, &item.CoverArt, &item.Path, &item.BestQualityFile, &item.MobileQualityFile, &item.AlbumId, &item.ArtistId)
	if err != nil {
		return Track{}, err
	}

	return item, nil
}

func (db *Database) UpdateTrack(ctx context.Context, id, bestQualityFile, mobileQualityFile string) error {
	ds := dialect.Update("tracks").Set(goqu.Record{
		"best_quality_file":   bestQualityFile,
		"mobile_quality_file": mobileQualityFile,
	}).
	Where(goqu.C("id").Eq(id)).
	Prepared(true)

	tag, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	fmt.Printf("tag: %v\n", tag)

	return nil
}
