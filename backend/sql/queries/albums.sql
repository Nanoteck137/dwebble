-- name: GetAllAlbums :many
SELECT * FROM albums;

-- name: GetAlbum :one
SELECT * FROM albums WHERE id=$1;
