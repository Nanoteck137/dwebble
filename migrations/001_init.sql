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

CREATE TABLE tags (
    id TEXT,
    name TEXT,

    CONSTRAINT tags_pk PRIMARY KEY(id),
    CONSTRAINT tags_name_unique UNIQUE(name)
);

CREATE TABLE tracks_to_tags (
    track_id TEXT,
    tag_id TEXT,

    CONSTRAINT tracks_to_tags_pk PRIMARY KEY(track_id, tag_id),

    CONSTRAINT tracks_to_tags_track_id_fk FOREIGN KEY (track_id)
        REFERENCES tracks(id),
    CONSTRAINT tracks_to_tags_tag_id_fk FOREIGN KEY (tag_id)
        REFERENCES tags(id)
);

-- +goose Down
DROP TABLE tracks_to_tags;
DROP TABLE tags;
DROP TABLE tracks;
DROP TABLE albums;
DROP TABLE artists;
