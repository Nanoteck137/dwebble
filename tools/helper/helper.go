package helper

import (
	"context"
	"database/sql"
	"errors"
	"os"

	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type ImportTrackData struct {
	InputFile string

	Name      string
	OtherName string

	AlbumId  string
	ArtistId string

	Number int64
	Year   int64
}

func ImportTrack(ctx context.Context, db *database.Database, workDir types.WorkDir, data ImportTrackData) (string, error) {
	trackId := utils.CreateTrackId()

	trackDir := workDir.Track(trackId)

	dirs := []string{
		trackDir.String(),
		trackDir.Media(),
	}

	for _, dir := range dirs {
		err := os.Mkdir(dir, 0755)
		if err != nil && !os.IsExist(err) {
			return "", err
		}
	}

	// original, err := utils.ProcessOriginalVersion(data.InputFile, trackDir, "track.original")
	// if err != nil {
	// 	return "", err
	// }

	// mobile, err := utils.ProcessMobileVersion(data.InputFile, trackDir, "track.mobile")
	// if err != nil {
	// 	return "", err
	// }

	var mediaType types.MediaType

	res, err := ffprobe.ProbeURL(ctx, data.InputFile)
	if err != nil {
		return "", err
	}
	_ = res

	// TODO(patrik): Check for nil
	stream := res.FirstAudioStream()

	// TODO(patrik): Move to helper
	switch res.Format.FormatName {
	case "flac":
		mediaType = types.MediaTypeFlac
	case "ogg":
		switch stream.CodecName {
		case "opus":
			mediaType = types.MediaTypeOggOpus
		case "vorbis":
			mediaType = types.MediaTypeOggVorbis
		}
	}

	if mediaType == "" {
		return "", errors.New("unknown media type")
	}

	log.Info("Found media type", "type", mediaType)

	// // TODO(patrik): Replace this with ffprobe.v2 package
	// trackInfo, err := utils.GetTrackInfo(data.InputFile)
	// if err != nil {
	// 	return "", err
	// }

	// TODO(patrik): Get from the probe
	var duration int64 = 0

	trackId, err = db.CreateTrack(ctx, database.CreateTrackParams{
		Id:   trackId,
		Name: data.Name,
		OtherName: sql.NullString{
			String: data.OtherName,
			Valid:  data.OtherName != "",
		},
		AlbumId:  data.AlbumId,
		ArtistId: data.ArtistId,
		Duration: duration,
		Number: sql.NullInt64{
			Int64: data.Number,
			Valid: data.Number != 0,
		},
		Year: sql.NullInt64{
			Int64: data.Year,
			Valid: data.Year != 0,
		},
	})
	if err != nil {
		return "", err
	}

	mediaId := utils.CreateTrackMediaId()
	mediaDir := trackDir.MediaItem(mediaId)

	{
		dirs := []string{
			mediaDir,
		}

		for _, dir := range dirs {
			err := os.Mkdir(dir, 0755)
			if err != nil && !os.IsExist(err) {
				return "", err
			}
		}
	}

	original, err := utils.ProcessOriginalVersion(data.InputFile, mediaDir, "track")
	if err != nil {
		return "", err
	}

	err = db.CreateTrackMedia(ctx, database.CreateTrackMediaParams{
		Id:         mediaId,
		TrackId:    trackId,
		Filename:   original,
		MediaType:  mediaType,
		IsOriginal: true,
	})
	if err != nil {
		return "", err
	}

	return trackId, nil
}

func CreateAlbum(ctx context.Context, db *database.Database, workDir types.WorkDir, params database.CreateAlbumParams) (database.Album, error) {
	album, err := db.CreateAlbum(ctx, params)
	if err != nil {
		return database.Album{}, err
	}

	albumDir := workDir.Album(album.Id)

	err = os.Mkdir(albumDir, 0755)
	if err != nil {
		return database.Album{}, err
	}

	return album, nil
}

func CreateArtist(ctx context.Context, db *database.Database, workDir types.WorkDir, params database.CreateArtistParams) (database.Artist, error) {
	artist, err := db.CreateArtist(ctx, params)
	if err != nil {
		return database.Artist{}, err
	}

	artistDir := workDir.Artist(artist.Id)

	err = os.Mkdir(artistDir, 0755)
	if err != nil {
		return database.Artist{}, err
	}

	return artist, nil
}
