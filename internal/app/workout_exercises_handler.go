package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type createWorkoutExerciseBody struct {
	ExerciseID int64 `form:"exercise_id,required"`
	WorkoutID  int64 `form:"workout_id,required"`
}

type WorkoutsExercisesHandler struct {
	q  *db.Queries
	sm *SessionManager
	fp *formparser.FormParser
}

func (h *WorkoutsExercisesHandler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /workouts/{workout_id}/exercises/edit", ghttp.Adapt(h.edit))
	mux.HandleFunc("POST /workouts_exercises", h.create)
	mux.HandleFunc("POST /workouts_exercises/{id}/delete", h.delete)
}

func (h *WorkoutsExercisesHandler) edit(w http.ResponseWriter, r *http.Request) (Node, error) {
	workoutID, _ := pathInt(r, "workout_id")
	userID := h.sm.UserID(r.Context())
	_, err := h.q.GetWorkoutByID(r.Context(), db.GetWorkoutByIDParams{
		ID:     workoutID,
		UserID: userID,
	})
	if err != nil {
		return nil, StatusCodeError{http.StatusBadRequest}
	}

	exercises, _ := h.q.ListAllExercises(r.Context(), userID)
	return Layout(
		H1(Text("Choose Exercise")),
		Map(exercises, func(exercise db.Exercise) Node {
			return Form(
				Method("POST"),
				Input(Type("hidden"), Name("workout_id"), Value(strconv.FormatInt(workoutID, 10))),
				Input(Type("hidden"), Name("exercise_id"), Value(strconv.FormatInt(exercise.ID, 10))),
				Action("/workouts_exercises"),
				Button(Type("submit"), Text("Add")),
			)
		}),
	), nil
}

func (h *WorkoutsExercisesHandler) create(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	var data createWorkoutExerciseBody
	if err := h.fp.Parse(&data, r); err != nil {
		log.Printf("Error parsing body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := h.q.GetWorkoutByID(r.Context(), db.GetWorkoutByIDParams{
		ID:     data.WorkoutID,
		UserID: userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.q.GetExerciseByID(r.Context(), db.GetExerciseByIDParams{
		ID:     data.ExerciseID,
		UserID: userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	workout, _ := h.q.CreateWorkoutExercise(r.Context(), db.CreateWorkoutExerciseParams{
		UserID:     userID,
		WorkoutID:  data.WorkoutID,
		ExerciseID: data.ExerciseID,
	})

	http.Redirect(w, r, fmt.Sprintf("/workouts/%d/exercises/edit", workout.ID), http.StatusFound)
}

func (h *WorkoutsExercisesHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, _ := pathInt(r, "id")
	userID := h.sm.UserID(r.Context())
	workoutExercise, err := h.q.GetWorkoutExerciseById(r.Context(), db.GetWorkoutExerciseByIdParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.q.DeleteWorkoutExerciseByID(r.Context(), db.DeleteWorkoutExerciseByIDParams{
		ID:     id,
		UserID: userID,
	})
	http.Redirect(w, r, fmt.Sprintf("/workouts/%d/exercises/edit", workoutExercise.WorkoutID), http.StatusFound)
}
