-- name: GetAllSongs :many
SELECT * FROM songs;

-- name: GetSong :one
SELECT * FROM songs WHERE id=$1;

-- name: GetSongsByAlbum :many
SELECT * FROM songs WHERE album_id=$1;
