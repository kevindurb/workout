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
	sessionsHandler          *SessionsHandler
	workoutsExercisesHandler *WorkoutsExercisesHandler
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
		sessionsHandler:          &SessionsHandler{q, sm, fp},
		workoutsExercisesHandler: &WorkoutsExercisesHandler{q, sm, fp},
	}
}

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()
	r.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Files))))

	// a.sessionsHandler.Routes(r)
	// a.workoutsHandler.Routes(r)
	// a.exercisesHandler.Routes(r)
	// a.workoutsExercisesHandler.Routes(r)
	// a.entriesHandler.Routes(r)
	// a.homeHandler.Routes(r)

	return a.sm.LoadAndSave(r)
}
