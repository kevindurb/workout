package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kevindurb/planner/internal/db"
	. "github.com/kevindurb/planner/internal/html"
	ihttp "github.com/kevindurb/planner/internal/http"
	"github.com/kevindurb/planner/internal/middleware"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

var entriesPaths = Paths{"entries"}

type EntriesHandler struct {
	queries *db.Queries
	sm      *SessionManager
}

func (h *EntriesHandler) Route(r chi.Router) {
	r.Get("/", ghttp.Adapt(h.list))
	r.Get("/new", ghttp.Adapt(h.new))
	r.Post("/", h.create)

	r.Route("/{entry_id}", func(r chi.Router) {
		r.Use(middleware.EntityCtx(func(r *http.Request) (db.Entry, error) {
			return h.queries.GetEntryByID(r.Context(), db.GetEntryByIDParams{
				ID:     ihttp.PathInt(r, "entry_id"),
				UserID: h.sm.UserID(r.Context()),
			})
		}))
		r.Get("/", ghttp.Adapt(h.show))
		r.Get("/edit", ghttp.Adapt(h.edit))
		r.Post("/", h.update)
		r.Post("/delete", h.delete)
	})
}

func (h *EntriesHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	entry := middleware.FromContext[db.Entry](r.Context())
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
	entry := middleware.FromContext[db.Entry](r.Context())
	return Layout(
		H1(Text("Edit " + entry.Name)),
	), nil
}

func (h *EntriesHandler) create(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, entriesPaths.List(), http.StatusFound)
}

func (h *EntriesHandler) update(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, entriesPaths.List(), http.StatusFound)
}

func (h *EntriesHandler) delete(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	entry := middleware.FromContext[db.Entry](r.Context())
	h.queries.DeleteEntryByID(r.Context(), db.DeleteEntryByIDParams{ID: entry.ID, UserID: userID})
	http.Redirect(w, r, "/", http.StatusFound)
}
