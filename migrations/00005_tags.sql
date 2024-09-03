-- +goose Up
CREATE TABLE tags (
    slug TEXT,
    name TEXT NOT NULL,

    CONSTRAINT tags_pk PRIMARY KEY(slug)
);

CREATE TABLE tracks_to_tags (
    track_id TEXT,
    tag_slug TEXT,

    CONSTRAINT tracks_to_tags_pk PRIMARY KEY(track_id, tag_slug),

    CONSTRAINT tracks_to_tags_track_id_fk FOREIGN KEY (track_id)
        REFERENCES tracks(id),
    CONSTRAINT tracks_to_tags_tag_id_fk FOREIGN KEY (tag_slug)
        REFERENCES tags(slug)
);

-- +goose Down
DROP TABLE tracks_to_tags;
DROP TABLE tags;
