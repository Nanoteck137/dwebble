package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
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
	"github.com/nanoteck137/pyrin/api"
	pyrinapi "github.com/nanoteck137/pyrin/api"
	vld "github.com/tiendc/go-validator"
)

type albumApi struct {
	app core.App
}

func (api *albumApi) HandleGetAlbums(c echo.Context) error {
	filter := c.QueryParam("filter")
	sort := c.QueryParam("sort")
	includeAll := ParseQueryBool(c.QueryParam("includeAll"))

	albums, err := api.app.DB().GetAllAlbums(c.Request().Context(), filter, sort, includeAll)
	if err != nil {
		return err
	}

	res := types.GetAlbums{
		Albums: make([]types.Album, len(albums)),
	}

	for i, album := range albums {
		res.Albums[i] = types.Album{
			Id:         album.Id,
			Name:       album.Name,
			CoverArt:   utils.ConvertAlbumCoverURL(c, album.CoverArt),
			ArtistId:   album.ArtistId,
			ArtistName: album.ArtistName,
		}
	}

	return c.JSON(200, SuccessResponse(res))
}

func (api *albumApi) HandleGetAlbumById(c echo.Context) error {
	id := c.Param("id")
	album, err := api.app.DB().GetAlbumById(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return AlbumNotFound()
		}

		return err
	}

	return c.JSON(200, SuccessResponse(types.GetAlbumById{
		Album: types.Album{
			Id:       album.Id,
			Name:     album.Name,
			CoverArt: utils.ConvertAlbumCoverURL(c, album.CoverArt),
			ArtistId: album.ArtistId,
		},
	}))
}

func (api *albumApi) HandleGetAlbumTracksById(c echo.Context) error {
	id := c.Param("id")

	album, err := api.app.DB().GetAlbumById(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return AlbumNotFound()
		}

		return err
	}

	tracks, err := api.app.DB().GetTracksByAlbum(c.Request().Context(), album.Id)
	if err != nil {
		return err
	}

	res := types.GetAlbumTracksById{
		Tracks: make([]types.Track, len(tracks)),
	}

	for i, track := range tracks {
		res.Tracks[i] = ConvertDBTrack(c, track)
	}

	return c.JSON(200, SuccessResponse(res))
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

// TODO(patrik): Move this
func processMobileVersion(input string, outputDir, name string) (string, error) {
	inputExt := path.Ext(input)
	isLossyInput := inputExt == ".opus" || inputExt == ".mp3"

	outputExt := ".opus"
	if isLossyInput {
		outputExt = inputExt
	}

	filename := utils.Slug(name) + outputExt

	var args []string
	args = append(args, "-i", input, "-map_metadata", "-1")

	if outputExt == inputExt {
		args = append(args, "-codec", "copy")
	}

	args = append(args, path.Join(outputDir, filename))

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return filename, nil
}

func processOriginalVersion(input string, outputDir, name string) (string, utils.TrackInfo, error) {
	inputExt := path.Ext(input)
	// isLossyInput := inputExt == ".opus"

	outputExt := inputExt
	switch inputExt {
	case ".wav":
		outputExt = ".flac"
	}

	// .opus -> .opus
	// .mp3  -> .mp3
	// .wav  -> .flac
	// .flac -> .flac

	filename := utils.Slug(name) + outputExt

	var args []string
	args = append(args, "-i", input, "-map_metadata", "-1")

	if outputExt == inputExt {
		args = append(args, "-codec", "copy")
	}

	out := path.Join(outputDir, filename)
	args = append(args, out)

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", utils.TrackInfo{}, err
	}

	info, err := utils.GetTrackInfo(input)
	if err != nil {
		return "", utils.TrackInfo{}, err
	}

	return filename, info, nil
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
		Name:      body.Name,
		ArtistId:  artist.Id,

		CoverArt:  sql.NullString{},
		Year:      sql.NullInt64{},

		Available: false,
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

		mobileFile, err := processMobileVersion(file.Name(), albumDir.MobileFiles(), filename)
		if err != nil {
			return err
		}

		originalFile, trackInfo, err := processOriginalVersion(file.Name(), albumDir.OriginalFiles(), filename)
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
				year.Valid = true;
			}
		}

		_, err = db.CreateTrack(ctx, database.CreateTrackParams{
			Name:             name,
			AlbumId:          album.Id,
			ArtistId:         artist.Id,
			Number:           sql.NullInt64{},
			Duration:         sql.NullInt64{
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

func InstallAlbumHandlers(app core.App, group Group) {
	a := albumApi{app: app}

	group.Register(
		Handler{
			Name:        "GetAlbums",
			Path:        "/albums",
			Method:      http.MethodGet,
			DataType:    types.GetAlbums{},
			BodyType:    nil,
			HandlerFunc: a.HandleGetAlbums,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "GetAlbumById",
			Method:      http.MethodGet,
			Path:        "/albums/:id",
			DataType:    types.GetAlbumById{},
			BodyType:    nil,
			Errors:      []api.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: a.HandleGetAlbumById,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "GetAlbumTracks",
			Method:      http.MethodGet,
			Path:        "/albums/:id/tracks",
			DataType:    types.GetAlbumTracksById{},
			BodyType:    nil,
			Errors:      []api.ErrorType{ErrTypeAlbumNotFound},
			HandlerFunc: a.HandleGetAlbumTracksById,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "ImportAlbum",
			Method:      http.MethodPost,
			Path:        "/albums/import",
			DataType:    PostAlbumImport{},
			BodyType:    PostAlbumImportBody{},
			IsMultiForm: true,
			HandlerFunc: a.HandlePostAlbumImport,
		},
	)
}
