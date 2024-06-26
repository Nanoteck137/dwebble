-- NOTE(patrik): Final

-- +goose Up
CREATE TABLE artists (
    id TEXT,
    name TEXT NOT NULL,
    picture TEXT,

    path TEXT NOT NULL,
    available BOOL NOT NULL,

    CONSTRAINT artists_pk PRIMARY KEY (id)
);

CREATE TABLE albums (
    id TEXT,
    name TEXT NOT NULL,
    cover_art TEXT,

    artist_id TEXT NOT NULL,

    path TEXT NOT NULL,
    available BOOL NOT NULL,

    CONSTRAINT albums_pk PRIMARY KEY(id),

    CONSTRAINT albums_artist_id_fk FOREIGN KEY (artist_id)
        REFERENCES artists(id),

    CONSTRAINT albums_unique_path UNIQUE(path)
);

CREATE TABLE tracks (
    id TEXT,
    track_number INT NOT NULL,
    name TEXT NOT NULL,
    cover_art TEXT,

    duration INT,

    best_quality_file TEXT NOT NULL,
    mobile_quality_file TEXT NOT NULL,

    album_id TEXT NOT NULL,
    artist_id TEXT NOT NULL,

    path TEXT NOT NULL,
    available BOOL NOT NULL,

    CONSTRAINT tracks_pk PRIMARY KEY(id),

    -- CONSTRAINT tracks_track_number_unique UNIQUE(track_number, album_id) INITIALLY DEFERRED,

    CONSTRAINT tracks_album_id_fk FOREIGN KEY (album_id)
        REFERENCES albums(id),

    CONSTRAINT tracks_artist_id_fk FOREIGN KEY (artist_id)
        REFERENCES artists(id)
);


-- +goose Down
DROP TABLE tracks;
DROP TABLE albums;
DROP TABLE artists;
