package middleware

import (
	"context"
	"net/http"
)

type ctxKey[T any] struct{}

func FromContext[T any](ctx context.Context) *T {
	if val, ok := ctx.Value(ctxKey[T]{}).(*T); ok {
		return val
	}
	return nil
}

func EntityCtx[T any](fetchFn func(r *http.Request) (T, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entity, err := fetchFn(r)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), ctxKey[T]{}, &entity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
