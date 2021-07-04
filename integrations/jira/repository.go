package jira

import (
	"net/http"
	"time"

	"github.com/dominikbraun/timetrace/core"
)

const jiraV3BaseURL = "rest/api/3"

// Repository provides a set of methods for integrating JIRA as a third party
// provider into timetrace. It aims to fufill the command timetrace push jira,
// which will upload all of a users "unsyncronised" worklogs to JIRA. It is
// consumed through the interface in cli/push.go
type Repository struct {
	client      *http.Client
	authToken   string
	email       string
	jiraAddress string
	formatter   core.Formatter
}

// RepositoryConfig passes in all the config needed to connect to JIRA
type RepositoryConfig struct {
	AuthToken   string
	Email       string
	JIRAAddress string
}

type issue struct {
	Worklogs []worklog `json:"worklogs"`
}

type worklog struct {
	ID string `json:"id"`
}

// New returns an instantiated JIRA repository
func New(cfg RepositoryConfig) *Repository {
	return &Repository{
		authToken:   cfg.AuthToken,
		email:       cfg.Email,
		jiraAddress: cfg.JIRAAddress,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

// Name returns the name of the integration
func (r *Repository) Name() string {
	return "jira"
}
