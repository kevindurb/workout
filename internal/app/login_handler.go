package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"
	"github.com/kevindurb/planner/internal/session"
	"golang.org/x/crypto/bcrypt"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type loginBody struct {
	Email    string `form:"email,required" validate:"email"`
	Password string `form:"password,required"`
}

type LoginHandler struct {
	q  *db.Queries
	sm *session.Manager
	fp *formparser.FormParser
}

func (h *LoginHandler) Route(r chi.Router) {
	r.Get("/", ghttp.Adapt(h.show))
	r.Post("/", h.login)
}

func (h *LoginHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("Login")),
		Form(
			Method("POST"),
			Action("/login"),
			Label(For("email"), Text("Email")),
			Input(Type("email"), ID("email"), Name("email"), Required()),
			Label(For("password"), Text("Password")),
			Input(Type("password"), ID("password"), Name("password"), Required()),
			Button(Type("submit"), Text("Login")),
			A(Href("/signup"), Text("Signup")),
		),
	), nil
}

func (h *LoginHandler) login(w http.ResponseWriter, r *http.Request) {
	var data loginBody
	if err := h.fp.Parse(&data, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.q.GetUserByEmail(r.Context(), data.Email)
	if err != nil {
		log.Printf("Error getting user by email (%s): %v", data.Email, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = bcrypt.CompareHashAndPassword(user.Hash, []byte(data.Password)); err != nil {
		log.Printf("Error comparing password (%s) for user (%s): %v", data.Password, data.Email, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.sm.SetUserID(r.Context(), user.ID)

	http.Redirect(w, r, "/", http.StatusFound)
}
