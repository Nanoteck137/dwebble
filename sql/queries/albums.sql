-- name: GetAllAlbums :many
SELECT * FROM albums;

-- name: GetAlbum :one
SELECT * FROM albums WHERE id=$1;

-- name: GetAlbumsByArtist :many
SELECT * FROM albums WHERE artist_id=$1;

-- name: GetAlbumByPath :one
SELECT * FROM albums WHERE path=$1;

-- name: GetAlbumsByArtistAndName :many
SELECT * FROM albums WHERE artist_id=$1 AND name LIKE $2;

-- name: CreateAlbum :one
INSERT INTO albums (id, name, path, cover_art, artist_id) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: DeleteAllAlbums :exec
DELETE FROM albums;
