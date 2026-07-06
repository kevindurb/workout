package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"
	"github.com/kevindurb/planner/internal/routes"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type HomeHandler struct {
	q  *db.Queries
	sm *SessionManager
	fp *formparser.FormParser
}

func (h *HomeHandler) Route(r chi.Router) {
	r.Get("/", ghttp.Adapt(h.show))
}

func (h *HomeHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	userID := h.sm.UserID(r.Context())
	workouts, _ := h.q.ListAllWorkouts(r.Context(), userID)
	return Layout(
		H1(Text("Home")),
		A(Href(routes.Workouts.New()), Text("New Workout")),
		A(Href(routes.Exercises.New()), Text("New Exercise")),
		Ul(
			Map(workouts, func(workout db.Workout) Node {
				return Li(A(Href(routes.Workouts.Show(workout.ID)), Text(workout.Name)))
			}),
		),
	), nil
}
