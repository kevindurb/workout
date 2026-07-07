package entries

import (
	"net/http"

	"github.com/kevindurb/planner/internal/database/sqlcgen"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"
	"github.com/kevindurb/planner/internal/middleware"
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
	entry := middleware.FromContext[sqlcgen.Entry](r.Context())
	return Layout(
		H1(Text(entry.Name)),
	), nil
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) (Node, error) {
	entries, _ := h.q.ListAllEntries(r.Context(), h.sm.UserID(r.Context()))
	return Layout(
		H1(Text("Entries")),
		Map(entries, func(entry sqlcgen.Entry) Node {
			return P(Text(entry.Name))
		}),
	), nil
}

func (h *Handler) new(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("New Entry")),
	), nil
}

func (h *Handler) edit(w http.ResponseWriter, r *http.Request) (Node, error) {
	entry := middleware.FromContext[sqlcgen.Entry](r.Context())
	return Layout(
		H1(Text("Edit " + entry.Name)),
	), nil
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, routes.Entries.List(), http.StatusFound)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, routes.Entries.List(), http.StatusFound)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	entry := middleware.FromContext[sqlcgen.Entry](r.Context())
	h.q.DeleteEntryByID(r.Context(), sqlcgen.DeleteEntryByIDParams{ID: entry.ID, UserID: userID})
	http.Redirect(w, r, "/", http.StatusFound)
}
