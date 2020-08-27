package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// HookResponse defines the shape of our HookRouter responses
type HookResponse struct {
	Received bool `json:"received"`
}

// HookRouter returns a new webhook router that can be plugged into an HTTP server to receive webhooks
func HookRouter(path string) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		// Read the incoming request
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}

		// Log the request body
		log.Printf("Received: %s\n", string(body))

		// Send output as application/json
		resp := HookResponse{
			Received: true,
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	return r
}
