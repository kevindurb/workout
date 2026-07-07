package workouts

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"
	"github.com/kevindurb/planner/internal/middleware"
	"github.com/kevindurb/planner/internal/session"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type createWorkoutBody struct {
	Name string `form:"name,required"`
}

type updateWorkoutBody struct {
	Name string `form:"name,required"`
}

type Handler struct {
	q  *db.Queries
	sm *session.Manager
	fp *formparser.FormParser
}

func NewHandler(
	q *db.Queries,
	sm *session.Manager,
	fp *formparser.FormParser,
) *Handler {
	return &Handler{q, sm, fp}
}

func (h *Handler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	userID := h.sm.UserID(r.Context())
	workout := middleware.FromContext[db.Workout](r.Context())
	exercises, _ := h.q.ListExercisesByWorkoutId(r.Context(), db.ListExercisesByWorkoutIdParams{
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

func (h *Handler) list(w http.ResponseWriter, r *http.Request) (Node, error) {
	workouts, _ := h.q.ListAllWorkouts(r.Context(), h.sm.UserID(r.Context()))
	return Layout(
		H1(Text("Workouts")),
		Map(workouts, func(workout db.Workout) Node {
			return P(Text(workout.Name))
		}),
	), nil
}

func (h *Handler) new(w http.ResponseWriter, r *http.Request) (Node, error) {
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

func (h *Handler) edit(w http.ResponseWriter, r *http.Request) (Node, error) {
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

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	var data createWorkoutBody
	if err := h.fp.Parse(&data, r); err != nil {
		log.Printf("Error parsing body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	workout, err := h.q.CreateWorkout(r.Context(), db.CreateWorkoutParams{
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

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	workout := middleware.FromContext[db.Workout](r.Context())
	userID := h.sm.UserID(r.Context())
	var data updateWorkoutBody
	if err := h.fp.Parse(&data, r); err != nil {
		log.Printf("Error parsing body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.q.UpdateWorkout(r.Context(), db.UpdateWorkoutParams{
		ID:     workout.ID,
		UserID: userID,
		Name:   data.Name,
	})

	http.Redirect(w, r, fmt.Sprintf("/workouts/%d", workout.ID), http.StatusFound)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	workout := middleware.FromContext[db.Workout](r.Context())
	userID := h.sm.UserID(r.Context())
	h.q.DeleteWorkoutByID(r.Context(), db.DeleteWorkoutByIDParams{
		ID:     workout.ID,
		UserID: userID,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
