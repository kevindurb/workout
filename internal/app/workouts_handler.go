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

var workoutPaths = Paths{"workouts"}

type createWorkoutBody struct {
	Name string `form:"name,required"`
}

type updateWorkoutBody struct {
	Name string `form:"name,required"`
}

type WorkoutsHandler struct {
	queries *db.Queries
	sm      *SessionManager
	fp      *formparser.FormParser
}

func (h *WorkoutsHandler) Route(r chi.Router) {
	workoutExercisesHandler := WorkoutsExercisesHandler{h.queries, h.sm, h.fp}
	r.Get("/", ghttp.Adapt(h.list))
	r.Get("/new", ghttp.Adapt(h.new))
	r.Post("/", h.create)

	r.Route("/{workout_id}", func(r chi.Router) {
		r.Use(middleware.EntityCtx(func(r *http.Request) (db.Workout, error) {
			id, _ := pathInt(r, "workout_id")
			userID := h.sm.UserID(r.Context())
			return h.queries.GetWorkoutByID(r.Context(), db.GetWorkoutByIDParams{
				ID:     id,
				UserID: userID,
			})
		}))
		r.Get("/", ghttp.Adapt(h.show))
		r.Get("/edit", ghttp.Adapt(h.edit))
		r.Post("/", h.update)
		r.Post("/delete", h.delete)

		r.Route("/exercises", workoutExercisesHandler.Route)
	})
}

func (h *WorkoutsHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	userID := h.sm.UserID(r.Context())
	workout := middleware.FromContext[db.Workout](r.Context())
	exercises, _ := h.queries.ListExercisesByWorkoutId(r.Context(), db.ListExercisesByWorkoutIdParams{
		WorkoutID: workout.ID,
		UserID:    userID,
	})
	return Layout(
		H1(Text(workout.Name)),
		A(Href(fmt.Sprintf("/workouts/%d/edit", workout.ID)), Text("Edit")),
		A(Href(fmt.Sprintf("/workouts/%d/exercises/edit", workout.ID)), Text("Exercises")),
		Ul(
			Map(exercises, func(exercise db.ListExercisesByWorkoutIdRow) Node {
				return Li(A(Href(fmt.Sprintf("/exercises/%d", exercise.ID)), Text(exercise.Name)))
			}),
		),
	), nil
}

func (h *WorkoutsHandler) list(w http.ResponseWriter, r *http.Request) (Node, error) {
	workouts, _ := h.queries.ListAllWorkouts(r.Context(), h.sm.UserID(r.Context()))
	return Layout(
		H1(Text("Workouts")),
		Map(workouts, func(workout db.Workout) Node {
			return P(Text(workout.Name))
		}),
	), nil
}

func (h *WorkoutsHandler) new(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("New Workout")),
		Form(
			Method("POST"),
			Action("/workouts"),
			Label(For("name"), Text("Name")),
			Input(Type("text"), ID("name"), Name("name"), Required()),
			Button(Type("submit"), Text("Create")),
		),
	), nil
}

func (h *WorkoutsHandler) edit(w http.ResponseWriter, r *http.Request) (Node, error) {
	workout := middleware.FromContext[db.Workout](r.Context())
	return Layout(
		H1(Text("Edit "+workout.Name)),
		Form(
			Method("POST"),
			Action(fmt.Sprintf("/workouts/%d", workout.ID)),
			Label(For("name"), Text("Name")),
			Input(Type("text"), ID("name"), Name("name"), Value(workout.Name), Required()),
			Button(Type("submit"), Text("Save")),
		),
	), nil
}

func (h *WorkoutsHandler) create(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	var data createWorkoutBody
	if err := h.fp.Parse(&data, r); err != nil {
		log.Printf("Error parsing body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	workout, err := h.queries.CreateWorkout(r.Context(), db.CreateWorkoutParams{
		UserID: userID,
		Name:   data.Name,
	})

	if err != nil {
		log.Printf("Error creating workout(user_id: %d, name: %s): %v", userID, data.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/workouts/%d", workout.ID), http.StatusFound)
}

func (h *WorkoutsHandler) update(w http.ResponseWriter, r *http.Request) {
	workout := middleware.FromContext[db.Workout](r.Context())
	userID := h.sm.UserID(r.Context())
	var data updateWorkoutBody
	if err := h.fp.Parse(&data, r); err != nil {
		log.Printf("Error parsing body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.queries.UpdateWorkout(r.Context(), db.UpdateWorkoutParams{
		ID:     workout.ID,
		UserID: userID,
		Name:   data.Name,
	})

	http.Redirect(w, r, fmt.Sprintf("/workouts/%d", workout.ID), http.StatusFound)
}

func (h *WorkoutsHandler) delete(w http.ResponseWriter, r *http.Request) {
	workout := middleware.FromContext[db.Workout](r.Context())
	userID := h.sm.UserID(r.Context())
	h.queries.DeleteWorkoutByID(r.Context(), db.DeleteWorkoutByIDParams{
		ID:     workout.ID,
		UserID: userID,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
