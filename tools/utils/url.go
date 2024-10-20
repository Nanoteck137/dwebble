package utils

import (
	"database/sql"
	"fmt"

	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

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

func ConvertArtistPictureURL(c pyrin.Context, val sql.NullString) string {
	return ConvertURL(c, "/images/"+DefaultArtistPictureName)
}

func ConvertAlbumCoverURL(c pyrin.Context, albumId string, val sql.NullString) types.CoverArt {
	if val.Valid && val.String != "" {
		coverArt := val.String
		return types.CoverArt{
			Original: ConvertURL(c, "/files/albums/images/"+albumId+"/"+coverArt),
			Small:    ConvertURL(c, "/files/albums/images/"+albumId+"/"+"cover-128.png"),
			Medium:   ConvertURL(c, "/files/albums/images/"+albumId+"/"+"cover-256.png"),
			Large:    ConvertURL(c, "/files/albums/images/"+albumId+"/"+"cover-512.png"),
		}
	}

	url := ConvertURL(c, "/images/"+DefaultAlbumCoverArtName)
	return types.CoverArt{
		Original: url,
		Small:    url,
		Medium:   url,
		Large:    url,
	}
}
