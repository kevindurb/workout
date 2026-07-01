package app

import (
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type HomeHandler struct {
	queries *db.Queries
	sm      *SessionManager
	fp      *formparser.FormParser
}

func (h *HomeHandler) Routes(mux *http.ServeMux) {
	registerAuthRoutes(mux, h.sm, []Route{
		{"GET /", ghttp.Adapt(h.show)},
	})
}

func (h *HomeHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("Home")),
		A(Href("/workouts/new"), Text("New Workout")),
		A(Href("/exercises/new"), Text("New Exercise")),
	), nil
}
