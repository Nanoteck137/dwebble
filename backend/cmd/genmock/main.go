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
	gen, err := cuid2.Init(cuid2.WithLength(32))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		artistId := gen()

		queries.CreateArtist(context.Background(), database.CreateArtistParams{
			ID:      artistId,
			Name:    fmt.Sprintf("Test #%v", i + 1),
			Picture: pgtype.Text{},
		})
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

	db, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

    queries.DeleteAllArtists(context.Background());

    generateMockData(queries)
}
