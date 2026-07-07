package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kevindurb/planner/internal/database/sqlcgen"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"
	"github.com/kevindurb/planner/internal/httpx"
	"github.com/kevindurb/planner/internal/middleware"
	"github.com/kevindurb/planner/internal/session"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type WorkoutsExercisesHandler struct {
	q  *sqlcgen.Queries
	sm *session.Manager
	fp *formparser.FormParser
}

type createWorkoutExerciseBody struct {
	ExerciseID int64
}

func (h *WorkoutsExercisesHandler) Route(r chi.Router) {
	r.Get("/new", ghttp.Adapt(h.new))
	r.Post("/", h.create)

	r.Route("/{workout_exercise_id}", func(r chi.Router) {
		r.Use(middleware.EntityCtx(func(r *http.Request) (sqlcgen.WorkoutsExercise, error) {
			return h.q.GetWorkoutExerciseById(r.Context(), sqlcgen.GetWorkoutExerciseByIdParams{
				ID:     httpx.PathInt(r, "workout_exercise_id"),
				UserID: h.sm.UserID(r.Context()),
			})
		}))
		r.Post("/delete", h.delete)
	})
}

func (h *WorkoutsExercisesHandler) new(w http.ResponseWriter, r *http.Request) (Node, error) {
	userID := h.sm.UserID(r.Context())
	workout := middleware.FromContext[sqlcgen.Workout](r.Context())
	exercises, _ := h.q.ListAllExercises(r.Context(), userID)
	workoutExercises, _ := h.q.ListExercisesByWorkoutId(r.Context(), sqlcgen.ListExercisesByWorkoutIdParams{
		WorkoutID: workout.ID,
		UserID:    userID,
	})
	return Layout(
		H1(Text("Choose Exercise")),
		A(Href("/exercises/new"), Text("Create Exercise")),
		H2(Text("Existing")),
		Map(workoutExercises, func(exercise sqlcgen.ListExercisesByWorkoutIdRow) Node {
			return Form(
				Method("POST"),
				Action(fmt.Sprintf("/workouts/%d/exercises/%d/delete", workout.ID, exercise.ID)),
				Label(Text(exercise.Name)),
				Button(Type("submit"), Text("Remove")),
			)
		}),
		H2(Text("Ones to add")),
		Map(exercises, func(exercise sqlcgen.Exercise) Node {
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
	workout := middleware.FromContext[sqlcgen.Workout](r.Context())
	var data createWorkoutExerciseBody
	err := h.fp.Parse(&data, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.q.CreateWorkoutExercise(r.Context(), sqlcgen.CreateWorkoutExerciseParams{
		UserID:     userID,
		WorkoutID:  workout.ID,
		ExerciseID: data.ExerciseID,
	})
	if err != nil {
		log.Printf("Error creating workout_exercise: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/workouts/%d/exercises/new", workout.ID), http.StatusFound)
}

func (h *WorkoutsExercisesHandler) delete(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	workoutExercise := middleware.FromContext[sqlcgen.WorkoutsExercise](r.Context())
	h.q.DeleteWorkoutExerciseByID(r.Context(), sqlcgen.DeleteWorkoutExerciseByIDParams{
		ID:     workoutExercise.ID,
		UserID: userID,
	})
	http.Redirect(w, r, fmt.Sprintf("/workouts/%d/exercises/edit", workoutExercise.WorkoutID), http.StatusFound)
}
