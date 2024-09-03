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

func (d WorkDir) Albums() string {
	return path.Join(d.String(), "albums")
}

func (d WorkDir) Album(slug string) AlbumDir {
	return AlbumDir(path.Join(d.String(), "albums", slug))
}

func (d WorkDir) Artists() string {
	return path.Join(d.String(), "artists")
}

type AlbumDir string

func (d AlbumDir) String() string {
	return string(d)
}

func (d AlbumDir) Images() string {
	return path.Join(d.String(), "images")
}

func (d AlbumDir) OriginalFiles() string {
	return path.Join(d.String(), "original")
}

func (d AlbumDir) MobileFiles() string {
	return path.Join(d.String(), "mobile")
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
