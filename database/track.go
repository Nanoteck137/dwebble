package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
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
