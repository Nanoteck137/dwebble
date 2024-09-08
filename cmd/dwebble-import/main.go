package main

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/pelletier/go-toml/v2"
)

func importAlbum(ctx context.Context, db *database.Database, workDir types.WorkDir, albumPath string) error {
	d, err := os.ReadFile(path.Join(albumPath, "album.toml"))
	if err != nil {
		return err
	}

	var metadata library.AlbumMetadata
	err = toml.Unmarshal(d, &metadata)
	if err != nil {
		return err
	}

	pretty.Println(metadata)

	getOrCreateArtist := func(name string) (database.Artist, error) {
		artist, err := db.GetArtistByName(ctx, name)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				artist, err := db.CreateArtist(ctx, database.CreateArtistParams{
					Name:    name,
					Picture: sql.NullString{},
				})
				if err != nil {
					return database.Artist{}, err
				}

				return artist, nil
			} else {
				return database.Artist{}, err
			}
		}

		return artist, nil
	}

	artist, err := getOrCreateArtist(metadata.Artist)
	pretty.Println(artist)

	album, err := db.CreateAlbum(ctx, database.CreateAlbumParams{
		Name:     metadata.Album,
		ArtistId: artist.Id,

		CoverArt: sql.NullString{},
		Year:     sql.NullInt64{},

		Available: true,
	})
	if err != nil {
		return err
	}

	albumDir := workDir.Album(album.Id)

	dirs := []string{
		albumDir.String(),
		albumDir.Images(),
		albumDir.OriginalFiles(),
		albumDir.MobileFiles(),
	}

	for _, dir := range dirs {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}

	pretty.Println(album)

	if metadata.CoverArt != "" {
		p := path.Join(albumPath, metadata.CoverArt)

		ext := path.Ext(metadata.CoverArt)
		filename := "cover" + ext

		// TODO(patrik): Close file
		file, err := os.OpenFile(path.Join(albumDir.Images(), filename), os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		ff, err := os.Open(p)
		if err != nil {
			return err
		}

		_, err = io.Copy(file, ff)
		if err != nil {
			return err
		}

		err = db.UpdateAlbum(ctx, album.Id, database.AlbumChanges{
			CoverArt: types.Change[sql.NullString]{
				Value: sql.NullString{
					String: filename,
					Valid:  true,
				},
				Changed: true,
			},
		})

		if err != nil {
			return err
		}
	}

	for _, track := range metadata.Tracks {
		now := time.Now()
		_ = now

		artist, err := getOrCreateArtist(track.Artist)
		if err != nil {
			return err
		}

		p := track.File.Lossless
		if p == "" {
			p = track.File.Lossy
			if p == "" {
				// TODO(patrik): Better error
				return errors.New("Track contains no file")
			}
		}

		p = path.Join(albumPath, p)

		trackFilename := path.Base(p)

		// TODO(patrik): Maybe save the original filename to use when exporting
		ext := path.Ext(trackFilename)
		originalName := strings.TrimSuffix(trackFilename, ext)
		filename := originalName + "-" + strconv.FormatInt(now.Unix(), 10)

		file, err := os.CreateTemp("", "track.*"+ext)
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())

		ff, err := os.Open(p)
		if err != nil {
			return err
		}

		_, err = io.Copy(file, ff)
		if err != nil {
			return err
		}

		file.Close()

		mobileFile, err := utils.ProcessMobileVersion(file.Name(), albumDir.MobileFiles(), filename)
		if err != nil {
			return err
		}

		originalFile, trackInfo, err := utils.ProcessOriginalVersion(file.Name(), albumDir.OriginalFiles(), filename)
		if err != nil {
			return err
		}

		trackId, err := db.CreateTrack(ctx, database.CreateTrackParams{
			Name:     track.Name,
			AlbumId:  album.Id,
			ArtistId: artist.Id,
			Number: sql.NullInt64{
				Int64: int64(track.Num),
				Valid: track.Num != 0,
			},
			Duration: sql.NullInt64{
				Int64: int64(trackInfo.Duration),
				Valid: true,
			},
			Year: sql.NullInt64{
				Int64: int64(track.Year),
				Valid: track.Year != 0,
			},
			ExportName:       originalName,
			OriginalFilename: originalFile,
			MobileFilename:   mobileFile,
			Available:        true,
		})
		if err != nil {
			return err
		}

		// TODO(patrik): Combine genres with tags
		for _, tag := range track.Tags {
			slug := utils.Slug(tag)

			err := db.CreateTag(ctx, slug, tag)
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}

			err = db.AddTagToTrack(ctx, slug, trackId)
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}
		}

		for _, genre := range track.Genres {
			slug := utils.Slug(genre)

			err := db.CreateTag(ctx, slug, genre)
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}

			err = db.AddTagToTrack(ctx, slug, trackId)
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}
		}
	}

	return nil
}

func main() {
	out := "./work"
	err := os.MkdirAll(out, 0755)
	if err != nil {
		log.Fatal(err)
	}

	workDir := types.WorkDir(out)

	dirs := []string{
		workDir.Albums(),
		workDir.Artists(),
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	db, err := database.Open(workDir)
	if err != nil {
		log.Fatal(err)
	}

	err = database.RunMigrateUp(db)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.TODO()

	var albums []string
	err = filepath.WalkDir("/Volumes/media/music/Anime/Abyss", func(p string, d fs.DirEntry, err error) error {
		if d.Name() == "album.toml" {
			albums = append(albums, path.Dir(p))
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, album := range albums {
		err = importAlbum(ctx, db, workDir, album)
		if err != nil {
			log.Fatal(err)
		}
	}

}
