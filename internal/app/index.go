package app

import (
	"net/http"

	"github.com/kevindurb/planner/internal/html"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func (a *App) HandleIndex(w http.ResponseWriter, r *http.Request) (g.Node, error) {
	return html.Layout(
		h.H1(g.Text("Workout")),
	), nil
}
