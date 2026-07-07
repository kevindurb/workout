package workouts

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kevindurb/planner/internal/db"
	ihttp "github.com/kevindurb/planner/internal/http"
	"github.com/kevindurb/planner/internal/middleware"

	ghttp "maragu.dev/gomponents/http"
)

func (h *Handler) Routes(r chi.Router) {
	r.Get("/", ghttp.Adapt(h.list))
	r.Get("/new", ghttp.Adapt(h.new))
	r.Post("/", h.create)

	r.Route("/{workout_id}", func(r chi.Router) {
		r.Use(middleware.EntityCtx(func(r *http.Request) (db.Workout, error) {
			return h.queries.GetWorkoutByID(r.Context(), db.GetWorkoutByIDParams{
				ID:     ihttp.PathInt(r, "workout_id"),
				UserID: h.sm.UserID(r.Context()),
			})
		}))
		r.Get("/", ghttp.Adapt(h.show))
		r.Get("/edit", ghttp.Adapt(h.edit))
		r.Post("/", h.update)
		r.Post("/delete", h.delete)
	})
}
