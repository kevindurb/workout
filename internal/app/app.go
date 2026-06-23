package app

import (
	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"

	"github.com/kevindurb/planner/internal/db"
	"github.com/kevindurb/planner/static"
)

type App struct {
	DB        *sql.DB
	queries   *db.Queries
	decoder   *schema.Decoder
	validator *validator.Validate
	home      *HomeHandler
	workouts  *WorkoutsHandler
	exercises *ExercisesHandler
	entries   *EntriesHandler
}

func New(conn *sql.DB) *App {
	q := db.New(conn)
	return &App{
		DB:        conn,
		queries:   q,
		decoder:   schema.NewDecoder(),
		validator: validator.New(),
		home:      NewHomeHandler(),
		workouts:  NewWorkoutsHandler(q),
		exercises: NewExercisesHandler(q),
		entries:   NewEntriesHandler(q),
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Files))))

	mux.Handle("/workouts/", http.StripPrefix("/workouts", a.workouts.Routes()))
	mux.Handle("/", a.home.Routes())

	return mux
}
