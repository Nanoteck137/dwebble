package collection

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/gosimple/slug"
	"github.com/nanoteck137/dwebble/v2/utils"
)

var NotFoundErr = fmt.Errorf("Not found")

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

func (artist *ArtistDef) FindAlbumByName(name string) (*AlbumDef, error) {
	for i, album := range artist.Albums {
		if album.Name == name {
			return &artist.Albums[i], nil
		}
	}

	return nil, NotFoundErr
}

func (artist *ArtistDef) CreateAlbum(name string) *AlbumDef {
	artist.Albums = append(artist.Albums, AlbumDef{
		Id:       utils.CreateId(),
		Name:     name,
		CoverArt: "",
		Tracks:   []TrackDef{},
	})

	return &artist.Albums[len(artist.Albums)-1]
}

type collectionEntry struct {
	def ArtistDef

	originalPath string
	created      bool
}

type Collection struct {
	path string
	entries []collectionEntry
}

func Read(collectionPath string) (*Collection, error) {
	entries, err := os.ReadDir(collectionPath)
	if err != nil {
		return nil, err
	}

	colEntries := []collectionEntry{}

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

		colEntries = append(colEntries, collectionEntry{
			def:          artistDef,
			originalPath: p,
			created:      false,
		})
	}

	return &Collection{path: collectionPath, entries: colEntries}, nil
}

func (col *Collection) FindArtistByName(name string) (*ArtistDef, error) {
	for i, entry := range col.entries {
		if entry.def.Name == name {
			return &col.entries[i].def, nil
		}
	}

	return nil, NotFoundErr
}

func (col *Collection) CreateArtist(name string) *ArtistDef {
	col.entries = append(col.entries, collectionEntry{
		def: ArtistDef{
			Id:      utils.CreateId(),
			Name:    name,
			Changed: 0,
			Picture: "",
			Albums:  []AlbumDef{},
		},
		created: true,
	})

	return &col.entries[len(col.entries)-1].def
}

func (col *Collection) Flush() error {
	for _, entry := range col.entries {
		if entry.created {
			name := slug.Make(entry.def.Name)
			p := path.Join(col.path, name)
			err := os.MkdirAll(p, 0755)
			if err != nil {
				// TODO(patrik): Should we return errors in the loop?
				return err
			}

			data, err := json.MarshalIndent(entry.def, "", "  ")
			if err != nil {
				// TODO(patrik): Should we return errors in the loop?
				return err
			}

			err = os.WriteFile(path.Join(p, "artist.json"), data, 0644)
			if err != nil {
				// TODO(patrik): Should we return errors in the loop?
				return err
			}
		} else {
			// TODO(patrik): Only write dirty data

			data, err := json.MarshalIndent(entry.def, "", "  ")
			if err != nil {
				// TODO(patrik): Should we return errors in the loop?
				return err
			}

			err = os.WriteFile(path.Join(entry.originalPath, "artist.json"), data, 0644)
			if err != nil {
				// TODO(patrik): Should we return errors in the loop?
				return err
			}
		}
	}

	return nil
}
