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
	q  *db.Queries
	sm *SessionManager
	fp *formparser.FormParser
}

func (h *HomeHandler) Routes(mux *http.ServeMux) {
	registerAuthRoutes(mux, h.sm, []Route{
		{"GET /", ghttp.Adapt(h.show)},
	})
}

func (h *HomeHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	userID := h.sm.UserID(r.Context())
	workouts, _ := h.q.ListAllWorkouts(r.Context(), userID)
	return Layout(
		H1(Text("Home")),
		A(Href(workoutsPathBuilder.New()), Text("New Workout")),
		A(Href(exercisesPathBuilder.New()), Text("New Exercise")),
		Ul(
			Map(workouts, func(workout db.Workout) Node {
				return Li(A(Href(workoutsPathBuilder.Show(workout.ID)), Text(workout.Name)))
			}),
		),
	), nil
}
