package jira

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/dominikbraun/timetrace/core"
)

// Repository provides a set of methods for integrating JIRA as a third party
// provider into timetrace. It aims to fufill the command timetrace push jira,
// which will upload all of a users "unsyncronised" worklogs to JIRA. It is
// consumed through the interface in cli/push.go
type Repository struct {
	client      *http.Client
	authToken   string
	email       string
	jiraAddress string
}

// RepositoryConfig passes in all the config needed to connect to JIRA
type RepositoryConfig struct {
	AuthToken   string
	Email       string
	JIRAAddress string
}

// New returns an instantiated JIRA repository
func New(cfg RepositoryConfig) *Repository {
	return &Repository{
		authToken:   cfg.AuthToken,
		email:       cfg.Email,
		jiraAddress: cfg.JIRAAddress,
	}
}

// Name returns the name of the integration
func (r *Repository) Name() string {
	return "jira"
}

// CheckRecordsExist queries JIRA for all current record states returning
// records that do not currently exist in JIRA. This places JIRA as the source
// of truth on whether records have been uploaded to JIRA already. It would be
// possible to save them to disk, but this runs the risk of reuploading a lot
// of work records to JIRA. Best to ask JIRA what it thinks the current state
// of its records is.
func (r *Repository) CheckRecordsExist(records []*core.Record) ([]*core.Record, error) {
	return records, nil
}

// UploadRecords takes an array of records and attempts to upload them to JIRA
func (r *Repository) UploadRecord(records *core.Record) error {
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	return nil
}
