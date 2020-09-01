package gh

import (
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

// App encapsulates the fields needed to define a GitHub App
type App struct {
	GitHubEnterpriseBaseURL string
	ID                      int64
	Key                     []byte
}

// Installation encapsulates the fields needed to define an installation of a GitHub App
type Installation struct {
	ID int64
}

// NewEnterpriseClient instantiates a new GitHub Client using the App and Installation
func NewEnterpriseClient(app App, installation Installation) (*github.Client, error) {
	// Shared transport to reuse TCP connections.
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, app.ID, installation.ID, app.Key)
	if err != nil {
		return nil, err
	}

	itr.BaseURL = app.GitHubEnterpriseBaseURL
	client, err := github.NewEnterpriseClient(app.GitHubEnterpriseBaseURL, app.GitHubEnterpriseBaseURL, &http.Client{Transport: itr})
	if err != nil {
		return nil, err
	}

	// Overwrite User-Agent, for logging
	// See: https://developer.github.com/v3/#user-agent-required
	client.UserAgent = "swinton/example-golang-github-app"

	return client, nil
}
