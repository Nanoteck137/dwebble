package library

type TrackFile struct {
	Lossless string `toml:"lossless"`
	Lossy    string `toml:"lossy"`
}

type TrackMetadata struct {
	Num       int       `toml:"num"`
	Name      string    `toml:"name"`
	Duration  int       `toml:"duration"`
	Artist    string    `toml:"artist"`
	Year      int       `toml:"year"`
	Tags      []string  `toml:"tags"`
	Genres    []string  `toml:"genres"`
	Featuring []string  `toml:"featuring"`
	File      TrackFile `toml:"file,inline"`
}

type AlbumMetadata struct {
	Album    string          `toml:"album"`
	Artist   string          `toml:"artist"`
	CoverArt string          `toml:"coverart"`
	Tracks   []TrackMetadata `toml:"tracks"`
}
