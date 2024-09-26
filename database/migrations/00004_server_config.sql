-- +goose Up

CREATE TABLE config (
    id INTEGER NOT NULL,
    owner_id TEXT NOT NULL,

    CONSTRAINT config_pk PRIMARY KEY(id),

    CONSTRAINT config_owner_id_fk FOREIGN KEY (owner_id)
        REFERENCES users(id)
);

-- +goose Down
DROP TABLE config;
