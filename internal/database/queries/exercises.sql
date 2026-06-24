-- name: ListAllExercises :many
SELECT *
FROM exercises
WHERE user_id = ?;

-- name: GetExerciseByID :one
SELECT *
FROM exercises
WHERE id = ?;

-- name: CreateExercise :one
INSERT INTO exercises (
  user_id,
  name
) VALUES (?, ?)
RETURNING *;

-- name: UpdateExercise :one
UPDATE exercises
SET name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteExerciseByID :exec
DELETE FROM exercises
WHERE id = ?;
