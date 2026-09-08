-- +goose Up
CREATE TABLE events (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name       text NOT NULL,
    date       timestamptz NOT NULL,
    location   text,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE events;
