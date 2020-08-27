package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/swinton/example-golang-github-app/web"
)

func main() {
	// Webhook router
	router := web.HookRouter("/")

	// Server, listening on port 8000
	port := 8000
	fmt.Fprintf(os.Stdout, "Server running at: http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), router))
}
