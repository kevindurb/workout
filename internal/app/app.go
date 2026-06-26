package app

import (
	"database/sql"
	"net/http"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	"github.com/kevindurb/planner/static"
)

type App struct {
	DB        *sql.DB
	sm        *SessionManager
	fp        *formparser.FormParser
	queries   *db.Queries
	home      *HomeHandler
	workouts  *WorkoutsHandler
	exercises *ExercisesHandler
	entries   *EntriesHandler
	sessions  *SessionsHandler
}

func New(conn *sql.DB) *App {
	q := db.New(conn)
	fp := formparser.New()
	sm := NewSessionManager()
	sm.Store = sqlite3store.New(conn)
	return &App{
		DB:        conn,
		sm:        sm,
		fp:        fp,
		queries:   q,
		home:      &HomeHandler{},
		workouts:  &WorkoutsHandler{q, sm},
		exercises: &ExercisesHandler{q, sm, fp},
		entries:   &EntriesHandler{q, sm},
		sessions:  &SessionsHandler{q, sm, fp},
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Files))))

	handleAndStrip(mux, "/workouts", a.sm.RequireAuth(a.workouts.Routes()))
	handleAndStrip(mux, "/exercises", a.sm.RequireAuth(a.exercises.Routes()))
	handleAndStrip(mux, "/entries", a.sm.RequireAuth(a.entries.Routes()))

	mux.Handle("/login", a.sessions.Routes())
	mux.Handle("/signup", a.sessions.Routes())
	mux.Handle("/", a.sm.RequireAuth(a.home.Routes()))

	return a.sm.LoadAndSave(mux)
}

func handleAndStrip(mux *http.ServeMux, pattern string, h http.Handler) {
	mux.Handle(pattern+"/", http.StripPrefix(pattern, h))
	mux.Handle(pattern, http.RedirectHandler(pattern+"/", http.StatusMovedPermanently))
}
