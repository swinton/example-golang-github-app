package probot

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
)

// PayloadInstallation represents the incoming installation part of the payload
type PayloadInstallation struct {
	Installation *github.Installation `json:"installation"`
}

type contextKey string

const probotKey contextKey = "probot"

// NewMiddleware returns a mux.MiddlewareFunc encapsulating Probot features
func NewMiddleware() mux.MiddlewareFunc {
	app := NewApp()
	log.Printf("loaded GitHub App ID %d\n", app.ID)

	// Middleware for Probot features
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Reset will return a new ReadCloser for the body that can be passed to subsequent handlers
			reset := func(old io.ReadCloser, b []byte) io.ReadCloser {
				old.Close()
				return ioutil.NopCloser(bytes.NewBuffer(b))
			}

			// Validate the payload
			// Per the docs: https://docs.github.com/en/developers/webhooks-and-events/securing-your-webhooks#validating-payloads-from-github
			payloadBytes, err := github.ValidatePayload(r, []byte(os.Getenv("GITHUB_APP_WEBHOOK_SECRET")))
			if err != nil {
				log.Println(err)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			log.Printf("signature validates: %s\n", r.Header.Get("X-Hub-Signature"))

			// Get the installation from the payload
			payload := &PayloadInstallation{}
			json.Unmarshal(payloadBytes, payload)
			log.Printf("installation: %d\n", payload.Installation.GetID())

			// Reset the body for subsequent handlers to access
			r.Body = reset(r.Body, payloadBytes)

			// Call the next handler, with modified context.
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), probotKey, app)))
		})
	}
}

// FromContext exposes probot features to request handlers
func FromContext(ctx context.Context) (*App, bool) {
	probot, ok := ctx.Value(probotKey).(*App)
	return probot, ok
}
