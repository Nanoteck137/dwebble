package apis

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/helper"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
)

type Name struct {
	Default string  `json:"default"`
	Other   *string `json:"other"`
}

type TrackDetails struct {
	Id   string `json:"id"`
	Name Name   `json:"name"`

	Duration int64  `json:"duration"`
	Number   *int64 `json:"number"`
	Year     *int64 `json:"year"`

	CoverArt types.Images `json:"coverArt"`

	AlbumId   string `json:"albumId"`
	AlbumName Name   `json:"albumName"`

	ArtistId   string `json:"artistId"`
	ArtistName Name   `json:"artistName"`

	Tags             []string     `json:"tags"`
	FeaturingArtists []ArtistInfo `json:"featuringArtists"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

type Track struct {
	Id   string `json:"id"`
	Name Name   `json:"name"`

	Duration int64  `json:"duration"`
	Number   *int64 `json:"number"`
	Year     *int64 `json:"year"`

	CoverArt types.Images `json:"coverArt"`

	AlbumId   string `json:"albumId"`
	AlbumName Name   `json:"albumName"`

	Artists []ArtistInfo `json:"artists"`

	Tags []string `json:"tags"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

func ConvertDBTrack(c pyrin.Context, track database.Track) Track {
	artists := make([]ArtistInfo, len(track.FeaturingArtists)+1)

	artists[0] = ArtistInfo{
		Id: track.ArtistId,
		Name: Name{
			Default: track.ArtistName,
			Other:   ConvertSqlNullString(track.ArtistOtherName),
		},
	}

	featuringArtists := ConvertDBFeaturingArtists(track.FeaturingArtists)
	for i, v := range featuringArtists {
		artists[i+1] = v
	}

	return Track{
		Id: track.Id,
		Name: Name{
			Default: track.Name,
			Other:   ConvertSqlNullString(track.OtherName),
		},
		Duration: track.Duration,
		Number:   ConvertSqlNullInt64(track.Number),
		Year:     ConvertSqlNullInt64(track.Year),
		CoverArt: ConvertAlbumCoverURL(c, track.AlbumId, track.AlbumCoverArt),
		AlbumId:  track.AlbumId,
		AlbumName: Name{
			Default: track.AlbumName,
			Other:   ConvertSqlNullString(track.AlbumOtherName),
		},
		Artists: artists,
		Tags:    utils.SplitString(track.Tags.String),
		Created: track.Created,
		Updated: track.Updated,
	}
}

func ConvertDBTrackToDetails(c pyrin.Context, track database.Track) TrackDetails {
	artists := make([]ArtistInfo, len(track.FeaturingArtists))

	featuringArtists := ConvertDBFeaturingArtists(track.FeaturingArtists)
	for i, v := range featuringArtists {
		artists[i] = v
	}

	return TrackDetails{
		Id: track.Id,
		Name: Name{
			Default: track.Name,
			Other:   ConvertSqlNullString(track.OtherName),
		},
		Duration: track.Duration,
		Number:   ConvertSqlNullInt64(track.Number),
		Year:     ConvertSqlNullInt64(track.Year),
		CoverArt: ConvertAlbumCoverURL(c, track.AlbumId, track.AlbumCoverArt),
		AlbumId:  track.AlbumId,
		AlbumName: Name{
			Default: track.AlbumName,
			Other:   ConvertSqlNullString(track.AlbumOtherName),
		},
		ArtistId: track.ArtistId,
		ArtistName: Name{
			Default: track.ArtistName,
			Other:   ConvertSqlNullString(track.ArtistOtherName),
		},
		Tags:             utils.SplitString(track.Tags.String),
		FeaturingArtists: featuringArtists,
		Created:          track.Created,
		Updated:          track.Updated,
	}
}

type GetTracks struct {
	Page   types.Page `json:"page"`
	Tracks []Track    `json:"tracks"`
}

type GetDetailedTracks struct {
	Page   types.Page     `json:"page"`
	Tracks []TrackDetails `json:"tracks"`
}

type GetTrackById struct {
	Track
}

type GetDetailedTrackById struct {
	TrackDetails
}

type EditTrackBody struct {
	Name             *string   `json:"name,omitempty"`
	OtherName        *string   `json:"otherName,omitempty"`
	ArtistId         *string   `json:"artistId,omitempty"`
	Year             *int64    `json:"year,omitempty"`
	Number           *int64    `json:"number,omitempty"`
	Tags             *[]string `json:"tags,omitempty"`
	FeaturingArtists *[]string `json:"featuringArtists,omitempty"`
}

func DiscardEntriesStringArray(arr []string) []string {
	if arr == nil {
		return nil
	}

	var res []string
	for _, s := range arr {
		if s != "" {
			res = append(res, s)
		}
	}

	return res
}

// TODO(patrik): Should this be pyrins default
func DiscardEntriesStringArrayPtr(arr *[]string) *[]string {
	if arr == nil {
		return nil
	}

	var res []string

	v := *arr
	for _, s := range v {
		if s != "" {
			res = append(res, s)
		}
	}

	return &res
}

func (b *EditTrackBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
	b.OtherName = transform.StringPtr(b.OtherName)
	b.Tags = transform.StringArrayPtr(b.Tags)
	b.Tags = DiscardEntriesStringArrayPtr(b.Tags)
	b.FeaturingArtists = DiscardEntriesStringArrayPtr(b.FeaturingArtists)
}

func (b EditTrackBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name,
			validate.Required.When(b.Name != nil),
		),
		validate.Field(&b.ArtistId,
			validate.Required.When(b.ArtistId != nil),
		),
		// TODO(patrik): Add this to all numbers (album, artists)
		validate.Field(&b.Year, validate.Min(0)),
		validate.Field(&b.Number, validate.Min(0)),
	)
}

type UploadTrack struct {
	Id string `json:"id"`
}

type UploadTrackBody struct {
	Name      string `json:"name"`
	OtherName string `json:"otherName"`

	Number int64 `json:"number"`
	Year   int64 `json:"year"`

	AlbumId  string `json:"albumId"`
	ArtistId string `json:"artistId"`

	Tags             []string `json:"tags"`
	FeaturingArtists []string `json:"featuringArtists"`
}

// TODO(patrik): Move to pyrin
func StringArray(arr []string) []string {
	if arr == nil {
		return nil
	}

	for i, s := range arr {
		arr[i] = transform.String(s)
	}

	return arr
}

func (b *UploadTrackBody) Transform() {
	b.Name = transform.String(b.Name)
	b.OtherName = transform.String(b.OtherName)

	b.Tags = StringArray(b.Tags)
	b.Tags = DiscardEntriesStringArray(b.Tags)
	b.FeaturingArtists = StringArray(b.FeaturingArtists)
	b.FeaturingArtists = DiscardEntriesStringArray(b.FeaturingArtists)
}

func (b UploadTrackBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
		// validate.Field(&b.OtherName, validate.Required.When(b.OtherName != nil)),
	)
}

type GetTrackDetails struct {
	TrackDetails
}

// TODO(patrik): Move
func getPageOptions(q url.Values) database.FetchOptions {
	perPage := 100
	page := 0

	if s := q.Get("perPage"); s != "" {
		i, _ := strconv.Atoi(s)
		if i > 0 {
			perPage = i
		}
	}

	if s := q.Get("page"); s != "" {
		i, _ := strconv.Atoi(s)
		page = i
	}

	return database.FetchOptions{
		Filter:  q.Get("filter"),
		Sort:    q.Get("sort"),
		PerPage: perPage,
		Page:    page,
	}
}

func InstallTrackHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetTracks",
			Method:       http.MethodGet,
			Path:         "/tracks",
			ResponseType: GetTracks{},
			Errors:       []pyrin.ErrorType{ErrTypeInvalidFilter, ErrTypeInvalidSort},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				opts := getPageOptions(q)

				tracks, p, err := app.DB().GetPagedTracks(c.Request().Context(), opts)
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					if errors.Is(err, database.ErrInvalidSort) {
						return nil, InvalidSort(err)
					}

					return nil, err
				}

				res := GetTracks{
					Page:   p,
					Tracks: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetDetailedTracks",
			Method:       http.MethodGet,
			Path:         "/tracks/detailed",
			ResponseType: GetDetailedTracks{},
			Errors:       []pyrin.ErrorType{ErrTypeInvalidFilter, ErrTypeInvalidSort},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				opts := getPageOptions(q)
				tracks, p, err := app.DB().GetPagedTracks(c.Request().Context(), opts)
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					if errors.Is(err, database.ErrInvalidSort) {
						return nil, InvalidSort(err)
					}

					return nil, err
				}

				res := GetDetailedTracks{
					Page:   p,
					Tracks: make([]TrackDetails, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrackToDetails(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SearchTracks",
			Method:       http.MethodGet,
			Path:         "/tracks/search",
			ResponseType: GetTracks{},
			Errors:       []pyrin.ErrorType{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				query := strings.TrimSpace(q.Get("query"))

				tracks, err := app.DB().SearchTracks(query)
				if err != nil {
					return nil, err
				}

				res := GetTracks{
					Tracks: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Tracks[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetTrackById",
			Method:       http.MethodGet,
			Path:         "/tracks/:id",
			ResponseType: GetTrackById{},
			Errors:       []pyrin.ErrorType{ErrTypeTrackNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				track, err := app.DB().GetTrackById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TrackNotFound()
					}

					return nil, err
				}

				return GetTrackById{
					Track: ConvertDBTrack(c, track),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetTrackDetails",
			Method:       http.MethodGet,
			Path:         "/tracks/:id/details",
			ResponseType: GetTrackDetails{},
			Errors:       []pyrin.ErrorType{ErrTypeTrackNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				track, err := app.DB().GetTrackById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TrackNotFound()
					}

					return nil, err
				}

				return GetTrackDetails{
					TrackDetails: ConvertDBTrackToDetails(c, track),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "EditTrack",
			Method:   http.MethodPatch,
			Path:     "/tracks/:id",
			BodyType: EditTrackBody{},
			Errors:   []pyrin.ErrorType{ErrTypeTrackNotFound, ErrTypeArtistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[EditTrackBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				track, err := app.DB().GetTrackById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TrackNotFound()
					}

					return nil, err
				}

				changes := database.TrackChanges{}

				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != track.Name,
					}
				}

				if body.OtherName != nil {
					changes.OtherName = types.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.OtherName,
							Valid:  *body.OtherName != "",
						},
						Changed: *body.OtherName != track.OtherName.String,
					}
				}

				if body.ArtistId != nil {
					// TODO(patrik): Should we silently continue if
					// we don't find the artist
					_, err := app.DB().GetArtistById(ctx, *body.ArtistId)
					if err != nil {
						if errors.Is(err, database.ErrItemNotFound) {
							return nil, ArtistNotFound()
						}

						return nil, err
					}

					changes.ArtistId = types.Change[string]{
						Value:   *body.ArtistId,
						Changed: *body.ArtistId != track.ArtistId,
					}
				}

				if body.Year != nil {
					changes.Year = types.Change[sql.NullInt64]{
						Value: sql.NullInt64{
							Int64: *body.Year,
							Valid: *body.Year != 0,
						},
						Changed: *body.Year != track.Year.Int64,
					}
				}

				if body.Number != nil {
					changes.Number = types.Change[sql.NullInt64]{
						Value: sql.NullInt64{
							Int64: *body.Number,
							Valid: *body.Number != 0,
						},
						Changed: *body.Number != track.Number.Int64,
					}
				}

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				err = db.UpdateTrack(ctx, track.Id, changes)
				if err != nil {
					return nil, err
				}

				if body.Tags != nil {
					tags := *body.Tags

					err = db.RemoveAllTagsFromTrack(ctx, track.Id)
					if err != nil {
						return nil, err
					}

					for _, tag := range tags {
						slug := utils.Slug(tag)

						err := db.CreateTag(ctx, slug)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}

						err = db.AddTagToTrack(ctx, slug, track.Id)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}
				}

				if body.FeaturingArtists != nil {
					featuringArtists := *body.FeaturingArtists

					err := db.RemoveAllTrackFeaturingArtists(ctx, track.Id)
					if err != nil {
						return nil, err
					}

					for _, artistId := range featuringArtists {
						artist, err := db.GetArtistById(ctx, artistId)
						if err != nil {
							if errors.Is(err, database.ErrItemNotFound) {
								// TODO(patrik): Should we just be
								// silently continuing
								// NOTE(patrik): If we don't find the artist
								// we just silently continue
								continue
							}

							return nil, err
						}

						err = db.AddFeaturingArtistToTrack(ctx, track.Id, artist.Id)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}
				}

				err = tx.Commit()
				if err != nil {
					return nil, err
				}

				{
					track, err := app.DB().GetTrackById(ctx, track.Id)
					if err != nil {
						return nil, err
					}

					err = app.DB().UpdateSearchTrack(ctx, track)
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeleteTrack",
			Method: http.MethodDelete,
			Path:   "/tracks/:id",
			Errors: []pyrin.ErrorType{ErrTypeTrackNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				track, err := app.DB().GetTrackById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TrackNotFound()
					}

					return nil, err
				}

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				err = DeleteTrack(ctx, db, app.WorkDir(), track)
				if err != nil {
					return nil, err
				}

				err = tx.Commit()
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.FormApiHandler{
			Name:   "UploadTrack",
			Method: http.MethodPost,
			Path:   "/tracks",
			Spec: pyrin.FormSpec{
				BodyType: UploadTrackBody{},
				Files: map[string]pyrin.FormFileSpec{
					"track": {
						NumExpected: 1,
					},
				},
			},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[UploadTrackBody](c)
				if err != nil {
					pretty.Println(err)
					return nil, err
				}

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				// db := app.DB()

				ctx := context.TODO()

				album, err := db.GetAlbumById(ctx, body.AlbumId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, AlbumNotFound()
					}

					return nil, err
				}

				artist, err := db.GetArtistById(ctx, body.ArtistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ArtistNotFound()
					}

					return nil, err
				}

				copyFormFileToTemp := func(file *multipart.FileHeader) (string, error) {
					ext := path.Ext(file.Filename)

					src, err := file.Open()
					if err != nil {
						return "", err
					}
					defer src.Close()

					// TODO(patrik): Copy the file to $trackDir/raw.flac instead
					dst, err := os.CreateTemp("", "track.*"+ext)
					if err != nil {
						return "", err
					}
					defer dst.Close()

					_, err = io.Copy(dst, src)
					if err != nil {
						return "", err
					}

					return dst.Name(), nil
				}

				files, err := pyrin.FormFiles(c, "track")
				if err != nil {
					return nil, err
				}

				file := files[0]

				// ext := path.Ext(file.Filename)
				// name := strings.TrimSuffix(file.Filename, ext)

				filename, err := copyFormFileToTemp(file)
				if err != nil {
					return nil, err
				}
				defer os.Remove(filename)

				data := helper.ImportTrackData{
					InputFile: filename,
					Name:      body.Name,
					OtherName: body.OtherName,
					AlbumId:   album.Id,
					ArtistId:  artist.Id,
					Number:    body.Number,
					Year:      body.Year,
				}
				trackId, err := helper.ImportTrack(ctx, db, app.WorkDir(), data)
				if err != nil {
					return nil, err
				}

				for _, tag := range body.Tags {
					// TODO(patrik): Move slugify to transform function
					slug := utils.Slug(tag)

					err := db.CreateTag(ctx, slug)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}

					err = db.AddTagToTrack(ctx, slug, trackId)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				for _, artistId := range body.FeaturingArtists {
					artist, err := db.GetArtistById(ctx, artistId)
					if err != nil {
						if errors.Is(err, database.ErrItemNotFound) {
							// TODO(patrik): Should we just be
							// silently continuing
							// NOTE(patrik): If we don't find the artist
							// we just silently continue
							continue
						}

						return nil, err
					}

					err = db.AddFeaturingArtistToTrack(ctx, trackId, artist.Id)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				err = tx.Commit()
				if err != nil {
					return nil, err
				}

				{
					track, err := app.DB().GetTrackById(ctx, trackId)
					if err != nil {
						return nil, err
					}

					err = app.DB().InsertTrackToSearch(ctx, track)
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},
	)
}
