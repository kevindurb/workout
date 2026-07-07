package signup

import (
	"log"
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"
	"github.com/kevindurb/planner/internal/session"
	"golang.org/x/crypto/bcrypt"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type signupBody struct {
	Email    string `form:"email,required" validate:"email"`
	Password string `form:"password,required"`
}

type Handler struct {
	q  *db.Queries
	sm *session.Manager
	fp *formparser.FormParser
}

func NewHandler(
	q *db.Queries,
	sm *session.Manager,
	fp *formparser.FormParser,
) *Handler {
	return &Handler{q, sm, fp}
}

func (h *Handler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("Signup")),
		Form(
			Method("POST"),
			Action("/signup"),
			Label(For("email"), Text("Email")),
			Input(Type("email"), ID("email"), Name("email"), Required()),
			Label(For("password"), Text("Password")),
			Input(Type("password"), ID("password"), Name("password"), Required()),
			Button(Type("submit"), Text("Signup")),
			A(Href("/login"), Text("Login")),
		),
	), nil
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	var data signupBody
	if err := h.fp.Parse(&data, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = h.q.CreateUser(r.Context(), db.CreateUserParams{
		Email: data.Email,
		Hash:  hash,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
