package sync

import (
	"context"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
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


var dialect = goqu.Dialect("postgres")

type ArtistUpdateRequest struct {
	Id string `goqu:"skipupdate"`
	Name *string `db:"name" goqu:"omitnil"`
}

func RunArtistUpdate(ctx context.Context, db *pgxpool.Pool, req *ArtistUpdateRequest) error {
	sql, params, err := dialect.Update("artists").
		Set(req).
		Where(goqu.C("id").Eq(req.Id)).
		Prepared(true).ToSQL()
	if err != nil {
		return err
	}

	pretty.Println(sql)
	pretty.Println(params)

	_, err = db.Exec(ctx, sql, params...)

	return nil
}

func SyncCollection(col *collection.Collection, db *pgxpool.Pool) error {
	queries := database.New(db)

	ctx := context.Background()

	// artists, err := queries.GetAllArtists(ctx)
	// if err != nil {
	// 	return err
	// }
	//
	// pretty.Println(artists)

	sql, _, err := dialect.From("artists").Select(goqu.C("id")).ToSQL()
	if err != nil {
		return err
	}

	rows, err := db.Query(ctx, sql)
	if err != nil {
		return err
	}

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}

		ids = append(ids, id)
	}

	pretty.Println(ids)

	for artistId, artist := range col.Artists {
		dbArtist, err := queries.GetArtist(ctx, artistId)
		if err != nil {
			if err == pgx.ErrNoRows {
				_, err := queries.CreateArtist(ctx, database.CreateArtistParams{
					ID:      artistId,
					Name:    artist.Name,
					Picture: "TODO",
				})

				if err != nil {
					return err
				}

				continue
			}

			return err
		} else {
			updateReq := ArtistUpdateRequest{}
			if dbArtist.Name != artist.Name {
				updateReq.Name = &artist.Name
			}

			if !reflect.DeepEqual(updateReq, ArtistUpdateRequest{}) {
				RunArtistUpdate(ctx, db, &updateReq)
			}

			_ = dbArtist
		}
	}

	for albumId, album := range col.Albums {
		dbAlbum, err := queries.GetAlbum(ctx, albumId)
		if err != nil {
			if err == pgx.ErrNoRows {
				_, err := queries.CreateAlbum(ctx, database.CreateAlbumParams{
					ID:       albumId,
					Name:     album.Name,
					CoverArt: "TODO",
					ArtistID: album.ArtistId,
				})

				if err != nil {
					return err
				}

				continue
			}

			return err
		} else {
			// updateReq := ArtistUpdateRequest{}
			// if dbArtist.Name != artist.Name {
			// 	updateReq.Name = &artist.Name
			// }
			//
			// if !reflect.DeepEqual(updateReq, ArtistUpdateRequest{}) {
			// 	RunArtistUpdate(ctx, db, &updateReq)
			// }
			//
			_ = dbAlbum
		}

	}

	return nil
}
