package types

import "path"

type MediaType string

const (
	MediaTypeFlac      MediaType = "flac"
	MediaTypeOggOpus   MediaType = "ogg-opus"
	MediaTypeOggVorbis MediaType = "ogg-vorbis"
	MediaTypeMp3       MediaType = "mp3"
	MediaTypeAcc       MediaType = "acc"
)

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

func (d WorkDir) Track(id string) TrackDir {
	return TrackDir(path.Join(d.Tracks(), id))
}

type TrackDir string

func (d TrackDir) String() string {
	return string(d)
}

func (d TrackDir) Media() string {
	return path.Join(d.String(), "media")
}

func (d TrackDir) MediaItem(id string) string {
	return path.Join(d.Media(), id)
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
