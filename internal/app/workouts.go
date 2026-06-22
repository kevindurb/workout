package app

import (
	"net/http"

	. "github.com/kevindurb/planner/internal/html"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type WorkoutsHandler struct {
}

func NewWorkoutsHandler() *WorkoutsHandler {
	return &WorkoutsHandler{}
}

func (h *WorkoutsHandler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", ghttp.Adapt(h.list))
	mux.Handle("GET /new", ghttp.Adapt(h.new))
	mux.HandleFunc("POST /", h.create)

	return mux
}

func (h *WorkoutsHandler) list(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("Workouts")),
	), nil
}

func (h *WorkoutsHandler) new(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("New Workout")),
	), nil
}

func (h *WorkoutsHandler) create(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/workouts/", http.StatusFound)
}
