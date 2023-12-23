-- name: GetAllArtists :many
SELECT * FROM artists;

-- name: GetArtist :one
SELECT * FROM artists WHERE id=$1;
