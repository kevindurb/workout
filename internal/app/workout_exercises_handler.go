package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"
	"github.com/kevindurb/planner/internal/middleware"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type WorkoutsExercisesHandler struct {
	q  *db.Queries
	sm *SessionManager
	fp *formparser.FormParser
}

func (h *WorkoutsExercisesHandler) Route(r chi.Router) {
	r.Get("/new", ghttp.Adapt(h.new))
	r.Post("/", h.create)

	r.Route("/{workout_exercise_id}", func(r chi.Router) {
		r.Use(middleware.EntityCtx(func(r *http.Request) (db.WorkoutsExercise, error) {
			id, _ := pathInt(r, "workout_exercise_id")
			userID := h.sm.UserID(r.Context())
			return h.q.GetWorkoutExerciseById(r.Context(), db.GetWorkoutExerciseByIdParams{
				ID:     id,
				UserID: userID,
			})
		}))
		r.Post("/delete", h.delete)
	})
}

func (h *WorkoutsExercisesHandler) new(w http.ResponseWriter, r *http.Request) (Node, error) {
	userID := h.sm.UserID(r.Context())
	workout := middleware.FromContext[db.Workout](r.Context())
	exercises, _ := h.q.ListAllExercises(r.Context(), userID)
	workoutExercises, _ := h.q.ListExercisesByWorkoutId(r.Context(), db.ListExercisesByWorkoutIdParams{
		WorkoutID: workout.ID,
		UserID:    userID,
	})
	return Layout(
		H1(Text("Choose Exercise")),
		A(Href("/exercises/new"), Text("Create Exercise")),
		H2(Text("Existing")),
		Map(workoutExercises, func(exercise db.ListExercisesByWorkoutIdRow) Node {
			return Form(
				Method("POST"),
				Action(fmt.Sprintf("/workouts/%d/exercises/%d/delete", workout.ID, exercise.ID)),
				Label(Text(exercise.Name)),
				Button(Type("submit"), Text("Remove")),
			)
		}),
		H2(Text("Ones to add")),
		Map(exercises, func(exercise db.Exercise) Node {
			return Form(
				Method("POST"),
				Action(fmt.Sprintf("/workouts/%d/exercises/%d", workout.ID, exercise.ID)),
				Label(Text(exercise.Name)),
				Button(Type("submit"), Text("Add")),
			)
		}),
	), nil
}

func (h *WorkoutsExercisesHandler) create(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	workout := middleware.FromContext[db.Workout](r.Context())
	exercise := middleware.FromContext[db.Exercise](r.Context())
	_, err := h.q.CreateWorkoutExercise(r.Context(), db.CreateWorkoutExerciseParams{
		UserID:     userID,
		WorkoutID:  workout.ID,
		ExerciseID: exercise.ID,
	})
	if err != nil {
		log.Printf("Error creating workout_exercise: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/workouts/%d/exercises/new", workout.ID), http.StatusFound)
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
