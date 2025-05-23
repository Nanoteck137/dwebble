package cmd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path"

	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/spf13/cobra"

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

func (d OldWorkDir) Track(id string) TrackDir {
	return TrackDir(path.Join(d.Tracks(), id))
}

type TrackDir string

func (d TrackDir) String() string {
	return string(d)
}

func (d TrackDir) Media() string {
	return path.Join(d.String(), "media")
}

func (d TrackDir) MediaItem(id string) string {
	return path.Join(d.Media(), id)
}

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")

		err := os.MkdirAll(output, 0755)
		if err != nil {
			log.Fatal("Failed to create output", "err", err, "output", output)
		}

		app := core.NewBaseApp(&config.Config{
			RunMigrations: true,
			DataDir:       output,
		})

		err = app.Bootstrap()
		if err != nil {
			log.Fatal("Failed to bootstrap app", "err", err)
		}

		oldWorkDir := OldWorkDir(input)

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

			for _, track := range tracks {
				trackDir := app.WorkDir().Track(track.Id)

				dirs := []string{
					trackDir,
				}

				for _, dir := range dirs {
					err := os.Mkdir(dir, 0755)
					if err != nil && !os.IsExist(err) {
						log.Fatal("Failed to create track dir", "err", err, "track", track)
					}
				}

				var mediaType types.MediaType
				var filename string

				formats := *track.Formats.Get()
				for _, tf := range formats {
					if tf.IsOriginal {
						dir := oldWorkDir.Track(track.Id).MediaItem(tf.Id)
						ext := path.Ext(tf.Filename)
						p := path.Join(dir, tf.Filename)

						name := "track" + ext
						dst := path.Join(trackDir, name)

						cmd := exec.Command("cp", p, dst)
						err := cmd.Run()
						if err != nil {
							log.Fatal("Failed to copy original track file to new directory", "err", err, "track", p, "new", dst)
						}

						filename = name
						mediaType = tf.MediaType

						break
					}
				}

				_, err = app.DB().CreateTrack(ctx, database.CreateTrackParams{
					Id:        track.Id,
					Filename:  filename,
					MediaType: mediaType,
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
			}

			{
				users, err := olddb.GetAllUsers(ctx)
				if err != nil {
					log.Fatal("Failed to get users", "err", err)
				}

				for _, user := range users {
					newUser, err := app.DB().CreateUser(ctx, database.CreateUserParams{
						Id:       user.Id,
						Username: user.Username,
						Password: user.Password,
						Role:     user.Role,
						Created:  user.Created,
						Updated:  user.Updated,
					})
					if err != nil {
						log.Fatal("Failed to create user", "err", err, "user", user)
					}

					playlists, err := olddb.GetPlaylistsByUser(ctx, user.Id)
					if err != nil {
						log.Fatal("Failed to get user playlists", "err", err, "user", user)
					}

					for _, playlist := range playlists {
						_, err := app.DB().CreatePlaylist(ctx, database.CreatePlaylistParams{
							Id:      playlist.Id,
							Name:    playlist.Name,
							Picture: playlist.Picture,
							OwnerId: playlist.OwnerId,
							Created: playlist.Created,
							Updated: playlist.Updated,
						})
						if err != nil {
							log.Fatal("Failed to create user playlist", "err", err, "user", user, "playlist", playlist)
						}

						items, err := olddb.GetPlaylistItems(ctx, playlist.Id)
						if err != nil {
							log.Fatal("Failed to get playlist items", "err", err, "user", user, "playlist", playlist)
						}

						for _, item := range items {
							err = app.DB().AddItemToPlaylist(ctx, playlist.Id, item.TrackId)
							if err != nil {
								log.Fatal("Failed to add item to playlist", "err", err, "user", user, "playlist", playlist, "item", item)
							}
						}
					}

					taglists, err := olddb.GetTaglistByUser(ctx, user.Id)
					if err != nil {
						log.Fatal("Failed to get user taglists", "err", err, "user", user)
					}

					for _, taglist := range taglists {
						_, err := app.DB().CreateTaglist(ctx, database.CreateTaglistParams{
							Name:    taglist.Name,
							Filter:  taglist.Filter,
							OwnerId: taglist.OwnerId,
							Created: taglist.Created,
							Updated: taglist.Updated,
						})
						if err != nil {
							log.Fatal("Failed to create user taglist", "err", err, "user", user, "taglist", taglist)
						}
					}

					apiTokens, err := olddb.GetAllApiTokensForUser(ctx, user.Id)
					if err != nil {
						log.Fatal("Failed to get user api tokens", "err", err, "user", user)
					}

					for _, token := range apiTokens {
						_, err := app.DB().CreateApiToken(ctx, database.CreateApiTokenParams{
							Id:      token.Id,
							UserId:  token.UserId,
							Name:    token.Name,
							Created: token.Created,
							Updated: token.Updated,
						})
						if err != nil {
							log.Fatal("Failed to create api token", "err", err, "user", user, "token", token)
						}
					}

					// NOTE(patrik): This needs to happen after playlists
					err = app.DB().UpdateUserSettings(ctx, database.UserSettings{
						Id:            newUser.Id,
						DisplayName:   user.DisplayName,
						QuickPlaylist: user.QuickPlaylist,
					})
					if err != nil {
						log.Fatal("Failed to update user settings", "err", err, "user", user)
					}
				}
			}

			f, err := os.Create(app.WorkDir().SetupFile())
			if err != nil {
				log.Fatal("Failed to create setup file", "err", err)
			}
			f.Close()
		}
	},
}

func init() {
	migrateCmd.Flags().String("input", "", "Input Directory")
	migrateCmd.MarkFlagRequired("input")
	migrateCmd.MarkFlagDirname("input")

	migrateCmd.Flags().String("output", "", "Output Directory")
	migrateCmd.MarkFlagRequired("output")
	migrateCmd.MarkFlagDirname("output")

	rootCmd.AddCommand(migrateCmd)
}
