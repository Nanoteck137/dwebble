-- +goose Up
CREATE TABLE playlists (
    id TEXT,
    name TEXT NOT NULL,

    owner_id TEXT NOT NULL,

    CONSTRAINT playlists_pk PRIMARY KEY(id),

    CONSTRAINT playlists_owner_id_fk FOREIGN KEY (owner_id)
        REFERENCES users(id)
);

CREATE TABLE playlist_items (
    playlist_id TEXT NOT NULL,
    track_id TEXT NOT NULL,
    item_index INTEGER NOT NULL, 

    CONSTRAINT playlist_items_pk PRIMARY KEY(playlist_id, track_id),

    CONSTRAINT playlist_items_playlist_id_fk FOREIGN KEY (playlist_id)
        REFERENCES playlists(id),

    CONSTRAINT playlist_items_track_id_fk FOREIGN KEY (track_id)
        REFERENCES tracks(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE playlist_items;
DROP TABLE playlists;
