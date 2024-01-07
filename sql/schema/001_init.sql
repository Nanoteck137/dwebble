-- +goose Up
CREATE TABLE artists (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    picture TEXT NOT NULL
);

CREATE TABLE albums (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    cover_art TEXT NOT NULL,

    artist_id TEXT NOT NULL REFERENCES artists(id)
);

CREATE TABLE tracks (
    id TEXT PRIMARY KEY,
    track_number INT NOT NULL,
    name TEXT NOT NULL,
    cover_art TEXT NOT NULL,

    filename TEXT NOT NULL,

    album_id TEXT NOT NULL REFERENCES albums(id),
    artist_id TEXT NOT NULL REFERENCES artists(id),

    UNIQUE(album_id, track_number)
);

-- +goose Down
DROP TABLE tracks;
DROP TABLE albums;
DROP TABLE artists;
