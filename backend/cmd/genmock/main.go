package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dwebble/v2/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/nrednav/cuid2"
)

func generateMockData(queries *database.Queries) {
	ctx := context.Background()

	gen, err := cuid2.Init(cuid2.WithLength(32))
	if err != nil {
		log.Fatal(err)
	}

	for artistIndex := 0; artistIndex < 100; artistIndex++ {
		artistId := gen()

		err := queries.CreateArtist(ctx, database.CreateArtistParams{
			ID:      artistId,
			Name:    fmt.Sprintf("Test #%v", artistIndex+1),
			Picture: pgtype.Text{},
		})

		if err != nil {
			log.Printf("Failed to create artist %v", err)
		}

		for albumIndex := 0; albumIndex < 10; albumIndex++ {
			albumId := gen()
			err := queries.CreateAlbum(ctx, database.CreateAlbumParams{
				ID:       albumId,
				Name:     fmt.Sprintf("Album #%v", albumIndex+1),
				CoverArt: pgtype.Text{},
				ArtistID: artistId,
			})

			if err != nil {
				log.Printf("Failed to create album %v", err)
			}

			for songIndex := 0; songIndex < 15; songIndex++ {
				songId := gen()
				err := queries.CreateSong(ctx, database.CreateSongParams{
					ID:       songId,
					Name:     fmt.Sprintf("Song #%v", songIndex+1),
					CoverArt: pgtype.Text{},
					AlbumID:  albumId,
					ArtistID: artistId,
				})

				if err != nil {
					log.Printf("Failed to create song %v", err)
				}
			}
		}
	}
}

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

	db, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

	queries.DeleteAllSongs(ctx)
	queries.DeleteAllAlbums(ctx)
	queries.DeleteAllArtists(ctx)

	generateMockData(queries)
}
