package database

import "context"

type Album struct {
	Id       string
	Name     string
	CoverArt string
	ArtistId string
	Path     string
}

func (db *Database) GetAllAlbums(ctx context.Context) ([]Album, error) {
	ds := dialect.From("albums").Select("id", "name", "cover_art", "artist_id", "path")

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	var items []Album
	for rows.Next() {
		var item Album
		err := rows.Scan(&item.Id, &item.Name, &item.CoverArt, &item.ArtistId, &item.Path)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
