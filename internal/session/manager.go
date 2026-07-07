package session

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type Manager struct {
	*scs.SessionManager
}

func New() *Manager {
	return &Manager{
		SessionManager: scs.New(),
	}
}

func (m *Manager) SetUserID(ctx context.Context, id int64) {
	m.Put(ctx, "user_id", id)
}

func (m *Manager) UserID(ctx context.Context) int64 {
	return m.GetInt64(ctx, "user_id")
}

func (m *Manager) IsLoggedIn(ctx context.Context) bool {
	return m.UserID(ctx) != 0
}

func (m *Manager) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.IsLoggedIn(r.Context()) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
