-- name: GetAllAlbums :many
SELECT * FROM albums;

-- name: GetAlbum :one
SELECT * FROM albums WHERE id=$1;

-- name: GetAlbumsByArtist :many
SELECT * FROM albums WHERE artist_id=$1;
