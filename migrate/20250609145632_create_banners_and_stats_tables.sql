-- +goose Up
CREATE TABLE IF NOT EXISTS banners (
                                       id   SERIAL PRIMARY KEY,
                                       name TEXT NOT NULL,
                                       CONSTRAINT banners_name_unique UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS stats (
    banner_id INTEGER NOT NULL REFERENCES banners(id),
    timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    count INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (banner_id, timestamp)
);

-- +goose Down
DROP TABLE IF EXISTS stats;
DROP TABLE IF EXISTS banners;