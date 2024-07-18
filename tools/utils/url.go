package utils

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
)

const (
	DefaultArtistPictureName = "default/default_artist.png"
	DefaultAlbumCoverArtName = "default/default_album.png"
	DefaultTrackCoverArtName = "default/default_album.png"
)

func ConvertURL(c echo.Context, path string) string {
	host := c.Request().Host

	scheme := "http"

	h := c.Request().Header.Get("X-Forwarded-Proto")
	if h != "" {
		scheme = h
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

func ConvertImageURL(c echo.Context, val sql.NullString, def string) string {
	coverArt := def
	if val.Valid && val.String != "" {
		coverArt = val.String
	}

	return ConvertURL(c, "/images/"+coverArt)
}

func ConvertArtistPictureURL(c echo.Context, val sql.NullString) string {
	return ConvertImageURL(c, val, DefaultArtistPictureName)
}

func ConvertAlbumCoverURL(c echo.Context, val sql.NullString) string {
	return ConvertImageURL(c, val, DefaultAlbumCoverArtName)
}

func ConvertTrackCoverURL(c echo.Context, val sql.NullString) string {
	return ConvertImageURL(c, val, DefaultTrackCoverArtName)
}
