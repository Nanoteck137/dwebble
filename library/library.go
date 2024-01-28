package library

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/utils"
)

type Track struct {
	Name     string
	Number   int
	FileName string
}

type Album struct {
	Name       string
	ArtistName string
	Tracks     []Track
}

type Artist struct {
	Name        string
	PicturePath string
	Albums      []Album
}

type Library struct {
	Artists map[string]*Artist
}

func ReadFromFS(fsys fs.FS) (*Library, error) {
	reg := regexp.MustCompile(`(\d*).*`)
	_ = reg

	dir := "."
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		log.Fatal(err)
	}


	getAllTrackFromDir := func(fss fs.FS, dir string) ([]Track, error) {
		entries, err := fs.ReadDir(fss, dir)
		if err != nil {
			return nil, err
		}

		var tracks []Track

		for _, entry := range entries {
			p := path.Join(dir, entry.Name())

			if utils.IsMusicFile(p) {
				ext := path.Ext(p)
				name := strings.TrimSuffix(entry.Name(), ext)
				// reg.FindAllStringSubmatchIndex

				captures := reg.FindStringSubmatch(entry.Name())
				if captures == nil {
					continue
				}

				if captures[1] == "" {
					continue
				}

				trackNum, err := strconv.Atoi(captures[1])
				if err != nil {
					log.Fatal(err)
				}

				tracks = append(tracks, Track{
					Name:     name,
					Number:   trackNum,
					FileName: p,
				})
			}
		}

		return tracks, nil
	}

	artists := make(map[string]*Artist)

	for _, entry := range entries {
		p := path.Join(dir, entry.Name())
		if entry.IsDir() {
			entries, err := fs.ReadDir(fsys, p)
			if err != nil {
				log.Fatal(err)
			}

			if !utils.HasMusic(entries) && !utils.IsMultiDisc(entries) {
				name := entry.Name()
				artists[name] = &Artist{
					Name: name,
				}
			}
		}
	}

	pretty.Println(artists)

	for _, entry := range entries {
		p := path.Join(dir, entry.Name())
		if entry.IsDir() {
			fmt.Printf("Valid: %v\n", p)

			entries, err := fs.ReadDir(fsys, p)
			if err != nil {
				log.Fatal(err)
			}

			artistName := "Various Artists"

			if utils.HasMusic(entries) {
				fmt.Printf("%v is an album\n", p)

				tracks, err := getAllTrackFromDir(fsys, p)
				if err != nil {
					log.Fatal(err)
				}

				artist := artists[artistName]

				albumName := entry.Name()
				artist.Albums = append(artist.Albums, Album{
					Name:       albumName,
					ArtistName: artistName,
					Tracks:     tracks,
				})
			} else if utils.IsMultiDisc(entries) {
				fmt.Printf("%v is an multidisc album\n", p)

				// albumName := entry.Name()
				// albums = append(albums, Album{
				// 	Name:       albumName,
				// 	ArtistName: artistName,
				// 	Tracks:     []Track{},
				// })
			} else {
				artistName = entry.Name()
				fmt.Printf("%v is an artist\n", p)

				for _, entry := range entries {
					p := path.Join(p, entry.Name())
					fmt.Printf("p: %v\n", p)

					if !entry.IsDir() {
						continue
					}

					entries, err := fs.ReadDir(fsys, p)
					if err != nil {
						log.Fatal(err)
					}

					if utils.HasMusic(entries) {
						tracks, err := getAllTrackFromDir(fsys, p)
						if err != nil {
							log.Fatal(err)
						}

						artist := artists[artistName]

						albumName := entry.Name()
						artist.Albums = append(artist.Albums, Album{
							Name:       albumName,
							ArtistName: artistName,
							Tracks:     tracks,
						})
					} else if utils.IsMultiDisc(entries) {
						// albumName := entry.Name()
						// albums = append(albums, Album{
						// 	Name:       albumName,
						// 	ArtistName: artistName,
						// 	Tracks:     []Track{},
						// })
					} else {
						log.Printf("Warning: No music found at '%v'", p)
					}
				}
			}

			fmt.Printf("artistName: %v\n", artistName)
			fmt.Println()
		}
	}

	// pretty.Println(albums)

	return &Library{
		Artists: artists,
	}, nil
}

func GetOrCreateArtist(queries *database.Queries, name string) (database.Artist, error) {
	ctx := context.Background()

	artists, err := queries.GetArtistByName(ctx, name)
	if err != nil {
		return database.Artist{}, err
	}

	var artist database.Artist

	if len(artists) == 0 {
		a, err := queries.CreateArtist(ctx, database.CreateArtistParams{
			ID:      utils.CreateId(),
			Name:    name,
			Picture: "TODO",
		})

		if err != nil {
			return database.Artist{}, err
		}

		artist = a
	} else if len(artists) == 1 {
		artist = artists[0]
	} else {
		return database.Artist{}, fmt.Errorf("Returned multiple artists for '%v'", name)
	}

	return artist, nil
}

func (lib *Library) Sync(db *pgxpool.Pool) error {
	queries := database.New(db)
	ctx := context.Background()

	_ = ctx
	_ = queries

	// for _, album := range lib.Albums {
	// 	name := album.ArtistName
	// 	artist, err := GetOrCreateArtist(queries, name)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	dbAlbum, err := queries.CreateAlbum(ctx, database.CreateAlbumParams{
	// 		ID:       utils.CreateId(),
	// 		Name:     album.Name,
	// 		CoverArt: "TODO",
	// 		ArtistID: artist.ID,
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	for _, track := range album.Tracks {
	// 		_, err := queries.CreateTrack(ctx, database.CreateTrackParams{
	// 			ID:                utils.CreateId(),
	// 			TrackNumber:       int32(track.Number),
	// 			Name:              track.Name,
	// 			CoverArt:          "TODO",
	// 			BestQualityFile:   "TODO",
	// 			MobileQualityFile: "TODO",
	// 			AlbumID:           dbAlbum.ID,
	// 			ArtistID:          artist.ID,
	// 		})
	//
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}
