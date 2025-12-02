package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

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

func init() {
	initCmd.Flags().String("dir", ".", "input directory")
	initCmd.Flags().StringP("output", "o", "album.toml", "write result to file")

	rootCmd.AddCommand(initCmd)
}
