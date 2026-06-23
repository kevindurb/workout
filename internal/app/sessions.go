package app

import (
	"context"
	"log"
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"
	"golang.org/x/crypto/bcrypt"

	"github.com/alexedwards/scs/v2"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type credsBody struct {
	Email    string `form:"email,required" validate:"email"`
	Password string `form:"password,required"`
}

type SessionsHandler struct {
	queries *db.Queries
	sm      *scs.SessionManager
	fp      *formparser.FormParser
}

func NewSessionsHandler(queries *db.Queries, fp *formparser.FormParser) *SessionsHandler {
	sm := scs.New()
	return &SessionsHandler{
		queries: queries,
		sm:      sm,
		fp:      fp,
	}
}

func (h *SessionsHandler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /login", ghttp.Adapt(h.showLogin))
	mux.Handle("POST /login", http.HandlerFunc(h.login))
	mux.Handle("GET /signup", ghttp.Adapt(h.showSignup))
	mux.Handle("POST /signup", http.HandlerFunc(h.signup))

	return mux
}

func (h *SessionsHandler) showLogin(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("Login")),
		Form(
			Method("POST"),
			Action("/login"),
			Label(For("email"), Text("Email")),
			Input(Type("email"), ID("email"), Name("email")),
			Label(For("password"), Text("Password")),
			Input(Type("password"), ID("password"), Name("password")),
			Button(Type("submit"), Text("Login")),
			A(Href("/signup"), Text("Signup")),
		),
	), nil
}

func (h *SessionsHandler) login(w http.ResponseWriter, r *http.Request) {
	var data credsBody
	if err := h.fp.Parse(&data, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUserByEmail(r.Context(), data.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = bcrypt.CompareHashAndPassword(user.Hash, []byte(data.Password)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.SetUserID(r.Context(), user.ID)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *SessionsHandler) showSignup(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("Signup")),
		Form(
			Method("POST"),
			Action("/signup"),
			Label(For("email"), Text("Email")),
			Input(Type("email"), ID("email"), Name("email")),
			Label(For("password"), Text("Password")),
			Input(Type("password"), ID("password"), Name("password")),
			Button(Type("submit"), Text("Signup")),
			A(Href("/login"), Text("Login")),
		),
	), nil
}

func (h *SessionsHandler) signup(w http.ResponseWriter, r *http.Request) {
	var data credsBody
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

	_, err = h.queries.CreateUser(r.Context(), db.CreateUserParams{
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

func (h *SessionsHandler) SetUserID(ctx context.Context, id int64) {
	h.sm.Put(ctx, "user_id", id)
}

func (h *SessionsHandler) UserID(ctx context.Context) int64 {
	return h.sm.GetInt64(ctx, "user_id")
}

func (h *SessionsHandler) IsLoggedIn(ctx context.Context) bool {
	return h.UserID(ctx) != 0
}

func (h *SessionsHandler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.IsLoggedIn(r.Context()) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
