package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nanoteck137/dwebble/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL not set")
	}

	ctx := context.Background()

	db, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	_ = db

	queries := database.New(db)
	_ = queries

	queries.DeleteAllTracks(ctx)
	queries.DeleteAllAlbums(ctx)
	queries.DeleteAllArtists(ctx)

	artist, err := queries.CreateArtist(ctx, database.CreateArtistParams{
		ID:      "test",
		Name:    "test",
		Picture: "TODO",
	})

	if err != nil {
		log.Fatal(err)
	}

	album, err := queries.CreateAlbum(ctx, database.CreateAlbumParams{
		ID:       "test",
		Name:     "test",
		CoverArt: "TODO",
		ArtistID: artist.ID,
	})

	if err != nil {
		log.Fatal(err)
	}

	var tracks []database.Track
	for i := 0; i < 5; i++ {
		track, err := queries.CreateTrack(ctx, database.CreateTrackParams{
			ID:                "track" + strconv.Itoa(i + 1),
			TrackNumber:       int32(i) + 1,
			Name:              "track" + strconv.Itoa(i + 1),
			AlbumID:           album.ID,
			ArtistID:          artist.ID,
		})

		if err != nil {
			log.Fatal(err)
		}

		tracks = append(tracks, track)
	}

	// tx, err := db.BeginTx(ctx, pgx.TxOptions{ DeferrableMode: pgx.NotDeferrable })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	dialect := goqu.Dialect("postgres")
	for i := 0; i < len(tracks); i++ {
		track := tracks[i]
		
		sql, _, err := dialect.Update("tracks").Set(goqu.Record{
			"track_number": len(tracks) - i,
		}).Where(goqu.C("id").Eq(track.ID)).ToSQL()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(sql)
		tag, err := db.Exec(ctx, sql)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Tag", tag)
	}

	// err = tx.Commit(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// col, err := collection.ReadFromDir("./mockdata")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// _ = col

	// err = sync.SyncCollection(col, db)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return

	// otherartist := collection.ArtistMetadata{
	// 	Id:     utils.CreateId(),
	// 	Name:   "Some Other Artist",
	// 	Albums: []collection.AlbumMetadata{},
	// }

	// artist := collection.ArtistMetadata{
	// 	Id:   utils.CreateId(),
	// 	Name: "Ado",
	// 	Albums: []collection.AlbumMetadata{
	// 		{
	// 			Id:   utils.CreateId(),
	// 			Name: "unravel",
	// 			Tracks: []collection.TrackMetadata{
	// 				{
	// 					Id:       utils.CreateId(),
	// 					Name:     "unravel",
	// 					FilePath: "unravel/01 - unravel.flac",
	// 				},
	// 			},
	// 		},
	// 		{
	// 			Id:   utils.CreateId(),
	// 			Name: "Kura Kura",
	// 			Tracks: []collection.TrackMetadata{
	// 				{
	// 					Id:       utils.CreateId(),
	// 					Name:     "Kura Kura",
	// 					FilePath: "Kura Kura/01. クラクラ.flac",
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	//
	// data, err := json.MarshalIndent(artist, "", "  ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Println(string(data))
	//
	// err = os.WriteFile("../testcol/Ado/artist.json", data, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// col := collection.New()
	//
	// data, err := json.Marshal(col)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Println(string(data))
}
