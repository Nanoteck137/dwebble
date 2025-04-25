package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/pelletier/go-toml/v2"
)

type GetSystemInfo struct {
	Version string `json:"version"`
}

// TODO(patrik): Add testing for this
func FixMetadata(metadata *library.Metadata) error {
	album := &metadata.Album

	album.Name = transform.String(album.Name)

	if album.Year == 0 {
		album.Year = metadata.General.Year
	}

	if len(album.Artists) == 0 {
		album.Artists = []string{UNKNOWN_ARTIST_NAME}
	}

	for i := range metadata.Tracks {
		t := &metadata.Tracks[i]

		if t.Year == 0 {
			t.Year = metadata.General.Year
		}

		t.Name = transform.String(t.Name)

		t.Tags = append(t.Tags, metadata.General.Tags...)
		t.Tags = append(t.Tags, metadata.General.TrackTags...)

		if len(t.Artists) == 0 {
			t.Artists = []string{UNKNOWN_ARTIST_NAME}
		}

		for i, tag := range t.Tags {
			t.Tags[i] = utils.Slug(strings.TrimSpace(tag))
		}
	}

	// err := validate.ValidateStruct(&metadata.Album,
	// 	validate.Field(&metadata.Album.Name, validate.Required),
	// 	validate.Field(&metadata.Album.Artists, validate.Length(1, 0)),
	// )
	// if err != nil {
	// 	return err
	// }
	//
	// for _, track := range metadata.Tracks {
	// 	err := validate.ValidateStruct(&track,
	// 		validate.Field(&track.File, validate.Required),
	// 		validate.Field(&track.Name, validate.Required),
	// 		validate.Field(&track.Artists, validate.Length(1, 0)),
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type SyncHelper struct {
	artists map[string]string
}

func (helper *SyncHelper) getOrCreateArtist(ctx context.Context, db *database.Database, name string) (string, error) {
	slug := utils.Slug(name)

	if artist, exists := helper.artists[slug]; exists {
		return artist, nil
	}

	dbArtist, err := db.GetArtistBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			dbArtist, err = db.CreateArtist(ctx, database.CreateArtistParams{
				Slug: slug,
				Name: name,
			})
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	helper.artists[slug] = dbArtist.Id
	return dbArtist.Id, nil
}

func (helper *SyncHelper) setAlbumFeaturingArtists(ctx context.Context, db *database.Database, albumId string, artists []string) error {
	err := db.RemoveAllAlbumFeaturingArtists(ctx, albumId)
	if err != nil {
		return err
	}

	for _, artistName := range artists {
		artistId, err := helper.getOrCreateArtist(ctx, db, artistName)
		if err != nil {
			return err
		}

		err = db.AddFeaturingArtistToAlbum(ctx, albumId, artistId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (helper *SyncHelper) setTrackFeaturingArtists(ctx context.Context, db *database.Database, trackId string, artists []string) error {
	err := db.RemoveAllTrackFeaturingArtists(ctx, trackId)
	if err != nil {
		return err
	}

	for _, artistName := range artists {
		artistId, err := helper.getOrCreateArtist(ctx, db, artistName)
		if err != nil {
			return err
		}

		err = db.AddFeaturingArtistToTrack(ctx, trackId, artistId)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO(patrik): Update the errors for album
func (helper *SyncHelper) syncAlbum(ctx context.Context, metadata *library.Metadata, db *database.Database) error {
	err := FixMetadata(metadata)
	if err != nil {
		return err
	}

	dbAlbum, err := db.GetAlbumById(ctx, metadata.Album.Id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			artist, err := helper.getOrCreateArtist(ctx, db, metadata.Album.Artists[0])
			if err != nil {
				return fmt.Errorf("failed to create artist for album: %w", err)
			}

			dbAlbum, err = db.CreateAlbum(ctx, database.CreateAlbumParams{
				Id:       metadata.Album.Id,
				Name:     metadata.Album.Name,
				ArtistId: artist,
			})
			if err != nil {
				return fmt.Errorf("failed to create album: %w", err)
			}
		} else {
			return err
		}
	}

	changes := database.AlbumChanges{}

	// TODO(patrik): More updates
	changes.CoverArt = types.Change[sql.NullString]{
		Value: sql.NullString{
			String: metadata.General.Cover,
			Valid:  metadata.General.Cover != "",
		},
		Changed: metadata.General.Cover != dbAlbum.CoverArt.String,
	}

	err = db.UpdateAlbum(ctx, dbAlbum.Id, changes)
	if err != nil {
		return fmt.Errorf("failed to update album: %w", err)
	}

	err = helper.setAlbumFeaturingArtists(
		ctx,
		db,
		dbAlbum.Id,
		metadata.Album.Artists[1:],
	)
	if err != nil {
		return fmt.Errorf("failed to set album featuring artists: %w", err)
	}

	for i, track := range metadata.Tracks {
		artist, err := helper.getOrCreateArtist(ctx, db, track.Artists[0])
		if err != nil {
			return fmt.Errorf("failed to set create artist for track[%d]: %w", i, err)
		}

		dbTrack, err := db.GetTrackById(ctx, track.Id)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				probeResult, err := utils.ProbeTrack(track.File)
				if err != nil {
					return fmt.Errorf("failed to probe track[%d] file (%s): %w", i, track.File, err)
				}

				trackId, err := db.CreateTrack(ctx, database.CreateTrackParams{
					Id:           track.Id,
					Filename:     track.File,
					ModifiedTime: track.ModifiedTime,
					MediaType:    probeResult.MediaType,
					Name:         track.Name,
					OtherName:    sql.NullString{},
					AlbumId:      dbAlbum.Id,
					ArtistId:     artist,
					Duration:     int64(probeResult.Duration),
					Number: sql.NullInt64{
						Int64: int64(track.Number),
						Valid: track.Number != 0,
					},
					Year: sql.NullInt64{
						Int64: int64(track.Year),
						Valid: track.Year != 0,
					},
				})
				if err != nil {
					return fmt.Errorf("failed to create track[%d]: %w", i, err)
				}

				err = helper.setTrackFeaturingArtists(
					ctx,
					db,
					trackId,
					track.Artists[1:],
				)
				if err != nil {
					return fmt.Errorf("failed to set track[%d] featuring artists: %w", i, err)
				}

				continue
			}
		}

		err = helper.setTrackFeaturingArtists(
			ctx,
			db,
			dbTrack.Id,
			track.Artists[1:],
		)
		if err != nil {
			return fmt.Errorf("failed to set track[%d] featuring artists: %w", i, err)
		}

		// TODO(patrik): Check modified time and probe again
		// TODO(patrik): Update track

		changes := database.TrackChanges{}

		if track.ModifiedTime > dbTrack.ModifiedTime {
			probeResult, err := utils.ProbeTrack(track.File)
			if err != nil {
				return fmt.Errorf("failed to probe track[%d] file (%s): %w", i, track.File, err)
			}

			dur := int64(probeResult.Duration)
			changes.Duration = types.Change[int64]{
				Value:   dur,
				Changed: dur != dbTrack.Duration,
			}

			changes.MediaType = types.Change[types.MediaType]{
				Value:   probeResult.MediaType,
				Changed: probeResult.MediaType != dbTrack.MediaType,
			}

			changes.ModifiedTime = types.Change[int64]{
				Value:   track.ModifiedTime,
				Changed: track.ModifiedTime != dbTrack.ModifiedTime,
			}
		}

		// TODO(patrik): Implement all the changes here

		changes.Filename = types.Change[string]{
			Value:   track.File,
			Changed: dbTrack.Filename != track.File,
		}

		changes.Name = types.Change[string]{
			Value:   track.Name,
			Changed: dbTrack.Name != track.Name,
		}

		changes.Year = types.Change[sql.NullInt64]{
			Value: sql.NullInt64{
				Int64: int64(track.Year),
				Valid: track.Year != 0,
			},
			Changed: dbTrack.Year.Int64 != int64(track.Year),
		}

		err = db.UpdateTrack(ctx, dbTrack.Id, changes)
		if err != nil {
			return fmt.Errorf("failed to update track[%d]: %w", i, err)
		}
	}

	return nil
}

type SyncStatus struct {
	IsSyncing   bool     `json:"isSyncing"`
	LastReports []Report `json:"lastReports"`
}

type ReportType string

const (
	ReportTypeSearch ReportType = "search"
	ReportTypeSync   ReportType = "sync"
)

type Report struct {
	Type        ReportType `json:"type"`
	Message     string     `json:"message"`
	FullMessage *string    `json:"fullMessage,omitempty"`
}

type SyncHandler struct {
	mutex sync.RWMutex

	isSyncing   bool
	lastReports []Report
}

func (s *SyncHandler) GetSyncStatus() SyncStatus {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return SyncStatus{
		IsSyncing:   s.isSyncing,
		LastReports: s.lastReports,
	}
}

func (s *SyncHandler) SetSyncing(syncing bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.isSyncing = syncing
}

func (s *SyncHandler) GetSyncing(syncing bool) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.isSyncing
}

func (s *SyncHandler) RunSync(app core.App) error {
	s.SetSyncing(true)
	defer s.SetSyncing(false)

	// TODO(patrik): Check for duplicated ids
	search, err := library.FindAlbums("/Volumes/media/test/Ado/")
	if err != nil {
		return err
	}

	ctx := context.TODO()

	err = EnsureUnknownArtistExists(ctx, app.DB(), app.WorkDir())
	if err != nil {
		return err
	}

	helper := SyncHelper{
		artists: map[string]string{},
	}

	var syncErrors []error

	for _, album := range search.Albums {
		err := helper.syncAlbum(ctx, &album.Metadata, app.DB())
		if err != nil {
			syncErrors = append(syncErrors, err)
		}
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, err := range search.Errors {
		var fullMessage *string

		var tomlError *toml.DecodeError
		if errors.As(err, &tomlError) {
			m := tomlError.String()
			fullMessage = &m
		}

		s.lastReports = append(s.lastReports, Report{
			Type:        ReportTypeSearch,
			Message:     err.Error(),
			FullMessage: fullMessage,
		})
	}

	for _, err := range syncErrors {
		s.lastReports = append(s.lastReports, Report{
			Type:    ReportTypeSync,
			Message: err.Error(),
		})
	}

	return nil
}

var syncHandler = SyncHandler{}

func InstallSystemHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetSystemInfo",
			Path:         "/system/info",
			Method:       http.MethodGet,
			ResponseType: GetSystemInfo{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				return GetSystemInfo{
					Version: dwebble.Version,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "RefillSearch",
			Path:   "/system/search",
			Method: http.MethodPost,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()
				err = app.DB().RefillSearchTables(ctx)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetSyncStatus",
			Method:       http.MethodGet,
			Path:         "/system/library",
			ResponseType: SyncStatus{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				res := syncHandler.GetSyncStatus()
				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SyncLibrary",
			Method:       http.MethodPost,
			Path:         "/system/library",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik):
				//  - Handle Errors
				//  - Handle Single Album Syncing
				//  - Handle album modified syncing

				go func() {
					log.Info("Started library sync")

					err := syncHandler.RunSync(app)
					if err != nil {
						log.Error("Failed to run sync", "err", err)
					}

					log.Info("Library sync done")
				}()

				return nil, nil
			},
		},
	)
}
