// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: tracks.sql

package database

import (
	"context"
)

const createTrack = `-- name: CreateTrack :one
INSERT INTO tracks (id, track_number, name, cover_art, best_quality_file, mobile_quality_file, album_id, artist_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, track_number, name, cover_art, best_quality_file, mobile_quality_file, album_id, artist_id
`

type CreateTrackParams struct {
	ID                string `json:"id"`
	TrackNumber       int32  `json:"trackNumber"`
	Name              string `json:"name"`
	CoverArt          string `json:"coverArt"`
	BestQualityFile   string `json:"bestQualityFile"`
	MobileQualityFile string `json:"mobileQualityFile"`
	AlbumID           string `json:"albumId"`
	ArtistID          string `json:"artistId"`
}

func (q *Queries) CreateTrack(ctx context.Context, arg CreateTrackParams) (Track, error) {
	row := q.db.QueryRow(ctx, createTrack,
		arg.ID,
		arg.TrackNumber,
		arg.Name,
		arg.CoverArt,
		arg.BestQualityFile,
		arg.MobileQualityFile,
		arg.AlbumID,
		arg.ArtistID,
	)
	var i Track
	err := row.Scan(
		&i.ID,
		&i.TrackNumber,
		&i.Name,
		&i.CoverArt,
		&i.BestQualityFile,
		&i.MobileQualityFile,
		&i.AlbumID,
		&i.ArtistID,
	)
	return i, err
}

const deleteAllTracks = `-- name: DeleteAllTracks :exec
DELETE FROM tracks
`

func (q *Queries) DeleteAllTracks(ctx context.Context) error {
	_, err := q.db.Exec(ctx, deleteAllTracks)
	return err
}

const getAllTracks = `-- name: GetAllTracks :many
SELECT
    tracks.id, tracks.track_number, tracks.name, tracks.cover_art, tracks.best_quality_file, tracks.mobile_quality_file, tracks.album_id, tracks.artist_id, 
    albums.name as album_name,
    artists.name as artist_name
FROM tracks
JOIN albums ON albums.id = tracks.album_id 
JOIN artists ON artists.id = tracks.artist_id
`

type GetAllTracksRow struct {
	ID                string `json:"id"`
	TrackNumber       int32  `json:"trackNumber"`
	Name              string `json:"name"`
	CoverArt          string `json:"coverArt"`
	BestQualityFile   string `json:"bestQualityFile"`
	MobileQualityFile string `json:"mobileQualityFile"`
	AlbumID           string `json:"albumId"`
	ArtistID          string `json:"artistId"`
	AlbumName         string `json:"albumName"`
	ArtistName        string `json:"artistName"`
}

func (q *Queries) GetAllTracks(ctx context.Context) ([]GetAllTracksRow, error) {
	rows, err := q.db.Query(ctx, getAllTracks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllTracksRow{}
	for rows.Next() {
		var i GetAllTracksRow
		if err := rows.Scan(
			&i.ID,
			&i.TrackNumber,
			&i.Name,
			&i.CoverArt,
			&i.BestQualityFile,
			&i.MobileQualityFile,
			&i.AlbumID,
			&i.ArtistID,
			&i.AlbumName,
			&i.ArtistName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTrack = `-- name: GetTrack :one
SELECT
    tracks.id, tracks.track_number, tracks.name, tracks.cover_art, tracks.best_quality_file, tracks.mobile_quality_file, tracks.album_id, tracks.artist_id, 
    albums.name as album_name,
    artists.name as artist_name
FROM tracks 
JOIN albums ON albums.id = tracks.album_id 
JOIN artists ON artists.id = tracks.artist_id
WHERE tracks.id=$1
`

type GetTrackRow struct {
	ID                string `json:"id"`
	TrackNumber       int32  `json:"trackNumber"`
	Name              string `json:"name"`
	CoverArt          string `json:"coverArt"`
	BestQualityFile   string `json:"bestQualityFile"`
	MobileQualityFile string `json:"mobileQualityFile"`
	AlbumID           string `json:"albumId"`
	ArtistID          string `json:"artistId"`
	AlbumName         string `json:"albumName"`
	ArtistName        string `json:"artistName"`
}

func (q *Queries) GetTrack(ctx context.Context, id string) (GetTrackRow, error) {
	row := q.db.QueryRow(ctx, getTrack, id)
	var i GetTrackRow
	err := row.Scan(
		&i.ID,
		&i.TrackNumber,
		&i.Name,
		&i.CoverArt,
		&i.BestQualityFile,
		&i.MobileQualityFile,
		&i.AlbumID,
		&i.ArtistID,
		&i.AlbumName,
		&i.ArtistName,
	)
	return i, err
}

const getTracksByAlbum = `-- name: GetTracksByAlbum :many
SELECT 
    tracks.id, tracks.track_number, tracks.name, tracks.cover_art, tracks.best_quality_file, tracks.mobile_quality_file, tracks.album_id, tracks.artist_id, 
    albums.name as album_name,
    artists.name as artist_name
FROM tracks 
JOIN albums ON albums.id = tracks.album_id 
JOIN artists ON artists.id = tracks.artist_id 
WHERE album_id=$1 
ORDER BY track_number
`

type GetTracksByAlbumRow struct {
	ID                string `json:"id"`
	TrackNumber       int32  `json:"trackNumber"`
	Name              string `json:"name"`
	CoverArt          string `json:"coverArt"`
	BestQualityFile   string `json:"bestQualityFile"`
	MobileQualityFile string `json:"mobileQualityFile"`
	AlbumID           string `json:"albumId"`
	ArtistID          string `json:"artistId"`
	AlbumName         string `json:"albumName"`
	ArtistName        string `json:"artistName"`
}

func (q *Queries) GetTracksByAlbum(ctx context.Context, albumID string) ([]GetTracksByAlbumRow, error) {
	rows, err := q.db.Query(ctx, getTracksByAlbum, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTracksByAlbumRow{}
	for rows.Next() {
		var i GetTracksByAlbumRow
		if err := rows.Scan(
			&i.ID,
			&i.TrackNumber,
			&i.Name,
			&i.CoverArt,
			&i.BestQualityFile,
			&i.MobileQualityFile,
			&i.AlbumID,
			&i.ArtistID,
			&i.AlbumName,
			&i.ArtistName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
