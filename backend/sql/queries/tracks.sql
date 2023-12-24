-- name: GetAllTracks :many
SELECT * FROM tracks;

-- name: GetTrack :one
SELECT * FROM tracks WHERE id=$1;

-- name: GetTracksByAlbum :many
SELECT * FROM tracks WHERE album_id=$1;

-- name: CreateTrack :exec
INSERT INTO tracks (id, name, cover_art, album_id, artist_id) VALUES ($1, $2, $3, $4, $5);

-- name: DeleteAllTracks :exec
DELETE FROM tracks;
