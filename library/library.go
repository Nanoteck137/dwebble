package library

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kr/pretty"
	"github.com/pelletier/go-toml/v2"
)

type MetadataGeneral struct {
	Cover     string   `json:"cover" toml:"cover"`
	Tags      []string `json:"tags" toml:"tags"`
	TrackTags []string `json:"trackTags" toml:"trackTags"`
	Year      int      `json:"year" toml:"year"`
}

type MetadataAlbum struct {
	Id      string   `json:"id" toml:"id"`
	Name    string   `json:"name" toml:"name"`
	Year    int      `json:"year" toml:"year"`
	Tags    []string `json:"tags" toml:"tags"`
	Artists []string `json:"artists" toml:"artists"`
}

type MetadataTrack struct {
	Id      string   `json:"id" toml:"id"`
	File    string   `json:"file" toml:"file"`
	Name    string   `json:"name" toml:"name"`
	Number  int      `json:"number" toml:"number"`
	Year    int      `json:"year" toml:"year"`
	Tags    []string `json:"tags" toml:"tags"`
	Artists []string `json:"artists" toml:"artists"`

	ModifiedTime int64 `json:"-" toml:"-"`
}

type Metadata struct {
	General MetadataGeneral `json:"general" toml:"general"`
	Album   MetadataAlbum   `json:"album" toml:"album"`
	Tracks  []MetadataTrack `json:"tracks" toml:"tracks"`
}

type Album struct {
	Path     string
	Metadata Metadata
}

type Library struct {
	Albums []Metadata
}

func FindAlbums(p string) ([]Album, error) {
	var albums []string

	err := filepath.WalkDir(p, func(p string, d fs.DirEntry, err error) error {
		if d == nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		name := d.Name()

		if strings.HasPrefix(name, ".") {
			return nil
		}

		fmt.Printf("p: %v\n", p)

		if name == "album.toml" {
			albums = append(albums, path.Dir(p))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	pretty.Println(albums)

	res := make([]Album, 0, len(albums))

	for _, p := range albums {
		metadataPath := path.Join(p, "album.toml")
		data, err := os.ReadFile(metadataPath)
		if err != nil {
			return nil, err
		}

		var metadata Metadata
		err = toml.Unmarshal(data, &metadata)
		if err != nil {
			return nil, err
		}

		if metadata.General.Cover != "" {
			metadata.General.Cover = path.Join(p, metadata.General.Cover)
		}

		for i, t := range metadata.Tracks {
			metadata.Tracks[i].File = path.Join(p, t.File)

			info, err := os.Stat(metadata.Tracks[i].File)
			if err != nil {
				return nil, err
			}

			metadata.Tracks[i].ModifiedTime = info.ModTime().UnixMilli()
		}

		pretty.Println(metadata)

		res = append(res, Album{
			Path:     p,
			Metadata: metadata,
		})
	}

	return res, nil
}
