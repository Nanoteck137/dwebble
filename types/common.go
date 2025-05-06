package types

import "path"

type MediaType string

const (
	MediaTypeUnknown   MediaType = "unknown"
	MediaTypeFlac      MediaType = "flac"
	MediaTypeOggOpus   MediaType = "ogg-opus"
	MediaTypeOggVorbis MediaType = "ogg-vorbis"
	MediaTypeMp3       MediaType = "mp3"
	MediaTypeAcc       MediaType = "acc"
)

func GetMediaTypeFromExt(ext string) MediaType {
	switch ext {
	case ".flac":
		return MediaTypeFlac
	case ".opus":
		return MediaTypeOggOpus
	case ".ogg":
		return MediaTypeOggVorbis
	case ".mp3":
		return MediaTypeMp3
	case ".acc":
		return MediaTypeAcc
	}

	return MediaTypeUnknown
}

func (m MediaType) ToExt() (string, bool) {
	switch m {
	case MediaTypeFlac:
		return ".flac", true
	case MediaTypeOggOpus:
		return ".opus", true
	case MediaTypeOggVorbis:
		return ".ogg", true
	case MediaTypeMp3:
		return ".mp3", true
	case MediaTypeAcc:
		return ".acc", true
	}

	return "", false
}

type Map map[string]any

type WorkDir string

func (d WorkDir) String() string {
	return string(d)
}

func (d WorkDir) DatabaseFile() string {
	return path.Join(d.String(), "data.db")
}

func (d WorkDir) ExportFile() string {
	return path.Join(d.String(), "export.json")
}

func (d WorkDir) SetupFile() string {
	return path.Join(d.String(), "setup")
}

func (d WorkDir) Trash() string {
	return path.Join(d.String(), "trash")
}

func (d WorkDir) Artists() string {
	return path.Join(d.String(), "artists")
}

func (d WorkDir) Artist(id string) string {
	return path.Join(d.Artists(), id)
}

func (d WorkDir) Albums() string {
	return path.Join(d.String(), "albums")
}

func (d WorkDir) Album(id string) string {
	return path.Join(d.Albums(), id)
}

func (d WorkDir) Tracks() string {
	return path.Join(d.String(), "tracks")
}

func (d WorkDir) Track(id string) string {
	return path.Join(d.Tracks(), id)
}

func (d WorkDir) Cache() CacheDir {
	return CacheDir(path.Join(d.String(), "cache"))
}

type CacheDir string

func (d CacheDir) String() string {
	return string(d)
}

func (d CacheDir) Tracks() string {
	return path.Join(d.String(), "tracks")
}

func (d CacheDir) Track(id string) string {
	return path.Join(d.Tracks(), id)
}

type Change[T any] struct {
	Value   T
	Changed bool
}

type Error struct {
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

type ErrorList []*Error

func (p *ErrorList) Add(message string) {
	*p = append(*p, &Error{
		Message: message,
	})
}
