-- name: ListAllEntries :many
SELECT *
FROM entries
WHERE user_id = ?;

-- name: GetEntryByID :one
SELECT *
FROM entries
WHERE id = ?
AND user_id = ?;

-- name: CreateEntry :one
INSERT INTO entries (
  user_id,
  workout_id,
  name
) VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateEntry :one
UPDATE entries
SET name = ?
WHERE id = ?
AND user_id = ?
RETURNING *;

-- name: DeleteEntryByID :exec
DELETE FROM entries
WHERE id = ?
AND user_id = ?;
