-- name: GetAllAlbums :many
SELECT * FROM albums;

-- name: GetAlbum :one
SELECT * FROM albums WHERE id=$1;

-- name: GetAlbumsByArtist :many
SELECT * FROM albums WHERE artist_id=$1;

-- name: CreateAlbum :exec
INSERT INTO albums (id, name, cover_art, artist_id) VALUES ($1, $2, $3, $4);

-- name: DeleteAllAlbums :exec
DELETE FROM albums;
