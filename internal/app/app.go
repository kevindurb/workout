package app

import (
	"database/sql"
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	"github.com/kevindurb/planner/static"
)

type App struct {
	DB        *sql.DB
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
	return &App{
		DB:        conn,
		fp:        fp,
		queries:   q,
		home:      NewHomeHandler(),
		workouts:  NewWorkoutsHandler(q),
		exercises: NewExercisesHandler(q),
		entries:   NewEntriesHandler(q),
		sessions:  NewSessionsHandler(q, fp),
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Files))))

	handleAndStrip(mux, "/workouts", a.sessions.RequireAuth(a.workouts.Routes()))
	handleAndStrip(mux, "/exercises", a.sessions.RequireAuth(a.exercises.Routes()))
	handleAndStrip(mux, "/entries", a.sessions.RequireAuth(a.entries.Routes()))

	mux.Handle("/login", a.sessions.Routes())
	mux.Handle("/signup", a.sessions.Routes())
	mux.Handle("/", a.sessions.RequireAuth(a.home.Routes()))

	return a.sessions.sm.LoadAndSave(mux)
}

func handleAndStrip(mux *http.ServeMux, pattern string, h http.Handler) {
	mux.Handle(pattern+"/", http.StripPrefix(pattern, h))
	mux.Handle(pattern, http.RedirectHandler(pattern+"/", http.StatusMovedPermanently))
}
