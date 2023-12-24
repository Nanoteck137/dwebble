-- +goose Up
CREATE TABLE artists (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    picture TEXT
);

CREATE TABLE albums (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    cover_art TEXT,

    artist_id TEXT NOT NULL REFERENCES artists(id)
);

CREATE TABLE tracks (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    cover_art TEXT,

    album_id TEXT NOT NULL REFERENCES albums(id),
    artist_id TEXT NOT NULL REFERENCES artists(id)
);

-- +goose Down
DROP TABLE tracks;
DROP TABLE albums;
DROP TABLE artists;
