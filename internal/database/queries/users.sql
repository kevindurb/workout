-- name: ListAllUsers :many
SELECT *
FROM users;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;

-- name: CreateUser :one
INSERT INTO users (
  email,
  hash
) VALUES (?, ?)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = ?
AND hash = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = ?;
