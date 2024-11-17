package types

import "path"

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

func (d WorkDir) Artists() string {
	return path.Join(d.String(), "artists")
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
