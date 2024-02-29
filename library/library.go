package library

import (
	"context"
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
}

type Album struct {
	Path       string
	Name       string
	ArtistName string
	Tracks     []Track
}

type Artist struct {
	Path        string
	Name        string
	PicturePath string
	Albums      []Album
}

type Library struct {
	Artists map[string]*Artist
}

var fileReg = regexp.MustCompile(`(\d*).*`)

func getAllTrackFromDir(dir string) ([]Track, error) {
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

			tracks = append(tracks, Track{
				Path:   p,
				Name:   name,
				Number: trackNum,
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

			if !utils.HasMusic(entries) && !utils.IsMultiDisc(entries) {
				override := path.Join(p, "override.txt")

				name := entry.Name()

				data, err := os.ReadFile(override)
				if err == nil {
					name = strings.TrimSpace(string(data))
				}

				artists[name] = &Artist{
					Path: p,
					Name: name,
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

			if utils.HasMusic(entries) {
				fmt.Printf("%v is an album\n", p)

				tracks, err := getAllTrackFromDir(p)
				if err != nil {
					log.Fatal(err)
				}

				artist := artists[artistName]

				albumName := entry.Name()
				artist.Albums = append(artist.Albums, Album{
					Path:       p,
					Name:       albumName,
					ArtistName: artistName,
					Tracks:     tracks,
				})
			} else if utils.IsMultiDisc(entries) {
				fmt.Printf("%v is an multidisc album\n", p)

				// albumName := entry.Name()
				// albums = append(albums, Album{
				// 	Name:       albumName,
				// 	ArtistName: artistName,
				// 	Tracks:     []Track{},
				// })
			} else {
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

					if utils.HasMusic(entries) {
						tracks, err := getAllTrackFromDir(p)
						if err != nil {
							log.Fatal(err)
						}

						artist := artists[artistName]

						albumName := entry.Name()
						artist.Albums = append(artist.Albums, Album{
							Path:       p,
							Name:       albumName,
							ArtistName: artistName,
							Tracks:     tracks,
						})
					} else if utils.IsMultiDisc(entries) {
						// albumName := entry.Name()
						// albums = append(albums, Album{
						// 	Name:       albumName,
						// 	ArtistName: artistName,
						// 	Tracks:     []Track{},
						// })
					} else {
						log.Printf("Warning: No music found at '%v'", p)
					}
				}
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
				Name:    artist.Name,
				Picture: "",
				Path:    artist.Path,
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
				CoverArt: "",
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
				CoverArt:          "",
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

	ctx := context.Background()

	for _, artist := range lib.Artists {
		fmt.Println("Syncing:", artist.Name)

		dbArtist, err := GetOrCreateArtist(ctx, db, artist)
		if err != nil {
			return err
		}

		for _, album := range artist.Albums {
			dbAlbum, err := GetOrCreateAlbum(ctx, db, &album, dbArtist.Id)
			if err != nil {
				return err
			}

			for _, track := range album.Tracks {
				dbTrack, err := GetOrCreateTrack(ctx, db, &track, dbAlbum.Id, dbArtist.Id)
				if err != nil {
					return err
				}

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

				err = db.UpdateTrack(ctx, dbTrack.Id, name, transcodeName)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
