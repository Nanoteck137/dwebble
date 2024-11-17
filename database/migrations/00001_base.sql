-- +goose Up
CREATE TABLE artists (
    id TEXT NOT NULL,
    name TEXT NOT NULL CHECK(name<>''),
    picture TEXT,

    CONSTRAINT artists_pk PRIMARY KEY (id),
    UNIQUE (name COLLATE NOCASE)
);

CREATE TABLE albums (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    artist_id TEXT NOT NULL,

    cover_art TEXT,
    year INT,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    available BOOL NOT NULL,

    CONSTRAINT albums_pk PRIMARY KEY(id),

    CONSTRAINT albums_artist_id_fk FOREIGN KEY (artist_id)
        REFERENCES artists(id)
);

CREATE TABLE tracks (
    id TEXT NOT NULL,
    name TEXT NOT NULL,

    album_id TEXT NOT NULL,
    artist_id TEXT NOT NULL,

    track_number INT,
    duration INT,
    year INT,

    export_name TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    mobile_filename TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    available BOOL NOT NULL,

    CONSTRAINT tracks_pk PRIMARY KEY(id),

    CONSTRAINT tracks_album_id_fk FOREIGN KEY (album_id)
        REFERENCES albums(id),

    CONSTRAINT tracks_artist_id_fk FOREIGN KEY (artist_id)
        REFERENCES artists(id)
);


-- +goose Down
DROP TABLE tracks;
DROP TABLE albums;
DROP TABLE artists;
