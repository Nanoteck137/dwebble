-- +goose Up
CREATE TABLE genres (
    id TEXT,
    name TEXT,

    CONSTRAINT genres_pk PRIMARY KEY(id),
    CONSTRAINT genres_name_unique UNIQUE(name)
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
