-- name: ListAllExercises :many
SELECT *
FROM exercises;

-- name: GetExerciseByID :one
SELECT *
FROM exercises
WHERE id = ?;

-- name: CreateExercise :one
INSERT INTO exercises (
  name
) VALUES (?)
RETURNING *;

-- name: UpdateExercise :one
UPDATE exercises
SET name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteExerciseByID :exec
DELETE FROM exercises
WHERE id = ?;
