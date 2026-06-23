package app

import (
	"context"
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	. "github.com/kevindurb/planner/internal/html"

	"github.com/alexedwards/scs/v2"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type SessionsHandler struct {
	queries *db.Queries
	sm      *scs.SessionManager
}

func NewSessionsHandler(queries *db.Queries) *SessionsHandler {
	sm := scs.New()
	return &SessionsHandler{
		queries: queries,
		sm:      sm,
	}
}

func (h *SessionsHandler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /login", ghttp.Adapt(h.showLogin))
	mux.Handle("GET /signup", ghttp.Adapt(h.showSignup))

	return mux
}

func (h *SessionsHandler) showLogin(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("Login")),
	), nil
}

func (h *SessionsHandler) showSignup(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("Signup")),
	), nil
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
