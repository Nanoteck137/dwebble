CREATE TABLE IF NOT EXISTS artists (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    picture TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS albums (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    artist_id TEXT NOT NULL,
    picture TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at BIGINT NOT NULL DEFAULT 0,

    FOREIGN KEY(artist_id) REFERENCES artists(id)
);

CREATE TABLE IF NOT EXISTS tracks (
    id TEXT PRIMARY KEY NOT NULL,
    num INT NOT NULL,
    name TEXT NOT NULL,
    artist_id TEXT NOT NULL,
    album_id TEXT NOT NULL,

    cover_art TEXT,
    
    file_quality TEXT NOT NULL,
    file_mobile TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at BIGINT NOT NULL DEFAULT 0,

    FOREIGN KEY(artist_id) REFERENCES artists(id),
    FOREIGN KEY(album_id) REFERENCES albums(id)
);

