package helper

import (
	"context"
	"database/sql"
	"os"

	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
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
	err := os.Mkdir(trackDir, 0755)
	if err != nil {
		return "", err
	}

	original, err := utils.ProcessOriginalVersion(data.InputFile, trackDir, "track.original")
	if err != nil {
		return "", err
	}

	mobile, err := utils.ProcessMobileVersion(data.InputFile, trackDir, "track.mobile")
	if err != nil {
		return "", err
	}

	trackInfo, err := utils.GetTrackInfo(data.InputFile)
	if err != nil {
		return "", err
	}

	// dateRegex := regexp.MustCompile(`^([12]\d\d\d)`)
	// if tag, exists := trackInfo.Tags["title"]; exists {
	// 	name = tag
	// }
	//
	// var year sql.NullInt64
	// if tag, exists := trackInfo.Tags["date"]; exists {
	// 	match := dateRegex.FindStringSubmatch(tag)
	// 	if len(match) > 0 {
	// 		y, _ := strconv.Atoi(match[1])
	//
	// 		year.Int64 = int64(y)
	// 		year.Valid = true
	// 	}
	// }
	//
	// var number int
	// if tag, exists := trackInfo.Tags["track"]; exists {
	// 	y, _ := strconv.Atoi(tag)
	// 	number = y
	// }
	//
	// if number == 0 || data.ForceExtractNumber {
	// 	number = utils.ExtractNumber(originalName)
	// }

	trackId, err = db.CreateTrack(ctx, database.CreateTrackParams{
		Id:       trackId,
		Name:     data.Name,
		OtherName: sql.NullString{
			String: data.OtherName,
			Valid:  data.OtherName != "",
		},
		AlbumId:  data.AlbumId,
		ArtistId: data.ArtistId,
		Duration: int64(trackInfo.Duration),
		Number: sql.NullInt64{
			Int64: data.Number,
			Valid: data.Number != 0,
		},
		Year: sql.NullInt64{
			Int64: data.Year,
			Valid: data.Year != 0,
		},
		OriginalFilename: original,
		MobileFilename:   mobile,
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
