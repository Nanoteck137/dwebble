package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/nanoteck137/dwebble/v2/internal/database"
)

const defaultImagePath = "/images/default_artist.png"

type TrackDef struct {
	Id       string `json:"id"`
	Number   uint   `json:"number"`
	CoverArt string `json:"coverArt"`
	Name     string `json:"name"`
}

type AlbumDef struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	CoverArt string     `json:"coverArt"`
	Tracks   []TrackDef `json:"tracks"`
}

type ArtistDef struct {
	Id      string     `json:"id"`
	Name    string     `json:"name"`
	Changed int        `json:"changed"`
	Picture string     `json:"picture"`
	Albums  []AlbumDef `json:"albums"`
}

func updateArtist(queries *database.Queries, dbArtist database.Artist, artist ArtistDef) error {
	log.Printf("Updating '%v'\n", dbArtist.Name)
	return nil
}

func createArtist(queries *database.Queries, artist ArtistDef) error {
	ctx := context.Background()

	// TODO(patrik): Check for picture
	err := queries.CreateArtist(ctx, database.CreateArtistParams{
		ID:      artist.Id,
		Name:    artist.Name,
		Picture: defaultImagePath, // artist.Picture,
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, album := range artist.Albums {
		err := queries.CreateAlbum(ctx, database.CreateAlbumParams{
			ID:       album.Id,
			Name:     album.Name,
			CoverArt: defaultImagePath,
			ArtistID: artist.Id,
		})

		if err != nil {
			log.Fatal(err)
		}

		for _, track := range album.Tracks {
			err := queries.CreateTrack(ctx, database.CreateTrackParams{
				ID:          track.Id,
				TrackNumber: int32(track.Number),
				Name:        track.Name,
				CoverArt:    defaultImagePath,
				Filename:    "",
				AlbumID:     album.Id,
				ArtistID:    artist.Id,
			})

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	collection := os.Getenv("COLLECTION")
	fmt.Printf("collection: %v\n", collection)

	entries, err := os.ReadDir(collection)
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

	if len(os.Args) == 2 {
		queries.DeleteAllTracks(ctx)
		queries.DeleteAllAlbums(ctx)
		queries.DeleteAllArtists(ctx)
	}

	for _, entry := range entries {
		p := path.Join(collection, entry.Name())
		fmt.Printf("p: %v\n", p)

		data, err := os.ReadFile(path.Join(p, "artist.json"))
		if err != nil {
			log.Fatal(err)
		}

		var artist ArtistDef
		err = json.Unmarshal(data, &artist)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("artist: %+v\n", artist)

		dbArtist, err := queries.GetArtist(ctx, artist.Id)
		if err != nil {
			if err == pgx.ErrNoRows {
				log.Printf("Create artist '%v'\n", artist.Name)

				err := createArtist(queries, artist)
				if err != nil {
					log.Fatal(err)
				}

				continue
			}

			log.Fatal(err)
		}

		err = updateArtist(queries, dbArtist, artist)
		if err != nil {
			log.Fatal(err)
		}
	}
}
