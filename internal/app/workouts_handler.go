package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

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

func (h *WorkoutsHandler) Routes(mux *http.ServeMux) {
	registerAuthRoutes(mux, h.sm, []Route{
		{"GET /workouts/{id}", ghttp.Adapt(h.show)},
		{"GET /workouts/{id}/edit", ghttp.Adapt(h.edit)},
		{"GET /workouts", ghttp.Adapt(h.list)},
		{"GET /workouts/new", ghttp.Adapt(h.new)},
		{"POST /workouts", http.HandlerFunc(h.create)},
		{"POST /workouts/{id}", http.HandlerFunc(h.update)},
		{"POST /workouts/{id}/delete", http.HandlerFunc(h.delete)},
	})
}

func (h *WorkoutsHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	id, _ := pathInt(r, "id")
	userID := h.sm.UserID(r.Context())
	workout, err := h.queries.GetWorkoutByID(r.Context(), db.GetWorkoutByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, StatusCodeError{http.StatusNotFound}
	}
	return Layout(
		H1(Text(workout.Name)),
		A(Href(fmt.Sprintf("/workouts/%d/edit", workout.ID)), Text("Edit")),
		A(Href(fmt.Sprintf("/workouts/%d/exercises", workout.ID)), Text("Exercises")),
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
	id, _ := pathInt(r, "id")
	userID := h.sm.UserID(r.Context())
	workout, err := h.queries.GetWorkoutByID(r.Context(), db.GetWorkoutByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, StatusCodeError{http.StatusNotFound}
	}
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
	id, _ := pathInt(r, "id")
	userID := h.sm.UserID(r.Context())
	var data updateWorkoutBody
	if err := h.fp.Parse(&data, r); err != nil {
		log.Printf("Error parsing body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.queries.UpdateWorkout(r.Context(), db.UpdateWorkoutParams{
		ID:     id,
		UserID: userID,
		Name:   data.Name,
	})

	http.Redirect(w, r, fmt.Sprintf("/workouts/%d", id), http.StatusFound)
}

func (h *WorkoutsHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, _ := pathInt(r, "id")
	userID := h.sm.UserID(r.Context())
	h.queries.DeleteWorkoutByID(r.Context(), db.DeleteWorkoutByIDParams{
		ID:     id,
		UserID: userID,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
