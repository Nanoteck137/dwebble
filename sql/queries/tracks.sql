-- name: GetAllTracks :many
SELECT
    tracks.*, 
    albums.name as album_name,
    artists.name as artist_name
FROM tracks
LEFT JOIN albums ON albums.id = tracks.album_id 
LEFT JOIN artists ON artists.id = tracks.artist_id;

-- name: GetTrack :one
SELECT
    tracks.*, 
    albums.name as album_name,
    artists.name as artist_name
FROM tracks 
LEFT JOIN albums ON albums.id = tracks.album_id 
LEFT JOIN artists ON artists.id = tracks.artist_id
WHERE tracks.id=$1;

-- name: GetTracksByAlbum :many
SELECT 
    tracks.*, 
    albums.name as album_name,
    artists.name as artist_name
FROM tracks 
LEFT JOIN albums ON albums.id = tracks.album_id 
LEFT JOIN artists ON artists.id = tracks.artist_id 
WHERE album_id=$1 
ORDER BY track_number;

-- name: CreateTrack :exec
INSERT INTO tracks (id, track_number, name, cover_art, filename, album_id, artist_id) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: DeleteAllTracks :exec
DELETE FROM tracks;
