package probot

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
)

type contextKey string

const probotKey contextKey = "probot"

// Middleware for Probot features
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reset will return a new ReadCloser for the body that can be passed to subsequent handlers
		reset := func(old io.ReadCloser, b []byte) io.ReadCloser {
			old.Close()
			return ioutil.NopCloser(bytes.NewBuffer(b))
		}

		// Validate the payload
		// Per the docs: https://docs.github.com/en/developers/webhooks-and-events/securing-your-webhooks#validating-payloads-from-github
		payload, err := github.ValidatePayload(r, []byte(os.Getenv("GITHUB_APP_WEBHOOK_SECRET")))
		if err != nil {
			log.Println(err)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		log.Printf("signature validates: %s\n", r.Header.Get("X-Hub-Signature"))

		// Reset the body for subsequent handlers to access
		r.Body = reset(r.Body, payload)

		// Call the next handler, with modified context.
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), probotKey, "alive")))
	})
}

func FromContext(ctx context.Context) (string, bool) {
	probot, ok := ctx.Value(probotKey).(string)
	return probot, ok
}
