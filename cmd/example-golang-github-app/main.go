package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/swinton/example-golang-github-app/gh"
	"github.com/swinton/example-golang-github-app/web"
)

func main() {
	// Read GitHub App credentials from environment
	privateKey, err := ioutil.ReadFile(os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to load GitHub App private key from file: %s", os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH")))
	}

	id, err := strconv.Atoi(os.Getenv("GITHUB_APP_ID"))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to load GitHub App: %s", os.Getenv("GITHUB_APP_ID")))
	}

	// Instantiate GitHub App
	app := gh.App{ID: id, Key: privateKey}

	// Webhook router
	router := web.HookRouter(app, "/")

	// Server, listening on port 8000
	port := 8000
	fmt.Fprintf(os.Stdout, "Server running at: http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), router))
}
