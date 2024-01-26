package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nanoteck137/dwebble/collection"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/sync"
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

	// queries.DeleteAllTracks(ctx)
	// queries.DeleteAllAlbums(ctx)
	// queries.DeleteAllArtists(ctx)

	col, err := collection.ReadFromDir("./mockdata")
	if err != nil {
		log.Fatal(err)
	}

	_ = col

	err = sync.SyncCollection(col, db)
	if err != nil {
		log.Fatal(err)
	}

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
