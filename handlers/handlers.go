package handlers

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nanoteck137/dwebble/v2/database"
	"github.com/nanoteck137/dwebble/v2/types"
	"github.com/nanoteck137/dwebble/v2/utils"
)

const defaultArtistImage = "default_artist.png"
const defaultAlbumImage = "default_album.png"
const defaultTrackImage = defaultAlbumImage

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (apiError ApiError) Error() string {
	return apiError.Message
}

type ApiData struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

type ApiConfig struct {
	workDir  string
	validate *validator.Validate

	db      *pgxpool.Pool
	queries *database.Queries
}

func New(db *pgxpool.Pool, queries *database.Queries, workDir string) ApiConfig {
	var validate = validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return ApiConfig{
		workDir:  workDir,
		queries:  queries,
		db:       db,
		validate: validate,
	}
}

func (api *ApiConfig) validateBody(body any) map[string]string {
	err := api.validate.Struct(body)
	if err != nil {
		type ValidationError struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}

		validationErrs := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			validationErrs[field] = fmt.Sprintf("'%v' not satisfying tags '%v'", field, err.Tag())
		}

		return validationErrs
	}

	return nil
}

func (apiConfig *ApiConfig) HandlerGetAllArtists(c *fiber.Ctx) error {
	artists, err := apiConfig.queries.GetAllArtists(c.Context())
	if err != nil {
		return err
	}

	for i := range artists {
		a := &artists[i]
		a.Picture = ConvertURL(c, a.Picture)
	}

	return c.JSON(fiber.Map{"artists": artists})
}

type CreateArtistBody struct {
	Name string `json:"name" form:"name" validate:"required"`
}

func (api *ApiConfig) HandlerCreateArtist(c *fiber.Ctx) error {
	var body CreateArtistBody
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	errs := api.validateBody(body)
	if errs != nil {
		return ApiError{
			Status:  400,
			Message: "Failed to create artist",
			Data:    errs,
		}
	}

	artist, err := api.queries.CreateArtist(c.Context(), database.CreateArtistParams{
		ID:      utils.CreateId(),
		Name:    body.Name,
		Picture: defaultArtistImage,
	})

	if err != nil {
		return err
	}

	artist.Picture = ConvertURL(c, fmt.Sprintf("/images/%v", artist.Picture))

	return c.JSON(ApiData{
		Status: 200,
		Data:   artist,
	})
}

func (apiConfig *ApiConfig) HandlerGetArtist(c *fiber.Ctx) error {
	id := c.Params("id")

	artist, err := apiConfig.queries.GetArtist(c.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": fmt.Sprintf("No artist with id: %s", id)})
		} else {
			return err
		}
	}

	artist.Picture = ConvertURL(c, artist.Picture)

	albums, err := apiConfig.queries.GetAlbumsByArtist(c.Context(), id)
	if err != nil {
		return err
	}

	for i := range albums {
		a := &albums[i]
		a.CoverArt = ConvertURL(c, a.CoverArt)
	}

	res := struct {
		Artist database.Artist  `json:"artist"`
		Albums []database.Album `json:"albums"`
	}{
		Artist: artist,
		Albums: albums,
	}

	return c.JSON(res)
}

func (apiConfig *ApiConfig) HandlerGetAllAlbums(c *fiber.Ctx) error {
	albums, err := apiConfig.queries.GetAllAlbums(c.Context())
	if err != nil {
		return err
	}

	for i := range albums {
		a := &albums[i]
		a.CoverArt = ConvertURL(c, a.CoverArt)
	}

	return c.JSON(fiber.Map{"albums": albums})
}

type CreateAlbumBody struct {
	Name   string `json:"name" form:"name" validate:"required"`
	Artist string `json:"artist" form:"artist" validate:"required"`
}

func (api *ApiConfig) HandlerCreateAlbum(c *fiber.Ctx) error {
	var body CreateAlbumBody
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	errs := api.validateBody(body)
	if errs != nil {
		return ApiError{
			Status:  400,
			Message: "Failed to create album",
			Data:    errs,
		}
	}

	album, err := api.queries.CreateAlbum(c.Context(), database.CreateAlbumParams{
		ID:       utils.CreateId(),
		Name:     body.Name,
		CoverArt: defaultAlbumImage,
		ArtistID: body.Artist,
	})

	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			switch err.Code {
			case "23503":
				switch err.ConstraintName {
				case "albums_artist_id_fk":
					return ApiError{
						Status:  400,
						Message: fmt.Sprintf("No artist with id: '%v'", body.Artist),
						Data:    nil,
					}
				}
			}
		}

		return err
	}

	album.CoverArt = ConvertURL(c, fmt.Sprintf("/images/%v", album.CoverArt))

	return c.JSON(ApiData{
		Status: 200,
		Data:   album,
	})
}

func (apiConfig *ApiConfig) HandlerGetAlbum(c *fiber.Ctx) error {
	id := c.Params("id")

	album, err := apiConfig.queries.GetAlbum(c.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": fmt.Sprintf("No album with id: %s", id)})
		} else {
			return err
		}
	}

	album.CoverArt = ConvertURL(c, album.CoverArt)

	tracks, err := apiConfig.queries.GetTracksByAlbum(c.Context(), id)
	if err != nil {
		return err
	}

	for i := range tracks {
		t := &tracks[i]
		t.CoverArt = ConvertURL(c, t.CoverArt)
		t.BestQualityFile = ConvertURL(c, "/tracks/"+t.BestQualityFile)
		t.MobileQualityFile = ConvertURL(c, "/tracks/"+t.MobileQualityFile)
	}

	res := struct {
		Album  database.Album                 `json:"album"`
		Tracks []database.GetTracksByAlbumRow `json:"tracks"`
	}{
		Album:  album,
		Tracks: tracks,
	}

	return c.JSON(res)
}

func ConvertURL(c *fiber.Ctx, path string) string {
	return c.BaseURL() + path
}

func (apiConfig *ApiConfig) HandlerGetAllTracks(c *fiber.Ctx) error {
	tracks, err := apiConfig.queries.GetAllTracks(c.Context())
	if err != nil {
		return err
	}

	for i := range tracks {
		track := &tracks[i]
		track.CoverArt = ConvertURL(c, "/images/"+track.CoverArt)
		track.BestQualityFile = ConvertURL(c, "/tracks/"+track.BestQualityFile)
		track.MobileQualityFile = ConvertURL(c, "/tracks/"+track.MobileQualityFile)
	}

	return c.JSON(fiber.Map{"tracks": tracks})
}

type CreateTrackBody struct {
	Name   string `json:"name" form:"name" validate:"required"`
	Number int    `json:"number" form:"number" validate:"required"`
	Album  string `json:"album" form:"album" validate:"required"`
	Artist string `json:"artist" form:"artist" validate:"required"`
}

func getExtForBestQuality(contentType string) (string, error) {
	switch contentType {
	case "audio/x-flac", "audio/flac":
		return "flac", nil
	case "audio/ogg":
		return "ogg", nil
	case "audio/mpeg":
		return "mp3", nil
	case "":
		return "", errors.New("No 'Content-Type' is not set")
	default:
		return "", fmt.Errorf("Unsupported Content-Type: %v", contentType)
	}
}

func checkMobileQualityType(file *multipart.FileHeader) error {
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		return ApiError{
			Status:  400,
			Message: "Failed to create track",
			Data: map[string]any{
				"mobileQualityFile": "No Content-Type set",
			},
		}
	} else if contentType != "audio/mpeg" {
		return ApiError{
			Status:  400,
			Message: "Failed to create track",
			Data: map[string]any{
				"mobileQualityFile": "Content-Type needs to be audio/mpeg",
			},
		}
	}

	return nil
}

func (api *ApiConfig) writeTrackFile(file *multipart.FileHeader, name string) error {
	fileName := path.Join(api.workDir, "tracks", name)
	return utils.WriteFormFile(file, fileName)
}

func (api *ApiConfig) writeImageFile(file *multipart.FileHeader, name string) error {
	fileName := path.Join(api.workDir, "images", name)
	return utils.WriteFormFile(file, fileName)
}

func (api *ApiConfig) HandlerCreateTrack(c *fiber.Ctx) error {
	tx, err := api.db.BeginTx(c.Context(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(c.Context())

	queries := database.New(tx)

	form, err := c.MultipartForm()
	if err != nil {
		return ApiError{
			Status:  400,
			Message: fmt.Sprintf("Failed to create track: %v", err),
		}
	}

	var body CreateTrackBody
	err = c.BodyParser(&body)
	if err != nil {
		return types.ApiBadRequestError("Failed to create track: " + err.Error())
	}

	errs := api.validateBody(body)
	if errs != nil {
		return ApiError{
			Status:  400,
			Message: "Failed to create track",
			Data:    errs,
		}
	}

	id := utils.CreateId()
	log.Println("New ID:", id)

	err = os.MkdirAll(path.Join(api.workDir, "tracks"), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(api.workDir, "images"), 0755)
	if err != nil {
		return err
	}

	bestQualityFile, err := utils.GetSingleFile(form, "bestQualityFile")
	if err != nil {
		return ApiError{
			Status:  400,
			Message: "Failed to create track",
			Data: map[string]any{
				"bestQualityFile": err.Error(),
			},
		}
	}

	ext, err := getExtForBestQuality(bestQualityFile.Header.Get("Content-Type"))
	if err != nil {
		return ApiError{
			Status:  400,
			Message: "Failed to create track",
			Data: map[string]any{
				"bestQualityFile": err.Error(),
			},
		}
	}

	coverArtFile, err := utils.GetSingleFile(form, "coverArt")
	if err != nil {
		return ApiError{
			Status:  400,
			Message: "Failed to create track",
			Data: map[string]any{
				"coverArt": err.Error(),
			},
		}
	}

	fmt.Printf("Headers: %v\n", coverArtFile.Header)

	contentType := coverArtFile.Header.Get("Content-Type")

	var imageExt string
	switch contentType {
	case "image/jpeg":
		imageExt = "jpg"
	case "image/png":
		imageExt = "png"
	case "":
		return ApiError{
			Status:  400,
			Message: "Failed to create track",
			Data: map[string]any{
				"coverArt": errors.New("No 'Content-Type' is not set"),
			},
		}
	default:
		return ApiError{
			Status:  400,
			Message: "Failed to create track",
			Data: map[string]any{
				"coverArt": fmt.Errorf("Unsupported Content-Type: %v", contentType),
			},
		}
	}

	fmt.Printf("Image Ext: %v\n", imageExt)

	bestQualityFileName := fmt.Sprintf("%v.best.%v", id, ext)
	mobileQualityFileName := fmt.Sprintf("%v.mobile.mp3", id)
	coverArtFileName := fmt.Sprintf("%v.cover.%v", id, imageExt)

	track, err := queries.CreateTrack(c.Context(), database.CreateTrackParams{
		ID:                id,
		TrackNumber:       int32(body.Number),
		Name:              body.Name,
		CoverArt:          coverArtFileName,
		BestQualityFile:   bestQualityFileName,
		MobileQualityFile: mobileQualityFileName,
		AlbumID:           body.Album,
		ArtistID:          body.Artist,
	})

	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			switch err.Code {
			case "23503":
				switch err.ConstraintName {
				case "tracks_album_id_fk":
					return ApiError{
						Status:  400,
						Message: fmt.Sprintf("Failed to create track: No album with id: '%v'", body.Album),
						Data:    nil,
					}
				case "tracks_artist_id_fk":
					return ApiError{
						Status:  400,
						Message: fmt.Sprintf("Failed to create track: No artist with id: '%v'", body.Artist),
						Data:    nil,
					}
				}
			case "23505":
				if err.ConstraintName == "tracks_track_number_unique" {
					return ApiError{
						Status:  400,
						Message: "Failed to create track: Number need to be unique",
						Data:    nil,
					}
				}
			}
		}

		return err
	}

	track.CoverArt = ConvertURL(c, "/images/"+track.CoverArt)
	track.BestQualityFile = ConvertURL(c, "/tracks/"+track.BestQualityFile)
	track.MobileQualityFile = ConvertURL(c, "/tracks/"+track.MobileQualityFile)

	err = api.writeTrackFile(bestQualityFile, bestQualityFileName)
	if err != nil {
		return err
	}

	err = api.writeImageFile(coverArtFile, coverArtFileName)
	if err != nil {
		return err
	}

	fmt.Println("Ext:", ext)
	fmt.Println("Len:", len(form.File["mobileQualityFile"]))
	if ext == "mp3" && len(form.File["mobileQualityFile"]) == 0 {
		err := os.Symlink(bestQualityFileName, path.Join(api.workDir, "tracks", mobileQualityFileName))
		if err != nil {
			return err
		}
	} else {
		mobileQualityFile, err := utils.GetSingleFile(form, "mobileQualityFile")
		if err != nil {
			return ApiError{
				Status:  400,
				Message: "Failed to create track",
				Data: map[string]any{
					"mobileQualityFile": err.Error(),
				},
			}
		}

		err = api.writeTrackFile(mobileQualityFile, mobileQualityFileName)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(ApiData{
		Status: 200,
		Data:   track,
	})
}

func (apiConfig *ApiConfig) HandlerGetTrack(c *fiber.Ctx) error {
	id := c.Params("id")

	track, err := apiConfig.queries.GetTrack(c.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": fmt.Sprintf("No track with id: %s", id)})
		} else {
			return err
		}
	}

	track.CoverArt = ConvertURL(c, track.CoverArt)

	return c.JSON(track)
}

func (apiConfig *ApiConfig) HandlerCreateQueueFromAlbum(c *fiber.Ctx) error {
	body := struct {
		AlbumId string `json:"albumId"`
	}{}

	fmt.Println("Hello")

	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	tracks, err := apiConfig.queries.GetTracksByAlbum(c.Context(), body.AlbumId)
	if err != nil {
		return err
	}

	res := struct {
		Queue []database.GetTracksByAlbumRow `json:"queue"`
	}{
		Queue: tracks,
	}

	return c.JSON(res)
}
