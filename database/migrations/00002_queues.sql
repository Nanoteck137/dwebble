-- +goose Up
CREATE TABLE players (
    id TEXT PRIMARY KEY, /* dwebble-web-frontend */

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE queues (
    id TEXT PRIMARY KEY,
    player_id TEXT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK(name<>''),
    item_index INT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE default_queues (
    player_id TEXT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    queue_id TEXT NOT NULL REFERENCES queues(id) ON DELETE CASCADE,

    PRIMARY KEY(player_id, user_id)
);

CREATE TABLE queue_items (
    id TEXT PRIMARY KEY,
    queue_id TEXT NOT NULL REFERENCES queues(id) ON DELETE CASCADE,

    order_number INT NOT NULL,
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

-- +goose Down
DROP TABLE queue_items;
DROP TABLE default_queues;
DROP TABLE queues;
DROP TABLE players;
