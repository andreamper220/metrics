-- +goose Up
CREATE TABLE IF NOT EXISTS metrics_gauge (
    id varchar(128) PRIMARY KEY NOT NULL,
    value double precision NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS metrics_gauge;