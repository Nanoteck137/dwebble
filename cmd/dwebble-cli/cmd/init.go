package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/doug-martin/goqu/v9"
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

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
			log.Fatal("Failed to read dir", "err", err)
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
			log.Warn("No tracks found... Quitting")
			return
		}

		p := tracks[0]

		probe, err := utils.ProbeTrack(p)
		if err != nil {
			log.Fatal("Failed to probe track", "err", err)
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
				log.Fatal("Failed to get track info", "err", err)
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
			log.Fatal("Failed to marshal metadata", "err", err)
		}

		err = os.WriteFile(output, data, 0644)
		if err != nil {
			log.Fatal("Failed to write output", "err", err)
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

		db, err := database.OpenRaw("../data.db")
		if err != nil {
			log.Fatal("Failed to open db", "err", err)
		}

		query := AlbumQuery()

		var albums []Album
		err = db.Select(&albums, query)
		if err != nil {
			log.Fatal("Failed to get albums", "err", err)
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
			log.Fatal("Failed", "err", err)
		}

		{
			query = AlbumQuery().Where(goqu.I("albums.id").Eq(value))
			var album Album
			err = db.Get(&album, query)
			if err != nil {
				log.Fatal("Failed to get album", "err", err)
			}

			pretty.Println(album)

			query = TrackQuery().Where(goqu.I("tracks.album_id").Eq(album.Id))
			var dbTracks []Track
			err = db.Select(&dbTracks, query)
			if err != nil {
				log.Fatal("Failed to get album tracks", "err", err)
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
				log.Fatal("Failed to read dir", "err", err)
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
				log.Warn("No tracks found... Quitting")
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
				log.Fatal("Failed to marshal metadata", "err", err)
			}

			err = os.WriteFile(output, data, 0644)
			if err != nil {
				log.Fatal("Failed to write output", "err", err)
			}
		}
	},
}

func init() {
	initCmd.Flags().String("dir", ".", "input directory")
	initCmd.Flags().StringP("output", "o", "album.toml", "write result to file")

	oldInitCmd.Flags().String("dir", ".", "input directory")
	oldInitCmd.Flags().StringP("output", "o", "album.toml", "write result to file")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(oldInitCmd)
}
