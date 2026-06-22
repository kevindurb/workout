package tz

import (
	"net/http"
	"time"
)

func FromRequest(r *http.Request) *time.Location {
	loc := time.UTC
	if c, err := r.Cookie("tz"); err == nil {
		if l, err := time.LoadLocation(c.Value); err == nil {
			loc = l
		}
	}
	return loc
}
