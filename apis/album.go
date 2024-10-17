package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/faceair/jio"
	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	pyrinapi "github.com/nanoteck137/pyrin/api"
	vld "github.com/tiendc/go-validator"
)

type albumApi struct {
	app core.App
}

func ConvertDBAlbum(c pyrin.Context, album database.Album) types.Album {
	var year *int64
	if album.Year.Valid {
		year = &album.Year.Int64
	}

	return types.Album{
		Id:         album.Id,
		Name:       album.Name,
		CoverArt:   utils.ConvertAlbumCoverURL(c, album.Id, album.CoverArt),
		ArtistId:   album.ArtistId,
		ArtistName: album.ArtistName,
		Year:       year,
		Created:    album.Created,
		Updated:    album.Updated,
	}
}

func (api *albumApi) HandleGetAlbums(c pyrin.Context) (any, error) {
	// TODO(patrik): Add to pyrin
	// filter := c.QueryParam("filter")
	// sort := c.QueryParam("sort")
	// includeAll := ParseQueryBool(c.QueryParam("includeAll"))

	filter := ""
	sort := ""
	includeAll := false

	albums, err := api.app.DB().GetAllAlbums(c.Request().Context(), filter, sort, includeAll)
	if err != nil {
		return nil, err
	}

	res := types.GetAlbums{
		Albums: make([]types.Album, len(albums)),
	}

	for i, album := range albums {
		res.Albums[i] = ConvertDBAlbum(c, album)
	}

	return res, nil
}

func (api *albumApi) HandleGetAlbumById(c pyrin.Context) (any, error) {
	id := c.Param("id")
	album, err := api.app.DB().GetAlbumById(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return nil, AlbumNotFound()
		}

		return nil, err
	}

	return types.GetAlbumById{
		Album: ConvertDBAlbum(c, album),
	}, nil
}

func (api *albumApi) HandleGetAlbumTracksById(c pyrin.Context) (any, error) {
	id := c.Param("id")

	album, err := api.app.DB().GetAlbumById(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return nil, AlbumNotFound()
		}

		return nil, err
	}

	tracks, err := api.app.DB().GetTracksByAlbum(c.Request().Context(), album.Id)
	if err != nil {
		return nil, err
	}

	res := types.GetAlbumTracksById{
		Tracks: make([]types.Track, len(tracks)),
	}

	for i, track := range tracks {
		res.Tracks[i] = ConvertDBTrack(c, track)
	}

	return res, nil
}

var _ types.Body = (*PatchAlbumBody)(nil)

type PatchAlbumBody struct {
	Name       *string `json:"name"`
	ArtistId   *string `json:"artistId"`
	ArtistName *string `json:"artistName"`
	Year       *int64  `json:"year"`
}

func (PatchAlbumBody) Schema() jio.Schema {
	panic("unimplemented")
}

func (api *albumApi) HandlePatchAlbum(c echo.Context) error {
	id := c.Param("id")

	var body PatchAlbumBody
	d := json.NewDecoder(c.Request().Body)
	err := d.Decode(&body)
	if err != nil {
		return err
	}

	album, err := api.app.DB().GetAlbumById(c.Request().Context(), id)
	if err != nil {
		// TODO(patrik): Handle error
		return err
	}

	var name types.Change[string]
	if body.Name != nil {
		n := strings.TrimSpace(*body.Name)
		name.Value = n
		name.Changed = n != album.Name
	}

	ctx := context.TODO()

	var artistId types.Change[string]
	if body.ArtistId != nil {
		artistId.Value = *body.ArtistId
		artistId.Changed = *body.ArtistId != album.ArtistId
	} else if body.ArtistName != nil {
		artistName := strings.TrimSpace(*body.ArtistName)

		artist, err := api.app.DB().GetArtistByName(ctx, artistName)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				artist, err = api.app.DB().CreateArtist(ctx, database.CreateArtistParams{
					Name:    artistName,
					Picture: sql.NullString{},
				})

				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		artistId.Value = artist.Id
		artistId.Changed = artist.Id != album.ArtistId
	}

	var year types.Change[sql.NullInt64]
	if body.Year != nil {
		year.Value = sql.NullInt64{
			Int64: *body.Year,
			Valid: *body.Year != 0,
		}
		year.Changed = *body.Year != album.Year.Int64
	}

	err = api.app.DB().UpdateAlbum(c.Request().Context(), album.Id, database.AlbumChanges{
		Name:     name,
		ArtistId: artistId,
		Year:     year,
	})
	if err != nil {
		return err
	}

	return c.JSON(200, pyrinapi.SuccessResponse(nil))
}

// TODO(patrik): Move the album folder to trash can system
func (api *albumApi) HandleDeleteAlbum(c echo.Context) error {
	id := c.Param("id")

	db, tx, err := api.app.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = db.RemoveAlbumTracks(c.Request().Context(), id)
	if err != nil {
		return fmt.Errorf("Failed to remove album tracks: %w", err)
	}

	err = db.RemoveAlbum(c.Request().Context(), id)
	if err != nil {
		return fmt.Errorf("Failed to remove album: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return c.JSON(200, pyrinapi.SuccessResponse(nil))
}

var _ types.Body = (*PostAlbumImportBody)(nil)

type PostAlbumImportBody struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

func (PostAlbumImportBody) Schema() jio.Schema {
	panic("unimplemented")
}

type PostAlbumImport struct {
	AlbumId string `json:"albumId"`
}

func (api *albumApi) HandlePostAlbumImport(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	data := form.Value["data"][0]
	var body PostAlbumImportBody
	err = json.Unmarshal(([]byte)(data), &body)
	if err != nil {
		return err
	}

	errs := vld.Validate(
		vld.Required(&body.Name).OnError(
			vld.SetField("name", nil),
		),
		vld.Required(&body.Artist).OnError(
			vld.SetField("artist", nil),
		),
	)
	if errs != nil {
		return errs
	}

	db, tx, err := api.app.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.TODO()

	artist, err := db.GetArtistByName(ctx, body.Artist)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			artist, err = db.CreateArtist(ctx, database.CreateArtistParams{
				Name:    body.Artist,
				Picture: sql.NullString{},
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	pretty.Println(artist)

	album, err := db.CreateAlbum(ctx, database.CreateAlbumParams{
		Name:     body.Name,
		ArtistId: artist.Id,

		CoverArt: sql.NullString{},
		Year:     sql.NullInt64{},

		Available: true,
	})
	if err != nil {
		return err
	}

	albumDir := api.app.WorkDir().Album(album.Id)

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

	coverArt := form.File["coverArt"]
	if len(coverArt) > 0 {
		f := coverArt[0]

		ext := path.Ext(f.Filename)
		filename := "cover-original" + ext

		dst := path.Join(albumDir.Images(), filename)

		// TODO(patrik): Close file
		file, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		ff, err := f.Open()
		if err != nil {
			return err
		}

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

	files := form.File["files"]
	for _, f := range files {
		now := time.Now()

		// TODO(patrik): Maybe save the original filename to use when exporting
		ext := path.Ext(f.Filename)
		originalName := strings.TrimSuffix(f.Filename, ext)
		filename := originalName + "-" + strconv.FormatInt(now.Unix(), 10)

		file, err := os.CreateTemp("", "track.*"+ext)
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())

		ff, err := f.Open()
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

		name := originalName
		dateRegex := regexp.MustCompile(`^([12]\d\d\d)`)

		if tag, exists := trackInfo.Tags["title"]; exists {
			name = tag
		}

		var year sql.NullInt64
		if tag, exists := trackInfo.Tags["date"]; exists {
			match := dateRegex.FindStringSubmatch(tag)
			if len(match) > 0 {
				y, _ := strconv.Atoi(match[1])

				year.Int64 = int64(y)
				year.Valid = true
			}
		}

		var number int
		if tag, exists := trackInfo.Tags["track"]; exists {
			y, _ := strconv.Atoi(tag)
			number = y
		}

		if number == 0 {
			number = utils.ExtractNumber(originalName)
		}

		_, err = db.CreateTrack(ctx, database.CreateTrackParams{
			Name:     name,
			AlbumId:  album.Id,
			ArtistId: artist.Id,
			Number: sql.NullInt64{
				Int64: int64(number),
				Valid: number != 0,
			},
			Duration: sql.NullInt64{
				Int64: int64(trackInfo.Duration),
				Valid: true,
			},
			Year:             year,
			ExportName:       originalName,
			OriginalFilename: originalFile,
			MobileFilename:   mobileFile,
			Available:        true,
		})
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return c.JSON(200, pyrinapi.SuccessResponse(PostAlbumImport{
		AlbumId: album.Id,
	}))
}

func (api *albumApi) HandlePostAlbumImportTrackById(c echo.Context) error {
	id := c.Param("id")

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	db, tx, err := api.app.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.TODO()
	album, err := db.GetAlbumById(ctx, id)
	if err != nil {
		// TODO(patrik): Handle error
		return err
	}

	albumDir := api.app.WorkDir().Album(album.Id)

	files := form.File["files"]
	for _, f := range files {
		now := time.Now()

		// TODO(patrik): Maybe save the original filename to use when exporting
		ext := path.Ext(f.Filename)
		originalName := strings.TrimSuffix(f.Filename, ext)
		filename := originalName + "-" + strconv.FormatInt(now.Unix(), 10)

		file, err := os.CreateTemp("", "track.*"+ext)
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())

		ff, err := f.Open()
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

		name := originalName
		dateRegex := regexp.MustCompile(`^([12]\d\d\d)`)

		if tag, exists := trackInfo.Tags["title"]; exists {
			name = tag
		}

		var year sql.NullInt64
		if tag, exists := trackInfo.Tags["date"]; exists {
			match := dateRegex.FindStringSubmatch(tag)
			if len(match) > 0 {
				y, _ := strconv.Atoi(match[1])

				year.Int64 = int64(y)
				year.Valid = true
			}
		}

		var number int
		if tag, exists := trackInfo.Tags["track"]; exists {
			y, _ := strconv.Atoi(tag)
			number = y
		}

		if number == 0 {
			number = utils.ExtractNumber(originalName)
		}

		_, err = db.CreateTrack(ctx, database.CreateTrackParams{
			Name:     name,
			AlbumId:  album.Id,
			ArtistId: album.ArtistId,
			Number: sql.NullInt64{
				Int64: int64(number),
				Valid: number != 0,
			},
			Duration: sql.NullInt64{
				Int64: int64(trackInfo.Duration),
				Valid: true,
			},
			Year:             year,
			ExportName:       originalName,
			OriginalFilename: originalFile,
			MobileFilename:   mobileFile,
			Available:        true,
		})
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return c.JSON(200, pyrinapi.SuccessResponse(nil))
}

func InstallAlbumHandlers(app core.App, group pyrin.Group) {
	a := albumApi{app: app}

	group.Register(
		pyrin.ApiHandler{
			Name:        "GetAlbums",
			Path:        "/albums",
			Method:      http.MethodGet,
			DataType:    types.GetAlbums{},
			HandlerFunc: a.HandleGetAlbums,
		},

		pyrin.ApiHandler{
			Name:        "GetAlbumById",
			Method:      http.MethodGet,
			Path:        "/albums/:id",
			DataType:    types.GetAlbumById{},
			Errors:      []ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: a.HandleGetAlbumById,
		},

		pyrin.ApiHandler{
			Name:        "GetAlbumTracks",
			Method:      http.MethodGet,
			Path:        "/albums/:id/tracks",
			DataType:    types.GetAlbumTracksById{},
			Errors:      []ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: a.HandleGetAlbumTracksById,
		},
	)

	// group.Register(
	// 	Handler{
	// 		Name:        "EditAlbum",
	// 		Method:      http.MethodPatch,
	// 		Path:        "/albums/:id",
	// 		DataType:    nil,
	// 		BodyType:    PatchAlbumBody{},
	// 		HandlerFunc: a.HandlePatchAlbum,
	// 	},
	//
	// 	Handler{
	// 		Name:        "DeleteAlbum",
	// 		Method:      http.MethodDelete,
	// 		Path:        "/albums/:id",
	// 		DataType:    nil,
	// 		BodyType:    nil,
	// 		HandlerFunc: a.HandleDeleteAlbum,
	// 	},
	//
	// 	Handler{
	// 		Name:        "ImportAlbum",
	// 		Method:      http.MethodPost,
	// 		Path:        "/albums/import",
	// 		DataType:    PostAlbumImport{},
	// 		BodyType:    PostAlbumImportBody{},
	// 		IsMultiForm: true,
	// 		HandlerFunc: a.HandlePostAlbumImport,
	// 	},
	//
	// 	Handler{
	// 		Name:        "ImportTrackToAlbum",
	// 		Method:      http.MethodPost,
	// 		Path:        "/albums/:id/import/track",
	// 		DataType:    nil,
	// 		BodyType:    nil,
	// 		IsMultiForm: true,
	// 		HandlerFunc: a.HandlePostAlbumImportTrackById,
	// 	},
	// )
}
