-- name: ListAllWorkouts :many
SELECT *
FROM workouts
WHERE user_id = ?;

-- name: GetWorkoutByID :one
SELECT *
FROM workouts
WHERE id = ?;

-- name: CreateWorkout :one
INSERT INTO workouts (
  user_id,
  name
) VALUES (?, ?)
RETURNING *;

-- name: UpdateWorkout :one
UPDATE workouts
SET name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteWorkoutByID :exec
DELETE FROM workouts
WHERE id = ?;
