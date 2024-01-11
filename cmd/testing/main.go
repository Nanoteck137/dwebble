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
	"github.com/nanoteck137/dwebble/v2/collection"
	"github.com/nanoteck137/dwebble/v2/internal/database"
)

const defaultImagePath = "/images/default_artist.png"

func updateArtist(queries *database.Queries, dbArtist database.Artist, artist collection.ArtistDef) error {
	log.Printf("Updating '%v'\n", dbArtist.Name)
	return nil
}

func createArtist(workDir string, queries *database.Queries, artistDir string, artist collection.ArtistDef) error {
	ctx := context.Background()

	picture := defaultImagePath

	// pretty.Println(artist)
	// if artist.Picture != "" {
	// 	ext := path.Ext(artist.Picture)[1:]
	// 	typ := "picture"
	// 	name := fmt.Sprintf("%v.%v.%v", artist.Id, typ, ext)
	// 	dst := path.Join(workDir, name)
	//
	// 	fmt.Printf("dst: %v\n", dst)
	//
	// 	a, err := filepath.Abs(path.Join(artistDir, artist.Picture))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	err = os.Symlink(a, dst)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	picture = fmt.Sprintf("/images/%v", name)
	// }

	// TODO(patrik): Check for picture
	err := queries.CreateArtist(ctx, database.CreateArtistParams{
		ID:      artist.Id,
		Name:    artist.Name,
		Picture: picture, 
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

	collectionPath := os.Getenv("COLLECTION")
	fmt.Printf("collection: %v\n", collectionPath)

	entries, err := os.ReadDir(collectionPath)
	if err != nil {
		log.Fatal(err)
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL not set")
	}

	workingDir := "./work"
	err = os.MkdirAll(workingDir, 0755)
	if err != nil {
		log.Fatal(err)
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
		p := path.Join(collectionPath, entry.Name())
		fmt.Printf("p: %v\n", p)

		data, err := os.ReadFile(path.Join(p, "artist.json"))
		if err != nil {
			log.Fatal(err)
		}

		var artist collection.ArtistDef
		err = json.Unmarshal(data, &artist)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("artist: %+v\n", artist)

		dbArtist, err := queries.GetArtist(ctx, artist.Id)
		if err != nil {
			if err == pgx.ErrNoRows {
				log.Printf("Create artist '%v'\n", artist.Name)

				err := createArtist(workingDir, queries, p, artist)
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
