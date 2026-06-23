package app

import (
	"net/http"

	. "github.com/kevindurb/planner/internal/html"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type HomeHandler struct {
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", ghttp.Adapt(h.show))

	return mux
}

func (h *HomeHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("Home")),
		A(Href("/workouts/new"), Text("New Workout")),
		A(Href("/exercises/new"), Text("New Exercise")),
	), nil
}
