-- name: GetAllSongs :many
SELECT * FROM songs;

-- name: GetSong :one
SELECT * FROM songs WHERE id=$1;

-- name: GetSongsByAlbum :many
SELECT * FROM songs WHERE album_id=$1;

-- name: CreateSong :exec
INSERT INTO songs (id, name, cover_art, album_id, artist_id) VALUES ($1, $2, $3, $4, $5);

-- name: DeleteAllSongs :exec
DELETE FROM songs;
