package collection

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

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

type Collection struct {
	artists []ArtistDef
}

func Read(collectionPath string) (*Collection, error) {
	entries, err := os.ReadDir(collectionPath)
	if err != nil {
		return nil, err
	}

	artists := []ArtistDef{}

	// TODO(patrik): Check for duplicated arists name
	for _, entry := range entries {
		p := path.Join(collectionPath, entry.Name())
		_ = p

		data, err := os.ReadFile(path.Join(p, "artist.json"))
		if err != nil {
			log.Printf("Skipping '%v': %v\n", p, err)
			continue
		}

		var artistDef ArtistDef
		err = json.Unmarshal(data, &artistDef)
		if err != nil {
			log.Printf("Failed to parse json '%v': %v\n", p, err)
			// TODO(patrik): What to do here?
			return nil, err
		}

		artists = append(artists, artistDef)
	}

	return &Collection{artists: artists}, nil
}

func (col *Collection) FindArtistByName(name string) *ArtistDef {
	for _, artist := range col.artists {
		if artist.Name == name {
			return &artist
		}
	}

	return nil
}
