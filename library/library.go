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

	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/pelletier/go-toml/v2"
)

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
	ctx context.Context
	db  *database.Database

	ArtistMapping map[*LibraryArtist]database.Artist
	AlbumMapping  map[*LibraryAlbum]database.Album
	TagMapping    map[string]database.Tag
	GenreMapping  map[string]database.Genre
}

func (sync *SyncContext) GetOrCreateArtist(artist *LibraryArtist) (database.Artist, error) {
	p := artist.Name
	dbArtist, err := sync.db.GetArtistByName(sync.ctx, p)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			dbArtist, err := sync.db.CreateArtist(sync.ctx, database.CreateArtistParams{
				Name:    artist.Name,
				Picture: sql.NullString{},
			})

			if err != nil {
				return database.Artist{}, err
			}

			sync.ArtistMapping[artist] = dbArtist

			return dbArtist, nil
		}

		return database.Artist{}, err
	}

	sync.ArtistMapping[artist] = dbArtist

	return dbArtist, nil
}

func (sync *SyncContext) GetOrCreateAlbum(album *LibraryAlbum) (database.Album, error) {
	artist, exists := sync.ArtistMapping[album.Artist]
	if !exists {
		return database.Album{}, fmt.Errorf("Artist not mapped '%s'", album.Artist.Name)
	}

	p := artist.Id + "-" + album.Name
	dbAlbum, err := sync.db.GetAlbumByPath(sync.ctx, p)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			dbAlbum, err := sync.db.CreateAlbum(sync.ctx, database.CreateAlbumParams{
				Name:      album.Name,
				ArtistId:  artist.Id,
				Available: false,
			})

			if err != nil {
				return database.Album{}, err
			}

			sync.AlbumMapping[album] = dbAlbum

			return dbAlbum, nil
		}

		return database.Album{}, err
	}

	sync.AlbumMapping[album] = dbAlbum

	return dbAlbum, nil
}

func (sync *SyncContext) GetOrCreateTrack(track *LibraryTrack) (database.Track, error) {
	album, exists := sync.AlbumMapping[track.Album]
	if !exists {
		return database.Track{}, fmt.Errorf("Album not mapped '%s'", track.Album.Name)
	}

	artist, exists := sync.ArtistMapping[track.Artist]
	if !exists {
		return database.Track{}, fmt.Errorf("Artist not mapped '%s'", track.Artist.Name)
	}

	dbTrack, err := sync.db.GetTrackByNameAndAlbum(sync.ctx, track.Name, album.Id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			id, err := sync.db.CreateTrack(sync.ctx, database.CreateTrackParams{
				Name:              track.Name,
				AlbumId:           album.Id,
				ArtistId:          artist.Id,
			})
			if err != nil {
				return database.Track{}, err
			}

			dbTrack, err := sync.db.GetTrackById(sync.ctx, id)
			if err != nil {
				return database.Track{}, err
			}

			return dbTrack, nil
		}

		return database.Track{}, err
	}

	return dbTrack, nil
}

func (sync *SyncContext) GetOrCreateTag(tag string) (database.Tag, error) {
	dbTag, err := sync.db.GetTagByName(sync.ctx, tag)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			dbTag, err := sync.db.CreateTag(sync.ctx, tag)

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

func (sync *SyncContext) GetOrCreateGenre(genre string) (database.Genre, error) {
	dbGenre, err := sync.db.GetGenreByName(sync.ctx, genre)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			dbGenre, err := sync.db.CreateGenre(sync.ctx, genre)

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
		ctx:           ctx,
		db:            db,
		ArtistMapping: make(map[*LibraryArtist]database.Artist),
		AlbumMapping:  make(map[*LibraryAlbum]database.Album),
		TagMapping:    make(map[string]database.Tag),
		GenreMapping:  make(map[string]database.Genre),
	}

	// err := db.MarkAllArtistsUnavailable(ctx)
	// if err != nil {
	// 	return err
	// }

	err := db.MarkAllAlbumsUnavailable(ctx)
	if err != nil {
		return err
	}

	err = db.MarkAllTracksUnavailable(ctx)
	if err != nil {
		return err
	}

	// removeAll := func(dir string) error {
	// 	entries, err := os.ReadDir(dir)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	for _, entry := range entries {
	// 		p := path.Join(dir, entry.Name())
	// 		err = os.Remove(p)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	//
	// 	return nil
	// }

	// err = removeAll(workDir.ImagesDir())
	// if err != nil {
	// 	return err
	// }
	//
	// err = removeAll(workDir.OriginalTracksDir())
	// if err != nil {
	// 	return err
	// }
	//
	// err = removeAll(workDir.MobileTracksDir())
	// if err != nil {
	// 	return err
	// }

	for _, artist := range lib.Artists {
		_, err := syncContext.GetOrCreateArtist(artist)
		if err != nil {
			return err
		}

		// err = db.UpdateArtist(ctx, dbArtist.Id, database.ArtistChanges{
		// 	Available: true,
		// })
		// if err != nil {
		// 	return err
		// }
	}

	for _, artist := range lib.Artists {
		for _, album := range artist.Albums {
			dbAlbum, err := syncContext.GetOrCreateAlbum(album)
			if err != nil {
				return err
			}

			coverArt := ""
			// if album.CoverArt != "" {
			// 	p, err := filepath.Abs(album.CoverArt)
			// 	if err != nil {
			// 		return err
			// 	}
			//
			// 	name := dbAlbum.Id + path.Ext(p)
			// 	sym := path.Join(workDir.ImagesDir(), name)
			// 	err = utils.SymlinkReplace(p, sym)
			// 	if err != nil {
			// 		return err
			// 	}
			//
			// 	coverArt = name
			// }

			for _, track := range album.Tracks {
				dbTrack, err := syncContext.GetOrCreateTrack(track)
				if err != nil {
					return err
				}

				// originalMedia, err := filepath.Abs(track.BestQualityFile)
				// if err != nil {
				// 	return err
				// }
				//
				// mobileMedia, err := filepath.Abs(track.MobileQualityFile)
				// if err != nil {
				// 	return err
				// }

				// originalMediaSymlink := path.Join(workDir.OriginalTracksDir(), dbTrack.Id+path.Ext(originalMedia))
				// err = utils.SymlinkReplace(originalMedia, originalMediaSymlink)
				// if err != nil {
				// 	return err
				// }
				//
				// mobileMediaSymlink := path.Join(workDir.MobileTracksDir(), dbTrack.Id+path.Ext(mobileMedia))
				// err = utils.SymlinkReplace(mobileMedia, mobileMediaSymlink)
				// if err != nil {
				// 	return err
				// }

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
						dbTag, err := syncContext.GetOrCreateTag(tag)
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
						dbGenre, err := syncContext.GetOrCreateGenre(genre)
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

				_ = artist

				// TODO(patrik): Change all 'Changed' to conditions
				changes := database.TrackChanges{}
				// changes.BestQualityFile.Value = path.Base(originalMediaSymlink)
				// changes.BestQualityFile.Changed = true
				// changes.MobileQualityFile.Value = path.Base(mobileMediaSymlink)
				// changes.MobileQualityFile.Changed = true
				// changes.CoverArt.Value = sql.NullString{
				// 	String: coverArt,
				// 	Valid:  coverArt != "",
				// }
				// changes.CoverArt.Changed = true
				// changes.Duration.Value = track.Duration
				// changes.Duration.Changed = true
				// changes.ArtistId = types.Change[string]{
				// 	Value:   artist.Id,
				// 	Changed: dbTrack.ArtistId != artist.Id,
				// }
				// changes.Available = true

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
