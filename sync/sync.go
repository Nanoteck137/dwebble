package sync

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/collection"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
)

// rootdir/
//   [Artist]/
//     [Albums]/
//     artist.json

var dialect = goqu.Dialect("postgres")

type ArtistUpdateRequest struct {
	Name *string `db:"name" goqu:"omitnil"`
}

func RunArtistUpdate(ctx context.Context, db *pgxpool.Pool, artistId string, req *ArtistUpdateRequest) error {
	sql, params, err := dialect.Update("artists").
		Set(req).
		Where(goqu.C("id").Eq(artistId)).
		Prepared(true).ToSQL()
	if err != nil {
		return err
	}

	pretty.Println(sql)
	pretty.Println(params)

	_, err = db.Exec(ctx, sql, params...)

	return nil
}

type AlbumUpdateRequest struct {
	Name     *string `db:"name" goqu:"omitnil"`
	ArtistId *string `db:"artist_id" goqu:"omitnil"`
}

func RunAlbumUpdate(ctx context.Context, db *pgxpool.Pool, albumId string, req *AlbumUpdateRequest) error {
	sql, params, err := dialect.Update("albums").
		Set(req).
		Where(goqu.C("id").Eq(albumId)).
		Prepared(true).ToSQL()
	if err != nil {
		return err
	}

	pretty.Println(sql)
	pretty.Println(params)

	_, err = db.Exec(ctx, sql, params...)

	return nil
}

// type FileType int
//
// const (
// 	FileTypeOriginal FileType = iota
// 	FileTypeMobile
// 	FileTypeImage
// )
//
// func fileSync(typ FileType) error {
// 	p := path.Join(col.BasePath, track.FilePath)
// 	fmt.Printf("Path: %v\n", p)
//
// 	dst := workDir.OriginalTracksDir()
// 	err := os.MkdirAll(dst, 0755)
// 	if err != nil {
// 		return err
// 	}
//
// 	ext := path.Ext(p)[1:]
// 	fileName := fmt.Sprintf("%v.%v", track.Id, ext)
// 	fileDst := path.Join(dst, fileName)
// 	test, err := filepath.Rel(dst, p)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("Test: %v\n", test)
//
// 	err = os.Symlink(test, fileDst)
// 	if err != nil {
// 		if os.IsExist(err) {
// 			err := os.Remove(fileDst)
// 			if err != nil {
// 				return err
// 			}
//
// 			err = os.Symlink(test, fileDst)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			return err
// 		}
// 	}
//
// 	return nil
// }

func SyncCollection(col *collection.Collection, db *pgxpool.Pool, workDir types.WorkDir) error {
	queries := database.New(db)

	ctx := context.Background()

	sql, _, err := dialect.From("artists").Select(goqu.C("id")).ToSQL()
	if err != nil {
		return err
	}

	rows, err := db.Query(ctx, sql)
	if err != nil {
		return err
	}

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}

		ids = append(ids, id)
	}

	pretty.Println(ids)

	for artistId, artist := range col.Artists {
		dbArtist, err := queries.GetArtist(ctx, artistId)
		if err != nil {
			if err == pgx.ErrNoRows {
				p := path.Join(col.BasePath, artist.PicturePath)
				fmt.Printf("Image: %v\n", p)

				dst := workDir.ImagesDir()
				err := os.MkdirAll(dst, 0755)
				if err != nil {
					return err
				}

				ext := path.Ext(p)[1:]
				fileName := fmt.Sprintf("%v.%v", artist.Id, ext)
				fileDst := path.Join(dst, fileName)
				test, err := filepath.Rel(dst, p)
				if err != nil {
					return err
				}
				fmt.Printf("Test: %v\n", test)

				err = os.Symlink(test, fileDst)
				if err != nil {
					if os.IsExist(err) {
						err := os.Remove(fileDst)
						if err != nil {
							return err
						}

						err = os.Symlink(test, fileDst)
						if err != nil {
							return err
						}
					} else {
						return err
					}
				}

				_, err = queries.CreateArtist(ctx, database.CreateArtistParams{
					ID:      artistId,
					Name:    artist.Name,
					Picture: fileName,
				})

				if err != nil {
					return err
				}

				continue
			}

			return err
		} else {
			updateReq := ArtistUpdateRequest{}
			if dbArtist.Name != artist.Name {
				updateReq.Name = &artist.Name
			}

			if !reflect.DeepEqual(updateReq, ArtistUpdateRequest{}) {
				RunArtistUpdate(ctx, db, dbArtist.ID, &updateReq)
			}

			_ = dbArtist
		}
	}

	for albumId, album := range col.Albums {
		dbAlbum, err := queries.GetAlbum(ctx, albumId)
		if err != nil {
			if err == pgx.ErrNoRows {
				p := path.Join(col.BasePath, album.CoverArtPath)
				fmt.Printf("Image: %v\n", p)

				dst := workDir.ImagesDir()
				err := os.MkdirAll(dst, 0755)
				if err != nil {
					return err
				}

				ext := path.Ext(p)[1:]
				fileName := fmt.Sprintf("%v.%v", album.Id, ext)
				fileDst := path.Join(dst, fileName)
				test, err := filepath.Rel(dst, p)
				if err != nil {
					return err
				}
				fmt.Printf("Test: %v\n", test)

				err = os.Symlink(test, fileDst)
				if err != nil {
					if os.IsExist(err) {
						err := os.Remove(fileDst)
						if err != nil {
							return err
						}

						err = os.Symlink(test, fileDst)
						if err != nil {
							return err
						}
					} else {
						return err
					}
				}

				_, err = queries.CreateAlbum(ctx, database.CreateAlbumParams{
					ID:       albumId,
					Name:     album.Name,
					CoverArt: fileName,
					ArtistID: album.ArtistId,
				})

				if err != nil {
					return err
				}

				continue
			}

			return err
		} else {
			updateReq := AlbumUpdateRequest{}
			if dbAlbum.Name != album.Name {
				updateReq.Name = &album.Name
			}

			if dbAlbum.ArtistID != album.ArtistId {
				updateReq.ArtistId = &album.ArtistId
			}

			pretty.Println(updateReq)

			if !reflect.DeepEqual(updateReq, ArtistUpdateRequest{}) {
				RunAlbumUpdate(ctx, db, dbAlbum.ID, &updateReq)
			}

			_ = dbAlbum
		}

	}

	for trackId, track := range col.Tracks {
		dbTrack, err := queries.GetTrack(ctx, trackId)
		if err != nil {
			if err == pgx.ErrNoRows {
				p := path.Join(col.BasePath, track.FilePath)
				fmt.Printf("Path: %v\n", p)

				dst := workDir.OriginalTracksDir()
				err := os.MkdirAll(dst, 0755)
				if err != nil {
					return err
				}

				ext := path.Ext(p)[1:]
				fileName := fmt.Sprintf("%v.%v", track.Id, ext)
				fileDst := path.Join(dst, fileName)
				test, err := filepath.Rel(dst, p)
				if err != nil {
					return err
				}
				fmt.Printf("Test: %v\n", test)

				err = os.Symlink(test, fileDst)
				if err != nil {
					if os.IsExist(err) {
						err := os.Remove(fileDst)
						if err != nil {
							return err
						}

						err = os.Symlink(test, fileDst)
						if err != nil {
							return err
						}
					} else {
						return err
					}
				}

				_, err = queries.CreateTrack(ctx, database.CreateTrackParams{
					ID:                trackId,
					TrackNumber:       int32(track.Number),
					Name:              track.Name,
					CoverArt:          "TODO",
					BestQualityFile:   fileName,
					MobileQualityFile: "TODO",
					AlbumID:           track.AlbumId,
					ArtistID:          track.ArtistId,
				})

				if err != nil {
					return err
				}

				continue
			}

			return err
		} else {
			// updateReq := ArtistUpdateRequest{}
			// if dbArtist.Name != artist.Name {
			// 	updateReq.Name = &artist.Name
			// }
			//
			// if !reflect.DeepEqual(updateReq, ArtistUpdateRequest{}) {
			// 	RunArtistUpdate(ctx, db, &updateReq)
			// }
			//
			_ = dbTrack
		}

	}

	return nil
}
