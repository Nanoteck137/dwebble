-- +goose Up

ALTER TABLE artists ADD COLUMN created INTEGER DEFAULT 0 NOT NULL;
ALTER TABLE artists ADD COLUMN updated INTEGER DEFAULT 0 NOT NULL;
UPDATE artists SET created = cast((unixepoch("subsec") * 1000.0) as INTEGER);
UPDATE artists SET updated = cast((unixepoch("subsec") * 1000.0) as INTEGER);

ALTER TABLE playlists ADD COLUMN created INTEGER DEFAULT 0 NOT NULL;
ALTER TABLE playlists ADD COLUMN updated INTEGER DEFAULT 0 NOT NULL;
UPDATE playlists SET created = cast((unixepoch("subsec") * 1000.0) as INTEGER);
UPDATE playlists SET updated = cast((unixepoch("subsec") * 1000.0) as INTEGER);

ALTER TABLE playlist_items ADD COLUMN created INTEGER DEFAULT 0 NOT NULL;
ALTER TABLE playlist_items ADD COLUMN updated INTEGER DEFAULT 0 NOT NULL;
UPDATE playlist_items SET created = cast((unixepoch("subsec") * 1000.0) as INTEGER);
UPDATE playlist_items SET updated = cast((unixepoch("subsec") * 1000.0) as INTEGER);

-- +goose Down

ALTER TABLE playlist_items DROP COLUMN updated;
ALTER TABLE playlist_items DROP COLUMN created;

ALTER TABLE playlists DROP COLUMN updated;
ALTER TABLE playlists DROP COLUMN created;

ALTER TABLE artists DROP COLUMN updated;
ALTER TABLE artists DROP COLUMN created;


