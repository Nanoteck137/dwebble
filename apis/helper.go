package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/helper"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

const defaultMemory = 32 << 20 // 32 MB

func User(app core.App, c pyrin.Context) (*database.User, error) {
	apiTokenHeader := c.Request().Header.Get("X-Api-Token")
	if apiTokenHeader != "" {
		ctx := context.TODO()
		token, err := app.DB().GetApiTokenById(ctx, apiTokenHeader)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				return nil, InvalidAuth("invalid api token")
			}

			return nil, err
		}

		user, err := app.DB().GetUserById(c.Request().Context(), token.UserId)
		if err != nil {
			// TODO(patrik): Should we handle error here?
			return nil, err
		}

		return &user, nil
	}

	authHeader := c.Request().Header.Get("Authorization")
	tokenString := utils.ParseAuthHeader(authHeader)
	if tokenString == "" {
		return nil, InvalidAuth("invalid authorization header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(app.Config().JwtSecret), nil
	})

	if err != nil {
		// TODO(patrik): Handle error better
		return nil, InvalidAuth("invalid authorization token")
	}

	jwtValidator := jwt.NewValidator(jwt.WithIssuedAt())

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := jwtValidator.Validate(token.Claims); err != nil {
			return nil, InvalidAuth("invalid authorization token")
		}

		userId := claims["userId"].(string)
		user, err := app.DB().GetUserById(c.Request().Context(), userId)
		if err != nil {
			return nil, InvalidAuth("invalid authorization token")
		}

		return &user, nil
	}

	return nil, InvalidAuth("invalid authorization token")
}

func ConvertSqlNullString(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}

	return nil
}

func ConvertSqlNullInt64(value sql.NullInt64) *int64 {
	if value.Valid {
		return &value.Int64
	}

	return nil
}

const (
	UNKNOWN_ARTIST_ID   = "unknown"
	UNKNOWN_ARTIST_NAME = "UNKNOWN"
)

func EnsureUnknownArtistExists(ctx context.Context, db *database.Database, workDir types.WorkDir) error {
	_, err := db.GetArtistById(ctx, UNKNOWN_ARTIST_ID)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			log.Info("Creating 'unknown' artist")
			_, err = helper.CreateArtist(ctx, db, workDir, database.CreateArtistParams{
				Id:   UNKNOWN_ARTIST_ID,
				Name: UNKNOWN_ARTIST_NAME,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

const (
	DefaultArtistPictureName = "default/default_artist.png"
	DefaultAlbumCoverArtName = "default/default_album.png"
	DefaultTrackCoverArtName = "default/default_album.png"
)

func ConvertURL(c pyrin.Context, path string) string {
	host := c.Request().Host

	scheme := "http"

	h := c.Request().Header.Get("X-Forwarded-Proto")
	if h != "" {
		scheme = h
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

func ConvertImageURL(c pyrin.Context, albumId string, val sql.NullString, def string) string {
	coverArt := def
	if val.Valid && val.String != "" {
		coverArt = val.String
	}

	return ConvertURL(c, "/files/albums/images/"+albumId+"/"+coverArt)
}

func ConvertArtistPicture(c pyrin.Context, artistId string, val sql.NullString) types.Images {
	if val.Valid && val.String != "" {
		first := "/files/artists/" + artistId + "/"
		return types.Images{
			Original: ConvertURL(c, first+val.String),
			Small:    ConvertURL(c, first+"picture-128.png"),
			Medium:   ConvertURL(c, first+"picture-256.png"),
			Large:    ConvertURL(c, first+"picture-512.png"),
		}
	}

	url := ConvertURL(c, "/files/images/"+DefaultArtistPictureName)
	return types.Images{
		Original: url,
		Small:    url,
		Medium:   url,
		Large:    url,
	}
}

func ConvertAlbumCoverURL(c pyrin.Context, albumId string, val sql.NullString) types.Images {
	if val.Valid && val.String != "" {
		coverArt := val.String
		return types.Images{
			// TODO(patrik): Move to const (cover-128.png, cover-256.png, cover-512.png)
			Original: ConvertURL(c, "/files/albums/images/"+albumId+"/"+coverArt),
			Small:    ConvertURL(c, "/files/albums/images/"+albumId+"/"+"cover-128.png"),
			Medium:   ConvertURL(c, "/files/albums/images/"+albumId+"/"+"cover-256.png"),
			Large:    ConvertURL(c, "/files/albums/images/"+albumId+"/"+"cover-512.png"),
		}
	}

	url := ConvertURL(c, "/files/images/"+DefaultAlbumCoverArtName)
	return types.Images{
		Original: url,
		Small:    url,
		Medium:   url,
		Large:    url,
	}
}

func DeleteArtist(ctx context.Context, db *database.Database, workDir types.WorkDir, artist database.Artist) error {
	dir := workDir.Artist(artist.Id)
	targetName := fmt.Sprintf("artist-%s-%d", artist.Id, time.Now().UnixMilli())
	target := path.Join(workDir.Trash(), targetName)

	err := os.Rename(dir, target)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = db.DeleteArtistFromSearch(ctx, artist)
	if err != nil {
		return err
	}

	err = db.DeleteArtist(ctx, artist.Id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAlbum(ctx context.Context, db *database.Database, workDir types.WorkDir, album database.Album) error {
	dir := workDir.Album(album.Id)
	targetName := fmt.Sprintf("album-%s-%d", album.Id, time.Now().UnixMilli())
	target := path.Join(workDir.Trash(), targetName)

	err := os.Rename(dir, target)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = db.DeleteAlbumFromSearch(ctx, album)
	if err != nil {
		return err
	}

	err = db.DeleteAlbum(ctx, album.Id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTrack(ctx context.Context, db *database.Database, workDir types.WorkDir, track database.Track) error {
	dir := workDir.Track(track.Id)
	targetName := fmt.Sprintf("track-%s-%d", track.Id, time.Now().UnixMilli())
	target := path.Join(workDir.Trash(), targetName)

	err := os.Rename(dir.String(), target)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = db.DeleteTrackFromSearch(ctx, track)
	if err != nil {
		return err
	}

	err = db.DeleteTrack(ctx, track.Id)
	if err != nil {
		return err
	}

	return nil
}
