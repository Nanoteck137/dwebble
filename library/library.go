package library

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
	"github.com/nanoteck137/parasect"
)

type Track struct {
	Path     string
	Name     string
	Number   int
	Artist   string
	Duration int
	Tags     []string
}

type Album struct {
	Path       string
	Name       string
	ArtistName string
	CoverArt   string
	Tracks     []Track
}

type Artist struct {
	Path    string
	Name    string
	Picture string
	Albums  []Album
}

type Library struct {
	Artists map[string]*Artist
}

var fileReg = regexp.MustCompile(`(\d*).*`)

func getAllTrackFromDir(dir string, albumArtist string) ([]Track, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var tracks []Track

	for _, entry := range entries {
		p := path.Join(dir, entry.Name())

		if utils.IsMusicFile(p) {
			info, err := parasect.GetTrackInfo(p)
			if err != nil {
				continue
			}

			captures := fileReg.FindStringSubmatch(entry.Name())
			if captures == nil {
				continue
			}

			if captures[1] == "" {
				continue
			}

			trackNum, err := strconv.Atoi(captures[1])
			if err != nil {
				log.Fatal(err)
			}

			name := ""
			tagTitle := info.Tags["title"]
			if tagTitle != "" {
				name = tagTitle
			} else {
				ext := path.Ext(p)
				name = strings.TrimSuffix(entry.Name(), ext)
			}

			artist := albumArtist
			tagArtist := info.Tags["artist"]
			if tagArtist != "" {
				artist = tagArtist
			}

			tags := strings.Split(info.Tags["tags"], ",")

			var realTags []string
			for _, tag := range tags {
				if tag == "" {
					continue
				}

				realTags = append(realTags, strings.TrimSpace(tag))
			}

			duration := info.Duration

			tracks = append(tracks, Track{
				Path:     p,
				Name:     name,
				Number:   trackNum,
				Artist:   artist,
				Duration: duration,
				Tags:     realTags,
			})
		}
	}

	return tracks, nil
}

func ReadFromDir(dir string) (*Library, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	artists := make(map[string]*Artist)

	for _, entry := range entries {
		p := path.Join(dir, entry.Name())
		if entry.IsDir() {
			entries, err := os.ReadDir(p)
			if err != nil {
				log.Fatal(err)
			}

			if !utils.HasMusic(entries) {
				override := path.Join(p, "override.txt")

				name := entry.Name()

				data, err := os.ReadFile(override)
				if err == nil {
					name = strings.TrimSpace(string(data))
				}

				entries, err := os.ReadDir(p)
				if err != nil {
					log.Fatal(err)
				}

				var picture string
				for _, entry := range entries {
					if entry.IsDir() {
						continue
					}

					name := entry.Name()
					ext := path.Ext(name)
					if ext == "" {
						continue
					}

					name = strings.TrimSuffix(name, ext)
					if name != "picture" {
						continue
					}

					if utils.IsValidImageExt(ext[1:]) {
						picture = entry.Name()
						break
					}
				}

				artists[name] = &Artist{
					Path:    p,
					Name:    name,
					Picture: picture,
				}
			}
		}
	}

	for _, entry := range entries {
		p := path.Join(dir, entry.Name())
		if entry.IsDir() {
			entries, err := os.ReadDir(p)
			if err != nil {
				log.Fatal(err)
			}

			artistName := "Various Artists"

			artistName = entry.Name()

			data, err := os.ReadFile(path.Join(p, "override.txt"))
			if err == nil {
				artistName = strings.TrimSpace(string(data))
			}

			for _, entry := range entries {
				p := path.Join(p, entry.Name())

				if !entry.IsDir() {
					continue
				}

				entries, err := os.ReadDir(p)
				if err != nil {
					log.Fatal(err)
				}

				if !utils.HasMusic(entries) {
					continue
				}

				albumName := entry.Name()
				data, err := os.ReadFile(path.Join(p, "override.txt"))
				if err == nil {
					albumName = strings.TrimSpace(string(data))
				}

				var coverArt string
				for _, entry := range entries {
					if entry.IsDir() {
						continue
					}

					name := entry.Name()
					ext := path.Ext(name)
					if ext == "" {
						continue
					}

					name = strings.TrimSuffix(name, ext)
					if name != "cover" {
						continue
					}

					if utils.IsValidImageExt(ext[1:]) {
						coverArt = entry.Name()
						break
					}
				}

				tracks, err := getAllTrackFromDir(p, artistName)
				if err != nil {
					log.Fatal(err)
				}

				for _, track := range tracks {
					if _, exists := artists[track.Artist]; !exists {
						artists[track.Artist] = &Artist{
							Name: track.Artist,
						}
					}
				}

				artist := artists[artistName]

				artist.Albums = append(artist.Albums, Album{
					Path:       p,
					Name:       albumName,
					ArtistName: artistName,
					CoverArt:   coverArt,
					Tracks:     tracks,
				})
			}
		}
	}

	return &Library{
		Artists: artists,
	}, nil
}

func GetOrCreateArtist(ctx context.Context, db *database.Database, artist *Artist) (database.Artist, error) {
	var dbArtist database.Artist
	var err error

	if artist.Path != "" {
		dbArtist, err = db.GetArtistByPath(ctx, artist.Path)
	} else {
		dbArtist, err = db.GetArtistByName(ctx, artist.Name)
	}

	if err != nil {
		if err == types.ErrNoArtist {
			artist, err := db.CreateArtist(ctx, database.CreateArtistParams{
				Name: artist.Name,
				Path: artist.Path,
			})
			if err != nil {
				return database.Artist{}, err
			}

			return artist, nil
		} else {
			return database.Artist{}, err
		}
	}

	return dbArtist, nil
}

func GetOrCreateAlbum(ctx context.Context, db *database.Database, album *Album, artistId string) (database.Album, error) {
	dbAlbum, err := db.GetAlbumByPath(ctx, album.Path)
	if err != nil {
		if err == types.ErrNoAlbum {
			album, err := db.CreateAlbum(ctx, database.CreateAlbumParams{
				Name:     album.Name,
				ArtistId: artistId,
				Path:     album.Path,
			})

			if err != nil {
				return database.Album{}, err
			}

			return album, nil
		} else {
			return database.Album{}, err
		}
	}

	return dbAlbum, nil
}

func GetOrCreateTrack(ctx context.Context, db *database.Database, track *Track, albumId string, artistId string) (database.Track, error) {
	dbTrack, err := db.GetTrackByPath(ctx, track.Path)
	if err != nil {
		if err == types.ErrNoTrack {
			track, err := db.CreateTrack(ctx, database.CreateTrackParams{
				TrackNumber: track.Number,
				Name:        track.Name,
				Path:        track.Path,
				Duration:    track.Duration,
				AlbumId:     albumId,
				ArtistId:    artistId,
			})

			if err != nil {
				return database.Track{}, err
			}

			return track, nil
		} else {
			return database.Track{}, err
		}
	}

	return dbTrack, nil
}

func GetOrCreateTag(ctx context.Context, db *database.Database, name string) (database.Tag, error) {
	tag, err := db.GetTagByName(ctx, name)
	if err != nil {
		if err == pgx.ErrNoRows {
			tag, err = db.CreateTag(ctx, name)
			if err != nil {
				return database.Tag{}, err
			}
		} else {
			return database.Tag{}, err
		}
	}

	return tag, nil
}

func (lib *Library) syncArtists(ctx context.Context, db *database.Database, imageDir string) (map[string]database.Artist, error) {
	artists := make(map[string]database.Artist)

	for _, artist := range lib.Artists {
		fmt.Println("Syncing:", artist.Name)

		dbArtist, err := GetOrCreateArtist(ctx, db, artist)
		if err != nil {
			return nil, err
		}

		var artistChanges database.ArtistChanges
		artistChanges.Available = true
		artistChanges.Name.Value = artist.Name
		artistChanges.Name.Changed = dbArtist.Name != artist.Name

		if artist.Picture != "" {
			p := path.Join(artist.Path, artist.Picture)
			ext := path.Ext(artist.Picture)
			name := "artist_" + dbArtist.Id + ext
			dst := path.Join(imageDir, name)
			err := utils.SymlinkReplace(p, dst)
			if err != nil {
				return nil, err
			}

			artistChanges.Picture.Value = sql.NullString{
				String: name,
				Valid:  true,
			}

			artistChanges.Picture.Changed = true
		}

		err = db.UpdateArtist(ctx, dbArtist.Id, artistChanges)
		if err != nil {
			return nil, err
		}

		artists[artist.Name] = dbArtist
	}

	return artists, nil
}

func (lib *Library) Sync(workDir types.WorkDir, dir string, db *database.Database) error {
	trackDir := workDir.OriginalTracksDir()
	err := os.MkdirAll(trackDir, 0755)
	if err != nil {
		return err
	}

	mobileTrackDir := workDir.MobileTracksDir()
	err = os.MkdirAll(mobileTrackDir, 0755)
	if err != nil {
		return err
	}

	transcodeDir := workDir.TranscodeDir()
	err = os.MkdirAll(transcodeDir, 0755)
	if err != nil {
		return err
	}

	imageDir := workDir.ImagesDir()
	err = os.MkdirAll(imageDir, 0755)
	if err != nil {
		return err
	}

	ctx := context.Background()

	err = db.MarkAllTracksUnavailable(ctx)
	if err != nil {
		return err
	}

	err = db.MarkAllAlbumsUnavailable(ctx)
	if err != nil {
		return err
	}

	err = db.MarkAllArtistsUnavailable(ctx)
	if err != nil {
		return err
	}

	artists, err := lib.syncArtists(ctx, db, imageDir)

	for _, artist := range lib.Artists {
		dbArtist := artists[artist.Name]

		for _, album := range artist.Albums {
			dbAlbum, err := GetOrCreateAlbum(ctx, db, &album, dbArtist.Id)
			if err != nil {
				return err
			}

			var albumChanges database.AlbumChanges
			albumChanges.Available = true
			albumChanges.Name.Value = album.Name
			albumChanges.Name.Changed = dbAlbum.Name != album.Name

			albumCoverRealPath := path.Join(album.Path, album.CoverArt)
			if album.CoverArt != "" {
				ext := path.Ext(album.CoverArt)
				name := "album_" + dbAlbum.Id + ext
				dst := path.Join(imageDir, name)
				err := utils.SymlinkReplace(albumCoverRealPath, dst)
				if err != nil {
					return err
				}

				albumChanges.CoverArt.Value = sql.NullString{
					String: name,
					Valid:  true,
				}

				albumChanges.CoverArt.Changed = true

				dbAlbum.CoverArt = sql.NullString{
					String: name,
					Valid:  true,
				}
			}

			err = db.UpdateAlbum(ctx, dbAlbum.Id, albumChanges)
			if err != nil {
				return err
			}

			for _, track := range album.Tracks {
				artist := artists[track.Artist]
				coverArt := dbAlbum.CoverArt

				dbTrack, err := GetOrCreateTrack(ctx, db, &track, dbAlbum.Id, artist.Id)
				if err != nil {
					log.Fatal(err)
				}

				tags, err := db.GetTrackTags(ctx, dbTrack.Id)
				if err != nil {
					log.Fatal(err)
				}

				for _, tag := range track.Tags {
					tag, err := GetOrCreateTag(ctx, db, tag)
					if err != nil {
						log.Fatal(err)
					}

					err = db.AddTagToTrack(ctx, tag.Id, dbTrack.Id)
					if err != nil {
						var pgErr *pgconn.PgError
						if errors.As(err, &pgErr) {
							if pgErr.Code == pgerrcode.UniqueViolation {
								continue
							}
						}

						log.Fatal(err)
					}
				}

				hasTag := func(tag string) bool {
					for _, t := range track.Tags {
						if tag == t {
							return true
						}
					}

					return false
				}

				for _, tag := range tags {
					if !hasTag(tag.Name) {
						err := db.RemoveTagFromTrack(ctx, tag.Id, dbTrack.Id)
						if err != nil {
							log.Fatal(err)
						}
					}
				}

				var trackChanges database.TrackChanges
				trackChanges.Available = true
				trackChanges.Number.Value = track.Number
				trackChanges.Number.Changed = dbTrack.Number != track.Number
				trackChanges.Name.Value = track.Name
				trackChanges.Name.Changed = dbTrack.Name != track.Name
				trackChanges.CoverArt.Value = coverArt
				trackChanges.CoverArt.Changed = dbTrack.CoverArt != coverArt
				trackChanges.Duration.Value = track.Duration
				trackChanges.Duration.Changed = dbTrack.Duration != track.Duration

				trackPath := track.Path
				trackExt := path.Ext(trackPath)
				name := dbTrack.Id + trackExt
				dst := path.Join(trackDir, name)

				err = utils.SymlinkReplace(trackPath, dst)
				if err != nil {
					log.Fatal(err)
				}

				var mobileSrc string

				if utils.IsLossyFormatExt(trackExt) {
					mobileSrc = trackPath
				} else {
					transcodeName := fmt.Sprintf("%v.opus", dbTrack.Id)
					dstTranscode := path.Join(transcodeDir, transcodeName)

					_, err = os.Stat(dstTranscode)
					if err != nil {
						if os.IsNotExist(err) {
							err = parasect.RunFFmpeg(true, "-y", "-i", trackPath, "-vbr", "on", "-b:a", "128k", dstTranscode)
							if err != nil {
								log.Fatal(err)
							}
						} else {
							log.Fatal(err)
						}
					}

					src, err := filepath.Abs(dstTranscode)
					if err != nil {
						log.Fatal(err)
					}

					mobileSrc = src
				}

				mobileSrcExt := filepath.Ext(mobileSrc)
				mobileSrcName := dbTrack.Id + mobileSrcExt
				dst = path.Join(mobileTrackDir, mobileSrcName)
				err = utils.SymlinkReplace(mobileSrc, dst)
				if err != nil {
					log.Fatal(err)
				}

				trackChanges.BestQualityFile.Value = name
				trackChanges.BestQualityFile.Changed = dbTrack.BestQualityFile != name
				trackChanges.MobileQualityFile.Value = mobileSrcName
				trackChanges.MobileQualityFile.Changed = dbTrack.MobileQualityFile != mobileSrcName
				trackChanges.ArtistId.Value = artist.Id
				trackChanges.ArtistId.Changed = dbTrack.ArtistId != artist.Id

				err = db.UpdateTrack(ctx, dbTrack.Id, trackChanges)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	return nil
}
