package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"

	"github.com/swinton/example-golang-github-app/gh"
)

// HookResponse defines the shape of our HookRouter responses
type HookResponse struct {
	Received bool `json:"received"`
}

// HookRouter returns a new webhook router that can be plugged into an HTTP server to receive webhooks
func HookRouter(app gh.App, path string) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		// Read the incoming request
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		// Parse the incoming request into an event
		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			log.Println(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Log the request headers and body
		log.Printf("headers: %v\n", r.Header)
		log.Printf("event type: %T\n", event)

		// Handle the event
		switch e := event.(type) {
		case *github.IssuesEvent:
			log.Printf("installation id: %d\n", *e.Installation.ID)
			log.Printf("issue id: %d\n", *e.Issue.ID)
		default:
			log.Printf("Unknown event type: %s\n", github.WebHookType(r))
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Send response as application/json
		resp := HookResponse{
			Received: true,
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	return r
}
