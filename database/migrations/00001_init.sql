-- +goose Up
CREATE TABLE artists (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(name<>''),
    other_name TEXT,

    picture TEXT,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL

    -- UNIQUE (name COLLATE NOCASE)
);

CREATE TABLE albums (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(name<>''),
    other_name TEXT,

    artist_id TEXT NOT NULL REFERENCES artists(id),

    cover_art TEXT,
    year INT,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE albums_extra_artists (
    album_id TEXT NOT NULL REFERENCES albums(id),
    artist_id TEXT NOT NULL REFERENCES artists(id),

    PRIMARY KEY(album_id, artist_id)
);

CREATE TABLE tracks (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(name<>''),
    other_name TEXT,

    album_id TEXT NOT NULL REFERENCES albums(id),
    artist_id TEXT NOT NULL REFERENCES artists(id),

    track_number INT,
    duration INT,
    year INT,

    export_name TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    mobile_filename TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE tracks_extra_artists (
    track_id TEXT NOT NULL REFERENCES tracks(id),
    artist_id TEXT NOT NULL REFERENCES artists(id),

    PRIMARY KEY(track_id, artist_id)
);

CREATE TABLE tags (
    slug TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE tracks_to_tags (
    track_id TEXT REFERENCES tracks(id) ON DELETE CASCADE,
    tag_slug TEXT REFERENCES tags(slug),

    PRIMARY KEY(track_id, tag_slug)
);

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL COLLATE NOCASE CHECK(username<>'') UNIQUE,
    password TEXT NOT NULL CHECK(password<>''),
    role TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE users_settings (
    id TEXT PRIMARY KEY REFERENCES users(id),
    display_name TEXT,

    quick_playlist TEXT REFERENCES playlists(id)
);

CREATE TABLE playlists (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(name<>''),

    owner_id TEXT NOT NULL REFERENCES users(id),

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE playlist_items (
    playlist_id TEXT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    item_index INTEGER NOT NULL, 

    PRIMARY KEY(playlist_id, track_id)
);

CREATE TABLE taglists (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(name<>''),

    filter TEXT NOT NULL,

    -- TODO(patrik): Nullable for system lists
    owner_id TEXT NOT NULL REFERENCES users(id),

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);


-- +goose Down
DROP TABLE taglists;

DROP TABLE playlist_items;
DROP TABLE playlists;

DROP TABLE users_settings;
DROP TABLE users;

DROP TABLE tracks_to_tags;
DROP TABLE tags;

DROP TABLE tracks;
DROP TABLE albums;
DROP TABLE artists;
