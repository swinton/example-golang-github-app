package probotize

import (
	"context"
	"net/http"
)

type contextKey string

const probotKey contextKey = "probot"

func Probotize(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Modify the context
		// https://stackoverflow.com/a/49247940
		f(w, r.WithContext(context.WithValue(r.Context(), probotKey, "alive")))
	}
}

func FromContext(ctx context.Context) (string, bool) {
	probot, ok := ctx.Value(probotKey).(string)
	return probot, ok
}
