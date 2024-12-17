package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database/adapter"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type ExtraArtist struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	OtherName *string `json:"other_name"`
}

type ExtraArtists []ExtraArtist

func (s ExtraArtists) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (v *ExtraArtists) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch value := value.(type) {
	case string:
		return json.Unmarshal([]byte(value), &v)
	case []byte:
		return json.Unmarshal(value, &v)
	default:
		return errors.New("type assertion to ([]byte / string) failed")
	}
}

type Track struct {
	Id        string         `db:"id"`
	Name      string         `db:"name"`
	OtherName sql.NullString `db:"other_name"`

	AlbumId  string `db:"album_id"`
	ArtistId string `db:"artist_id"`

	Number   sql.NullInt64 `db:"track_number"`
	Duration sql.NullInt64 `db:"duration"`
	Year     sql.NullInt64 `db:"year"`

	ExportName       string `db:"export_name"`
	OriginalFilename string `db:"original_filename"`
	MobileFilename   string `db:"mobile_filename"`

	AlbumName     string         `db:"album_name"`
	AlbumCoverArt sql.NullString `db:"album_cover_art"`
	ArtistName    string         `db:"artist_name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`

	ExtraArtists   ExtraArtists   `db:"extra_artists"`
}

// TODO(patrik): Move down to the UpdateTrack
type TrackChanges struct {
	Name      types.Change[string]
	OtherName types.Change[sql.NullString]

	AlbumId  types.Change[string]
	ArtistId types.Change[string]

	Number   types.Change[sql.NullInt64]
	Duration types.Change[sql.NullInt64]
	Year     types.Change[sql.NullInt64]

	ExportName       types.Change[string]
	OriginalFilename types.Change[string]
	MobileFilename   types.Change[string]

	Created types.Change[int64]
}

func TrackQuery() *goqu.SelectDataset {
	// TODO(patrik): Add Back
	tags := dialect.From("tracks_to_tags").
		Select(
			goqu.I("tracks_to_tags.track_id").As("track_id"),
			goqu.Func("group_concat", goqu.I("tags.name"), ",").As("tags"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(goqu.I("tracks_to_tags.tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(goqu.I("tracks_to_tags.track_id"))

	extraArtists := dialect.From("tracks_extra_artists").
		Select(
			goqu.I("tracks_extra_artists.track_id").As("track_id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",
					"id",
					goqu.I("artists.id"),
					"name",
					goqu.I("artists.name"),
					"other_name",
					goqu.I("artists.other_name"),
				),
			).As("artists"),
		).
		Join(
			goqu.I("artists"),
			goqu.On(goqu.I("tracks_extra_artists.artist_id").Eq(goqu.I("artists.id"))),
		).
		GroupBy(goqu.I("tracks_extra_artists.track_id"))

	query := dialect.From("tracks").
		Select(
			"tracks.id",
			"tracks.name",
			"tracks.other_name",

			"tracks.album_id",
			"tracks.artist_id",

			"tracks.track_number",
			"tracks.duration",
			"tracks.year",

			"tracks.export_name",
			"tracks.original_filename",
			"tracks.mobile_filename",

			"tracks.created",
			"tracks.updated",

			goqu.I("albums.name").As("album_name"),
			goqu.I("albums.cover_art").As("album_cover_art"),
			goqu.I("artists.name").As("artist_name"),

			goqu.I("tags.tags").As("tags"),

			goqu.I("extra_artists.artists").As("extra_artists"),
		).
		Prepared(true).
		Join(
			goqu.I("albums"),
			goqu.On(goqu.I("tracks.album_id").Eq(goqu.I("albums.id"))),
		).
		Join(
			goqu.I("artists"),
			goqu.On(goqu.I("tracks.artist_id").Eq(goqu.I("artists.id"))),
		).
		LeftJoin(
			tags.As("tags"),
			goqu.On(goqu.I("tracks.id").Eq(goqu.I("tags.track_id"))),
		).
		LeftJoin(
			extraArtists.As("extra_artists"),
			goqu.On(goqu.I("tracks.id").Eq(goqu.I("extra_artists.track_id"))),
		)

	return query
}

// TODO(patrik): Remove?
type Scan interface {
	Scan(dest ...any) error
}

// TODO(patrik): Move
type FetchOption struct {
	Filter  string
	Sort    string
	PerPage int
	Page    int
}

func (db *Database) GetAllTracks(ctx context.Context, filterStr, sortStr string) ([]Track, error) {
	query := TrackQuery()

	var err error

	a := adapter.TrackResolverAdapter{}
	resolver := filter.New(&a)

	query, err = applyFilter(query, resolver, filterStr)
	if err != nil {
		return nil, err
	}

	query, err = applySort(query, resolver, filterStr)
	if err != nil {
		return nil, err
	}

	var items []Track
	err = db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetPagedTracks(ctx context.Context, opts FetchOption) ([]Track, types.Page, error) {
	query := TrackQuery()

	var err error

	a := adapter.TrackResolverAdapter{}
	resolver := filter.New(&a)

	query, err = applyFilter(query, resolver, opts.Filter)
	if err != nil {
		return nil, types.Page{}, err
	}

	query, err = applySort(query, resolver, opts.Sort)
	if err != nil {
		return nil, types.Page{}, err
	}

	countQuery := query.
		Select(goqu.COUNT("tracks.id"))

	if opts.PerPage > 0 {
		query = query.
			Limit(uint(opts.PerPage)).
			Offset(uint(opts.Page * opts.PerPage))
	}

	var totalItems int
	err = db.Get(&totalItems, countQuery)
	if err != nil {
		return nil, types.Page{}, err
	}

	totalPages := utils.TotalPages(opts.PerPage, totalItems)
	page := types.Page{
		Page:       opts.Page,
		PerPage:    opts.PerPage,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	var items []Track
	err = db.Select(&items, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db *Database) GetTracksByAlbum(ctx context.Context, albumId string) ([]Track, error) {
	query := TrackQuery().
		Where(goqu.I("tracks.album_id").Eq(albumId)).
		Order(goqu.I("tracks.track_number").Asc().NullsLast(), goqu.I("tracks.name").Asc())

	var items []Track
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetTrackById(ctx context.Context, id string) (Track, error) {
	query := TrackQuery().
		Where(goqu.I("tracks.id").Eq(id))

	var item Track
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Track{}, ErrItemNotFound
		}

		return Track{}, err
	}

	return item, nil
}

func (db *Database) GetTrackByNameAndAlbum(ctx context.Context, name string, albumId string) (Track, error) {
	query := TrackQuery().
		Where(
			goqu.And(
				goqu.I("tracks.name").Eq(name),
				goqu.I("tracks.album_id").Eq(albumId),
			),
		)

	var item Track
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Track{}, ErrItemNotFound
		}

		return Track{}, err
	}

	return item, nil
}

type CreateTrackParams struct {
	Id        string
	Name      string
	OtherName sql.NullString

	AlbumId  string
	ArtistId string

	Number   sql.NullInt64
	Duration sql.NullInt64
	Year     sql.NullInt64

	ExportName       string
	OriginalFilename string
	MobileFilename   string

	Created int64
	Updated int64
}

func (db *Database) CreateTrack(ctx context.Context, params CreateTrackParams) (string, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateId()
	}

	ds := dialect.Insert("tracks").Rows(goqu.Record{
		"id":         id,
		"name":       params.Name,
		"other_name": params.OtherName,

		"album_id":  params.AlbumId,
		"artist_id": params.ArtistId,

		"track_number": params.Number,
		"duration":     params.Duration,
		"year":         params.Year,

		"export_name":       params.ExportName,
		"original_filename": params.OriginalFilename,
		"mobile_filename":   params.MobileFilename,

		"created": created,
		"updated": updated,
	}).
		Returning("id").
		Prepared(true)

	var item string

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return "", err
	}

	err = row.Scan(&item)
	if err != nil {
		return "", err
	}

	return item, nil
}

// TODO(patrik): Move
func addToRecord[T any](record goqu.Record, name string, change types.Change[T]) {
	if change.Changed {
		record[name] = change.Value
	}
}

func (db *Database) UpdateTrack(ctx context.Context, id string, changes TrackChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "other_name", changes.OtherName)

	addToRecord(record, "album_id", changes.AlbumId)
	addToRecord(record, "artist_id", changes.ArtistId)

	addToRecord(record, "track_number", changes.Number)
	addToRecord(record, "duration", changes.Duration)
	addToRecord(record, "year", changes.Year)

	addToRecord(record, "export_name", changes.ExportName)
	addToRecord(record, "original_filename", changes.OriginalFilename)
	addToRecord(record, "mobile_filename", changes.MobileFilename)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("tracks").
		Set(record).
		Where(goqu.I("tracks.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAlbumTracks(ctx context.Context, albumId string) error {
	query := dialect.Delete("tracks").
		Prepared(true).
		Where(goqu.I("tracks.album_id").Eq(albumId))

	// TODO(patrik): Check result?
	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveTrack(ctx context.Context, id string) error {
	query := dialect.Delete("tracks").
		Prepared(true).
		Where(goqu.I("tracks.id").Eq(id))

	// TODO(patrik): Check result?
	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): This need cleanup
var tags []Tag

func (db *Database) Invalidate() {
	log.Debug("Database.Invalidate")

	ctx := context.Background()

	tags, _ = db.GetAllTags(ctx)
}

func TrackMapNameToId(typ string, name string) string {
	switch typ {
	case "tags":
		for _, t := range tags {
			if t.Name == strings.ToLower(name) {
				return t.Name
			}
		}
	}

	return ""
}

func (db *Database) RemoveAllExtraArtistsFromTrack(ctx context.Context, trackId string) error {
	query := dialect.Delete("tracks_extra_artists").
		Where(goqu.I("tracks_extra_artists.track_id").Eq(trackId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddExtraArtistToTrack(ctx context.Context, trackId, artistId string) error {
	query := dialect.Insert("tracks_extra_artists").
		Rows(goqu.Record{
			"track_id":  trackId,
			"artist_id": artistId,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) {
			if sqlErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}
