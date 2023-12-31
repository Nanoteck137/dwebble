package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/nanoteck137/dwebble/v2/internal/database"
	"github.com/nrednav/cuid2"
)

const defaultImagePath = "/images/default_artist.png"

func generateMockData(queries *database.Queries, dir string) {
	ctx := context.Background()

	gen, err := cuid2.Init(cuid2.WithLength(32))
	if err != nil {
		log.Fatal(err)
	}

	for artistIndex := 0; artistIndex < 100; artistIndex++ {
		artistId := gen()

		err := queries.CreateArtist(ctx, database.CreateArtistParams{
			ID:      artistId,
			Name:    fmt.Sprintf("Artist #%v", artistIndex+1),
			Picture: defaultImagePath,
		})

		if err != nil {
			log.Printf("Failed to create artist %v", err)
			continue
		}

		for albumIndex := 0; albumIndex < 4; albumIndex++ {
			albumId := gen()

			err := queries.CreateAlbum(ctx, database.CreateAlbumParams{
				ID:       albumId,
				Name:     fmt.Sprintf("Album #%v", albumIndex+1),
				CoverArt: defaultImagePath,
				ArtistID: artistId,
			})

			if err != nil {
				log.Printf("Failed to create album %v", err)
				continue
			}

			for trackIndex := 0; trackIndex < 8; trackIndex++ {
				trackId := gen()

				err := queries.CreateTrack(ctx, database.CreateTrackParams{
					ID:          trackId,
					TrackNumber: int32(trackIndex) + 1,
					Name:        fmt.Sprintf("Track #%v", trackIndex+1),
					CoverArt:    defaultImagePath,
					Filename:    "123.flac",
					AlbumID:     albumId,
					ArtistID:    artistId,
				})

				if err != nil {
					log.Printf("Failed to create track %v", err)
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

	dir := "test"

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	queries.DeleteAllTracks(ctx)
	queries.DeleteAllAlbums(ctx)
	queries.DeleteAllArtists(ctx)

	generateMockData(queries, dir)
}
