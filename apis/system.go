package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type systemApi struct {
	app core.App
}

func (api *systemApi) HandleGetSystemInfo(c echo.Context) error {
	return c.JSON(200, SuccessResponse(types.GetSystemInfo{
		Version: dwebble.Version,
		IsSetup: api.app.IsSetup(),
	}))
}

type ExportArtist struct {
	Id      string
	Name    string
	Picture string
}

type ExportTrack struct {
	Id   string
	Name string

	AlbumId  string
	ArtistId string

	Number   int64
	Duration int64
	Year     int64

	ExportName       string
	OriginalFilename string
	MobileFilename   string

	Created int64
	Updated int64

	Tags   []string
	Genres []string
}

type ExportAlbum struct {
	Id       string
	Name     string
	ArtistId string

	CoverArt string
	Year     int64
}

type Export struct {
	Artists []ExportArtist
	Albums  []ExportAlbum
	Tracks  []ExportTrack
}

func (api *systemApi) HandlePostSystemExport(c echo.Context) error {
	db := api.app.DB()

	ctx := context.TODO()

	export := Export{}

	artists, err := db.GetAllArtists(ctx)
	if err != nil {
		return err
	}

	export.Artists = make([]ExportArtist, len(artists))
	for i, artist := range artists {
		export.Artists[i] = ExportArtist{
			Id:      artist.Id,
			Name:    artist.Name,
			Picture: artist.Picture.String,
		}
	}

	albums, err := db.GetAllAlbums(ctx, "", "", false)
	if err != nil {
		return err
	}

	export.Albums = make([]ExportAlbum, len(albums))
	for i, album := range albums {
		export.Albums[i] = ExportAlbum{
			Id:       album.Id,
			Name:     album.Name,
			ArtistId: album.ArtistId,
			CoverArt: album.CoverArt.String,
			Year:     album.Year.Int64,
		}
	}

	tracks, err := db.GetAllTracks(ctx, "", "", false)
	if err != nil {
		return err
	}

	export.Tracks = make([]ExportTrack, len(tracks))
	for i, track := range tracks {
		export.Tracks[i] = ExportTrack{
			Id:               track.Id,
			Name:             track.Name,
			AlbumId:          track.AlbumId,
			ArtistId:         track.ArtistId,
			Number:           track.Number.Int64,
			Duration:         track.Duration.Int64,
			Year:             track.Year.Int64,
			ExportName:       track.ExportName,
			OriginalFilename: track.OriginalFilename,
			MobileFilename:   track.MobileFilename,
			Created:          track.Created,
			Updated:          track.Updated,
			Tags:             utils.SplitString(track.Tags.String),
			Genres:           utils.SplitString(track.Genres.String),
		}
	}

	pretty.Println(export)

	d, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(api.app.WorkDir().ExportFile(), d, 0644)
	if err != nil {
		return err
	}

	return c.JSON(200, SuccessResponse(nil))
}

func (api *systemApi) HandlePostSystemImport(c echo.Context) error {
	db, tx, err := api.app.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	d, err := os.ReadFile(api.app.WorkDir().ExportFile())
	if err != nil {
		return err
	}

	var export Export
	err = json.Unmarshal(d, &export)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	for _, artist := range export.Artists {
		_, err := db.CreateArtist(ctx, database.CreateArtistParams{
			Id:   artist.Id,
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
			}

			return err
		}
	}

	for _, album := range export.Albums {
		_, err := db.CreateAlbum(ctx, database.CreateAlbumParams{
			Id:       album.Id,
			Name:     album.Name,
			ArtistId: album.ArtistId,
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

			return err
		}
	}

	for _, track := range export.Tracks {
		_, err := db.CreateTrack(ctx, database.CreateTrackParams{
			Id:               track.Id,
			Name:             track.Name,
			AlbumId:          track.AlbumId,
			ArtistId:         track.ArtistId,
			Number:           sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Duration:         sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Year:             sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
			ExportName:       track.ExportName,
			OriginalFilename: track.OriginalFilename,
			MobileFilename:   track.MobileFilename,
			Created:          track.Created,
			Updated:          track.Updated,
			Available:        true,
		})
		if err != nil {
			var e sqlite3.Error
			if errors.As(err, &e) {
				if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
					continue
				}
			}

			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return c.JSON(200, SuccessResponse(nil))
}

func InstallSystemHandlers(app core.App, group Group) {
	api := systemApi{app: app}

	group.Register(
		Handler{
			Name:        "GetSystemInfo",
			Path:        "/system/info",
			Method:      http.MethodGet,
			DataType:    types.GetSystemInfo{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetSystemInfo,
		},

		Handler{
			Name:        "SystemExport",
			Path:        "/system/export",
			Method:      http.MethodPost,
			DataType:    nil,
			BodyType:    nil,
			HandlerFunc: api.HandlePostSystemExport,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "SystemImport",
			Path:        "/system/import",
			Method:      http.MethodPost,
			DataType:    nil,
			BodyType:    nil,
			HandlerFunc: api.HandlePostSystemImport,
			Middlewares: []echo.MiddlewareFunc{},
		},
	)
}
