package main

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

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

var AppName = dwebble.AppName + "-import"

var rootCmd = &cobra.Command{
	Use:     AppName + " <OUT>",
	Version: dwebble.Version,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out := args[0]
		dir, _ := cmd.Flags().GetString("dir")

		workDir := types.WorkDir(out)

		dirs := []string{
			workDir.Albums(),
			workDir.Artists(),
		}

		for _, dir := range dirs {
			err := os.Mkdir(dir, 0755)
			if err != nil && !os.IsExist(err) {
				log.Fatal("Failed", "err", err)
			}
		}

		db, err := database.Open(workDir)
		if err != nil {
			log.Fatal("Failed", "err", err)
		}

		err = database.RunMigrateUp(db)
		if err != nil {
			log.Fatal("Failed", "err", err)
		}

		ctx := context.TODO()

		var albums []string
		err = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if d == nil {
				return nil
			}

			if d.Name() == "album.toml" {
				albums = append(albums, path.Dir(p))
			}

			return nil
		})
		if err != nil {
			log.Fatal("Failed", "err", err)
		}

		for _, album := range albums {
			err = importAlbum(ctx, db, workDir, album)
			if err != nil {
				log.Fatal("Failed", "err", err)
			}
		}
	},
}

func init() {
	rootCmd.SetVersionTemplate(dwebble.VersionTemplate(AppName))
	rootCmd.Flags().String("dir", ".", "Search directory")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to run root command", "err", err)
	}
}

func importAlbum(ctx context.Context, db *database.Database, workDir types.WorkDir, albumPath string) error {
	db, tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	d, err := os.ReadFile(path.Join(albumPath, "album.toml"))
	if err != nil {
		return err
	}

	var metadata AlbumMetadata
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
		filename := "cover-original" + ext

		dst := path.Join(albumDir.Images(), filename)

		// TODO(patrik): Close file
		file, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		ff, err := os.Open(p)
		if err != nil {
			return err
		}
		defer ff.Close()

		_, err = io.Copy(file, ff)
		if err != nil {
			return err
		}

		i := path.Join(albumDir.Images(), "cover-128.png")
		err = utils.CreateResizedImage(dst, i, 128)
		if err != nil {
			return err
		}

		i = path.Join(albumDir.Images(), "cover-256.png")
		err = utils.CreateResizedImage(dst, i, 256)
		if err != nil {
			return err
		}

		i = path.Join(albumDir.Images(), "cover-512.png")
		err = utils.CreateResizedImage(dst, i, 512)
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

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	Execute()
}
