package home

import (
	"net/http"

	"github.com/kevindurb/planner/internal/database/sqlcgen"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"
	"github.com/kevindurb/planner/internal/routes"
	"github.com/kevindurb/planner/internal/session"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type Handler struct {
	q  *sqlcgen.Queries
	sm *session.Manager
	fp *formparser.FormParser
}

func NewHandler(
	q *sqlcgen.Queries,
	sm *session.Manager,
	fp *formparser.FormParser,
) *Handler {
	return &Handler{q, sm, fp}
}

func (h *Handler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	userID := h.sm.UserID(r.Context())
	workouts, _ := h.q.ListAllWorkouts(r.Context(), userID)
	return Layout(
		H1(Text("Home")),
		A(Href(routes.Workouts.New()), Text("New Workout")),
		A(Href(routes.Exercises.New()), Text("New Exercise")),
		Ul(
			Map(workouts, func(workout sqlcgen.Workout) Node {
				return Li(A(Href(routes.Workouts.Show(workout.ID)), Text(workout.Name)))
			}),
		),
	), nil
}
