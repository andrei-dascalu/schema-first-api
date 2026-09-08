-- +goose Up
CREATE TYPE rsvp_status AS ENUM ('pending', 'accepted', 'declined');

CREATE TABLE guests (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    uuid NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    name        text NOT NULL,
    email       text,
    rsvp_status rsvp_status NOT NULL DEFAULT 'pending',
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX guests_event_id_idx ON guests (event_id);

-- +goose Down
DROP TABLE guests;
DROP TYPE rsvp_status;
