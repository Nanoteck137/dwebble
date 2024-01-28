package collection

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/kennygrant/sanitize"
	"github.com/nanoteck137/dwebble/utils"
)

type TrackMetadata struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Number   int    `json:"number"`
	FileName string `json:"fileName"`
	ArtistId string `json:"artistId,omitempty"`
}

type AlbumMetadata struct {
	Id       string          `json:"id"`
	Name     string          `json:"name"`
	CoverArt string          `json:"coverArt"`
	Dir      string          `json:"dir"`
	Tracks   []TrackMetadata `json:"tracks"`
}

type ArtistMetadata struct {
	Id      string          `json:"id"`
	Name    string          `json:"name"`
	Picture string          `json:"picture"`
	Albums  []AlbumMetadata `json:"albums"`
}

type Track struct {
	Id       string
	Number   int
	Name     string
	FilePath string
	AlbumId  string
	ArtistId string
}

type Album struct {
	Id           string
	Name         string
	CoverArtPath string
	Dir          string
	TrackIds     []string
	ArtistId     string
}

type Artist struct {
	Id          string
	Dir         string
	Name        string
	PicturePath string
	AlbumIds    []string
}

type Collection struct {
	BasePath string
	Artists  map[string]*Artist
	Albums   map[string]*Album
	Tracks   map[string]*Track
}

func ReadEntryFromDir(dir string) (*ArtistMetadata, error) {
	artistFile := path.Join(dir, "artist.json")

	data, err := os.ReadFile(artistFile)
	if err != nil {
		return nil, err
	}

	var artistMetadata ArtistMetadata
	err = json.Unmarshal(data, &artistMetadata)
	if err != nil {
		return nil, err
	}

	return &artistMetadata, nil
}

func ReadFromDir(dir string) (*Collection, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var artists []*ArtistMetadata

	artistMap := make(map[string]*Artist)
	albumMap := make(map[string]*Album)
	trackMap := make(map[string]*Track)

	for _, entry := range entries {
		p := path.Join(dir, entry.Name())
		fmt.Printf("Path: %v\n", p)

		artist, err := ReadEntryFromDir(p)
		if err != nil {
			log.Printf("Warning: %v\n", err)
			continue
		}

		artistMap[artist.Id] = &Artist{
			Id:          artist.Id,
			Dir:         entry.Name(),
			Name:        artist.Name,
			PicturePath: path.Join(entry.Name(), artist.Picture),
			AlbumIds:    []string{},
		}

		artists = append(artists, artist)
	}

	for _, artistMetadata := range artists {
		artist, ok := artistMap[artistMetadata.Id]
		if !ok {
			log.Fatalf("No artist with id: '%v'\n", artistMetadata.Id)
		}

		for _, albumMetadata := range artistMetadata.Albums {
			album := &Album{
				Id:           albumMetadata.Id,
				Name:         albumMetadata.Name,
				CoverArtPath: path.Join(artist.Dir, albumMetadata.Dir, "picture.png"),
				Dir:          albumMetadata.Dir,
				TrackIds:     []string{},
				ArtistId:     artistMetadata.Id,
			}

			_, exists := albumMap[albumMetadata.Id]
			if exists {
				log.Fatalf("Duplicated album id: '%v'\n", albumMetadata.Id)
			}

			albumMap[albumMetadata.Id] = album

			artist.AlbumIds = append(artist.AlbumIds, albumMetadata.Id)

			for _, trackMetadata := range albumMetadata.Tracks {
				filePath := path.Join(artist.Dir, album.Dir, trackMetadata.FileName)

				track := &Track{
					Id:       trackMetadata.Id,
					Number:   trackMetadata.Number,
					Name:     trackMetadata.Name,
					FilePath: filePath,
					AlbumId:  album.Id,
					ArtistId: artist.Id,
				}

				_, exists := trackMap[trackMetadata.Id]
				if exists {
					log.Fatalf("Duplicated track id: '%v'\n", albumMetadata.Id)
				}

				trackMap[trackMetadata.Id] = track
				album.TrackIds = append(album.TrackIds, track.Id)
			}
		}
	}

	return &Collection{
		BasePath: dir,
		Artists:  artistMap,
		Albums:   albumMap,
		Tracks:   trackMap,
	}, nil
}

func NewEmpty(dir string) *Collection {
	return &Collection{
		BasePath: dir,
		Artists:  map[string]*Artist{},
		Albums:   map[string]*Album{},
		Tracks:   map[string]*Track{},
	}
}

type File struct {
	Content io.Reader
	ContentType string
}

var ErrUnsupportedContentType = errors.New("Unsupported content type")

func (file *File) Ext() (string, error) {
	switch file.ContentType {
	case "image/png":
		return "png", nil
	case "image/jpeg":
		return "jpg", nil
	case "audio/flac":
		return "flac", nil
	}

	return "", ErrUnsupportedContentType
}

func (col *Collection) CreateArtist(name string, picture File) (*Artist, error) {
	newId := utils.CreateId()
	a, exists := col.Artists[newId]
	if exists {
		return nil, fmt.Errorf("Failed to create new artist: Generated id '%v' matched with '%v' id", newId, a.Name)
	}

	for _, artist := range col.Artists {
		if artist.Name == name {
			return nil, fmt.Errorf("Failed to create new artist: Artist with name '%v' already exists", name)
		}
	}

	dirName := sanitize.Name(name)
	p := path.Join(col.BasePath, dirName)

	_, err := os.Stat(p)
	if err == nil {
		return nil, fmt.Errorf("Failed to create new artist: Directory with name '%v' already exists", dirName)
	}

	err = os.MkdirAll(p, 0755)
	if err != nil {
		return nil, err
	}

	ext, err := picture.Ext()
	if err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("picture.%v", ext)
	dstFile, err := os.Create(path.Join(p, fileName))
	if err != nil {
		return nil, err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, picture.Content)
	if err != nil {
		return nil, err
	}

	artist := new(Artist)
	*artist = Artist{
		Id:          newId,
		Dir:         dirName,
		Name:        name,
		PicturePath: path.Join(dirName, fileName),
		AlbumIds:    []string{},
	}

	col.Artists[newId] = artist

	return artist, nil
}

func (col *Collection) CreateAlbum(name string, coverArt File, artist *Artist) (*Album, error) {
	newId := utils.CreateId()
	a, exists := col.Albums[newId]
	if exists {
		return nil, fmt.Errorf("Failed to create new album: Generated id '%v' matched with '%v' id", newId, a.Name)
	}

	for _, artist := range col.Albums {
		if artist.Name == name {
			return nil, fmt.Errorf("Failed to create new album: '%v' already exists", name)
		}
	}

	dirName := sanitize.Name(name)
	p := path.Join(col.BasePath, artist.Dir, dirName)

	_, err := os.Stat(p)
	if err == nil {
		return nil, fmt.Errorf("Failed to create new album: Directory with name '%v' already exists", dirName)
	}

	err = os.MkdirAll(p, 0755)
	if err != nil {
		return nil, err
	}

	ext, err := coverArt.Ext()
	if err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("coverart.%v", ext)
	dstFile, err := os.Create(path.Join(p, fileName))
	if err != nil {
		return nil, err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, coverArt.Content)
	if err != nil {
		return nil, err
	}

	album := new(Album)
	*album = Album{
		Id:           newId,
		Name:         name,
		CoverArtPath: path.Join(artist.Dir, dirName, fileName),
		Dir:          dirName,
		TrackIds:     []string{},
		ArtistId:     artist.Id,
	}

	artist.AlbumIds = append(artist.AlbumIds, album.Id)

	fmt.Printf("%v\n", album.CoverArtPath)

	col.Albums[newId] = album

	return album, nil
}

func (col* Collection) GetAllTrackFromAlbum(album *Album) []*Track {
	var result []*Track

	for _, trackId := range album.TrackIds {
		// TODO(patrik): Check if it exists
		track := col.Tracks[trackId]
		result = append(result, track)
	}

	return result
}

func (col *Collection) CreateTrack(name string, number int, trackFile File, album *Album, artist *Artist) (*Track, error) {
	newId := utils.CreateId()
	t, exists := col.Tracks[newId]
	if exists {
		return nil, fmt.Errorf("Generated id '%v' matched with '%v' id", newId, t.Name)
	}

	for _, t := range col.GetAllTrackFromAlbum(album) {
		if t.Number == number {
			return nil, fmt.Errorf("Track with number '%v' already exists '%v' inside '%v'", number, t.Name, album.Name)
		}
	}

	p := path.Join(col.BasePath, artist.Dir, album.Dir)

	ext, err := trackFile.Ext()
	if err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("%v.%v", number, ext)
	filePath := path.Join(p, fileName)
	dstFile, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, trackFile.Content)
	if err != nil {
		return nil, err
	}

	track := new(Track)
	*track = Track{
		Id:       newId,
		Number:   number,
		Name:     name,
		FilePath: filePath,
		AlbumId:  album.Id,
		ArtistId: artist.Id,
	}

	album.TrackIds = append(album.TrackIds, track.Id)

	col.Tracks[newId] = track

	return track, nil
}

func (col *Collection) FlushToDisk() error {
	for _, artist := range col.Artists {
		var albums []AlbumMetadata
		for _, albumId := range artist.AlbumIds {
			// TODO(patrik): Check if it exists
			album := col.Albums[albumId]

			albums = append(albums, AlbumMetadata{
				Id:       albumId,
				Name:     album.Name,
				CoverArt: path.Base(album.CoverArtPath),
				Dir:      album.Dir,
				Tracks:   []TrackMetadata{},
			})
		}

		metadata := ArtistMetadata{
			Id:      artist.Id,
			Name:    artist.Name,
			Picture: path.Base(artist.PicturePath),
			Albums:  albums,
		}

		data, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(data))
		p := path.Join(col.BasePath, artist.Dir, "artist.json")

		err = os.WriteFile(p, data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
