package sync

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/collection"
	"github.com/nanoteck137/dwebble/database"
)

// rootdir/
//   [Artist]/
//     [Albums]/
//     artist.json

func SyncCollection(col *collection.Collection, db *pgxpool.Pool) error {
	queries := database.New(db)

	ctx := context.Background()

	artists, err := queries.GetAllArtists(ctx)
	if err != nil {
		return err
	}

	pretty.Println(artists)

	for k, v := range col.Artists {
		dbArtist, err := queries.GetArtist(ctx, k)
		if err != nil {
			if err == pgx.ErrNoRows {
				_, err := queries.CreateArtist(ctx, database.CreateArtistParams{
					ID:      k,
					Name:    v.Name,
					Picture: "TODO",
				})

				if err != nil {
					return err
				}

				continue
			}

			return err
		} else {
			_ = dbArtist
		}
	}

	return nil
}
