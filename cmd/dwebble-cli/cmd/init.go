package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/doug-martin/goqu/v9"
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/cmd/dwebble-cli/api"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var dateRegex = regexp.MustCompile(`^([12]\d\d\d)`)

func parseArtist(s string) []string {
	if s == "" {
		return []string{}
	}

	splits := strings.Split(s, ",")

	artists := make([]string, 0, len(splits))
	for _, s := range splits {
		a := strings.TrimSpace(s)

		if a != "" {
			artists = append(artists, a)
		}
	}

	return artists
}

type TrackInfo struct {
	Name   string
	Artist string
	Number int
	Year   int
}

func getTrackInfo(p string) (TrackInfo, error) {
	probe, err := utils.ProbeTrack(p)
	if err != nil {
		return TrackInfo{}, err
	}

	var res TrackInfo

	res.Name, _ = probe.Tags.GetString("title")
	res.Artist, _ = probe.Tags.GetString("artist")

	if tag, err := probe.Tags.GetString("date"); err == nil {
		match := dateRegex.FindStringSubmatch(tag)
		if len(match) > 0 {
			res.Year, _ = strconv.Atoi(match[1])
		}
	}

	if tag, err := probe.Tags.GetInt("track"); err == nil {
		res.Number = int(tag)
	}

	return res, nil
}

var initCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		output, _ := cmd.Flags().GetString("output")

		// TODO(patrik): Features

		metadata := library.Metadata{}

		extract := true

		// TODO(patrik): Discard hidden files (starts with .)
		entries, err := os.ReadDir(dir)
		if err != nil {
			slog.Error("Failed to read dir", "err", err)
			os.Exit(-1)
		}

		var tracks []string
		var images []string

		for _, e := range entries {
			if e.IsDir() {
				continue
			}

			p := path.Join(dir, e.Name())

			ext := path.Ext(p)

			if utils.IsValidTrackExt(ext) {
				tracks = append(tracks, p)
			}

			if utils.IsValidImageExt(ext) {
				images = append(images, p)
			}
		}

		if len(tracks) <= 0 {
			slog.Warn("No tracks found... Quitting")
			return
		}

		p := tracks[0]

		probe, err := utils.ProbeTrack(p)
		if err != nil {
			slog.Error("Failed to probe track", "err", err)
			os.Exit(-1)
		}

		isSingle := len(tracks) == 1

		metadata.Album.Id = utils.CreateAlbumId()

		if len(images) > 0 {
			// TODO(patrik): Better selection?
			metadata.General.Cover = images[0]
		}

		if !isSingle {
			metadata.Album.Name, _ = probe.Tags.GetString("album")
		} else {
			// NOTE(patrik): If we only have one track then we make the
			// album name the same as the track name
			metadata.Album.Name, _ = probe.Tags.GetString("title")
		}

		if !isSingle {
			if tag, err := probe.Tags.GetString("album_artist"); err == nil {
				metadata.Album.Artists = parseArtist(tag)
			} else {
				if tag, err := probe.Tags.GetString("artist"); err == nil {
					metadata.Album.Artists = parseArtist(tag)
				}
			}
		} else {
			if tag, err := probe.Tags.GetString("artist"); err == nil {
				metadata.Album.Artists = parseArtist(tag)
			}
		}

		if tag, err := probe.Tags.GetString("date"); err == nil {
			match := dateRegex.FindStringSubmatch(tag)
			if len(match) > 0 {
				metadata.General.Year, _ = strconv.ParseInt(match[1], 10, 64)
			}
		}

		for _, p := range tracks {
			filename := path.Base(p)
			fmt.Printf("Found track: %s\n", filename)

			trackInfo, err := getTrackInfo(p)
			if err != nil {
				slog.Error("Failed to get track info", "err", err)
				os.Exit(-1)
			}

			if trackInfo.Name == "" {
				trackInfo.Name = strings.TrimSuffix(filename, path.Ext(p))
			}

			if trackInfo.Number == 0 || extract {
				trackInfo.Number = utils.ExtractNumber(filename)
			}

			// TODO(patrik): If artist is empty then use album maybe
			artists := parseArtist(trackInfo.Artist)

			metadata.Tracks = append(metadata.Tracks, library.MetadataTrack{
				Id:      utils.CreateTrackId(),
				File:    filename,
				Name:    trackInfo.Name,
				Number:  int64(trackInfo.Number),
				Year:    0,
				Tags:    []string{},
				Artists: artists,
			})
		}

		data, err := toml.Marshal(&metadata)
		if err != nil {
			slog.Error("Failed to marshal metadata", "err", err)
			os.Exit(-1)
		}

		err = os.WriteFile(output, data, 0644)
		if err != nil {
			slog.Error("Failed to write output", "err", err)
			os.Exit(-1)
		}
	},
}

type Album struct {
	Id string `db:"id"`

	Name string `db:"name"`

	ArtistId string `db:"artist_id"`

	CoverArt sql.NullString `db:"cover_art"`
	Year     sql.NullInt64  `db:"year"`

	ArtistName      string         `db:"artist_name"`
	ArtistOtherName sql.NullString `db:"artist_other_name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`

	FeaturingArtists database.FeaturingArtists `db:"featuring_artists"`
}

var dialect = goqu.Dialect("sqlite_returning")

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

func AlbumQuery() *goqu.SelectDataset {
	tags := dialect.From("albums_tags").
		Select(
			goqu.I("albums_tags.album_id").As("album_id"),
			goqu.Func("group_concat", goqu.I("tags.slug"), ",").As("tags"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(goqu.I("albums_tags.tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(goqu.I("albums_tags.album_id"))

	query := dialect.From("albums").
		Select(
			"albums.id",

			"albums.name",

			"albums.artist_id",

			"albums.cover_art",
			"albums.year",

			"albums.created",
			"albums.updated",

			goqu.I("artists.name").As("artist_name"),
			goqu.I("artists.other_name").As("artist_other_name"),

			goqu.I("tags.tags").As("tags"),

			goqu.I("featuring_artists.artists").As("featuring_artists"),
		).
		Prepared(true).
		Join(
			goqu.I("artists"),
			goqu.On(
				goqu.I("albums.artist_id").Eq(goqu.I("artists.id")),
			),
		).
		LeftJoin(
			tags.As("tags"),
			goqu.On(goqu.I("albums.id").Eq(goqu.I("tags.album_id"))),
		).
		LeftJoin(
			FeaturingArtistsQuery("albums_featuring_artists", "album_id").As("featuring_artists"),
			goqu.On(goqu.I("albums.id").Eq(goqu.I("featuring_artists.id"))),
		)

	return query
}

type Track struct {
	Id string `db:"id"`

	Name string `db:"name"`

	AlbumId  string `db:"album_id"`
	ArtistId string `db:"artist_id"`

	Number sql.NullInt64 `db:"number"`
	Year   sql.NullInt64 `db:"year"`

	AlbumName string `db:"album_name"`

	ArtistName string `db:"artist_name"`

	Tags sql.NullString `db:"tags"`

	FeaturingArtists database.FeaturingArtists `db:"featuring_artists"`
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
			"tracks.id",

			"tracks.name",

			"tracks.album_id",
			"tracks.artist_id",

			"tracks.number",
			"tracks.year",

			goqu.I("albums.name").As("album_name"),

			goqu.I("artists.name").As("artist_name"),

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

var oldInitCmd = &cobra.Command{
	Use: "old-init",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		output, _ := cmd.Flags().GetString("output")
		dbPath, _ := cmd.Flags().GetString("db")

		fmt.Printf("dbPath: %v\n", dbPath)

		db, err := database.OpenRaw(dbPath)
		if err != nil {
			slog.Error("Failed to open db", "err", err)
			os.Exit(-1)
		}

		query := AlbumQuery()

		var albums []Album
		err = db.Select(&albums, query)
		if err != nil {
			slog.Error("Failed to get albums", "err", err)
			os.Exit(-1)
		}

		// pretty.Println(albums)

		options := make([]huh.Option[string], len(albums))

		for i, a := range albums {
			name := fmt.Sprintf("%s: %s", a.ArtistName, a.Name)
			options[i] = huh.NewOption[string](name, a.Id)
		}

		var value string
		err = huh.NewSelect[string]().
			Title("Pick album").
			Options(options...).
			Value(&value).
			WithHeight(10).
			Run()
		if err != nil {
			slog.Error("Failed", "err", err)
			os.Exit(-1)
		}

		{
			query = AlbumQuery().Where(goqu.I("albums.id").Eq(value))
			var album Album
			err = db.Get(&album, query)
			if err != nil {
				slog.Error("Failed to get album", "err", err)
				os.Exit(-1)
			}

			pretty.Println(album)

			query = TrackQuery().Where(goqu.I("tracks.album_id").Eq(album.Id))
			var dbTracks []Track
			err = db.Select(&dbTracks, query)
			if err != nil {
				slog.Error("Failed to get album tracks", "err", err)
				os.Exit(-1)
			}

			pretty.Println(dbTracks)

			metadata := library.Metadata{}

			metadata.Album.Id = album.Id
			metadata.Album.Name = album.Name
			metadata.Album.Year = album.Year.Int64
			metadata.Album.Tags = strings.Split(album.Tags.String, ",")
			metadata.Album.Artists = []string{album.ArtistName}

			for _, v := range album.FeaturingArtists {
				metadata.Album.Artists = append(metadata.Album.Artists, v.Name)
			}

			// TODO(patrik): Discard hidden files (starts with .)
			entries, err := os.ReadDir(dir)
			if err != nil {
				slog.Error("Failed to read dir", "err", err)
				os.Exit(-1)
			}

			var tracks []string
			var images []string

			for _, e := range entries {
				if e.IsDir() {
					continue
				}

				p := path.Join(dir, e.Name())

				ext := path.Ext(p)

				if utils.IsValidTrackExt(ext) {
					tracks = append(tracks, p)
				}

				if utils.IsValidImageExt(ext) {
					images = append(images, p)
				}
			}

			if len(tracks) <= 0 {
				slog.Warn("No tracks found... Quitting")
				return
			}

			if len(images) > 0 {
				// TODO(patrik): Better selection?
				metadata.General.Cover = images[0]
			}

			for _, p := range tracks {
				filename := path.Base(p)
				num := utils.ExtractNumber(filename)

				var track *Track

				for _, t := range dbTracks {
					if t.Number.Valid && t.Number.Int64 == int64(num) {
						track = &t
						break
					}
				}

				if track == nil {
					slog.Warn("Failed to find match for track", "track", p)
					continue
				}

				res := library.MetadataTrack{
					Id:      track.Id,
					File:    filename,
					Name:    track.Name,
					Number:  track.Number.Int64,
					Year:    track.Year.Int64,
					Tags:    strings.Split(track.Tags.String, ","),
					Artists: []string{track.ArtistName},
				}

				for _, a := range track.FeaturingArtists {
					res.Artists = append(res.Artists, a.Name)
				}

				metadata.Tracks = append(metadata.Tracks, res)
			}

			pretty.Println(metadata)

			data, err := toml.Marshal(&metadata)
			if err != nil {
				slog.Error("Failed to marshal metadata", "err", err)
				os.Exit(-1)
			}

			err = os.WriteFile(output, data, 0644)
			if err != nil {
				slog.Error("Failed to write output", "err", err)
				os.Exit(-1)
			}
		}
	},
}

type Playlist struct {
	Id      string         `db:"id"`
	Name    string         `db:"name"`
	Picture sql.NullString `db:"picture"`

	OwnerId string `db:"owner_id"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

type PlaylistItem struct {
	PlaylistId string `db:"playlist_id"`
	TrackId    string `db:"track_id"`
}

func PlaylistQuery() *goqu.SelectDataset {
	query := dialect.From("playlists").
		Select(
			"playlists.id",
			"playlists.name",
			"playlists.picture",

			"playlists.owner_id",

			"playlists.created",
			"playlists.updated",
		).
		Prepared(true)

	return query
}

type ExportedPlaylistItem struct {
	TrackId    string `json:"trackId"`
	TrackName  string `json:"trackName"`
	AlbumName  string `json:"albumName"`
	ArtistName string `json:"artistName"`
}

type ExportedPlaylist struct {
	Name  string                 `json:"name"`
	Items []ExportedPlaylistItem `json:"items"`
}

var exportPlaylistsCmd = &cobra.Command{
	Use: "export-playlists",
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		dbPath, _ := cmd.Flags().GetString("db")

		_ = output

		db, err := database.OpenRaw(dbPath)
		if err != nil {
			slog.Error("Failed to open db", "err", err)
			os.Exit(-1)
		}

		query := PlaylistQuery()

		var playlists []Playlist
		err = db.Select(&playlists, query)
		if err != nil {
			slog.Error("Failed to get playlists", "err", err)
			os.Exit(-1)
		}

		pretty.Println(playlists)

		for _, playlist := range playlists {
			query := dialect.From("playlist_items").
				Select(
					"playlist_items.track_id",
				).
				Where(goqu.I("playlist_id").Eq(playlist.Id))

			var items []string
			err := db.Select(&items, query)
			if err != nil {
				slog.Error("Failed to get playlist items", "err", err, "playlistId", playlist.Id)
				os.Exit(-1)
			}

			res := ExportedPlaylist{
				Name:  playlist.Name,
				Items: make([]ExportedPlaylistItem, len(items)),
			}

			for i, trackId := range items {
				query := TrackQuery().Where(goqu.I("tracks.id").Eq(trackId))

				var track Track
				err = db.Get(&track, query)
				if err != nil {
					slog.Error("Failed to get track for playlist", "err", err, "playlistId", playlist.Id, "trackId", trackId)
					os.Exit(-1)
				}

				res.Items[i] = ExportedPlaylistItem{
					TrackId:    trackId,
					TrackName:  track.Name,
					AlbumName:  track.AlbumName,
					ArtistName: track.ArtistName,
				}
			}

			pretty.Println(res)

			d, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				slog.Error("Failed to marshal json", "err", err, "playlistId", playlist.Id)
				os.Exit(-1)
			}

			name := utils.Slug(playlist.Name) + ".json"
			err = os.WriteFile(path.Join(output, name), d, 0644)
			if err != nil {
				slog.Error("Failed to write json", "err", err, "playlistId", playlist.Id)
				os.Exit(-1)
			}
		}
	},
}

var importPlaylistCmd = &cobra.Command{
	Use:  "import-playlist <PLAYLIST_FILE>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		server, _ := cmd.Flags().GetString("server")
		token, _ := cmd.Flags().GetString("token")

		client := api.New(server)
		if token != "" {
			client.SetApiToken(token)
		}

		d, err := os.ReadFile(file)
		if err != nil {
			slog.Error("Failed", "err", err)
			os.Exit(-1)
		}

		var p ExportedPlaylist
		err = json.Unmarshal(d, &p)
		if err != nil {
			slog.Error("Failed", "err", err)
			os.Exit(-1)
		}

		pretty.Println(p)

		playlists, err := client.GetPlaylists(api.Options{})
		if err != nil {
			slog.Error("Failed", "err", err)
			os.Exit(-1)
		}

		pretty.Println(playlists)

		var playlistId *string

		for _, playlist := range playlists.Playlists {
			if playlist.Name == p.Name {
				playlistId = &playlist.Id
				break
			}
		}

		if playlistId == nil {
			slog.Info("Failed to find playlist creating new playlist", "name", p.Name)

			res, err := client.CreatePlaylist(api.CreatePlaylistBody{
				Name: p.Name,
			}, api.Options{})
			if err != nil {
				slog.Error("Failed", "err", err)
				os.Exit(-1)
			}

			playlistId = &res.Id
		} else {
			_, err = client.ClearPlaylist(*playlistId, api.Options{})
			if err != nil {
				slog.Error("Failed", "err", err)
				os.Exit(-1)
			}
		}

		for _, item := range p.Items {
			_, err := client.AddItemToPlaylist(*playlistId, api.AddItemToPlaylistBody{
				TrackId: item.TrackId,
			}, api.Options{})
			if err != nil {
				slog.Error("Failed to add track to playlist", "err", err, "track", item)
				os.Exit(-1)
			}
		}
	},
}

func init() {
	initCmd.Flags().String("dir", ".", "input directory")
	initCmd.Flags().StringP("output", "o", "album.toml", "write result to file")

	oldInitCmd.Flags().String("dir", ".", "input directory")
	oldInitCmd.Flags().StringP("output", "o", "album.toml", "write result to file")
	oldInitCmd.Flags().String("db", "", "database to use")
	oldInitCmd.MarkFlagRequired("db")

	exportPlaylistsCmd.Flags().StringP("output", "o", ".", "write result to file")
	exportPlaylistsCmd.Flags().String("db", "", "database to use")
	exportPlaylistsCmd.MarkFlagRequired("db")

	importPlaylistCmd.Flags().StringP("token", "t", "", "Api Token to use")
	importPlaylistCmd.Flags().String("server", "", "Server Address")
	importPlaylistCmd.MarkFlagRequired("server")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(oldInitCmd)
	rootCmd.AddCommand(exportPlaylistsCmd)
	rootCmd.AddCommand(importPlaylistCmd)
}
