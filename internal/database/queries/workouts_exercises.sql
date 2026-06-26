-- name: CreateWorkoutExercise :one
INSERT INTO workouts_exercises (
  user_id,
  workout_id,
  exercise_id
) VALUES (?, ?, ?)
RETURNING *;

-- name: DeleteWorkoutExerciseByID :exec
DELETE FROM workouts_exercises
WHERE id = ?
AND user_id = ?;

