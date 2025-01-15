package cmd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/spf13/cobra"
	"gopkg.in/vansante/go-ffprobe.v2"

	olddb "github.com/nanoteck137/dwebble/cmd/dwebble-migrate/database"
)

type OldWorkDir string

func (d OldWorkDir) String() string {
	return string(d)
}

func (d OldWorkDir) DatabaseFile() string {
	return path.Join(d.String(), "data.db")
}

func (d OldWorkDir) ExportFile() string {
	return path.Join(d.String(), "export.json")
}

func (d OldWorkDir) SetupFile() string {
	return path.Join(d.String(), "setup")
}

func (d OldWorkDir) Trash() string {
	return path.Join(d.String(), "trash")
}

func (d OldWorkDir) Artists() string {
	return path.Join(d.String(), "artists")
}

func (d OldWorkDir) Artist(id string) string {
	return path.Join(d.Artists(), id)
}

func (d OldWorkDir) Albums() string {
	return path.Join(d.String(), "albums")
}

func (d OldWorkDir) Album(id string) string {
	return path.Join(d.Albums(), id)
}

func (d OldWorkDir) Tracks() string {
	return path.Join(d.String(), "tracks")
}

func (d OldWorkDir) Track(id string) string {
	return path.Join(d.Tracks(), id)
}

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command, args []string) {
		app := core.NewBaseApp(&config.Config{
			DataDir: "./work",
		})

		err := app.Bootstrap()
		if err != nil {
			log.Fatal("Failed to bootstrap app", "err", err)
		}

		oldWorkDir := OldWorkDir("../dwebble_data_work")

		olddb, err := olddb.Open(oldWorkDir.DatabaseFile())
		if err != nil {
			log.Fatal("Failed to connect to old database", "err", err)
		}

		ctx := context.TODO()

		{
			artists, err := olddb.GetAllArtists(ctx, "", "")
			if err != nil {
				log.Fatal("Failed to get artists", "err", err)
			}

			pretty.Println(artists)

			for _, artist := range artists {
				src := oldWorkDir.Artist(artist.Id)
				dst := app.WorkDir().Artist(artist.Id)
				cmd := exec.Command("cp", "-r", src, dst)
				err := cmd.Run()
				if err != nil {
					log.Fatal("Failed to copy artist dir", "err", err)
				}

				_, err = app.DB().CreateArtist(ctx, database.CreateArtistParams{
					Id:        artist.Id,
					Name:      artist.Name,
					OtherName: artist.OtherName,
					Picture:   artist.Picture,
					Created:   artist.Created,
					Updated:   artist.Updated,
				})
				if err != nil {
					log.Fatal("Failed to create artist", "err", err, "artist", artist)
				}

				tags := utils.SplitString(artist.Tags.String)
				for _, tag := range tags {
					slug := utils.Slug(tag)

					err := app.DB().CreateTag(ctx, slug)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						log.Fatal("Failed to create tag", "err", err, "tag", tag, "artist", artist)
					}

					err = app.DB().AddTagToArtist(ctx, slug, artist.Id)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						log.Fatal("Failed to add tag to artist", "err", err, "tag", tag, "artist", artist)
					}
				}
			}
		}

		{
			albums, err := olddb.GetAllAlbums(ctx, "", "")
			if err != nil {
				log.Fatal("Failed to get albums", "err", err)
			}

			pretty.Println(albums)

			for _, album := range albums {
				src := oldWorkDir.Album(album.Id)
				dst := app.WorkDir().Album(album.Id)
				cmd := exec.Command("cp", "-r", src, dst)
				err := cmd.Run()
				if err != nil {
					log.Fatal("Failed to copy album dir", "err", err)
				}

				_, err = app.DB().CreateAlbum(ctx, database.CreateAlbumParams{
					Id:        album.Id,
					Name:      album.Name,
					OtherName: album.OtherName,
					ArtistId:  album.ArtistId,
					CoverArt:  album.CoverArt,
					Year:      album.Year,
					Created:   album.Created,
					Updated:   album.Updated,
				})
				if err != nil {
					log.Fatal("Failed to create album", "err", err, "album", album)
				}

				tags := utils.SplitString(album.Tags.String)
				for _, tag := range tags {
					slug := utils.Slug(tag)

					err := app.DB().CreateTag(ctx, slug)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						log.Fatal("Failed to create tag", "err", err, "tag", tag, "album", album)
					}

					err = app.DB().AddTagToAlbum(ctx, slug, album.Id)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						log.Fatal("Failed to add tag to album", "err", err, "tag", tag, "album", album)
					}
				}

				for _, artist := range album.FeaturingArtists {
					err := app.DB().AddFeaturingArtistToAlbum(ctx, album.Id, artist.Id)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						log.Fatal("Failed to add featuring artist to album", "err", err, "artist", artist, "album", album)
					}
				}
			}
		}

		{
			tracks, err := olddb.GetAllTracks(ctx, "", "")
			if err != nil {
				log.Fatal("Failed to get tracks", "err", err)
			}

			pretty.Println(tracks)

			for _, track := range tracks {
				// src := oldWorkDir.Album(album.Id)
				// dst := app.WorkDir().Album(album.Id)
				// cmd := exec.Command("cp", "-r", src, dst)
				// err := cmd.Run()
				// if err != nil {
				// 	log.Fatal("Failed to copy album dir", "err", err)
				// }

				trackDir := app.WorkDir().Track(track.Id)

				dirs := []string{
					trackDir.String(),
					trackDir.Media(),
				}

				for _, dir := range dirs {
					err := os.Mkdir(dir, 0755)
					if err != nil && !os.IsExist(err) {
						log.Fatal("Failed to create track dir", "err", err, "track", track)
					}
				}

				_, err = app.DB().CreateTrack(ctx, database.CreateTrackParams{
					Id:        track.Id,
					Name:      track.Name,
					OtherName: track.OtherName,
					AlbumId:   track.AlbumId,
					ArtistId:  track.ArtistId,
					Duration:  track.Duration,
					Number:    track.Number,
					Year:      track.Year,
					Created:   track.Created,
					Updated:   track.Updated,
				})
				if err != nil {
					log.Fatal("Failed to create track", "err", err, "track", track)
				}

				tags := utils.SplitString(track.Tags.String)
				for _, tag := range tags {
					slug := utils.Slug(tag)

					err := app.DB().CreateTag(ctx, slug)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						log.Fatal("Failed to create tag", "err", err, "tag", tag, "track", track)
					}

					err = app.DB().AddTagToTrack(ctx, slug, track.Id)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						log.Fatal("Failed to add tag to track", "err", err, "tag", tag, "track", track)
					}
				}

				for _, artist := range track.FeaturingArtists {
					err := app.DB().AddFeaturingArtistToTrack(ctx, track.Id, artist.Id)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						log.Fatal("Failed to add featuring artist to track", "err", err, "artist", artist, "track", track)
					}
				}

				// TODO(patrik): Replace with utils.CreateTrackMediaId()
				mediaId := utils.CreateId()
				mediaDir := trackDir.MediaItem(mediaId)

				{
					dirs := []string{
						mediaDir,
					}

					for _, dir := range dirs {
						err := os.Mkdir(dir, 0755)
						if err != nil && !os.IsExist(err) {
							log.Fatal("Failed to create track media directory", "err", err, "dir", mediaDir, "track", track)
						}
					}
				}

				p := path.Join(oldWorkDir.Track(track.Id), track.OriginalFilename)

				original, err := utils.ProcessOriginalVersion(p, mediaDir, "track")
				if err != nil {
					log.Fatal("Failed to process media", "err", err, "track", track)
				}

				var mediaType types.MediaType

				res, err := ffprobe.ProbeURL(ctx, p)
				if err != nil {
					log.Fatal("Failed to probe media", "err", err, "path", p)
				}
				_ = res

				// TODO(patrik): Check for nil
				stream := res.FirstAudioStream()

				// TODO(patrik): Move to helper
				switch res.Format.FormatName {
				case "flac":
					mediaType = types.MediaTypeFlac
				case "ogg":
					switch stream.CodecName {
					case "opus":
						mediaType = types.MediaTypeOggOpus
					case "vorbis":
						mediaType = types.MediaTypeOggVorbis
					}
				default:
					log.Fatal("Unknown media type", "format", res.Format)
				}

				log.Info("Found media type", "type", mediaType)

				err = app.DB().CreateTrackMedia(ctx, database.CreateTrackMediaParams{
					Id:         mediaId,
					TrackId:    track.Id,
					Filename:   original,
					MediaType:  mediaType,
					IsOriginal: true,
				})
				if err != nil {
					log.Fatal("Failed to create track media", "err", err, "track", track, "path", p)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
