package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/apis"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/spf13/cobra"
)

type OldWorkDir string

func (d OldWorkDir) String() string {
	return string(d)
}

func (d OldWorkDir) DatabaseFile() string {
	return path.Join(d.String(), "data.db")
}

func (d OldWorkDir) ExportFile() string {
	return path.Join(d.String(), "export.json")
}

func (d OldWorkDir) Albums() string {
	return path.Join(d.String(), "albums")
}

func (d OldWorkDir) Album(slug string) OldAlbumDir {
	return OldAlbumDir(path.Join(d.String(), "albums", slug))
}

func (d OldWorkDir) Artists() string {
	return path.Join(d.String(), "artists")
}

type OldAlbumDir string

func (d OldAlbumDir) String() string {
	return string(d)
}

func (d OldAlbumDir) Images() string {
	return path.Join(d.String(), "images")
}

func (d OldAlbumDir) OriginalFiles() string {
	return path.Join(d.String(), "original")
}

func (d OldAlbumDir) MobileFiles() string {
	return path.Join(d.String(), "mobile")
}

var AppName = dwebble.AppName + "-migrate"

var rootCmd = &cobra.Command{
	Use:     AppName,
	Version: dwebble.Version,
	Run: func(cmd *cobra.Command, args []string) {
		workDir := types.WorkDir("./work")

		// TODO(patrik): Move to function (used in base_app.Bootstrap, dwebble-import)
		dirs := []string{
			workDir.Albums(),
			workDir.Artists(),
			workDir.Tracks(),
		}

		for _, dir := range dirs {
			err := os.Mkdir(dir, 0755)
			if err != nil && !os.IsExist(err) {
				log.Fatal("Failed to bootstrap directories", "err", err)
			}
		}

		db, err := database.Open(workDir)
		if err != nil {
			log.Fatal("Failed to open database", "err", err)
		}

		err = database.RunMigrateUp(db)
		if err != nil {
			log.Fatal("Failed migrate database", "err", err)
		}

		db, tx, err := db.Begin()
		if err != nil {
			log.Fatal("Failed to start database transaction", "err", err)
		}
		defer tx.Rollback()

		oldWorkDir := OldWorkDir(".")

		d, err := os.ReadFile(oldWorkDir.ExportFile())
		if err != nil {
			log.Fatal("Failed to read export data", "err", err)
		}

		var export apis.Export
		err = json.Unmarshal(d, &export)
		if err != nil {
			log.Fatal("Failed to unmarshal export data", "err", err)
		}

		export, err = sanitizeExport(export)
		if err != nil {
			log.Fatal("Failed to sanitize export", "err", err)
		}

		ctx := context.TODO()
		_ = ctx

		artistMapping := make(map[string]string)
		_ = artistMapping

		for _, artist := range export.Artists {
			dbArtist, err := db.CreateArtist(ctx, database.CreateArtistParams{
				Id:   utils.CreateArtistId(),
				Name: artist.Name,
				Picture: sql.NullString{
					String: artist.Picture,
					Valid:  artist.Picture != "",
				},
			})

			if err != nil {
				var e sqlite3.Error
				if errors.As(err, &e) {
					if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
						continue
					}

					if e.ExtendedCode == sqlite3.ErrConstraintUnique {
						log.Warn("Duplicate artist name", "name", artist.Name)
						continue
					}
				}

				log.Fatal("Failed to create artist", "err", err)
			}

			artistMapping[artist.Id] = dbArtist.Id
		}

		albumMapping := make(map[string]string)

		for _, album := range export.Albums {
			artistId, exists := artistMapping[album.ArtistId]
			if !exists {
				log.Fatal("Missing album artist id", "artistId", album.ArtistId)
			}

			dbAlbum, err := db.CreateAlbum(ctx, database.CreateAlbumParams{
				Id:       utils.CreateAlbumId(),
				Name:     album.Name,
				ArtistId: artistId,
				CoverArt: sql.NullString{
					String: album.CoverArt,
					Valid:  album.CoverArt != "",
				},
				Year: sql.NullInt64{
					Int64: album.Year,
					Valid: album.Year != 0,
				},
				Available: true,
			})
			if err != nil {
				var e sqlite3.Error
				if errors.As(err, &e) {
					if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
						continue
					}
				}

				log.Fatal("Failed", "err", err)
			}

			albumMapping[album.Id] = dbAlbum.Id
		}

		for _, track := range export.Tracks {
			artistId, exists := artistMapping[track.ArtistId]
			if !exists {
				log.Fatal("Missing track artist id", "artistId", track.ArtistId)
			}

			albumId, exists := albumMapping[track.AlbumId]
			if !exists {
				log.Fatal("Missing track album id", "albumId", track.AlbumId)
			}

			dbTrackId, err := db.CreateTrack(ctx, database.CreateTrackParams{
				Id:       utils.CreateTrackId(),
				Name:     track.Name,
				AlbumId:  albumId,
				ArtistId: artistId,
				Number: sql.NullInt64{
					Int64: track.Number,
					Valid: track.Number != 0,
				},
				Duration: sql.NullInt64{
					Int64: track.Duration,
					Valid: track.Duration != 0,
				},
				Year: sql.NullInt64{
					Int64: track.Year,
					Valid: track.Year != 0,
				},
				ExportName:       track.ExportName,
				OriginalFilename: track.OriginalFilename,
				MobileFilename:   track.MobileFilename,
				Created:          track.Created,
				Available:        true,
			})
			if err != nil {
				var e sqlite3.Error
				if errors.As(err, &e) {
					if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
						continue
					}
				}

				log.Fatal("Failed", "err", err)
			}

			for _, tag := range track.Tags {
				slug := utils.Slug(tag)

				err := db.CreateTag(ctx, slug, tag)
				if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
					log.Fatal("Failed", "err", err)
				}

				err = db.AddTagToTrack(ctx, slug, dbTrackId)
				if err != nil {
					log.Fatal("Failed", "err", err)
				}
			}
		}

		err = tx.Commit()
		if err != nil {
			log.Fatal("Failed to commit transaction", "err", err)
		}
	},
}

func sanitizeString(s string) string {
	return strings.TrimSpace(s)
}

func sanitizeExport(export apis.Export) (apis.Export, error) {
	res := apis.Export{}
	res.Albums = export.Albums
	res.Tracks = export.Tracks

	artistWanted := make(map[string]struct{})

	artists := make(map[string]apis.ExportArtist)

	for _, album := range export.Albums {
		artistWanted[album.ArtistId] = struct{}{}
	}

	for _, album := range export.Tracks {
		artistWanted[album.ArtistId] = struct{}{}
	}

	for _, artist := range export.Artists {
		_, exists := artistWanted[artist.Id]
		if !exists {
			continue
		}

		sanitizedName := sanitizeString(artist.Name)
		k := strings.ToLower(sanitizedName)

		if k == "" {
			k = "unknown"
		}

		a, exists := artists[k]
		if exists {
			for i, t := range res.Tracks {
				if t.ArtistId == artist.Id {
					res.Tracks[i].ArtistId = a.Id
				}
			}

			for i, album := range res.Albums {
				if album.ArtistId == artist.Id {
					res.Albums[i].ArtistId = a.Id
				}
			}

		} else {
			artists[k] = artist
		}
	}

	res.Artists = make([]apis.ExportArtist, 0, len(artists))
	for _, artist := range artists {
		res.Artists = append(res.Artists, artist)
	}

	return res, nil
}

func init() {
	rootCmd.SetVersionTemplate(dwebble.VersionTemplate(AppName))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to run root command", "err", err)
	}
}

func main() {
	Execute()
}
