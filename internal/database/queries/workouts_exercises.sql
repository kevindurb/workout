-- name: GetWorkoutExerciseById :one
SELECT *
FROM workouts_exercises
WHERE id = ?
AND user_id = ?;

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

-- name: ListExercisesByWorkoutId :many
SELECT exercises.*
FROM exercises
JOIN workouts_exercises ON workouts_exercises.exercise_id = exercises.id
WHERE workouts_exercises.workout_id = ?
AND workouts_exercises.user_id = ?;
