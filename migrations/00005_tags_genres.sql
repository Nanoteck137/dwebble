-- +goose Up
CREATE TABLE tags (
    id TEXT,
    name TEXT NOT NULL,
    display_name TEXT NOT NULL,

    CONSTRAINT tags_pk PRIMARY KEY(id),
    CONSTRAINT tags_name_unique UNIQUE(name COLLATE NOCASE)
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

CREATE TABLE genres (
    id TEXT,
    name TEXT NOT NULL,
    display_name TEXT NOT NULL,

    CONSTRAINT genres_pk PRIMARY KEY(id),
    CONSTRAINT genres_name_unique UNIQUE(name COLLATE NOCASE)
);

CREATE TABLE tracks_to_genres (
    track_id TEXT,
    genre_id TEXT,

    CONSTRAINT tracks_to_genres_pk PRIMARY KEY(track_id, genre_id),

    CONSTRAINT tracks_to_genres_track_id_fk FOREIGN KEY (track_id)
        REFERENCES tracks(id),
    CONSTRAINT tracks_to_genres_tag_id_fk FOREIGN KEY (genre_id)
        REFERENCES genres(id)
);

-- +goose Down
DROP TABLE tracks_to_genres;
DROP TABLE genres;
DROP TABLE tracks_to_tags;
DROP TABLE tags;
