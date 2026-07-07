package app

import (
	"database/sql"
	"net/http"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/go-chi/chi/v5"
	"github.com/kevindurb/planner/internal/db"
	"github.com/kevindurb/planner/internal/entries"
	"github.com/kevindurb/planner/internal/exercises"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	"github.com/kevindurb/planner/internal/home"
	"github.com/kevindurb/planner/internal/login"
	"github.com/kevindurb/planner/internal/middleware"
	"github.com/kevindurb/planner/internal/session"
	"github.com/kevindurb/planner/internal/signup"
	"github.com/kevindurb/planner/internal/workouts"
	"github.com/kevindurb/planner/static"
)

type App struct {
	db               *sql.DB
	sm               *session.Manager
	fp               *formparser.FormParser
	q                *db.Queries
	homeHandler      *home.Handler
	workoutsHandler  *workouts.Handler
	exercisesHandler *exercises.Handler
	entriesHandler   *entries.Handler
	loginHandler     *login.Handler
	signupHandler    *signup.Handler
}

func New(conn *sql.DB) *App {
	q := db.New(conn)
	fp := formparser.New()
	sm := session.New()
	sm.Store = sqlite3store.New(conn)
	return &App{
		db:               conn,
		sm:               sm,
		fp:               fp,
		q:                q,
		homeHandler:      &home.Handler{q, sm, fp},
		workoutsHandler:  &workouts.Handler{q, sm, fp},
		exercisesHandler: &exercises.Handler{q, sm, fp},
		entriesHandler:   entries.NewHandler(q, sm, fp),
		loginHandler:     &login.Handler{q, sm, fp},
		signupHandler:    &signup.Handler{q, sm, fp},
	}
}

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.MethodOverride)
	r.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Files))))
	r.Route("/login", a.loginHandler.Routes)
	r.Route("/signup", a.signupHandler.Routes)
	r.With(a.sm.RequireAuth).Route("/workouts", a.workoutsHandler.Routes)
	r.With(a.sm.RequireAuth).Route("/exercises", a.exercisesHandler.Routes)
	r.With(a.sm.RequireAuth).Route("/entries", a.entriesHandler.Routes)
	r.With(a.sm.RequireAuth).Route("/", a.homeHandler.Routes)

	return a.sm.LoadAndSave(r)
}
