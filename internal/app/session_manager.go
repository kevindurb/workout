package app

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type SessionManager struct {
	*scs.SessionManager
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		SessionManager: scs.New(),
	}
}

func (sm *SessionManager) SetUserID(ctx context.Context, id int64) {
	sm.Put(ctx, "user_id", id)
}

func (sm *SessionManager) UserID(ctx context.Context) int64 {
	return sm.GetInt64(ctx, "user_id")
}

func (sm *SessionManager) IsLoggedIn(ctx context.Context) bool {
	return sm.UserID(ctx) != 0
}

func (sm *SessionManager) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !sm.IsLoggedIn(r.Context()) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
