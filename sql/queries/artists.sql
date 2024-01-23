-- name: GetAllArtists :many
SELECT * FROM artists ORDER BY name;

-- name: GetArtist :one
SELECT * FROM artists WHERE id=$1;

-- name: GetArtistByName :many
SELECT * FROM artists WHERE name LIKE $1;

-- name: CreateArtist :one
INSERT INTO artists (id, name, picture) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteAllArtists :exec
DELETE FROM artists;
