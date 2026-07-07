package entries

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kevindurb/planner/internal/db"
	"github.com/kevindurb/planner/internal/httpx"
	"github.com/kevindurb/planner/internal/middleware"

	ghttp "maragu.dev/gomponents/http"
)

func (h *Handler) Routes(r chi.Router) {
	r.Get("/", ghttp.Adapt(h.list))
	r.Get("/new", ghttp.Adapt(h.new))
	r.Post("/", h.create)

	r.Route("/{entry_id}", func(r chi.Router) {
		r.Use(middleware.EntityCtx(func(r *http.Request) (db.Entry, error) {
			return h.q.GetEntryByID(r.Context(), db.GetEntryByIDParams{
				ID:     httpx.PathInt(r, "entry_id"),
				UserID: h.sm.UserID(r.Context()),
			})
		}))
		r.Get("/", ghttp.Adapt(h.show))
		r.Get("/edit", ghttp.Adapt(h.edit))
		r.Post("/", h.update)
		r.Post("/delete", h.delete)
	})
}
