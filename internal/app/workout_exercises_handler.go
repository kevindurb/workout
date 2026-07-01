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
	registerAuthRoutes(mux, h.sm, []Route{
		{"GET /workouts/{workout_id}/exercises/edit", ghttp.Adapt(h.edit)},
		{"POST /workouts_exercises", http.HandlerFunc(h.create)},
		{"POST /workouts_exercises/{id}/delete", http.HandlerFunc(h.delete)},
	})
}

func (h *WorkoutsExercisesHandler) edit(w http.ResponseWriter, r *http.Request) (Node, error) {
	workoutID, _ := pathInt(r, "workout_id")
	userID := h.sm.UserID(r.Context())
	_, err := h.q.GetWorkoutByID(r.Context(), db.GetWorkoutByIDParams{
		ID:     workoutID,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error getting workout(%d): %v", workoutID, err)
		return nil, StatusCodeError{http.StatusBadRequest}
	}

	exercises, _ := h.q.ListAllExercises(r.Context(), userID)
	workoutExercises, _ := h.q.ListExercisesByWorkoutId(r.Context(), db.ListExercisesByWorkoutIdParams{
		WorkoutID: workoutID,
		UserID:    userID,
	})
	return Layout(
		H1(Text("Choose Exercise")),
		A(Href("/exercises/new"), Text("Create Exercise")),
		H2(Text("Existing")),
		Map(workoutExercises, func(exercise db.ListExercisesByWorkoutIdRow) Node {
			return Form(
				Method("POST"),
				Action(fmt.Sprintf("/workouts_exercises/%d/delete", exercise.WorkoutExerciseID)),
				Label(Text(exercise.Name)),
				Button(Type("submit"), Text("Remove")),
			)
		}),
		H2(Text("Ones to add")),
		Map(exercises, func(exercise db.Exercise) Node {
			return Form(
				Method("POST"),
				Action("/workouts_exercises"),
				Input(Type("hidden"), Name("workout_id"), Value(strconv.FormatInt(workoutID, 10))),
				Input(Type("hidden"), Name("exercise_id"), Value(strconv.FormatInt(exercise.ID, 10))),
				Label(Text(exercise.Name)),
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
		log.Printf("Error getting workout (%d): %v", data.WorkoutID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.q.GetExerciseByID(r.Context(), db.GetExerciseByIDParams{
		ID:     data.ExerciseID,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error getting exercise (%d): %v", data.ExerciseID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.q.CreateWorkoutExercise(r.Context(), db.CreateWorkoutExerciseParams{
		UserID:     userID,
		WorkoutID:  data.WorkoutID,
		ExerciseID: data.ExerciseID,
	})
	if err != nil {
		log.Printf("Error creating workout_exercise: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/workouts/%d/exercises/edit", data.WorkoutID), http.StatusFound)
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
