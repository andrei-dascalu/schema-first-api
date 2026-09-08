-- name: CreateEvent :one
INSERT INTO events (name, date, location)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListEvents :many
SELECT * FROM events ORDER BY date;

-- name: GetEvent :one
SELECT * FROM events WHERE id = $1;

-- name: UpdateEvent :one
UPDATE events
SET name = $2, date = $3, location = $4
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :execrows
DELETE FROM events WHERE id = $1;
