package library

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type MetadataGeneral struct {
	Cover     string   `json:"cover" toml:"cover"`
	Tags      []string `json:"tags" toml:"tags"`
	TrackTags []string `json:"trackTags" toml:"trackTags"`
	Year      int64    `json:"year" toml:"year"`
}

type MetadataAlbum struct {
	Id      string   `json:"id" toml:"id"`
	Name    string   `json:"name" toml:"name"`
	Year    int64    `json:"year" toml:"year"`
	Tags    []string `json:"tags" toml:"tags"`
	Artists []string `json:"artists" toml:"artists"`
}

type MetadataTrack struct {
	Id      string   `json:"id" toml:"id"`
	File    string   `json:"file" toml:"file"`
	Name    string   `json:"name" toml:"name"`
	Number  int64    `json:"number" toml:"number"`
	Year    int64    `json:"year" toml:"year"`
	Tags    []string `json:"tags" toml:"tags"`
	Artists []string `json:"artists" toml:"artists"`
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

type SearchError struct {
	Path string
	Err  error
}

type Search struct {
	Albums []Album
	Errors map[string]error
}

func readAlbum(p string) (Album, error) {
	metadataPath := path.Join(p, "album.toml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return Album{}, err
	}

	var metadata Metadata
	err = toml.Unmarshal(data, &metadata)
	if err != nil {
		return Album{}, err
	}

	if metadata.General.Cover != "" {
		metadata.General.Cover = path.Join(p, metadata.General.Cover)
	}

	for i, t := range metadata.Tracks {
		metadata.Tracks[i].File = path.Join(p, t.File)
	}

	return Album{
		Path:     p,
		Metadata: metadata,
	}, nil
}

func FindAlbums(p string) (*Search, error) {
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

		if name == "album.toml" {
			albums = append(albums, path.Dir(p))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	errors := map[string]error{}
	res := make([]Album, 0, len(albums))

	for _, p := range albums {
		album, err := readAlbum(p)
		if err != nil {
			errors[p] = err
			continue
		}

		res = append(res, album)
	}

	return &Search{
		Albums: res,
		Errors: errors,
	}, nil
}

type FileNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"isDir"`
	Children []FileNode `json:"children,omitempty"`
}

func includeFile(name string) bool {
	return false
	// return name == "album.toml"
}

// BuildDirTree builds a directory tree starting from `root`
func BuildDirTree(root string) (FileNode, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return FileNode{}, err
	}

	return buildTreeRecursive(absRoot, absRoot)
}

func buildTreeRecursive(root, current string) (FileNode, error) {
	info, err := os.Stat(current)
	if err != nil {
		return FileNode{}, err
	}

	relPath, _ := filepath.Rel(root, current)
	if relPath == "." {
		relPath = "/" // root
	} else {
		relPath = "/" + filepath.ToSlash(relPath) // always start with /
	}

	node := FileNode{
		Name:  info.Name(),
		Path:  relPath,
		IsDir: info.IsDir(),
	}

	if info.IsDir() {
		entries, err := os.ReadDir(current)
		if err != nil {
			return node, err
		}

		for _, entry := range entries {
			childPath := filepath.Join(current, entry.Name())

			if entry.IsDir() || includeFile(entry.Name()) {
				childNode, err := buildTreeRecursive(root, childPath)
				if err != nil {
					return node, err
				}
				node.Children = append(node.Children, childNode)
			}
		}
	}

	return node, nil
}
