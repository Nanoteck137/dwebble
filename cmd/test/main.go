package main

import (
	"context"
	"log"
	"os"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/library"
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

	queries := database.New(db)

	queries.DeleteAllTracks(ctx)
	queries.DeleteAllAlbums(ctx)
	queries.DeleteAllArtists(ctx)

	fsys := os.DirFS("/Volumes/media/music")
	lib, err := library.ReadFromFS(fsys)
	if err != nil {
		log.Fatal(err)
	}

	_ = lib

	err = lib.Sync(db)
	if err != nil {
		log.Fatal(err)
	}

	// pretty.Println(lib)

	pretty.Println(queries.GetAllArtists(ctx))

}
