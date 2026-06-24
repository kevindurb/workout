package app

import (
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	. "github.com/kevindurb/planner/internal/html"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type WorkoutsHandler struct {
	queries *db.Queries
	sm      *SessionManager
}

func (h *WorkoutsHandler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /{id}", ghttp.Adapt(h.show))
	mux.Handle("GET /{id}/edit", ghttp.Adapt(h.edit))
	mux.Handle("GET /", ghttp.Adapt(h.list))
	mux.Handle("GET /new", ghttp.Adapt(h.new))
	mux.HandleFunc("POST /", h.create)
	mux.HandleFunc("POST /{id}", h.update)
	mux.HandleFunc("POST /{id}/delete", h.delete)

	return mux
}

func (h *WorkoutsHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	id, _ := pathInt(r, "id")
	workout, err := h.queries.GetWorkoutByID(r.Context(), id)
	if err != nil {
		return nil, StatusCodeError{http.StatusNotFound}
	}
	return Layout(
		H1(Text(workout.Name)),
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
	), nil
}

func (h *WorkoutsHandler) edit(w http.ResponseWriter, r *http.Request) (Node, error) {
	id, _ := pathInt(r, "id")
	workout, err := h.queries.GetWorkoutByID(r.Context(), id)
	if err != nil {
		return nil, StatusCodeError{http.StatusNotFound}
	}
	return Layout(
		H1(Text("Edit " + workout.Name)),
	), nil
}

func (h *WorkoutsHandler) create(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/workouts/", http.StatusFound)
}

func (h *WorkoutsHandler) update(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/workouts/", http.StatusFound)
}

func (h *WorkoutsHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, _ := pathInt(r, "id")
	h.queries.DeleteWorkoutByID(r.Context(), id)
	http.Redirect(w, r, "/", http.StatusFound)
}
