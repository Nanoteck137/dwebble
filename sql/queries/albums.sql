-- name: GetAllAlbums :many
SELECT * FROM albums;

-- name: GetAlbum :one
SELECT * FROM albums WHERE id=$1;

-- name: GetAlbumsByArtist :many
SELECT * FROM albums WHERE artist_id=$1;

-- name: GetAlbumsByArtistAndName :many
SELECT * FROM albums WHERE artist_id=$1 AND name LIKE $2;

-- name: CreateAlbum :one
INSERT INTO albums (id, name, cover_art, artist_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: DeleteAllAlbums :exec
DELETE FROM albums;
