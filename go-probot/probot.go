package probot

import (
	"context"
	"net/http"
)

type contextKey string

const probotKey contextKey = "probot"

func ProbotMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the next handler, with modified context.
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), probotKey, "alive")))
	})
}

func FromContext(ctx context.Context) (string, bool) {
	probot, ok := ctx.Value(probotKey).(string)
	return probot, ok
}
