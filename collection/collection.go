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
	Number   int    `json:"number"`
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
	Id       string
	Name     string
	FileName string
	AlbumId  string
	ArtistId string
}

type Album struct {
	Id       string
	Name     string
	TrackIds []string
	ArtistId string
}

type Artist struct {
	Id       string
	Name     string
	AlbumIds []string
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
			Id:       artist.Id,
			Name:     artist.Name,
			AlbumIds: []string{},
		}
	}

	for _, artistMetadata := range artists {
		artist, ok := artistMap[artistMetadata.Id]
		if !ok {
			log.Fatalf("No artist with id: '%v'\n", artistMetadata.Id)
		}

		for _, albumMetadata := range artistMetadata.Albums {
			album := &Album{
				Id:       albumMetadata.Id,
				Name:     albumMetadata.Name,
				TrackIds: []string{},
				ArtistId: artistMetadata.Id,
			}

			_, exists := albumMap[albumMetadata.Id]
			if exists {
				log.Fatalf("Duplicated album id: '%v'\n", albumMetadata.Id)
			}

			albumMap[albumMetadata.Id] = album

			artist.AlbumIds = append(artist.AlbumIds, albumMetadata.Id)

			for _, trackMetadata := range albumMetadata.Tracks {
				track := &Track{
					Id:       trackMetadata.Id,
					Name:     trackMetadata.Name,
					FileName: trackMetadata.FileName,
					AlbumId:  album.Id,
					ArtistId: artist.Id,
				}

				_, exists := trackMap[trackMetadata.Id]
				if exists {
					log.Fatalf("Duplicated track id: '%v'\n", albumMetadata.Id)
				}

				trackMap[trackMetadata.Id] = track
				album.TrackIds = append(album.TrackIds, track.Id)
			}
		}
	}

	return &Collection{
		Artists: artistMap,
		Albums:  albumMap,
		Tracks:  trackMap,
	}, nil
}
