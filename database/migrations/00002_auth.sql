-- NOTE: Final

-- +goose Up
CREATE TABLE users (
    id TEXT,
    username TEXT NOT NULL CHECK(username<>''),
    password TEXT NOT NULL CHECK(password<>''),

    CONSTRAINT users_pk PRIMARY KEY(id),
    UNIQUE (username COLLATE NOCASE)
);

-- +goose Down
DROP TABLE users;
