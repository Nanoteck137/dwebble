package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/spf13/cobra"

	_ "github.com/mattn/go-sqlite3"
)

var AppName = dwebble.AppName + "-migrate"

func rowsToArrayMap(rows *sql.Rows) ([]map[string]any, error) {

	cols, _ := rows.Columns()

	var res []map[string]any
	for rows.Next() {
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]any)
		for i, colName := range cols {
			val := columnPointers[i].(*any)
			m[colName] = *val
		}

		res = append(res, m)
	}

	return res, nil
}

var rootCmd = &cobra.Command{
	Use:     AppName + " <IN>",
	Version: dwebble.Version,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		oldDb, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", args[0]))
		if err != nil {
			log.Fatal("Failed", "err", err)
		}

		workDir := types.WorkDir(dir)

		db, err := database.Open(workDir)
		if err != nil {
			log.Fatal("Failed", "err", err)
		}

		err = db.RunMigrateUp()
		if err != nil {
			log.Fatal("Failed", "err", err)
		}

		valueToSqlString := func(v any) sql.NullString {
			if v == nil {
				return sql.NullString{}
			}

			return sql.NullString{
				String: v.(string),
				Valid:  true,
			}
		}

		valueToSqlInt64 := func(v any) sql.NullInt64 {
			if v == nil {
				return sql.NullInt64{}
			}

			return sql.NullInt64{
				Int64: v.(int64),
				Valid: true,
			}
		}

		{
			rows, err := oldDb.Query("SELECT * FROM artists")
			if err != nil {
				log.Fatal("Failed", "err", err)
			}
			defer rows.Close()

			artists, err := rowsToArrayMap(rows)
			if err != nil {
				log.Fatal("Failed", "err", err)
			}

			pretty.Println(artists)

			ctx := context.TODO()

			for _, v := range artists {
				_, err := db.CreateArtist(ctx, database.CreateArtistParams{
					Id:      v["id"].(string),
					Name:    v["name"].(string),
					Picture: valueToSqlString(v["picture"]),
				})
				if err != nil {
					log.Fatal("Failed", "err", err)
				}
			}
		}

		{
			rows, err := oldDb.Query("SELECT * FROM albums")
			if err != nil {
				log.Fatal("Failed", "err", err)
			}
			defer rows.Close()

			albums, err := rowsToArrayMap(rows)
			if err != nil {
				log.Fatal("Failed", "err", err)
			}

			pretty.Println(albums)

			ctx := context.TODO()

			for _, v := range albums {
				_, err := db.CreateAlbum(ctx, database.CreateAlbumParams{
					Id:       v["id"].(string),
					Name:     v["name"].(string),
					ArtistId: v["artist_id"].(string),
					CoverArt: valueToSqlString(v["cover_art"]),
					Year:     valueToSqlInt64(v["year"]),
					Created:  v["created"].(int64),
					Updated:  v["updated"].(int64),
				})
				if err != nil {
					log.Fatal("Failed", "err", err)
				}
			}
		}

		{
			rows, err := oldDb.Query("SELECT * FROM tracks")
			if err != nil {
				log.Fatal("Failed", "err", err)
			}
			defer rows.Close()

			tracks, err := rowsToArrayMap(rows)
			if err != nil {
				log.Fatal("Failed", "err", err)
			}

			pretty.Println(tracks)

			ctx := context.TODO()

			for _, v := range tracks {
				_, err := db.CreateTrack(ctx, database.CreateTrackParams{
					Id:               v["id"].(string),
					Name:             v["name"].(string),
					AlbumId:          v["album_id"].(string),
					ArtistId:         v["artist_id"].(string),
					Duration:         valueToSqlInt64(v["duration"]).Int64,
					Number:           valueToSqlInt64(v["track_number"]),
					Year:             valueToSqlInt64(v["year"]),
					OriginalFilename: v["original_filename"].(string),
					MobileFilename:   v["mobile_filename"].(string),
					Created:          v["created"].(int64),
					Updated:          v["updated"].(int64),
				})
				if err != nil {
					log.Fatal("Failed", "err", err)
				}
			}
		}

		{
			rows, err := oldDb.Query("SELECT * FROM tags")
			if err != nil {
				log.Fatal("Failed", "err", err)
			}
			defer rows.Close()

			tags, err := rowsToArrayMap(rows)
			if err != nil {
				log.Fatal("Failed", "err", err)
			}

			pretty.Println(tags)

			ctx := context.TODO()

			for _, v := range tags {
				err := db.CreateTag(ctx, v["slug"].(string))
				if err != nil {
					log.Fatal("Failed", "err", err)
				}
			}
		}

		{
			rows, err := oldDb.Query("SELECT * FROM tracks_to_tags")
			if err != nil {
				log.Fatal("Failed", "err", err)
			}
			defer rows.Close()

			tags, err := rowsToArrayMap(rows)
			if err != nil {
				log.Fatal("Failed", "err", err)
			}

			pretty.Println(tags)

			ctx := context.TODO()

			for _, v := range tags {
				err := db.AddTagToTrack(ctx, v["tag_slug"].(string), v["track_id"].(string))
				if err != nil {
					log.Fatal("Failed", "err", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.Flags().StringP("dir", "d", ".", "Work Directory")
	rootCmd.SetVersionTemplate(dwebble.VersionTemplate(AppName))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to run root command", "err", err)
	}
}

func main() {
	Execute()
}
