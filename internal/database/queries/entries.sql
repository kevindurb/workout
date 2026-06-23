-- name: ListAllEntries :many
SELECT *
FROM entries;

-- name: GetEntryByID :one
SELECT *
FROM entries
WHERE id = ?;

-- name: CreateEntry :one
INSERT INTO entries (
  name
) VALUES (?)
RETURNING *;

-- name: UpdateEntry :one
UPDATE entries
SET name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteEntryByID :exec
DELETE FROM entries
WHERE id = ?;
