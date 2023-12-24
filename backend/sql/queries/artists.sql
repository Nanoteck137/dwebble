-- name: GetAllArtists :many
SELECT * FROM artists ORDER BY name;

-- name: GetArtist :one
SELECT * FROM artists WHERE id=$1;

-- name: CreateArtist :exec
INSERT INTO artists (id, name, picture) VALUES ($1, $2, $3);

-- name: DeleteAllArtists :exec
DELETE FROM artists;
