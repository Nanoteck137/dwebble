package sync

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
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
			sql, params, err := goqu.Update("artists").Set(goqu.Record{
				"name": v.Name,
			}).Where(goqu.C("id").Eq(k)).ToSQL()
			if err != nil {
				return err
			}

			pretty.Println(sql)
			pretty.Println(params)

			tag, err := db.Exec(ctx, sql)
			fmt.Println(tag.String())

			_ = dbArtist
		}
	}

	return nil
}
