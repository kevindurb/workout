package app

import (
	"database/sql"
	"net/http"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/go-chi/chi/v5"
	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	"github.com/kevindurb/planner/static"
)

type App struct {
	db                       *sql.DB
	sm                       *SessionManager
	fp                       *formparser.FormParser
	q                        *db.Queries
	homeHandler              *HomeHandler
	workoutsHandler          *WorkoutsHandler
	exercisesHandler         *ExercisesHandler
	entriesHandler           *EntriesHandler
	workoutsExercisesHandler *WorkoutsExercisesHandler
	loginHandler             *LoginHandler
	signupHandler            *SignupHandler
}

func New(conn *sql.DB) *App {
	q := db.New(conn)
	fp := formparser.New()
	sm := NewSessionManager()
	sm.Store = sqlite3store.New(conn)
	return &App{
		db:                       conn,
		sm:                       sm,
		fp:                       fp,
		q:                        q,
		homeHandler:              &HomeHandler{q, sm, fp},
		workoutsHandler:          &WorkoutsHandler{q, sm, fp},
		exercisesHandler:         &ExercisesHandler{q, sm, fp},
		entriesHandler:           &EntriesHandler{q, sm},
		workoutsExercisesHandler: &WorkoutsExercisesHandler{q, sm, fp},
		loginHandler:             &LoginHandler{q, sm, fp},
		signupHandler:            &SignupHandler{q, sm, fp},
	}
}

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()
	r.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Files))))
	r.Route("/login", a.loginHandler.Route)
	r.Route("/signup", a.signupHandler.Route)
	r.With(a.sm.RequireAuth).Route("/workouts", a.workoutsHandler.Route)
	r.With(a.sm.RequireAuth).Route("/exercises", a.exercisesHandler.Route)
	r.With(a.sm.RequireAuth).Route("/entries", a.entriesHandler.Route)
	r.With(a.sm.RequireAuth).Route("/", a.homeHandler.Route)

	// a.workoutsExercisesHandler.Routes(r)

	return a.sm.LoadAndSave(r)
}
