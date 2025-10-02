package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/pelletier/go-toml/v2"
)

type GetSystemInfo struct {
	Version string `json:"version"`
}

func fixArr(arr []string) []string {
	seen := map[string]bool{}
	res := make([]string, 0, len(arr))

	for _, value := range arr {
		value = anvil.String(value)
		if value == "" {
			continue
		}

		if !seen[value] {
			seen[value] = true
			res = append(res, value)
		}
	}

	return res
}

// TODO(patrik): Add testing for this
func FixMetadata(metadata *library.Metadata) error {
	album := &metadata.Album

	album.Name = anvil.String(album.Name)

	if album.Year == 0 {
		album.Year = metadata.General.Year
	}

	if len(album.Artists) == 0 {
		album.Artists = []string{UNKNOWN_ARTIST_NAME}
	}

	album.Artists = fixArr(album.Artists)

	for i := range metadata.Tracks {
		t := &metadata.Tracks[i]

		if t.Year == 0 {
			t.Year = metadata.General.Year
		}

		t.Name = anvil.String(t.Name)

		t.Tags = append(t.Tags, metadata.General.Tags...)
		t.Tags = append(t.Tags, metadata.General.TrackTags...)

		if len(t.Artists) == 0 {
			t.Artists = []string{UNKNOWN_ARTIST_NAME}
		}

		t.Artists = fixArr(t.Artists)

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

	albums map[string]struct{}
	tracks map[string]struct{}
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
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (helper *SyncHelper) setAlbumTags(ctx context.Context, db *database.Database, albumId string, tags []string) error {
	err := db.RemoveAllTagsFromAlbum(ctx, albumId)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		slug := utils.Slug(tag)

		err := db.CreateTag(ctx, slug)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddTagToAlbum(ctx, slug, albumId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
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
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (helper *SyncHelper) setTrackTags(ctx context.Context, db *database.Database, trackId string, tags []string) error {
	err := db.RemoveAllTagsFromTrack(ctx, trackId)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		slug := utils.Slug(tag)

		err := db.CreateTag(ctx, slug)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddTagToTrack(ctx, slug, trackId)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
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

	helper.albums[dbAlbum.Id] = struct{}{}

	changes := database.AlbumChanges{}

	// TODO(patrik): More updates

	changes.Name = types.Change[string]{
		Value:   metadata.Album.Name,
		Changed: metadata.Album.Name != dbAlbum.Name,
	}

	artist, err := helper.getOrCreateArtist(ctx, db, metadata.Album.Artists[0])
	if err != nil {
		return fmt.Errorf("failed to create artist for album: %w", err)
	}

	changes.ArtistId = types.Change[string]{
		Value:   artist,
		Changed: artist != dbAlbum.ArtistId,
	}

	changes.CoverArt = types.Change[sql.NullString]{
		Value: sql.NullString{
			String: metadata.General.Cover,
			Valid:  metadata.General.Cover != "",
		},
		Changed: metadata.General.Cover != dbAlbum.CoverArt.String,
	}

	changes.Year = types.Change[sql.NullInt64]{
		Value: sql.NullInt64{
			Int64: metadata.Album.Year,
			Valid: metadata.Album.Year != 0,
		},
		Changed: metadata.Album.Year != dbAlbum.Year.Int64,
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

	err = helper.setAlbumTags(
		ctx,
		db,
		dbAlbum.Id,
		metadata.Album.Tags,
	)
	if err != nil {
		return fmt.Errorf("failed to set album tags: %w", err)
	}

	for i, track := range metadata.Tracks {
		stat, err := os.Stat(track.File)
		if err != nil {
			return fmt.Errorf("failed to stat track[%d] file (%s): %w", i, track.File, err)
		}

		modifiedTime := stat.ModTime().UnixMilli()

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
					ModifiedTime: modifiedTime,
					MediaType:    probeResult.MediaType,
					Name:         track.Name,
					OtherName:    sql.NullString{},
					AlbumId:      dbAlbum.Id,
					ArtistId:     artist,
					Duration:     int64(probeResult.Duration),
					Number: sql.NullInt64{
						Int64: track.Number,
						Valid: track.Number != 0,
					},
					Year: sql.NullInt64{
						Int64: track.Year,
						Valid: track.Year != 0,
					},
				})
				if err != nil {
					return fmt.Errorf("failed to create track[%d]: %w", i, err)
				}

				helper.tracks[trackId] = struct{}{}

				err = helper.setTrackFeaturingArtists(
					ctx,
					db,
					trackId,
					track.Artists[1:],
				)
				if err != nil {
					return fmt.Errorf("failed to set track[%d] featuring artists: %w", i, err)
				}

				err = helper.setTrackTags(ctx, db, trackId, track.Tags)
				if err != nil {
					return fmt.Errorf("failed to set track[%d] tags: %w", i, err)
				}

				continue
			}
		}

		helper.tracks[dbTrack.Id] = struct{}{}

		err = helper.setTrackFeaturingArtists(
			ctx,
			db,
			dbTrack.Id,
			track.Artists[1:],
		)
		if err != nil {
			return fmt.Errorf("failed to set track[%d] featuring artists: %w", i, err)
		}

		err = helper.setTrackTags(ctx, db, dbTrack.Id, track.Tags)
		if err != nil {
			return fmt.Errorf("failed to set track[%d] tags: %w", i, err)
		}

		// TODO(patrik): Check modified time and probe again
		// TODO(patrik): Update track

		changes := database.TrackChanges{}

		if modifiedTime > dbTrack.ModifiedTime {
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
				Value:   modifiedTime,
				Changed: modifiedTime != dbTrack.ModifiedTime,
			}
		}

		// TODO(patrik): Implement all the changes here

		changes.Filename = types.Change[string]{
			Value:   track.File,
			Changed: track.File != dbTrack.Filename,
		}

		changes.Name = types.Change[string]{
			Value:   track.Name,
			Changed: track.Name != dbTrack.Name,
		}

		changes.ArtistId = types.Change[string]{
			Value:   artist,
			Changed: artist != dbTrack.ArtistId,
		}

		changes.Number = types.Change[sql.NullInt64]{
			Value: sql.NullInt64{
				Int64: track.Number,
				Valid: track.Number != 0,
			},
			Changed: track.Number != dbTrack.Number.Int64,
		}

		changes.Year = types.Change[sql.NullInt64]{
			Value: sql.NullInt64{
				Int64: track.Year,
				Valid: track.Year != 0,
			},
			Changed: track.Year != dbTrack.Year.Int64,
		}

		err = db.UpdateTrack(ctx, dbTrack.Id, changes)
		if err != nil {
			return fmt.Errorf("failed to update track[%d]: %w", i, err)
		}
	}

	return nil
}

type ReportType string

const (
	ReportTypeSearch ReportType = "search"
	ReportTypeSync   ReportType = "sync"
)

type SyncError struct {
	Type        ReportType `json:"type"`
	Message     string     `json:"message"`
	FullMessage *string    `json:"fullMessage,omitempty"`
}

type MissingAlbum struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	ArtistName string `json:"artistName"`
}

type MissingTrack struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	AlbumName  string `json:"albumName"`
	ArtistName string `json:"artistName"`
}

type SyncHandler struct {
	broker *Broker

	mutex sync.RWMutex

	isSyncing atomic.Bool

	syncErrors    []SyncError
	missingAlbums []MissingAlbum
	missingTracks []MissingTrack

	SyncChannel chan bool
}

type Report struct {
	SyncErrors    []SyncError    `json:"syncErrors"`
	MissingAlbums []MissingAlbum `json:"missingAlbums"`
	MissingTracks []MissingTrack `json:"missingTracks"`
}

func (s *SyncHandler) GetReport() Report {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return Report{
		SyncErrors:    s.syncErrors,
		MissingAlbums: s.missingAlbums,
		MissingTracks: s.missingTracks,
	}
}

func (s *SyncHandler) Cleanup(app core.App) error {
	tx, err := app.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.TODO()

	for _, track := range s.missingTracks {
		err := tx.DeleteTrack(ctx, track.Id)
		if err != nil {
			return err
		}

		slog.Info("Deleted track", "track", track)
	}

	for _, album := range s.missingAlbums {
		err := tx.DeleteAlbum(ctx, album.Id)
		if err != nil {
			return err
		}

		slog.Info("Deleted album", "album", album)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.missingAlbums = []MissingAlbum{}
	s.missingTracks = []MissingTrack{}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncHandler) RunSync(app core.App) error {
	s.isSyncing.Store(true)
	defer s.isSyncing.Store(false)

	slog.Debug("Searching for albums", "libraryDir", app.Config().LibraryDir)

	// TODO(patrik): Check for duplicated ids
	search, err := library.FindAlbums(app.Config().LibraryDir)
	if err != nil {
		return err
	}

	slog.Debug("Done searching for albums")

	ctx := context.TODO()

	err = EnsureUnknownArtistExists(ctx, app.DB(), app.WorkDir())
	if err != nil {
		return err
	}

	helper := SyncHelper{
		artists: map[string]string{},
		albums:  map[string]struct{}{},
		tracks:  map[string]struct{}{},
	}

	var syncErrors []error

	for _, album := range search.Albums {
		slog.Debug("Syncing album", "path", album.Path)

		err := helper.syncAlbum(ctx, &album.Metadata, app.DB())
		if err != nil {
			syncErrors = append(syncErrors, err)
		}
	}

	var missingAlbums []MissingAlbum
	var missingTracks []MissingTrack

	{
		ids, err := app.DB().GetAllAlbumIds(ctx)
		if err != nil {
			return err
		}

		for _, id := range ids {
			_, exists := helper.albums[id]
			if !exists {
				album, err := app.DB().GetAlbumById(ctx, id)
				if err != nil {
					// TODO(patrik): How should we handle the error?
					continue
				}

				missingAlbums = append(missingAlbums, MissingAlbum{
					Id:         id,
					Name:       album.Name,
					ArtistName: album.ArtistName,
				})
			}
		}
	}

	{
		ids, err := app.DB().GetAllTrackIds(ctx)
		if err != nil {
			return err
		}

		for _, id := range ids {
			_, exists := helper.tracks[id]
			if !exists {
				track, err := app.DB().GetTrackById(ctx, id)
				if err != nil {
					// TODO(patrik): How should we handle the error?
					continue
				}

				missingTracks = append(missingTracks, MissingTrack{
					Id:         id,
					Name:       track.Name,
					AlbumName:  track.AlbumName,
					ArtistName: track.ArtistName,
				})
			}
		}
	}

	errs := make([]SyncError, 0, len(search.Errors)+len(syncErrors))

	for _, err := range search.Errors {
		var fullMessage *string

		var tomlError *toml.DecodeError
		if errors.As(err, &tomlError) {
			m := tomlError.String()
			fullMessage = &m
		}

		errs = append(errs, SyncError{
			Type:        ReportTypeSearch,
			Message:     err.Error(),
			FullMessage: fullMessage,
		})
	}

	for _, err := range syncErrors {
		errs = append(errs, SyncError{
			Type:    ReportTypeSync,
			Message: err.Error(),
		})
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.syncErrors = errs
	s.missingAlbums = missingAlbums
	s.missingTracks = missingTracks

	return nil
}

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type EventData interface {
	GetEventType() string
}

// NOTE(patrik): Based on: https://gist.github.com/Ananto30/8af841f250e89c07e122e2a838698246
type Broker struct {
	Notifier chan EventData

	newClients     chan chan EventData
	closingClients chan chan EventData
	clients        map[chan EventData]bool
}

func NewServer() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan EventData, 1),
		newClients:     make(chan chan EventData),
		closingClients: make(chan chan EventData),
		clients:        make(map[chan EventData]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:
			broker.clients[s] = true
			slog.Debug("Client added", "numClients", len(broker.clients))
		case s := <-broker.closingClients:
			delete(broker.clients, s)
			slog.Debug("Removed client", "numClients", len(broker.clients))
		case event := <-broker.Notifier:
			for clientMessageChan := range broker.clients {
				clientMessageChan <- event
			}
		}
	}
}

func (broker *Broker) EmitEvent(event EventData) {
	broker.Notifier <- event
}

var syncHandler = SyncHandler{
	SyncChannel: make(chan bool),
	broker:      NewServer(),
}

const (
	EventSyncing string = "syncing"
	EventReport  string = "report"
)

type SyncEvent struct {
	Syncing bool `json:"syncing"`
}

func (s SyncEvent) GetEventType() string {
	return EventSyncing
}

type ReportEvent struct {
	Report
}

func (s ReportEvent) GetEventType() string {
	return EventReport
}

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
			Name:   "GetSyncStatus",
			Method: http.MethodGet,
			Path:   "/system/library",
			// ResponseType: SyncStatus{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// res := syncHandler.GetSyncStatus()
				return nil, nil
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
					if syncHandler.isSyncing.Load() {
						slog.Info("Syncing already")
						return
					}

					slog.Info("Started library sync")

					syncHandler.broker.EmitEvent(SyncEvent{
						Syncing: true,
					})

					err := syncHandler.RunSync(app)
					if err != nil {
						slog.Error("Failed to run sync", "err", err)
					}

					syncHandler.broker.EmitEvent(SyncEvent{
						Syncing: false,
					})

					syncHandler.broker.EmitEvent(ReportEvent{
						Report: syncHandler.GetReport(),
					})

					slog.Info("Library sync done")
				}()

				return nil, nil
			},
		},

		// TODO(patrik): Better name?
		pyrin.ApiHandler{
			Name:   "CleanupLibrary",
			Method: http.MethodPost,
			Path:   "/system/library/cleanup",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				if syncHandler.isSyncing.Load() {
					return nil, errors.New("library is syncing")
				}

				err := syncHandler.Cleanup(app)
				if err != nil {
					return nil, err
				}

				syncHandler.broker.EmitEvent(ReportEvent{
					Report: syncHandler.GetReport(),
				})

				return nil, nil
			},
		},

		pyrin.NormalHandler{
			Name:   "SseHandler",
			Method: http.MethodGet,
			Path:   "/system/library/sse",
			HandlerFunc: func(c pyrin.Context) error {
				r := c.Request()
				w := c.Response()

				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")

				w.Header().Set("Access-Control-Allow-Origin", "*")

				rc := http.NewResponseController(w)

				eventChan := make(chan EventData)
				syncHandler.broker.newClients <- eventChan

				defer func() {
					syncHandler.broker.closingClients <- eventChan
				}()

				sendEvent := func(eventData EventData) {
					fmt.Fprintf(w, "data: ")

					event := Event{
						Type: eventData.GetEventType(),
						Data: eventData,
					}

					encode := json.NewEncoder(w)
					encode.Encode(event)

					fmt.Fprintf(w, "\n\n")
					rc.Flush()
				}

				sendEvent(SyncEvent{
					Syncing: syncHandler.isSyncing.Load(),
				})

				sendEvent(ReportEvent{
					Report: syncHandler.GetReport(),
				})

				for {
					select {
					case <-r.Context().Done():
						syncHandler.broker.closingClients <- eventChan
						return nil

					case event := <-eventChan:
						sendEvent(event)
					}
				}
			},
		},
	)
}
