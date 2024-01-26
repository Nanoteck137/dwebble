package collection

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

type TrackMetadata struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	FileName string `json:"fileName"`
	ArtistId string `json:"artistId,omitempty"`
}

type AlbumMetadata struct {
	Id     string          `json:"id"`
	Name   string          `json:"name"`
	Dir    string          `json:"dir"`
	Tracks []TrackMetadata `json:"tracks"`
}

type ArtistMetadata struct {
	Id     string          `json:"id"`
	Name   string          `json:"name"`
	Albums []AlbumMetadata `json:"albums"`
}

type Track struct {
	Name     string
	FileName string
	Album    *Album
	Artist   *Artist
}

type Album struct {
	Name     string
	Tracks   []*Track
	ArtistId string
}

type Artist struct {
	Name   string
	Albums []*Album
}

type Collection struct {
	Artists map[string]*Artist
	Albums  map[string]*Album
	Tracks  map[string]*Track
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
			continue
		}

		artists = append(artists, artist)
	}

	artistMap := make(map[string]*Artist)
	albumMap := make(map[string]*Album)
	trackMap := make(map[string]*Track)

	for _, artist := range artists {
		if _, ok := artistMap[artist.Id]; ok {
			// TODO(patrik): Better error message
			log.Fatalf("Duplicated artist ids: %v\n", artist.Id)
		}

		artistMap[artist.Id] = &Artist{
			Name:   artist.Name,
			Albums: []*Album{},
		}
	}

	for _, artistMetadata := range artists {
		artist, ok := artistMap[artistMetadata.Id]
		if !ok {
			log.Fatalf("No artist with id: '%v'\n", artistMetadata.Id)
		}

		for _, albumMetadata := range artistMetadata.Albums {
			album := &Album{
				Name:     albumMetadata.Name,
				ArtistId: artistMetadata.Id,
				Tracks:   []*Track{},
			}

			_, exists := albumMap[albumMetadata.Id]
			if exists {
				log.Fatalf("Duplicated album id: '%v'\n", albumMetadata.Id)
			}

			albumMap[albumMetadata.Id] = album

			artist.Albums = append(artist.Albums, album)

			for _, trackMetadata := range albumMetadata.Tracks {
				track := &Track{
					Name:     trackMetadata.Name,
					FileName: trackMetadata.FileName,
					Album:    album,
					Artist:   artist,
				}

				_, exists := trackMap[trackMetadata.Id]
				if exists {
					log.Fatalf("Duplicated track id: '%v'\n", albumMetadata.Id)
				}

				trackMap[trackMetadata.Id] = track
				album.Tracks = append(album.Tracks, track)
			}
		}
	}

	return &Collection{
		Artists: artistMap,
		Albums:  albumMap,
		Tracks:  trackMap,
	}, nil
}
