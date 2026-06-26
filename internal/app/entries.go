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
	mux.Handle("GET /entries/{id}", ghttp.Adapt(h.show))
	mux.Handle("GET /entries/{id}/edit", ghttp.Adapt(h.edit))
	mux.Handle("GET /entries", ghttp.Adapt(h.list))
	mux.Handle("GET /entries/new", ghttp.Adapt(h.new))
	mux.HandleFunc("POST /entries", h.create)
	mux.HandleFunc("POST /entries/{id}", h.update)
	mux.HandleFunc("POST /entries/{id}/delete", h.delete)
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
