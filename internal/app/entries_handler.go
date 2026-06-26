package app

import (
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	. "github.com/kevindurb/planner/internal/html"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type EntriesHandler struct {
	queries *db.Queries
	sm      *SessionManager
}

func (h *EntriesHandler) Routes(mux *http.ServeMux) {
	registerAuthRoutes(mux, h.sm, []Route{
		{"GET /entries/{id}", ghttp.Adapt(h.show)},
		{"GET /entries/{id}/edit", ghttp.Adapt(h.edit)},
		{"GET /entries", ghttp.Adapt(h.list)},
		{"GET /entries/new", ghttp.Adapt(h.new)},

		{"POST /entries", http.HandlerFunc(h.create)},
		{"POST /entries/{id}", http.HandlerFunc(h.update)},
		{"POST /entries/{id}/delete", http.HandlerFunc(h.delete)},
	})
}

func (h *EntriesHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	id, _ := pathInt(r, "id")
	entry, err := h.queries.GetEntryByID(r.Context(), id)
	if err != nil {
		return nil, StatusCodeError{http.StatusNotFound}
	}
	return Layout(
		H1(Text(entry.Name)),
	), nil
}

func (h *EntriesHandler) list(w http.ResponseWriter, r *http.Request) (Node, error) {
	entries, _ := h.queries.ListAllEntries(r.Context(), h.sm.UserID(r.Context()))
	return Layout(
		H1(Text("Entries")),
		Map(entries, func(entry db.Entry) Node {
			return P(Text(entry.Name))
		}),
	), nil
}

func (h *EntriesHandler) new(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("New Entry")),
	), nil
}

func (h *EntriesHandler) edit(w http.ResponseWriter, r *http.Request) (Node, error) {
	id, _ := pathInt(r, "id")
	entry, err := h.queries.GetEntryByID(r.Context(), id)
	if err != nil {
		return nil, StatusCodeError{http.StatusNotFound}
	}
	return Layout(
		H1(Text("Edit " + entry.Name)),
	), nil
}

func (h *EntriesHandler) create(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/entries/", http.StatusFound)
}

func (h *EntriesHandler) update(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/entries/", http.StatusFound)
}

func (h *EntriesHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, _ := pathInt(r, "id")
	h.queries.DeleteEntryByID(r.Context(), id)
	http.Redirect(w, r, "/", http.StatusFound)
}
