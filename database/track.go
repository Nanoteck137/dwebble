package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/dwebble/database/adapter"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type Track struct {
	RowId int `db:"rowid"`

	Id string `db:"id"`

	Filename  string `db:"filename"`
	MediaType types.MediaType `db:"media_type"`

	Name      string         `db:"name"`
	OtherName sql.NullString `db:"other_name"`

	AlbumId  string `db:"album_id"`
	ArtistId string `db:"artist_id"`

	Duration int64         `db:"duration"`
	Number   sql.NullInt64 `db:"number"`
	Year     sql.NullInt64 `db:"year"`

	OriginalFilename string `db:"original_filename"`
	MobileFilename   string `db:"mobile_filename"`

	AlbumName      string         `db:"album_name"`
	AlbumOtherName sql.NullString `db:"album_other_name"`
	AlbumCoverArt  sql.NullString `db:"album_cover_art"`

	ArtistName      string         `db:"artist_name"`
	ArtistOtherName sql.NullString `db:"artist_other_name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`

	FeaturingArtists FeaturingArtists `db:"featuring_artists"`
}

// TODO(patrik): Move
func FeaturingArtistsQuery(table, idColName string) *goqu.SelectDataset {
	tbl := goqu.T(table)

	return dialect.From(tbl).
		Select(
			tbl.Col(idColName).As("id"),
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
			goqu.On(tbl.Col("artist_id").Eq(goqu.I("artists.id"))),
		).
		GroupBy(tbl.Col(idColName))
}

// TODO(patrik): Use goqu.T more
func TrackQuery() *goqu.SelectDataset {
	tags := dialect.From("tracks_tags").
		Select(
			goqu.I("tracks_tags.track_id").As("track_id"),
			goqu.Func("group_concat", goqu.I("tags.slug"), ",").As("tags"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(goqu.I("tracks_tags.tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(goqu.I("tracks_tags.track_id"))

	query := dialect.From("tracks").
		Select(
			"tracks.rowid",

			"tracks.id",

			"tracks.filename",
			"tracks.media_type",

			"tracks.name",
			"tracks.other_name",

			"tracks.album_id",
			"tracks.artist_id",

			"tracks.number",
			"tracks.duration",
			"tracks.year",

			"tracks.created",
			"tracks.updated",

			goqu.I("albums.name").As("album_name"),
			goqu.I("albums.other_name").As("album_other_name"),
			goqu.I("albums.cover_art").As("album_cover_art"),

			goqu.I("artists.name").As("artist_name"),
			goqu.I("artists.other_name").As("artist_other_name"),

			goqu.I("tags.tags").As("tags"),

			goqu.I("featuring_artists.artists").As("featuring_artists"),
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
			FeaturingArtistsQuery("tracks_featuring_artists", "track_id").As("featuring_artists"),
			goqu.On(goqu.I("tracks.id").Eq(goqu.I("featuring_artists.id"))),
		)

	return query
}

func (db *Database) GetAllTracksByArtistId(ctx context.Context, artistId string) ([]Track, error) {
	query := TrackQuery().
		Where(goqu.I("tracks.artist_id").Eq(artistId))

	var items []Track
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetTracksByIds(ctx context.Context, ids []string) ([]Track, error) {
	query := TrackQuery().Where(goqu.I("tracks.id").In(ids))

	var items []Track
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetTracksByAlbumForPlay(ctx context.Context, albumId string) ([]NewTrackQueryItem, error) {
	query := NewTrackQuery().
		Where(goqu.I("tracks.album_id").Eq(albumId)).
		Order(goqu.I("tracks.number").Asc().NullsLast(), goqu.I("tracks.name").Asc()).
		As("tracks")

	var items []NewTrackQueryItem
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// TODO(patrik): Move
type FetchOptions struct {
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

	query, err = applySort(query, resolver, sortStr)
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

func (db *Database) GetPagedTracks(ctx context.Context, opts FetchOptions) ([]Track, types.Page, error) {
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
		Where(
			goqu.I("tracks.id").In(
				goqu.From("tracks").
					Select("tracks.id").
					Where(goqu.I("tracks.album_id").Eq(albumId)),
			),
		).
		Order(
			goqu.I("tracks.number").Asc().NullsLast(),
			goqu.I("tracks.name").Asc(),
		)

	var items []Track
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func AlbumTrackSubquery(albumId string) *goqu.SelectDataset {
	return goqu.From("tracks").
		Select("tracks.id").
		Where(goqu.I("tracks.album_id").Eq(albumId))
}

func ArtistTrackSubquery(artistId string) *goqu.SelectDataset {
	tbl := goqu.T("tracks_featuring_artists")
	return goqu.From("tracks").
		Select("tracks.id").
		FullOuterJoin(
			tbl,
			goqu.On(
				goqu.I("tracks.id").Eq(tbl.Col("track_id")),
			),
		).
		Where(
			goqu.Or(
				goqu.I("tracks.artist_id").Eq(artistId),
				tbl.Col("artist_id").Eq(artistId),
			),
		)
}

func PlaylistTrackSubquery(playlistId string) *goqu.SelectDataset {
	return goqu.From("playlist_items").
		Select("playlist_items.track_id").
		Where(
			goqu.I("playlist_items.playlist_id").Eq(playlistId),
		)
}

func (db *Database) GetTracksIn(ctx context.Context, in any, sort string) ([]Track, error) {
	query := TrackQuery().
		Where(
			goqu.I("tracks.id").In(in),
		)

	a := adapter.TrackResolverAdapter{}
	resolver := filter.New(&a)

	query, err := applySort(query, resolver, sort)
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
	Id string

	Filename  string
	MediaType types.MediaType

	Name      string
	OtherName sql.NullString

	AlbumId  string
	ArtistId string

	Duration int64
	Number   sql.NullInt64
	Year     sql.NullInt64

	// OriginalFilename string
	// MobileFilename   string

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
		"id": id,

		"filename":   params.Filename,
		"media_type": params.MediaType,

		"name":       params.Name,
		"other_name": params.OtherName,

		"album_id":  params.AlbumId,
		"artist_id": params.ArtistId,

		"duration": params.Duration,
		"number":   params.Number,
		"year":     params.Year,

		"created": created,
		"updated": updated,
	}).
		Prepared(true).
		Returning("id")

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

type TrackChanges struct {
	Filename  types.Change[string]
	MediaType types.Change[types.MediaType]

	Name      types.Change[string]
	OtherName types.Change[sql.NullString]

	AlbumId  types.Change[string]
	ArtistId types.Change[string]

	Duration types.Change[int64]
	Number   types.Change[sql.NullInt64]
	Year     types.Change[sql.NullInt64]

	Created types.Change[int64]
}

func (db *Database) UpdateTrack(ctx context.Context, id string, changes TrackChanges) error {
	record := goqu.Record{}

	addToRecord(record, "filename", changes.Filename)
	addToRecord(record, "media_type", changes.MediaType)

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "other_name", changes.OtherName)

	addToRecord(record, "album_id", changes.AlbumId)
	addToRecord(record, "artist_id", changes.ArtistId)

	addToRecord(record, "duration", changes.Duration)
	addToRecord(record, "number", changes.Number)
	addToRecord(record, "year", changes.Year)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("tracks").
		Prepared(true).
		Set(record).
		Where(goqu.I("tracks.id").Eq(id))

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) ChangeAllTrackArtist(ctx context.Context, artistId, newArtistId string) error {
	record := goqu.Record{
		"artist_id": newArtistId,
		"updated":   time.Now().UnixMilli(),
	}

	query := goqu.Update("tracks").
		Prepared(true).
		Set(record).
		Where(goqu.I("tracks.artist_id").Eq(artistId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return err
}

func (db *Database) DeleteTrackMedia(ctx context.Context, trackId string) error {
	query := dialect.Delete("tracks_media").
		Prepared(true).
		Where(goqu.I("tracks_media.track_id").Eq(trackId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteTrack(ctx context.Context, id string) error {
	query := dialect.Delete("tracks").
		Prepared(true).
		Where(goqu.I("tracks.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func TrackMapNameToId(typ string, name string) string {
	switch typ {
	case "tags":
		return name
	}

	return ""
}

func (db *Database) AddTagToTrack(ctx context.Context, tagSlug, trackId string) error {
	ds := dialect.Insert("tracks_tags").
		Prepared(true).
		Rows(goqu.Record{
			"track_id": trackId,
			"tag_slug": tagSlug,
		})

	_, err := db.Exec(ctx, ds)
	if err != nil {
		var e sqlite3.Error
		if errors.As(err, &e) {
			if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (db *Database) RemoveAllTagsFromTrack(ctx context.Context, trackId string) error {
	query := dialect.Delete("tracks_tags").
		Prepared(true).
		Where(goqu.I("track_id").Eq(trackId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db *Database) RemoveAllTrackFeaturingArtists(ctx context.Context, trackId string) error {
	query := dialect.Delete("tracks_featuring_artists").
		Prepared(true).
		Where(goqu.I("tracks_featuring_artists.track_id").Eq(trackId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db *Database) AddFeaturingArtistToTrack(ctx context.Context, trackId, artistId string) error {
	query := dialect.Insert("tracks_featuring_artists").
		Prepared(true).
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

func (db *Database) RemoveFeaturingArtistFromTrack(ctx context.Context, trackId, artistId string) error {
	query := goqu.Delete("tracks_featuring_artists").
		Prepared(true).
		Where(
			goqu.And(
				goqu.I("tracks_featuring_artists.track_id").Eq(trackId),
				goqu.I("tracks_featuring_artists.artist_id").Eq(artistId),
			),
		)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type CreateTrackMediaParams struct {
	Id      string
	TrackId string

	Filename   string
	MediaType  types.MediaType
	Rank       int
	IsOriginal bool
}

func (db *Database) CreateTrackMedia(ctx context.Context, params CreateTrackMediaParams) error {
	query := dialect.Insert("tracks_media").
		Rows(goqu.Record{
			"id":       params.Id,
			"track_id": params.TrackId,

			"filename":    params.Filename,
			"media_type":  params.MediaType,
			"rank":        params.Rank,
			"is_original": params.IsOriginal,
		}).
		Prepared(true)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
