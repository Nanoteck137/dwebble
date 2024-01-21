package handlers

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/dwebble/utils"
)

// HandleGetTracks godoc
//
//	@Summary		Get all tracks
//	@Description	Get all tracks
//	@Tags			tracks
//	@Produce		json
//	@Success		200	{object}	types.ApiResponse[types.ApiGetTracksData]
//	@Failure		400	{object}	types.ApiError
//	@Failure		500	{object}	types.ApiError
//	@Router			/tracks [get]
func (api *ApiConfig) HandleGetTracks(c *fiber.Ctx) error {
	tracks, err := api.queries.GetAllTracks(c.UserContext())
	if err != nil {
		return err
	}

	result := types.ApiGetTracksData{
		Tracks: make([]types.ApiGetTracksDataTrackItem, len(tracks)),
	}

	for i, track := range tracks {
		result.Tracks[i] = types.ApiGetTracksDataTrackItem{
			ApiTrack: types.ApiTrack{
				Id:                track.ID,
				Number:            track.TrackNumber,
				Name:              track.Name,
				CoverArt:          ConvertURL(c, "/images/"+track.CoverArt),
				BestQualityFile:   ConvertURL(c, "/tracks/"+track.BestQualityFile),
				MobileQualityFile: ConvertURL(c, "/tracks/"+track.MobileQualityFile),
				AlbumId:           track.AlbumID,
				ArtistId:          track.ArtistID,
			},
			AlbumName:  track.AlbumName,
			ArtistName: track.ArtistName,
		}
	}

	return c.JSON(types.NewApiResponse(result))
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
		return types.ApiBadRequestError("Failed to create track", map[string]any{
			"mobileQualityFile": "No Content-Type set",
		})
	} else if contentType != "audio/mpeg" {
		return types.ApiBadRequestError("Failed to create track", map[string]any{
			"mobileQualityFile": "Content-Type needs to be audio/mpeg",
		})
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

// HandlePostTrack godoc
//
//	@Summary		Create new track
//	@Description	Create new track
//	@Tags			tracks
//	@Accept			mpfd
//	@Produce		json
//	@Param			name				formData	string	true	"Track name"
//	@Param			number				formData	string	true	"Track number"
//	@Param			albumId				formData	string	true	"Album Id"
//	@Param			artistId			formData	string	true	"Artist Id"
//	@Param			bestQualityFile		formData	file	true	"Best Quality File"
//	@Param			mobileQualityFile	formData	file	true	"Mobile Quality File"
//	@Param			coverArt			formData	file	true	"Cover Art"
//	@Success		200					{object}	types.ApiResponse[types.ApiPostTrackData]
//	@Failure		400					{object}	types.ApiError
//	@Failure		500					{object}	types.ApiError
//	@Router			/tracks [post]
func (api *ApiConfig) HandlePostTrack(c *fiber.Ctx) error {
	tx, err := api.db.BeginTx(c.UserContext(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(c.UserContext())

	queries := database.New(tx)

	form, err := c.MultipartForm()
	if err != nil {
		return types.ApiBadRequestError(fmt.Sprintf("Failed to create track: %v", err))
	}

	var body CreateTrackBody
	err = c.BodyParser(&body)
	if err != nil {
		return types.ApiBadRequestError("Failed to create track: " + err.Error())
	}

	errs := api.validateBody(body)
	if errs != nil {
		return types.ApiBadRequestError("Failed to create track", errs)
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
		return types.ApiBadRequestError("Failed to create track", types.Map{
			"bestQualityFile": err.Error(),
		})
	}

	ext, err := getExtForBestQuality(bestQualityFile.Header.Get("Content-Type"))
	if err != nil {
		return types.ApiBadRequestError("Failed to create track", types.Map{
			"bestQualityFile": err.Error(),
		})
	}

	coverArtFile, err := utils.GetSingleFile(form, "coverArt")
	if err != nil {
		return types.ApiBadRequestError("Failed to create track", types.Map{
			"coverArt": err.Error(),
		})
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
		return types.ApiBadRequestError("Failed to create track", types.Map{
			"coverArt": "No 'Content-Type' is not set",
		})
	default:
		return types.ApiBadRequestError("Failed to create track", types.Map{
			"coverArt": "Unsupported Content-Type: " + contentType,
		})
	}

	fmt.Printf("Image Ext: %v\n", imageExt)

	bestQualityFileName := fmt.Sprintf("%v.best.%v", id, ext)
	mobileQualityFileName := fmt.Sprintf("%v.mobile.mp3", id)
	coverArtFileName := fmt.Sprintf("%v.cover.%v", id, imageExt)

	track, err := queries.CreateTrack(c.UserContext(), database.CreateTrackParams{
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
					return types.ApiBadRequestError(fmt.Sprintf("Failed to create track: No album with id: '%v'", body.Album))
				case "tracks_artist_id_fk":
					return types.ApiBadRequestError(fmt.Sprintf("Failed to create track: No artist with id: '%v'", body.Artist))
				}
			case "23505":
				if err.ConstraintName == "tracks_track_number_unique" {
					return types.ApiBadRequestError("Failed to create track: Number need to be unique")
				}
			}
		}

		return err
	}

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
			return types.ApiBadRequestError("Failed to create track", types.Map{
				"mobileQualityFile": err.Error(),
			})
		}

		err = api.writeTrackFile(mobileQualityFile, mobileQualityFileName)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(c.UserContext())
	if err != nil {
		return err
	}

	return c.JSON(types.NewApiResponse(types.ApiPostTrackData{
		Id:                track.ID,
		Number:            track.TrackNumber,
		Name:              track.Name,
		CoverArt:          ConvertURL(c, "/images/"+track.CoverArt),
		BestQualityFile:   ConvertURL(c, "/tracks/"+track.BestQualityFile),
		MobileQualityFile: ConvertURL(c, "/tracks/"+track.MobileQualityFile),
		AlbumId:           track.AlbumID,
		ArtistId:          track.ArtistID,
	}))
}

// HandleGetTrackById godoc
//
//	@Summary		Get track by id
//	@Description	Get track by id
//	@Tags			tracks
//	@Produce		json
//	@Param			id	path		string	true	"Track Id"
//	@Success		200	{object}	types.ApiResponse[types.ApiGetTrackByIdData]
//	@Failure		400	{object}	types.ApiError
//	@Failure		404	{object}	types.ApiError
//	@Failure		500	{object}	types.ApiError
//	@Router			/tracks/{id} [get]
func (api *ApiConfig) HandleGetTrackById(c *fiber.Ctx) error {
	id := c.Params("id")

	track, err := api.queries.GetTrack(c.UserContext(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return types.ApiNotFoundError(fmt.Sprintf("No track with id: '%s'", id))
		} else {
			return err
		}
	}

	return c.JSON(types.NewApiResponse(types.ApiGetTrackByIdData{
		Id:                track.ID,
		Number:            track.TrackNumber,
		Name:              track.Name,
		CoverArt:          ConvertURL(c, "/images/"+track.CoverArt),
		BestQualityFile:   ConvertURL(c, "/tracks/"+track.BestQualityFile),
		MobileQualityFile: ConvertURL(c, "/tracks/"+track.MobileQualityFile),
		AlbumId:           track.AlbumID,
		ArtistId:          track.ArtistID,
	}))
}

func InstallTrackHandlers(router *fiber.App, apiConfig *ApiConfig) {
	router.Get("/tracks", apiConfig.HandleGetTracks)
	router.Post("/tracks", apiConfig.HandlePostTrack)
	router.Get("/tracks/:id", apiConfig.HandleGetTrackById)
}
