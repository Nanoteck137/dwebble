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

	// return ConvertURL(c, "/files/images/"+DefaultArtistPictureName)
}

func ConvertAlbumCoverURL(c pyrin.Context, albumId string, val sql.NullString) types.Images {
	if val.Valid && val.String != "" {
		coverArt := val.String
		return types.Images{
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
