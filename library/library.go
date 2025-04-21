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
	Cover     string   `toml:"cover"`
	Tags      []string `toml:"tags"`
	TrackTags []string `toml:"trackTags"`
	Year      int      `toml:"year"`
}

type MetadataAlbum struct {
	Name    string   `toml:"name"`
	Year    int      `toml:"year"`
	Tags    []string `toml:"tags"`
	Artists []string `toml:"artists"`
}

type MetadataTrack struct {
	File    string   `toml:"file"`
	Name    string   `toml:"name"`
	Number  int      `toml:"number"`
	Year    int      `toml:"year"`
	Tags    []string `toml:"tags"`
	Artists []string `toml:"artists"`
}

type Metadata struct {
	General MetadataGeneral `toml:"general"`
	Album   MetadataAlbum   `toml:"album"`
	Tracks  []MetadataTrack `toml:"tracks"`
}

type Library struct {
	Albums []Metadata
}

func FindAlbums(p string) ([]Metadata, error) {
	var albums []string

	err := filepath.WalkDir(p, func(p string, d fs.DirEntry, err error) error {
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

	res := make([]Metadata, 0, len(albums))

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
		}

		pretty.Println(metadata)

		res = append(res, metadata)
	}

	return res, nil
}
