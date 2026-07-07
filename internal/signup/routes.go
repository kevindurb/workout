package signup

import (
	"github.com/go-chi/chi/v5"
	ghttp "maragu.dev/gomponents/http"
)

func (h *Handler) Routes(r chi.Router) {
	r.Get("/", ghttp.Adapt(h.show))
	r.Post("/", h.signup)
}
