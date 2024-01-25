package collection

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/kr/pretty"
)

type TrackMetadata struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"filePath"`
	Artist   string `json:"artist,omitempty"`
}

type AlbumMetadata struct {
	Id     string          `json:"id"`
	Name   string          `json:"name"`
	Tracks []TrackMetadata `json:"tracks"`
}

type ArtistMetadata struct {
	Id     string          `json:"id"`
	Name   string          `json:"name"`
	Albums []AlbumMetadata `json:"albums"`
}

type Artist struct {
	Name   string
	Albums []Album
}

type Album struct {
	Name   string
	Artist *Artist
}

type Collection struct {
	artists map[string]*Artist
}

func ReadEntryFromDir(dir string) (*ArtistMetadata, error) {
	artistFile := path.Join(dir, "artist.json")

	data, err := os.ReadFile(artistFile)
	if err != nil {
		return nil, err
	}

	var artistMetadata ArtistMetadata
	err = json.Unmarshal(data, &artistMetadata)
	if err != nil {
		return nil, err
	}

	return &artistMetadata, nil
}

func ReadFromDir(dir string) (*Collection, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var artists []*ArtistMetadata

	for _, entry := range entries {
		p := path.Join(dir, entry.Name())
		fmt.Printf("Path: %v\n", p)

		artist, err := ReadEntryFromDir(p)
		if err != nil {
			log.Printf("Warning: %v\n", err)
		}

		artists = append(artists, artist)
	}

	artistsMap := make(map[string]*Artist)

	for _, artist := range artists {
		if _, ok := artistsMap[artist.Id]; ok {
			// TODO(patrik): Better error message
			log.Fatalf("Duplicated artist ids: %v\n", artist.Id)
		}

		artistsMap[artist.Id] = &Artist{
			Name:   artist.Name,
			Albums: []Album{},
		}
	}

	for _, artist := range artists {
		a, ok := artistsMap[artist.Id]
		if !ok {
			log.Fatalf("No artist with id: '%v'\n", artist.Id)
		}

		for _, album := range artist.Albums {
			a.Albums = append(a.Albums, Album{
				Name:   album.Name,
				Artist: a,
			})
		}
	}

	pretty.Println(artistsMap)

	return nil, nil
}
