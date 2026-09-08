-- name: CreateGuest :one
INSERT INTO guests (event_id, name, email, rsvp_status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListGuestsByEvent :many
SELECT * FROM guests WHERE event_id = $1 ORDER BY created_at;

-- name: GetGuest :one
SELECT * FROM guests WHERE id = $1 AND event_id = $2;

-- name: UpdateGuest :one
UPDATE guests
SET name = $3, email = $4, rsvp_status = $5
WHERE id = $1 AND event_id = $2
RETURNING *;

-- name: DeleteGuest :execrows
DELETE FROM guests WHERE id = $1 AND event_id = $2;
