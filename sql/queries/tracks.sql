-- name: GetAllTracks :many
SELECT
    tracks.*, 
    albums.name as album_name,
    artists.name as artist_name
FROM tracks
JOIN albums ON albums.id = tracks.album_id 
JOIN artists ON artists.id = tracks.artist_id;

-- name: GetTrack :one
SELECT
    tracks.*, 
    albums.name as album_name,
    artists.name as artist_name
FROM tracks 
JOIN albums ON albums.id = tracks.album_id 
JOIN artists ON artists.id = tracks.artist_id
WHERE tracks.id=$1;

-- name: GetTracksByAlbum :many
SELECT 
    tracks.*, 
    albums.name as album_name,
    artists.name as artist_name
FROM tracks 
JOIN albums ON albums.id = tracks.album_id 
JOIN artists ON artists.id = tracks.artist_id 
WHERE album_id=$1 
ORDER BY track_number;

-- name: CreateTrack :one
INSERT INTO tracks (id, track_number, name, cover_art, best_quality_file, mobile_quality_file, album_id, artist_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: DeleteAllTracks :exec
DELETE FROM tracks;
