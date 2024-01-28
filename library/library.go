package library

import (
	"fmt"
	"io/fs"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/kr/pretty"
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

type RawTrack struct {
	Name     string
	Number   int
	FileName string
}

type RawAlbum struct {
	Name       string
	ArtistName string
	Tracks     []RawTrack
}

type Library struct {
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

	getAllTrackFromDir := func(fss fs.FS, dir string) ([]RawTrack, error) {
		entries, err := fs.ReadDir(fss, dir)
		if err != nil {
			return nil, err
		}

		var tracks []RawTrack

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

				tracks = append(tracks, RawTrack{
					Name:     name,
					Number:   trackNum,
					FileName: p,
				})
			}
		}

		return tracks, nil
	}

	var albums []RawAlbum

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
				albums = append(albums, RawAlbum{
					Name:       albumName,
					ArtistName: artistName,
					Tracks:     tracks,
				})
			} else if isMultiDisc(entries) {
				fmt.Printf("%v is an multidisc album\n", p)

				albumName := entry.Name()
				albums = append(albums, RawAlbum{
					Name:       albumName,
					ArtistName: artistName,
					Tracks:     []RawTrack{},
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
						albums = append(albums, RawAlbum{
							Name:       albumName,
							ArtistName: artistName,
							Tracks:     tracks,
						})
					} else if isMultiDisc(entries) {
						albumName := entry.Name()
						albums = append(albums, RawAlbum{
							Name:       albumName,
							ArtistName: artistName,
							Tracks:     []RawTrack{},
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

	pretty.Println(albums)

	return &Library{}, nil
}
