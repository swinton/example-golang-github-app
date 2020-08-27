package web

import (
	"encoding/json"
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
		// Read the incoming request as JSON
		body := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Log the request headers and body
		log.Printf("headers: %v\n", r.Header)
		log.Printf("body: %v\n", body)

		// Send output as application/json
		resp := HookResponse{
			Received: true,
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	return r
}
