package probot

import (
	"bytes"
	ctx "context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
)

// PayloadInstallation represents the incoming installation part of the payload
type PayloadInstallation struct {
	Installation *github.Installation `json:"installation"`
}

type contextKey string

const probotContextKey contextKey = "probotContext"

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

			context := NewContext(app)

			// Validate the payload
			// Per the docs: https://docs.github.com/en/developers/webhooks-and-events/securing-your-webhooks#validating-payloads-from-github
			payloadBytes, err := github.ValidatePayload(r, []byte(app.Secret))
			if err != nil {
				log.Println(err)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Log the request headers
			log.Printf("signature validates: %s\n", r.Header.Get("X-Hub-Signature"))
			log.Printf("headers: %v\n", r.Header)

			// Get the installation from the payload
			payload := &PayloadInstallation{}
			json.Unmarshal(payloadBytes, payload)
			log.Printf("installation: %d\n", payload.Installation.GetID())
			log.Printf("received GitHub App ID %d\n", app.ID)

			// Parse the incoming request into an event
			context.Payload, err = github.ParseWebHook(github.WebHookType(r), payloadBytes)
			if err != nil {
				log.Println(err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			log.Printf("event type: %T\n", context.Payload)

			// Instantiate client
			installation := Installation{ID: payload.Installation.GetID()}
			context.GitHub, err = NewEnterpriseClient(app, installation)
			if err != nil {
				log.Println(err)
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}
			log.Printf("client %s instantiated for %s\n", context.GitHub.UserAgent, context.GitHub.BaseURL)

			// Reset the body for subsequent handlers to access
			r.Body = reset(r.Body, payloadBytes)

			// Call the next handler, with modified context.
			next.ServeHTTP(w, r.WithContext(ctx.WithValue(r.Context(), probotContextKey, context)))
		})
	}
}

// FromContext exposes probot features to request handlers
func FromContext(context ctx.Context) (*Context, bool) {
	probot, ok := context.Value(probotContextKey).(*Context)
	return probot, ok
}
