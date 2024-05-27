-- +goose Up
CREATE TABLE users (
    id TEXT,
    username TEXT NOT NULL,
    password TEXT NOT NULL,

    CONSTRAINT users_pk PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE users;
