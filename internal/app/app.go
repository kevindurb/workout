package app

import (
	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"

	"github.com/kevindurb/planner/static"
)

type App struct {
	DB        *sql.DB
	decoder   *schema.Decoder
	validator *validator.Validate
	home      *HomeHandler
}

func New(db *sql.DB) *App {
	return &App{
		DB:        db,
		decoder:   schema.NewDecoder(),
		validator: validator.New(),
		home:      NewHomeHandler(),
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Files))))

	mux.Handle("/", a.home.Routes())

	return mux
}
