package app

import (
	"net/http"
)

type Route struct {
	Pattern string
	Handler http.Handler
}

func registerAuthRoutes(mux *http.ServeMux, sm *SessionManager, routes []Route) {
	for _, route := range routes {
		mux.Handle(route.Pattern, sm.RequireAuth(route.Handler))
	}
}

func registerPublicRoutes(mux *http.ServeMux, routes []Route) {
	for _, route := range routes {
		mux.Handle(route.Pattern, route.Handler)
	}
}
