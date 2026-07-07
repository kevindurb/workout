package middleware

import (
	"net/http"
	"strings"
)

func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err == nil {
				if m := r.PostFormValue("_method"); m != "" {
					m = strings.ToUpper(m)
					if m == http.MethodPut || m == http.MethodPatch || m == http.MethodDelete {
						r.Method = m
					}
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
