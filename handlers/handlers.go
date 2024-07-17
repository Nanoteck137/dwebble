package handlers

import (
	"database/sql"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
)

type Handler struct {
	Name        string
	Method      string
	Path        string
	DataType    any
	BodyType    types.Body
	HandlerFunc echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

type Group interface {
	Register(handlers ...Handler)
}

const (
	DefaultArtistPictureName = "default/default_artist.png"
	DefaultAlbumCoverArtName = "default/default_album.png"
	DefaultTrackCoverArtName = "default/default_album.png"
)

type Handlers struct {
	libraryDir string
	workDir    types.WorkDir
	db         *database.Database

	jwtValidator *jwt.Validator
}

func New(db *database.Database, libraryDir string, workDir types.WorkDir) *Handlers {
	return &Handlers{
		libraryDir:   libraryDir,
		workDir:      workDir,
		db:           db,
		jwtValidator: jwt.NewValidator(jwt.WithIssuedAt()),
	}
}

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
