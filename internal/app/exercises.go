package app

import (
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	. "github.com/kevindurb/planner/internal/html"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type ExercisesHandler struct {
	queries *db.Queries
	sm      *SessionManager
}

func (h *ExercisesHandler) Routes() http.Handler {
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

func (h *ExercisesHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	id, _ := pathInt(r, "id")
	exercise, err := h.queries.GetExerciseByID(r.Context(), id)
	if err != nil {
		return nil, StatusCodeError{http.StatusNotFound}
	}
	return Layout(
		H1(Text(exercise.Name)),
	), nil
}

func (h *ExercisesHandler) list(w http.ResponseWriter, r *http.Request) (Node, error) {
	exercises, _ := h.queries.ListAllExercises(r.Context(), h.sm.UserID(r.Context()))
	return Layout(
		H1(Text("Exercises")),
		Map(exercises, func(exercise db.Exercise) Node {
			return P(Text(exercise.Name))
		}),
	), nil
}

func (h *ExercisesHandler) new(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("New Exercise")),
	), nil
}

func (h *ExercisesHandler) edit(w http.ResponseWriter, r *http.Request) (Node, error) {
	id, _ := pathInt(r, "id")
	exercise, err := h.queries.GetExerciseByID(r.Context(), id)
	if err != nil {
		return nil, StatusCodeError{http.StatusNotFound}
	}
	return Layout(
		H1(Text("Edit " + exercise.Name)),
	), nil
}

func (h *ExercisesHandler) create(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/exercises/", http.StatusFound)
}

func (h *ExercisesHandler) update(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/exercises/", http.StatusFound)
}

func (h *ExercisesHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, _ := pathInt(r, "id")
	h.queries.DeleteExerciseByID(r.Context(), id)
	http.Redirect(w, r, "/", http.StatusFound)
}
