package probot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

// App encapsulates the fields needed to define a GitHub App
type App struct {
	BaseURL string
	ID      int64
	Key     []byte
	Secret  string
	Client  *github.Client
}

// Installation encapsulates the fields needed to define an installation of a GitHub App
type Installation struct {
	ID int64
}

// NewApp instantiates a GitHub App from environment variables
func NewApp() *App {
	// Read GitHub App credentials from environment
	baseURL, exists := os.LookupEnv("GITHUB_BASE_URL")
	if !exists {
		log.Fatal("Unable to load GitHub Base URL from environment")
	}

	privateKey, err := ioutil.ReadFile(os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to load GitHub App private key from file: %s", os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH")))
	}

	id, err := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to load GitHub App: %s", os.Getenv("GITHUB_APP_ID")))
	}

	secret, exists := os.LookupEnv("GITHUB_APP_WEBHOOK_SECRET")
	if !exists {
		log.Fatal("Unable to load webhook secret from environment")
	}

	// Instantiate GitHub App
	app := &App{BaseURL: baseURL, ID: id, Key: privateKey, Secret: secret}

	return app
}

// NewEnterpriseClient instantiates a new GitHub Client using the App and Installation
func NewEnterpriseClient(app *App, installation Installation) (*github.Client, error) {
	// Shared transport to reuse TCP connections.
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, app.ID, installation.ID, app.Key)
	if err != nil {
		return nil, err
	}

	itr.BaseURL = app.BaseURL
	client, err := github.NewEnterpriseClient(app.BaseURL, app.BaseURL, &http.Client{Transport: itr})
	if err != nil {
		return nil, err
	}

	// Overwrite User-Agent, for logging
	// See: https://developer.github.com/v3/#user-agent-required
	client.UserAgent = "swinton/example-golang-github-app"

	return client, nil
}
