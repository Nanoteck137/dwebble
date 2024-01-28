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
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/utils"
)

var validExts = []string{
	"flac",
	"opus",
	"mp3",
}

func isMusicFile(p string) bool {
	ext := path.Ext(p)
	if ext == "" {
		return false
	}

	ext = strings.ToLower(ext[1:])

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}

	return false
}

var validMultiDiscPrefix = []string{
	"cd",
	"disc",
}

func isMultiDiscName(name string) bool {
	for _, prefix := range validMultiDiscPrefix {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

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

type Library struct {
	Albums []Album
}

func ReadFromFS(fsys fs.FS) (*Library, error) {
	reg := regexp.MustCompile(`(\d*).*`)
	_ = reg

	dir := "."
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		log.Fatal(err)
	}

	hasMusic := func(entries []fs.DirEntry) bool {
		for _, entry := range entries {
			if isMusicFile(entry.Name()) {
				return true
			}
		}

		return false
	}

	isMultiDisc := func(entries []fs.DirEntry) bool {
		for _, entry := range entries {
			if isMultiDiscName(entry.Name()) {
				return true
			}
		}
		return false
	}

	getAllTrackFromDir := func(fss fs.FS, dir string) ([]Track, error) {
		entries, err := fs.ReadDir(fss, dir)
		if err != nil {
			return nil, err
		}

		var tracks []Track

		for _, entry := range entries {
			p := path.Join(dir, entry.Name())

			if isMusicFile(p) {
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

	var albums []Album

	for _, entry := range entries {
		p := path.Join(dir, entry.Name())
		if entry.IsDir() {
			fmt.Printf("Valid: %v\n", p)

			entries, err := fs.ReadDir(fsys, p)
			if err != nil {
				log.Fatal(err)
			}

			artistName := "Various Artists"

			if hasMusic(entries) {
				fmt.Printf("%v is an album\n", p)

				tracks, err := getAllTrackFromDir(fsys, p)
				if err != nil {
					log.Fatal(err)
				}

				albumName := entry.Name()
				albums = append(albums, Album{
					Name:       albumName,
					ArtistName: artistName,
					Tracks:     tracks,
				})
			} else if isMultiDisc(entries) {
				fmt.Printf("%v is an multidisc album\n", p)

				albumName := entry.Name()
				albums = append(albums, Album{
					Name:       albumName,
					ArtistName: artistName,
					Tracks:     []Track{},
				})
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

					if hasMusic(entries) {
						tracks, err := getAllTrackFromDir(fsys, p)
						if err != nil {
							log.Fatal(err)
						}

						albumName := entry.Name()
						albums = append(albums, Album{
							Name:       albumName,
							ArtistName: artistName,
							Tracks:     tracks,
						})
					} else if isMultiDisc(entries) {
						albumName := entry.Name()
						albums = append(albums, Album{
							Name:       albumName,
							ArtistName: artistName,
							Tracks:     []Track{},
						})
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
		Albums: albums,
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

	for _, album := range lib.Albums {
		_ = album

		name := album.ArtistName
		artist, err := GetOrCreateArtist(queries, name)
		if err != nil {
			return err
		}

		dbAlbum, err := queries.CreateAlbum(ctx, database.CreateAlbumParams{
			ID:       utils.CreateId(),
			Name:     album.Name,
			CoverArt: "TODO",
			ArtistID: artist.ID,
		})
		if err != nil {
			return err
		}

		for _, track := range album.Tracks {
			_, err := queries.CreateTrack(ctx, database.CreateTrackParams{
				ID:                utils.CreateId(),
				TrackNumber:       int32(track.Number),
				Name:              track.Name,
				CoverArt:          "TODO",
				BestQualityFile:   "TODO",
				MobileQualityFile: "TODO",
				AlbumID:           dbAlbum.ID,
				ArtistID:          artist.ID,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}
