package httpx

import (
	"net/http"
	"strconv"
)

func PathInt(r *http.Request, name string) int64 {
	id, _ := strconv.ParseInt(r.PathValue(name), 10, 64)
	return id
}
