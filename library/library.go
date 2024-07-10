package library

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
	"github.com/pelletier/go-toml/v2"
)

// type Track struct {
// 	Path     string
// 	Name     string
// 	Number   int
// 	Artist   string
// 	Duration int
// 	Tags     []string
// }
//
// type Album struct {
// 	Path       string
// 	Name       string
// 	ArtistName string
// 	CoverArt   string
// 	Tracks     []Track
// }
//
// type Artist struct {
// 	Path    string
// 	Name    string
// 	Picture string
// 	Albums  []Album
// }
//
// type Library struct {
// 	Artists map[string]*Artist
// }
//
// var fileReg = regexp.MustCompile(`(\d*).*`)
//
// func getAllTrackFromDir(dir string, albumArtist string) ([]Track, error) {
// 	entries, err := os.ReadDir(dir)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var tracks []Track
//
// 	for _, entry := range entries {
// 		p := path.Join(dir, entry.Name())
//
// 		if utils.IsMusicFile(p) {
// 			info, err := parasect.GetTrackInfo(p)
// 			if err != nil {
// 				continue
// 			}
//
// 			captures := fileReg.FindStringSubmatch(entry.Name())
// 			if captures == nil {
// 				continue
// 			}
//
// 			if captures[1] == "" {
// 				continue
// 			}
//
// 			trackNum, err := strconv.Atoi(captures[1])
// 			if err != nil {
// 				log.Fatal(err)
// 			}
//
// 			name := ""
// 			tagTitle := info.Tags["title"]
// 			if tagTitle != "" {
// 				name = tagTitle
// 			} else {
// 				ext := path.Ext(p)
// 				name = strings.TrimSuffix(entry.Name(), ext)
// 			}
//
// 			artist := albumArtist
// 			tagArtist := info.Tags["artist"]
// 			if tagArtist != "" {
// 				artist = tagArtist
// 			}
//
// 			tags := strings.Split(info.Tags["tags"], ",")
//
// 			var realTags []string
// 			for _, tag := range tags {
// 				if tag == "" {
// 					continue
// 				}
//
// 				realTags = append(realTags, strings.TrimSpace(tag))
// 			}
//
// 			duration := info.Duration
//
// 			tracks = append(tracks, Track{
// 				Path:     p,
// 				Name:     name,
// 				Number:   trackNum,
// 				Artist:   artist,
// 				Duration: duration,
// 				Tags:     realTags,
// 			})
// 		}
// 	}
//
// 	return tracks, nil
// }
//
// func ReadFromDir(dir string) (*Library, error) {
// 	entries, err := os.ReadDir(dir)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	artists := make(map[string]*Artist)
//
// 	for _, entry := range entries {
// 		p := path.Join(dir, entry.Name())
// 		if entry.IsDir() {
// 			entries, err := os.ReadDir(p)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
//
// 			if !utils.HasMusic(entries) {
// 				override := path.Join(p, "override.txt")
//
// 				name := entry.Name()
//
// 				data, err := os.ReadFile(override)
// 				if err == nil {
// 					name = strings.TrimSpace(string(data))
// 				}
//
// 				entries, err := os.ReadDir(p)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
//
// 				var picture string
// 				for _, entry := range entries {
// 					if entry.IsDir() {
// 						continue
// 					}
//
// 					name := entry.Name()
// 					ext := path.Ext(name)
// 					if ext == "" {
// 						continue
// 					}
//
// 					name = strings.TrimSuffix(name, ext)
// 					if name != "picture" {
// 						continue
// 					}
//
// 					if utils.IsValidImageExt(ext[1:]) {
// 						picture = entry.Name()
// 						break
// 					}
// 				}
//
// 				artists[name] = &Artist{
// 					Path:    p,
// 					Name:    name,
// 					Picture: picture,
// 				}
// 			}
// 		}
// 	}
//
// 	for _, entry := range entries {
// 		p := path.Join(dir, entry.Name())
// 		if entry.IsDir() {
// 			entries, err := os.ReadDir(p)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
//
// 			artistName := "Various Artists"
//
// 			artistName = entry.Name()
//
// 			data, err := os.ReadFile(path.Join(p, "override.txt"))
// 			if err == nil {
// 				artistName = strings.TrimSpace(string(data))
// 			}
//
// 			for _, entry := range entries {
// 				p := path.Join(p, entry.Name())
//
// 				if !entry.IsDir() {
// 					continue
// 				}
//
// 				entries, err := os.ReadDir(p)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
//
// 				if !utils.HasMusic(entries) {
// 					continue
// 				}
//
// 				albumName := entry.Name()
// 				data, err := os.ReadFile(path.Join(p, "override.txt"))
// 				if err == nil {
// 					albumName = strings.TrimSpace(string(data))
// 				}
//
// 				var coverArt string
// 				for _, entry := range entries {
// 					if entry.IsDir() {
// 						continue
// 					}
//
// 					name := entry.Name()
// 					ext := path.Ext(name)
// 					if ext == "" {
// 						continue
// 					}
//
// 					name = strings.TrimSuffix(name, ext)
// 					if name != "cover" {
// 						continue
// 					}
//
// 					if utils.IsValidImageExt(ext[1:]) {
// 						coverArt = entry.Name()
// 						break
// 					}
// 				}
//
// 				tracks, err := getAllTrackFromDir(p, artistName)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
//
// 				for _, track := range tracks {
// 					if _, exists := artists[track.Artist]; !exists {
// 						artists[track.Artist] = &Artist{
// 							Name: track.Artist,
// 						}
// 					}
// 				}
//
// 				artist := artists[artistName]
//
// 				artist.Albums = append(artist.Albums, Album{
// 					Path:       p,
// 					Name:       albumName,
// 					ArtistName: artistName,
// 					CoverArt:   coverArt,
// 					Tracks:     tracks,
// 				})
// 			}
// 		}
// 	}
//
// 	return &Library{
// 		Artists: artists,
// 	}, nil
// }
//
// func GetOrCreateArtist(ctx context.Context, db *database.Database, artist *Artist) (database.Artist, error) {
// 	var dbArtist database.Artist
// 	var err error
//
// 	if artist.Path != "" {
// 		dbArtist, err = db.GetArtistByPath(ctx, artist.Path)
// 	} else {
// 		dbArtist, err = db.GetArtistByName(ctx, artist.Name)
// 	}
//
// 	if err != nil {
// 		if err == types.ErrNoArtist {
// 			artist, err := db.CreateArtist(ctx, database.CreateArtistParams{
// 				Name: artist.Name,
// 				Path: artist.Path,
// 			})
// 			if err != nil {
// 				return database.Artist{}, err
// 			}
//
// 			return artist, nil
// 		} else {
// 			return database.Artist{}, err
// 		}
// 	}
//
// 	return dbArtist, nil
// }
//
// func GetOrCreateAlbum(ctx context.Context, db *database.Database, album *Album, artistId string) (database.Album, error) {
// 	dbAlbum, err := db.GetAlbumByPath(ctx, album.Path)
// 	if err != nil {
// 		if err == types.ErrNoAlbum {
// 			album, err := db.CreateAlbum(ctx, database.CreateAlbumParams{
// 				Name:     album.Name,
// 				ArtistId: artistId,
// 				Path:     album.Path,
// 			})
//
// 			if err != nil {
// 				return database.Album{}, err
// 			}
//
// 			return album, nil
// 		} else {
// 			return database.Album{}, err
// 		}
// 	}
//
// 	return dbAlbum, nil
// }
//
// func GetOrCreateTrack(ctx context.Context, db *database.Database, track *Track, albumId string, artistId string) (database.Track, error) {
// 	dbTrack, err := db.GetTrackByPath(ctx, track.Path)
// 	if err != nil {
// 		if err == types.ErrNoTrack {
// 			track, err := db.CreateTrack(ctx, database.CreateTrackParams{
// 				TrackNumber: track.Number,
// 				Name:        track.Name,
// 				Path:        track.Path,
// 				Duration:    track.Duration,
// 				AlbumId:     albumId,
// 				ArtistId:    artistId,
// 			})
//
// 			if err != nil {
// 				return database.Track{}, err
// 			}
//
// 			return track, nil
// 		} else {
// 			return database.Track{}, err
// 		}
// 	}
//
// 	return dbTrack, nil
// }
//
// func GetOrCreateTag(ctx context.Context, db *database.Database, name string) (database.Tag, error) {
// 	tag, err := db.GetTagByName(ctx, name)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			tag, err = db.CreateTag(ctx, name)
// 			if err != nil {
// 				return database.Tag{}, err
// 			}
// 		} else {
// 			return database.Tag{}, err
// 		}
// 	}
//
// 	return tag, nil
// }
//
// func (lib *Library) syncArtists(ctx context.Context, db *database.Database, imageDir string) (map[string]database.Artist, error) {
// 	artists := make(map[string]database.Artist)
//
// 	for _, artist := range lib.Artists {
// 		fmt.Println("Syncing:", artist.Name)
//
// 		dbArtist, err := GetOrCreateArtist(ctx, db, artist)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		var artistChanges database.ArtistChanges
// 		artistChanges.Available = true
// 		artistChanges.Name.Value = artist.Name
// 		artistChanges.Name.Changed = dbArtist.Name != artist.Name
//
// 		if artist.Picture != "" {
// 			p := path.Join(artist.Path, artist.Picture)
// 			ext := path.Ext(artist.Picture)
// 			name := "artist_" + dbArtist.Id + ext
// 			dst := path.Join(imageDir, name)
// 			err := utils.SymlinkReplace(p, dst)
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			artistChanges.Picture.Value = sql.NullString{
// 				String: name,
// 				Valid:  true,
// 			}
//
// 			artistChanges.Picture.Changed = true
// 		}
//
// 		err = db.UpdateArtist(ctx, dbArtist.Id, artistChanges)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		artists[artist.Name] = dbArtist
// 	}
//
// 	return artists, nil
// }
//
// func (lib *Library) Sync(workDir types.WorkDir, dir string, db *database.Database) error {
// 	trackDir := workDir.OriginalTracksDir()
// 	err := os.MkdirAll(trackDir, 0755)
// 	if err != nil {
// 		return err
// 	}
//
// 	mobileTrackDir := workDir.MobileTracksDir()
// 	err = os.MkdirAll(mobileTrackDir, 0755)
// 	if err != nil {
// 		return err
// 	}
//
// 	transcodeDir := workDir.TranscodeDir()
// 	err = os.MkdirAll(transcodeDir, 0755)
// 	if err != nil {
// 		return err
// 	}
//
// 	imageDir := workDir.ImagesDir()
// 	err = os.MkdirAll(imageDir, 0755)
// 	if err != nil {
// 		return err
// 	}
//
// 	ctx := context.Background()
//
// 	err = db.MarkAllTracksUnavailable(ctx)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = db.MarkAllAlbumsUnavailable(ctx)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = db.MarkAllArtistsUnavailable(ctx)
// 	if err != nil {
// 		return err
// 	}
//
// 	artists, err := lib.syncArtists(ctx, db, imageDir)
//
// 	for _, artist := range lib.Artists {
// 		dbArtist := artists[artist.Name]
//
// 		for _, album := range artist.Albums {
// 			dbAlbum, err := GetOrCreateAlbum(ctx, db, &album, dbArtist.Id)
// 			if err != nil {
// 				return err
// 			}
//
// 			var albumChanges database.AlbumChanges
// 			albumChanges.Available = true
// 			albumChanges.Name.Value = album.Name
// 			albumChanges.Name.Changed = dbAlbum.Name != album.Name
//
// 			albumCoverRealPath := path.Join(album.Path, album.CoverArt)
// 			if album.CoverArt != "" {
// 				ext := path.Ext(album.CoverArt)
// 				name := "album_" + dbAlbum.Id + ext
// 				dst := path.Join(imageDir, name)
// 				err := utils.SymlinkReplace(albumCoverRealPath, dst)
// 				if err != nil {
// 					return err
// 				}
//
// 				albumChanges.CoverArt.Value = sql.NullString{
// 					String: name,
// 					Valid:  true,
// 				}
//
// 				albumChanges.CoverArt.Changed = true
//
// 				dbAlbum.CoverArt = sql.NullString{
// 					String: name,
// 					Valid:  true,
// 				}
// 			}
//
// 			err = db.UpdateAlbum(ctx, dbAlbum.Id, albumChanges)
// 			if err != nil {
// 				return err
// 			}
//
// 			for _, track := range album.Tracks {
// 				artist := artists[track.Artist]
// 				coverArt := dbAlbum.CoverArt
//
// 				dbTrack, err := GetOrCreateTrack(ctx, db, &track, dbAlbum.Id, artist.Id)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
//
// 				tags, err := db.GetTrackTags(ctx, dbTrack.Id)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
//
// 				for _, tag := range track.Tags {
// 					tag, err := GetOrCreateTag(ctx, db, tag)
// 					if err != nil {
// 						log.Fatal(err)
// 					}
//
// 					err = db.AddTagToTrack(ctx, tag.Id, dbTrack.Id)
// 					if err != nil {
// 						var pgErr *pgconn.PgError
// 						if errors.As(err, &pgErr) {
// 							if pgErr.Code == pgerrcode.UniqueViolation {
// 								continue
// 							}
// 						}
//
// 						log.Fatal(err)
// 					}
// 				}
//
// 				hasTag := func(tag string) bool {
// 					for _, t := range track.Tags {
// 						if tag == t {
// 							return true
// 						}
// 					}
//
// 					return false
// 				}
//
// 				for _, tag := range tags {
// 					if !hasTag(tag.Name) {
// 						err := db.RemoveTagFromTrack(ctx, tag.Id, dbTrack.Id)
// 						if err != nil {
// 							log.Fatal(err)
// 						}
// 					}
// 				}
//
// 				var trackChanges database.TrackChanges
// 				trackChanges.Available = true
// 				trackChanges.Number.Value = track.Number
// 				trackChanges.Number.Changed = dbTrack.Number != track.Number
// 				trackChanges.Name.Value = track.Name
// 				trackChanges.Name.Changed = dbTrack.Name != track.Name
// 				trackChanges.CoverArt.Value = coverArt
// 				trackChanges.CoverArt.Changed = dbTrack.CoverArt != coverArt
// 				trackChanges.Duration.Value = track.Duration
// 				trackChanges.Duration.Changed = dbTrack.Duration != track.Duration
//
// 				trackPath := track.Path
// 				trackExt := path.Ext(trackPath)
// 				name := dbTrack.Id + trackExt
// 				dst := path.Join(trackDir, name)
//
// 				err = utils.SymlinkReplace(trackPath, dst)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
//
// 				var mobileSrc string
//
// 				if utils.IsLossyFormatExt(trackExt) {
// 					mobileSrc = trackPath
// 				} else {
// 					transcodeName := fmt.Sprintf("%v.opus", dbTrack.Id)
// 					dstTranscode := path.Join(transcodeDir, transcodeName)
//
// 					_, err = os.Stat(dstTranscode)
// 					if err != nil {
// 						if os.IsNotExist(err) {
// 							err = parasect.RunFFmpeg(true, "-y", "-i", trackPath, "-vbr", "on", "-b:a", "128k", dstTranscode)
// 							if err != nil {
// 								log.Fatal(err)
// 							}
// 						} else {
// 							log.Fatal(err)
// 						}
// 					}
//
// 					src, err := filepath.Abs(dstTranscode)
// 					if err != nil {
// 						log.Fatal(err)
// 					}
//
// 					mobileSrc = src
// 				}
//
// 				mobileSrcExt := filepath.Ext(mobileSrc)
// 				mobileSrcName := dbTrack.Id + mobileSrcExt
// 				dst = path.Join(mobileTrackDir, mobileSrcName)
// 				err = utils.SymlinkReplace(mobileSrc, dst)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
//
// 				trackChanges.BestQualityFile.Value = name
// 				trackChanges.BestQualityFile.Changed = dbTrack.BestQualityFile != name
// 				trackChanges.MobileQualityFile.Value = mobileSrcName
// 				trackChanges.MobileQualityFile.Changed = dbTrack.MobileQualityFile != mobileSrcName
// 				trackChanges.ArtistId.Value = artist.Id
// 				trackChanges.ArtistId.Changed = dbTrack.ArtistId != artist.Id
//
// 				err = db.UpdateTrack(ctx, dbTrack.Id, trackChanges)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 			}
// 		}
// 	}
//
// 	return nil
// }

type LibraryTrack struct {
	Name              string
	Number            int
	Duration          int
	BestQualityFile   string
	MobileQualityFile string
	Tags              []string
	Genres            []string
	Artist            *LibraryArtist
	Album             *LibraryAlbum
}

type LibraryAlbum struct {
	Name     string
	CoverArt string
	Artist   *LibraryArtist
	Tracks   []*LibraryTrack
}

type LibraryArtist struct {
	Name   string
	Albums []*LibraryAlbum
}

type Library struct {
	Artists []*LibraryArtist
}

func ReadFromDir(dir string) (*Library, error) {
	var albums []string

	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if d == nil {
			return filepath.SkipDir
		}

		if d.Name() == "album.toml" {
			albums = append(albums, p)
		}

		return nil
	})

	artists := make(map[string]*LibraryArtist)

	GetOrCreateArtist := func(artistName string) *LibraryArtist {
		artist, exists := artists[artistName]
		if !exists {
			artist = &LibraryArtist{
				Name:   artistName,
				Albums: []*LibraryAlbum{},
			}

			artists[artistName] = artist
		}

		return artist
	}

	for _, p := range albums {
		base := path.Dir(p)

		d, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}

		var metadata AlbumMetadata
		err = toml.Unmarshal(d, &metadata)
		if err != nil {
			return nil, err
		}

		artist := GetOrCreateArtist(metadata.Artist)

		coverArt := ""

		if metadata.CoverArt != "" {
			coverArt = path.Join(base, metadata.CoverArt)
		}

		album := &LibraryAlbum{
			Name:     metadata.Album,
			CoverArt: coverArt,
			Artist:   artist,
			Tracks:   []*LibraryTrack{},
		}
		artist.Albums = append(artist.Albums, album)

		for _, t := range metadata.Tracks {
			// for _, featureArtist := range t.Featuring {
			// 	artist := GetOrCreateArtist(featureArtist)
			// 	_ = artist
			// }

			// TODO(patrik): Add more checks here
			bestQualityFile := t.File.Lossless
			mobileQualityFile := t.File.Lossy

			if bestQualityFile == "" {
				bestQualityFile = mobileQualityFile
			}

			artist := GetOrCreateArtist(t.Artist)

			album.Tracks = append(album.Tracks, &LibraryTrack{
				Name:              t.Name,
				Number:            t.Num,
				Duration:          t.Duration,
				BestQualityFile:   path.Join(base, bestQualityFile),
				MobileQualityFile: path.Join(base, mobileQualityFile),
				Tags:              t.Tags,
				Genres:            t.Genres,
				Artist:            artist,
				Album:             album,
			})
		}
	}

	res := make([]*LibraryArtist, 0, len(artists))
	for _, a := range artists {
		res = append(res, a)
	}

	return &Library{
		Artists: res,
	}, nil
}

type SyncContext struct {
	ArtistMapping map[*LibraryArtist]database.Artist
	AlbumMapping  map[*LibraryAlbum]database.Album
	TagMapping    map[string]database.Tag
	GenreMapping  map[string]database.Genre
}

func (sync *SyncContext) GetOrCreateArtist(ctx context.Context, db *database.Database, artist *LibraryArtist) (database.Artist, error) {
	dbArtist, err := db.GetArtistByName(ctx, artist.Name)
	if err != nil {
		if errors.Is(err, types.ErrNoArtist) {
			dbArtist, err := db.CreateArtist(ctx, database.CreateArtistParams{
				Name:    artist.Name,
				Picture: sql.NullString{},
				Path:    artist.Name,
			})

			if err != nil {
				return database.Artist{}, nil
			}

			sync.ArtistMapping[artist] = dbArtist

			return dbArtist, nil
		}

		return database.Artist{}, err
	}

	sync.ArtistMapping[artist] = dbArtist

	return dbArtist, nil
}

func (sync *SyncContext) GetOrCreateAlbum(ctx context.Context, db *database.Database, album *LibraryAlbum) (database.Album, error) {
	dbAlbum, err := db.GetAlbumByName(ctx, album.Name)
	if err != nil {
		if errors.Is(err, types.ErrNoAlbum) {
			artist, exists := sync.ArtistMapping[album.Artist]
			if !exists {
				return database.Album{}, fmt.Errorf("Artist not mapped '%s'", album.Artist.Name)
			}

			dbAlbum, err := db.CreateAlbum(ctx, database.CreateAlbumParams{
				Name:     album.Name,
				CoverArt: sql.NullString{},
				ArtistId: artist.Id,
				Path:     album.Name,
			})

			if err != nil {
				return database.Album{}, nil
			}

			sync.AlbumMapping[album] = dbAlbum

			return dbAlbum, nil
		}

		return database.Album{}, err
	}

	sync.AlbumMapping[album] = dbAlbum

	return dbAlbum, nil
}

func (sync *SyncContext) GetOrCreateTrack(ctx context.Context, db *database.Database, track *LibraryTrack) (database.Track, error) {
	album, exists := sync.AlbumMapping[track.Album]
	if !exists {
		return database.Track{}, fmt.Errorf("Album not mapped '%s'", track.Album.Name)
	}

	dbTrack, err := db.GetTrackByNameAndAlbum(ctx, track.Name, album.Id)
	if err != nil {
		if errors.Is(err, types.ErrNoTrack) {
			artist, exists := sync.ArtistMapping[track.Artist]
			if !exists {
				return database.Track{}, fmt.Errorf("Artist not mapped '%s'", track.Artist.Name)
			}

			// TODO(patrik): Path need fixing
			dbTrack, err := db.CreateTrack(ctx, database.CreateTrackParams{
				TrackNumber:       track.Number,
				Name:              track.Name,
				CoverArt:          sql.NullString{},
				Path:              artist.Id + "-" + album.Id + "-" + track.Name,
				Duration:          track.Duration,
				BestQualityFile:   "",
				MobileQualityFile: "",
				AlbumId:           album.Id,
				ArtistId:          artist.Id,
			})

			if err != nil {
				return database.Track{}, nil
			}

			return dbTrack, nil
		}

		return database.Track{}, err
	}

	return dbTrack, nil
}

func (sync *SyncContext) GetOrCreateTag(ctx context.Context, db *database.Database, tag string) (database.Tag, error) {
	dbTag, err := db.GetTagByName(ctx, tag)
	if err != nil {
		if errors.Is(err, types.ErrNoTag) {
			dbTag, err := db.CreateTag(ctx, tag)

			if err != nil {
				return database.Tag{}, nil
			}

			sync.TagMapping[tag] = dbTag

			return dbTag, nil
		}

		return database.Tag{}, err
	}

	sync.TagMapping[tag] = dbTag

	return dbTag, nil
}

func (sync *SyncContext) GetOrCreateGenre(ctx context.Context, db *database.Database, genre string) (database.Genre, error) {
	dbGenre, err := db.GetGenreByName(ctx, genre)
	if err != nil {
		if errors.Is(err, types.ErrNoGenre) {
			dbGenre, err := db.CreateGenre(ctx, genre)

			if err != nil {
				return database.Genre{}, nil
			}

			sync.GenreMapping[genre] = dbGenre

			return dbGenre, nil
		}

		return database.Genre{}, err
	}

	sync.GenreMapping[genre] = dbGenre

	return dbGenre, nil
}

func (lib *Library) Sync(workDir types.WorkDir, db *database.Database) error {
	ctx := context.Background()

	syncContext := SyncContext{
		ArtistMapping: make(map[*LibraryArtist]database.Artist),
		AlbumMapping:  make(map[*LibraryAlbum]database.Album),
		TagMapping:    make(map[string]database.Tag),
		GenreMapping:  make(map[string]database.Genre),
	}

	err := db.MarkAllArtistsUnavailable(ctx)
	if err != nil {
		return err
	}

	err = db.MarkAllAlbumsUnavailable(ctx)
	if err != nil {
		return err
	}

	err = db.MarkAllTracksUnavailable(ctx)
	if err != nil {
		return err
	}

	removeAll := func(dir string) error {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			p := path.Join(dir, entry.Name())
			err = os.Remove(p)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err = removeAll(workDir.ImagesDir())
	if err != nil {
		return err
	}

	err = removeAll(workDir.OriginalTracksDir())
	if err != nil {
		return err
	}

	err = removeAll(workDir.MobileTracksDir())
	if err != nil {
		return err
	}

	for _, artist := range lib.Artists {
		dbArtist, err := syncContext.GetOrCreateArtist(ctx, db, artist)
		if err != nil {
			return err
		}

		err = db.UpdateArtist(ctx, dbArtist.Id, database.ArtistChanges{
			Available: true,
		})
		if err != nil {
			return err
		}
	}

	for _, artist := range lib.Artists {
		for _, album := range artist.Albums {
			dbAlbum, err := syncContext.GetOrCreateAlbum(ctx, db, album)
			if err != nil {
				return err
			}

			coverArt := ""
			fmt.Printf("album.CoverArt: %v\n", album.CoverArt)
			if album.CoverArt != "" {
				p, err := filepath.Abs(album.CoverArt)
				if err != nil {
					return err
				}

				name := dbAlbum.Id + path.Ext(p)
				sym := path.Join(workDir.ImagesDir(), name)
				err = utils.SymlinkReplace(p, sym)
				if err != nil {
					return err
				}

				coverArt = name
			}

			fmt.Printf("coverArt: %v\n", coverArt)

			for _, track := range album.Tracks {
				dbTrack, err := syncContext.GetOrCreateTrack(ctx, db, track)
				if err != nil {
					return err
				}

				// db.GetTrackTags()

				originalMedia, err := filepath.Abs(track.BestQualityFile)
				if err != nil {
					return err
				}

				mobileMedia, err := filepath.Abs(track.MobileQualityFile)
				if err != nil {
					return err
				}

				fmt.Printf("originalMedia: %v\n", originalMedia)
				fmt.Printf("mobileMedia: %v\n", mobileMedia)

				originalMediaSymlink := path.Join(workDir.OriginalTracksDir(), dbTrack.Id+path.Ext(originalMedia))
				err = utils.SymlinkReplace(originalMedia, originalMediaSymlink)
				if err != nil {
					return err
				}

				mobileMediaSymlink := path.Join(workDir.MobileTracksDir(), dbTrack.Id+path.Ext(mobileMedia))
				err = utils.SymlinkReplace(mobileMedia, mobileMediaSymlink)
				if err != nil {
					return err
				}

				currentTrackTags, err := db.GetTrackTags(ctx, dbTrack.Id)
				if err != nil {
					return err
				}

				for _, tag := range track.Tags {
					hasTag := false
					for _, t := range currentTrackTags {
						if t.Name == strings.ToLower(tag) {
							hasTag = true
							break
						}
					}

					if !hasTag {
						dbTag, err := syncContext.GetOrCreateTag(ctx, db, tag)
						if err != nil {
							return err
						}

						err = db.AddTagToTrack(ctx, dbTag.Id, dbTrack.Id)
						if err != nil {
							return err
						}
					}
				}

				for _, t := range currentTrackTags {
					hasTag := false
					for _, trackTag := range track.Tags {
						if t.Name == strings.ToLower(trackTag) {
							hasTag = true
							break
						}
					}

					if !hasTag {
						err := db.RemoveTagFromTrack(ctx, t.Name, dbTrack.Id)
						if err != nil {
							return err
						}
					}
				}

				currentTrackGenres, err := db.GetTrackGenres(ctx, dbTrack.Id)
				if err != nil {
					return err
				}

				for _, genre := range track.Genres {
					hasGenre := false
					for _, g := range currentTrackGenres {
						if g.Name == strings.ToLower(genre) {
							hasGenre = true
							break
						}
					}

					if !hasGenre {
						dbGenre, err := syncContext.GetOrCreateGenre(ctx, db, genre)
						if err != nil {
							return err
						}

						db.AddGenreToTrack(ctx, dbGenre.Id, dbTrack.Id)
						if err != nil {
							return err
						}
					}
				}

				for _, g := range currentTrackGenres {
					hasGenre := false
					for _, trackGenre := range track.Genres {
						if g.Name == strings.ToLower(trackGenre) {
							hasGenre = true
							break
						}
					}

					if !hasGenre {
						err = db.RemoveGenreFromTrack(ctx, g.Id, dbTrack.Id)
						if err != nil {
							return err
						}
					}
				}

				artist, exists := syncContext.ArtistMapping[track.Artist]
				if !exists {
					return fmt.Errorf("Artist not mapped '%s'", track.Artist.Name)
				}

				// TODO(patrik): Change all 'Changed' to conditions
				changes := database.TrackChanges{}
				changes.BestQualityFile.Value = path.Base(originalMediaSymlink)
				changes.BestQualityFile.Changed = true
				changes.MobileQualityFile.Value = path.Base(mobileMediaSymlink)
				changes.MobileQualityFile.Changed = true
				changes.CoverArt.Value = sql.NullString{
					String: coverArt,
					Valid:  coverArt != "",
				}
				changes.CoverArt.Changed = true
				changes.Duration.Value = track.Duration
				changes.Duration.Changed = true
				changes.ArtistId = types.Change[string]{
					Value:   artist.Id,
					Changed: dbTrack.ArtistId != artist.Id,
				}
				changes.Available = true

				pretty.Println(changes)

				err = db.UpdateTrack(ctx, dbTrack.Id, changes)
				if err != nil {
					return err
				}
			}

			// TODO(patrik): Temporary
			artist, exists := syncContext.ArtistMapping[album.Artist]
			if !exists {
				return fmt.Errorf("Artist not mapped '%s'", album.Artist.Name)
			}

			changes := database.AlbumChanges{}
			changes.CoverArt.Value = sql.NullString{
				String: coverArt,
				Valid:  coverArt != "",
			}
			changes.CoverArt.Changed = true
			changes.ArtistId = types.Change[string]{
				Value:   artist.Id,
				Changed: dbAlbum.ArtistId != artist.Id,
			}
			changes.Available = true
			err = db.UpdateAlbum(ctx, dbAlbum.Id, changes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
