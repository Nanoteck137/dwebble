-- +goose Up
CREATE TABLE taglists (
    id TEXT,
    name TEXT NOT NULL,

    filter TEXT NOT NULL,

    owner_id TEXT NOT NULL,

    CONSTRAINT taglists_pk PRIMARY KEY(id),

    CONSTRAINT taglists_owner_id_fk FOREIGN KEY (owner_id)
        REFERENCES users(id)
);

-- +goose Down
DROP TABLE taglists;
