-- +goose Up
CREATE TABLE IF NOT EXISTS metrics_counter (
    id varchar(128) PRIMARY KEY NOT NULL,
    value bigint NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS metrics_counter;