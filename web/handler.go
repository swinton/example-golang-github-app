package web

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"

	"github.com/swinton/example-golang-github-app/gh"
)

var ctx = context.Background()

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
			log.Printf("issue owner: %s\n", *e.Repo.Owner.Login)
			log.Printf("issue repo: %s\n", *e.Repo.Name)
			log.Printf("issue id: %d\n", *e.Issue.ID)

			// Instantiate client
			installation := gh.Installation{ID: *e.Installation.ID}
			client, err := gh.NewEnterpriseClient(app, installation)
			if err != nil {
				log.Println(err)
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}
			log.Printf("client %s instantiated for %s\n", client.UserAgent, client.BaseURL)

			// Create a comment back on the issue
			// https://github.com/google/go-github/blob/d57a3a84ba041135efb6b7ad3991f827c93c306a/github/issues_comments.go#L101-L117
			newComment := &github.IssueComment{Body: github.String("## :wave: :earth_americas:\n\n![fellowshipoftheclaps](https://user-images.githubusercontent.com/27806/91333726-91c46f00-e793-11ea-9724-dc2e18ca28d0.gif)")}
			comment, _, err := client.Issues.CreateComment(ctx, *e.Repo.Owner.Login, *e.Repo.Name, int(*e.Issue.ID), newComment)
			if err != nil {
				log.Println(err)
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}
			log.Printf("comment created: %+v\n", comment)
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
