package library

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
)

type Track struct {
	Path   string
	Name   string
	Number int
	Artist string
	Tags   []string
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
			res, err := utils.CheckFile(p)
			if err != nil {
				continue
			}

			pretty.Println(res)

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
			tagTitle := res.Tags["title"]
			if tagTitle != "" {
				name = tagTitle
			} else {
				ext := path.Ext(p)
				name = strings.TrimSuffix(entry.Name(), ext)
			}

			artist := albumArtist
			tagArtist := res.Tags["artist"]
			if tagArtist != "" {
				artist = tagArtist
			}

			tags := strings.Split(res.Tags["tags"], ",")

			var realTags []string
			for _, tag := range tags {
				if tag == "" {
					continue
				}

				realTags = append(realTags, strings.TrimSpace(tag))
			}

			tracks = append(tracks, Track{
				Path:   p,
				Name:   name,
				Number: trackNum,
				Artist: artist,
				Tags:   realTags,
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

	pretty.Println(artists)

	for _, entry := range entries {
		p := path.Join(dir, entry.Name())
		if entry.IsDir() {
			fmt.Printf("Valid: %v\n", p)

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
			fmt.Printf("%v is an artist\n", p)

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

				artist := artists[artistName]

				artist.Albums = append(artist.Albums, Album{
					Path:       p,
					Name:       albumName,
					ArtistName: artistName,
					CoverArt:   coverArt,
					Tracks:     tracks,
				})
			}

			fmt.Printf("artistName: %v\n", artistName)
			fmt.Println()
		}
	}

	// pretty.Println(albums)

	return &Library{
		Artists: artists,
	}, nil
}

func GetOrCreateArtist(ctx context.Context, db *database.Database, artist *Artist) (database.Artist, error) {
	dbArtist, err := db.GetArtistByPath(ctx, artist.Path)
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
				TrackNumber:       track.Number,
				Name:              track.Name,
				Path:              track.Path,
				BestQualityFile:   "",
				MobileQualityFile: "",
				AlbumId:           albumId,
				ArtistId:          artistId,
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

	artists := make(map[string]database.Artist)

	for _, artist := range lib.Artists {
		fmt.Println("Syncing:", artist.Name)

		dbArtist, err := GetOrCreateArtist(ctx, db, artist)
		if err != nil {
			return err
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
				return err
			}

			artistChanges.Picture.Value = sql.NullString{
				String: name,
				Valid:  true,
			}

			artistChanges.Picture.Changed = true
		}

		err = db.UpdateArtist(ctx, dbArtist.Id, artistChanges)
		if err != nil {
			return err
		}

		artists[artist.Name] = dbArtist
	}

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

			if album.CoverArt != "" {
				p := path.Join(album.Path, album.CoverArt)
				ext := path.Ext(album.CoverArt)
				name := "album_" + dbAlbum.Id + ext
				dst := path.Join(imageDir, name)
				err := utils.SymlinkReplace(p, dst)
				if err != nil {
					return err
				}

				albumChanges.CoverArt.Value = sql.NullString{
					String: name,
					Valid:  true,
				}

				albumChanges.CoverArt.Changed = true
			}

			err = db.UpdateAlbum(ctx, dbAlbum.Id, albumChanges)
			if err != nil {
				return err
			}

			for _, track := range album.Tracks {
				artist := artists[track.Artist]
				dbTrack, err := GetOrCreateTrack(ctx, db, &track, dbAlbum.Id, artist.Id)
				if err != nil {
					return err
				}

				var trackChanges database.TrackChanges
				trackChanges.Available = true
				trackChanges.Number.Value = track.Number
				trackChanges.Number.Changed = dbTrack.Number != track.Number
				trackChanges.Name.Value = track.Name
				trackChanges.Name.Changed = dbTrack.Name != track.Name

				p := track.Path
				ext := path.Ext(p)
				name := fmt.Sprintf("%v%v", dbTrack.Id, ext)
				dst := path.Join(trackDir, name)
				fmt.Printf("p: %v\n", p)

				err = utils.SymlinkReplace(p, dst)

				transcodeName := fmt.Sprintf("%v.mp3", dbTrack.Id)
				dstTranscode := path.Join(transcodeDir, transcodeName)

				_, err = os.Stat(dstTranscode)
				if err != nil {
					if os.IsNotExist(err) {
						// TODO(patrik): Temp
						err := utils.RunFFmpeg(true, "-y", "-i", p, dstTranscode)
						if err != nil {
							return err
						}
					} else {
						return err
					}
				}

				src, err := filepath.Rel(mobileTrackDir, dstTranscode)
				if err != nil {
					return err
				}

				dst = path.Join(mobileTrackDir, transcodeName)
				err = utils.SymlinkReplace(src, dst)

				trackChanges.BestQualityFile.Value = name
				trackChanges.BestQualityFile.Changed = dbTrack.BestQualityFile != name
				trackChanges.MobileQualityFile.Value = transcodeName
				trackChanges.MobileQualityFile.Changed = dbTrack.MobileQualityFile != transcodeName

				err = db.UpdateTrack(ctx, dbTrack.Id, trackChanges)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
