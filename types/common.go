package types

import "path"

type Map map[string]any

type WorkDir string

func (d WorkDir) String() string {
	return string(d)
}

func (d WorkDir) OriginalTracksDir() string {
	return path.Join(d.String(), "original-tracks")
}

func (d WorkDir) MobileTracksDir() string {
	return path.Join(d.String(), "mobile-tracks")
}

func (d WorkDir) ImagesDir() string {
	return path.Join(d.String(), "images")
}
